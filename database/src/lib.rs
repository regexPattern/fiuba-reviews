mod dolly;

use std::collections::HashMap;

use dolly::sql::create_tables;
use dolly::{Catedra, ComentariosDocentePorCuatri, Materia};
use http_cache_reqwest::{CACacheManager, Cache, CacheMode, HttpCache};
use reqwest::Client;
use reqwest_middleware::ClientBuilder;
use uuid::Uuid;

pub async fn run() -> anyhow::Result<String> {
    let client = ClientBuilder::new(Client::new())
        .with(Cache(HttpCache {
            mode: CacheMode::ForceCache,
            manager: CACacheManager::default(),
            options: None,
        }))
        .build();

    let materias = Materia::fetch_all(&client).await?;
    let comentarios = ComentariosDocentePorCuatri::fetch_all(&client).await?;

    let mut docente_to_uuid = HashMap::new();

    let mut output_buffer: Vec<String> = vec![
        create_tables::MATERIAS.into(),
        create_tables::CATEDRAS.into(),
        create_tables::DOCENTES.into(),
        create_tables::COMENTARIOS.into(),
        create_tables::CATEDRA_DOCENTE.into(),
    ];

    for materia in materias {
        let mut query_buffer = vec![materia.insert_query()];

        let catedras = match Catedra::fetch_for_materia(&client, materia.codigo).await {
            Ok(catedras) => catedras.map(|c| (Uuid::new_v4(), c)),
            Err(err) => {
                tracing::error!("error descargando catedras de materia {}", materia.codigo);
                tracing::debug!("descripcion error: {err}");
                continue;
            }
        };

        for (codigo_catedra, catedra) in catedras {
            query_buffer.push(Catedra::insert_query(codigo_catedra, materia.codigo));

            for (nombre_docente, calificacion) in catedra.docentes {
                let codigo_docente = docente_to_uuid
                    .entry((materia.codigo, nombre_docente.to_owned()))
                    .or_insert_with(|| {
                        let codigo_docente = Uuid::new_v4();
                        query_buffer
                            .push(calificacion.insert_query(&codigo_docente, nombre_docente));
                        codigo_docente
                    });

                query_buffer.push(dolly::sql::catedra_docente_rel_query(
                    &codigo_catedra,
                    &codigo_docente,
                ));
            }
        }

        output_buffer.push(query_buffer.join("\n"));
    }

    for ((codigo_materia, nombre_docente), comentarios) in comentarios.into_iter() {
        if let Some(codigo_docente) = docente_to_uuid.get(&(codigo_materia, nombre_docente)) {
            output_buffer.push(comentarios.insert_query(codigo_docente));
        }
    }

    Ok(output_buffer.join("\n\n"))
}
