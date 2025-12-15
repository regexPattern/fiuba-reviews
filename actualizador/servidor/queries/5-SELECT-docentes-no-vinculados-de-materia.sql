-- Parámetros
-- $1: Código de la materia.
-- $2: Arreglo de strings con los nombres de los docentes de la materia del SIU.
WITH nombres_siu AS (
    SELECT
        unnest($2::text[]) AS nombre
),
con_match_exacto AS (
    SELECT
        ns.nombre
    FROM
        nombres_siu ns
    WHERE
        EXISTS (
            SELECT
                1
            FROM
                docente d
            WHERE
                d.codigo_materia = $1
                AND d.nombre_siu = ns.nombre)
),
sin_match_exacto AS (
    SELECT
        nombre
    FROM
        nombres_siu
    EXCEPT
    SELECT
        nombre
    FROM
        con_match_exacto
)
SELECT
    sme.nombre AS nombre_siu,
    d.codigo::text AS codigo,
    d.nombre AS nombre_db,
    similarity (d.nombre, sme.nombre) AS similitud
FROM
    sin_match_exacto sme
    INNER JOIN LATERAL (
        SELECT
            codigo,
            nombre
        FROM
            docente
        WHERE
            codigo_materia = $1
            AND nombre_siu IS NULL
            AND similarity (nombre, sme.nombre) >= 0.5) d ON TRUE
ORDER BY
    sme.nombre,
    similitud DESC
