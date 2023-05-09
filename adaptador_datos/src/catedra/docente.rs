use serde::Deserialize;
use uuid::Uuid;

pub const TABLA: &'static str = r#"
CREATE TABLE IF NOT EXISTS docente (
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
);"#;

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
        format!(
            r#"INSERT INTO docente (codigo, nombre, acepta_critica, asistencia, buen_trato, claridad, clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails, respuestas)
VALUES ('{}', '{}', {}, {}, {}, {}, {}, {}, {}, {}, {}, {});"#,
            codigo_docente,
            nombre_docente.replace("'", "''"),
            self.acepta_critica
                .map(|v| v.to_string())
                .unwrap_or("NULL".to_string()),
            self.asistencia
                .map(|v| v.to_string())
                .unwrap_or("NULL".to_string()),
            self.buen_trato
                .map(|v| v.to_string())
                .unwrap_or("NULL".to_string()),
            self.claridad
                .map(|v| v.to_string())
                .unwrap_or("NULL".to_string()),
            self.clase_organizada
                .map(|v| v.to_string())
                .unwrap_or("NULL".to_string()),
            self.cumple_horarios
                .map(|v| v.to_string())
                .unwrap_or("NULL".to_string()),
            self.fomenta_participacion
                .map(|v| v.to_string())
                .unwrap_or("NULL".to_string()),
            self.panorama_amplio
                .map(|v| v.to_string())
                .unwrap_or("NULL".to_string()),
            self.responde_mails
                .map(|v| v.to_string())
                .unwrap_or("NULL".to_string()),
            self.respuestas.unwrap_or(0.0)
        )
    }
}
