-- Resumen:
-- Obtiene las materias que tienen códigos temporales o sin asociar.
-- Busca materias cuyo código comience con 'COD%', lo cual indica que
-- son códigos provisionales que necesitan ser asociados a códigos definitivos.
SELECT
    codigo,  -- Código temporal de la materia
    nombre   -- Nombre de la materia
FROM
    materia
WHERE
    codigo LIKE 'COD%'  -- Filtro para códigos temporales que comienzan con 'COD'
