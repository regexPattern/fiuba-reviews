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
    conexion: &PgPool,
    modelo_gpt: M,
) -> anyhow::Result<Option<String>>
where
    M: Modelo + Send + Sync + 'static,
{
    let cliente_http = Client::new();
    let modelo_gpt = Arc::new(modelo_gpt);

    let comentarios_por_docente = comentarios_por_docente(conexion).await?;
    let cantidad_docentes = comentarios_por_docente.len();

    let semaphore = Arc::new(Semaphore::new(MAX_SOLICITUDES_CONCURRENTES));
    let mut handles = Vec::with_capacity(cantidad_docentes);

    for (codigo_docente, comentarios) in comentarios_por_docente {
        let cliente_http = Client::clone(&cliente_http);
        let modelo = Arc::clone(&modelo_gpt);

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

    for task in handles {
        if let Ok(tupla_values) = task.await.unwrap() {
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
            "\
UPDATE docente AS doc
SET resumen_comentarios = val.resumen_comentarios,
    comentarios_ultimo_resumen = val.comentarios_ultimo_resumen
FROM (
    VALUES
        {}
)
AS val(codigo_docente, resumen_comentarios, comentarios_ultimo_resumen)
WHERE doc.codigo::text = val.codigo_docente;
",
            values_a_actualizar.join(",\n        ")
        ))
    } else {
        None
    };

    Ok(query)
}

async fn comentarios_por_docente(conexion: &PgPool) -> anyhow::Result<HashMap<Uuid, Vec<String>>> {
    let comentarios: Vec<Comentario> = sqlx::query_as(const_format::formatcp!(
        "\
SELECT com.codigo_docente, com.contenido
FROM comentario com
WHERE com.codigo_docente IN (
  SELECT doc.codigo
  FROM docente doc
  INNER JOIN comentario com
  ON com.codigo_docente = doc.codigo
  GROUP BY doc.codigo
  HAVING COUNT(com) > (doc.comentarios_ultimo_resumen * {})
  AND COUNT(com) > {}
);
",
        PROPORCION_COMENTARIOS_ACTUALIZACION,
        MIN_COMENTARIOS_ACTUALIZACION
    ))
    .fetch_all(conexion)
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
    modelo_gpt: Arc<M>,
) -> anyhow::Result<String>
where
    M: Modelo + Send + Sync + 'static,
{
    let cantidad_comentarios_actual = comentarios.len();

    let resumen_comentarios = modelo_gpt
        .resumir_comentarios(cliente_http, comentarios)
        .await
        .map_err(|err| {
            semaphore.close();
            err
        })?;

    let tupla_values = format!(
        "('{}', '{}', {})",
        codigo_docente,
        resumen_comentarios.replace('\'', "''"),
        cantidad_comentarios_actual
    );

    Ok(tupla_values)
}
