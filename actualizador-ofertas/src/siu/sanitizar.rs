use super::normalizar::normalizar_rol;
use super::{Materia, OfertaCarrera};
use tracing::{debug, warn};

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

#[cfg(test)]
mod tests {
    use tracing_test::traced_test;

    use super::*;

    #[test]
    fn se_conservan_las_materias_que_quedan_sin_comisiones() {
        let mut ofertas: Vec<_> = serde_json::from_value(serde_json::json!([{
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

        sanitizar_ofertas_carreras(&mut ofertas);

        assert_eq!(ofertas[0].materias.len(), 2);
        assert_eq!(ofertas[0].materias[0].codigo, "CB002");
        assert_eq!(ofertas[0].materias[1].codigo, "CB003");
        assert!(ofertas[0].materias.iter().all(|m| m.comisiones.is_empty()));
    }

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
        let docente =
            serde_json::json!([{"nombre": "VARGAS GABRIELA", "rol": "Profesor/a Adjunto/a"}]);
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
    fn se_considera_una_comision_con_multiples_turnos_como_multiples_comisiones_distintas() {
        let docentes =
            serde_json::json!([{"nombre": "VARGAS GABRIELA", "rol": "Profesor/a Adjunto/a"}]);
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
