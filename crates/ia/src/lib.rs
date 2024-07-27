pub mod llm;
mod sql;

use std::{collections::HashMap, sync::Arc};

use llm::ResumidorComentarios;
use sql::Sql;
use sqlx::types::Uuid;
use tokio::{sync::Semaphore, task::JoinHandle};
use tracing::Instrument;

const MAX_SOLICITUDES_CONCURRENTES: usize = 5;

pub async fn actualizar_comentarios<L>(
    llm: L,
    comentarios_de_docente: HashMap<Uuid, (String, Vec<String>)>,
) -> anyhow::Result<Option<String>>
where
    L: ResumidorComentarios + Send + Sync + 'static,
{
    let llm = Arc::new(llm);
    let semaphore = Arc::new(Semaphore::new(MAX_SOLICITUDES_CONCURRENTES));

    let cantidad_docentes = comentarios_de_docente.len();

    let mut handles: Vec<JoinHandle<anyhow::Result<(Uuid, String, usize)>>> =
        Vec::with_capacity(cantidad_docentes);

    for (codigo_doc, (nombre_doc, comentarios)) in comentarios_de_docente {
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

    let resumenes: Vec<_> = resumenes.into_iter().map(|r| r.sql()).collect();

    dbg!(&resumenes);

    Ok(if !resumenes.is_empty() {
        Some(format!(
            r"
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
            resumenes.join(",\n        ")
        ))
    } else {
        None
    })
}
