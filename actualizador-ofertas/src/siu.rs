mod consolidar;
mod normalizar;
mod sanitizar;

pub use consolidar::consolidar_ofertas_materias;
pub use sanitizar::sanitizar_ofertas_carreras;

use serde::Deserialize;

#[derive(Deserialize, PartialEq, Debug)]
pub struct OfertaCarrera {
    carrera: String,
    anio: u16,
    cuatrimestre: u8,
    pub materias: Vec<Materia>,
}

#[derive(Deserialize, PartialEq, Debug)]
pub struct Materia {
    pub codigo: String,
    pub nombre: String,
    comisiones: Vec<Comision>,
}

#[derive(Deserialize, PartialEq, Debug)]
struct Comision {
    nombre: String,
    docentes: Vec<Docente>,
    clases: Vec<Clase>,
}

#[derive(Deserialize, PartialEq, Debug)]
struct Docente {
    nombre: String,
    rol: Option<String>,
}

#[derive(Deserialize, PartialEq, Debug)]
struct Clase {
    tipo: String,
    dia: String,
    inicio: String,
    fin: String,
}

#[derive(Deserialize, PartialEq, Debug)]
pub struct OfertaMateria {
    carreras: Vec<String>,
    anio: u16,
    cuatrimestre: u8,
    materia: Materia,
}
