use std::{fs::File, io::Write};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();

    let sql = database::run().await?;
    let mut init_sql_file = File::create("init.sql")?;
    init_sql_file.write_all(sql.as_bytes())?;

    Ok(())
}
