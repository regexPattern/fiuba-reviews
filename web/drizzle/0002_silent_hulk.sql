DROP INDEX "comentario_codigo_docente_codigo_cuatrimestre_fecha_creacio_idx";--> statement-breakpoint
DROP INDEX "comentario_codigo_docente_idx";--> statement-breakpoint
CREATE INDEX "catedra_codigo_materia_activa_idx" ON "catedra" USING btree ("codigo_materia","activa");--> statement-breakpoint
CREATE INDEX "catedra_docente_codigo_docente_idx" ON "catedra_docente" USING btree ("codigo_docente");--> statement-breakpoint
CREATE INDEX "comentario_docente_cuatrimestre_idx" ON "comentario" USING btree ("codigo_docente","codigo_cuatrimestre" DESC NULLS LAST,"codigo" DESC NULLS LAST);--> statement-breakpoint
CREATE INDEX "docente_codigo_materia_nombre_siu_idx" ON "docente" USING btree ("codigo_materia","nombre_siu");--> statement-breakpoint
CREATE INDEX "equivalencia_codigo_materia_plan_anterior_idx" ON "equivalencia" USING btree ("codigo_materia_plan_anterior");--> statement-breakpoint
CREATE INDEX "plan_materia_codigo_materia_idx" ON "plan_materia" USING btree ("codigo_materia");
