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
use uuid::Uuid;

use materia::Materia;
use sql::InsertTuplesBuffer;

pub async fn init_query() -> anyhow::Result<String> {
    let cliente_http = Arc::new(crear_cliente_http());

    let materias = Materia::descargar_todas(&cliente_http).await?;
    let comentarios = comentario::descargar_todos(&cliente_http).await?;

    let mut data = InsertTuplesBuffer {
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

    let mut handles: Vec<_> = Vec::with_capacity(materias.len());

    for mat in materias {
        data.materias.push(mat.sql());

        let cliente_http = Arc::clone(&cliente_http);
        handles.push(tokio::spawn(
            async move { mat.indexar(&cliente_http).await },
        ));
    }

    let mut codigos_docentes: HashMap<(i16, String), Uuid> = HashMap::new();

    for task in handles {
        if let Ok(mut materia) = task.await.unwrap() {
            codigos_docentes.extend(std::mem::take(&mut materia.codigos_docentes));
            data.extend(materia);
        }
    }

    for (meta, comentarios) in comentarios {
        if let Some(codigo_docente) =
            codigos_docentes.get(&(meta.codigo_materia, meta.nombre_docente))
        {
            data.comentarios.extend(
                comentarios.iter().map(|c| {
                    comentario::sql_comentario(c, codigo_docente, &meta.nombre_cuatrimestre)
                }),
            );
        }
    }

    Ok(["BEGIN;\n", &data.sql(), "\nCOMMIT;\n"].join("\n"))
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
