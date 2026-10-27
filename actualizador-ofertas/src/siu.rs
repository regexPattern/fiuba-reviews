mod consolidar;
mod normalizar;
mod sanitizar;

pub use consolidar::consolidar_ofertas_materias;
pub use sanitizar::sanitizar_ofertas_carreras;

use serde::Deserialize;
use std::collections::BTreeSet;

#[derive(Debug, Deserialize, PartialEq)]
pub struct OfertaCarrera {
    carrera: String,
    anio: u16,
    cuatrimestre: u8,
    pub materias: Vec<Materia>,
}

#[derive(Debug, Deserialize, PartialEq)]
pub struct Materia {
    pub codigo: String,
    pub nombre: String,
    comisiones: Vec<Comision>,
}

#[derive(Debug, Deserialize, PartialEq)]
struct Comision {
    nombre: String,
    docentes: Vec<Docente>,
    clases: Vec<Clase>,
}

#[derive(Debug, Deserialize, PartialEq)]
struct Docente {
    nombre: String,
    rol: Option<String>,
}

#[derive(Debug, Deserialize, PartialEq)]
struct Clase {
    tipo: String,
    dia: String,
    inicio: String,
    fin: String,
}

#[derive(Debug, Deserialize, PartialEq)]
pub struct OfertaMateria {
    carreras: Vec<String>,
    anio: u16,
    cuatrimestre: u8,
    materia: Materia,
}

impl OfertaMateria {
    pub fn codigo(&self) -> &str {
        &self.materia.codigo
    }

    /// Nombres únicos y no vacíos de los docentes de todas las comisiones, en orden alfabético.
    pub fn nombres_docentes(&self) -> BTreeSet<&str> {
        self.materia
            .comisiones
            .iter()
            .flat_map(|c| &c.docentes)
            .map(|d| d.nombre.trim())
            .filter(|nombre| !nombre.is_empty())
            .collect()
    }
}
