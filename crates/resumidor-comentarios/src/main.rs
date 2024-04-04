use clap::Parser;
use resumidor_comentarios::gpt::OpenAIClient;
use sqlx::PgPool;
use std::{env, io::Write};

#[derive(Parser)]
struct Cli {
    /// Ejecutar la query directamente en la base de datos.
    #[clap(short, long)]
    commit: bool,

    /// Regenerar resumenes de comentarios para todos los docentes.
    #[clap(short, long)]
    force: bool,
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let cli = Cli::parse();
    let _ = dotenvy::dotenv();

    tracing_subscriber::fmt::init();

    let modelo_gpt = OpenAIClient {
        api_key: env::var("OPENAI_API_KEY")
            .expect("variable de entorno `OPENAI_API_KEY` necesaria para conectar con OpenAI API"),
    };

    let base_de_datos_url = env::var("DATABASE_URL")
        .expect("variable de entorno `DATABASE_URL` necesaria para conectar con la base de datos");

    let conexion = PgPool::connect(&base_de_datos_url).await?;
    tracing::info!("conexion establecida con la base de datos");

    let query =
        resumidor_comentarios::query_actualizacion(&conexion, modelo_gpt, cli.force).await?;

    if let Some(query) = query {
        if cli.commit {
            if let Err(err) = sqlx::query(&query).execute(&conexion).await {
                tracing::error!("error actualizando la base de datos");
                tracing::debug!("descripcion error: {}", err);
            } else {
                tracing::info!("base de datos actualizada exitosamente");
                return Ok(());
            }
        }

        let mut archivo = std::fs::File::create("update.sql")?;
        archivo.write_all(query.as_bytes())?;

        tracing::info!("query guardada en archivo `update.sql`");
    } else {
        tracing::info!("ningún docente se ha actualizado");
    }

    Ok(())
}
