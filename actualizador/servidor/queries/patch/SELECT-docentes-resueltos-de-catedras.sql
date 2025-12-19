-- DESCRIPCIÓN
-- Retorna todos los docentes de las cátedras del SIU de una materia con su estado de resolución.
-- Un docente está resuelto si ya existe un docente en la base de datos con el mismo nombre_siu.
--
-- PARÁMETROS
-- $1: Código de la materia.
-- $2: Arreglo JSONB con las cátedras de la materia del SIU.
--
WITH docentes_siu AS (
    SELECT
        (cat_elem ->> 'codigo')::int AS codigo_catedra_siu,
        doc_elem ->> 'nombre' AS nombre_docente_siu
    FROM
        jsonb_array_elements($2::jsonb) AS cat_elem,
        jsonb_array_elements(cat_elem -> 'docentes') AS doc_elem
),
docentes_resueltos AS (
    SELECT DISTINCT
        d.nombre_siu
    FROM
        docente d
    WHERE
        d.codigo_materia = $1
        AND d.nombre_siu IS NOT NULL
)
SELECT
    ds.codigo_catedra_siu,
    ds.nombre_docente_siu,
    EXISTS (
        SELECT 1
        FROM docentes_resueltos dr
        WHERE dr.nombre_siu = ds.nombre_docente_siu
    ) AS resuelto
FROM
    docentes_siu ds;