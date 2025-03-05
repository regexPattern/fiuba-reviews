package main

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
)

func getCodigosMaterias() (map[string]string, error) {
	logger := log.Default().WithPrefix("DB 💽")

	res, err := conn.Query(context.Background(), `
SELECT m.codigo, lower(unaccent(m.nombre))
FROM materia m
INNER JOIN plan_materia pm
ON m.codigo = pm.codigo_materia
INNER JOIN plan p
ON p.codigo = pm.codigo_plan
WHERE p.esta_vigente = true;
		`)

	if err != nil {
		return nil, err
	}

	cods := make(map[string]string, 0)

	for res.Next() {
		var cod, nombre string

		err := res.Scan(&cod, &nombre)
		if err != nil {
			return nil, err
		}

		cods[nombre] = cod
	}

	logger.Info(fmt.Sprintf("Obtenidas %v materias", len(cods)))

	return cods, nil
}
