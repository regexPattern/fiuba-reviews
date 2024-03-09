use resumidor_comentarios::gpt::OpenAIClient;
use sqlx::PgPool;
use std::env;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();
    let _ = dotenvy::dotenv();

    let modelo = OpenAIClient {
        api_key: env::var("OPENAI_API_KEY")
            .expect("variable de entorno `OPENAI_API_KEY` necesaria para conectar con OpenAI API"),
    };

    let database_url = env::var("DATABASE_URL")
        .expect("variable de entorno `DATABASE_URL` necesaria para conectar con la base de datos");

    let conexion_db = PgPool::connect(&database_url).await.unwrap();

    tracing::info!("conexion establecida con la base de datos");

    let query = resumidor_comentarios::query_actualizacion(&conexion_db, modelo).await?;

    if let Some(query_actualizacion) = query {
        sqlx::query(&query_actualizacion)
            .execute(&conexion_db)
            .await?;

        tracing::info!("base de datos actualizada");
    } else {
        tracing::info!("ningún docente se ha actualizado");
    }

    Ok(())
}
