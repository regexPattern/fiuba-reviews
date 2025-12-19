-- Parámetros
-- $1: Arreglo de strings con los codigos de las materias del SIU.
SELECT DISTINCT
    mat.codigo,
    mat.nombre
FROM
    materia mat
    INNER JOIN plan_materia pm ON pm.codigo_materia = mat.codigo
    INNER JOIN plan ON plan.codigo = pm.codigo_plan
WHERE
    plan.esta_vigente
    AND mat.codigo = ANY ($1::text[])
    AND mat.cuatrimestre_ultima_actualizacion IS DISTINCT FROM (
        SELECT
            max(codigo)
        FROM
            cuatrimestre);

