use std::{collections::HashMap, hash::Hash};

use base64::{engine::general_purpose, Engine};
use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;
use uuid::Uuid;

const URL_DESCARGA: &'static str = "https://dollyfiuba.com/analitics/comentarios_docentes.json";

pub const TABLA: &'static str = r#"
CREATE TABLE IF NOT EXISTS comentario (
    codigo         TEXT PRIMARY KEY,
    codigo_docente TEXT REFERENCES docente(codigo) NOT NULL,
    cuatrimestre   TEXT NOT NULL,
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
    pub async fn descargar_comentarios(
        client: &ClientWithMiddleware,
    ) -> anyhow::Result<HashMap<Cuatrimestre, Vec<String>>> {
        #[derive(Deserialize)]
        struct ComentariosCuatrimestre {
            #[serde(flatten)]
            cuatrimestre: Cuatrimestre,
            comentarios: Vec<Option<String>>,
        }

        tracing::info!("descargando listado de comentarios");

        let respuesta = client.get(URL_DESCARGA).send().await?.text().await?;
        let cuatrimestres: Vec<ComentariosCuatrimestre> =
            serde_json::from_str(&respuesta).map_err(|err| SerdeError::new(respuesta, err))?;

        Ok(cuatrimestres
            .into_iter()
            .map(|cuatrimestre| {
                let comentarios_decoded: Vec<_> = cuatrimestre
                    .comentarios
                    .into_iter()
                    .filter_map(|c| c)
                    .filter_map(|c| {
                        String::from_utf8(general_purpose::STANDARD.decode(c).ok()?).ok()
                    })
                    .collect();

                (cuatrimestre.cuatrimestre, comentarios_decoded)
            })
            .collect())
    }

    pub fn sql(&self, codigo_docente: &Uuid, comentarios: &[String]) -> String {
        comentarios
            .iter()
            .map(|contenido| {
                format!(
                    r#"
INSERT INTO comentario (codigo, codigo_docente, cuatrimestre, contenido)
VALUES ('{}', '{}', '{}', '{}');
"#,
                    Uuid::new_v4(),
                    codigo_docente,
                    self.nombre,
                    contenido.replace("'", "''")
                )
            })
            .collect::<Vec<_>>()
            .join("")
    }
}
