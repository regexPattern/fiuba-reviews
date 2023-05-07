mod dolly;

use std::collections::{HashMap, HashSet};

use dolly::sql::create_tables;
use dolly::{Catedra, Materia};
use http_cache_reqwest::{CACacheManager, Cache, CacheMode, HttpCache};
use reqwest::Client;
use reqwest_middleware::ClientBuilder;
use uuid::Uuid;

use crate::dolly::Comentarios;

pub async fn run() -> anyhow::Result<String> {
    let client = ClientBuilder::new(Client::new())
        .with(Cache(HttpCache {
            mode: CacheMode::ForceCache,
            manager: CACacheManager::default(),
            options: None,
        }))
        .build();

    let materias = Materia::fetch_all(&client).await?;
    let comentarios = Comentarios::fetch_all(&client).await?;

    let mut docente_to_uuid = HashMap::new();

    let mut output_buffer: Vec<String> = vec![
        create_tables::MATERIAS.into(),
        create_tables::CATEDRAS.into(),
        create_tables::DOCENTES.into(),
        create_tables::COMENTARIOS.into(),
        create_tables::CATEDRA_DOCENTE.into(),
    ];

    for materia in materias.take(1) {
        let mut query_buffer = vec![materia.insert_query()];

        let catedras = match Catedra::fetch_for_materia(&client, materia.codigo).await {
            Ok(catedras) => catedras.map(|c| (Uuid::new_v4(), c)),
            Err(err) => {
                tracing::error!("error descargando catedras de materia {}", materia.codigo);
                tracing::debug!("descripcion error: {err}");
                continue;
            }
        };

        let mut catedras_guardadas = HashSet::new();

        for (codigo_catedra, catedra) in catedras {
            let mut nombres_docentes_catedra = catedra
                .docentes
                .keys()
                .map(|nombre| nombre.to_owned())
                .collect::<Vec<_>>();

            // TODO: Esta no s la forma mas eficiente de hacer esto claramente, ni en tiempo ni en
            // memoria, pero como se me fueron ocurriendo estas ideas sobre la marcha, entonces por
            // el momento se queda asi.
            //
            nombres_docentes_catedra.sort();
            let nombre_catedra = nombres_docentes_catedra.join("-");
            if catedras_guardadas.contains(&nombre_catedra) {
                continue;
            } else {
                catedras_guardadas.insert(nombre_catedra.clone());
            }

            query_buffer.push(Catedra::insert_query(
                codigo_catedra,
                nombre_catedra,
                materia.codigo,
            ));

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

    let comentarios_vera = comentarios
        .into_iter()
        .filter(|(c, _)| c.codigo_materia == 7507 && c.nombre_docente == "Justo Narcizo")
        .collect::<Vec<_>>();

    dbg!(&comentarios_vera);
    dbg!(comentarios_vera.len());

    Ok(output_buffer.join("\n\n"))
}
