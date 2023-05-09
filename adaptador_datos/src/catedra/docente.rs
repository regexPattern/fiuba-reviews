use serde::Deserialize;
use uuid::Uuid;

pub const TABLA: &'static str = r#"
CREATE TABLE IF NOT EXISTS docente (
    codigo                TEXT PRIMARY KEY,
    nombre                TEXT NOT NULL,
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

mod pesos_calificaciones {
    pub const ACEPTA_CRITICA: f64 = 0.5;
    pub const ASISTENCIA: f64 = 1.0;
    pub const BUEN_TRATO: f64 = 0.5;
    pub const CLARIDAD: f64 = 0.7;
    pub const CLASE_ORGANIZADA: f64 = 0.7;
    pub const CUMPLE_HORARIOS: f64 = 1.0;
    pub const FOMENTA_PARTICIPACION: f64 = 0.5;
    pub const PANORAMA_AMPLIO: f64 = 0.5;
    pub const RESPONDE_MAILS: f64 = 0.5;
}

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
    respuestas: Option<u32>,
}

impl Docente {
    pub fn promedio_calificaciones(&self) -> f64 {
        let calificaciones_con_peso = [
            self.acepta_critica.unwrap_or_default() * pesos_calificaciones::ACEPTA_CRITICA,
            self.asistencia.unwrap_or_default() * pesos_calificaciones::ASISTENCIA,
            self.buen_trato.unwrap_or_default() * pesos_calificaciones::BUEN_TRATO,
            self.claridad.unwrap_or_default() * pesos_calificaciones::CLARIDAD,
            self.clase_organizada.unwrap_or_default() * pesos_calificaciones::CLASE_ORGANIZADA,
            self.cumple_horarios.unwrap_or_default() * pesos_calificaciones::CUMPLE_HORARIOS,
            self.fomenta_participacion.unwrap_or_default()
                * pesos_calificaciones::FOMENTA_PARTICIPACION,
            self.panorama_amplio.unwrap_or_default() * pesos_calificaciones::PANORAMA_AMPLIO,
            self.responde_mails.unwrap_or_default() * pesos_calificaciones::RESPONDE_MAILS,
        ];

        let cantidad_calificaciones = calificaciones_con_peso.len();

        if cantidad_calificaciones > 0 {
            calificaciones_con_peso.into_iter().sum::<f64>() / cantidad_calificaciones as f64
        } else {
            0.0
        }
    }

    pub fn sql(&self, nombre_docente: &str, codigo_docente: Uuid) -> String {
        format!(
            r#"
INSERT INTO docente (codigo, nombre, respuestas, promedio, acepta_critica, asistencia, buen_trato, claridad, clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails)
VALUES ('{}', '{}', {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {});
"#,
            codigo_docente,
            nombre_docente.replace("'", "''"),
            self.respuestas.unwrap_or(0),
            self.promedio_calificaciones(),
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
