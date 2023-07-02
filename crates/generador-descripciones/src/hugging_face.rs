use std::collections::VecDeque;

use format_serde_error::SerdeError;
use reqwest::Client;
use serde::{Deserialize, Serialize};

const HF_INFERENCE_API_URL: &str =
    "https://api-inference.huggingface.co/models/facebook/bart-large-cnn";

pub async fn generar_descripcion(
    cliente_http: Client,
    codigo_docente: &str,
    comentarios_docente: &[String],
    api_key: &str,
) -> anyhow::Result<String> {
    #[derive(Serialize)]
    struct Payload {
        inputs: String,
    }

    #[derive(Deserialize)]
    struct Resumen {
        #[serde(alias = "summary_text")]
        descripcion: String,
    }

    tracing::info!("generando descripcion para docente '{}'", codigo_docente);

    let res = cliente_http
        .post(HF_INFERENCE_API_URL)
        .bearer_auth(api_key)
        .json(&Payload {
            inputs: comentarios_docente.join("."),
        })
        .send()
        .await?;

    let data = res.text().await?;

    let mut resumenes: VecDeque<Resumen> = serde_json::from_str(&data).map_err(|err| {
        tracing::error!("{data}");
        err
    })?;

    if let Some(Resumen { descripcion }) = resumenes.pop_front() {
        Ok(descripcion)
    } else {
        Err(anyhow::anyhow!(
            "respuesta de Inference API no incluye resultados"
        ))
    }
}
