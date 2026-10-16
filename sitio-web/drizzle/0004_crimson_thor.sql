DROP MATERIALIZED VIEW mv_materias_planes_vigentes;
--> statement-breakpoint
DROP MATERIALIZED VIEW mv_catedras_mostradas_materia;
--> statement-breakpoint
DROP MATERIALIZED VIEW mv_codigos_fuentes_catedras;
--> statement-breakpoint

CREATE MATERIALIZED VIEW materia_vigente AS
WITH vigentes AS (
	SELECT pm.codigo_materia AS codigo, count(DISTINCT pm.codigo_plan)::int AS cantidad_planes
	FROM plan_materia pm
	JOIN plan p ON p.codigo = pm.codigo_plan
	WHERE p.vigente
	GROUP BY pm.codigo_materia
),
base AS (
	SELECT
		m.codigo,
		m.nombre,
		v.cantidad_planes,
		cu.numero AS cuatrimestre_numero,
		cu.anio AS cuatrimestre_anio,
		m.codigo_cuatrimestre_ultima_actualizacion IS NULL AS usa_equivalencias,
		COALESCE((
			SELECT jsonb_agg(
				jsonb_build_object('codigo', a.codigo, 'nombre', a.nombre) ORDER BY a.codigo
			)
			FROM equivalencia e
			JOIN materia a ON a.codigo = e.codigo_materia_plan_anterior
			WHERE e.codigo_materia_plan_vigente = m.codigo
		), '[]'::jsonb) AS equivalencias
	FROM vigentes v
	JOIN materia m ON m.codigo = v.codigo
	LEFT JOIN cuatrimestre cu ON cu.codigo = m.codigo_cuatrimestre_ultima_actualizacion
),
con_fuentes AS (
	SELECT
		b.*,
		CASE
			WHEN b.usa_equivalencias THEN ARRAY(
				SELECT e.codigo_materia_plan_anterior
				FROM equivalencia e
				WHERE e.codigo_materia_plan_vigente = b.codigo
				ORDER BY e.codigo_materia_plan_anterior
			)
			ELSE ARRAY[b.codigo]
		END AS codigos_materias_fuente
	FROM base b
),
con_catedras AS (
	SELECT
		f.*,
		(
			SELECT count(*)
			FROM catedra c
			WHERE c.codigo_materia = ANY (f.codigos_materias_fuente)
				AND (f.usa_equivalencias OR c.activa)
		)::int AS cantidad_catedras
	FROM con_fuentes f
)
SELECT
	codigo,
	nombre,
	cantidad_planes,
	cantidad_catedras,
	cuatrimestre_numero,
	cuatrimestre_anio,
	usa_equivalencias,
	codigos_materias_fuente,
	equivalencias,
	(row_number() OVER (ORDER BY cantidad_planes DESC, cantidad_catedras DESC, nombre))::int AS orden
FROM con_catedras;
--> statement-breakpoint
CREATE UNIQUE INDEX materia_vigente_codigo_idx ON materia_vigente (codigo);
--> statement-breakpoint
CREATE INDEX materia_vigente_orden_idx ON materia_vigente (orden);
--> statement-breakpoint

CREATE MATERIALIZED VIEW catedra_mostrada AS
SELECT
	mv.codigo AS codigo_materia,
	c.codigo AS codigo_catedra,
	c.codigo_materia AS codigo_materia_catedra,
	string_agg(
		d.nombre, '-' ORDER BY COALESCE(pr.prioridad, 2147483647), d.nombre
	) AS nombre
FROM materia_vigente mv
CROSS JOIN LATERAL unnest(mv.codigos_materias_fuente) AS f (codigo)
JOIN catedra c ON c.codigo_materia = f.codigo AND (mv.usa_equivalencias OR c.activa)
JOIN catedra_docente cd ON cd.codigo_catedra = c.codigo
JOIN docente d ON d.codigo = cd.codigo_docente
LEFT JOIN prioridad_rol pr ON pr.rol = d.rol
GROUP BY mv.codigo, c.codigo, c.codigo_materia;
--> statement-breakpoint
CREATE UNIQUE INDEX catedra_mostrada_codigo_materia_codigo_catedra_idx
	ON catedra_mostrada (codigo_materia, codigo_catedra);
--> statement-breakpoint

CREATE VIEW calificacion_catedra AS
SELECT cd.codigo_catedra, avg(r.promedio)::double precision AS calificacion
FROM catedra_docente cd
JOIN resumen_calificaciones_docente r ON r.codigo_docente = cd.codigo_docente
GROUP BY cd.codigo_catedra;
--> statement-breakpoint

CREATE FUNCTION refresh_materialized_views() RETURNS void LANGUAGE plpgsql AS $$
BEGIN
	REFRESH MATERIALIZED VIEW CONCURRENTLY materia_vigente;
	REFRESH MATERIALIZED VIEW CONCURRENTLY catedra_mostrada;
END $$;
