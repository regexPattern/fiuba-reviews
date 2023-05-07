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

pub async fn get_all(State(pool): State<PgPool>) -> Result<Json<Vec<Materia>>, StatusCode> {
    let materias = sqlx::query_as::<_, Materia>("SELECT * FROM materias")
        .fetch_all(&pool)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(materias))
}

pub async fn by_codigo(
    State(pool): State<PgPool>,
    Path(codigo): Path<u32>,
) -> Result<Json<Materia>, StatusCode> {
    let materia = sqlx::query_as::<_, Materia>("SELECT * FROM materias WHERE codigo = $1")
        .bind(codigo as i32)
        .fetch_one(&pool)
        .await
        .map_err(|err| match err {
            sqlx::Error::RowNotFound => StatusCode::NOT_FOUND,
            _ => StatusCode::INTERNAL_SERVER_ERROR,
        })?;

    Ok(Json(materia))
}
