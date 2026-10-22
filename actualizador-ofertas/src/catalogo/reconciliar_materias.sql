-- Reconcializar de materias.
--
-- Imposible hacer esto a mano, asi que gracias Claude Sonnet 5.5 (high).
--
-- Sincroniza los códigos provisionales ('XX###') de las materias vigentes con los códigos
-- oficiales del SIU y migra los datos base (docentes, calificaciones y comentarios) desde sus
-- equivalencias de planes anteriores.
--
-- Se apoya en `materia_vigente`, que se asume al día al momento de ejecutarse. Una materia solo se
-- procesa una vez: al cambiar su código deja de cumplir la condición `^XX[0-9]{3}$`.
--
-- Parámetros:
--   $1: text[] con los nombres de las materias del SIU.
--   $2: text[] con los códigos de las materias del SIU (misma posición que $1).
WITH catalogo AS (
	SELECT nombre, codigo
	FROM unnest($1::text[], $2::text[]) AS c (nombre, codigo)
),

reemplazos AS (
	SELECT
		mv.codigo AS codigo_anterior,
		c.codigo AS codigo_nuevo,
		mv.nombre,
		mv.codigos_materias_fuente AS codigos_equivalencias
	FROM catalogo c
	JOIN materia_vigente mv ON lower(unaccent(mv.nombre)) = lower(unaccent(c.nombre))
	WHERE mv.codigo ~ '^XX[0-9]{3}$'
		AND mv.usa_equivalencias
		AND mv.codigo <> c.codigo
),

-- Las FKs de `equivalencia`, `plan_materia` y `catedra` propagan el cambio con ON UPDATE CASCADE.
materias_actualizadas AS (
	UPDATE materia m
	SET codigo = r.codigo_nuevo
	FROM reemplazos r
	WHERE m.codigo = r.codigo_anterior
	RETURNING m.codigo, r.codigo_anterior, r.nombre, r.codigos_equivalencias
),

-- Si varias equivalencias tienen un docente con el mismo nombre se conserva una sola fuente:
-- la que tenga más calificaciones, luego más comentarios y, por último, el menor código.
docentes_elegidos AS (
	SELECT DISTINCT ON (r.codigo_nuevo, d.nombre)
		r.codigo_nuevo AS codigo_materia,
		d.codigo AS codigo_docente_anterior,
		d.nombre,
		d.resumen_comentarios,
		d.comentarios_ultimo_resumen
	FROM reemplazos r
	JOIN docente d ON d.codigo_materia = ANY (r.codigos_equivalencias)
	ORDER BY
		r.codigo_nuevo,
		d.nombre,
		(SELECT count(*) FROM calificacion_dolly c WHERE c.codigo_docente = d.codigo) DESC,
		(SELECT count(*) FROM comentario c WHERE c.codigo_docente = d.codigo) DESC,
		d.codigo
),

-- Los códigos nuevos se generan acá para poder mapear docente anterior -> docente nuevo.
docentes_a_migrar AS MATERIALIZED (
	SELECT *, gen_random_uuid() AS codigo_docente_nuevo
	FROM docentes_elegidos
),

docentes_insertados AS (
	INSERT INTO docente (codigo, nombre, codigo_materia, resumen_comentarios, comentarios_ultimo_resumen)
	SELECT
		codigo_docente_nuevo, nombre, codigo_materia, resumen_comentarios, comentarios_ultimo_resumen
	FROM docentes_a_migrar
	RETURNING codigo, codigo_materia
),

-- Se reserva el código de cada calificación nueva para poder remapear `comentario`.
calificaciones_a_migrar AS MATERIALIZED (
	SELECT
		nextval(pg_get_serial_sequence('calificacion_dolly', 'codigo'))::int AS codigo_nuevo,
		c.codigo AS codigo_anterior,
		dm.codigo_docente_nuevo,
		c.acepta_critica,
		c.asistencia,
		c.buen_trato,
		c.claridad,
		c.clase_organizada,
		c.cumple_horarios,
		c.fomenta_participacion,
		c.panorama_amplio,
		c.responde_mails,
		c.fecha_creacion
	FROM calificacion_dolly c
	JOIN docentes_a_migrar dm ON dm.codigo_docente_anterior = c.codigo_docente
),

calificaciones_insertadas AS (
	INSERT INTO calificacion_dolly (
		codigo, codigo_docente, acepta_critica, asistencia, buen_trato, claridad, clase_organizada,
		cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails, fecha_creacion
	)
	SELECT
		codigo_nuevo, codigo_docente_nuevo, acepta_critica, asistencia, buen_trato, claridad,
		clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails,
		fecha_creacion
	FROM calificaciones_a_migrar
	RETURNING codigo_docente
),

comentarios_insertados AS (
	INSERT INTO comentario (
		codigo_docente, codigo_cuatrimestre, contenido, de_dolly, fecha_creacion,
		codigo_calificacion_dolly
	)
	SELECT
		dm.codigo_docente_nuevo,
		c.codigo_cuatrimestre,
		c.contenido,
		c.de_dolly,
		c.fecha_creacion,
		cm.codigo_nuevo
	FROM comentario c
	JOIN docentes_a_migrar dm ON dm.codigo_docente_anterior = c.codigo_docente
	LEFT JOIN calificaciones_a_migrar cm ON cm.codigo_anterior = c.codigo_calificacion_dolly
	RETURNING codigo_docente
)

SELECT
	ma.codigo,
	ma.codigo_anterior,
	ma.nombre,
	ma.codigos_equivalencias,
	(
		SELECT count(*) FROM docentes_insertados di WHERE di.codigo_materia = ma.codigo
	)::int AS docentes_migrados,
	(
		SELECT count(*)
		FROM calificaciones_insertadas ci
		JOIN docentes_insertados di ON di.codigo = ci.codigo_docente
		WHERE di.codigo_materia = ma.codigo
	)::int AS calificaciones_migradas,
	(
		SELECT count(*)
		FROM comentarios_insertados ci
		JOIN docentes_insertados di ON di.codigo = ci.codigo_docente
		WHERE di.codigo_materia = ma.codigo
	)::int AS comentarios_migrados
FROM materias_actualizadas ma
ORDER BY ma.codigo;
