DO $$
DECLARE
    tables text;
BEGIN
    SELECT
        string_agg(format('%I.%I', schemaname, tablename), ', ') INTO tables
    FROM
        pg_tables
    WHERE
        schemaname NOT IN ('pg_catalog', 'information_schema');
    EXECUTE 'TRUNCATE TABLE ' || tables || ' RESTART IDENTITY CASCADE';
END
$$;

WITH materias_a_actualizar AS (
    SELECT DISTINCT
        m.codigo,
        'COD' || LPAD((453 + ROW_NUMBER() OVER (ORDER BY m.codigo))::text, 3, '0') AS nuevo_codigo
    FROM
        materia m
        JOIN plan_materia pm ON m.codigo = pm.codigo_materia
        JOIN plan p ON pm.codigo_plan = p.codigo
    WHERE
        p.esta_vigente = TRUE
        AND m.codigo NOT LIKE 'COD%')
UPDATE
    materia
SET
    codigo = nuevo_codigo
FROM
    materias_a_actualizar
WHERE
    materia.codigo = materias_a_actualizar.codigo;

