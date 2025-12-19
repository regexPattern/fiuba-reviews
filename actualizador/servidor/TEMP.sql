-- CREATE
CREATE TABLE oferta_comisiones (
    codigo_carrera int REFERENCES carrera (codigo) NOT NULL,
    codigo_cuatrimestre int REFERENCES cuatrimestre (codigo) NOT NULL,
    contenido jsonb NOT NULL
);

INSERT INTO oferta_comisiones (codigo_carrera, codigo_cuatrimestre, contenido)
    VALUES (1, 22, pg_read_file('/ingenieria-civil.json')::json),
    (2, 22, pg_read_file('/ingenieria-electronica.json')::json),
    (10, 22, pg_read_file('/ingenieria-en-energia-electrica.json')::json),
    (4, 22, pg_read_file('/ingenieria-en-informatica.json')::json),
    (11, 22, pg_read_file('/ingenieria-en-petroleo.json')::json),
    (3, 22, pg_read_file('/ingenieria-industrial.json')::json),
    (5, 22, pg_read_file('/ingenieria-mecanica.json')::json),
    (7, 21, pg_read_file('/ingenieria-quimica.json')::json);

ALTER TABLE docente
    ADD COLUMN nombre_siu text DEFAULT NULL;

ALTER TABLE docente
    DROP CONSTRAINT docente_nombre_codigo_materia_key;

--
--
--
-- DESTROY
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

