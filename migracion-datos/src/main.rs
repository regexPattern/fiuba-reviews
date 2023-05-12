use std::{fs::File, io::Write};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();

    let init_sql = migrador_datos::generar_sql().await?;
    let mut archivo = File::create("init.sql")?;
    archivo.write_all(init_sql.as_bytes())?;

    Ok(())
}
