use sqlx::postgres::PgPoolOptions;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();

    let url_conexion = std::env::var("DATABASE_URL").unwrap();
    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(&url_conexion)
        .await?;

    let query_sql = generador_descripciones::query_actualizacion_db(&pool).await?;

    tracing::info!("actualizando registros en la base de datos");
    sqlx::query(&query_sql).execute(&pool).await?;

    Ok(())
}
