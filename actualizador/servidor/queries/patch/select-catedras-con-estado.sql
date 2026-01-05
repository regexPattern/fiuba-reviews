-- DESCRIPCIÓN
-- Retorna todas las cátedras del SIU de una materia con su estado (ya
-- existe o no existe).
--
-- Una cátedra del SIU ya existe en la base de datos si existe una cátedra que
-- tenga el mismo nombre o firma. La firma de una cátedra es la concatenación
-- de los nombres (normalizados) de los docentes de la cátedra.
--
-- PARÁMETROS
-- $1: Código de la materia.
-- $2: Arreglo JSONB con las cátedras de la materia del SIU.
--
WITH catedras_siu AS (
    SELECT
        cat_elem ->> 'codigo' AS codigo_siu,
        string_agg(trim(regexp_replace(lower(unaccent (doc_elem ->> 'nombre')), '\s+', ' ', 'g')), '-' ORDER BY trim(regexp_replace(lower(unaccent (doc_elem ->> 'nombre')), '\s+', ' ', 'g'))) AS firma_docentes
    FROM
        jsonb_array_elements($2::jsonb) AS cat_elem,
        jsonb_array_elements(cat_elem -> 'docentes') AS doc_elem
    GROUP BY
        cat_elem ->> 'codigo'
),
firmas_catedras_db AS (
    SELECT DISTINCT
        string_agg(trim(regexp_replace(lower(unaccent (COALESCE(d.nombre_siu, d.nombre))), '\s+', ' ', 'g')), '-' ORDER BY trim(regexp_replace(lower(unaccent (COALESCE(d.nombre_siu, d.nombre))), '\s+', ' ', 'g'))) AS firma_docentes
    FROM
        catedra c
        INNER JOIN catedra_docente cd ON cd.codigo_catedra = c.codigo
        INNER JOIN docente d ON d.codigo = cd.codigo_docente
    WHERE
        c.codigo_materia = $1
    GROUP BY
        c.codigo
)
SELECT
    cs.codigo_siu::int AS codigo,
    EXISTS (
        SELECT
            1
        FROM
            firmas_catedras_db fdb
        WHERE
            fdb.firma_docentes = cs.firma_docentes) AS ya_existente
FROM
    catedras_siu cs;

