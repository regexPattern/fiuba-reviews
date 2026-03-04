-- Sincroniza cátedras de una materia
-- $1: código de la materia (text)
-- $2: JSON array de cátedras del SIU con estructura [{codigo, docentes: [{nombre, rol}]}]
WITH docentes_siu AS (
    SELECT
        (cat_elem ->> 'codigo')::int AS codigo_catedra_siu,
        doc_elem ->> 'nombre' AS nombre_siu,
        d.codigo AS codigo_docente
    FROM
        jsonb_array_elements($2::jsonb) AS cat_elem,
        jsonb_array_elements(cat_elem -> 'docentes') AS doc_elem
        LEFT JOIN docente d ON d.codigo_materia = $1
            AND d.nombre_siu = doc_elem ->> 'nombre'
),
firmas_siu AS (
    SELECT
        codigo_catedra_siu,
        string_agg(trim(regexp_replace(lower(unaccent (nombre_siu)), '\s+', ' ', 'g')), '-' ORDER BY trim(regexp_replace(lower(unaccent (nombre_siu)), '\s+', ' ', 'g'))) AS firma,
        array_agg(codigo_docente) AS codigos_docentes,
        count(*) AS docentes_resueltos
    FROM
        docentes_siu
    WHERE
        codigo_docente IS NOT NULL
    GROUP BY
        codigo_catedra_siu
),
conteo_docentes_siu AS (
    SELECT
        (cat_elem ->> 'codigo')::int AS codigo_catedra_siu,
        jsonb_array_length(cat_elem -> 'docentes') AS total_docentes
    FROM
        jsonb_array_elements($2::jsonb) AS cat_elem
),
catedras_resueltas AS (
    SELECT
        fs.codigo_catedra_siu,
        fs.firma,
        fs.codigos_docentes
    FROM
        firmas_siu fs
        JOIN conteo_docentes_siu cs ON cs.codigo_catedra_siu = fs.codigo_catedra_siu
    WHERE
        fs.docentes_resueltos = cs.total_docentes
),
firmas_db AS (
    SELECT
        c.codigo,
        string_agg(trim(regexp_replace(lower(unaccent (COALESCE(d.nombre_siu, d.nombre))), '\s+', ' ', 'g')), '-' ORDER BY trim(regexp_replace(lower(unaccent (COALESCE(d.nombre_siu, d.nombre))), '\s+', ' ', 'g'))) AS firma
    FROM
        catedra c
        JOIN catedra_docente cd ON cd.codigo_catedra = c.codigo
        JOIN docente d ON d.codigo = cd.codigo_docente
    WHERE
        c.codigo_materia = $1
    GROUP BY
        c.codigo
),
catedras_match AS (
    SELECT
        cr.codigo_catedra_siu,
        cr.firma,
        cr.codigos_docentes,
        fdb.codigo AS codigo_catedra_existente
    FROM
        catedras_resueltas cr
        LEFT JOIN firmas_db fdb ON fdb.firma = cr.firma
),
desactivadas AS (
    UPDATE
        catedra
    SET
        activa = FALSE
    WHERE
        codigo_materia = $1
),
activadas AS (
    UPDATE
        catedra
    SET
        activa = TRUE
    WHERE
        codigo IN (
            SELECT
                codigo_catedra_existente
            FROM
                catedras_match
            WHERE
                codigo_catedra_existente IS NOT NULL)
        RETURNING
            codigo
),
nuevas AS (
INSERT INTO catedra (codigo_materia, activa)
    SELECT
        $1,
        TRUE
    FROM
        catedras_match
    WHERE
        codigo_catedra_existente IS NULL
    RETURNING
        codigo
),
nuevas_con_orden AS (
    SELECT
        codigo,
        row_number() OVER () AS orden
    FROM
        nuevas
),
catedras_nuevas_match AS (
    SELECT
        cm.codigos_docentes,
        nco.codigo AS codigo_catedra_nueva
    FROM (
        SELECT
            codigos_docentes,
            row_number() OVER () AS orden
    FROM
        catedras_match
    WHERE
        codigo_catedra_existente IS NULL) cm
    JOIN nuevas_con_orden nco ON nco.orden = cm.orden
),
insertar_docentes AS (
INSERT INTO catedra_docente (codigo_catedra, codigo_docente)
    SELECT
        cnm.codigo_catedra_nueva,
        unnest(cnm.codigos_docentes)
    FROM
        catedras_nuevas_match cnm
)
SELECT
    (
        SELECT
            count(*)
        FROM
            activadas) AS catedras_activadas,
    (
        SELECT
            count(*)
        FROM
            nuevas) AS catedras_creadas;

