-- Migra los docentes, sus comentarios y calificaciones hacia una materia desde
-- sus equivalencias. Es decir, si una nueva materia A (de un nuevo plan) tiene
-- como equivalencias a las materias B y C (de un plan anterior), la materia A
-- todavía no tiene información sobre las comisiones propias de esta materia,
-- por lo que lo más útil es asumir que es un renombramiento de sus materias
-- equivalentes del plan anterior. En el peor de los casos la materia A tiene
-- sus propias cátedras y no utiliza estos datos copiados.
--
-- Claramente generada con inteligencia artificial, pero si, hace lo que se
-- pretende.
WITH materias_equivalentes AS (
    SELECT
        e.codigo_materia_plan_anterior AS codigo_materia_equivalente
    FROM
        equivalencia e
    WHERE
        e.codigo_materia_plan_vigente = $1
),
docentes_equivalencias AS (
    SELECT
        d.codigo AS codigo_antiguo,
        gen_random_uuid () AS codigo_nuevo,
    d.nombre,
    d.resumen_comentarios,
    d.comentarios_ultimo_resumen
FROM
    docente d
    JOIN materias_equivalentes me ON d.codigo_materia = me.codigo_materia_equivalente
),
docentes_copiados AS (
INSERT INTO docente (codigo, nombre, codigo_materia, resumen_comentarios, comentarios_ultimo_resumen)
    SELECT
        de.codigo_nuevo,
        de.nombre,
        $1,
        de.resumen_comentarios,
        de.comentarios_ultimo_resumen
    FROM
        docentes_equivalencias de
    ON CONFLICT (nombre,
        codigo_materia)
        DO NOTHING
    RETURNING
        codigo,
        nombre
),
mapeo_codigos_docentes AS (
    SELECT
        de.codigo_antiguo,
        de.codigo_nuevo
    FROM
        docentes_equivalencias de
    WHERE
        EXISTS (
            SELECT
                1
            FROM
                docente d
            WHERE
                d.codigo = de.codigo_nuevo
                AND d.nombre = de.nombre
                AND d.codigo_materia = $1)
),
calificaciones_dolly_copiadas AS (
INSERT INTO calificacion_dolly (codigo_docente, acepta_critica, asistencia, buen_trato, claridad, clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails)
    SELECT
        m.codigo_nuevo,
        c.acepta_critica,
        c.asistencia,
        c.buen_trato,
        c.claridad,
        c.clase_organizada,
        c.cumple_horarios,
        c.fomenta_participacion,
        c.panorama_amplio,
        c.responde_mails
    FROM
        calificacion_dolly c
        JOIN mapeo_codigos_docentes m ON c.codigo_docente = m.codigo_antiguo
),
comentarios_copiados AS (
INSERT INTO comentario (codigo_docente, codigo_cuatrimestre, contenido, es_de_dolly, fecha_creacion)
    SELECT
        m.codigo_nuevo,
        cm.codigo_cuatrimestre,
        cm.contenido,
        cm.es_de_dolly,
        cm.fecha_creacion
    FROM
        comentario cm
        JOIN mapeo_codigos_docentes m ON cm.codigo_docente = m.codigo_antiguo
),
materia_actualizada AS (
    UPDATE
        materia
    SET
        docentes_migrados_de_equivalencia = TRUE
    WHERE
        codigo = $1
)
SELECT
    count(*)
FROM
    mapeo_codigos_docentes
