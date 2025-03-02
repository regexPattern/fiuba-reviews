package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func getCodigosMaterias() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	res, _ := conn.Query(context.Background(), `
SELECT m.codigo, lower(unaccent(m.nombre))
FROM materia m
INNER JOIN plan_materia pm
ON m.codigo = pm.codigo_materia
INNER JOIN plan p
ON p.codigo = pm.codigo_plan
WHERE p.esta_vigente = true;
		`)

	materias := make(map[string]string, 0)

	for res.Next() {
		var codigo, nombre string
		_ = res.Scan(&codigo, &nombre)
		materias[nombre] = codigo
	}
}
