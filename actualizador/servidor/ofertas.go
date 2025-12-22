package main

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/regexPattern/fiuba-reviews/actualizador/queries"
)

type ofertaCarrera struct {
	CodigoCarrera   string          `db:"codigo_carrera"`
	NombreCarrera   string          `db:"nombre_carrera"`
	Cuatrimestre    cuatrimestre    `db:"cuatrimestre"`
	OfertasMaterias []ofertaMateria `db:"contenido"`
}

type ofertaMateria struct {
	materia
	Catedras []catedra `json:"catedras"`
}

type materia struct {
	Codigo string `db:"codigo" json:"codigo"`
	Nombre string `db:"nombre" json:"nombre"`
}

type catedra struct {
	Codigo   int       `json:"codigo"`
	Docentes []docente `json:"docentes"`
}

type docente struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

type cuatrimestre struct {
	Numero int `json:"numero"`
	Anio   int `json:"anio"`
}

type ofertaMateriaMasReciente struct {
	NombreCarrera string
	ofertaMateria
	cuatrimestre
}

func newOfertasMaterias(conn *pgx.Conn) (map[string]ofertaMateriaMasReciente, error) {
	rows, err := conn.Query(context.TODO(), queries.OfertasCarreras)
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

	ofertasMaterias := make(map[string]ofertaMateriaMasReciente)
	materiasPorCuatri := make(map[cuatrimestre]int)

	for _, ofCarr := range ofertasCarreras {
		for _, ofMat := range ofCarr.OfertasMaterias {
			logger := slog.Default().
				With("carrera", ofCarr.NombreCarrera, "cuatrimestre", ofCarr.Cuatrimestre)

			var docentesCatedra int
			for _, cat := range ofMat.Catedras {
				if n := len(cat.Docentes); n == 0 {
					logger.Warn(
						fmt.Sprintf(
							"cátedra de materia %v no tiene docentes",
							ofMat.Codigo,
						),
					)
				} else {
					docentesCatedra += len(cat.Docentes)
				}
			}

			if docentesCatedra == 0 {
				logger.Warn(
					fmt.Sprintf(
						"oferta de comisiones de materia %v no tiene docentes",
						ofMat.Codigo,
					),
				)
				continue
			}

			if _, ok := ofertasMaterias[ofMat.Codigo]; !ok {
				ofertasMaterias[ofMat.Codigo] = ofertaMateriaMasReciente{
					NombreCarrera: ofCarr.NombreCarrera,
					ofertaMateria: ofMat,
					cuatrimestre:  ofCarr.Cuatrimestre,
				}

				materiasPorCuatri[ofCarr.Cuatrimestre]++
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
