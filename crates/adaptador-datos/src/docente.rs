use serde::Deserialize;
use uuid::Uuid;

use crate::sql::Sql;

#[derive(Deserialize, Debug)]
#[cfg_attr(test, derive(Default, Clone, PartialEq))]
pub struct Calificacion {
    pub respuestas: usize,
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

pub fn sql_docente(codigo_docente: &Uuid, nombre_docente: &str, codigo_materia: i16) -> String {
    let codigo = codigo_docente.sanitize();
    let nombre = nombre_docente.sanitize();
    format!("({codigo}, {nombre}, {codigo_materia})")
}

pub fn sql_rel_catedra_docente(codigo_catedra: &Uuid, codigo_docente: &Uuid) -> String {
    let codigo_catedra = codigo_catedra.sanitize();
    let codigo_docente = codigo_docente.sanitize();
    format!("({codigo_catedra}, {codigo_docente})")
}

pub fn sql_calificacion(calificacion: &Calificacion, codigo_docente: &Uuid) -> String {
    assert!(calificacion.respuestas > 0);

    let codigo_docente = codigo_docente.sanitize();
    let mut calificaciones = Vec::with_capacity(calificacion.respuestas);

    for _ in 0..calificacion.respuestas {
        let acepta_critica = calificacion.acepta_critica.unwrap_or_default();
        let asistencia = calificacion.asistencia.unwrap_or_default();
        let buen_trato = calificacion.buen_trato.unwrap_or_default();
        let claridad = calificacion.claridad.unwrap_or_default();
        let clase_organizada = calificacion.clase_organizada.unwrap_or_default();
        let cumple_horarios = calificacion.cumple_horarios.unwrap_or_default();
        let fomenta_participacion = calificacion.fomenta_participacion.unwrap_or_default();
        let panorama_amplio = calificacion.panorama_amplio.unwrap_or_default();
        let responde_mails = calificacion.responde_mails.unwrap_or_default();

        calificaciones.push(format!(
            "({codigo_docente}::uuid, \
            {acepta_critica}, \
            {asistencia}, \
            {buen_trato}, \
            {claridad}, \
            {clase_organizada}, \
            {cumple_horarios}, \
            {fomenta_participacion}, \
            {panorama_amplio}, \
            {responde_mails})"
        ))
    }

    calificaciones.sanitize()
}

pub fn bulk_insert_docentes(insert_tuples: &Vec<String>) -> String {
    format!(
        "\
INSERT INTO docente (codigo, nombre, codigo_materia)
VALUES
    {}
ON CONFLICT (nombre, codigo_materia)
DO NOTHING;",
        insert_tuples.sanitize()
    )
}

pub fn bulk_insert_rel_catedras_docentes(insert_tuples: &Vec<String>) -> String {
    format!(
        "\
INSERT INTO catedra_docente (codigo_catedra, codigo_docente)
VALUES
    {};",
        insert_tuples.sanitize()
    )
}

pub fn bulk_insert_calificaciones(insert_tuples: &Vec<String>) -> String {
    format!(
        "\
WITH temp_calificacion(codigo_docente) AS (
    VALUES
        {}
)
INSERT INTO calificacion
SELECT gen_random_uuid (), t.*
FROM temp_calificacion t
WHERE NOT EXISTS (
    SELECT 1
    FROM calificacion c
    WHERE c.codigo_docente = t.codigo_docente
);",
        insert_tuples.sanitize().replace('\n', "\n    ")
    )
}
