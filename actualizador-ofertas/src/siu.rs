use serde::Deserialize;
use std::collections::{HashMap, hash_map::Entry};
use tracing::{debug, warn};

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

/// Limpia en el lugar las ofertas de las carreras, excluyendo comisiones condicionales
/// o sin docentes y unificando los roles docentes. Conserva las materias sin comisiones.
pub fn sanitizar_ofertas_carreras(ofertas_carreras: &mut [OfertaCarrera]) {
    for oferta in ofertas_carreras {
        for materia in &mut oferta.materias {
            sanitizar_comisiones_materia(materia);
        }
    }
}

fn sanitizar_comisiones_materia(materia: &mut Materia) {
    let codigo = &materia.codigo;
    materia.comisiones.retain(|c| {
        if c.nombre.trim().to_uppercase().contains("CONDICIONAL") {
            debug!(materia = %codigo, comision = %c.nombre, "comisión condicional excluida");
            return false;
        }
        if c.docentes.is_empty() {
            debug!(materia = %codigo, comision = %c.nombre, "comisión sin docentes excluida");
            return false;
        }
        true
    });

    if materia.comisiones.is_empty() {
        warn!(
            materia = %materia.codigo,
            nombre = %materia.nombre,
            "materia sin comisiones"
        );
    }

    for comision in &mut materia.comisiones {
        if comision.docentes.iter().any(|d| d.nombre.trim().is_empty()) {
            warn!(
                materia = %materia.codigo,
                comision = %comision.nombre,
                "comisión con docente con nombre vacío"
            );
        }
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
/// Los roles sin variantes conocidas, incluidos los desconocidos, se mantienen tal cual.
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
    mod sanitizar {
        #[test]
        fn se_conservan_las_materias_que_quedan_sin_comisiones() {
            let mut ofertas: Vec<super::super::OfertaCarrera> =
                serde_json::from_value(serde_json::json!([{
                    "carrera": "Ingeniería Civil",
                    "anio": 2026,
                    "cuatrimestre": 1,
                    "materias": [
                        {"codigo": "CB002", "nombre": "Álgebra Lineal", "comisiones": []},
                        {"codigo": "CB003", "nombre": "Análisis Matemático", "comisiones": [
                            {"nombre": "A", "docentes": [], "clases": []}
                        ]}
                    ]
                }]))
                .unwrap();

            super::super::sanitizar_ofertas_carreras(&mut ofertas);

            assert_eq!(ofertas[0].materias.len(), 2);
            assert_eq!(ofertas[0].materias[0].codigo, "CB002");
            assert_eq!(ofertas[0].materias[1].codigo, "CB003");
            assert!(ofertas[0].materias.iter().all(|m| m.comisiones.is_empty()));
        }

        mod sanitizar_comisiones_materia {
            use super::super::super::*;
            use tracing_test::traced_test;

            #[test]
            #[traced_test]
            fn se_advierten_las_materias_sin_comisiones() {
                let mut materia: Materia = serde_json::from_value(serde_json::json!({
                    "codigo": "CB002",
                    "nombre": "Álgebra Lineal",
                    "comisiones": []
                }))
                .unwrap();

                sanitizar_comisiones_materia(&mut materia);

                assert!(logs_contain("materia sin comisiones"));
                assert!(logs_contain("materia=CB002"));
                assert!(logs_contain("nombre=Álgebra Lineal"));
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
                assert!(logs_contain("materia=CB002"));
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
                assert!(logs_contain("materia=CB002"));
            }

            #[test]
            #[traced_test]
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
                assert!(logs_contain("comisión con docente con nombre vacío"));
                assert!(logs_contain("materia=CB002"));
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
        use tracing_test::traced_test;

        #[test]
        fn las_ofertas_vacias_producen_un_resultado_vacio() {
            let ofertas_carreras = vec![];
            let ofertas_materias = consolidar_ofertas_materias(ofertas_carreras);
            assert_eq!(ofertas_materias.len(), 0);
        }

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

            let esperadas: Vec<OfertaMateria> = serde_json::from_value(serde_json::json!([
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
            for periodos in [[2025, 2026], [2026, 2025]] {
                let ofertas: Vec<OfertaCarrera> = periodos
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
        fn una_carrera_sin_comisiones_no_excluye_las_de_otra_del_mismo_periodo() {
            let ofertas = serde_json::from_value(serde_json::json!([
                {
                    "carrera": "Ingeniería Civil", "anio": 2026, "cuatrimestre": 1,
                    "materias": [{"codigo": "CB002", "nombre": "Álgebra Lineal", "comisiones": []}]
                },
                {
                    "carrera": "Ingeniería Industrial", "anio": 2026, "cuatrimestre": 1,
                    "materias": [{"codigo": "CB002", "nombre": "Álgebra Lineal", "comisiones": [
                        {"nombre": "A", "docentes": [], "clases": []}
                    ]}]
                }
            ]))
            .unwrap();

            let consolidadas = consolidar_ofertas_materias(ofertas);

            assert_eq!(consolidadas.len(), 1);
            assert_eq!(consolidadas[0].carreras, ["Ingeniería Industrial"]);
            assert_eq!(consolidadas[0].materia.comisiones.len(), 1);
            assert_eq!(consolidadas[0].materia.comisiones[0].nombre, "A");
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
}
