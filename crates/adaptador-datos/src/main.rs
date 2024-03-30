use std::{fs::File, io::Write};

use clap::{Parser, Subcommand};
use sqlx::PgPool;

#[derive(Parser)]
struct Cli {
    #[clap(subcommand)]
    command: Option<Command>,
}

#[derive(Subcommand)]
enum Command {
    Update {
        #[arg(short, long)]
        output: bool,
    },
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let cli = Cli::parse();

    tracing_subscriber::fmt::init();

    match cli.command {
        None => {
            let query = adaptador_datos::init_query().await?;
            let mut archivo = File::create("init.sql")?;
            archivo.write_all(query.as_bytes())?;
        }
        Some(Command::Update { output }) => {
            let db = PgPool::connect("postgres://postgres:postgres@localhost:5432").await?;
            let query = adaptador_datos::update_query(&db).await?;

            if output {
                let mut archivo = File::create("update.sql")?;
                archivo.write_all(query.as_bytes())?;
            } else {
                sqlx::query(&query).execute(&db).await?;
            }
        }
    };

    Ok(())
}
