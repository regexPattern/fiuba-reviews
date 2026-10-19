/// Devuelve la forma canónica de un rol de docente.
///
/// El SIU escribe el mismo rol de distintas formas (por ejemplo `JTP` y
/// `Jefe/a Trabajos Practicos`, o `Ayudante 1ra` y `Ayudante 1ro/a`). Se unifican para que
/// el rol pueda compararse y ordenarse sin tratar variantes del mismo rol como distintas.
/// Los roles sin variantes conocidas, incluidos los desconocidos, se mantienen tal cual.
pub(super) fn normalizar_rol(rol: &str) -> String {
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

#[cfg(test)]
mod tests {
    use super::normalizar_rol;

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
