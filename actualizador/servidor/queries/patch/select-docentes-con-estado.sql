-- DESCRIPCIÓN
-- Retorna todos los docentes de las cátedras del SIU de una materia con su código de la base de datos.
-- Un docente del SIU está resuelto si ya existe un docente en la base
-- de datos con el campo nombre_siu igual al nombre del docente de SIU.
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
)
SELECT
    ds.codigo_catedra_siu,
    ds.nombre_docente_siu,
    d.codigo AS codigo_docente
FROM
    docentes_siu ds
    LEFT JOIN docente d ON d.codigo_materia = $1
        AND d.nombre_siu = ds.nombre_docente_siu
        AND d.nombre_siu IS NOT NULL;

