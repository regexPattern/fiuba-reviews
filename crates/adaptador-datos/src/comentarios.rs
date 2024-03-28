use std::{collections::HashMap, hash::Hash};

use base64::{engine::general_purpose, Engine};
use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;

use crate::sql::Sql;

const URL_DESCARGA_COMENTARIOS: &str = "https://dollyfiuba.com/analitics/comentarios_docentes.json";

#[derive(Debug)]
pub struct Comentario;

#[derive(Deserialize, PartialEq, Eq, Hash, Debug)]
pub struct MetaData {
    #[serde(alias = "mat")]
    pub codigo_materia: u32,

    #[serde(alias = "doc")]
    pub nombre_docente: String,

    #[serde(alias = "cuat")]
    pub nombre_cuatrimestre: String,
}

impl Comentario {
    pub async fn descargar_todos(
        cliente_http: &ClientWithMiddleware,
    ) -> anyhow::Result<HashMap<MetaData, Vec<String>>> {
        #[derive(Deserialize)]
        struct RespuestaDolly {
            #[serde(flatten)]
            metadata: MetaData,
            comentarios: Vec<Option<String>>,
        }

        tracing::info!("descargando listado de comentarios");

        let res = cliente_http.get(URL_DESCARGA_COMENTARIOS).send().await?;
        let data = res.text().await?;

        let respuestas: Vec<RespuestaDolly> =
            serde_json::from_str(&data).map_err(|err| SerdeError::new(data, err))?;

        Ok(respuestas
            .into_iter()
            .map(|res| {
                let comentarios: Vec<_> = res
                    .comentarios
                    .into_iter()
                    .flatten()
                    .filter_map(|com| {
                        String::from_utf8(general_purpose::STANDARD.decode(com).ok()?).ok()
                    })
                    .collect();

                (res.metadata, comentarios)
            })
            .collect())
    }

    pub fn query_sql(metadata: &MetaData, comentarios: &[String]) -> String {
        comentarios
            .iter()
            .map(|contenido| {
                format!(
                    r#"
INSERT INTO comentario(cuatrimestre, codigo_docente, contenido)
VALUES ('{}', '{}', '{}');
"#,
                    metadata.nombre_cuatrimestre.sanitizar_sql(),
                    metadata.codigo_docente.sanitizar_sql(),
                    contenido.sanitizar_sql()
                )
            })
            .collect::<Vec<_>>()
            .join("")
    }
    pub fn cuatrimestre_query_sql(nombre_cuatrimestre: &str) -> String {
        format!(
            r#"
INSERT INTO cuatrimestre(nombre)
VALUES ('{}');
        "#,
            nombre_cuatrimestre.sanitizar_sql()
        )
    }
}
