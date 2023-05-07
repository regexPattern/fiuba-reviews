use axum::{
    extract::{Path, State},
    http::StatusCode,
    Json,
};
use serde::Serialize;
use sqlx::{FromRow, PgPool};

#[derive(Serialize, FromRow)]
pub struct Catedra {
    codigo: String,
    nombre: String,
}

pub async fn by_materia(
    State(pool): State<PgPool>,
    Path(codigo): Path<u32>,
) -> Result<Json<Vec<Catedra>>, StatusCode> {
    let catedras = sqlx::query_as::<_, Catedra>("SELECT * FROM catedras WHERE codigo_materia = $1")
        .bind(codigo as i32)
        .fetch_all(&pool)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(catedras))
}

#[derive(Serialize, FromRow)]
pub struct CatedraDocente {
    codigo_catedra: String,
    codigo_docente: String,
}

pub async fn by_docente(
    State(pool): State<PgPool>,
    Path(codigo): Path<String>,
) -> Result<Json<Vec<String>>, StatusCode> {
    let catedras = sqlx::query_as::<_, CatedraDocente>(
        "SELECT * FROM catedra_docente WHERE codigo_docente = $1",
    )
    .bind(codigo)
    .fetch_all(&pool)
    .await
    .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let codigos_catedras = catedras
        .into_iter()
        .map(|catedra| catedra.codigo_catedra)
        .collect();

    Ok(Json(codigos_catedras))
}
