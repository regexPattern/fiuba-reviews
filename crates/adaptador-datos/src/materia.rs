use std::{cmp::Ordering, collections::HashMap};

use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;
use uuid::Uuid;

use crate::{catedra::Catedra, docente::Calificacion};

const URL_DESCARGA_EQUIVALENCIAS: &str = "https://raw.githubusercontent.com/lugfi/dolly/7e105810fadd340aa4f89f9ae58160a2fea6e7ae/data/equivalencias.json";
const URL_DESCARGA_MATERIAS: &str = "https://raw.githubusercontent.com/lugfi/dolly/7e105810fadd340aa4f89f9ae58160a2fea6e7ae/data/comun.json";
const URL_DESCARGA_CATEDRAS: &str = "https://dollyfiuba.com/analitics/cursos";

#[serde_with::serde_as]
#[derive(Deserialize, Debug)]
pub struct Materia {
    #[serde_as(as = "serde_with::DisplayFromStr")]
    pub codigo: i16,

    pub nombre: String,

    #[serde(skip)]
    pub codigo_equivalencia: Option<i16>,
}

impl Materia {
    pub async fn descargar_todas(cliente: &ClientWithMiddleware) -> anyhow::Result<Vec<Self>> {
        #[derive(Deserialize)]
        struct RespuestaDolly {
            materias: Vec<Materia>,
        }

        tracing::info!("descargando listado de materias");

        let res = cliente.get(URL_DESCARGA_MATERIAS).send().await?;
        let data = res.text().await?;

        let RespuestaDolly { mut materias } =
            serde_json::from_str(&data).map_err(|err| SerdeError::new(data, err))?;

        Self::asignas_equivalencias(cliente, &mut materias).await?;

        materias.sort_by(
            |a, b| match (a.codigo_equivalencia, b.codigo_equivalencia) {
                (None, Some(_)) => Ordering::Less,
                (Some(_), None) => Ordering::Greater,
                _ => a.codigo.cmp(&b.codigo),
            },
        );

        Ok(materias)
    }

    pub async fn descargar_catedras(
        &self,
        cliente_http: &ClientWithMiddleware,
    ) -> anyhow::Result<Vec<Catedra>> {
        #[derive(Deserialize)]
        struct RespuestaDolly {
            #[serde(alias = "opciones")]
            catedras: Vec<CatedraDolly>,
        }

        #[derive(Deserialize)]
        struct CatedraDolly {
            docentes: HashMap<String, Calificacion>,
        }

        tracing::info!("descargando catedras de materia {}", self.codigo);

        let res = cliente_http
            .get(format!("{}/{}", URL_DESCARGA_CATEDRAS, self.codigo))
            .send()
            .await?;

        let data = res.text().await?;

        let RespuestaDolly { catedras } =
            serde_json::from_str(&data).map_err(|err| SerdeError::new(data, err))?;

        let catedras = catedras
            .into_iter()
            .map(|c| Catedra {
                codigo: Uuid::new_v4(),
                docentes: c.docentes,
            })
            .collect();

        let _catedars = Catedra::consumir_repetidas(catedras);

        Ok([].into())
    }

    pub fn sql(&self) -> String {
        format!(
            "({}, '{}', {}, {})",
            self.codigo,
            self.nombre,
            self.codigo_equivalencia
                .map(|c| c.to_string())
                .unwrap_or("NULL".into()),
            false
        )
    }

    async fn asignas_equivalencias(
        cliente_http: &ClientWithMiddleware,
        materias: &mut [Self],
    ) -> anyhow::Result<()> {
        tracing::info!("descargando listado de equivalencias");

        let res = cliente_http.get(URL_DESCARGA_EQUIVALENCIAS).send().await?;
        let data = res.text().await?;

        let codigos_equivalencias: Vec<Vec<i16>> =
            serde_json::from_str(&data).map_err(|err| SerdeError::new(data, err))?;

        let mut equivalencias = HashMap::new();

        for codigo in codigos_equivalencias {
            let mut codigos = codigo.into_iter();

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
}
