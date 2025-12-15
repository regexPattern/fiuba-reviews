-- Parámetros
-- $1: Código de la materia.
-- $2: Arreglo de strings con los nombres de los docentes de la materia del SIU.
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
    -- Scores: word_similarity para comparar cadenas de distinta longitud
    word_similarity (d.nombre_norm, sme.nombre_norm) AS similitud
FROM
    sin_match_exacto sme
    LEFT JOIN LATERAL (
        SELECT
            codigo,
            nombre,
            trim(regexp_replace(lower(unaccent (nombre)), '\s+', ' ', 'g')) AS nombre_norm
        FROM
            docente
        WHERE
            codigo_materia = $1
            AND nombre_siu IS NULL
            AND word_similarity (trim(regexp_replace(lower(unaccent (nombre)), '\s+', ' ', 'g')), sme.nombre_norm) >= 0.3) d ON TRUE
WHERE
    d.codigo IS NOT NULL
UNION ALL
SELECT
    sme.nombre AS nombre_siu,
    NULL::text AS codigo,
    NULL::text AS nombre_db,
    NULL::float AS similitud
FROM
    sin_match_exacto sme
WHERE
    NOT EXISTS (
        SELECT
            1
        FROM
            docente d
        WHERE
            d.codigo_materia = $1
            AND nombre_siu IS NULL
            AND word_similarity (trim(regexp_replace(lower(unaccent (d.nombre)), '\s+', ' ', 'g')), sme.nombre_norm) >= 0.3)
ORDER BY
    nombre_siu,
    similitud DESC NULLS LAST;

