use std::{
    cmp::Ordering,
    collections::{HashMap, HashSet},
};

use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;
use uuid::Uuid;

use crate::{
    catedras::Catedra,
    docentes::{self, Calificacion, Docente},
    sql::Sql,
};

const URL_DESCARGA_EQUIVALENCIAS: &str = "https://raw.githubusercontent.com/lugfi/dolly/f47f553a89dc7c7cbf8192277c9f2e3e1e826bf0/data/equivalencias.json";
const URL_DESCARGA_MATERIAS: &str =
    "https://raw.githubusercontent.com/lugfi/dolly/master/data/comun.json";
const URL_DESCARGA_CATEDRAS: &str = "https://dollyfiuba.com/analitics/cursos";

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
    pub async fn descargar_todas(
        cliente_http: &ClientWithMiddleware,
    ) -> anyhow::Result<impl Iterator<Item = Self>> {
        #[derive(Deserialize)]
        struct RespuestaDolly {
            materias: Vec<Materia>,
        }

        tracing::info!("descargando listado de materias");
        let res = cliente_http.get(URL_DESCARGA_MATERIAS).send().await?;
        let data = res.text().await?;

        let RespuestaDolly { mut materias } =
            serde_json::from_str(&data).map_err(|err| SerdeError::new(data, err))?;

        Self::asignas_equivalencias(cliente_http, &mut materias).await?;

        materias.sort_by(
            |a, b| match (a.codigo_equivalencia, b.codigo_equivalencia) {
                (None, Some(_)) => Ordering::Less,
                (Some(_), None) => Ordering::Greater,
                _ => a.codigo.cmp(&b.codigo),
            },
        );

        Ok(materias.into_iter())
    }

    async fn asignas_equivalencias(
        cliente_http: &ClientWithMiddleware,
        materias: &mut [Self],
    ) -> anyhow::Result<()> {
        tracing::info!("descargando listado de equivalencias");

        let res = cliente_http.get(URL_DESCARGA_EQUIVALENCIAS).send().await?;
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

    pub async fn descargar_catedras(
        &self,
        cliente_http: &ClientWithMiddleware,
    ) -> anyhow::Result<impl Iterator<Item = Catedra>> {
        #[derive(Deserialize)]
        struct RespuestaDolly {
            #[serde(alias = "opciones")]
            catedras: Vec<CatedraDolly>,
        }

        #[derive(Deserialize)]
        struct CatedraDolly {
            pub nombre: String,
            pub docentes: HashMap<String, Calificacion>,
        }

        tracing::info!("descargando catedras de materia {}", self.codigo);

        let res = cliente_http
            .get(format!("{}/{}", URL_DESCARGA_CATEDRAS, self.codigo))
            .send()
            .await?;

        let data = res.text().await?;

        let RespuestaDolly { mut catedras } =
            serde_json::from_str(&data).map_err(|err| SerdeError::new(data, err))?;

        for catedra in &mut catedras {
            let mut nombres_docentes: Vec<_> = catedra.nombre.split('-').collect();
            nombres_docentes.sort();
            catedra.nombre = nombres_docentes.join("-").to_uppercase();
        }

        let catedras: HashSet<_> = catedras
            .into_iter()
            .map(|catedra| {
                let docentes = catedra
                    .docentes
                    .into_iter()
                    .map(|(nombre, calificacion)| Docente {
                        codigo: docentes::generar_codigo_docente(self.codigo, &nombre),
                        nombre,
                        calificacion,
                    });

                Catedra {
                    codigo: Uuid::new_v4(),
                    nombre: catedra.nombre,
                    docentes: docentes.collect(),
                }
            })
            .collect();

        Ok(catedras.into_iter())
    }

    pub fn query_sql(&self) -> String {
        format!(
            r#"
INSERT INTO materia(codigo, nombre, codigo_equivalencia)
VALUES ({}, '{}', {});
"#,
            self.codigo,
            self.nombre.sanitizar_sql(),
            self.codigo_equivalencia
                .map(|cod| cod.to_string())
                .unwrap_or("NULL".to_string())
        )
    }
}
