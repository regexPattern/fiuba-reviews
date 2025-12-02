use std::{env, io::Write};

use clap::{Parser, Subcommand};
use reqwest::Client;
use resumidor_comentarios::llm::OllamaClient;
use sqlx::postgres::PgPoolOptions;

const DATABASE_URL_ENV: &str = "DATABASE_URL";
const OPENAI_API_KEY_ENV: &str = "OPENAI_API_KEY";

#[derive(Parser)]
struct Cli {
    #[command(subcommand)]
    comand: Subcomando,
}

#[derive(Subcommand)]
enum Subcomando {
    Resumir,
    Sanitizar,
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let cli = Cli::parse();

    tracing_subscriber::fmt::init();
    dotenvy::dotenv()?;

    let db_url = env::var(DATABASE_URL_ENV).expect(const_format::formatcp!(
        "variable de entorno `{}` necesaria para conectar con la base de datos",
        DATABASE_URL_ENV
    ));

    let db = PgPoolOptions::new().connect(&db_url).await?;

    tracing::info!("conexion establecida con la base de datos");

    let _api_key = env::var(OPENAI_API_KEY_ENV).expect(const_format::formatcp!(
        "variable de entorno `{}` necesaria para conectar con OpenAI API",
        OPENAI_API_KEY_ENV
    ));

    let llm = OllamaClient {
        //api_key,
        http_client: Client::new(),
    };

    match cli.comand {
        Subcomando::Resumir => {
            if let Some(query) =
                resumidor_comentarios::query_actualizacion_resumenes(llm, &db).await?
            {
                let mut file = std::fs::File::create("update.sql")?;
                file.write_all(query.as_bytes())?;

                tracing::info!("query guardada en archivo `update.sql`");
            } else {
                tracing::info!("ningún docente se ha actualizado");
            }
        }
        Subcomando::Sanitizar => {}
    };

    Ok(())
}
