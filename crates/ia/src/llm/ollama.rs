use format_serde_error::SerdeError;
use reqwest::Client;
use serde::Deserialize;

use super::{ResumidorComentarios, Sanitizador};

const OLLAMA_ENDPOINT_URL: &str = "http://localhost:11434/api/chat";

// Podés revisar los modelos disponibles en la documentación de Ollama:
// https://ollama.com/library
const MODELO: &str = "llama3.1";

#[derive(Debug)]
pub struct OllamaClient {
    pub http_client: Client,
}

impl ResumidorComentarios for OllamaClient {
    // Documentación de los endpoints de la API de Chat Completions de Ollama:
    // https://github.com/ollama/ollama/blob/main/docs/api.md#generate-a-chat-completion
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
                    "content": include_str!("../../prompts/ollama.prompt"),
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
            "options": {
                "temperature": 0.5,
            },
            "stream": false,
        });

        tracing::debug!("enviando request a Ollama");

        let res = self
            .http_client
            .post(OLLAMA_ENDPOINT_URL)
            .json(&prompt)
            .send()
            .await?
            .error_for_status()?;

        let res_body = res.text().await?;

        let res: OllamaResponse =
            serde_json::from_str(&res_body).map_err(|e| SerdeError::new(res_body, e))?;

        tracing::info!("generado resumen de comentarios de docente");

        Ok(res.message.content)
    }
}

impl Sanitizador for OllamaClient {}

#[derive(Debug, Deserialize)]
struct OllamaResponse {
    message: OllamaResponseMsg,
}

#[derive(Debug, Deserialize)]
struct OllamaResponseMsg {
    content: String,
}
