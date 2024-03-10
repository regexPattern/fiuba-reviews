mod catedras;
mod comentarios;
mod docentes;
mod materias;
mod sql;

use std::collections::HashSet;

use comentarios::Comentario;
use http_cache_reqwest::{CACacheManager, Cache, CacheMode, HttpCache};
use materias::Materia;
use reqwest::Client;
use reqwest_middleware::ClientBuilder;

pub async fn indexar_dolly() -> anyhow::Result<String> {
    let cliente_http = ClientBuilder::new(Client::new())
        .with(Cache(HttpCache {
            mode: CacheMode::Default,
            manager: CACacheManager::default(),
            options: None,
        }))
        .build();

    let mut output_buffer: Vec<String> = vec![
        String::from_utf8_lossy(include_bytes!("../sql/schema.sql")).into(),
        format!("BEGIN;"),
    ];

    let materias = Materia::descargar_todas(&cliente_http).await?;
    let comentarios_docentes_cuatrimestres = Comentario::descargar_todos(&cliente_http).await?;

    let mut codigos_docentes = HashSet::new();

    for materia in materias.into_iter().take(1) {
        output_buffer.push(materia.query_sql());

        let catedras = match materia.descargar_catedras(&cliente_http).await {
            Ok(catedras) => catedras,
            Err(err) => {
                tracing::error!("error descargando catedras de materia {}", materia.codigo);
                tracing::debug!("descripcion error: {err}");
                continue;
            }
        };

        for catedra in catedras {
            output_buffer.push(catedra.query_sql(materia.codigo));

            for docente in &catedra.docentes {
                if !codigos_docentes.contains(&docente.codigo) {
                    output_buffer.push(docente.query_sql());
                    codigos_docentes.insert(docente.codigo.clone());
                }

                output_buffer.push(catedra.relacion_con_docente_query_sql(&docente));
            }
        }
    }

    let nombre_cuatrimestres: HashSet<&str> = comentarios_docentes_cuatrimestres
        .keys()
        .map(|md| md.nombre_cuatrimestre.as_str())
        .collect();

    for nombre in nombre_cuatrimestres {
        output_buffer.push(Comentario::cuatrimestre_query_sql(nombre));
    }

    for (metadata, comentarios) in &comentarios_docentes_cuatrimestres {
        if codigos_docentes.contains(&metadata.codigo_docente) {
            output_buffer.push(Comentario::query_sql(metadata, comentarios));
        }
    }

    output_buffer.push("COMMIT;".into());

    Ok(output_buffer.join(""))
}
