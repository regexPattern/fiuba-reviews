use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;
use serde_with::{serde_as, DisplayFromStr};

const URL_DESCARGA: &'static str =
    "https://raw.githubusercontent.com/lugfi/dolly/master/data/comun.json";

pub const TABLA: &'static str = "\
CREATE TABLE IF NOT EXISTS materias (
    codigo INTEGER PRIMARY KEY,
    nombre TEXT NOT NULL
);";

#[serde_as]
#[derive(Deserialize)]
pub struct Materia {
    #[serde_as(as = "DisplayFromStr")]
    pub codigo: u32,
    pub nombre: String,
}

impl Materia {
    pub async fn descargar(cliente: &ClientWithMiddleware) -> anyhow::Result<Vec<Self>> {
        #[derive(Deserialize)]
        struct Respuesta {
            materias: Vec<Materia>,
        }

        tracing::info!("descargando listado de materias");

        let respuesta = cliente.get(URL_DESCARGA).send().await?.text().await?;
        let Respuesta { materias } =
            serde_json::from_str(&respuesta).map_err(|err| SerdeError::new(respuesta, err))?;

        Ok(materias)
    }

    pub fn sql(&self) -> String {
        format!(
            "INSERT INTO materias (codigo, nombre) VALUES ({}, '{}');",
            self.codigo,
            self.nombre.replace("'", "''")
        )
    }
}
