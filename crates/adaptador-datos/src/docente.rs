use serde::Deserialize;

#[derive(Deserialize)]
#[cfg_attr(test, derive(Default, Clone, PartialEq, Debug))]
pub struct Calificacion {
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
