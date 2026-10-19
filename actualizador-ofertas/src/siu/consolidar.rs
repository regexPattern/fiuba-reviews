use super::{OfertaCarrera, OfertaMateria};
use std::collections::{HashMap, hash_map::Entry};
use tracing::debug;

/// Consolida las ofertas de varias carreras en una oferta por materia.
///
/// Para cada código de materia se conserva el período más reciente con comisiones.
/// Las materias sin comisiones se ignoran antes de comparar períodos.
/// Si varias carreras ofrecen la misma materia en ese período con comisiones, se devuelve
/// una sola oferta que lista esas carreras y reúne las comisiones de cada una.
pub fn consolidar_ofertas_materias(ofertas_carreras: Vec<OfertaCarrera>) -> Vec<OfertaMateria> {
    let mut por_codigo: HashMap<String, OfertaMateria> = HashMap::new();

    for oferta in ofertas_carreras {
        let periodo = (oferta.anio, oferta.cuatrimestre);
        for materia in oferta.materias {
            if materia.comisiones.is_empty() {
                continue;
            }

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
                        debug!(
                            materia = %actual.materia.codigo,
                            carrera = %oferta.carrera,
                            periodo = ?(oferta.cuatrimestre, oferta.anio),
                            periodo_anterior = ?(periodo_actual.1, periodo_actual.0),
                            "oferta sobreescrita por ser de un cuatrimestre más reciente"
                        );
                        *actual = OfertaMateria {
                            carreras: vec![oferta.carrera.clone()],
                            anio: oferta.anio,
                            cuatrimestre: oferta.cuatrimestre,
                            materia,
                        };
                    } else if periodo < periodo_actual {
                        debug!(
                            materia = %actual.materia.codigo,
                            carrera = %oferta.carrera,
                            periodo = ?(oferta.cuatrimestre, oferta.anio),
                            periodo_vigente = ?(periodo_actual.1, periodo_actual.0),
                            "oferta descartada por ser de un cuatrimestre más antiguo"
                        );
                    } else {
                        let cantidad_anterior = actual.materia.comisiones.len();
                        if !actual.carreras.contains(&oferta.carrera) {
                            actual.carreras.push(oferta.carrera.clone());
                        }
                        for comision in materia.comisiones {
                            if !actual.materia.comisiones.contains(&comision) {
                                actual.materia.comisiones.push(comision);
                            }
                        }
                        let comisiones_agregadas =
                            actual.materia.comisiones.len() - cantidad_anterior;
                        if comisiones_agregadas > 0 {
                            debug!(
                                materia = %actual.materia.codigo,
                                carrera = %oferta.carrera,
                                periodo = ?(oferta.cuatrimestre, oferta.anio),
                                comisiones_agregadas,
                                "agregado comisiones a la oferta existente"
                            );
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
    use tracing_test::traced_test;

    #[test]
    fn las_ofertas_con_materias_distintas_conservan_las_materias_con_comisiones() {
        let ofertas_carreras = serde_json::from_value(serde_json::json!([
            {
                "carrera": "Ingeniería Civil",
                "anio": 2026,
                "cuatrimestre": 1,
                "materias": [{
                    "codigo": "TA073",
                    "nombre": "AEROPUERTOS I",
                    "comisiones": [{"nombre": "A", "docentes": [], "clases": []}]
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

        let esperadas: Vec<_> = serde_json::from_value(serde_json::json!([
            {
                "carreras": ["Ingeniería Civil"],
                "anio": 2026,
                "cuatrimestre": 1,
                "materia": {
                    "codigo": "TA073",
                    "nombre": "AEROPUERTOS I",
                    "comisiones": [{"nombre": "A", "docentes": [], "clases": []}]
                }
            }
        ]))
        .unwrap();

        let mut ofertas_materias = consolidar_ofertas_materias(ofertas_carreras);

        ofertas_materias.sort_by(|a, b| b.materia.codigo.cmp(&a.materia.codigo));
        assert_eq!(ofertas_materias, esperadas);
    }

    #[test]
    fn una_oferta_reciente_vacia_conserva_la_anterior_con_comisiones() {
        // Misma materia pero en 2025 con comisión A y 2026 vacía, en ambos órdenes de llegada
        // (primero se procesa 2025 seguido de 2026, y luego el caso inverso también).
        for periodos in [[2025, 2026], [2026, 2025]] {
            let ofertas: Vec<_> = periodos
                .into_iter()
                .map(|anio| {
                    serde_json::from_value(serde_json::json!({
                        "carrera": "Ingeniería Civil",
                        "anio": anio,
                        "cuatrimestre": 1,
                        "materias": [{
                            "codigo": "CB002",
                            "nombre": "Álgebra Lineal",
                            "comisiones": if anio == 2025 {
                                serde_json::json!([{"nombre": "A", "docentes": [], "clases": []}])
                            } else {
                                serde_json::json!([])
                            }
                        }]
                    }))
                    .unwrap()
                })
                .collect();

            let consolidadas = consolidar_ofertas_materias(ofertas);

            assert_eq!(consolidadas.len(), 1);
            assert_eq!(consolidadas[0].anio, 2025);
            assert_eq!(consolidadas[0].cuatrimestre, 1);
            assert_eq!(consolidadas[0].materia.comisiones.len(), 1);
            assert_eq!(consolidadas[0].materia.comisiones[0].nombre, "A");
        }
    }

    #[test]
    #[traced_test]
    fn las_ofertas_con_la_misma_materia_y_el_mismo_periodo_unen_sus_comisiones() {
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

        let esperadas: Vec<_> = serde_json::from_value(serde_json::json!([{
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
        assert!(logs_contain("agregado comisiones a la oferta existente"));
        assert!(logs_contain("materia=CB002"));
        assert!(logs_contain("carrera=Ingeniería Industrial"));
        assert!(logs_contain("periodo=(1, 2026)"));
        assert!(logs_contain("comisiones_agregadas=1"));
    }

    #[test]
    #[traced_test]
    fn las_ofertas_con_la_misma_materia_y_distinto_periodo_conservan_la_mas_reciente() {
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

        let esperadas: Vec<_> = serde_json::from_value(serde_json::json!([{
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
        assert!(logs_contain(
            "oferta sobreescrita por ser de un cuatrimestre más reciente"
        ));
        assert!(logs_contain(
            "oferta descartada por ser de un cuatrimestre más antiguo"
        ));
        assert!(logs_contain("materia=CB002"));
        assert!(logs_contain("carrera=Ingeniería Civil"));
        assert!(logs_contain("periodo=(2, 2026)"));
        assert!(logs_contain("periodo_anterior=(2, 2025)"));
        assert!(logs_contain("periodo=(1, 2026)"));
        assert!(logs_contain("periodo_vigente=(2, 2026)"));
    }
}
