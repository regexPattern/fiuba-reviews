use std::{fs::File, io::Write};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();

    let data_inicializacion = adaptador::run().await?;
    let mut archivo_init_sql = File::create("init.sql")?;
    archivo_init_sql.write_all(data_inicializacion.as_bytes())?;

    Ok(())
}
