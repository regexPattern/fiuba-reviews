use std::io::Write;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();

    let query_sql = database::indexar_dolly().await?;
    let mut archivo_init_sql = std::fs::File::create("init.sql")?;
    archivo_init_sql.write_all(query_sql.as_bytes())?;

    Ok(())
}
