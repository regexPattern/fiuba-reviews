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

func getOfertasMaterias(conn *pgx.Conn) (map[string]UltimaOfertaMateria, error) {
	rows, err := conn.Query(context.Background(), selectOfertasCarrerasQuery)
	if err != nil {
		return nil, fmt.Errorf("error consultando ofertas de comisiones de carreras: %w", err)
	}

	ofertasCarreras, err := pgx.CollectRows(rows, pgx.RowToStructByName[OfertaCarrera])
	if err != nil {
		return nil, fmt.Errorf("error procesando ofertas de comisiones de carreras")
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
						fmt.Sprintf("cátedra de materia %v no tiene docentes", om.Codigo),
						"carrera",
						oc.Carrera,
						"cuatrimestre",
						oc.Cuatrimestre,
					)
				} else {
					docentesCatedra += len(cat.Docentes)
				}
			}

			if docentesCatedra == 0 {
				slog.Warn(
					fmt.Sprintf(
						"oferta de comisiones de materia %v (%v) no tiene docentes",
						om.Codigo,
						om.Nombre,
					),
					"carrera",
					oc.Carrera,
					"cuatrimestre",
					fmt.Sprintf("%vQ%v", oc.Cuatrimestre.Numero, oc.Cuatrimestre.Anio),
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

	for c, n := range materiasPorCuatri {
		slog.Debug(
			fmt.Sprintf(
				"encontradas %v ofertas de comisiones de materias de cuatrimestre %v",
				n,
				c,
			),
		)
	}

	slog.Debug(
		fmt.Sprintf(
			"encontradas %v ofertas de comisiones de materias en total",
			len(ofertasMaterias),
		),
	)

	return ofertasMaterias, nil
}
