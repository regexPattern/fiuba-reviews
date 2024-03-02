use resumidor_comentarios::gpt::OpenAIClient;
use sqlx::PgPool;

const DATABASE_URL_ENV_VAR: &str = "DATABASE_URL";

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();
    let _ = dotenvy::dotenv();

    let modelo = OpenAIClient::new();

    let database_url = std::env::var(DATABASE_URL_ENV_VAR).expect(const_format::concatcp!(
        "variable de entorno `",
        DATABASE_URL_ENV_VAR,
        "` necesaria para conectar con la base de datos"
    ));

    let conexion_db = PgPool::connect(&database_url).await.unwrap();

    tracing::info!("conexion establecida con la base de datos");

    let query = resumidor_comentarios::update_query(&conexion_db, modelo).await?;

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
