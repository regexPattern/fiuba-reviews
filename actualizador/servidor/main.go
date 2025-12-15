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

	patches, err := getPatchesActualizacionMaterias(conn)
	if err != nil {
		slog.Error(fmt.Sprintf("error generando patches de actualización de materias: %v", err))
		os.Exit(1)
	}

	if err := startServer(conn, patches); err != nil {
		slog.Error(
			fmt.Sprintf(
				"error ejecutando servidor de resolución de actualización de ofertas de materias: %v",
				err,
			),
		)
		os.Exit(1)
	}
	//
	// ofertasMaterias, err := getOfertasMaterias(conn)
	// if err != nil {
	// 	slog.Error(fmt.Sprintf("error obteniendo ofertas de comisiones de materias: %v", err))
	// 	os.Exit(1)
	// }
	//
	// codigosMaterias := make([]string, 0, len(ofertasMaterias))
	// nombresMaterias := make([]string, 0, len(ofertasMaterias))
	// for codigo, oferta := range ofertasMaterias {
	// 	codigosMaterias = append(codigosMaterias, codigo)
	// 	nombresMaterias = append(nombresMaterias, oferta.Nombre)
	// }
	//
	// if err := syncMateriasDb(conn, codigosMaterias, nombresMaterias); err != nil {
	// 	slog.Error(
	// 		fmt.Sprintf("error sincronizando materias de la base de datos con el siu: %v", err),
	// 	)
	// 	os.Exit(1)
	// }
	//
	// materiasPendientes, err := getPatchesMateria(conn, codigosMaterias, ofertasMaterias)
	// if err != nil {
	// 	slog.Error(fmt.Sprintf("error obteniendo materias a actualizar: %v", err))
	// 	os.Exit(1)
	// }
	//
	// slog.Debug(
	// 	fmt.Sprintf(
	// 		"encontradas %d materias con actualizaciones pendientes",
	// 		len(materiasPendientes),
	// 	),
	// )
	//
	// // for i, mat := range materiasAActualizar {
	// // 	slog.Debug(fmt.Sprintf("materia %d: %s - %d docentes pendientes - %d catedras nuevas",
	// // 		i+1, mat.Nombre, len(mat.DocentesPendientes), len(mat.CatedrasNuevas)))
	// // }
	//
	// for _, m := range materiasPendientes {
	// 	_ = m
	// }
}

func getPatchesActualizacionMaterias(conn *pgx.Conn) ([]patchActualizacionMateria, error) {
	ofertas, err := getOfertasMaterias(conn)
	if err != nil {
		return nil, fmt.Errorf(
			"error obteniendo ofertas de comisiones de materias: %w",
			err,
		)
	}

	codigos := make([]string, 0, len(ofertas))
	nombres := make([]string, 0, len(ofertas))

	for cod, om := range ofertas {
		codigos = append(codigos, cod)
		nombres = append(nombres, om.Nombre)
	}

	if err := syncDb(conn, codigos, nombres); err != nil {
		return nil, fmt.Errorf(
			"error sincronizando materias de la base de datos con el siu: %w",
			err,
		)
	}

	patches, err := buildPatchesActualizacionMaterias(conn, codigos, ofertas)
	if err != nil {
		return nil, fmt.Errorf(
			"error construyendo patches de actualización de ofertas de comisiones de materias: %w",
			err,
		)
	}

	return patches, nil
}

func startServer(_ *pgx.Conn, _ any) error {
	return nil
}
