use axum::http::StatusCode;
use serde::Serialize;
use sqlx::{FromRow, PgPool};

#[derive(Serialize, FromRow)]
pub struct Comentario {
    pub codigo: String,
    pub codigo_docente: String,
    pub cuatrimestre: String,
    pub contenido: String,
}

pub async fn comentarios_de_docente(
    pool: &PgPool,
    codigo_docente: &str,
) -> Result<Vec<Comentario>, StatusCode> {
    let comentarios = sqlx::query_as::<_, Comentario>(
        r#"
SELECT *
FROM comentario
WHERE codigo_docente = $1;
"#,
    )
    .bind(codigo_docente)
    .fetch_all(pool)
    .await
    .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(comentarios)
}
