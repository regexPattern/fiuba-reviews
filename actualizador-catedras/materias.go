package main

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/charmbracelet/log"
)

type materia struct {
	Codigo   string    `json:"codigo"`
	Nombre   string    `json:"nombre"`
	Catedras []catedra `json:"catedras"`
}

type catedra struct {
	Codigo   int       `json:"codigo"`
	Docentes []docente `json:"docentes"`
}

type docente struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

func fetchNombresACodigosMateriasDB() (map[string]string, error) {
	logger := log.Default().WithPrefix("DB 💽")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := conn.Query(ctx, `
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

	codigos := make(map[string]string, 0)

	for rows.Next() {
		var cod, nombre string

		err := rows.Scan(&cod, &nombre)
		if err != nil {
			return nil, err
		}

		codigos[nombre] = cod
	}

	logger.Info(fmt.Sprintf("Obtenidas %v materias", len(codigos)))

	return codigos, nil
}

func fetchCodigosMateriasDesactualizadas() {
}

func actualizarCodigosMaterias(materias []materia, nombresACodigosMaterias map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conn.Query(ctx, `
CREATE TEMP TABLE tmp_codigos_materias (
	nombre_materia_norm TEXT PRIMARY KEY,
	codigo_materia_siu TEXT NOT NULL
);
		`)

	return nil
}

func filtrarMateriasMasRecientes(planes []plan) []materia {
	maxMaterias := 0
	for _, p := range planes {
		maxMaterias += len(p.materias)
	}

	cuatris := make(map[string]cuatri, maxMaterias)
	materias := make(map[string]materia, maxMaterias)

	for _, p := range planes {
		for _, m := range p.materias {
			cuatriUltimoCambio, ok := cuatris[m.Nombre]

			if !ok || p.cuatri.esDespuesDe(cuatriUltimoCambio) {
				cuatris[m.Nombre] = p.cuatri
				materias[m.Nombre] = m
			}
		}
	}

	return slices.Collect(maps.Values(materias))
}
