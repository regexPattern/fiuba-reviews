-- DESCRIPCIÓN
-- Sincroniza las cátedras resueltas de una materia. Una cátedra está resuelta
-- si todos sus docentes están en el array de códigos resueltos.
-- Si la cátedra ya existe (misma firma), la activa.
-- Si no existe, la crea y asocia los docentes.
--
-- PARÁMETROS
-- $1: Código de la materia.
-- $2: Arreglo JSONB con las cátedras de la materia del SIU.
-- $3: Arreglo JSONB con los códigos de docentes resueltos.
--
WITH 
-- Expandir docentes del SIU con sus códigos resueltos
docentes_siu AS (
    SELECT
        (cat_elem ->> 'codigo')::int AS codigo_catedra_siu,
        doc_elem ->> 'nombre' AS nombre_siu,
        d.codigo AS codigo_docente
    FROM
        jsonb_array_elements($2::jsonb) AS cat_elem,
        jsonb_array_elements(cat_elem -> 'docentes') AS doc_elem
        LEFT JOIN docente d ON d.codigo_materia = $1 
            AND d.nombre_siu = doc_elem ->> 'nombre'
    WHERE
        d.codigo = ANY(SELECT jsonb_array_elements_text($3::jsonb)::uuid)
),
-- Calcular firma de cada cátedra del SIU (solo con docentes resueltos)
firmas_siu AS (
    SELECT
        codigo_catedra_siu,
        string_agg(
            trim(regexp_replace(lower(unaccent(nombre_siu)), '\s+', ' ', 'g')), 
            '-' ORDER BY trim(regexp_replace(lower(unaccent(nombre_siu)), '\s+', ' ', 'g'))
        ) AS firma,
        array_agg(codigo_docente) AS codigos_docentes
    FROM docentes_siu
    GROUP BY codigo_catedra_siu
),
-- Contar docentes totales por cátedra del SIU (para saber si está completa)
conteo_docentes_siu AS (
    SELECT
        (cat_elem ->> 'codigo')::int AS codigo_catedra_siu,
        jsonb_array_length(cat_elem -> 'docentes') AS total_docentes
    FROM jsonb_array_elements($2::jsonb) AS cat_elem
),
-- Filtrar solo cátedras completamente resueltas
catedras_resueltas AS (
    SELECT fs.codigo_catedra_siu, fs.firma, fs.codigos_docentes
    FROM firmas_siu fs
    JOIN conteo_docentes_siu cs ON cs.codigo_catedra_siu = fs.codigo_catedra_siu
    WHERE array_length(fs.codigos_docentes, 1) = cs.total_docentes
),
-- Buscar cátedras existentes en la DB con sus firmas
firmas_db AS (
    SELECT
        c.codigo,
        string_agg(
            trim(regexp_replace(lower(unaccent(COALESCE(d.nombre_siu, d.nombre))), '\s+', ' ', 'g')), 
            '-' ORDER BY trim(regexp_replace(lower(unaccent(COALESCE(d.nombre_siu, d.nombre))), '\s+', ' ', 'g'))
        ) AS firma
    FROM catedra c
    JOIN catedra_docente cd ON cd.codigo_catedra = c.codigo
    JOIN docente d ON d.codigo = cd.codigo_docente
    WHERE c.codigo_materia = $1
    GROUP BY c.codigo
),
-- Match entre cátedras del SIU y existentes
catedras_match AS (
    SELECT 
        cr.codigo_catedra_siu,
        cr.firma,
        cr.codigos_docentes,
        fdb.codigo AS codigo_catedra_existente
    FROM catedras_resueltas cr
    LEFT JOIN firmas_db fdb ON fdb.firma = cr.firma
),
-- Activar cátedras existentes
activadas AS (
    UPDATE catedra
    SET activa = true
    WHERE codigo IN (SELECT codigo_catedra_existente FROM catedras_match WHERE codigo_catedra_existente IS NOT NULL)
    RETURNING codigo
),
-- Crear cátedras nuevas
nuevas AS (
    INSERT INTO catedra (codigo_materia, activa)
    SELECT $1, true
    FROM catedras_match
    WHERE codigo_catedra_existente IS NULL
    RETURNING codigo
),
-- Asociar nuevas cátedras con el orden de inserción
nuevas_con_orden AS (
    SELECT 
        n.codigo,
        row_number() OVER () AS orden
    FROM nuevas n
),
catedras_nuevas_match AS (
    SELECT 
        cm.codigos_docentes,
        nco.codigo AS codigo_catedra_nueva
    FROM (
        SELECT 
            codigos_docentes,
            row_number() OVER () AS orden
        FROM catedras_match 
        WHERE codigo_catedra_existente IS NULL
    ) cm
    JOIN nuevas_con_orden nco ON nco.orden = cm.orden
),
-- Insertar relaciones catedra_docente para las nuevas
insertar_docentes AS (
    INSERT INTO catedra_docente (codigo_catedra, codigo_docente)
    SELECT 
        cnm.codigo_catedra_nueva,
        unnest(cnm.codigos_docentes)
    FROM catedras_nuevas_match cnm
)
-- Retornar resumen
SELECT 
    (SELECT count(*) FROM activadas) AS catedras_activadas,
    (SELECT count(*) FROM nuevas) AS catedras_creadas;