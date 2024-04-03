mod catedra;
mod comentario;
mod docente;
mod materia;
mod sql;

use http_cache::{CACacheManager, CacheMode, HttpCache};
use http_cache_reqwest::Cache;
use reqwest::Client;
use reqwest_middleware::{ClientBuilder, ClientWithMiddleware};
use sqlx::{FromRow, PgPool};
use std::{
    collections::{HashMap, HashSet},
    sync::Arc,
};
use tokio::sync::Semaphore;
use uuid::Uuid;

use sql::BulkInsertTuples;

const MAX_SOLICITUDES_CONCURRENTES: usize = 5;

pub async fn init_query() -> anyhow::Result<String> {
    generar_query(HashMap::new()).await
}

pub async fn update_query(db: &PgPool) -> anyhow::Result<String> {
    let mut codigos_docentes: HashMap<i16, HashMap<String, Uuid>> = HashMap::new();

    #[derive(FromRow)]
    struct Docente {
        codigo: Uuid,
        nombre: String,
        codigo_materia: i16,
    }

    let docentes_existentes =
        sqlx::query_as::<_, Docente>("SELECT codigo, nombre, codigo_materia FROM docente;")
            .fetch_all(db)
            .await
            .unwrap();

    for doc in docentes_existentes {
        let docentes_materia = codigos_docentes.entry(doc.codigo_materia).or_default();
        docentes_materia.insert(doc.nombre, doc.codigo);
    }

    generar_query(codigos_docentes).await
}

async fn generar_query(
    mut codigos_docentes: HashMap<i16, HashMap<String, Uuid>>,
) -> anyhow::Result<String> {
    let cliente_http = Arc::new(crear_cliente_http());

    let materias = materia::descargar_todas(&cliente_http).await?;
    let comentarios = comentario::descargar_todos(&cliente_http).await?;

    let mut bulk_inserts = BulkInsertTuples {
        materias: Vec::with_capacity(materias.len()),
        comentarios: Vec::with_capacity(comentarios.len()),
        cuatrimestres: comentarios
            .keys()
            .map(|c| c.nombre_cuatrimestre.as_str())
            .collect::<HashSet<_>>()
            .into_iter()
            .map(comentario::sql_cuatrimestre)
            .collect(),
        ..Default::default()
    };

    let mut handles = Vec::with_capacity(materias.len());
    let semaphore = Arc::new(Semaphore::new(MAX_SOLICITUDES_CONCURRENTES));

    for mat in materias {
        bulk_inserts.materias.push(mat.sql());

        let cliente_http = Arc::clone(&cliente_http);
        let semaphore = Arc::clone(&semaphore);
        let codigos_docentes = codigos_docentes.remove(&mat.codigo).unwrap_or_default();

        handles.push(tokio::spawn(async move {
            let _permit = semaphore.acquire().await;
            mat.scrape(&cliente_http, codigos_docentes).await
        }));
    }

    let mut codigos_docentes = HashMap::new();

    for task in handles {
        if let Ok(mut materia) = task.await.unwrap() {
            codigos_docentes.extend(std::mem::take(&mut materia.codigos_docentes));
            bulk_inserts.extend(materia);
        }
    }

    tracing::info!("adaptando datos de materias descargadas");
    tracing::info!("adaptando comentarios de docentes");

    for (meta, comentarios) in comentarios {
        if let Some(codigo_docente) =
            codigos_docentes.get(&(meta.nombre_docente, meta.codigo_materia))
        {
            bulk_inserts.comentarios.extend(
                comentarios.iter().map(|c| {
                    comentario::sql_comentario(c, codigo_docente, &meta.nombre_cuatrimestre)
                }),
            );
        }
    }

    tracing::info!("generando query sql");

    Ok([
        &String::from_utf8_lossy(include_bytes!("../sql/schema.sql")).to_string(),
        "BEGIN;",
        &bulk_inserts.sql(),
        "COMMIT;",
    ]
    .join("\n\n"))
}

fn crear_cliente_http() -> ClientWithMiddleware {
    let cache_mode = if cfg!(debug_assertions) {
        tracing::debug!("forzando utilización de cache en requests");
        CacheMode::ForceCache
    } else {
        CacheMode::Default
    };

    ClientBuilder::new(Client::new())
        .with(Cache(HttpCache {
            mode: cache_mode,
            manager: CACacheManager::default(),
            options: None,
        }))
        .build()
}
