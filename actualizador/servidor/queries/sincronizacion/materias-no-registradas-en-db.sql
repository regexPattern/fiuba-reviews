-- Parámetros
-- $1: Arreglo de strings con los nombres de las materias del SIU.
-- $2: Arreglo de strings con los codigos de las materias del SIU.
SELECT
    siu.codigo_siu AS codigo,
    siu.nombre_siu AS nombre
FROM
    unnest($1::text[], $2::text[]) AS siu (nombre_siu,
        codigo_siu)
WHERE
    NOT EXISTS (
        SELECT
            1
        FROM
            materia mat
        WHERE
            lower(unaccent (mat.nombre)) = lower(unaccent (siu.nombre_siu)));

