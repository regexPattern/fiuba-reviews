package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
)

func main() {
	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)
	slog.SetDefault(slog.New(logger))

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error(fmt.Sprintf("error estableciendo conexión con la base de datos: %v", err))
		os.Exit(1)
	}

	slog.Info("conexión establecida con la base de datos")

	ofertasMaterias, err := getOfertasMaterias(conn)
	if err != nil {
		slog.Error(fmt.Sprintf("error obteniendo ofertas de comisiones de materias: %v", err))
		os.Exit(1)
	}

	codigosMaterias := make([]string, 0, len(ofertasMaterias))
	nombresMaterias := make([]string, 0, len(ofertasMaterias))
	for codigo, oferta := range ofertasMaterias {
		codigosMaterias = append(codigosMaterias, codigo)
		nombresMaterias = append(nombresMaterias, oferta.Nombre)
	}

	if err := syncMateriasDb(conn, codigosMaterias, nombresMaterias); err != nil {
		slog.Error(
			fmt.Sprintf("error sincronizando materias de la base de datos con el siu: %v", err),
		)
		os.Exit(1)
	}

	materiasPendientes, err := getPatchesMateria(conn, codigosMaterias, ofertasMaterias)
	if err != nil {
		slog.Error(fmt.Sprintf("error obteniendo materias a actualizar: %v", err))
		os.Exit(1)
	}

	slog.Debug(
		fmt.Sprintf(
			"encontradas %d materias con actualizaciones pendientes",
			len(materiasPendientes),
		),
	)

	// for i, mat := range materiasAActualizar {
	// 	slog.Debug(fmt.Sprintf("materia %d: %s - %d docentes pendientes - %d catedras nuevas",
	// 		i+1, mat.Nombre, len(mat.DocentesPendientes), len(mat.CatedrasNuevas)))
	// }

	for _, m := range materiasPendientes {
		_ = m
	}
}

func getOfertasMaterias(conn *pgx.Conn) (map[string]UltimaOfertaMateria, error) {
	rows, err := conn.Query(context.Background(), `
SELECT
    oc.codigo_carrera,
    lower(unaccent (carr.nombre)) AS nombre_carrera,
    json_build_object('numero', cuat.numero, 'anio', cuat.anio) AS cuatrimestre,
    oc.contenido
FROM
    oferta_comisiones oc
    INNER JOIN cuatrimestre cuat ON cuat.codigo = oc.codigo_cuatrimestre
    INNER JOIN carrera carr ON carr.codigo = oc.codigo_carrera
ORDER BY
    cuat.codigo DESC;
		`)
	if err != nil {
		return nil, fmt.Errorf("error consultando ofertas de comisiones de carreras: %w", err)
	}

	ofertasCarreras, err := pgx.CollectRows(rows, pgx.RowToStructByName[OfertaCarrera])
	if err != nil {
		return nil, fmt.Errorf("error procesando ofertas de comisiones de carreras")
	}

	ofertasPorCuatrimestre := make(map[Cuatrimestre]int)
	for _, oc := range ofertasCarreras {
		ofertasPorCuatrimestre[oc.Cuatrimestre]++
	}

	slog.Debug(
		fmt.Sprintf("encontradas %v ofertas de comisiones de carreras", len(ofertasCarreras)),
	)

	ofertasMaterias := make(map[string]UltimaOfertaMateria)
	materiasPorCuatri := make(map[Cuatrimestre]int)

	for _, oc := range ofertasCarreras {
		for _, om := range oc.Materias {
			var docentesCatedra int

			for _, c := range om.Catedras {
				if n := len(c.Docentes); n == 0 {
					slog.Warn(
						fmt.Sprintf("cátedra de materia %v no tiene docentes", om.Codigo),
						"carrera",
						oc.Carrera,
						"cuatrimestre",
						fmt.Sprintf("%vQ%v", oc.Cuatrimestre.Numero, oc.Cuatrimestre.Anio),
					)
				} else {
					docentesCatedra += len(c.Docentes)
				}

				for _, d := range c.Docentes {
					_ = d
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
				"encontradas %v ofertas de comisiones de materias de cuatrimestre %vQ%v",
				n,
				c.Numero,
				c.Anio,
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
