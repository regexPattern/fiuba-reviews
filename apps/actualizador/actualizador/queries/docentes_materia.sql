SELECT
    codigo,
    nombre
    nombre_siu
FROM
    docente
WHERE
    codigo_materia = $1
