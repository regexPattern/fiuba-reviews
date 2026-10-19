use serde::Deserialize;
use std::collections::{HashMap, hash_map::Entry};
use tracing::{debug, warn};

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

pub fn sanitizar_ofertas_carreras(mut ofertas_carreras: Vec<OfertaCarrera>) -> Vec<OfertaCarrera> {
    for oferta in &mut ofertas_carreras {
        for materia in &mut oferta.materias {
            sanitizar_comisiones_materia(materia);
        }
    }
    ofertas_carreras
}

fn sanitizar_comisiones_materia(materia: &mut Materia) {
    if materia.comisiones.is_empty() {
        warn!(
            materia = %materia.codigo,
            nombre = %materia.nombre,
            "materia sin comisiones"
        );
        return;
    }

    let codigo = &materia.codigo;
    materia.comisiones.retain(|comision| {
        if comision
            .nombre
            .trim()
            .to_uppercase()
            .contains("CONDICIONAL")
        {
            debug!(materia = %codigo, comision = %comision.nombre, "comisión condicional excluida");
            return false;
        }
        if comision.docentes.is_empty() {
            debug!(materia = %codigo, comision = %comision.nombre, "comisión sin docentes excluida");
            return false;
        }
        true
    });

    for comision in &mut materia.comisiones {
        for docente in &mut comision.docentes {
            if let Some(rol) = &docente.rol {
                docente.rol = Some(normalizar_rol(rol));
            }
        }
    }
}

/// Devuelve la forma canónica de un rol de docente.
///
/// El SIU escribe el mismo rol de distintas formas (por ejemplo `JTP` y
/// `Jefe/a Trabajos Practicos`, o `Ayudante 1ra` y `Ayudante 1ro/a`). Se unifican para que
/// el rol pueda compararse y ordenarse sin tratar variantes del mismo rol como distintas.
/// También se quita el sufijo ` ad-honorem`. Los roles sin variantes conocidas, incluidos los desconocidos, se mantienen tal cual.
fn normalizar_rol(rol: &str) -> String {
    let rol = rol.strip_suffix("ad-honorem").unwrap_or(rol).trim();
    match rol {
        "JTP" => "Jefe/a Trabajos Practicos",
        "Titular" => "Profesor/a Titular/a",
        "Ayudante 1ra" | "Ayudante de 1ra" => "Ayudante 1ro/a",
        "Ayudante 2da" | "Ayudante de 2da" => "Ayudante 2do/a",
        otro => otro,
    }
    .to_string()
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
    mod sanitizar {
        mod sanitizar_comisiones_materia {
            use super::super::super::*;
            use tracing_test::traced_test;

            #[test]
            #[traced_test]
            fn se_filtran_las_materias_sin_comisiones() {
                let mut materia: Materia = serde_json::from_value(serde_json::json!({
                    "codigo": "CB002",
                    "nombre": "Álgebra Lineal",
                    "comisiones": []
                }))
                .unwrap();

                sanitizar_comisiones_materia(&mut materia);

                assert!(materia.comisiones.is_empty());
                assert!(logs_contain("materia sin comisiones"));
                assert!(logs_contain("CB002"));
                assert!(logs_contain("Álgebra Lineal"));
            }

            #[test]
            #[traced_test]
            fn se_filtran_las_comisiones_sin_docentes() {
                let mut materia: Materia = serde_json::from_value(serde_json::json!({
                    "codigo": "CB002",
                    "nombre": "Álgebra Lineal",
                    "comisiones": [
                        {"nombre": "A", "docentes": [], "clases": []},
                        {"nombre": "B", "docentes": [{"nombre": "VARGAS GABRIELA", "rol": "Profesor/a Adjunto/a"}], "clases": []}
                    ]
                }))
                .unwrap();

                sanitizar_comisiones_materia(&mut materia);

                assert_eq!(materia.comisiones.len(), 1);
                assert_eq!(materia.comisiones[0].nombre, "B");
                assert!(logs_contain("comisión sin docentes excluida"));
                assert!(logs_contain("CB002"));
            }

            #[test]
            #[traced_test]
            fn se_filtran_las_comisiones_condicionales() {
                let docente = serde_json::json!([{"nombre": "VARGAS GABRIELA", "rol": "Profesor/a Adjunto/a"}]);
                let mut materia: Materia = serde_json::from_value(serde_json::json!({
                    "codigo": "CB002",
                    "nombre": "Álgebra Lineal",
                    "comisiones": [
                        {"nombre": "CONDICIONALES", "docentes": docente, "clases": []},
                        {"nombre": "CONDICIONALES - GRYNBERG/ORECCHIA", "docentes": docente, "clases": []},
                        {"nombre": "Curso Conde Comisión 1", "docentes": docente, "clases": []},
                        {"nombre": "A", "docentes": docente, "clases": []}
                    ]
                }))
                .unwrap();

                sanitizar_comisiones_materia(&mut materia);

                assert_eq!(materia.comisiones.len(), 2);
                assert_eq!(materia.comisiones[0].nombre, "Curso Conde Comisión 1");
                assert_eq!(materia.comisiones[1].nombre, "A");
                assert!(logs_contain("comisión condicional excluida"));
                assert!(logs_contain("CB002"));
                assert!(logs_contain("CONDICIONALES - GRYNBERG/ORECCHIA"));
            }

            #[test]
            #[traced_test]
            #[ignore = "pendiente: warn de comisión incompleta"]
            fn se_conservan_las_comisiones_con_docentes_de_nombre_vacio() {
                let mut materia: Materia = serde_json::from_value(serde_json::json!({
                    "codigo": "CB002",
                    "nombre": "Álgebra Lineal",
                    "comisiones": [{
                        "nombre": "A",
                        "docentes": [
                            {"nombre": "VARGAS GABRIELA", "rol": "Profesor/a Adjunto/a"},
                            {"nombre": "", "rol": "Ayudante 1ro/a"}
                        ],
                        "clases": []
                    }]
                }))
                .unwrap();

                sanitizar_comisiones_materia(&mut materia);

                assert_eq!(materia.comisiones.len(), 1);
                assert_eq!(materia.comisiones[0].nombre, "A");
                assert_eq!(materia.comisiones[0].docentes.len(), 2);
                assert!(logs_contain("comisión incompleta"));
                assert!(logs_contain("CB002"));
            }

            #[test]
            fn se_normalizan_los_roles_de_los_docentes() {
                let mut materia: Materia = serde_json::from_value(serde_json::json!({
                    "codigo": "CB002",
                    "nombre": "Álgebra Lineal",
                    "comisiones": [{
                        "nombre": "A",
                        "docentes": [
                            {"nombre": "VARGAS GABRIELA", "rol": "Titular"},
                            {"nombre": "VIBRENTIS FRANCISCO ANTONIO", "rol": "JTP"},
                            {"nombre": "ORECCHIA MONICA INES", "rol": "Ayudante de 1ra"},
                            {"nombre": "GRYNBERG SEBASTIAN PABLO", "rol": "Ayudante 2da"},
                            {"nombre": "CONDE MARIA EUGENIA", "rol": null}
                        ],
                        "clases": []
                    }]
                }))
                .unwrap();

                sanitizar_comisiones_materia(&mut materia);

                let docentes = &materia.comisiones[0].docentes;

                assert_eq!(docentes.len(), 5);
                assert_eq!(docentes[0].rol, Some(normalizar_rol("Titular")));
                assert_eq!(docentes[1].rol, Some(normalizar_rol("JTP")));
                assert_eq!(docentes[2].rol, Some(normalizar_rol("Ayudante de 1ra")));
                assert_eq!(docentes[3].rol, Some(normalizar_rol("Ayudante 2da")));
                assert_eq!(docentes[4].rol, None);
            }

            #[test]
            fn se_considera_una_comision_con_multiples_turnos_como_multiples_comisiones_distintas()
            {
                let docentes = serde_json::json!([{"nombre": "VARGAS GABRIELA", "rol": "Profesor/a Adjunto/a"}]);
                let mut materia: Materia = serde_json::from_value(serde_json::json!({
                    "codigo": "CB002",
                    "nombre": "Álgebra Lineal",
                    "comisiones": [
                        {
                            "nombre": "Comisión 1",
                            "docentes": docentes,
                            "clases": [{"tipo": "Teórica", "dia": "Lunes", "inicio": "09:00", "fin": "13:00"}]
                        },
                        {
                            "nombre": "Comisión 1",
                            "docentes": docentes,
                            "clases": [{"tipo": "Teórica", "dia": "Martes", "inicio": "18:00", "fin": "22:00"}]
                        }
                    ]
                }))
                .unwrap();

                sanitizar_comisiones_materia(&mut materia);

                assert_eq!(materia.comisiones.len(), 2);
                assert_eq!(materia.comisiones[0].clases[0].dia, "Lunes");
                assert_eq!(materia.comisiones[1].clases[0].dia, "Martes");
            }
        }

        mod normalizar_rol_docente {
            use super::super::super::*;

            #[test]
            fn se_normalizan_las_variantes_de_cada_rol() {
                assert_eq!(normalizar_rol("JTP"), "Jefe/a Trabajos Practicos");
                assert_eq!(normalizar_rol("Titular"), "Profesor/a Titular/a");
                assert_eq!(normalizar_rol("Ayudante 1ra"), "Ayudante 1ro/a");
                assert_eq!(normalizar_rol("Ayudante de 1ra"), "Ayudante 1ro/a");
                assert_eq!(
                    normalizar_rol("Ayudante 1ro/a ad-honorem"),
                    "Ayudante 1ro/a"
                );
                assert_eq!(normalizar_rol("Ayudante 2da"), "Ayudante 2do/a");
                assert_eq!(normalizar_rol("Ayudante de 2da"), "Ayudante 2do/a");
                assert_eq!(
                    normalizar_rol("Ayudante 2do/a ad-honorem"),
                    "Ayudante 2do/a"
                );
            }

            #[test]
            fn los_roles_ya_normalizados_no_cambian() {
                for rol in [
                    "Profesor/a Titular/a",
                    "Profesor/a Asociado/a",
                    "Profesor/a Adjunto/a",
                    "Profesor/a Consulto/a",
                    "Jefe/a Trabajos Practicos",
                    "Ayudante 1ro/a",
                    "Ayudante 2do/a",
                    "Coordinador/a",
                ] {
                    assert_eq!(normalizar_rol(rol), rol);
                }
            }

            #[test]
            fn los_roles_desconocidos_se_mantienen() {
                assert_eq!(normalizar_rol("Rol Inventado"), "Rol Inventado");
            }
        }
    }

    mod consolidar {
        use super::super::*;

        #[test]
        fn las_ofertas_vacias_producen_un_resultado_vacio() {
            let ofertas_carreras = vec![];
            let ofertas_materias = consolidar_ofertas_materias(ofertas_carreras);
            assert_eq!(ofertas_materias.len(), 0);
        }

        #[test]
        fn las_ofertas_con_materias_distintas_conservan_todas_las_materias() {
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
}
