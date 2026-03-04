-- DESCRIPCIÓN
-- Actualiza el cuatrimestre de última actualización de una materia sin cambios.
--
-- Cuando una materia no tiene cambios en un cuatrimestre nuevo, se marca
-- como actualizada en ese cuatrimestre para evitar que se considere como
-- pendiente de actualización en el futuro.
--
-- PARÁMETROS
-- $1: código de la materia (text)
-- $2: número del cuatrimestre (int)
-- $3: año del cuatrimestre (int)
UPDATE
    materia
SET
    cuatrimestre_ultima_actualizacion = (
        SELECT
            codigo
        FROM
            cuatrimestre
        WHERE
            numero = $2
            AND anio = $3)
WHERE
    codigo = $1;
