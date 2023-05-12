use std::{
    collections::HashMap,
    hash::{Hash, Hasher},
};

use reqwest_middleware::ClientWithMiddleware;
use serde::Deserialize;
use uuid::Uuid;

use crate::{materias::Materia, sql::Sql};

pub const CREACION_TABLA_CATEDRAS: &str = r#"
CREATE TABLE IF NOT EXISTS Catedra(
    codigo         TEXT PRIMARY KEY,
    codigo_materia INTEGER REFERENCES Materia(codigo) NOT NULL
);
"#;

pub const CREACION_TABLA_DOCENTES: &str = r#"
CREATE TABLE IF NOT EXISTS Docente(
    -- Datos personales.
    codigo                TEXT PRIMARY KEY,
    nombre                TEXT NOT NULL,

    -- Datos calificacion.
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
);
"#;

pub const CREACION_TABLA_CATEDRA_DOCENTE: &str = r#"
CREATE TABLE IF NOT EXISTS CatedraDocente(
    codigo_catedra TEXT REFERENCES Catedra(codigo),
    codigo_docente TEXT REFERENCES Docente(codigo),
    CONSTRAINT catedra_docente_pkey PRIMARY KEY (codigo_catedra, codigo_docente)
);
"#;

const URL_DESCARGA_CATEDRAS: &str = "https://dollyfiuba.com/analitics/cursos";

#[derive(Debug)]
pub struct Catedra {
    pub codigo: Uuid,
    pub nombre: String,
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
    respuestas: Option<u32>,
}

impl PartialEq for Catedra {
    fn eq(&self, other: &Self) -> bool {
        self.nombre == other.nombre
    }
}

impl Eq for Catedra {}

impl Hash for Catedra {
    fn hash<H: Hasher>(&self, state: &mut H) {
        self.nombre.hash(state);
    }
}

impl Materia {
    pub async fn catedras(&self, http: &ClientWithMiddleware) -> anyhow::Result<Vec<Catedra>> {
        #[derive(Deserialize)]
        struct Catedras {
            #[serde(alias = "opciones")]
            catedras: Vec<CatedraDolly>,
        }

        #[derive(Deserialize)]
        struct CatedraDolly {
            pub nombre: String,
            pub docentes: HashMap<String, Calificacion>,
        }

        tracing::info!("descargando catedras de materia {}", self.codigo);

        let res = http
            .get(format!("{}/{}", URL_DESCARGA_CATEDRAS, self.codigo))
            .send()
            .await?;

        let Catedras { mut catedras } = res.json().await?;

        for catedra in &mut catedras {
            let mut nombres_docentes: Vec<_> = catedra.nombre.split('-').collect();
            nombres_docentes.sort();
            catedra.nombre = nombres_docentes.join("-").to_uppercase();
        }

        let catedras = catedras.into_iter().map(|catedra| Catedra {
            codigo: Uuid::new_v4(),
            nombre: catedra.nombre,
            docentes: catedra.docentes,
        });

        Ok(catedras.collect())
    }
}

impl Catedra {
    pub fn query_sql(&self, codigo_materia: u32) -> String {
        format!(
            r#"
INSERT INTO Catedra(codigo, codigo_materia)
VALUES ('{}', {});
"#,
            self.codigo, codigo_materia,
        )
    }

    pub fn relacion_con_docente_query_sql(&self, codigo_docente: &Uuid) -> String {
        format!(
            r#"
INSERT INTO CatedraDocente(codigo_catedra, codigo_docente)
VALUES ('{}', '{}');
"#,
            self.codigo, codigo_docente
        )
    }
}

impl Calificacion {
    pub fn query_sql(&self, nombre_docente: &str, codigo_docente: Uuid) -> String {
        format!(
            r#"
INSERT INTO Docente(codigo, nombre, respuestas, acepta_critica, asistencia, buen_trato, claridad, clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails)
VALUES ('{}', '{}', {}, {}, {}, {}, {}, {}, {}, {}, {}, {});
"#,
            codigo_docente,
            nombre_docente.sanitizar(),
            self.respuestas.unwrap_or(0),
            self.acepta_critica.map_or("NULL".into(), |v| v.to_string()),
            self.asistencia.map_or("NULL".into(), |v| v.to_string()),
            self.buen_trato.map_or("NULL".into(), |v| v.to_string()),
            self.claridad.map_or("NULL".into(), |v| v.to_string()),
            self.clase_organizada
                .map_or("NULL".into(), |v| v.to_string()),
            self.cumple_horarios
                .map_or("NULL".into(), |v| v.to_string()),
            self.fomenta_participacion
                .map_or("NULL".into(), |v| v.to_string()),
            self.panorama_amplio
                .map_or("NULL".into(), |v| v.to_string()),
            self.responde_mails.map_or("NULL".into(), |v| v.to_string()),
        )
    }
}
