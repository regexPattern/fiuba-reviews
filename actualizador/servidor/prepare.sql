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

