use anyhow::Context;
use sqlx::postgres::PgPoolOptions;

const DATABASE_URL: &str = "DATABASE_URL";
const HF_INFERENCE_API_KEY: &str = "HF_INFERENCE_API_KEY";

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();

    let url_conexion = obtener_variable_de_entorno(DATABASE_URL)?;
    let api_key = obtener_variable_de_entorno(HF_INFERENCE_API_KEY)?;

    tracing::info!("estableciendo conexion con la base de datos");
    let conexion = PgPoolOptions::new()
        .max_connections(5)
        .connect(&url_conexion)
        .await?;

    let query_sql = generador_descripciones::actualizar(&conexion, api_key).await?;

    if let Some(query_sql) = query_sql {
        tracing::info!("actualizando registros en la base de datos");
        sqlx::query(&query_sql).execute(&conexion).await?;
    } else {
        tracing::info!("todos los registros estan actualizados");
    }

    Ok(())
}

fn obtener_variable_de_entorno(variable: &str) -> anyhow::Result<String> {
    Ok(std::env::var(variable)
        .with_context(|| format!("variable de entorno {variable} no configurada"))?)
}
