use std::{collections::HashMap, hash::Hash};

use base64::{engine::general_purpose, Engine};
use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::{Deserialize, Serialize};
use uuid::Uuid;

use crate::sql::Sql;

const URL_DESCARGA: &str = "https://dollyfiuba.com/analitics/comentarios_docentes.json";
const URL_HF_MODELO_SUMARIZACION: &str =
    "https://api-inference.huggingface.co/models/facebook/bart-large-cnn";

pub const CREACION_TABLA_CUATRIMESTRES: &str = r#"
CREATE TABLE IF NOT EXISTS Cuatrimestre(
    nombre TEXT PRIMARY KEY
);
"#;

pub const CREACION_TABLA_COMENTARIOS: &str = r#"
CREATE TABLE IF NOT EXISTS Comentario(
    codigo         UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    codigo_docente TEXT REFERENCES Docente(codigo) NOT NULL,
    cuatrimestre   TEXT REFERENCES Cuatrimestre(nombre) NOT NULL,
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
INSERT INTO Cuatrimestre(nombre)
VALUES ('{}');
        "#,
            nombre.sanitizar()
        )
    }
}

pub struct Comentario;

impl Comentario {
    pub fn sql(
        cuatrimestre: &Cuatrimestre,
        codigo_docente: &Uuid,
        comentarios: &[String],
    ) -> String {
        comentarios
            .iter()
            .map(|contenido| {
                format!(
                    r#"
INSERT INTO Comentario(cuatrimestre, codigo_docente, contenido)
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

    pub async fn sql_descripcion_ia(
        http: &ClientWithMiddleware,
        api_key: &str,
        codigo_docente: Uuid,
        comentarios: &Vec<String>,
    ) -> anyhow::Result<Option<String>> {
        #[derive(Serialize)]
        struct Inputs {
            inputs: String,
        }

        #[derive(Deserialize)]
        struct Resumen {
            #[serde(alias = "summary_text")]
            descripcion: String,
        }

        if comentarios.is_empty() {
            return Ok(None);
        }

        tracing::info!("generando descripcion de docente {codigo_docente}");

        let res = http
            .post(URL_HF_MODELO_SUMARIZACION)
            .bearer_auth(api_key)
            .json(&Inputs {
                inputs: comentarios.join("."),
            })
            .send()
            .await?;

        let data = res.text().await?;

        let mut resumenes: Vec<Resumen> =
            serde_json::from_str(&data).map_err(|err| SerdeError::new(data, err))?;

        let Resumen { descripcion } = resumenes.remove(0);

        let query = format!(
            r#"
UPDATE Docente
SET descripcion = '{}', comentarios_ultima_descripcion = {}
WHERE codigo = '{codigo_docente}';
"#,
            descripcion.sanitizar(),
            comentarios.len()
        );

        Ok(Some(query))
    }
}
