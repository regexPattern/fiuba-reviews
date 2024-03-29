use crate::materia::ResIndexadoMateria;

#[derive(Default, Debug)]
pub struct InsertTuplesBuffer {
    pub materias: Vec<String>,
    pub catedras: Vec<String>,
    pub docentes: Vec<String>,
    pub rel_catedras_docentes: Vec<String>,
    pub calificaciones: Vec<String>,
    pub cuatrimestres: Vec<String>,
    pub comentarios: Vec<String>,
}

impl InsertTuplesBuffer {
    pub fn extend(&mut self, materia: ResIndexadoMateria) {
        self.catedras.extend(materia.catedras);
        self.docentes.extend(materia.docentes);
        self.rel_catedras_docentes
            .extend(materia.rel_catedras_docentes);
        self.calificaciones.extend(materia.calificaciones);
    }

    pub fn sql(&self) -> String {
        [
            String::from_utf8_lossy(include_bytes!("../sql/schema.sql")).to_string(),
            self.materias
                .bulk_insert("materia", "(codigo, nombre, codigo_equivalencia)"),
            self.catedras
                .bulk_insert("catedra", "(codigo, codigo_materia)"),
            self.docentes
                .bulk_insert("docente", "(codigo, nombre, codigo_materia)"),
            self.rel_catedras_docentes
                .bulk_insert("catedra_docente", "(codigo_catedra, codigo_docente)"),
            self.calificaciones.bulk_insert("calificacion", "(codigo_docente, acepta_critica, asistencia, buen_trato, claridad, clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails)"),
            self.cuatrimestres.bulk_insert("cuatrimestre", "(nombre)"),
            self.comentarios.bulk_insert("comentario", "(codigo, codigo_docente, cuatrimestre, contenido)")
        ]
        .join("\n")
    }
}

trait BulkInsertable {
    fn bulk_insert(&self, nombre_tabla: &str, nombres_columnas: &str) -> String;
}

impl BulkInsertable for Vec<String> {
    fn bulk_insert(&self, tabla: &str, columnas: &str) -> String {
        format!(
            "
INSERT INTO {} {}
VALUES
\t{};
            ",
            tabla,
            columnas,
            self.join(",\n\t")
        )
    }
}
