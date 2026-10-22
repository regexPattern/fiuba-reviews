-- Conserva los tres dígitos del código provisional; no asigna códigos oficiales ni copia historial.
-- Las FKs de plan_materia, equivalencia (plan vigente) y catedra usan ON UPDATE CASCADE.
UPDATE materia
SET codigo = 'XX' || substring(codigo FROM 4)
WHERE codigo ~ '^COD[0-9]{3}$';
--> statement-breakpoint
SELECT refresh_materialized_views();