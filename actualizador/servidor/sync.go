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
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		return fmt.Errorf("error iniciando transacción de sincronización de materias: %w", err)
	}

	defer func() {
		if err := tx.Rollback(context.TODO()); err != nil {
			slog.Error(
				fmt.Sprintf(
					"error haciendo rollback de transacción de sincronización de materias: %v", err,
				),
			)
		}
	}()

	type materiaSincronizadaRow struct {
		Codigo                 string   `db:"codigo"`
		Nombre                 string   `db:"nombre"`
		DocentesMigrados       int      `db:"docentes_migrados"`
		ComentariosMigrados    int      `db:"comentarios_migrados"`
		CalificacionesMigradas int      `db:"calificaciones_migradas"`
		CodigosEquivalencias   []string `db:"codigos_equivalencias"`
	}

	rows, err := tx.Query(context.TODO(), updateSyncMateriasQuery, nombres, codigos)
	if err != nil {
		return fmt.Errorf("error ejecutando query de sincronización de materias: %w", err)
	}

	materiasSincronizadas, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[materiaSincronizadaRow],
	)
	if err != nil {
		return fmt.Errorf("error serializando materias sincronizadas: %w", err)
	}

	for _, m := range materiasSincronizadas {
		slog.Debug(
			fmt.Sprintf("sincronizada materia %v (%v)", m.Codigo, m.Nombre),
			"docentes", m.DocentesMigrados,
			"calificaciones", m.CalificacionesMigradas,
			"comentarios", m.ComentariosMigrados,
			"equivalencias", m.CodigosEquivalencias,
		)
	}

	slog.Info(fmt.Sprintf("sincronizadas %d materias en total", len(materiasSincronizadas)))

	if err := tx.Commit(context.TODO()); err != nil {
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
		context.TODO(),
		selectMateriasNoRegistradasQuery,
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

	for _, m := range materiasNoRegistradas {
		slog.Warn(
			fmt.Sprintf(
				"materia %v (%v) no está registrada en la base de datos",
				m.Codigo,
				m.Nombre,
			),
		)
	}

	return nil
}
