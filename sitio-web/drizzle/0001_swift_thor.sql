ALTER TABLE "comentario" RENAME COLUMN "es_de_dolly" TO "de_dolly";--> statement-breakpoint
ALTER TABLE "materia" RENAME COLUMN "cuatrimestre_ultima_actualizacion" TO "codigo_cuatrimestre_ultima_actualizacion";--> statement-breakpoint
ALTER TABLE "plan" RENAME COLUMN "esta_vigente" TO "vigente";--> statement-breakpoint
ALTER TABLE "plan_materia" RENAME COLUMN "es_electiva" TO "electiva";--> statement-breakpoint
ALTER TABLE "materia" RENAME CONSTRAINT "materia_cuatrimestre_ultima_actualizacion_fkey" TO "materia_codigo_cuatrimestre_ultima_actualizacion_fkey";
