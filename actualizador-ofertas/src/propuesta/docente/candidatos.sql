-- Busca, para cada docente del SIU, el docente de la DB al que ya está vinculado y sus candidatos
-- entre los docentes de su misma materia.
--
-- Un docente del SIU está vinculado a todo docente de la materia cuyo `nombre_siu` es igual (sin
-- acentos, mayúsculas ni espacios repetidos). Los vinculados se devuelven siempre, primero y con
-- `vinculado` en true. Después se devuelven los demás docentes de la materia con similitud >= $3
-- (incluidos los que ya tienen otro `nombre_siu`), hasta completar $4 filas por docente del SIU.
-- Un docente del SIU sin vínculo ni candidatos se devuelve con las columnas del docente en NULL.
--
-- La DB suele guardar solo apellidos y el SIU "APELLIDO NOMBRE", por eso se toma la mejor de
-- las dos direcciones de `word_similarity` y de `similarity`.
--
-- Parámetros:
--   $1: text[] con el código de la materia de cada docente del SIU.
--   $2: text[] con el nombre de cada docente del SIU (misma posición que $1).
--   $3: real con la similitud mínima.
--   $4: bigint con la cantidad máxima de docentes por docente del SIU.
WITH entrantes AS (
	SELECT
		codigo_materia,
		nombre_siu,
		btrim(regexp_replace(lower(unaccent(nombre_siu)), '\s+', ' ', 'g')) AS nombre_norm
	FROM unnest($1::text[], $2::text[]) AS e (codigo_materia, nombre_siu)
),

docentes AS (
	SELECT
		codigo,
		codigo_materia,
		nombre,
		nombre_siu,
		btrim(regexp_replace(lower(unaccent(nombre)), '\s+', ' ', 'g')) AS nombre_norm,
		btrim(regexp_replace(lower(unaccent(nombre_siu)), '\s+', ' ', 'g')) AS nombre_siu_norm
	FROM docente
	WHERE codigo_materia = ANY($1::text[])
)

SELECT
	e.codigo_materia,
	e.nombre_siu,
	c.codigo,
	c.nombre,
	c.nombre_siu AS nombre_siu_actual,
	c.score,
	c.vinculado
FROM entrantes e
LEFT JOIN LATERAL (
	SELECT *
	FROM (
		SELECT
			d.codigo,
			d.nombre,
			d.nombre_siu,
			d.nombre_siu_norm = e.nombre_norm AS vinculado,
			greatest(
				word_similarity(d.nombre_norm, e.nombre_norm),
				word_similarity(e.nombre_norm, d.nombre_norm),
				similarity(d.nombre_norm, e.nombre_norm)
			) AS score
		FROM docentes d
		WHERE d.codigo_materia = e.codigo_materia
	) s
	WHERE s.vinculado OR s.score >= $3
	ORDER BY s.vinculado DESC, s.score DESC, s.nombre
	LIMIT $4
) c ON TRUE
ORDER BY e.codigo_materia, e.nombre_siu, c.vinculado DESC, c.score DESC, c.nombre;
