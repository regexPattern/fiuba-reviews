package main

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

//go:embed queries/2-UPDATE-sync-materias-query.sql
var updateSyncMateriasQuery string

//go:embed queries/3-SELECT-materias-no-registradas.sql
var selectMateriasNoRegistradasQuery string

func syncDb(conn *pgx.Conn, codigos, nombres []string) error {
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("error iniciando transacción de sincronización de materias: %w", err)
	}

	defer tx.Rollback(context.Background())

	type MateriaSincronizadaRow struct {
		Codigo                 string   `db:"codigo"`
		Nombre                 string   `db:"nombre"`
		DocentesMigrados       int      `db:"docentes_migrados"`
		ComentariosMigrados    int      `db:"comentarios_migrados"`
		CalificacionesMigradas int      `db:"calificaciones_migradas"`
		CodigosEquivalencias   []string `db:"codigos_equivalencias"`
	}

	rows, err := tx.Query(context.Background(), updateSyncMateriasQuery, nombres, codigos)
	if err != nil {
		return fmt.Errorf("error ejecutando query de sincronización de materias: %w", err)
	}

	materiasSincronizadas, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[MateriaSincronizadaRow],
	)
	if err != nil {
		return fmt.Errorf("error procesando materias sincronizadas: %w", err)
	}

	for _, m := range materiasSincronizadas {
		slog.Debug(
			fmt.Sprintf("sincronizado materia %s %s", m.Codigo, m.Nombre),
			"docentes", m.DocentesMigrados,
			"calificaciones", m.CalificacionesMigradas,
			"comentarios", m.ComentariosMigrados,
			"equivalencias", m.CodigosEquivalencias,
		)
	}

	slog.Info(fmt.Sprintf("sincronizadas %d materias en total", len(materiasSincronizadas)))

	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf(
			"error haciendo commit de la transacción de sincronización de materias: %w",
			err,
		)
	}

	if err := checkMateriasNoRegistradas(conn, codigos, nombres); err != nil {
		return fmt.Errorf("error verificando materias no registradas en la base de datos: %w", err)
	}

	return nil
}

func checkMateriasNoRegistradas(conn *pgx.Conn, codigos, nombres []string) error {
	rows, err := conn.Query(
		context.Background(),
		selectMateriasNoRegistradasQuery,
		nombres,
		codigos,
	)
	if err != nil {
		return fmt.Errorf("error consultando materias no registradas: %w", err)
	}

	materiasNoRegistradas, err := pgx.CollectRows(rows, pgx.RowToStructByName[Materia])
	if err != nil {
		return fmt.Errorf("error procesando materias no registradas: %v", err)
	}

	for _, m := range materiasNoRegistradas {
		slog.Warn(
			fmt.Sprintf("materia %v (%v) no registrada en la base de datos", m.Codigo, m.Nombre),
		)
	}

	return nil
}
