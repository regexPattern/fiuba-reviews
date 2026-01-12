CREATE SCHEMA IF NOT EXISTS "public";

CREATE TABLE IF NOT EXISTS "public"."carrera" (
    "codigo" serial PRIMARY KEY,
    "nombre" text NOT NULL
);

CREATE TABLE IF NOT EXISTS "public"."cuatrimestre" (
    "codigo" serial PRIMARY KEY,
    "numero" smallint NOT NULL CHECK (numero IN (1, 2)),
    "anio" smallint NOT NULL
);

CREATE TABLE IF NOT EXISTS "public"."materia" (
    "codigo" text PRIMARY KEY,
    "nombre" text NOT NULL,
    "cuatrimestre_ultima_actualizacion" integer REFERENCES "public"."cuatrimestre" ("codigo"),
    "docentes_migrados_de_equivalencia" boolean DEFAULT FALSE NOT NULL
);

CREATE TABLE IF NOT EXISTS "public"."plan" (
    "codigo" serial PRIMARY KEY,
    "codigo_carrera" integer NOT NULL REFERENCES "public"."carrera" ("codigo"),
    "anio" smallint NOT NULL,
    "esta_vigente" boolean NOT NULL
);

CREATE TABLE IF NOT EXISTS "public"."plan_materia" (
    "codigo_plan" integer NOT NULL REFERENCES "public"."plan" ("codigo"),
    "codigo_materia" text NOT NULL REFERENCES "public"."materia" ("codigo") ON UPDATE CASCADE,
    "es_electiva" boolean NOT NULL,
    PRIMARY KEY ("codigo_plan", "codigo_materia")
);

CREATE TABLE IF NOT EXISTS "public"."equivalencia" (
    "codigo_materia_plan_vigente" text NOT NULL REFERENCES "public"."materia" ("codigo") ON UPDATE CASCADE,
    "codigo_materia_plan_anterior" text NOT NULL REFERENCES "public"."materia" ("codigo"),
    CHECK (codigo_materia_plan_vigente <> codigo_materia_plan_anterior),
    PRIMARY KEY ("codigo_materia_plan_vigente", "codigo_materia_plan_anterior")
);

CREATE TABLE IF NOT EXISTS "public"."docente" (
    "codigo" uuid DEFAULT gen_random_uuid () PRIMARY KEY,
    "nombre" text NOT NULL,
    "codigo_materia" text NOT NULL REFERENCES "public"."materia" ("codigo"),
    "resumen_comentarios" text,
    "comentarios_ultimo_resumen" integer DEFAULT 0 NOT NULL,
    "nombre_siu" text DEFAULT NULL,
    "rol" text NULL,
    UNIQUE ("nombre", "codigo_materia")
);

CREATE TABLE IF NOT EXISTS "public"."catedra" (
    "codigo" uuid DEFAULT gen_random_uuid () PRIMARY KEY,
    "codigo_materia" text NOT NULL REFERENCES "public"."materia" ("codigo") ON UPDATE CASCADE,
    "activa" boolean DEFAULT FALSE NOT NULL
);

CREATE TABLE IF NOT EXISTS "public"."catedra_docente" (
    "codigo_catedra" uuid NOT NULL REFERENCES "public"."catedra" ("codigo") ON DELETE CASCADE,
    "codigo_docente" uuid NOT NULL REFERENCES "public"."docente" ("codigo"),
    PRIMARY KEY ("codigo_catedra", "codigo_docente")
);

CREATE TABLE IF NOT EXISTS "public"."comentario" (
    "codigo" serial PRIMARY KEY,
    "codigo_docente" uuid NOT NULL REFERENCES "public"."docente" ("codigo"),
    "codigo_cuatrimestre" integer NOT NULL REFERENCES "public"."cuatrimestre" ("codigo"),
    "contenido" text NOT NULL,
    "es_de_dolly" boolean DEFAULT FALSE NOT NULL,
    "fecha_creacion" timestamp with time zone DEFAULT (now() AT TIME ZONE 'America/Argentina/Buenos_Aires')
);

CREATE INDEX ON "public"."comentario" ("codigo_docente");

CREATE TABLE IF NOT EXISTS "public"."calificacion_dolly" (
    "codigo" serial PRIMARY KEY,
    "codigo_docente" uuid NOT NULL REFERENCES "public"."docente" ("codigo"),
    "acepta_critica" numeric(2, 1) NOT NULL CHECK (acepta_critica >= 1 AND acepta_critica <= 5),
    "asistencia" numeric(2, 1) NOT NULL CHECK (asistencia >= 1 AND asistencia <= 5),
    "buen_trato" numeric(2, 1) NOT NULL CHECK (buen_trato >= 1 AND buen_trato <= 5),
    "claridad" numeric(2, 1) NOT NULL CHECK (claridad >= 1 AND claridad <= 5),
    "clase_organizada" numeric(2, 1) NOT NULL CHECK (clase_organizada >= 1 AND clase_organizada <= 5),
    "cumple_horarios" numeric(2, 1) NOT NULL CHECK (cumple_horarios >= 1 AND cumple_horarios <= 5),
    "fomenta_participacion" numeric(2, 1) NOT NULL CHECK (fomenta_participacion >= 1 AND fomenta_participacion <= 5),
    "panorama_amplio" numeric(2, 1) NOT NULL CHECK (panorama_amplio >= 1 AND panorama_amplio <= 5),
    "responde_mails" numeric(2, 1) NOT NULL CHECK (responde_mails >= 1 AND responde_mails <= 5)
);

CREATE INDEX ON "public"."calificacion_dolly" ("codigo_docente");

CREATE TABLE IF NOT EXISTS "public"."oferta_comisiones" (
    "codigo_carrera" integer NOT NULL REFERENCES "public"."carrera" ("codigo"),
    "codigo_cuatrimestre" integer NOT NULL REFERENCES "public"."cuatrimestre" ("codigo"),
    "contenido" jsonb NOT NULL,
    PRIMARY KEY ("codigo_carrera", "codigo_cuatrimestre")
);

CREATE EXTENSION unaccent;

CREATE EXTENSION pg_trgm;

-- DESPUES DE INSERTAR

INSERT INTO oferta_comisiones (codigo_carrera, codigo_cuatrimestre, contenido)
    VALUES (1, 22, pg_read_file('/ingenieria-civil.json')::json),
    (2, 22, pg_read_file('/ingenieria-electronica.json')::json),
    (10, 22, pg_read_file('/ingenieria-en-energia-electrica.json')::json),
    (4, 22, pg_read_file('/ingenieria-en-informatica.json')::json),
    (11, 22, pg_read_file('/ingenieria-en-petroleo.json')::json),
    (3, 22, pg_read_file('/ingenieria-industrial.json')::json),
    (5, 22, pg_read_file('/ingenieria-mecanica.json')::json),
    (7, 21, pg_read_file('/ingenieria-quimica.json')::json);

SELECT
    setval(pg_get_serial_sequence('public.comentario', 'codigo'), (
            SELECT
                MAX(codigo)
            FROM public.comentario));

SELECT
    setval(pg_get_serial_sequence('public.carrera', 'codigo'), (
            SELECT
                MAX(codigo)
            FROM public.carrera));

SELECT
    setval(pg_get_serial_sequence('public.cuatrimestre', 'codigo'), (
            SELECT
                MAX(codigo)
            FROM public.cuatrimestre));

SELECT
    setval(pg_get_serial_sequence('public.plan', 'codigo'), (
            SELECT
                MAX(codigo)
            FROM public.plan));

SELECT
    setval(pg_get_serial_sequence('public.calificacion_dolly', 'codigo'), (
            SELECT
                MAX(codigo)
            FROM public.calificacion_dolly));

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

