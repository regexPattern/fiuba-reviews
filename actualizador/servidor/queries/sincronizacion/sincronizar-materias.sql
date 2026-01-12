-- Parámetros
-- $1: Arreglo de strings con los nombres de las materias del SIU.
-- $2: Arreglo de strings con los codigos de las materias del SIU.
WITH materias_a_actualizar AS (
    SELECT
        mat.codigo AS codigo_antiguo,
        siu.codigo_siu AS codigo_nuevo,
        mat.nombre
    FROM
        materia mat
        JOIN (
            SELECT
                *
            FROM
                unnest($1::text[], $2::text[]) AS t (nombre_siu,
                codigo_siu)) siu ON lower(unaccent (mat.nombre)) = lower(unaccent (siu.nombre_siu))
    WHERE
        mat.codigo IS DISTINCT FROM siu.codigo_siu
        AND EXISTS (
            SELECT
                1
            FROM
                plan_materia pm
                INNER JOIN plan ON plan.codigo = pm.codigo_plan
            WHERE
                pm.codigo_materia = mat.codigo
                AND plan.esta_vigente)
),
materias_actualizadas AS (
    UPDATE
        materia mat
    SET
        codigo = maa.codigo_nuevo
    FROM
        materias_a_actualizar maa
    WHERE
        mat.codigo = maa.codigo_antiguo
    RETURNING
        mat.codigo AS codigo_nuevo,
        maa.codigo_antiguo,
        mat.nombre
),
equivalencias_por_materia AS (
    SELECT
        ma.codigo_nuevo,
        array_agg(e.codigo_materia_plan_anterior) AS codigos_equivalencias
 FROM
    materias_actualizadas ma
    JOIN equivalencia e ON e.codigo_materia_plan_vigente = ma.codigo_antiguo
GROUP BY
    ma.codigo_nuevo
),
docentes_con_calificaciones AS (
    SELECT
        ma.codigo_nuevo,
        d.codigo AS codigo_docente_antiguo,
        d.nombre,
        d.resumen_comentarios,
        d.comentarios_ultimo_resumen,
        (
            SELECT
                count(*)
            FROM
                calificacion_dolly c
            WHERE
                c.codigo_docente = d.codigo) AS num_calificaciones
        FROM
            materias_actualizadas ma
            JOIN equivalencia e ON e.codigo_materia_plan_vigente = ma.codigo_antiguo
            JOIN docente d ON d.codigo_materia = e.codigo_materia_plan_anterior
),
docentes_a_migrar AS (
    SELECT DISTINCT ON (codigo_nuevo,
        nombre)
        codigo_nuevo,
        codigo_docente_antiguo,
        nombre,
        resumen_comentarios,
        comentarios_ultimo_resumen
    FROM
        docentes_con_calificaciones
    ORDER BY
        codigo_nuevo,
        nombre,
        num_calificaciones DESC
),
docentes_insertados AS (
INSERT INTO docente (nombre, codigo_materia, resumen_comentarios, comentarios_ultimo_resumen)
    SELECT
        nombre,
        codigo_nuevo,
        resumen_comentarios,
        comentarios_ultimo_resumen
    FROM
        docentes_a_migrar
    RETURNING
        codigo AS codigo_docente_nuevo,
        nombre,
        codigo_materia
),
mapeo_docentes AS (
    SELECT
        di.codigo_docente_nuevo,
        dm.codigo_docente_antiguo,
        di.codigo_materia
    FROM
        docentes_insertados di
        JOIN docentes_a_migrar dm ON di.nombre = dm.nombre
            AND di.codigo_materia = dm.codigo_nuevo
),
calificaciones_copiadas AS (
INSERT INTO calificacion_dolly (codigo_docente, acepta_critica, asistencia, buen_trato, claridad, clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails)
    SELECT
        m.codigo_docente_nuevo,
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
    JOIN mapeo_docentes m ON c.codigo_docente = m.codigo_docente_antiguo
RETURNING
    codigo_docente
),
comentarios_copiados AS (
INSERT INTO comentario (codigo_docente, codigo_cuatrimestre, contenido, es_de_dolly, fecha_creacion)
    SELECT
        m.codigo_docente_nuevo,
        cm.codigo_cuatrimestre,
        cm.contenido,
        cm.es_de_dolly,
        cm.fecha_creacion
    FROM
        comentario cm
        JOIN mapeo_docentes m ON cm.codigo_docente = m.codigo_docente_antiguo
    RETURNING
        codigo_docente
),
conteo_docentes AS (
    SELECT
        codigo_materia,
        count(*) AS docentes_migrados
    FROM
        docentes_insertados
    GROUP BY
        codigo_materia
),
conteo_calificaciones AS (
    SELECT
        md.codigo_materia,
        count(*) AS calificaciones_migradas
    FROM
        calificaciones_copiadas cc
        JOIN mapeo_docentes md ON cc.codigo_docente = md.codigo_docente_nuevo
    GROUP BY
        md.codigo_materia
),
conteo_comentarios AS (
    SELECT
        md.codigo_materia,
        count(*) AS comentarios_migrados
    FROM
        comentarios_copiados cmc
        JOIN mapeo_docentes md ON cmc.codigo_docente = md.codigo_docente_nuevo
    GROUP BY
        md.codigo_materia
)
SELECT
    ma.codigo_nuevo AS codigo,
    lower(unaccent (ma.nombre)) AS nombre,
    COALESCE(cd.docentes_migrados, 0)::int AS docentes_migrados,
    COALESCE(ccm.comentarios_migrados, 0)::int AS comentarios_migrados,
    COALESCE(ccal.calificaciones_migradas, 0)::int AS calificaciones_migradas,
    COALESCE(eq.codigos_equivalencias, ARRAY[]::text[]) AS codigos_equivalencias
FROM
    materias_actualizadas ma
    LEFT JOIN equivalencias_por_materia eq ON eq.codigo_nuevo = ma.codigo_nuevo
    LEFT JOIN conteo_docentes cd ON cd.codigo_materia = ma.codigo_nuevo
    LEFT JOIN conteo_comentarios ccm ON ccm.codigo_materia = ma.codigo_nuevo
    LEFT JOIN conteo_calificaciones ccal ON ccal.codigo_materia = ma.codigo_nuevo;

