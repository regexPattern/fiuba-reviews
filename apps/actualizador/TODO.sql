ALTER TABLE materia
ADD COLUMN codigo_cuatrimestre_actualizacion INTEGER,
ADD CONSTRAINT materia_codigo_cuatrimestre_actualizacion_fkey
FOREIGN KEY (codigo_cuatrimestre_actualizacion) REFERENCES cuatrimestre(codigo);

ALTER TABLE materia
DROP COLUMN cuatrimestre_ultima_actualizacion CASCADE;
