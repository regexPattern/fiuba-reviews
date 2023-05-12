mod catedras;
mod comentarios;
mod materias;
mod sql;

use std::collections::HashMap;

use comentarios::Cuatrimestre;
use http_cache_reqwest::{CACacheManager, Cache, CacheMode, HttpCache};
use materias::Materia;
use reqwest::Client;
use reqwest_middleware::ClientBuilder;
use uuid::Uuid;

pub async fn generar_sql() -> anyhow::Result<String> {
    let http = ClientBuilder::new(Client::new())
        .with(Cache(HttpCache {
            mode: CacheMode::Default,
            manager: CACacheManager::default(),
            options: None,
        }))
        .build();

    let mut queries: Vec<String> = vec![
        materias::CREACION_TABLA.into(),
        catedras::CREACION_TABLA_CATEDRAS.into(),
        catedras::CREACION_TABLA_DOCENTES.into(),
        comentarios::CREACION_TABLA.into(),
        catedras::CREACION_TABLA_CATEDRA_DOCENTE.into(),
    ];

    let materias = Materia::descargar(&http).await?;
    let mut codigos_docentes = HashMap::new();

    for materia in materias {
        queries.push(materia.query_sql());

        let catedras = match materia.catedras(&http).await {
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

    let comentarios = Cuatrimestre::descargar_comentarios(&http).await?;

    for (cuatrimestre, entradas) in comentarios {
        if let Some(codigo_docente) = codigos_docentes.get(&(
            cuatrimestre.codigo_materia,
            cuatrimestre.nombre_docente.clone(),
        )) {
            queries.push(cuatrimestre.sql(codigo_docente, &entradas));
        }
    }

    Ok(queries.join(""))
}
