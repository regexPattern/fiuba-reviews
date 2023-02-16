use super::{tablas, CatedraModelo, DocenteModelo, MateriaModelo};

pub fn construir_query(
    materias: Vec<MateriaModelo>,
    catedras: Vec<CatedraModelo>,
    docentes: Vec<DocenteModelo>,
) -> String {
    let mut secciones = vec![
        tablas::TABLA_MATERIAS,
        tablas::TABLA_CATEDRAS,
        tablas::TABLA_DOCENTES,
        tablas::TABLA_COMENTARIOS,
    ];

    secciones.join("\n\n")
}
