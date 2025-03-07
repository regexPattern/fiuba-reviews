package main

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

func initDBConn(logger *log.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var err error
	conn, err = pgx.Connect(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		logger.Error(err)
		return errors.New("error estableciendo conexión con la base de datos")
	}

	logger.Debug("establecida conexión con la base de datos")

	return nil
}

func selectNombresACodigosMaterias() (map[string]string, error) {
	logger := log.Default().WithPrefix("🛢️")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logger.Debug("obteniendo los códigos de las materias")

	rows, err := conn.Query(ctx, `
SELECT m.codigo, lower(unaccent(m.nombre))
FROM materia m
INNER JOIN plan_materia pm
ON m.codigo = pm.codigo_materia
INNER JOIN plan p
ON p.codigo = pm.codigo_plan
WHERE p.esta_vigente = true
		`)

	if err != nil {
		logger.Error(err)
		return nil, errors.New("error obteniendo los códigos de las materias")
	}

	codigos := make(map[string]string, 0)

	for rows.Next() {
		var cod, nombre string

		err := rows.Scan(&cod, &nombre)
		if err != nil {
			logger.Error(err, "codigo", cod, "nombre", nombre)
			return nil, errors.New("error serializando las materias")
		}

		codigos[nombre] = cod
	}

	logger.Info(fmt.Sprintf("encontrado los códigos de %v materias", len(codigos)))

	return codigos, nil
}

func syncCodigosMaterias(codigosMaterias map[string]string, ofertas []ofertaComisiones) error {
	logger := log.Default().WithPrefix("🛢️")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var n int

	err := conn.QueryRow(ctx, `SELECT count(*) FROM materia WHERE codigo LIKE 'COD%'`).Scan(&n)
	if err != nil {
		log.Error(err)
		return errors.New("error determinando la cantidad de materias sin actualizar")
	}

	logger.Debug(fmt.Sprintf("encontradas %v materias con códigos sin actualizar", n))

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	tx, err := conn.Begin(ctx)
	if err != nil {
		logger.Error(err)
		return errors.New("error iniciando transacción SQL de actualización")
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err = tx.Exec(ctx, `
CREATE TABLE tmp_codigos_materias (
	nombre_materia TEXT PRIMARY KEY,
	codigo_materia_actual TEXT NOT NULL,
	codigo_materia_siu TEXT NOT NULL
)
		`)

	if err != nil {
		logger.Error(err)
		return errors.New("error creando tabla temporal de asociación de códigos de materias")
	}

	rows := make(map[string][]any, len(codigosMaterias))
	for _, o := range ofertas {
		for _, m := range o.materias {
			if codActual, ok := codigosMaterias[m.Nombre]; ok {
				if _, ok := rows[m.Nombre]; !ok {
					rows[m.Nombre] = []any{m.Nombre, codActual, m.Codigo}
				}
			} else {
				log.Default().WithPrefix("🔎").Warn("materia no está en la base de datos", "codigo", m.Codigo, "nombre", m.Nombre)
			}
		}
	}

	logger.Debug(fmt.Sprintf("relacionados los códigos de %v materias", len(rows)))

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err = conn.CopyFrom(
		ctx,
		pgx.Identifier{"tmp_codigos_materias"},
		[]string{"nombre_materia", "codigo_materia_actual", "codigo_materia_siu"},
		pgx.CopyFromRows(slices.Collect(maps.Values(rows))),
	)

	if err != nil {
		logger.Error(err)
		return errors.New("error insertando códigos correctos de mateiras obtenidos del SIU")
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logger.Info("actualizando los códigos de las materias")

	updateRes, err := tx.Exec(ctx, `
WITH materias_a_actualizar AS (
    SELECT m.codigo as codigo_materia_actual, tcm.codigo_materia_siu
    FROM materia m
    JOIN tmp_codigos_materias tcm ON lower(unaccent(m.nombre)) = tcm.nombre_materia
    JOIN plan_materia pm ON m.codigo = pm.codigo_materia
    JOIN plan p ON pm.codigo_plan = p.codigo
    WHERE p.esta_vigente = TRUE
    AND tcm.codigo_materia_actual != tcm.codigo_materia_siu
)
UPDATE materia
SET codigo = ma.codigo_materia_siu
FROM materias_a_actualizar ma
WHERE materia.codigo = ma.codigo_materia_actual
		`)

	if err != nil {
		logger.Error(err)
		return errors.New("error actualizando códigos de materias")
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err = tx.Exec(ctx, `DROP TABLE tmp_codigos_materias`)
	if err != nil {
		logger.Error(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err = tx.Commit(ctx)
	if err != nil {
		logger.Error(err)
		return errors.New("error commiteando transacción SQL de actualización")
	}

	logger.Info(fmt.Sprintf("actualizados los códigos de %v materias", updateRes.RowsAffected()))

	return nil
}

// TODO:
func updateCatedrasMaterias() error {
	_, err := selectCodigosMateriasDesactualizadas()
	if err != nil {
		return err
	}

	return nil
}

// TODO:
func selectCodigosMateriasDesactualizadas() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, _ = conn.Query(ctx, `
SELECT m.codigo, m.nombre
FROM materia m
WHERE EXISTS (
    SELECT 1
    FROM plan_materia pm
    JOIN plan p ON pm.codigo_plan = p.codigo
    WHERE pm.codigo_materia = m.codigo
    AND p.esta_vigente = TRUE
)
AND NOT EXISTS (
    SELECT 1
    FROM actualizacion_catedras ac
    WHERE ac.codigo_materia = m.codigo
    AND ac.codigo_cuatrimestre = (SELECT MAX(codigo) FROM cuatrimestre)
)
		`)

	return nil, nil
}
