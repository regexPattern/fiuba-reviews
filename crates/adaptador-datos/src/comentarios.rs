use std::{collections::HashMap, hash::Hash};

use base64::{engine::general_purpose, Engine};
use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;

use crate::{docentes, sql::Sql};

const URL_DESCARGA: &str = "https://dollyfiuba.com/analitics/comentarios_docentes.json";

#[derive(Debug)]
pub struct Comentario;

#[derive(Deserialize, PartialEq, Eq, Hash, Debug)]
pub struct MetaDataComentario {
    pub codigo_docente: String,
    pub nombre_cuatrimestre: String,
}

impl Comentario {
    pub async fn descargar_todos(
        cliente_http: &ClientWithMiddleware,
    ) -> anyhow::Result<HashMap<MetaDataComentario, Vec<String>>> {
        #[derive(Deserialize)]
        struct RespuestaDolly {
            #[serde(flatten)]
            metadata: MetaDataComentariosDolly,
            comentarios: Vec<Option<String>>,
        }

        tracing::info!("descargando listado de comentarios");

        let res = cliente_http.get(URL_DESCARGA).send().await?;
        let data = res.text().await?;

        let respuestas: Vec<RespuestaDolly> =
            serde_json::from_str(&data).map_err(|err| SerdeError::new(data, err))?;

        Ok(respuestas
            .into_iter()
            .map(|respuesta| {
                let comentarios: Vec<_> = respuesta
                    .comentarios
                    .into_iter()
                    .flatten()
                    .filter_map(|comentario| {
                        String::from_utf8(general_purpose::STANDARD.decode(comentario).ok()?).ok()
                    })
                    .collect();

                (MetaDataComentario::from(respuesta.metadata), comentarios)
            })
            .collect())
    }

    pub fn query_sql(metadata: &MetaDataComentario, comentarios: &[String]) -> String {
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

#[derive(Deserialize)]
struct MetaDataComentariosDolly {
    #[serde(alias = "cuat")]
    nombre_cuatrimestre: String,

    #[serde(alias = "mat")]
    codigo_materia: u32,

    #[serde(alias = "doc")]
    nombre_docente: String,
}

impl From<MetaDataComentariosDolly> for MetaDataComentario {
    fn from(metadata: MetaDataComentariosDolly) -> Self {
        Self {
            codigo_docente: docentes::generar_codigo_docente(
                metadata.codigo_materia,
                &metadata.nombre_docente,
            ),
            nombre_cuatrimestre: metadata.nombre_cuatrimestre,
        }
    }
}
