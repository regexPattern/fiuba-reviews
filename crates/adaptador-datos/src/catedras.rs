use std::hash::{Hash, Hasher};

use uuid::Uuid;

use crate::{docentes::Docente, sql::Sql};

#[derive(Debug)]
pub struct Catedra {
    pub codigo: Uuid,
    pub nombre: String,
    pub docentes: Vec<Docente>,
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
INSERT INTO catedra(codigo, codigo_materia)
VALUES ('{}', {});
"#,
            self.codigo, codigo_materia,
        )
    }

    pub fn relacion_con_docente_query_sql(&self, docente: &Docente) -> String {
        format!(
            r#"
INSERT INTO catedra_docente(codigo_catedra, codigo_docente)
VALUES ('{}', '{}');
"#,
            self.codigo,
            docente.codigo.sanitizar_sql()
        )
    }
}
