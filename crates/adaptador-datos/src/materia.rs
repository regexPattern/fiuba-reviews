use std::{cmp::Ordering, collections::HashMap};

use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;
use uuid::Uuid;

use crate::{
    catedra::{self, Catedra},
    docente::{self, Calificacion},
    sql::Sql,
};

const URL_DESCARGA_EQUIVALENCIAS: &str = "https://raw.githubusercontent.com/lugfi/dolly/7e105810fadd340aa4f89f9ae58160a2fea6e7ae/data/equivalencias.json";
const URL_DESCARGA_MATERIAS: &str = "https://raw.githubusercontent.com/lugfi/dolly/7e105810fadd340aa4f89f9ae58160a2fea6e7ae/data/comun.json";
const URL_DESCARGA_CATEDRAS: &str = "https://dollyfiuba.com/analitics/cursos";

#[serde_with::serde_as]
#[derive(Deserialize, Debug)]
pub struct Materia {
    #[serde_as(as = "serde_with::DisplayFromStr")]
    codigo: i16,

    nombre: String,

    #[serde(skip)]
    codigo_equivalencia: Option<i16>,
}

#[derive(Default, Debug)]
pub struct MateriaScrapeResult {
    pub catedras: Vec<String>,
    pub docentes: Vec<String>,
    pub rel_catedras_docentes: Vec<String>,
    pub calificaciones: Vec<String>,
    pub codigos_docentes: HashMap<(String, i16), Uuid>,
}

pub async fn descargar_todas(cliente_http: &ClientWithMiddleware) -> anyhow::Result<Vec<Materia>> {
    #[derive(Deserialize)]
    struct ResDolly {
        materias: Vec<Materia>,
    }

    tracing::info!("descargando listado de materias");

    let res = cliente_http.get(URL_DESCARGA_MATERIAS).send().await?;
    let data = res.text().await?;

    let ResDolly { mut materias } =
        serde_json::from_str(&data).map_err(|err| SerdeError::new(data, err))?;

    asignas_equivalencias(cliente_http, &mut materias).await?;

    materias.sort_by(
        |a, b| match (a.codigo_equivalencia, b.codigo_equivalencia) {
            (None, Some(_)) => Ordering::Less,
            (Some(_), None) => Ordering::Greater,
            _ => a.codigo.cmp(&b.codigo),
        },
    );

    Ok(materias)
}

pub fn bulk_insert(insert_tuples: &Vec<String>) -> String {
    format!(
        "\
INSERT INTO materia (codigo, nombre, codigo_equivalencia)
VALUES
    {}
ON CONFLICT (codigo)
DO NOTHING;",
        insert_tuples.sanitize()
    )
}

impl Materia {
    pub async fn scape(
        &self,
        cliente_http: &ClientWithMiddleware,
    ) -> anyhow::Result<MateriaScrapeResult> {
        let catedras = self
            .descargar_catedras(&cliente_http)
            .await
            .inspect_err(|err| {
                tracing::error!("error descargando catedras de materia {}", self.codigo);
                tracing::debug!("descripcion error: {err}");
            })?;

        let mut materia = MateriaScrapeResult {
            catedras: Vec::with_capacity(catedras.len()),
            ..Default::default()
        };

        let mut codigos_docentes = HashMap::new();

        for cat in catedras {
            materia.catedras.push(cat.sql(self.codigo));

            for (nombre, calificacion) in cat.docentes {
                let codigo = codigos_docentes.entry(nombre).or_insert_with_key(|nombre| {
                    let codigo = Uuid::new_v4();
                    materia
                        .docentes
                        .push(docente::sql_docente(&codigo, &nombre, self.codigo));

                    codigo
                });

                materia
                    .rel_catedras_docentes
                    .push(docente::sql_rel_catedra_docente(&cat.codigo, &codigo));

                if calificacion.respuestas > 0 {
                    materia
                        .calificaciones
                        .push(docente::sql_calificacion(&calificacion, &codigo));
                }
            }
        }

        materia.codigos_docentes = codigos_docentes
            .into_iter()
            .map(|(n, c)| ((n, self.codigo), c))
            .collect();

        Ok(materia)
    }

    pub fn sql(&self) -> String {
        let codigo = self.codigo;
        let nombre = self.nombre.sanitize();
        let codigo_equivalencia = self
            .codigo_equivalencia
            .map(|c| c.to_string())
            .unwrap_or("NULL".to_string());

        format!("({codigo}, {nombre}, {codigo_equivalencia})",)
    }

    async fn descargar_catedras(
        &self,
        cliente_http: &ClientWithMiddleware,
    ) -> anyhow::Result<Vec<Catedra>> {
        #[derive(Deserialize)]
        struct ResDolly {
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

        let ResDolly { catedras } =
            serde_json::from_str(&data).map_err(|err| SerdeError::new(data, err))?;

        let mut catedras = catedras
            .into_iter()
            .map(|c| Catedra {
                codigo: Uuid::new_v4(),
                docentes: c.docentes,
            })
            .collect();

        catedra::eliminar_repetidas(&mut catedras);

        Ok(catedras)
    }
}

async fn asignas_equivalencias(
    cliente_http: &ClientWithMiddleware,
    materias: &mut [Materia],
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
