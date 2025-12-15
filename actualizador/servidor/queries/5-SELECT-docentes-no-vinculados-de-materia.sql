-- Parámetros
-- $1: Código de la materia.
-- $2: Arreglo de strings con los nombres de los docentes de la materia del SIU.
-- Normalización: case-insensitive, sin tildes, sin espacios extra al inicio/final
WITH nombres_siu AS (
    SELECT
        unnest($2::text[]) AS nombre,
        trim(regexp_replace(lower(unaccent (unnest($2::text[]))), '\s+', ' ', 'g')) AS nombre_norm
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
                AND trim(regexp_replace(lower(unaccent (d.nombre_siu)), '\s+', ' ', 'g')) = ns.nombre_norm)
),
sin_match_exacto AS (
    SELECT
        nombre,
        nombre_norm
    FROM
        nombres_siu
    EXCEPT
    SELECT
        nombre,
        trim(regexp_replace(lower(unaccent (nombre)), '\s+', ' ', 'g')) AS nombre_norm
    FROM
        con_match_exacto
)
SELECT
    sme.nombre AS nombre_siu,
    d.codigo::text AS codigo,
    d.nombre AS nombre_db,
    -- Scores: preferencia casi perfecta si hay substring normalizado, sino similarity normalizada
    CASE WHEN strpos(sme.nombre_norm, d.nombre_norm) > 0 THEN
        1.0 - 0.001 * ABS(LENGTH(sme.nombre_norm) - LENGTH(d.nombre_norm))
    ELSE
        similarity (d.nombre_norm, sme.nombre_norm)
    END AS similitud
FROM
    sin_match_exacto sme
    INNER JOIN LATERAL (
        SELECT
            codigo,
            nombre,
            trim(regexp_replace(lower(unaccent (nombre)), '\s+', ' ', 'g')) AS nombre_norm
        FROM
            docente
        WHERE
            codigo_materia = $1
            AND nombre_siu IS NULL
            AND similarity (trim(regexp_replace(lower(unaccent (nombre)), '\s+', ' ', 'g')), sme.nombre_norm) >= 0.4) d ON TRUE
ORDER BY
    sme.nombre,
    similitud DESC
