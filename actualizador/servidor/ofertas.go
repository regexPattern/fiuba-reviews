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

type ofertaCarrera struct {
	CodigoCarrera   string          `db:"codigo_carrera"`
	NombreCarrera   string          `db:"nombre_carrera"`
	Cuatrimestre    cuatrimestre    `db:"cuatrimestre"`
	OfertasMaterias []ofertaMateria `db:"contenido"`
}

type cuatrimestre struct {
	Numero int `json:"numero"`
	Anio   int `json:"anio"`
}

func (c cuatrimestre) String() string {
	return fmt.Sprintf("%vQ%v", c.Numero, c.Anio)
}

type materia struct {
	Codigo string `db:"codigo" json:"codigo"`
	Nombre string `db:"nombre" json:"nombre"`
}

type ofertaMateria struct {
	materia
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

type ultimaOfertaMateria struct {
	NombreCarrera string
	ofertaMateria
	cuatrimestre
}

func getOfertasMaterias(conn *pgx.Conn) (map[string]ultimaOfertaMateria, error) {
	rows, err := conn.Query(context.TODO(), selectOfertasCarrerasQuery)
	if err != nil {
		return nil, fmt.Errorf("error consultando ofertas de comisiones de carreras: %w", err)
	}

	ofertasCarreras, err := pgx.CollectRows(rows, pgx.RowToStructByName[ofertaCarrera])
	if err != nil {
		return nil, fmt.Errorf("error serializando ofertas de comisiones de carreras")
	}

	ofertasPorCuatri := make(map[cuatrimestre]int)
	for _, oc := range ofertasCarreras {
		ofertasPorCuatri[oc.Cuatrimestre]++
	}

	slog.Debug(
		fmt.Sprintf("encontradas %v ofertas de comisiones de carreras", len(ofertasCarreras)),
	)

	ofertasMaterias := make(map[string]ultimaOfertaMateria)
	materiasPorCuatri := make(map[cuatrimestre]int)

	for _, oc := range ofertasCarreras {
		for _, om := range oc.OfertasMaterias {
			logger := slog.Default().
				With("carrera", oc.NombreCarrera, "cuatrimestre", oc.Cuatrimestre)

			var docentesCatedra int

			for _, cat := range om.Catedras {
				if n := len(cat.Docentes); n == 0 {
					logger.Warn(
						fmt.Sprintf(
							"cátedra de materia %v (%v) no tiene docentes",
							om.Codigo,
							om.Nombre,
						),
					)
				} else {
					docentesCatedra += len(cat.Docentes)
				}

				for _, doc := range cat.Docentes {
					if doc.Nombre == "" {
						logger.Warn(
							fmt.Sprintf(
								"encontrado docente sin nombre en cátedra de materia %v (%v)",
								om.Codigo,
								om.Nombre,
							),
							"cátedra",
							cat.Codigo,
						)
						continue
					}
				}
			}

			if docentesCatedra == 0 {
				logger.Warn(
					fmt.Sprintf(
						"oferta de comisiones de materia %v (%v) no tiene docentes",
						om.Codigo,
						om.Nombre,
					),
				)
				continue
			}

			if _, ok := ofertasMaterias[om.Codigo]; !ok {
				ofertasMaterias[om.Codigo] = ultimaOfertaMateria{
					NombreCarrera: oc.NombreCarrera,
					ofertaMateria: om,
					cuatrimestre:  oc.Cuatrimestre,
				}

				materiasPorCuatri[oc.Cuatrimestre]++
			}
		}
	}

	for cuatri, n := range materiasPorCuatri {
		slog.Info(
			fmt.Sprintf(
				"encontradas %v ofertas de comisiones de materias de cuatrimestre %v",
				n,
				cuatri,
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
