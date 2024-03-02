mod hugging_face;
mod open_ia;

pub use hugging_face::HuggingFaceClient;
pub use open_ia::OpenAIClient;

use reqwest::Client;

pub trait Modelo {
    fn resumir_comentarios(
        &self,
        cliente_http: Client,
        comentarios: &[String],
    ) -> impl std::future::Future<Output = anyhow::Result<String>> + Send;
}
