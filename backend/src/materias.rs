use axum::{
    extract::{Path, State},
    http::StatusCode,
    Json,
};
use serde::Serialize;
use sqlx::{FromRow, PgPool};

#[derive(Serialize, FromRow)]
pub struct Materia {
    codigo: i32,
    nombre: String,
}

pub async fn listar(State(pool): State<PgPool>) -> Result<Json<Vec<Materia>>, StatusCode> {
    let materias = sqlx::query_as::<_, Materia>("SELECT * FROM materias")
        .fetch_all(&pool)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(materias))
}

pub async fn informacion(
    State(pool): State<PgPool>,
    Path(codigo_materia): Path<u32>,
) -> Result<Json<Materia>, StatusCode> {
    let materia = sqlx::query_as::<_, Materia>("SELECT * FROM materias WHERE codigo = $1")
        .bind(codigo_materia as i32)
        .fetch_one(&pool)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(materia))
}
