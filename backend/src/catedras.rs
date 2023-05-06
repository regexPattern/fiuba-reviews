use axum::{
    extract::{Path, State},
    http::StatusCode,
    Json,
};
use serde::Serialize;
use sqlx::PgPool;

#[derive(Serialize)]
pub struct Catedra {
    codigo: String,
    codigo_materia: i32,
}

pub async fn by_materia(
    State(pool): State<PgPool>,
    Path(codigo_materia): Path<u32>,
) -> Result<Json<Vec<String>>, StatusCode> {
    let catedras_de_materia = sqlx::query_as!(
        Catedra,
        "SELECT * FROM catedras WHERE codigo_materia = $1",
        codigo_materia as i32
    )
    .fetch_all(&pool)
    .await
    .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let codigos_catedras = catedras_de_materia
        .into_iter()
        .map(|catedra| catedra.codigo)
        .collect();

    Ok(Json(codigos_catedras))
}
