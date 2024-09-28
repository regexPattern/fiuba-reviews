use std::collections::VecDeque;

use format_serde_error::SerdeError;
use reqwest::Client;
use serde::Deserialize;

use super::{ResumidorComentarios, Sanitizador};

const OPENAI_API_ENDPOINT_URL: &str = "https://api.openai.com/v1/chat/completions";

// Podés revisar los modelos disponibles en la documentación de OpenAI:
// https://platform.openai.com/docs/models
const MODELO: &str = "gpt-4o-mini";

#[derive(Debug)]
pub struct OpenAiApiClient {
    pub api_key: String,
    pub http_client: Client,
}

impl ResumidorComentarios for OpenAiApiClient {
    // Documentación de los endpoints de la API de Chat Completions de OpenAI:
    // https://platform.openai.com/docs/api-reference/chat/create
    //
    async fn generar_resumen(
        &self,
        comentarios: &[String],
        nombre_docente: String,
    ) -> anyhow::Result<String> {
        let prompt = serde_json::json!({
            "model": MODELO,
            "messages": [
                {
                    "role": "system",
                    "content": include_str!("../../prompts/open_ai.prompt"),
                },
                {
                    "role": "user",
                    "content": format!(
                        r"
                        Nombre docente: {nombre_docente}
                        Listado de comentarios:
                        - {}",
                        comentarios.join("\n\n- ")
                    ),
                },
            ],
            "temperature": 0.7,
            "stream": false,
        });

        tracing::debug!("enviando request a OpenAI");

        let res = self
            .http_client
            .post(OPENAI_API_ENDPOINT_URL)
            .bearer_auth(&self.api_key)
            .json(&prompt)
            .send()
            .await?
            .error_for_status()?;

        let res_body = res.text().await?;

        let mut res: OpenAiApiResponse =
            serde_json::from_str(&res_body).map_err(|e| SerdeError::new(res_body, e))?;

        tracing::info!("generado resumen de comentarios de docente");

        let res = res.choices.pop_front().ok_or(anyhow::anyhow!(
            "respuesta de OpenAI API no incluye ningun resumen"
        ))?;

        Ok(res.message.content)
    }
}

impl Sanitizador for OpenAiApiClient {}

#[derive(Debug, Deserialize)]
struct OpenAiApiResponse {
    choices: VecDeque<OpenAiApiChoice>,
}

#[derive(Debug, Deserialize)]
struct OpenAiApiChoice {
    message: OpenAiApiResponseMsg,
}

#[derive(Debug, Deserialize)]
struct OpenAiApiResponseMsg {
    content: String,
}
