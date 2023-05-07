use axum::{
    extract::{Path, State},
    http::StatusCode,
    Json,
};
use serde::Serialize;
use sqlx::{FromRow, PgPool};

#[derive(Serialize, FromRow)]
pub struct Docente {
    codigo: String,
    nombre: String,
    respuestas: i32,
    acepta_critica: Option<f64>,
    asistencia: Option<f64>,
    buen_trato: Option<f64>,
    claridad: Option<f64>,
    clase_organizada: Option<f64>,
    cumple_horarios: Option<f64>,
    fomenta_participacion: Option<f64>,
    panorama_amplio: Option<f64>,
    responde_mails: Option<f64>,
}

#[derive(Serialize, FromRow)]
pub struct CatedraDocente {
    codigo_catedra: String,
    codigo_docente: String,
}

pub async fn by_codigo(
    State(pool): State<PgPool>,
    Path(codigo): Path<String>,
) -> Result<Json<Docente>, StatusCode> {
    let docente = sqlx::query_as::<_, Docente>("SELECT * FROM docentes WHERE codigo = $1")
        .bind(codigo)
        .fetch_one(&pool)
        .await
        .map_err(|err| match err {
            sqlx::Error::RowNotFound => StatusCode::NOT_FOUND,
            _ => StatusCode::INTERNAL_SERVER_ERROR,
        })?;

    Ok(Json(docente))
}

pub async fn by_catedra(
    State(pool): State<PgPool>,
    Path(codigo): Path<String>,
) -> Result<Json<Vec<String>>, StatusCode> {
    let docentes_de_catedra = sqlx::query_as::<_, CatedraDocente>(
        "SELECT * FROM catedra_docente WHERE codigo_catedra = $1",
    )
    .bind(codigo)
    .fetch_all(&pool)
    .await
    .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let codigos_docentes = docentes_de_catedra
        .into_iter()
        .map(|docente| docente.codigo_docente)
        .collect();

    Ok(Json(codigos_docentes))
}
