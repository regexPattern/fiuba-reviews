pub mod queries;

use sqlx::types::Uuid;

pub trait Sql {
    fn sql(&self) -> String;
}

impl Sql for Uuid {
    fn sql(&self) -> String {
        format!("'{}'", self)
    }
}

impl Sql for String {
    fn sql(&self) -> String {
        format!("'{}'", self.replace("'", "''"))
    }
}

impl Sql for usize {
    fn sql(&self) -> String {
        self.to_string()
    }
}

impl Sql for (Uuid, String, usize) {
    fn sql(&self) -> String {
        let (codigo_docente, resumen_comentarios, comentarios_ultimo_resumen) = self;
        format!(
            "({}, {}, {})",
            codigo_docente.sql(),
            resumen_comentarios.sql(),
            comentarios_ultimo_resumen.sql()
        )
    }
}
