mod dolly;

use std::collections::HashMap;

use dolly::{Calificacion, Catedra, ComentariosDocentePorCuatri, Materia};
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

    // Primero descargamos los datos de entrada "crudos". A partir e las materias luego se van a
    // descargar las catedras y docentes correspondientes a cada materia.
    //
    let materias = Materia::fetch_all(&client).await?;
    let comentarios = ComentariosDocentePorCuatri::fetch_all(&client).await?;

    // Ante la falta de un id unico por docente, el codigo de la materia que imparte y el nombre
    // del docente sirven como identificador unico. Sin embargo yo quiero trabajar con un verdadero
    // id, entonces relaciono estos campos con un UUID generado en runtime.
    //
    let mut docente_to_uuid = HashMap::new();

    let mut output_buffer: Vec<String> = vec![
        Materia::CREATE_TABLE.into(),
        Catedra::CREATE_TABLE.into(),
        Calificacion::CREATE_TABLE.into(),
        ComentariosDocentePorCuatri::CREATE_TABLE.into(),
    ];

    for materia in materias {
        let mut query_block_buffer = vec![materia.sql()];

        let catedras = match Catedra::fetch_for_materia(&client, materia.codigo).await {
            Ok(catedras) => catedras.map(|c| (Uuid::new_v4(), c)).collect::<Vec<_>>(),
            Err(err) => {
                tracing::error!("error descargando catedras de materia {}", materia.codigo);
                tracing::debug!("descripcion error: {err}");
                continue;
            }
        };

        for (codigo_catedra, catedra) in catedras {
            query_block_buffer.push(Catedra::sql(codigo_catedra, materia.codigo));

            // TODO: Me hace falta establecer una relacion entre docente y catedra. Ya tengo una
            // relacion entre el docente y la materia.
            //
            for (nombre_docente, calificacion) in catedra.docentes {
                docente_to_uuid
                    .entry((materia.codigo, nombre_docente.to_owned()))
                    .or_insert_with(|| {
                        let codigo_docente = Uuid::new_v4();
                        query_block_buffer.push(calificacion.sql(&codigo_docente, nombre_docente));
                        codigo_docente
                    });
            }
        }

        output_buffer.push(query_block_buffer.join("\n"));
    }

    for ((codigo_materia, nombre_docente), comentarios) in comentarios.into_iter() {
        if let Some(codigo_docente) = docente_to_uuid.get(&(codigo_materia, nombre_docente)) {
            output_buffer.push(comentarios.sql(codigo_docente));
        }
    }

    Ok(output_buffer.join("\n\n"))
}
