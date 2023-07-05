use std::collections::{HashMap, VecDeque};

use reqwest::{Client, StatusCode};

const INFERENCE_API_URL: &str =
    "https://api-inference.huggingface.co/models/facebook/bart-large-cnn";

#[derive(serde::Serialize)]
struct ParametrosModelo {
    inputs: String,
    options: HashMap<String, bool>,
}

#[derive(serde::Deserialize)]
struct RespuestaModelo {
    summary_text: String,
}

pub async fn generar_descripcion(
    cliente_http: Client,
    comentarios: &[String],
    api_key: &str,
) -> anyhow::Result<String> {
    tracing::debug!("enviando request a Inference API");

    let res = cliente_http
        .post(INFERENCE_API_URL)
        .bearer_auth(api_key)
        .json(&ParametrosModelo {
            inputs: comentarios.join("."),
            options: [("wait_for_model".to_string(), true)].into(),
        })
        .send()
        .await?;

    let status = res.status();

    tracing::debug!(
        "recibida respuesta de Inference API con status {}",
        status.as_u16()
    );

    if status != StatusCode::OK {
        let err = anyhow::anyhow!(if status == StatusCode::TOO_MANY_REQUESTS {
            "alcanzado el maximo de requests por hora de Inference API".to_string()
        } else {
            res.text().await?
        });

        tracing::error!("{err}");
        return Err(err);
    }

    let mut respuetas: VecDeque<RespuestaModelo> = res.json().await?;
    let respuesta = respuetas
        .pop_front()
        .ok_or(anyhow::anyhow!("respuesta no incluye ningun resumen"))?;

    tracing::info!("generada descripcion de docente");

    Ok(respuesta.summary_text)
}
