INSERT INTO docente (nombre, codigo_materia, nombre_siu)
    VALUES ($1, $2, $3)
RETURNING
    codigo;

