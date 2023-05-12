use std::{
    cmp::Ordering,
    collections::{HashMap, HashSet},
};

use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;
use uuid::Uuid;

use crate::catedras::{Calificacion, Catedra, NombreDocente};

const URL_DESCARGA_CATEDRAS: &str = "https://dollyfiuba.com/analitics/cursos";
const URL_DESCARGA_EQUIVALENCIAS: &str = "https://raw.githubusercontent.com/lugfi/dolly/f47f553a89dc7c7cbf8192277c9f2e3e1e826bf0/data/equivalencias.json";
const URL_DESCARGA_MATERIAS: &str =
    "https://raw.githubusercontent.com/lugfi/dolly/master/data/comun.json";

pub const CREACION_TABLA: &str = r#"
CREATE TABLE IF NOT EXISTS Materia(
    codigo              INTEGER PRIMARY KEY,
    nombre              TEXT NOT NULL,
    codigo_equivalencia INTEGER REFERENCES Materia(codigo)
);
"#;

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
    pub async fn descargar(cliente_http: &ClientWithMiddleware) -> anyhow::Result<Vec<Self>> {
        #[derive(Deserialize)]
        struct PayloadDolly {
            materias: Vec<Materia>,
        }

        tracing::info!("descargando listado de materias");

        let res = cliente_http
            .get(URL_DESCARGA_MATERIAS)
            .send()
            .await?
            .text()
            .await?;

        let PayloadDolly { mut materias } =
            serde_json::from_str(&res).map_err(|err| SerdeError::new(res, err))?;

        Self::asignas_equivalencias(cliente_http, &mut materias).await?;
        materias.sort_by(
            |a, b| match (a.codigo_equivalencia, b.codigo_equivalencia) {
                (None, Some(_)) => Ordering::Less,
                (Some(_), None) => Ordering::Greater,
                _ => a.codigo.cmp(&b.codigo),
            },
        );

        Ok(materias)
    }

    pub async fn catedras(
        &self,
        cliente_http: &ClientWithMiddleware,
    ) -> anyhow::Result<HashSet<Catedra>> {
        #[derive(Deserialize)]
        struct PayloadDolly {
            #[serde(alias = "opciones")]
            catedras: Vec<PayloadCatedra>,
        }

        #[derive(Deserialize)]
        struct PayloadCatedra {
            pub nombre: String,
            pub docentes: HashMap<NombreDocente, Calificacion>,
        }

        tracing::info!("descargando catedras de materia {}", self.codigo);

        let res = cliente_http
            .get(format!("{}/{}", URL_DESCARGA_CATEDRAS, self.codigo))
            .send()
            .await?
            .text()
            .await?;

        let PayloadDolly { mut catedras } =
            serde_json::from_str(&res).map_err(|err| SerdeError::new(res, err))?;

        for catedra in &mut catedras {
            let mut nombres_docentes: Vec<_> = catedra.nombre.split('-').collect();
            nombres_docentes.sort();
            catedra.nombre = nombres_docentes.join("-");
        }

        let catedras = catedras.into_iter().map(|catedra| {
            let acumulado: f64 = catedra
                .docentes
                .values()
                .map(|docente| docente.promedio())
                .sum();

            let cantidad_docentes = catedra.docentes.len();
            let promedio = if cantidad_docentes > 0 {
                acumulado / cantidad_docentes as f64
            } else {
                0.0
            };

            Catedra {
                codigo: Uuid::new_v4(),
                nombre: catedra.nombre,
                docentes: catedra.docentes,
                promedio,
            }
        });

        Ok(catedras.collect())
    }

    async fn asignas_equivalencias(
        cliente_http: &ClientWithMiddleware,
        materias: &mut [Self],
    ) -> anyhow::Result<()> {
        tracing::info!("descargando listado de equivalencias");

        let res = cliente_http.get(URL_DESCARGA_EQUIVALENCIAS).send().await?;
        let codigos_equivalencias: Vec<Vec<u32>> = res.json().await?;

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

    pub fn query_sql(&self) -> String {
        format!(
            r#"
INSERT INTO Materia(codigo, nombre, codigo_equivalencia)
VALUES ({}, '{}', {});
"#,
            self.codigo,
            self.nombre.replace('\'', "''"),
            self.codigo_equivalencia
                .map(|v| v.to_string())
                .unwrap_or("NULL".to_string())
        )
    }
}
