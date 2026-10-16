CREATE TABLE "resumen_calificaciones_docente" (
	"codigo_docente" uuid PRIMARY KEY NOT NULL,
	"cantidad" integer NOT NULL,
	"promedio" double precision NOT NULL,
	"acepta_critica" double precision NOT NULL,
	"asistencia" double precision NOT NULL,
	"buen_trato" double precision NOT NULL,
	"claridad" double precision NOT NULL,
	"clase_organizada" double precision NOT NULL,
	"cumple_horarios" double precision NOT NULL,
	"fomenta_participacion" double precision NOT NULL,
	"panorama_amplio" double precision NOT NULL,
	"responde_mails" double precision NOT NULL
);
--> statement-breakpoint
ALTER TABLE "resumen_calificaciones_docente" ADD CONSTRAINT "resumen_calificaciones_docente_codigo_docente_fkey" FOREIGN KEY ("codigo_docente") REFERENCES "public"."docente"("codigo") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
CREATE FUNCTION refresh_resumen_calificaciones_docente(codigos_docente uuid[]) RETURNS void
LANGUAGE sql AS $$
	DELETE FROM resumen_calificaciones_docente r
	WHERE r.codigo_docente = ANY (codigos_docente)
		AND NOT EXISTS (
			SELECT 1 FROM calificacion_dolly c WHERE c.codigo_docente = r.codigo_docente
		);

	INSERT INTO resumen_calificaciones_docente (
		codigo_docente, cantidad, promedio, acepta_critica, asistencia, buen_trato, claridad,
		clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails
	)
	SELECT
		c.codigo_docente,
		count(*)::int,
		(avg(c.acepta_critica) + avg(c.asistencia) + avg(c.buen_trato) + avg(c.claridad)
			+ avg(c.clase_organizada) + avg(c.cumple_horarios) + avg(c.fomenta_participacion)
			+ avg(c.panorama_amplio) + avg(c.responde_mails))::double precision / 9.0,
		avg(c.acepta_critica)::double precision,
		avg(c.asistencia)::double precision,
		avg(c.buen_trato)::double precision,
		avg(c.claridad)::double precision,
		avg(c.clase_organizada)::double precision,
		avg(c.cumple_horarios)::double precision,
		avg(c.fomenta_participacion)::double precision,
		avg(c.panorama_amplio)::double precision,
		avg(c.responde_mails)::double precision
	FROM calificacion_dolly c
	WHERE c.codigo_docente = ANY (codigos_docente)
	GROUP BY c.codigo_docente
	ON CONFLICT (codigo_docente) DO UPDATE SET
		cantidad = EXCLUDED.cantidad,
		promedio = EXCLUDED.promedio,
		acepta_critica = EXCLUDED.acepta_critica,
		asistencia = EXCLUDED.asistencia,
		buen_trato = EXCLUDED.buen_trato,
		claridad = EXCLUDED.claridad,
		clase_organizada = EXCLUDED.clase_organizada,
		cumple_horarios = EXCLUDED.cumple_horarios,
		fomenta_participacion = EXCLUDED.fomenta_participacion,
		panorama_amplio = EXCLUDED.panorama_amplio,
		responde_mails = EXCLUDED.responde_mails;
$$;
--> statement-breakpoint
CREATE FUNCTION trg_resumen_calificaciones_insert() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
	PERFORM refresh_resumen_calificaciones_docente(
		ARRAY(SELECT DISTINCT codigo_docente FROM new_rows)
	);
	RETURN NULL;
END $$;
--> statement-breakpoint
CREATE FUNCTION trg_resumen_calificaciones_update() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
	PERFORM refresh_resumen_calificaciones_docente(
		ARRAY(
			SELECT codigo_docente FROM new_rows
			UNION
			SELECT codigo_docente FROM old_rows
		)
	);
	RETURN NULL;
END $$;
--> statement-breakpoint
CREATE FUNCTION trg_resumen_calificaciones_delete() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
	PERFORM refresh_resumen_calificaciones_docente(
		ARRAY(SELECT DISTINCT codigo_docente FROM old_rows)
	);
	RETURN NULL;
END $$;
--> statement-breakpoint
CREATE FUNCTION trg_resumen_calificaciones_truncate() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
	TRUNCATE resumen_calificaciones_docente;
	RETURN NULL;
END $$;
--> statement-breakpoint
CREATE TRIGGER resumen_calificaciones_insert
	AFTER INSERT ON calificacion_dolly
	REFERENCING NEW TABLE AS new_rows
	FOR EACH STATEMENT EXECUTE FUNCTION trg_resumen_calificaciones_insert();
--> statement-breakpoint
CREATE TRIGGER resumen_calificaciones_update
	AFTER UPDATE ON calificacion_dolly
	REFERENCING OLD TABLE AS old_rows NEW TABLE AS new_rows
	FOR EACH STATEMENT EXECUTE FUNCTION trg_resumen_calificaciones_update();
--> statement-breakpoint
CREATE TRIGGER resumen_calificaciones_delete
	AFTER DELETE ON calificacion_dolly
	REFERENCING OLD TABLE AS old_rows
	FOR EACH STATEMENT EXECUTE FUNCTION trg_resumen_calificaciones_delete();
--> statement-breakpoint
CREATE TRIGGER resumen_calificaciones_truncate
	AFTER TRUNCATE ON calificacion_dolly
	FOR EACH STATEMENT EXECUTE FUNCTION trg_resumen_calificaciones_truncate();
--> statement-breakpoint
SELECT refresh_resumen_calificaciones_docente(ARRAY(SELECT DISTINCT codigo_docente FROM calificacion_dolly));
