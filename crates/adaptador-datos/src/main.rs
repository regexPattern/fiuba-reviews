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
        #[clap(short, long)]
        commit: bool,
    },
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let cli = Cli::parse();

    tracing_subscriber::fmt::init();

    let query = match cli.command {
        None => adaptador_datos::init_query().await?,
        Some(Command::Update { commit }) => {
            let db_url = std::env::var("DATABASE_URL").expect(
                "variable de entorno `DATABASE_URL` necesaria para conectar con la base de datos",
            );

            let db = PgPool::connect(&db_url).await?;
            tracing::info!("conexion establecida con la base de datos");

            let query = adaptador_datos::update_query(&db).await?;

            if commit {
                tracing::info!("actualizando base de datos");

                if let Err(err) = sqlx::query(&query).execute(&db).await {
                    tracing::error!("error actualizando la base de datos");
                    tracing::debug!("descripcion error: {}", err);
                } else {
                    tracing::info!("base de datos actualizada exitosamente");
                    return Ok(());
                }
            }

            query
        }
    };

    let mut archivo = File::create("init.sql")?;
    archivo.write_all(query.as_bytes())?;

    tracing::info!("query guardada en archivo `init.sql`");

    Ok(())
}
