mod ollama;
mod open_ai;

pub use ollama::OllamaClient;
pub use open_ai::OpenAiApiClient;

pub trait ResumidorComentarios {
    fn generar_resumen(
        &self,
        comentarios: &[String],
        nombre_docente: String,
    ) -> impl std::future::Future<Output = anyhow::Result<String>> + Send;
}

pub trait Sanitizador {}
