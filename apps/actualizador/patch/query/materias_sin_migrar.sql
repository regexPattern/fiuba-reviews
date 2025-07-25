SELECT m.codigo, m.nombre
FROM materia m
JOIN plan_materia pm ON m.codigo = pm.codigo_materia
JOIN plan p ON pm.codigo_plan = p.codigo
WHERE m.docentes_migrados_de_equivalencia = false
  AND p.esta_vigente = true
