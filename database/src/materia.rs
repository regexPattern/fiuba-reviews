use std::collections::HashMap;

use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;
use serde_with::{serde_as, DisplayFromStr};

const URL_DESCARGA_MATERIAS: &'static str =
    "https://raw.githubusercontent.com/lugfi/dolly/master/data/comun.json";
const URL_DESCARGA_EQUIVALENCIAS: &'static str = "https://raw.githubusercontent.com/lugfi/dolly/f47f553a89dc7c7cbf8192277c9f2e3e1e826bf0/data/equivalencias.json";

pub const TABLA: &'static str = "\
CREATE TABLE IF NOT EXISTS materias (
    codigo INTEGER PRIMARY KEY,
    nombre TEXT NOT NULL,
    codigo_equivalencia INTEGER
);";

#[serde_as]
#[derive(Deserialize, Debug)]
pub struct Materia {
    #[serde_as(as = "DisplayFromStr")]
    pub codigo: u32,
    nombre: String,

    #[serde(skip)]
    codigo_equivalencia: Option<u32>,
}

impl Materia {
    pub async fn descargar(cliente: &ClientWithMiddleware) -> anyhow::Result<Vec<Self>> {
        #[derive(Deserialize)]
        struct Respuesta {
            materias: Vec<Materia>,
        }

        tracing::info!("descargando listado de materias");

        let respuesta = cliente
            .get(URL_DESCARGA_MATERIAS)
            .send()
            .await?
            .text()
            .await?;

        let Respuesta { mut materias } =
            serde_json::from_str(&respuesta).map_err(|err| SerdeError::new(respuesta, err))?;

        Self::asignas_equivalencias(cliente, &mut materias).await?;

        Ok(materias)
    }

    async fn asignas_equivalencias(
        cliente: &ClientWithMiddleware,
        materias: &mut [Self],
    ) -> anyhow::Result<()> {
        tracing::info!("descargando listado de equivalencias");

        let respuesta = cliente.get(URL_DESCARGA_EQUIVALENCIAS).send().await?;
        let codigos_equivalencias: Vec<Vec<u32>> = respuesta.json().await?;

        let mut equivalencias = HashMap::new();
        for codigos in codigos_equivalencias {
            let mut codigos = codigos.into_iter();

            let codigo_materia_principal = match codigos.next() {
                Some(codigo) => codigo,
                None => continue,
            };

            for codigo in codigos {
                equivalencias.insert(codigo, codigo_materia_principal);
            }
        }

        tracing::info!("asignando equivalencias a materias");

        for materia in materias {
            materia.codigo_equivalencia = equivalencias.get(&materia.codigo).cloned();
        }

        Ok(())
    }

    pub fn sql(&self) -> String {
        let mut columnas = "codigo, nombre".to_string();
        let mut valores = format!("{}, '{}'", self.codigo, self.nombre.replace("'", "''"));

        if let Some(codigo_equivalencia) = self.codigo_equivalencia {
            columnas.push_str(", codigo_equivalencia");
            valores.push_str(&format!(", {codigo_equivalencia}"));
        }

        format!("INSERT INTO materias ({columnas}) VALUES ({valores});")
    }
}
