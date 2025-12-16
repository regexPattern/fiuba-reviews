package main

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

//go:embed queries/1-SELECT-ofertas-carreras.sql
var selectOfertasCarrerasQuery string

type OfertaCarrera struct {
	CodigoCarrera string          `db:"codigo_carrera"`
	Carrera       string          `db:"nombre_carrera"`
	Cuatrimestre  Cuatrimestre    `db:"cuatrimestre"`
	Materias      []OfertaMateria `db:"contenido"`
}

type Cuatrimestre struct {
	Numero int `json:"numero"`
	Anio   int `json:"anio"`
}

func (c Cuatrimestre) String() string {
	return fmt.Sprintf("%vQ%v", c.Numero, c.Anio)
}

type Materia struct {
	Codigo string `db:"codigo" json:"codigo"`
	Nombre string `db:"nombre" json:"nombre"`
}

type OfertaMateria struct {
	Materia
	Catedras []Catedra `json:"catedras"`
}

type Catedra struct {
	Codigo   int       `json:"codigo"`
	Docentes []Docente `json:"docentes"`
}

type Docente struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

type UltimaOfertaMateria struct {
	NombreCarrera string
	OfertaMateria
	Cuatrimestre
}

func getOfertasMaterias(conn *pgx.Conn) (map[string]UltimaOfertaMateria, error) {
	rows, err := conn.Query(context.TODO(), selectOfertasCarrerasQuery)
	if err != nil {
		return nil, fmt.Errorf("error consultando ofertas de comisiones de carreras: %w", err)
	}

	ofertasCarreras, err := pgx.CollectRows(rows, pgx.RowToStructByName[OfertaCarrera])
	if err != nil {
		return nil, fmt.Errorf("error serializando ofertas de comisiones de carreras")
	}

	ofertasPorCuatri := make(map[Cuatrimestre]int)
	for _, oc := range ofertasCarreras {
		ofertasPorCuatri[oc.Cuatrimestre]++
	}

	slog.Debug(
		fmt.Sprintf("encontradas %v ofertas de comisiones de carreras", len(ofertasCarreras)),
	)

	ofertasMaterias := make(map[string]UltimaOfertaMateria)
	materiasPorCuatri := make(map[Cuatrimestre]int)

	for _, oc := range ofertasCarreras {
		for _, om := range oc.Materias {
			var docentesCatedra int

			for _, cat := range om.Catedras {
				if n := len(cat.Docentes); n == 0 {
					slog.Warn(
						fmt.Sprintf(
							"cátedra de materia %v (%v) no tiene docentes",
							om.Codigo,
							om.Nombre,
						),
						"carrera",
						oc.Carrera,
						"cuatrimestre",
						oc.Cuatrimestre,
					)
				} else {
					docentesCatedra += len(cat.Docentes)
				}

				for _, doc := range cat.Docentes {
					if doc.Nombre == "" {
						slog.Warn(
							fmt.Sprintf(
								"docente sin nombre encontrado en cátedra de materia %v (%v)",
								om.Codigo,
								om.Nombre,
							),
							"carrera",
							oc.Carrera,
							"cuatrimestre",
							oc.Cuatrimestre,
							"cátedra",
							cat.Codigo,
						)
					}
				}
			}

			if docentesCatedra == 0 {
				slog.Warn(
					fmt.Sprintf(
						"oferta de comisiones de materia %v (%v) no tiene docentes",
						om.Codigo,
						om.Nombre,
					),
					"carrera", oc.Carrera,
					"cuatrimestre", oc.Cuatrimestre,
				)
				continue
			}

			if _, ok := ofertasMaterias[om.Codigo]; !ok {
				ofertasMaterias[om.Codigo] = UltimaOfertaMateria{
					NombreCarrera: oc.Carrera,
					OfertaMateria: om,
					Cuatrimestre:  oc.Cuatrimestre,
				}

				materiasPorCuatri[oc.Cuatrimestre]++
			}
		}
	}

	for cuat, n := range materiasPorCuatri {
		slog.Info(
			fmt.Sprintf(
				"encontradas %v ofertas de comisiones de materias de cuatrimestre %v",
				n,
				cuat,
			),
		)
	}

	slog.Info(
		fmt.Sprintf(
			"encontradas %v ofertas de comisiones de materias en total",
			len(ofertasMaterias),
		),
	)

	return ofertasMaterias, nil
}
