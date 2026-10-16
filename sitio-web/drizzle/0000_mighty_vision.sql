CREATE EXTENSION IF NOT EXISTS unaccent;--> statement-breakpoint
CREATE EXTENSION IF NOT EXISTS pg_trgm;--> statement-breakpoint
CREATE TABLE "calificacion_dolly" (
	"codigo" serial PRIMARY KEY NOT NULL,
	"codigo_docente" uuid NOT NULL,
	"acepta_critica" numeric(2, 1) NOT NULL,
	"asistencia" numeric(2, 1) NOT NULL,
	"buen_trato" numeric(2, 1) NOT NULL,
	"claridad" numeric(2, 1) NOT NULL,
	"clase_organizada" numeric(2, 1) NOT NULL,
	"cumple_horarios" numeric(2, 1) NOT NULL,
	"fomenta_participacion" numeric(2, 1) NOT NULL,
	"panorama_amplio" numeric(2, 1) NOT NULL,
	"responde_mails" numeric(2, 1) NOT NULL,
	"fecha_creacion" timestamp with time zone DEFAULT (now() AT TIME ZONE 'America/Argentina/Buenos_Aires'::text),
	CONSTRAINT "calificacion_dolly_acepta_critica_check" CHECK ((acepta_critica >= (1)::numeric) AND (acepta_critica <= (5)::numeric)),
	CONSTRAINT "calificacion_dolly_asistencia_check" CHECK ((asistencia >= (1)::numeric) AND (asistencia <= (5)::numeric)),
	CONSTRAINT "calificacion_dolly_buen_trato_check" CHECK ((buen_trato >= (1)::numeric) AND (buen_trato <= (5)::numeric)),
	CONSTRAINT "calificacion_dolly_claridad_check" CHECK ((claridad >= (1)::numeric) AND (claridad <= (5)::numeric)),
	CONSTRAINT "calificacion_dolly_clase_organizada_check" CHECK ((clase_organizada >= (1)::numeric) AND (clase_organizada <= (5)::numeric)),
	CONSTRAINT "calificacion_dolly_cumple_horarios_check" CHECK ((cumple_horarios >= (1)::numeric) AND (cumple_horarios <= (5)::numeric)),
	CONSTRAINT "calificacion_dolly_fomenta_participacion_check" CHECK ((fomenta_participacion >= (1)::numeric) AND (fomenta_participacion <= (5)::numeric)),
	CONSTRAINT "calificacion_dolly_panorama_amplio_check" CHECK ((panorama_amplio >= (1)::numeric) AND (panorama_amplio <= (5)::numeric)),
	CONSTRAINT "calificacion_dolly_responde_mails_check" CHECK ((responde_mails >= (1)::numeric) AND (responde_mails <= (5)::numeric))
);
--> statement-breakpoint
CREATE TABLE "carrera" (
	"codigo" serial PRIMARY KEY NOT NULL,
	"nombre" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE "catedra" (
	"codigo" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"codigo_materia" text NOT NULL,
	"activa" boolean DEFAULT false NOT NULL
);
--> statement-breakpoint
CREATE TABLE "catedra_docente" (
	"codigo_catedra" uuid NOT NULL,
	"codigo_docente" uuid NOT NULL,
	CONSTRAINT "catedra_docente_pkey" PRIMARY KEY("codigo_catedra","codigo_docente")
);
--> statement-breakpoint
CREATE TABLE "comentario" (
	"codigo" serial PRIMARY KEY NOT NULL,
	"codigo_docente" uuid NOT NULL,
	"codigo_cuatrimestre" integer NOT NULL,
	"contenido" text NOT NULL,
	"es_de_dolly" boolean DEFAULT false NOT NULL,
	"fecha_creacion" timestamp with time zone DEFAULT (now() AT TIME ZONE 'America/Argentina/Buenos_Aires'::text),
	"codigo_calificacion_dolly" integer
);
--> statement-breakpoint
CREATE TABLE "cuatrimestre" (
	"codigo" serial PRIMARY KEY NOT NULL,
	"numero" smallint NOT NULL,
	"anio" smallint NOT NULL,
	CONSTRAINT "cuatrimestre_numero_check" CHECK (numero = ANY (ARRAY[1, 2]))
);
--> statement-breakpoint
CREATE TABLE "docente" (
	"codigo" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"nombre" text NOT NULL,
	"codigo_materia" text NOT NULL,
	"resumen_comentarios" text,
	"comentarios_ultimo_resumen" integer DEFAULT 0 NOT NULL,
	"nombre_siu" text,
	"rol" text
);
--> statement-breakpoint
CREATE TABLE "equivalencia" (
	"codigo_materia_plan_vigente" text NOT NULL,
	"codigo_materia_plan_anterior" text NOT NULL,
	CONSTRAINT "equivalencia_pkey" PRIMARY KEY("codigo_materia_plan_vigente","codigo_materia_plan_anterior"),
	CONSTRAINT "equivalencia_check" CHECK (codigo_materia_plan_vigente <> codigo_materia_plan_anterior)
);
--> statement-breakpoint
CREATE TABLE "materia" (
	"codigo" text PRIMARY KEY NOT NULL,
	"nombre" text NOT NULL,
	"cuatrimestre_ultima_actualizacion" integer
);
--> statement-breakpoint
CREATE TABLE "oferta_comisiones" (
	"codigo_carrera" integer NOT NULL,
	"codigo_cuatrimestre" integer NOT NULL,
	"contenido" jsonb NOT NULL,
	CONSTRAINT "oferta_comisiones_pkey" PRIMARY KEY("codigo_carrera","codigo_cuatrimestre")
);
--> statement-breakpoint
CREATE TABLE "oferta_comisiones_raw" (
	"codigo_carrera" integer NOT NULL,
	"codigo_cuatrimestre" integer NOT NULL,
	"contenido" text NOT NULL,
	CONSTRAINT "oferta_comisiones_raw_pkey" PRIMARY KEY("codigo_carrera","codigo_cuatrimestre")
);
--> statement-breakpoint
CREATE TABLE "plan" (
	"codigo" serial PRIMARY KEY NOT NULL,
	"codigo_carrera" integer NOT NULL,
	"anio" smallint NOT NULL,
	"esta_vigente" boolean NOT NULL
);
--> statement-breakpoint
CREATE TABLE "plan_materia" (
	"codigo_plan" integer NOT NULL,
	"codigo_materia" text NOT NULL,
	"es_electiva" boolean NOT NULL,
	CONSTRAINT "plan_materia_pkey" PRIMARY KEY("codigo_plan","codigo_materia")
);
--> statement-breakpoint
CREATE TABLE "prioridad_rol" (
	"rol" text PRIMARY KEY NOT NULL,
	"prioridad" integer NOT NULL,
	CONSTRAINT "prioridad_rol_prioridad_check" CHECK ((prioridad >= 1) AND (prioridad <= 10))
);
--> statement-breakpoint
ALTER TABLE "calificacion_dolly" ADD CONSTRAINT "calificacion_dolly_codigo_docente_fkey" FOREIGN KEY ("codigo_docente") REFERENCES "public"."docente"("codigo") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "catedra" ADD CONSTRAINT "catedra_codigo_materia_fkey" FOREIGN KEY ("codigo_materia") REFERENCES "public"."materia"("codigo") ON DELETE no action ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "catedra_docente" ADD CONSTRAINT "catedra_docente_codigo_catedra_fkey" FOREIGN KEY ("codigo_catedra") REFERENCES "public"."catedra"("codigo") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "catedra_docente" ADD CONSTRAINT "catedra_docente_codigo_docente_fkey" FOREIGN KEY ("codigo_docente") REFERENCES "public"."docente"("codigo") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "comentario" ADD CONSTRAINT "comentario_codigo_docente_fkey" FOREIGN KEY ("codigo_docente") REFERENCES "public"."docente"("codigo") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "comentario" ADD CONSTRAINT "comentario_codigo_cuatrimestre_fkey" FOREIGN KEY ("codigo_cuatrimestre") REFERENCES "public"."cuatrimestre"("codigo") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "comentario" ADD CONSTRAINT "comentario_codigo_calificacion_dolly_fkey" FOREIGN KEY ("codigo_calificacion_dolly") REFERENCES "public"."calificacion_dolly"("codigo") ON DELETE set null ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "docente" ADD CONSTRAINT "docente_codigo_materia_fkey" FOREIGN KEY ("codigo_materia") REFERENCES "public"."materia"("codigo") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "equivalencia" ADD CONSTRAINT "equivalencia_codigo_materia_plan_vigente_fkey" FOREIGN KEY ("codigo_materia_plan_vigente") REFERENCES "public"."materia"("codigo") ON DELETE no action ON UPDATE cascade;--> statement-breakpoint
ALTER TABLE "equivalencia" ADD CONSTRAINT "equivalencia_codigo_materia_plan_anterior_fkey" FOREIGN KEY ("codigo_materia_plan_anterior") REFERENCES "public"."materia"("codigo") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "materia" ADD CONSTRAINT "materia_cuatrimestre_ultima_actualizacion_fkey" FOREIGN KEY ("cuatrimestre_ultima_actualizacion") REFERENCES "public"."cuatrimestre"("codigo") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "oferta_comisiones" ADD CONSTRAINT "oferta_comisiones_codigo_carrera_fkey" FOREIGN KEY ("codigo_carrera") REFERENCES "public"."carrera"("codigo") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "oferta_comisiones" ADD CONSTRAINT "oferta_comisiones_codigo_cuatrimestre_fkey" FOREIGN KEY ("codigo_cuatrimestre") REFERENCES "public"."cuatrimestre"("codigo") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "oferta_comisiones_raw" ADD CONSTRAINT "oferta_comisiones_raw_codigo_carrera_fkey" FOREIGN KEY ("codigo_carrera") REFERENCES "public"."carrera"("codigo") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "oferta_comisiones_raw" ADD CONSTRAINT "oferta_comisiones_raw_codigo_cuatrimestre_fkey" FOREIGN KEY ("codigo_cuatrimestre") REFERENCES "public"."cuatrimestre"("codigo") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "plan" ADD CONSTRAINT "plan_codigo_carrera_fkey" FOREIGN KEY ("codigo_carrera") REFERENCES "public"."carrera"("codigo") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "plan_materia" ADD CONSTRAINT "plan_materia_codigo_plan_fkey" FOREIGN KEY ("codigo_plan") REFERENCES "public"."plan"("codigo") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "plan_materia" ADD CONSTRAINT "plan_materia_codigo_materia_fkey" FOREIGN KEY ("codigo_materia") REFERENCES "public"."materia"("codigo") ON DELETE no action ON UPDATE cascade;--> statement-breakpoint
CREATE INDEX "calificacion_dolly_codigo_docente_idx" ON "calificacion_dolly" USING btree ("codigo_docente" uuid_ops);--> statement-breakpoint
CREATE INDEX "comentario_codigo_docente_codigo_cuatrimestre_fecha_creacio_idx" ON "comentario" USING btree ("codigo_docente","codigo_cuatrimestre" DESC NULLS FIRST,"fecha_creacion" DESC NULLS LAST,"codigo" DESC NULLS FIRST);--> statement-breakpoint
CREATE INDEX "comentario_codigo_docente_idx" ON "comentario" USING btree ("codigo_docente" uuid_ops);
--> statement-breakpoint
CREATE MATERIALIZED VIEW mv_codigos_fuentes_catedras AS
 SELECT m.codigo AS codigo_materia,
    m.nombre AS nombre_materia,
        CASE
            WHEN (m.cuatrimestre_ultima_actualizacion IS NULL) THEN NULL::text
            ELSE ((cuat.numero || 'C'::text) || cuat.anio)
        END AS cuatrimestre_ultima_actualizacion,
        CASE
            WHEN (m.cuatrimestre_ultima_actualizacion IS NOT NULL) THEN ARRAY[m.codigo]
            ELSE ARRAY( SELECT e.codigo_materia_plan_anterior
               FROM equivalencia e
              WHERE (e.codigo_materia_plan_vigente = m.codigo))
        END AS codigos_materias_fuente
   FROM (materia m
     LEFT JOIN cuatrimestre cuat ON ((cuat.codigo = m.cuatrimestre_ultima_actualizacion)))
  WITH NO DATA;
--> statement-breakpoint
CREATE UNIQUE INDEX v_codigos_materias_con_catedras_codigo_materia_idx ON mv_codigos_fuentes_catedras USING btree (codigo_materia);
--> statement-breakpoint
CREATE MATERIALIZED VIEW mv_catedras_mostradas_materia AS
 WITH docentes_catedra AS (
         SELECT DISTINCT mvcfc.codigo_materia,
            c.codigo AS codigo_catedra,
            d.nombre AS nombre_docente,
            p.prioridad
           FROM (((((mv_codigos_fuentes_catedras mvcfc
             CROSS JOIN LATERAL unnest(mvcfc.codigos_materias_fuente) codigo_materia_fuente(codigo_materia_fuente))
             JOIN catedra c ON ((c.codigo_materia = codigo_materia_fuente.codigo_materia_fuente)))
             JOIN catedra_docente cd ON ((cd.codigo_catedra = c.codigo)))
             JOIN docente d ON ((d.codigo = cd.codigo_docente)))
             LEFT JOIN prioridad_rol p ON ((p.rol = d.rol)))
          WHERE ((mvcfc.codigos_materias_fuente <> ARRAY[mvcfc.codigo_materia]) OR c.activa)
        ), catedras_agrupadas AS (
         SELECT dc.codigo_materia,
            dc.codigo_catedra,
            string_agg(dc.nombre_docente, '-'::text ORDER BY dc.prioridad, dc.nombre_docente) AS nombre_catedra
           FROM docentes_catedra dc
          GROUP BY dc.codigo_materia, dc.codigo_catedra
        ), catedras_deduplicadas AS (
         SELECT DISTINCT ON (catedras_agrupadas.codigo_materia, catedras_agrupadas.nombre_catedra) catedras_agrupadas.codigo_materia,
            catedras_agrupadas.codigo_catedra,
            catedras_agrupadas.nombre_catedra
           FROM catedras_agrupadas
          ORDER BY catedras_agrupadas.codigo_materia, catedras_agrupadas.nombre_catedra, catedras_agrupadas.codigo_catedra
        )
 SELECT codigo_materia,
    codigo_catedra,
    nombre_catedra
   FROM catedras_deduplicadas
  WITH NO DATA;
--> statement-breakpoint
CREATE UNIQUE INDEX mv_catedras_mostradas_materia_codigo_materia_codigo_catedra_idx ON mv_catedras_mostradas_materia USING btree (codigo_materia, codigo_catedra);
--> statement-breakpoint
CREATE MATERIALIZED VIEW mv_materias_planes_vigentes AS
 SELECT DISTINCT m.codigo,
    m.nombre
   FROM ((materia m
     JOIN plan_materia pm ON ((pm.codigo_materia = m.codigo)))
     JOIN plan p ON ((p.codigo = pm.codigo_plan)))
  WHERE (p.esta_vigente = true)
  ORDER BY m.codigo
  WITH NO DATA;
