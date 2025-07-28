-- Resumen:
-- Obtiene las materias de planes vigentes que aún no han migrado
-- sus docentes desde las materias equivalentes. Esta query identifica
-- materias que requieren ejecutar el proceso de migración de docentes.
SELECT
    m.codigo,  -- Código de la materia
    m.nombre   -- Nombre de la materia
FROM
    materia m
    JOIN plan_materia pm ON m.codigo = pm.codigo_materia  -- Relación materia-plan
    JOIN plan p ON pm.codigo_plan = p.codigo              -- Información del plan
WHERE
    m.docentes_migrados_de_equivalencia = FALSE  -- Materias sin migrar docentes
    AND p.esta_vigente = TRUE                     -- Solo planes vigentes
