use std::{fs::File, io::Write};

use clap::{Parser, Subcommand};
use sqlx::PgPool;

#[derive(Parser)]
struct Cli {
    #[clap(subcommand)]
    comando: Comando,
}

#[derive(Subcommand)]
enum Comando {
    Inicializar,
    Actualizar {
        #[clap(short, long)]
        commit: bool,
    },
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let cli = Cli::parse();

    tracing_subscriber::fmt::init();

    let (nombre_archivo, query) = match cli.comando {
        Comando::Inicializar => ("init.sql", adaptador_datos::query_inicializacion().await?),
        Comando::Actualizar { commit } => {
            let base_de_datos_url = std::env::var("DATABASE_URL").expect(
                "variable de entorno `DATABASE_URL` necesaria para conectar con la base de datos",
            );

            let conexion = PgPool::connect(&base_de_datos_url).await?;
            tracing::info!("conexion establecida con la base de datos");

            let query = adaptador_datos::query_actualizacion(&conexion).await?;

            if commit {
                tracing::info!("actualizando base de datos");

                if let Err(err) = sqlx::query(&query).execute(&conexion).await {
                    tracing::error!("error actualizando la base de datos");
                    tracing::debug!("descripcion error: {}", err);
                } else {
                    tracing::info!("base de datos actualizada exitosamente");
                    return Ok(());
                }
            }

            ("update.sql", query)
        }
    };

    let mut archivo = File::create(nombre_archivo)?;
    archivo.write_all(query.as_bytes())?;

    tracing::info!("query guardada en archivo `{}`", nombre_archivo);

    Ok(())
}
