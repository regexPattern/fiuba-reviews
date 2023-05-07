pub mod sql;

use std::collections::HashMap;

use base64::{engine::general_purpose, Engine};
use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;
use serde_with::{serde_as, DisplayFromStr};

#[serde_as]
#[derive(Deserialize, Debug)]
pub struct Materia {
    #[serde_as(as = "DisplayFromStr")]
    pub codigo: u32,
    pub nombre: String,
}

#[derive(Deserialize, Debug)]
pub struct Catedra {
    pub docentes: HashMap<String, Calificacion>,
}

#[derive(Deserialize, Default, Debug)]
pub struct Calificacion {
    acepta_critica: Option<f64>,
    asistencia: Option<f64>,
    buen_trato: Option<f64>,
    claridad: Option<f64>,
    clase_organizada: Option<f64>,
    cumple_horarios: Option<f64>,
    fomenta_participacion: Option<f64>,
    panorama_amplio: Option<f64>,
    responde_mails: Option<f64>,
    respuestas: Option<f64>,
}

impl Materia {
    const URL: &'static str =
        "https://raw.githubusercontent.com/lugfi/dolly/master/data/comun.json";

    pub async fn fetch_all(
        client: &ClientWithMiddleware,
    ) -> anyhow::Result<impl Iterator<Item = Self>> {
        #[derive(Deserialize)]
        struct Response {
            materias: Vec<Materia>,
        }

        tracing::info!("descargando listado de materias");
        let response = client.get(Self::URL).send().await?.text().await?;
        let Response { materias, .. } =
            serde_json::from_str(&response).map_err(|err| SerdeError::new(response, err))?;

        Ok(materias.into_iter())
    }
}

impl Catedra {
    const URL: &'static str = "https://dollyfiuba.com/analitics/cursos";

    pub async fn fetch_for_materia(
        client: &ClientWithMiddleware,
        codigo_materia: u32,
    ) -> anyhow::Result<impl Iterator<Item = Catedra>> {
        #[derive(Deserialize)]
        struct Response {
            #[serde(alias = "opciones")]
            catedras: Vec<Catedra>,
        }

        tracing::info!("descargando catedras de materia {}", codigo_materia);
        let response = client
            .get(format!("{}/{}", Self::URL, codigo_materia))
            .send()
            .await?
            .text()
            .await?;

        let Response { catedras } =
            serde_json::from_str(&response).map_err(|e| SerdeError::new(response, e))?;

        Ok(catedras.into_iter())
    }
}

#[derive(Debug, PartialEq, Eq, Hash)]
pub struct Comentarios {
    pub codigo_materia: u32,
    pub nombre_docente: String,
    pub nombre_cuatrimestre: String,
}

impl Comentarios {
    const URL: &'static str = "https://dollyfiuba.com/analitics/comentarios_docentes.json";

    pub async fn fetch_all(
        client: &ClientWithMiddleware,
    ) -> anyhow::Result<HashMap<Comentarios, Vec<String>>> {
        #[derive(Deserialize, Debug)]
        struct Cuatrimestre {
            #[serde(alias = "mat")]
            codigo_materia: u32,

            #[serde(alias = "doc")]
            nombre_docente: String,

            #[serde(alias = "cuat")]
            nombre_cuatrimestre: String,

            comentarios: Vec<Option<String>>,
        }

        tracing::info!("descargando listado de comentarios");

        let response = client.get(Self::URL).send().await?.text().await?;
        let cuatrimestres: Vec<Cuatrimestre> =
            serde_json::from_str(&response).map_err(|e| SerdeError::new(response, e))?;

        let mut comentarios = HashMap::new();
        for cuatrimestre in cuatrimestres {
            let decoded_entries: Vec<_> = cuatrimestre
                .comentarios
                .into_iter()
                .filter_map(|c| c)
                .filter_map(|c| String::from_utf8(general_purpose::STANDARD.decode(c).ok()?).ok())
                .collect();

            comentarios.insert(
                Comentarios {
                    nombre_cuatrimestre: cuatrimestre.nombre_cuatrimestre,
                    nombre_docente: cuatrimestre.nombre_docente,
                    codigo_materia: cuatrimestre.codigo_materia,
                },
                decoded_entries,
            );
        }

        Ok(comentarios)
    }
}
