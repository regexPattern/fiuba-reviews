use uuid::Uuid;

use crate::{
    catedra, comentario, docente,
    materia::{self, MateriaScrapeResult},
};

#[derive(Default, Debug)]
pub struct BulkInsertTuples {
    pub materias: Vec<String>,
    pub catedras: Vec<String>,
    pub docentes: Vec<String>,
    pub rel_catedras_docentes: Vec<String>,
    pub calificaciones: Vec<String>,
    pub cuatrimestres: Vec<String>,
    pub comentarios: Vec<String>,
}

pub trait Sql {
    fn sanitize(&self) -> String;
}

impl BulkInsertTuples {
    pub fn extend(&mut self, materia: MateriaScrapeResult) {
        self.catedras.extend(materia.catedras);
        self.docentes.extend(materia.docentes);
        self.rel_catedras_docentes
            .extend(materia.rel_catedras_docentes);
        self.calificaciones.extend(materia.calificaciones);
    }

    pub fn sql(&self) -> String {
        [
            materia::bulk_insert(&self.materias),
            catedra::bulk_insert(&self.catedras),
            docente::bulk_insert_docentes(&self.docentes),
            docente::bulk_insert_rel_catedras_docentes(&self.rel_catedras_docentes),
            docente::bulk_insert_calificaciones(&self.calificaciones),
            comentario::bulk_insert_cuatrimestre(&self.cuatrimestres),
            comentario::bulk_insert_comentarios(&self.comentarios),
        ]
        .join("\n\n")
    }
}

impl Sql for &str {
    fn sanitize(&self) -> String {
        format!("'{}'", self.replace("'", "''"))
    }
}

impl Sql for String {
    fn sanitize(&self) -> String {
        self.as_str().sanitize()
    }
}

impl Sql for Uuid {
    fn sanitize(&self) -> String {
        format!("'{}'", self)
    }
}

impl Sql for Vec<String> {
    fn sanitize(&self) -> String {
        self.join(",\n\t")
    }
}
