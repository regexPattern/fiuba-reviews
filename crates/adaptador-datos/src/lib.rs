mod catedras;
mod comentarios;
mod materias;
mod sql;

use std::collections::{HashMap, HashSet};

use comentarios::{Comentario, Cuatrimestre};
use http_cache_reqwest::{CACacheManager, Cache, CacheMode, HttpCache};
use materias::Materia;
use reqwest::Client;
use reqwest_middleware::ClientBuilder;
use uuid::Uuid;

pub async fn indexar_dolly() -> anyhow::Result<String> {
    let cliente_http = ClientBuilder::new(Client::new())
        .with(Cache(HttpCache {
            mode: CacheMode::ForceCache,
            manager: CACacheManager::default(),
            options: None,
        }))
        .build();

    let mut queries: Vec<String> = vec![
        String::from_utf8_lossy(include_bytes!("../sql/tablas.sql")).into(),
        format!("BEGIN;"),
    ];

    let materias = Materia::descargar(&cliente_http).await?;
    let comentarios = Cuatrimestre::descargar(&cliente_http).await?;

    let mut codigos_docentes = HashMap::new();

    for materia in materias {
        queries.push(materia.query_sql());

        let catedras = match materia.catedras(&cliente_http).await {
            Ok(catedras) => catedras,
            Err(err) => {
                tracing::error!("error descargando catedras de materia {}", materia.codigo);
                tracing::debug!("descripcion error: {err}");
                continue;
            }
        };

        for catedra in catedras {
            queries.push(catedra.query_sql(materia.codigo));

            for (nombre_docente, calificacion) in &catedra.docentes {
                let codigo_docente = codigos_docentes
                    .entry((materia.codigo, nombre_docente.clone()))
                    .or_insert_with(|| {
                        let codigo_docente = Uuid::new_v4();
                        queries.push(calificacion.query_sql(nombre_docente, codigo_docente));
                        codigo_docente
                    });

                queries.push(catedra.relacion_con_docente_query_sql(codigo_docente));
            }
        }
    }

    let nombres_cuatrimestres: HashSet<&str> =
        comentarios.keys().map(|c| c.nombre.as_str()).collect();

    for nombre in nombres_cuatrimestres {
        queries.push(Cuatrimestre::sql(nombre));
    }

    for (cuatrimestre, comentarios) in &comentarios {
        if let Some(codigo_docente) = codigos_docentes.get(&(
            cuatrimestre.codigo_materia,
            cuatrimestre.nombre_docente.clone(),
        )) {
            queries.push(Comentario::query_sql(
                cuatrimestre,
                codigo_docente,
                comentarios,
            ));
        }
    }

    queries.push("COMMIT;".into());

    Ok(queries.join(""))
}
