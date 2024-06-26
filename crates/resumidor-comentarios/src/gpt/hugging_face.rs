use std::collections::{HashMap, VecDeque};

use reqwest::{Client, StatusCode};
use serde::{Deserialize, Serialize};

use super::ModeloGpt;

const INFERENCE_API_ENDPOINT: &str =
    "https://api-inference.huggingface.co/models/facebook/bart-large-cnn";

#[derive(Debug)]
pub struct HuggingFaceClient {
    pub api_key: String,
}

impl ModeloGpt for HuggingFaceClient {
    async fn resumir_comentarios(
        &self,
        cliente_http: Client,
        _nombre_docente: &str,
        comentarios: &[String],
    ) -> anyhow::Result<String> {
        let prompt = InferenceApiPrompt {
            inputs: comentarios.join("."),
            options: [("wait_for_model".to_string(), true)].into(),
        };

        tracing::debug!("enviando request a Inference API");

        let res = cliente_http
            .post(INFERENCE_API_ENDPOINT)
            .bearer_auth(&self.api_key)
            .json(&prompt)
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

        let mut res: VecDeque<InferenceApiResponse> = res
            .json()
            .await
            .expect("formato de respuesta exitosa de Inference API es diferente al esperado");

        let res = res
            .pop_front()
            .ok_or(anyhow::anyhow!("respuesta no incluye ningun resumen"))?;

        tracing::info!("generada resumen de comentarios de docente");

        Ok(res.summary_text)
    }
}

#[derive(Debug, Serialize)]
struct InferenceApiPrompt {
    inputs: String,
    options: HashMap<String, bool>,
}

#[derive(Debug, Deserialize)]
struct InferenceApiResponse {
    summary_text: String,
}
