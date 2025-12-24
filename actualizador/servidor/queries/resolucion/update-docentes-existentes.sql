-- Actualiza nombre_siu y nombre para docentes existentes
-- $1: array de códigos de docentes (uuid[])
-- $2: array de nombres_siu (text[])
-- $3: array de nombres_db (text[])
UPDATE
    docente
SET
    nombre_siu = u.nombre_siu,
    nombre = u.nombre_db
FROM
    unnest($1::uuid[], $2::text[], $3::text[]) AS u (codigo,
        nombre_siu,
        nombre_db)
WHERE
    docente.codigo = u.codigo;

