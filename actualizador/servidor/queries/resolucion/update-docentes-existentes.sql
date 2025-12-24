-- Actualiza nombre_siu, nombre y rol para docentes existentes
-- $1: array de códigos de docentes (uuid[])
-- $2: array de nombres_siu (text[])
-- $3: array de nombres_db (text[])
-- $4: array de roles (text[])
UPDATE
    docente
SET
    nombre_siu = u.nombre_siu,
    nombre = u.nombre_db,
    rol = u.rol
FROM
    unnest($1::uuid[], $2::text[], $3::text[], $4::text[]) AS u (codigo,
        nombre_siu,
        nombre_db,
        rol)
WHERE
    docente.codigo = u.codigo;
