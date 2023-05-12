use std::{cmp::Ordering, collections::HashMap};

use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;

use crate::sql::Sql;

pub const CREACION_TABLA: &str = r#"
CREATE TABLE IF NOT EXISTS Materia(
    codigo              INTEGER PRIMARY KEY,
    nombre              TEXT NOT NULL,
    codigo_equivalencia INTEGER REFERENCES Materia(codigo)
);
"#;

const URL_DESCARGA_EQUIVALENCIAS: &str = "https://raw.githubusercontent.com/lugfi/dolly/f47f553a89dc7c7cbf8192277c9f2e3e1e826bf0/data/equivalencias.json";
const URL_DESCARGA_MATERIAS: &str =
    "https://raw.githubusercontent.com/lugfi/dolly/master/data/comun.json";

#[serde_with::serde_as]
#[derive(Deserialize, Debug)]
pub struct Materia {
    #[serde_as(as = "serde_with::DisplayFromStr")]
    pub codigo: u32,

    nombre: String,

    #[serde(skip)]
    codigo_equivalencia: Option<u32>,
}

impl Materia {
    pub async fn descargar(http: &ClientWithMiddleware) -> anyhow::Result<Vec<Self>> {
        #[derive(Deserialize)]
        struct Materias {
            materias: Vec<Materia>,
        }

        tracing::info!("descargando listado de materias");
        let res = http.get(URL_DESCARGA_MATERIAS).send().await?;
        let data = res.text().await?;

        let Materias { mut materias } =
            serde_json::from_str(&data).map_err(|err| SerdeError::new(data, err))?;

        Self::asignas_equivalencias(http, &mut materias).await?;

        // Ordenamos las materias para que las que no tienen equivalencia se terminen insertando
        // antes en la query SQL que las que si tienen equivalencias. Esto porque hay una relacion
        // entre las tablas de materias.
        materias.sort_by(
            |a, b| match (a.codigo_equivalencia, b.codigo_equivalencia) {
                (None, Some(_)) => Ordering::Less,
                (Some(_), None) => Ordering::Greater,
                _ => a.codigo.cmp(&b.codigo),
            },
        );

        Ok(materias)
    }

    async fn asignas_equivalencias(
        http: &ClientWithMiddleware,
        materias: &mut [Self],
    ) -> anyhow::Result<()> {
        tracing::info!("descargando listado de equivalencias");

        let res = http.get(URL_DESCARGA_EQUIVALENCIAS).send().await?;
        let data = res.text().await?;

        let codigos_equivalencias: Vec<Vec<u32>> =
            serde_json::from_str(&data).map_err(|err| SerdeError::new(data, err))?;

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
            materia.codigo_equivalencia = equivalencias.remove(&materia.codigo);
        }

        Ok(())
    }

    pub fn query_sql(&self) -> String {
        format!(
            r#"
INSERT INTO Materia(codigo, nombre, codigo_equivalencia)
VALUES ({}, '{}', {});
"#,
            self.codigo,
            self.nombre.sanitizar(),
            self.codigo_equivalencia
                .map(|v| v.to_string())
                .unwrap_or("NULL".to_string())
        )
    }
}
