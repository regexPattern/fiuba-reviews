UPDATE
    docente
SET
    nombre = $1,
    nombre_siu = $2
WHERE
    codigo = $3
