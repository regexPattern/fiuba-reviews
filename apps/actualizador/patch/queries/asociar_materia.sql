-- Resumen:
-- Actualiza el código de una materia específica en la base de datos.
UPDATE
    materia
SET
    codigo = $1 -- Código de la materia traído del SIU
WHERE
    codigo = $2 -- Código actual de la materia a actualizar
