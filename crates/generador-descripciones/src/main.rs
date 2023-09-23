use std::env;

use sqlx::PgPool;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();

    let _ = dotenvy::dotenv();

    let connection_url = env::var("DATABASE_URL")
        .expect("variable de entorno DATABASE_URL necesaria para conectar con la base de datos");

    let inference_api_key = env::var("INFERENCE_API_KEY").expect("variable de entorno INFERENCE_API_KEY necesaria para conectar con Hugging Face Inference API");

    let conexion_db = PgPool::connect(&connection_url).await.unwrap();

    let query_actualizacion =
        generador_descripciones::query_actualizacion(&conexion_db, inference_api_key).await?;

    if let Some(query_actualizacion) = query_actualizacion {
        sqlx::query(&query_actualizacion)
            .execute(&conexion_db)
            .await?;

        tracing::info!("base de datos actualizada");
    } else {
        tracing::info!("ningun docente se ha actualizado");
    }

    Ok(())
}
