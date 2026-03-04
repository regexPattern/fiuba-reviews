-- Actualiza el cuatrimestre de última actualización de una materia
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

