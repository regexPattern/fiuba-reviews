-- DESCRIPCIÓN
-- Toma un listado de códigos de las materias del SIU y retorna el código
-- y el nombre de aquellas que son candidatas a actualización.
--
-- Una materia es candidata a actualización si el cuatrimestre de última
-- actualización de la misma es anterior al cuatrimestre más reciente
-- disponible. Esto no significa que si o si haya una actualización
-- disponible para la materia, sino que simplemente esto aún no se ha
-- verificado para el último cuatrimestre.
--
-- PARÁMETROS
-- $1: Arreglo de strings con los códigos de las materias del SIU.
--
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

