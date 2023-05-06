use std::collections::HashMap;

use base64::{engine::general_purpose, Engine};
use format_serde_error::SerdeError;
use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;
use serde_with::{serde_as, DisplayFromStr};
use uuid::Uuid;

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
    pub acepta_critica: Option<f64>,
    pub asistencia: Option<f64>,
    pub buen_trato: Option<f64>,
    pub claridad: Option<f64>,
    pub clase_organizada: Option<f64>,
    pub cumple_horarios: Option<f64>,
    pub fomenta_participacion: Option<f64>,
    pub panorama_amplio: Option<f64>,
    pub responde_mails: Option<f64>,
    pub respuestas: Option<f64>,
}

impl Materia {
    pub const CREATE_TABLE: &'static str = "\
CREATE TABLE IF NOT EXISTS materias (
    codigo INTEGER PRIMARY KEY,
    nombre TEXT NOT NULL
);";

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

    pub fn sql(&self) -> String {
        format!(
            "INSERT INTO materias (codigo, nombre) VALUES ({}, '{}');",
            self.codigo,
            self.nombre.replace("'", "''")
        )
    }
}

impl Catedra {
    pub const CREATE_TABLE: &'static str = "\
CREATE TABLE IF NOT EXISTS catedras (
    codigo         UUID PRIMARY KEY,
    codigo_materia INTEGER REFERENCES materias(codigo) NOT NULL
);";

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

    pub fn sql(codigo_catedra: Uuid, codigo_materia: u32) -> String {
        format!("INSERT INTO catedras (codigo, codigo_materia) VALUES ('{codigo_catedra}', {codigo_materia});")
    }
}

impl Calificacion {
    pub const CREATE_TABLE: &'static str = "\
CREATE TABLE IF NOT EXISTS docentes (
    codigo                UUID PRIMARY KEY,
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
    responde_mails        DOUBLE PRECISION
);";

    pub fn sql(&self, codigo_docente: &Uuid, nombre_docente: String) -> String {
        let acepta_critica = self.acepta_critica.unwrap_or(0.0);
        let asistencia = self.asistencia.unwrap_or(0.0);
        let buen_trato = self.buen_trato.unwrap_or(0.0);
        let claridad = self.claridad.unwrap_or(0.0);
        let clase_organizada = self.clase_organizada.unwrap_or(0.0);
        let cumple_horarios = self.cumple_horarios.unwrap_or(0.0);
        let fomenta_participacion = self.fomenta_participacion.unwrap_or(0.0);
        let panorama_amplio = self.panorama_amplio.unwrap_or(0.0);
        let responde_mails = self.responde_mails.unwrap_or(0.0);
        let respuestas = self.respuestas.unwrap_or(0.0);

        format!("INSERT INTO docentes (codigo, nombre, acepta_critica, asistencia, buen_trato, claridad, clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails, respuestas) \
VALUES ('{codigo_docente}', '{}', {acepta_critica}, {asistencia}, {buen_trato}, {claridad}, {clase_organizada}, {cumple_horarios}, {fomenta_participacion}, {panorama_amplio}, {responde_mails}, {respuestas});", nombre_docente.replace("'", "''"))
    }
}

#[derive(Deserialize, Debug)]
pub struct ComentariosDocentePorCuatri {
    pub cuatrimestre: String,
    pub entradas: Vec<String>,
}

impl ComentariosDocentePorCuatri {
    pub const CREATE_TABLE: &'static str = "\
CREATE TABLE IF NOT EXISTS comentarios (
    codigo         UUID PRIMARY KEY,
    codigo_docente UUID REFERENCES docentes(codigo) NOT NULL,
    cuatrimestre   TEXT NOT NULL,
    contenido      TEXT NOT NULL
);";

    const URL: &'static str = "https://dollyfiuba.com/analitics/comentarios_docentes.json";

    pub async fn fetch_all(
        client: &ClientWithMiddleware,
    ) -> anyhow::Result<HashMap<(u32, String), Self>> {
        #[derive(Deserialize, Debug)]
        struct ResponseComentario {
            mat: u32,
            doc: String,
            cuat: String,
            comentarios: Vec<Option<String>>,
        }

        tracing::info!("descargando listado de comentarios");
        let response = client.get(Self::URL).send().await?.text().await?;
        let metadata: Vec<ResponseComentario> =
            serde_json::from_str(&response).map_err(|e| SerdeError::new(response, e))?;

        let comentarios = metadata.into_iter().map(|metadata| {
            let decoded_comentarios = metadata
                .comentarios
                .into_iter()
                .filter_map(|c| c)
                .filter_map(|c| String::from_utf8(general_purpose::STANDARD.decode(c).ok()?).ok())
                .collect();

            (
                (metadata.mat, metadata.doc),
                Self {
                    cuatrimestre: metadata.cuat,
                    entradas: decoded_comentarios,
                },
            )
        });

        Ok(comentarios.collect())
    }

    pub fn sql(&self, codigo_docente: &Uuid) -> String {
        let mut buffer = vec![];

        for e in &self.entradas {
            buffer.push(format!("INSERT INTO comentarios (codigo, codigo_docente, cuatrimestre, contenido) VALUES ('{}', '{codigo_docente}', '{}', '{}');",
                Uuid::new_v4(), self.cuatrimestre, e.replace("'", "''")));
        }

        buffer.join("\n")
    }
}
