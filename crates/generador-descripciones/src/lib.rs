mod inference_api;

use std::{collections::HashMap, sync::Arc};

use reqwest::Client;
use sqlx::{types::Uuid, FromRow, PgPool};
use tokio::sync::Semaphore;
use tracing::Instrument;

const MAX_SOLICITUDES_CONCURRENTES: usize = 10;
const PROPORCION_COMENTARIOS_ACTUALIZACION: usize = 2;
const MIN_COMENTARIOS_ACTUALIZACION: usize = 3;

#[derive(FromRow)]
struct Comentario {
    codigo_docente: Uuid,
    contenido: String,
}

pub async fn query_actualizacion(
    conexion_db: &PgPool,
    inference_api_key: String,
) -> anyhow::Result<Option<String>> {
    let cliente_http = Client::new();
    let inference_api_key = Arc::new(inference_api_key);

    tracing::info!("conexion establecida con la base de datos");

    let comentarios: Vec<Comentario> = sqlx::query_as(const_format::formatcp!(
        r#"
SELECT com.codigo_docente, com.contenido
FROM comentario com
WHERE com.codigo_docente IN (
  SELECT doc.codigo
  FROM docente doc
  INNER JOIN comentario com
  ON com.codigo_docente = doc.codigo
  GROUP BY doc.codigo
  HAVING COUNT(com) > (doc.comentarios_ultima_descripcion * {})
  AND COUNT(com) > {}
);
"#,
        PROPORCION_COMENTARIOS_ACTUALIZACION,
        MIN_COMENTARIOS_ACTUALIZACION
    ))
    .fetch_all(conexion_db)
    .await?;

    tracing::info!("comentarios obtenidos de la base de datos");

    let mut comentarios_por_docente: HashMap<Uuid, Vec<String>> = HashMap::new();

    for comentario in comentarios {
        let comentarios_de_docente = comentarios_por_docente
            .entry(comentario.codigo_docente)
            .or_default();

        comentarios_de_docente.push(comentario.contenido);
    }

    let cantidad_docentes = comentarios_por_docente.len();

    tracing::info!("encontrados {cantidad_docentes} docentes que requiren actualizacion");

    let semaphore = Arc::new(Semaphore::new(MAX_SOLICITUDES_CONCURRENTES));
    let mut handles = Vec::with_capacity(cantidad_docentes);

    for (codigo_docente, comentarios) in comentarios_por_docente {
        let cliente_http = Client::clone(&cliente_http);
        let semaphore = Arc::clone(&semaphore);
        let inference_api_key = Arc::clone(&inference_api_key);

        let span = tracing::debug_span!("docente", codigo = codigo_docente.to_string());

        handles.push(tokio::spawn(
            async move {
                let _permit = semaphore.acquire().await?;
                generar_tupla_values(
                    cliente_http,
                    &semaphore,
                    codigo_docente,
                    &comentarios,
                    &inference_api_key,
                )
                .await
            }
            .instrument(span),
        ));
    }

    let mut values_a_actualizar = Vec::with_capacity(handles.len());

    for handle in handles {
        if let Ok(tupla_values) = handle.await.unwrap() {
            values_a_actualizar.push(tupla_values);
        }
    }

    tracing::info!(
        "query de actualizacion generada para {} de {} docentes",
        values_a_actualizar.len(),
        cantidad_docentes
    );

    let query = if !values_a_actualizar.is_empty() {
        Some(format!(
            r#"
UPDATE docente AS doc
SET descripcion = val.descripcion,
    comentarios_ultima_descripcion = val.comentarios_ultima_descripcion
FROM (VALUES
    {})
AS val(codigo_docente, descripcion, comentarios_ultima_descripcion)
WHERE doc.codigo::text = val.codigo_docente;
"#,
            values_a_actualizar.join(",\n    ")
        ))
    } else {
        None
    };

    Ok(query)
}

async fn generar_tupla_values(
    cliente_http: Client,
    semaphore: &Semaphore,
    codigo_docente: Uuid,
    comentarios: &[String],
    inference_api_key: &str,
) -> anyhow::Result<String> {
    let cantidad_comentarios_actual = comentarios.len();
    let descripcion =
        inference_api::generar_descripcion(cliente_http, comentarios, inference_api_key)
            .await
            .map_err(|err| {
                semaphore.close();
                err
            })?;

    let tupla_values = format!(
        "('{}', '{}', {})",
        codigo_docente,
        descripcion.replace('\'', "''"),
        cantidad_comentarios_actual
    );

    Ok(tupla_values)
}
