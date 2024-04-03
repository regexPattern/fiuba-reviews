use std::collections::VecDeque;

use reqwest::{Client, StatusCode};
use serde::{Deserialize, Serialize};

use super::Modelo;

const CHAT_COMPLETION_API_ENDPOINT: &str = "https://api.openai.com/v1/chat/completions";

#[derive(Debug)]
pub struct OpenAIClient {
    pub api_key: String,
}

impl Modelo for OpenAIClient {
    async fn resumir_comentarios(
        &self,
        cliente_http: Client,
        comentarios: &[String],
    ) -> anyhow::Result<String> {
        let payload = OpenAIApiPayload {
            model: "gpt-3.5-turbo",
            messages: [
                PromptMessage {
                    role: "system",
                    content:
r"Los siguientes parrafos son comentarios de estudiantes sobre un profesor o profesora especifico.
Genera un resumen de los mismos. No me digás que los resultados son variados, andá directo al grano."
                    .to_string(),
                },
                PromptMessage {
                    role: "user",
                    content: comentarios.join("\n\n"),
                },
            ],
        };

        tracing::debug!("enviando request a OpenAI API");

        let res = cliente_http
            .post(CHAT_COMPLETION_API_ENDPOINT)
            .bearer_auth(&self.api_key)
            .json(&payload)
            .send()
            .await?;

        let status = res.status();

        tracing::debug!(
            "recibida respuesta de OpenAI API con status {}",
            status.as_u16()
        );

        if status != StatusCode::OK {
            let err = res.text().await?;
            tracing::error!("{err}");
            return Err(anyhow::anyhow!(err));
        }

        let mut res: OpenAIApiResponse = res
            .json()
            .await
            .expect("formato de respuesta exitosa de OpenAI API es diferente al esperado");

        let res = res
            .choices
            .pop_front()
            .ok_or(anyhow::anyhow!("respuesta no incluye ningun resumen"))?;

        tracing::info!("generado resumen de comentarios de docente");

        Ok(res.message.content)
    }
}

#[derive(Debug, Serialize)]
struct OpenAIApiPayload {
    model: &'static str,
    messages: [PromptMessage; 2],
}

#[derive(Debug, Serialize)]
struct PromptMessage {
    role: &'static str,
    content: String,
}

#[derive(Debug, Deserialize)]
struct OpenAIApiResponse {
    choices: VecDeque<ChatCompletionChoice>,
}

#[derive(Debug, Deserialize)]
struct ChatCompletionChoice {
    message: ModelGeneratedMessage,
}

#[derive(Debug, Deserialize)]
struct ModelGeneratedMessage {
    content: String,
}
