use std::{collections::HashMap, hash::Hash};

use base64::{engine::general_purpose, Engine};
use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;
use uuid::Uuid;

use crate::sql::Sql;

const URL_DESCARGA: &str = "https://dollyfiuba.com/analitics/comentarios_docentes.json";

pub const CREACION_TABLA_CUATRIMESTRES: &str = r#"
CREATE TABLE IF NOT EXISTS cuatrimestres(
    nombre TEXT PRIMARY KEY
);
"#;

pub const CREACION_TABLA_COMENTARIOS: &str = r#"
CREATE TABLE IF NOT EXISTS comentarios(
    codigo         UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    codigo_docente TEXT REFERENCES docentes(codigo) NOT NULL,
    cuatrimestre   TEXT REFERENCES cuatrimestres(nombre) NOT NULL,
    contenido      TEXT NOT NULL
);
"#;

#[derive(Deserialize, PartialEq, Eq, Hash)]
pub struct Cuatrimestre {
    #[serde(alias = "cuat")]
    pub nombre: String,

    #[serde(alias = "mat")]
    pub codigo_materia: u32,

    #[serde(alias = "doc")]
    pub nombre_docente: String,
}

impl Cuatrimestre {
    pub async fn descargar(
        client: &ClientWithMiddleware,
    ) -> anyhow::Result<HashMap<Self, Vec<String>>> {
        #[derive(Deserialize)]
        struct Payload {
            #[serde(flatten)]
            cuatrimestre: Cuatrimestre,
            comentarios: Vec<Option<String>>,
        }

        tracing::info!("descargando listado de comentarios");
        let res = client.get(URL_DESCARGA).send().await?;
        let data = res.text().await?;

        let cuatrimestres: Vec<Payload> =
            serde_json::from_str(&data).map_err(|err| SerdeError::new(data, err))?;

        Ok(cuatrimestres
            .into_iter()
            .map(|cuatrimestre| {
                let comentarios_decoded: Vec<_> = cuatrimestre
                    .comentarios
                    .into_iter()
                    .flatten()
                    .filter_map(|c| {
                        String::from_utf8(general_purpose::STANDARD.decode(c).ok()?).ok()
                    })
                    .collect();

                (cuatrimestre.cuatrimestre, comentarios_decoded)
            })
            .collect())
    }

    pub fn sql(nombre: &str) -> String {
        format!(
            r#"
INSERT INTO cuatrimestres(nombre)
VALUES ('{}');
        "#,
            nombre.sanitizar()
        )
    }
}

pub struct Comentario;

impl Comentario {
    pub fn query_sql(
        cuatrimestre: &Cuatrimestre,
        codigo_docente: &Uuid,
        comentarios: &[String],
    ) -> String {
        comentarios
            .iter()
            .map(|contenido| {
                format!(
                    r#"
INSERT INTO comentarios(cuatrimestre, codigo_docente, contenido)
VALUES ('{}', '{}', '{}');
"#,
                    cuatrimestre.nombre.sanitizar(),
                    codigo_docente,
                    contenido.sanitizar()
                )
            })
            .collect::<Vec<_>>()
            .join("")
    }
}
