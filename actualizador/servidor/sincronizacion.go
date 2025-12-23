package main

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/regexPattern/fiuba-reviews/actualizador/queries"
)

func sincronizarMaterias(conn *pgx.Conn, codigos, nombres []string) error {
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		return fmt.Errorf("error iniciando transacción de sincronización de materias: %w", err)
	}
	defer func() { _ = tx.Rollback(context.TODO()) }()

	rows, err := tx.Query(context.TODO(), queries.SincronizarMaterias, nombres, codigos)
	if err != nil {
		return fmt.Errorf("error ejecutando query de sincronización de materias: %w", err)
	}

	type materiaSincronizadaRow struct {
		Codigo                 string   `db:"codigo"`
		Nombre                 string   `db:"nombre"`
		DocentesMigrados       int      `db:"docentes_migrados"`
		ComentariosMigrados    int      `db:"comentarios_migrados"`
		CalificacionesMigradas int      `db:"calificaciones_migradas"`
		CodigosEquivalencias   []string `db:"codigos_equivalencias"`
	}

	materiasSincronizadas, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[materiaSincronizadaRow],
	)
	if err != nil {
		return fmt.Errorf("error serializando materias sincronizadas: %w", err)
	}

	for _, mat := range materiasSincronizadas {
		slog.Debug("materia_sincronizada", "codigo_materia", mat.Codigo,
			"docentes_migrados", mat.DocentesMigrados,
			"calificaciones_migradas", mat.CalificacionesMigradas,
			"comentarios_migrados", mat.ComentariosMigrados,
			"equivalencias", mat.CodigosEquivalencias,
		)
	}

	slog.Info("materias_sincronizadas", "count", len(materiasSincronizadas))

	if err := tx.Commit(context.TODO()); err != nil {
		return fmt.Errorf(
			"error haciendo commit de la transacción de sincronización de materias: %w",
			err,
		)
	}

	if err := checkMateriasNoRegistradas(conn, codigos, nombres); err != nil {
		return fmt.Errorf("error checkeando materias no registradas en la base de datos: %w", err)
	}

	return nil
}

func checkMateriasNoRegistradas(conn *pgx.Conn, codigos, nombres []string) error {
	rows, err := conn.Query(
		context.TODO(),
		queries.MateriasNoRegistradasEnDb,
		nombres,
		codigos,
	)
	if err != nil {
		return fmt.Errorf("error consultando materias no registradas: %w", err)
	}

	materiasNoRegistradas, err := pgx.CollectRows(rows, pgx.RowToStructByName[materia])
	if err != nil {
		return fmt.Errorf("error serializando materias no registradas: %v", err)
	}

	for _, mat := range materiasNoRegistradas {
		slog.Warn("materia_no_registrada_en_db", "codigo_materia", mat.Codigo)
	}

	return nil
}
