use std::collections::VecDeque;

use reqwest::{Client, StatusCode};
use serde::{Deserialize, Serialize};

use super::ModeloGpt;

const CHAT_COMPLETION_API_ENDPOINT: &str = "https://api.openai.com/v1/chat/completions";

#[derive(Debug)]
pub struct OpenAIClient {
    pub api_key: String,
}

impl ModeloGpt for OpenAIClient {
    async fn resumir_comentarios(
        &self,
        cliente_http: Client,
        nombre_docente: &str,
        comentarios: &[String],
    ) -> anyhow::Result<String> {
        let prompt = OpenAIApiPrompt {
            model: "gpt-3.5-turbo",
            messages: [
                PromptMessage {
                    role: "system",
                    content:
"Vas a recibir el nombre o apellido de un profesor o profesora seguido de un listado de
comentarios suyos hechos por sus alumnos.

Genera un resumen de los mismos. Prioriza darle peso a las ideas que se repiten en varios
comentarios en vez de ideas que a lo mejor se repiten varias veces pero únicamente dentro de un
comentario largo. No tomes en cuenta comentarios sobre la vida personal del docente como su
religión o su postura política.

Si el docente tiene pocos comentarios no muy extensos, hace un resumen corto. Si el docente tiene
muchos comentarios muy extensos, limitate a un máximo de alrededor de 120 palabras.

El resumen debe sonar natural, como si lo estuviera diciendo una persona en una conversación
casual, es decir, que por ejemplo, no hay necesidad de introducirlo diciendo que esto es un resumen
ni nada por el estilo."
                    .to_string(),
                },
                PromptMessage {
                    role: "user",
                    content: format!("\
Nombre docente: {nombre_docente}

Listado de comentarios:
- {}", comentarios.join("\n\n- ")),
                },
            ],
        };

        tracing::debug!("enviando request a OpenAI API");

        let res = cliente_http
            .post(CHAT_COMPLETION_API_ENDPOINT)
            .bearer_auth(&self.api_key)
            .json(&prompt)
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
struct OpenAIApiPrompt {
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
