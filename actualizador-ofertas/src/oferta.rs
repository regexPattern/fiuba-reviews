use serde::Deserialize;
use std::collections::{HashMap, hash_map::Entry};

#[derive(Deserialize, PartialEq, Debug)]
pub struct OfertaCarrera {
    carrera: String,
    anio: u16,
    cuatrimestre: u8,
    materias: Vec<Materia>,
}

#[derive(Deserialize, PartialEq, Debug)]
struct Materia {
    codigo: String,
    nombre: String,
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

/// Consolida las ofertas de varias carreras en una oferta por materia.
///
/// Para cada código de materia se conserva solo el período más reciente.
/// Si varias carreras ofrecen la misma materia en ese período, se devuelve una
/// sola oferta que lista todas esas carreras y reúne las comisiones de cada una.
pub fn consolidar_ofertas_materias(ofertas_carreras: Vec<OfertaCarrera>) -> Vec<OfertaMateria> {
    let mut por_codigo: HashMap<String, OfertaMateria> = HashMap::new();

    for oferta in ofertas_carreras {
        let periodo = (oferta.anio, oferta.cuatrimestre);
        for materia in oferta.materias {
            match por_codigo.entry(materia.codigo.clone()) {
                Entry::Vacant(e) => {
                    e.insert(OfertaMateria {
                        carreras: vec![oferta.carrera.clone()],
                        anio: oferta.anio,
                        cuatrimestre: oferta.cuatrimestre,
                        materia,
                    });
                }
                Entry::Occupied(mut e) => {
                    let actual = e.get_mut();
                    let periodo_actual = (actual.anio, actual.cuatrimestre);
                    if periodo > periodo_actual {
                        *actual = OfertaMateria {
                            carreras: vec![oferta.carrera.clone()],
                            anio: oferta.anio,
                            cuatrimestre: oferta.cuatrimestre,
                            materia,
                        };
                    } else if periodo == periodo_actual {
                        if !actual.carreras.contains(&oferta.carrera) {
                            actual.carreras.push(oferta.carrera.clone());
                        }
                        for comision in materia.comisiones {
                            if !actual.materia.comisiones.contains(&comision) {
                                actual.materia.comisiones.push(comision);
                            }
                        }
                    }
                }
            }
        }
    }

    por_codigo.into_values().collect()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn ofertas_vacias_producen_un_resultado_vacio() {
        let ofertas_carreras = vec![];
        let ofertas_materias = consolidar_ofertas_materias(ofertas_carreras);
        assert_eq!(ofertas_materias.len(), 0);
    }

    #[test]
    fn ofertas_con_materias_distintas_conservan_todas_las_materias() {
        let ofertas_carreras = serde_json::from_value(serde_json::json!([
            {
                "carrera": "Ingeniería Civil",
                "anio": 2026,
                "cuatrimestre": 1,
                "materias": [{
                    "codigo": "TA073",
                    "nombre": "AEROPUERTOS I",
                    "comisiones": []
                }]
            },
            {
                "carrera": "Ingeniería Industrial",
                "anio": 2026,
                "cuatrimestre": 1,
                "materias": [{
                    "codigo": "TA030",
                    "nombre": "AGROINDUSTRIA Y BIOECONOMÍA",
                    "comisiones": []
                }]
            }
        ]))
        .unwrap();

        let esperadas: Vec<OfertaMateria> = serde_json::from_value(serde_json::json!([
            {
                "carreras": ["Ingeniería Civil"],
                "anio": 2026,
                "cuatrimestre": 1,
                "materia": {
                    "codigo": "TA073",
                    "nombre": "AEROPUERTOS I",
                    "comisiones": []
                }
            },
            {
                "carreras": ["Ingeniería Industrial"],
                "anio": 2026,
                "cuatrimestre": 1,
                "materia": {
                    "codigo": "TA030",
                    "nombre": "AGROINDUSTRIA Y BIOECONOMÍA",
                    "comisiones": []
                }
            }
        ]))
        .unwrap();

        let mut ofertas_materias = consolidar_ofertas_materias(ofertas_carreras);

        ofertas_materias.sort_by(|a, b| b.materia.codigo.cmp(&a.materia.codigo));
        assert_eq!(ofertas_materias, esperadas);
    }

    #[test]
    fn ofertas_con_la_misma_materia_y_periodo_unen_sus_comisiones() {
        let ofertas_carreras = serde_json::from_value(serde_json::json!([
            {
                "carrera": "Ingeniería Civil",
                "anio": 2026,
                "cuatrimestre": 1,
                "materias": [{
                    "codigo": "CB002",
                    "nombre": "ÁLGEBRA LINEAL",
                    "comisiones": [{"nombre": "A", "docentes": [], "clases": []}]
                }]
            },
            {
                "carrera": "Ingeniería Industrial",
                "anio": 2026,
                "cuatrimestre": 1,
                "materias": [{
                    "codigo": "CB002",
                    "nombre": "ÁLGEBRA LINEAL",
                    "comisiones": [{"nombre": "B", "docentes": [], "clases": []}]
                }]
            }
        ]))
        .unwrap();

        let esperadas: Vec<OfertaMateria> = serde_json::from_value(serde_json::json!([{
            "carreras": ["Ingeniería Civil", "Ingeniería Industrial"],
            "anio": 2026,
            "cuatrimestre": 1,
            "materia": {
                "codigo": "CB002",
                "nombre": "ÁLGEBRA LINEAL",
                "comisiones": [
                    {"nombre": "A", "docentes": [], "clases": []},
                    {"nombre": "B", "docentes": [], "clases": []}
                ]
            }
        }]))
        .unwrap();

        let mut ofertas_materias = consolidar_ofertas_materias(ofertas_carreras);

        assert_eq!(ofertas_materias.len(), 1);
        ofertas_materias[0].carreras.sort();
        ofertas_materias[0]
            .materia
            .comisiones
            .sort_by(|a, b| a.nombre.cmp(&b.nombre));
        assert_eq!(ofertas_materias, esperadas);
    }

    #[test]
    fn ofertas_con_la_misma_materia_y_distinto_periodo_conservan_la_mas_reciente() {
        let ofertas_carreras = serde_json::from_value(serde_json::json!([
            {
                "carrera": "Ingeniería Civil",
                "anio": 2025,
                "cuatrimestre": 2,
                "materias": [{
                    "codigo": "CB002",
                    "nombre": "ÁLGEBRA LINEAL",
                    "comisiones": [{"nombre": "Old", "docentes": [], "clases": []}]
                }]
            },
            {
                "carrera": "Ingeniería Civil",
                "anio": 2026,
                "cuatrimestre": 2,
                "materias": [{
                    "codigo": "CB002",
                    "nombre": "ÁLGEBRA LINEAL",
                    "comisiones": [{"nombre": "Latest", "docentes": [], "clases": []}]
                }]
            },
            {
                "carrera": "Ingeniería Civil",
                "anio": 2026,
                "cuatrimestre": 1,
                "materias": [{
                    "codigo": "CB002",
                    "nombre": "ÁLGEBRA LINEAL",
                    "comisiones": [{"nombre": "Previous", "docentes": [], "clases": []}]
                }]
            }
        ]))
        .unwrap();

        let esperadas: Vec<OfertaMateria> = serde_json::from_value(serde_json::json!([{
            "carreras": ["Ingeniería Civil"],
            "anio": 2026,
            "cuatrimestre": 2,
            "materia": {
                "codigo": "CB002",
                "nombre": "ÁLGEBRA LINEAL",
                "comisiones": [{"nombre": "Latest", "docentes": [], "clases": []}]
            }
        }]))
        .unwrap();

        let ofertas_materias = consolidar_ofertas_materias(ofertas_carreras);

        assert_eq!(ofertas_materias, esperadas);
    }
}
