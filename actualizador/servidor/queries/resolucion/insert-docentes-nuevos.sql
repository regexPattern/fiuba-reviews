-- Inserta nuevos docentes
-- $1: código de la materia (text)
-- $2: array de nombres_siu (text[])
-- $3: array de nombres_db (text[])
-- $4: array de roles (text[])
INSERT INTO docente (codigo_materia, nombre_siu, nombre, rol)
SELECT
    $1,
    u.nombre_siu,
    u.nombre_db,
    u.rol
FROM
    unnest($2::text[], $3::text[], $4::text[]) AS u (nombre_siu,
        nombre_db,
        rol);
