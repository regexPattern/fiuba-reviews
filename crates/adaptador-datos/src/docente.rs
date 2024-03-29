use serde::Deserialize;
use uuid::Uuid;

#[derive(Deserialize, Debug)]
#[cfg_attr(test, derive(Default, Clone, PartialEq))]
pub struct Calificacion {
    respuestas: usize,
    acepta_critica: Option<f64>,
    asistencia: Option<f64>,
    buen_trato: Option<f64>,
    claridad: Option<f64>,
    clase_organizada: Option<f64>,
    cumple_horarios: Option<f64>,
    fomenta_participacion: Option<f64>,
    panorama_amplio: Option<f64>,
    responde_mails: Option<f64>,
}

pub fn sql(codigo_docente: &Uuid, nombre_docente: &str, codigo_materia: i16) -> String {
    format!(
        "('{}', '{}', {})",
        codigo_docente,
        nombre_docente.replace("'", "''"),
        codigo_materia,
    )
}

pub fn sql_rel_catedra_docente(codigo_catedra: &Uuid, codigo_docente: &Uuid) -> String {
    format!("('{}', '{}')", codigo_catedra, codigo_docente)
}

impl Calificacion {
    pub fn sql(&self, codigo_docente: &Uuid) -> String {
        let buffer: Vec<()> = Vec::new();

        for _ in 0..self.respuestas {
            // TODO: insertar cada calificacion como individual...
        }

        format!(
            "('{}', {}, {}, {}, {}, {}, {}, {}, {}, {})",
            codigo_docente,
            self.acepta_critica.unwrap_or_default(),
            self.asistencia.unwrap_or_default(),
            self.buen_trato.unwrap_or_default(),
            self.claridad.unwrap_or_default(),
            self.clase_organizada.unwrap_or_default(),
            self.cumple_horarios.unwrap_or_default(),
            self.fomenta_participacion.unwrap_or_default(),
            self.panorama_amplio.unwrap_or_default(),
            self.responde_mails.unwrap_or_default(),
        )
    }
}
