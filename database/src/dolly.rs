use std::{
    collections::{HashMap, HashSet},
    hash::Hash,
};

use base64::{engine::general_purpose, Engine};
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;

const CATEDRAS_BASE_URL: &'static str = "https://dollyfiuba.com/analitics/cursos";
const COMENTARIOS_URL: &'static str = "https://dollyfiuba.com/analitics/comentarios_docentes.json";
const MATERIAS_URL: &'static str =
    "https://raw.githubusercontent.com/lugfi/dolly/master/data/comun.json";

#[derive(Deserialize, Debug)]
pub struct ComentariosCuatrimestre {
    cuatrimestre: String,
    comentarios: Vec<String>,
}

impl ComentariosCuatrimestre {
    pub async fn fetch_all(
        client: &ClientWithMiddleware,
    ) -> anyhow::Result<HashMap<(String, String), Self>> {
        #[derive(Deserialize, Debug)]
        struct ResponseComentario {
            mat: u32,
            doc: String,
            cuat: String,
            comentarios: Vec<Option<String>>,
        }

        tracing::info!("descargando listado de comentarios");
        let response = client.get(COMENTARIOS_URL).send().await?;
        let metadata: Vec<ResponseComentario> = response.json().await?;

        let comentarios = metadata.into_iter().map(|metadata| {
            let decoded_comentarios = metadata
                .comentarios
                .into_iter()
                .filter_map(|c| c)
                .filter_map(|c| String::from_utf8(general_purpose::STANDARD.decode(c).ok()?).ok())
                .collect();

            (
                (metadata.mat.to_string(), metadata.doc),
                Self {
                    cuatrimestre: metadata.cuat,
                    comentarios: decoded_comentarios,
                },
            )
        });

        Ok(comentarios.collect())
    }
}

#[derive(Deserialize, Debug)]
pub struct Materia {
    pub codigo: String,
    nombre: String,

    #[serde(skip)]
    pub catedras: Vec<Catedra>,
}

#[derive(Deserialize, Debug)]
pub struct Catedra {
    docentes: HashMap<String, Docente>,
}

#[derive(Deserialize, Default, Debug)]
pub struct Docente {
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
    pub async fn fetch_all(
        client: &ClientWithMiddleware,
    ) -> anyhow::Result<impl Iterator<Item = Self>> {
        #[derive(Deserialize)]
        struct Response {
            materias: Vec<Materia>,
        }

        tracing::info!("descargando listado de materias");
        let response = client.get(MATERIAS_URL).send().await?;
        let Response { mut materias, .. } = response.json().await?;

        // TODO: remove this in production.
        for materia in materias.iter_mut().take(3) {
            match materia.fetch_catedras(client).await {
                Ok(catedras) => materia.catedras = catedras.collect(),
                Err(err) => {
                    tracing::error!(
                        "error descargando catedras de materia {}: {err}",
                        materia.codigo
                    )
                }
            }
        }

        Ok(materias.into_iter())
    }

    pub async fn fetch_catedras(
        &self,
        client: &ClientWithMiddleware,
    ) -> anyhow::Result<impl Iterator<Item = Catedra>> {
        #[derive(Deserialize)]
        struct Response {
            #[serde(alias = "opciones")]
            catedras: Vec<Catedra>,
        }

        tracing::info!("descargando catedras de materia {}", self.codigo);
        let response = client
            .get(format!("{CATEDRAS_BASE_URL}/{}", self.codigo))
            .send()
            .await?;
        let Response { catedras } = response.json().await?;

        Ok(catedras.into_iter())
    }

    pub fn table_query() -> String {
        "\
CREATE TABLE IF NOT EXISTS materias (
    codigo INTEGER PRIMARY KEY,
    nombre TEXT NOT NULL
);"
        .into()
    }

    pub fn insert_query(&self) -> String {
        format!(
            "INSERT INTO materias (codigo, nombre) VALUES ({}, {});",
            self.codigo, self.nombre
        )
    }
}

impl Catedra {
    pub fn table_query() -> String {
        "\
CREATE TABLE IF NOT EXISTS catedras (
    nombre         TEXT NOT NULL,
    codigo_materia INTEGER REFERENCES materias(codigo) NOT NULL
);"
        .into()
    }

    pub fn insert_query(&self, codigo_materia: &str) -> String {
        format!(
            "INSERT INTO catedras (nombre, codigo_materia) VALUES ({}, {})",
            "nombre_catedra", codigo_materia
        )
    }
}

impl Docente {
    pub fn table_query() -> String {
        "\
CREATE TABLE IF NOT EXISTS docentes (
    nombre                TEXT NOT NULL,
    respuestas            INTEGER NOT NULL,
    acepta_critica        DOUBLE PRECISION,
    asistencia            DOUBLE PRECISION,
    buen_trato            DOUBLE PRECISION,
    claridad              DOUBLE PRECISION,
    clase_organizada      DOUBLE PRECISION,
    cumple_horarios       DOUBLE PRECISION,
    fomenta_participacion DOUBLE PRECISION,
    panorama_amplio       DOUBLE PRECISION,
    responde_mails        DOUBLE PRECISION,
    codigo_catedra        UUID REFERENCES catedras(codigo) NOT NULL
);"
        .into()
    }
}
