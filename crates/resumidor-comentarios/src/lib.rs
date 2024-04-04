pub mod gpt;

use std::{collections::HashMap, sync::Arc};

use gpt::ModeloGpt;
use reqwest::Client;
use sqlx::{types::Uuid, FromRow, PgPool};
use tokio::sync::Semaphore;
use tracing::Instrument;

const MAX_SOLICITUDES_CONCURRENTES: usize = 5;
const MIN_COMENTARIOS_ACTUALIZACION: usize = 3;
const PROPORCION_COMENTARIOS_ACTUALIZACION: usize = 2;

#[derive(FromRow)]
struct Comentario {
    codigo_docente: Uuid,
    nombre_docente: String,
    contenido: String,
}

pub async fn query_actualizacion<M>(
    conexion: &PgPool,
    modelo_gpt: M,
    forzar_actualizacion: bool,
) -> anyhow::Result<Option<String>>
where
    M: ModeloGpt + Send + Sync + 'static,
{
    let cliente_http = Client::new();
    let modelo_gpt = Arc::new(modelo_gpt);

    let comentarios_por_docente = comentarios_por_docente(conexion, forzar_actualizacion).await?;
    let cantidad_docentes = comentarios_por_docente.len();

    let semaphore = Arc::new(Semaphore::new(MAX_SOLICITUDES_CONCURRENTES));
    let mut handles = Vec::with_capacity(cantidad_docentes);

    for (codigo, (nombre, comentarios)) in comentarios_por_docente {
        let cliente_http = Client::clone(&cliente_http);
        let modelo = Arc::clone(&modelo_gpt);

        let semaphore = Arc::clone(&semaphore);
        let span = tracing::debug_span!("docente", codigo = codigo.to_string());

        handles.push(tokio::spawn(
            async move {
                let _permit = semaphore.acquire().await?;
                tupla_values(
                    cliente_http,
                    &semaphore,
                    codigo,
                    &nombre,
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
UPDATE docente AS d
SET resumen_comentarios = val.resumen_comentarios,
    comentarios_ultimo_resumen = val.comentarios_ultimo_resumen
FROM (
    VALUES
        {}
)
AS val(codigo_docente, resumen_comentarios, comentarios_ultimo_resumen)
WHERE d.codigo::text = val.codigo_docente;
",
            values_a_actualizar.join(",\n        ")
        ))
    } else {
        None
    };

    Ok(query)
}

async fn comentarios_por_docente(
    conexion: &PgPool,
    forzar_actualizacion: bool,
) -> anyhow::Result<HashMap<Uuid, (String, Vec<String>)>> {
    let comentarios: Vec<Comentario> = sqlx::query_as(&format!(
        "\
SELECT d.codigo AS codigo_docente, d.nombre AS nombre_docente, c.contenido
FROM comentario c
INNER JOIN docente d
ON c.codigo_docente = d.codigo
WHERE c.codigo_docente IN (
  SELECT d.codigo
  FROM docente d
  INNER JOIN comentario c
  ON c.codigo_docente = d.codigo
  GROUP BY d.codigo
  HAVING COUNT(c) > (d.comentarios_ultimo_resumen * {})
  AND COUNT(c) > {}
);",
        PROPORCION_COMENTARIOS_ACTUALIZACION,
        if forzar_actualizacion {
            1
        } else {
            MIN_COMENTARIOS_ACTUALIZACION
        }
    ))
    .fetch_all(conexion)
    .await?;

    tracing::info!("comentarios obtenidos de la base de datos");

    let mut comentarios_por_docente: HashMap<Uuid, (String, Vec<String>)> = HashMap::new();

    for com in comentarios {
        let comentarios_de_docente = comentarios_por_docente
            .entry(com.codigo_docente)
            .or_insert((com.nombre_docente, Vec::new()));

        comentarios_de_docente.1.push(com.contenido);
    }

    let cantidad_docentes = comentarios_por_docente.len();
    tracing::info!("encontrados {cantidad_docentes} docentes que requiren actualizacion");

    Ok(comentarios_por_docente)
}

async fn tupla_values<M>(
    cliente_http: Client,
    semaphore: &Semaphore,
    codigo_docente: Uuid,
    nombre_docente: &str,
    comentarios: &[String],
    modelo_gpt: Arc<M>,
) -> anyhow::Result<String>
where
    M: ModeloGpt + Send + Sync + 'static,
{
    let cantidad_comentarios_actual = comentarios.len();

    let resumen_comentarios = modelo_gpt
        .resumir_comentarios(cliente_http, nombre_docente, comentarios)
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
