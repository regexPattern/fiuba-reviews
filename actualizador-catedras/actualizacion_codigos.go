package main

import (
	"context"
	"errors"
	"maps"
	"slices"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
)

// ActualizarCodigosMaterias sincroniza los códigos de las materias en la
// base de datos con sus códigos correctos obtenidos del SIU.
func ActualizarCodigosMaterias(ofertas []oferta) error {
	logger := log.Default().WithPrefix("🛢️")

	if n, err := getCantMateriasDesactualizadas(logger); err != nil {
		return errors.New("error determinando la cantidad de materias sin actualizar")
	} else if n == 0 {
		logger.Info("no se encontraron materias con códigos sin actualizar")
		return nil
	} else {
		logger.Debugf("encontradas %v materias con códigos sin actualizar", n)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	tx, err := db.Begin(ctx)
	if err != nil {
		logger.Error(err)
		return errors.New("error iniciando transacción SQL de actualización de códigos")
	}

	if err := createTablaCodigos(logger, tx); err != nil {
		return errors.New("error creando tabla de asociación de códigos de materias")
	}

	if err := asociarCodigos(logger, tx, ofertas); err != nil {
		return errors.New("error asociando códigos de materias")
	}

	n, err := updateCodigosActuales(logger, tx)
	if err != nil {
		return errors.New("")
	}

	if err := deleteTablaCodigos(tx); err != nil {
		logger.Error("error eliminando tabla de asociación de códigos de materias", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.Error(err)
		return errors.New("error commiteando transacción SQL de actualización de códigos")
	}

	// INFO: Que no se hayan actualizados los códigos de ninguna materia de las
	// que estaban pendientes no es necesariamente un error, sino que a veces
	// hay cuatrimestres en los que no hay comisiones para algunas materias,
	// por lo que ni siquiera aparecen en el SIU.
	logger.Infof("actualizados los códigos de %v materias", n)

	return nil
}

// getCantMateriasDesactualizadas retorna la cantidad de materias cuyos códigos
// no han sido sincronizados con los códigos correctos del SIU.
func getCantMateriasDesactualizadas(logger *log.Logger) (int, error) {
	var n int

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// Cuando se crearon las masterias de los nuevos planes en FIUBA Reviews,
	// no se disponía de una fuente de información oficial de la cuál obtener
	// los nuevos códigos, por lo que se generaron códigos placeholder, que son
	// los que inician con el prefijo 'COD'.
	//
	// Si una materia aún tiene un código con este prefijo es porque su código
	// no ha sido reemplazado por el código oficial obtenido desde el SIU en
	// ejecuciones previas de esta utilidad.
	err := db.QueryRow(ctx, `
SELECT count(*) FROM materia WHERE codigo LIKE 'COD%'
		`).Scan(&n)

	if err != nil {
		logger.Error(err)
	}

	return n, err
}

// createTablaCodigos crea la tabla SQL para asociar los códigos actuales de
// las materias con los códigos correctos obtenidos desde el SIU.
func createTablaCodigos(logger *log.Logger, tx pgx.Tx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err := tx.Exec(ctx, `
	CREATE TABLE tmp_codigos_materias (
		nombre_materia TEXT PRIMARY KEY,
		codigo_materia_actual TEXT NOT NULL,
		codigo_materia_siu TEXT NOT NULL
	)
			`)

	if err != nil {
		logger.Error(err)
	}

	return err
}

// deleteTablaCodigos elimina la tabla SQL para asociar los códigos de las
// materias. En caso de error no hay mucho problema porque igual es una tabla
// temporal que se borra al final de la transacción. De lo único que habría que
// cuidarse es del caso en que no se cierre la transacción.
func deleteTablaCodigos(tx pgx.Tx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err := tx.Exec(ctx, `DROP TABLE tmp_codigos_materias`)

	return err
}

// asociarCodigos completa la tabla de asociación de códigos actuales de las
// materias con los códigos correctos obtenidos desde el SIU.
func asociarCodigos(logger *log.Logger, tx pgx.Tx, ofertas []oferta) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logger.Debug("obteniendo códigos de materias")

	rows, err := db.Query(ctx, `
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
		return err
	}

	codigosMaterias := make(map[string]string)

	for rows.Next() {
		var cod, nombre string

		err := rows.Scan(&cod, &nombre)
		if err != nil {
			logger.Error("error serializando las materias",
				"error", err, "codigo", cod, "nombre", nombre)
			return err
		}

		codigosMaterias[nombre] = cod
	}

	logger.Infof("encontrado los códigos de %v materias", len(codigosMaterias))

	materias := make(map[string][]any, len(codigosMaterias))

	materiaFaltanteLogger := log.Default().WithPrefix("🔎")

	for _, o := range ofertas {
		for _, m := range o.materias {
			if codActual, ok := codigosMaterias[m.Nombre]; ok {
				if _, ok := materias[m.Nombre]; !ok {
					materias[m.Nombre] = []any{m.Nombre, codActual, m.Codigo}
				}
			} else {
				materiaFaltanteLogger.Warn("materia no está en la base de datos",
					"codigo", m.Codigo, "nombre", m.Nombre)
			}
		}
	}

	logger.Debugf("obtenidos los códigos de %v materias desde el SIU", len(materias))

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"tmp_codigos_materias"},
		[]string{"nombre_materia", "codigo_materia_actual", "codigo_materia_siu"},
		pgx.CopyFromRows(slices.Collect(maps.Values(materias))),
	)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// updateCodigosActuales efectúa la actualización de los códigos de las
// materias con código desactualizado. Retorna la cantidad de registros que
// fueron afectados por la query de actualización.
func updateCodigosActuales(logger *log.Logger, tx pgx.Tx) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logger.Debug("actualizando los códigos de las materias")

	rows, err := tx.Exec(ctx, `
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
		return -1, errors.New("error actualizando códigos de materias")
	}

	return int(rows.RowsAffected()), nil
}
