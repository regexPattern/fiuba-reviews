mod deserializar;

use anyhow::Result;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;
use tracing::Level;
use uuid::Uuid;

use crate::entidad::{AdaptadorMateria, EntidadCatedra, EntidadDocente, EntidadMateria};

const URL_CATEDRAS: &'static str = "https://dollyfiuba.com/analitics/cursos";
const URL_COMENTARIOS: &'static str = "https://dollyfiuba.com/analitics/comentarios_docentes.json";
const URL_MATERIAS: &'static str =
    "https://raw.githubusercontent.com/lugfi/dolly/master/data/comun.json";

#[derive(Deserialize, Debug)]
pub struct MateriaRemoto {
    #[serde(deserialize_with = "self::deserializar::codigo")]
    pub codigo: u32,
    pub nombre: String,
}

#[derive(Debug, PartialEq)]
pub struct CatedraRemoto {
    nombre: String,
    calificaciones: Vec<CalificacionRemoto>,
}

#[derive(Deserialize, Debug, Default, PartialEq)]
pub struct CalificacionRemoto {
    pub nombre: String,
    pub respuestas: usize,

    pub acepta_critica: Option<f64>,
    pub asistencia: Option<f64>,
    pub buen_trato: Option<f64>,
    pub claridad: Option<f64>,
    pub clase_organizada: Option<f64>,
    pub cumple_horarios: Option<f64>,
    pub fomenta_participacion: Option<f64>,
    pub panorama_amplio: Option<f64>,
    pub responde_mails: Option<f64>,
}

#[derive(Deserialize, Debug, PartialEq)]
pub struct ComentarioRemoto {
    #[serde(rename(deserialize = "mat"))]
    pub codigo_materia: u32,

    #[serde(rename(deserialize = "doc"))]
    pub nombre_docente: String,

    #[serde(rename(deserialize = "cuat"))]
    pub cuatrimestre: String,

    #[serde(
        deserialize_with = "self::deserializar::comentarios",
        rename(deserialize = "comentarios")
    )]
    pub contenido_comentarios: Vec<String>,
}

pub async fn descargar_materias(cliente_http: &ClientWithMiddleware) -> Result<Vec<MateriaRemoto>> {
    #[derive(Deserialize)]
    struct Respuesta {
        materias: Vec<MateriaRemoto>,
    }

    tracing::event!(Level::DEBUG, "GET {URL_MATERIAS}");

    let respuesta = cliente_http.get(URL_MATERIAS).send().await?.text().await?;

    let Respuesta { materias } = serde_json::from_str(&respuesta)?;

    Ok(materias)
}

pub async fn descargar_comentarios(
    cliente_http: &ClientWithMiddleware,
) -> Result<Vec<ComentarioRemoto>> {
    let respuesta = cliente_http
        .get(URL_COMENTARIOS)
        .send()
        .await?
        .text()
        .await?;

    Ok(serde_json::from_str(&respuesta)?)
}

impl MateriaRemoto {
    pub async fn generar_adaptador(
        self,
        cliente_http: &ClientWithMiddleware,
    ) -> Result<AdaptadorMateria> {
        let (catedras, docentes) = match self.adaptar_catedras_y_docentes(cliente_http).await {
            Ok(data) => Ok(data),
            Err(err) => {
                tracing::event!(Level::ERROR, "MATERIA {} NO DISPONIBLE", self.codigo);
                Err(err)
            }
        }?;

        let materia = EntidadMateria {
            codigo: self.codigo,
            nombre: self.nombre,
        };

        Ok(AdaptadorMateria {
            entidad_materia: materia,
            entidades_catedras: catedras,
            entidades_docentes: docentes,
        })
    }

    async fn adaptar_catedras_y_docentes(
        &self,
        cliente_http: &ClientWithMiddleware,
    ) -> Result<(Vec<EntidadCatedra>, Vec<EntidadDocente>)> {
        let mut catedras = vec![];
        let mut docentes = vec![];

        for catedra in self.descargar_catedras(cliente_http).await? {
            let codigo_catedra = Uuid::new_v4();

            for calificacion in catedra.calificaciones {
                docentes.push(EntidadDocente {
                    codigo: Uuid::new_v4(),
                    calificacion,
                    codigo_catedra,
                });
            }

            catedras.push(EntidadCatedra {
                codigo: codigo_catedra,
                nombre: catedra.nombre,
                codigo_materia: self.codigo,
            });
        }

        let catedras = catedras.into_iter().collect();

        Ok((catedras, docentes))
    }

    async fn descargar_catedras(
        &self,
        cliente_http: &ClientWithMiddleware,
    ) -> Result<Vec<CatedraRemoto>> {
        #[derive(Deserialize)]
        struct Respuesta {
            #[serde(rename(deserialize = "opciones"))]
            catedras: Vec<CatedraRemoto>,
        }

        let url_catedras_materia = format!("{URL_CATEDRAS}/{}.json", self.codigo);

        tracing::event!(Level::DEBUG, "GET {url_catedras_materia}");

        let respuesta = cliente_http
            .get(url_catedras_materia)
            .send()
            .await?
            .text()
            .await?;

        let Respuesta { catedras } = serde_json::from_str(&respuesta)?;

        Ok(catedras)
    }
}
