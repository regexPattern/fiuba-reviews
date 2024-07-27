use std::{collections::HashMap, env, io::Write};

use reqwest::Client;
use sqlx::{postgres::PgPoolOptions, types::Uuid, PgPool};

const MIN_COMENTARIOS_ACTUALIZACION: usize = 3;
const PROPORCION_COMENTARIOS_ACTUALIZACION: usize = 2;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();
    dotenvy::dotenv()?;

    let db_conn_str = env::var("DATABASE_URL")
        .expect("variable de entorno `DATABASE_URL` necesaria para conectar con la base de datos");
    let db = PgPoolOptions::new().connect(&db_conn_str).await?;

    tracing::info!("conexion establecida con la base de datos");

    let api_key = env::var("OPENAI_API_KEY")
        .expect("variable de entorno `OPENAI_API_KEY` necesaria para conectar con OpenAI API");
    let http_client = Client::new();

    let llm = ia::llm::OpenAiApiClient {
        api_key,
        http_client,
    };

    let comentarios_por_docente = comentarios_de_docente(&db).await?;
    let query = ia::actualizar_comentarios(llm, comentarios_por_docente).await?;

    if let Some(query) = query {
        let mut file = std::fs::File::create("update.sql")?;
        file.write_all(query.as_bytes())?;

        tracing::info!("query guardada en archivo `update.sql`");
    } else {
        tracing::info!("ningún docente se ha actualizado");
    }

    Ok(())
}

async fn comentarios_de_docente(
    db: &PgPool,
) -> anyhow::Result<HashMap<Uuid, (String, Vec<String>)>> {
    let comentarios: Vec<(Uuid, String, String)> = sqlx::query_as(&format!(
        r"
SELECT d.codigo, d.nombre, c.contenido
FROM comentario c
INNER JOIN docente d
ON c.codigo_docente = d.codigo
WHERE c.codigo_docente IN (
  SELECT d.codigo
  FROM docente d
  INNER JOIN comentario c
  ON c.codigo_docente = d.codigo
  GROUP BY d.codigo
  HAVING COUNT(c) > (d.comentarios_ultimo_resumen * {})
  AND COUNT(c) > {}
);
",
        PROPORCION_COMENTARIOS_ACTUALIZACION, MIN_COMENTARIOS_ACTUALIZACION
    ))
    .fetch_all(db)
    .await?;

    tracing::info!("comentarios obtenidos de la base de datos");

    let mut comentarios_de_docente = HashMap::new();

    for (codigo_docente, nombre_docente, contenido) in comentarios {
        let comentarios_de_docente = comentarios_de_docente
            .entry(codigo_docente)
            .or_insert((nombre_docente, Vec::new()));

        comentarios_de_docente.1.push(contenido);
    }

    let cantidad_docentes = comentarios_de_docente.len();

    tracing::info!("encontrados {cantidad_docentes} docentes que requiren actualización");

    Ok(comentarios_de_docente.into_iter().take(5).collect())
}
