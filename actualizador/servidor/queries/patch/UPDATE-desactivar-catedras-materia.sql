-- DESCRIPCIÓN
-- Desactiva todas las cátedras de una materia.
--
-- PARÁMETROS
-- $1: Código de la materia.
--
UPDATE catedra SET activa = false WHERE codigo_materia = $1;