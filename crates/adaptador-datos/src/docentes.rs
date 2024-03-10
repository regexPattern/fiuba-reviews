use serde::Deserialize;

use crate::sql::Sql;

#[derive(Deserialize, Default, Debug)]
pub struct Docente {
    pub codigo: String,
    pub nombre: String,
    pub calificacion: Calificacion,
}

#[derive(Deserialize, Default, Debug)]
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

pub fn generar_codigo_docente(codigo_materia: u32, nombre_docente: &str) -> String {
    format!("{}-{}", codigo_materia, nombre_docente)
}

impl Docente {
    pub fn query_sql(&self) -> String {
        let query_docente = format!(
            r#"
INSERT INTO docente(codigo, nombre) VALUES ('{}', '{}');
"#,
            self.codigo,
            self.nombre.sanitizar_sql()
        );

        let acepta_critica = self.calificacion.acepta_critica.unwrap_or_default();
        let asistencia = self.calificacion.asistencia.unwrap_or_default();
        let buen_trato = self.calificacion.buen_trato.unwrap_or_default();
        let claridad = self.calificacion.claridad.unwrap_or_default();
        let clase_organizada = self.calificacion.clase_organizada.unwrap_or_default();
        let cumple_horarios = self.calificacion.cumple_horarios.unwrap_or_default();
        let fomenta_participacion = self.calificacion.fomenta_participacion.unwrap_or_default();
        let panorama_amplio = self.calificacion.panorama_amplio.unwrap_or_default();
        let responde_mails = self.calificacion.responde_mails.unwrap_or_default();

        let mut queries_calificaciones = String::new();

        for _ in 0..self.calificacion.respuestas {
            queries_calificaciones.push_str(&format!(
                r#"
INSERT INTO calificacion(codigo_docente, acepta_critica, asistencia, buen_trato, claridad, clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails)
VALUES ('{}', {}, {}, {}, {}, {}, {}, {}, {}, {});
"#,
            self.codigo.sanitizar_sql(),
            acepta_critica,
            asistencia,
            buen_trato,
            claridad,
            clase_organizada,
            cumple_horarios,
            fomenta_participacion,
            panorama_amplio,
            responde_mails,
        ));
        }

        query_docente + &queries_calificaciones
    }
}
