ALTER TABLE materia
    ADD COLUMN codigo_cuatrimestre_actualizacion integer,
    ADD CONSTRAINT materia_codigo_cuatrimestre_actualizacion_fkey FOREIGN KEY (codigo_cuatrimestre_actualizacion) REFERENCES cuatrimestre (codigo);

ALTER TABLE materia
    DROP COLUMN cuatrimestre_ultima_actualizacion CASCADE;

ALTER TABLE docente
    DROP CONSTRAINT docente_codigo_materia_fkey_cascade;

ALTER TABLE docente
    ADD CONSTRAINT docente_codigo_materia_fkey_cascade FOREIGN KEY (codigo_materia) REFERENCES materia (codigo) ON UPDATE CASCADE ON DELETE NO ACTION;

ALTER TABLE equivalencia
    DROP CONSTRAINT equivalencia_codigo_materia_plan_anterior_fkey;

ALTER TABLE equivalencia
    ADD CONSTRAINT equivalencia_codigo_materia_plan_anterior_fkey FOREIGN KEY (codigo_materia_plan_anterior) REFERENCES materia (codigo) ON UPDATE CASCADE ON DELETE NO ACTION;

