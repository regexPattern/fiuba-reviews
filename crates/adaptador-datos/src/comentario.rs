use std::collections::HashMap;

use base64::{engine::general_purpose, Engine};
use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;
use uuid::Uuid;

use crate::sql::Sql;

const URL_DESCARGA_COMENTARIOS: &str = "https://dollyfiuba.com/analitics/comentarios_docentes.json";

#[derive(Deserialize, PartialEq, Eq, Hash, Debug)]
pub struct MetaDataComentario {
    #[serde(alias = "mat")]
    pub codigo_materia: i16,

    #[serde(alias = "doc")]
    pub nombre_docente: String,

    #[serde(alias = "cuat")]
    pub nombre_cuatrimestre: String,
}

pub async fn descargar_todos(
    cliente_http: &ClientWithMiddleware,
) -> anyhow::Result<HashMap<MetaDataComentario, Vec<String>>> {
    #[derive(Deserialize)]
    struct RespuestaDolly {
        #[serde(flatten)]
        metadata: MetaDataComentario,
        comentarios: Vec<Option<String>>,
    }

    tracing::info!("descargando listado de comentarios");

    let res = cliente_http.get(URL_DESCARGA_COMENTARIOS).send().await?;
    let data = res.text().await?;

    let res: Vec<RespuestaDolly> =
        serde_json::from_str(&data).map_err(|err| SerdeError::new(data, err))?;

    let comentarios = res
        .into_iter()
        .map(|res| {
            let comentarios: Vec<_> = res
                .comentarios
                .into_iter()
                .flatten()
                .filter_map(|c| String::from_utf8(general_purpose::STANDARD.decode(c).ok()?).ok())
                .collect();

            (res.metadata, comentarios)
        })
        .collect();

    Ok(comentarios)
}

pub fn sql_cuatrimestre(nombre_cuatrimestre: &str) -> String {
    let nombre = nombre_cuatrimestre.sanitize();
    format!("({nombre})")
}

pub fn sql_comentario(comentario: &str, codigo_docente: &Uuid, cuatrimestre: &str) -> String {
    let codigo = Uuid::new_v4().sanitize();
    let codigo_docente = codigo_docente.sanitize();
    let cuatrimestre = cuatrimestre.sanitize();
    let contenido = comentario.sanitize();
    let es_de_dolly = true;

    format!("({codigo}, {codigo_docente}, {cuatrimestre}, {contenido}, {es_de_dolly})")
}

pub fn bulk_insert_cuatrimestre(insert_tuples: &Vec<String>) -> String {
    format!(
        "INSERT INTO cuatrimestre (nombre)
VALUES
\t{}
ON CONFLICT (nombre)
DO NOTHING;",
        insert_tuples.sanitize()
    )
}

pub fn bulk_insert_comentarios(insert_tuples: &Vec<String>) -> String {
    format!(
        "INSERT INTO comentario (codigo, codigo_docente, cuatrimestre, contenido, es_de_dolly)
VALUES
\t{};",
        insert_tuples.sanitize()
    )
}
