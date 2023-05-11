use std::{
    collections::HashMap,
    hash::{Hash, Hasher},
};

use serde::Deserialize;
use uuid::Uuid;

pub const CREACION_TABLA_CATEDRAS: &str = r#"
CREATE TABLE IF NOT EXISTS Catedras(
    codigo         TEXT PRIMARY KEY,
    nombre         TEXT NOT NULL,
    codigo_materia INTEGER REFERENCES Materias(codigo) NOT NULL,
    promedio       DOUBLE PRECISION NOT NULL
);
"#;

pub const CREACION_TABLA_DOCENTES: &str = r#"
CREATE TABLE IF NOT EXISTS Docentes(
    -- Datos personales.
    codigo                TEXT PRIMARY KEY,
    nombre                TEXT NOT NULL,

    -- Datos calificacion.
    respuestas            INTEGER NOT NULL,
    promedio              DOUBLE PRECISION NOT NULL,
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
    codigo_catedra TEXT REFERENCES Catedras(codigo),
    codigo_docente TEXT REFERENCES Docentes(codigo),
    CONSTRAINT catedra_docente_pkey PRIMARY KEY (codigo_catedra, codigo_docente)
);
"#;

#[derive(Debug)]
pub struct Catedra {
    pub codigo: Uuid,
    pub nombre: String,
    pub docentes: HashMap<NombreDocente, Calificacion>,
    pub promedio: f64,
}

#[derive(Clone, Deserialize, PartialEq, Eq, Hash, Debug)]
#[serde(transparent)]
pub struct NombreDocente(pub String);

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

impl Catedra {
    pub fn query_sql(&self, codigo_materia: u32) -> String {
        format!(
            r#"
INSERT INTO Catedras(codigo, codigo_materia, nombre, promedio)
VALUES ('{}', {}, '{}', {});
"#,
            self.codigo,
            codigo_materia,
            self.nombre.replace('\'', "''"),
            self.promedio
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
    pub fn promedio(&self) -> f64 {
        let calificaciones = [
            self.acepta_critica.unwrap_or_default(),
            self.asistencia.unwrap_or_default(),
            self.buen_trato.unwrap_or_default(),
            self.claridad.unwrap_or_default(),
            self.clase_organizada.unwrap_or_default(),
            self.cumple_horarios.unwrap_or_default(),
            self.fomenta_participacion.unwrap_or_default(),
            self.panorama_amplio.unwrap_or_default(),
            self.responde_mails.unwrap_or_default(),
        ];

        let cantidad_calificaciones = calificaciones.len();

        if cantidad_calificaciones > 0 {
            calificaciones.into_iter().sum::<f64>() / cantidad_calificaciones as f64
        } else {
            0.0
        }
    }

    pub fn query_sql(&self, nombre_docente: &NombreDocente, codigo_docente: Uuid) -> String {
        format!(
            r#"
INSERT INTO Docentes(codigo, nombre, respuestas, promedio, acepta_critica, asistencia, buen_trato, claridad, clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails)
VALUES ('{}', '{}', {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {});
"#,
            codigo_docente,
            nombre_docente.0.replace('\'', "''"),
            self.respuestas.unwrap_or(0),
            self.promedio(),
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
