pub mod llm;
mod sql;

use std::sync::Arc;

use llm::{ResumidorComentarios, Sanitizador};
use sql::queries;
use sqlx::{types::Uuid, PgPool};
use tokio::{sync::Semaphore, task::JoinHandle};
use tracing::Instrument;

const MAX_SOLICITUDES_CONCURRENTES: usize = 5;

pub async fn query_actualizacion_resumenes<L>(llm: L, db: &PgPool) -> anyhow::Result<Option<String>>
where
    L: ResumidorComentarios + Send + Sync + 'static,
{
    let llm = Arc::new(llm);
    let semaphore = Arc::new(Semaphore::new(MAX_SOLICITUDES_CONCURRENTES));

    let comentarios_docentes = queries::comentarios_docentes(db).await?;
    let cantidad_docentes = comentarios_docentes.len();

    let mut handles: Vec<JoinHandle<anyhow::Result<(Uuid, String, usize)>>> =
        Vec::with_capacity(cantidad_docentes);

    for (codigo_doc, (nombre_doc, comentarios)) in comentarios_docentes {
        let llm = Arc::clone(&llm);
        let semaphore = Arc::clone(&semaphore);
        let span = tracing::debug_span!("docente", codigo = codigo_doc.to_string());

        handles.push(tokio::spawn(
            async move {
                let _permit = semaphore.acquire().await?;
                let resumen = llm.generar_resumen(&comentarios, nombre_doc).await?;
                Ok((codigo_doc, resumen, comentarios.len()))
            }
            .instrument(span),
        ));
    }

    let mut resumenes = Vec::with_capacity(handles.len());

    for task in handles {
        match task.await.unwrap() {
            Ok(resumen) => resumenes.push(resumen),
            Err(err) => {
                tracing::error!("{err}");
            }
        }
    }

    tracing::info!(
        "query de actualizacion generada para {} de {} docentes",
        resumenes.len(),
        cantidad_docentes
    );

    Ok(queries::actualizar_resumen_docentes(resumenes).await)
}

pub async fn actualizar_nombres_materias<L>(
    _llm: L,
    _nombres_materias: Vec<(i16, String)>,
) -> anyhow::Result<Option<String>>
where
    L: Sanitizador,
{
    Ok(Some("".into()))
}

pub async fn actualizar_nombres_docentes<L>(
    _llm: L,
    _nombres_docentes: Vec<(Uuid, String)>,
) -> anyhow::Result<Option<String>>
where
    L: Sanitizador,
{
    Ok(Some("".into()))
}
