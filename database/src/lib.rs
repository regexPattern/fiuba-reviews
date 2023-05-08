mod catedra;
mod comentario;
mod materia;

use std::collections::HashMap;

use catedra::{docente, Catedra};
use comentario::Cuatrimestre;
use http_cache_reqwest::{CACacheManager, Cache, CacheMode, HttpCache};
use materia::Materia;
use reqwest::Client;
use reqwest_middleware::ClientBuilder;
use uuid::Uuid;

// use crate::comentario::Cuatrimestre;

pub async fn run() -> anyhow::Result<String> {
    let cliente = ClientBuilder::new(Client::new())
        .with(Cache(HttpCache {
            mode: CacheMode::ForceCache,
            manager: CACacheManager::default(),
            options: None,
        }))
        .build();

    let materias = Materia::descargar(&cliente).await?;
    let comentarios = Cuatrimestre::descargar_comentarios(&cliente).await?;

    let mut codigos_docentes = HashMap::new();

    let mut buffer_output: Vec<String> = vec![
        materia::TABLA.into(),
        catedra::TABLA.into(),
        docente::TABLA.into(),
        comentario::TABLA.into(),
        catedra::TABLA_RELACION_CATEDRA_DOCENTE.into(),
    ];

    for materia in materias.into_iter() {
        let mut bloque_query_buffer = vec![materia.sql()];

        let catedras = match Catedra::descargar_para_materia(&cliente, materia.codigo).await {
            Ok(catedras) => catedras,
            Err(err) => {
                tracing::error!("error descargando catedras de materia {}", materia.codigo);
                tracing::debug!("descripcion error: {err}");
                continue;
            }
        };

        for catedra in catedras {
            bloque_query_buffer.push(catedra.sql(materia.codigo));

            for (nombre_docente, calificacion) in &catedra.docentes {
                let codigo_docente = codigos_docentes
                    .entry((materia.codigo, nombre_docente.to_owned()))
                    .or_insert_with(|| {
                        let codigo_docente = Uuid::new_v4();
                        bloque_query_buffer.push(calificacion.sql(&nombre_docente, codigo_docente));
                        codigo_docente
                    });

                bloque_query_buffer.push(catedra.relacionar_docente_sql(&codigo_docente));
            }
        }

        buffer_output.push(bloque_query_buffer.join("\n"));
    }

    for (cuatrimestre, comentarios) in comentarios.into_iter() {
        if let Some(codigo_docente) = codigos_docentes.get(&(
            cuatrimestre.codigo_materia,
            cuatrimestre.nombre_docente.clone(),
        )) {
            buffer_output.push(cuatrimestre.sql(codigo_docente, &comentarios));
        }
    }

    Ok(buffer_output.join("\n\n"))
}
