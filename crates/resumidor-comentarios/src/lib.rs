pub mod gpt;

use std::{collections::HashMap, sync::Arc};

use gpt::Modelo;
use reqwest::Client;
use sqlx::{types::Uuid, FromRow, PgPool};
use tokio::sync::Semaphore;
use tracing::Instrument;

const MAX_SOLICITUDES_CONCURRENTES: usize = 5;
const MIN_COMENTARIOS_ACTUALIZACION: usize = 3;
const PROPORCION_COMENTARIOS_ACTUALIZACION: usize = 2 / 1;

#[derive(FromRow)]
struct Comentario {
    codigo_docente: Uuid,
    contenido: String,
}

pub async fn query_actualizacion<M>(
    conexion_db: &PgPool,
    modelo: M,
) -> anyhow::Result<Option<String>>
where
    M: Modelo + Send + Sync + 'static,
{
    let cliente_http = Client::new();
    let modelo = Arc::new(modelo);

    let comentarios_por_docente = comentarios_por_docente(conexion_db).await?;
    let cantidad_docentes = comentarios_por_docente.len();

    let semaphore = Arc::new(Semaphore::new(MAX_SOLICITUDES_CONCURRENTES));
    let mut handles = Vec::with_capacity(cantidad_docentes);

    for (codigo_docente, comentarios) in comentarios_por_docente {
        let cliente_http = Client::clone(&cliente_http);
        let modelo = Arc::clone(&modelo);

        let semaphore = Arc::clone(&semaphore);
        let span = tracing::debug_span!("docente", codigo = codigo_docente.to_string());

        handles.push(tokio::spawn(
            async move {
                let _permit = semaphore.acquire().await?;
                tupla_values(
                    cliente_http,
                    &semaphore,
                    codigo_docente,
                    &comentarios,
                    modelo,
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

async fn comentarios_por_docente(
    conexion_db: &PgPool,
) -> anyhow::Result<HashMap<Uuid, Vec<String>>> {
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

    Ok(comentarios_por_docente)
}

async fn tupla_values<M>(
    cliente_http: Client,
    semaphore: &Semaphore,
    codigo_docente: Uuid,
    comentarios: &[String],
    modelo: Arc<M>,
) -> anyhow::Result<String>
where
    M: Modelo + Send + Sync + 'static,
{
    let cantidad_comentarios_actual = comentarios.len();

    let descripcion = modelo
        .resumir_comentarios(cliente_http, comentarios)
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
