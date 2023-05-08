use std::collections::HashMap;

use serde::Deserialize;
use uuid::Uuid;

pub const TABLA: &'static str = "\
CREATE TABLE IF NOT EXISTS docentes (
    codigo                TEXT PRIMARY KEY,
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

#[derive(Deserialize, Default)]
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

impl Docente {
    pub fn sql(&self, nombre_docente: &str, codigo_docente: Uuid) -> String {
        let calificaciones = HashMap::from([
            ("acepta_critica", self.acepta_critica),
            ("asistencia", self.asistencia),
            ("buen_trato", self.buen_trato),
            ("claridad", self.claridad),
            ("clase_organizada", self.clase_organizada),
            ("cumple_horarios", self.cumple_horarios),
            ("fomenta_participacion", self.fomenta_participacion),
            ("panorama_amplio", self.panorama_amplio),
            ("responde_mails", self.responde_mails),
            ("respuestas", self.respuestas),
        ]);

        let mut columnas = Vec::with_capacity(10);
        let mut valores = Vec::with_capacity(10);

        for (columna, valor) in calificaciones {
            if let Some(value) = valor {
                columnas.push(columna);
                valores.push(value.to_string());
            }
        }

        format!(
            "INSERT INTO docentes (codigo, nombre, {}) VALUES ('{}', '{}', {});",
            columnas.join(", "),
            codigo_docente,
            nombre_docente.replace("'", "''"),
            valores.join(", ")
        )
    }
}
