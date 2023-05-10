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
    codigo_equivalencia: Option<i32>,
}

pub async fn index(State(pool): State<PgPool>) -> Result<Json<Vec<Materia>>, StatusCode> {
    let materias = sqlx::query_as::<_, Materia>(
        r#"
SELECT *
FROM materia;
"#,
    )
    .fetch_all(&pool)
    .await
    .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(materias))
}

#[derive(Serialize)]
pub struct CatedrasConNombreMateria {
    nombre_materia: String,
    catedras: Vec<Catedra>,
}

#[derive(Serialize, FromRow)]
pub struct Catedra {
    pub codigo: String,
    pub nombre: String,
    pub promedio: f64,
}

pub async fn catedras(
    State(pool): State<PgPool>,
    Path(codigo_materia): Path<u32>,
) -> Result<Json<CatedrasConNombreMateria>, StatusCode> {
    #[derive(FromRow)]
    struct NombreMateria(String);

    let nombre_materia = sqlx::query_as::<_, NombreMateria>(
        r#"
SELECT nombre
FROM materia
WHERE codigo = $1;
"#,
    )
    .bind(codigo_materia as i32)
    .fetch_one(&pool)
    .await
    .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let catedras = sqlx::query_as::<_, Catedra>(
        r#"
SELECT codigo, nombre, promedio
FROM catedra
WHERE codigo_materia = $1
ORDER BY promedio;
"#,
    )
    .bind(codigo_materia as i32)
    .fetch_all(&pool)
    .await
    .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(CatedrasConNombreMateria {
        nombre_materia: nombre_materia.0,
        catedras,
    }))
}
