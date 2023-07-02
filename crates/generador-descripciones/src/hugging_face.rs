use std::collections::{HashMap, VecDeque};

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
    struct Params {
        inputs: String,
        options: HashMap<String, bool>,
    }

    #[derive(Deserialize)]
    struct Resumen {
        #[serde(alias = "summary_text")]
        descripcion: String,
    }

    tracing::debug!("generando descripcion para docente '{}'", codigo_docente);

    let res = cliente_http
        .post(HF_INFERENCE_API_URL)
        .bearer_auth(api_key)
        .json(&Params {
            inputs: comentarios_docente.join("."),
            options: [("wait_for_model".to_string(), true)].into(),
        })
        .send()
        .await?;

    let data = res.text().await?;

    let mut resumenes: VecDeque<Resumen> = if let Ok(resumenes) = serde_json::from_str(&data) {
        resumenes
    } else {
        anyhow::bail!("error al deserializar respuesta de Inference API: {data}");
    };

    if let Some(Resumen { descripcion }) = resumenes.pop_front() {
        Ok(descripcion)
    } else {
        Err(anyhow::anyhow!(
            "respuesta de Inference API no incluye resultados"
        ))
    }
}
