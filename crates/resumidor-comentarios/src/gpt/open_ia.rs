use std::collections::VecDeque;

use reqwest::Client;
use serde::{Deserialize, Serialize};

use super::Modelo;

const API_KEY_ENV_VAR: &str = "OPENAI_API_KEY";
const CHAT_COMPLETION_API_ENDPOINT: &str = "https://api.openai.com/v1/chat/completions";

#[derive(Debug)]
pub struct OpenAIClient {
    api_key: String,
}

impl OpenAIClient {
    pub fn new() -> Self {
        Self {
            api_key: std::env::var(API_KEY_ENV_VAR).expect(const_format::concatcp!(
                "variable de entorno `",
                API_KEY_ENV_VAR,
                "` necesaria para conectar con OpenAI API",
            )),
        }
    }
}

impl Modelo for OpenAIClient {
    async fn resumen_comentarios(
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

        let mut res: OpenAIApiResponse = res.json().await?;
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
