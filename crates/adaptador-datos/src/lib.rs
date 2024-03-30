mod catedra;
mod comentario;
mod docente;
mod materia;
mod sql;

use http_cache::{CACacheManager, CacheMode, HttpCache};
use http_cache_reqwest::Cache;
use reqwest::Client;
use reqwest_middleware::{ClientBuilder, ClientWithMiddleware};
use std::{
    collections::{HashMap, HashSet},
    sync::Arc,
};

use sql::BulkInsertTuples;

pub async fn generar_query() -> anyhow::Result<String> {
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

    for mat in materias.into_iter().take(1) {
        bulk_inserts.materias.push(mat.sql());

        let cliente_http = Arc::clone(&cliente_http);
        handles.push(tokio::spawn(async move { mat.scape(&cliente_http).await }));
    }

    let mut codigos_docentes = HashMap::new();

    for task in handles {
        if let Ok(mut materia) = task.await.unwrap() {
            codigos_docentes.extend(std::mem::take(&mut materia.codigos_docentes));
            bulk_inserts.extend(materia);
        }
    }

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
