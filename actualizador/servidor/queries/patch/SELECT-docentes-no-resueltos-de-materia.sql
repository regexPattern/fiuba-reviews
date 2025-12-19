-- DESCRIPCIÓN
-- Retorna los docentes del SIU que no están resueltos y sus posibles
-- matches. Una match de un docente es una propuesta de qué docente ya
-- registrado en la base de datos puede corresponder a dicho docente del
-- SIU. Para determinar esto se buscan coincidencias exactas por
-- nombre_siu, luego para los restantes busca posibles matches con
-- docentes sin nombre_siu usando similitud de palabras (fuzzy match).
--
-- Un docente del SIU podría no tener matches (se considera que en este
-- caso se trata de un docente no registra), o uno o varios matches. En
-- estos casos el match con mayor similitud tendrá el score más alto, en
-- base al criterio de comparación utilizado.
--
-- PARÁMETROS
-- $1: Código de la materia.
-- $2: Arreglo de strings con los nombres de los docentes de la materia del SIU.
--
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
    word_similarity (d.nombre_norm, sme.nombre_norm) AS score
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
            AND word_similarity (trim(regexp_replace(lower(unaccent (nombre)), '\s+', ' ', 'g')), sme.nombre_norm) >= 0.5) d ON TRUE
WHERE
    d.codigo IS NOT NULL
UNION ALL
SELECT
    sme.nombre AS nombre_siu,
    NULL::text AS codigo,
    NULL::text AS nombre_db,
    NULL::float AS score
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
            AND word_similarity (trim(regexp_replace(lower(unaccent (d.nombre)), '\s+', ' ', 'g')), sme.nombre_norm) >= 0.5)
ORDER BY
    nombre_siu,
    score DESC NULLS LAST;

