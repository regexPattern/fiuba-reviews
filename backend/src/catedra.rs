use axum::{
    extract::{Path, State},
    http::StatusCode,
    Json,
};
use serde::Serialize;
use sqlx::{FromRow, PgPool};

use crate::comentario::{self, Comentario};

#[derive(Serialize, FromRow)]
struct Docente {
    codigo: String,
    nombre: String,
    respuestas: i32,
    promedio: f64,
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

#[derive(Serialize)]
pub struct DocenteConComentarios {
    #[serde(flatten)]
    docente: Docente,
    comentarios: Vec<Comentario>,
}

#[derive(Serialize)]
pub struct ComentariosCatedra {
    nombre_catedra: String,
    docentes_con_comentarios: Vec<DocenteConComentarios>,
}

pub async fn docentes_con_comentarios(
    State(pool): State<PgPool>,
    Path(codigo_catedra): Path<String>,
) -> Result<Json<ComentariosCatedra>, StatusCode> {
    #[derive(FromRow)]
    struct NombreCatedra(String);

    let nombre_catedra = sqlx::query_as::<_, NombreCatedra>(
        r#"
SELECT nombre
FROM catedra
WHERE codigo = $1;
"#,
    )
    .bind(&codigo_catedra)
    .fetch_one(&pool)
    .await
    .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let docentes = sqlx::query_as::<_, Docente>(
        r#"
SELECT docente.*
FROM docente
JOIN catedra_docente
ON catedra_docente.codigo_docente = docente.codigo
AND catedra_docente.codigo_catedra = $1;
"#,
    )
    .bind(codigo_catedra)
    .fetch_all(&pool)
    .await
    .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let mut docentes_con_comentarios = Vec::new();
    for docente in docentes {
        if let Ok(comentarios) = comentario::comentarios_de_docente(&pool, &docente.codigo).await {
            docentes_con_comentarios.push(DocenteConComentarios {
                docente,
                comentarios,
            });
        }
    }

    Ok(Json(ComentariosCatedra {
        nombre_catedra: nombre_catedra.0,
        docentes_con_comentarios,
    }))
}
