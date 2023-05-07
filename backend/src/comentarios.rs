use axum::{
    extract::{Path, State},
    http::StatusCode,
    Json,
};
use serde::Serialize;
use sqlx::{FromRow, PgPool};

#[derive(Serialize, FromRow)]
pub struct Comentario {
    codigo: String,
    codigo_docente: String,
    cuatrimestre: String,
    contenido: String,
}

pub async fn by_docente(
    State(pool): State<PgPool>,
    Path(codigo): Path<String>,
) -> Result<Json<Vec<Comentario>>, StatusCode> {
    let comentarios =
        sqlx::query_as::<_, Comentario>("SELECT * FROM comentarios WHERE codigo_docente = $1")
            .bind(codigo)
            .fetch_all(&pool)
            .await
            .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(comentarios))
}
