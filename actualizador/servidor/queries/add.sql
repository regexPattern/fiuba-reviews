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

ALTER TABLE docente
    ADD CONSTRAINT docente_nombre_siu_codigo_materia_key UNIQUE (nombre_siu, codigo_materia);

ALTER TABLE docente
    ADD COLUMN rol text NULL;

