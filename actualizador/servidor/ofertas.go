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

	slog.Debug("ofertas_carreras_encontradas", "count", len(ofertasCarreras))

	ofertasMaterias := make(map[string]ofertaMateriaMasReciente)
	materiasPorCuatri := make(map[cuatrimestre]int)

	for _, ofCarr := range ofertasCarreras {
		for _, ofMat := range ofCarr.OfertasMaterias {
			logger := slog.Default().
				With("codigo_materia", ofMat.Codigo, "carrera", ofCarr.NombreCarrera, "cuatrimestre", ofCarr.Cuatrimestre)

			// Si la oferta de la materia no tiene cátedras o solo tiene cátedras sin docentes, no
			// se agregan a listado de ofertas disponibles. Lo más probable es que haya sido
			// producto de un error con el parser de ofertas.

			if len(ofMat.Catedras) == 0 {
				logger.Warn("oferta_materia_sin_catedras")
				continue
			}

			var docentesCatedra int
			for _, cat := range ofMat.Catedras {
				if n := len(cat.Docentes); n == 0 {
					logger.Warn("catedra_sin_docentes")
				} else {
					docentesCatedra += len(cat.Docentes)
				}
			}

			if docentesCatedra == 0 {
				logger.Warn("oferta_materia_sin_docentes")
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
		slog.Info("ofertas_materias_cuatrimestre", "count", n, "cuatrimestre", cuatri)
	}

	slog.Info("ofertas_materias_total", "count", len(ofertasMaterias))

	return ofertasMaterias, nil
}
