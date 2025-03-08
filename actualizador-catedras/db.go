package main

import (
	"context"
	"errors"
	"maps"
	"os"
	"slices"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

var db *pgxpool.Pool

func InitDBPool(logger *log.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var err error
	db, err = pgxpool.New(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		logger.Error(err)
		return errors.New("error estableciendo conexión con la base de datos")
	}

	logger.Debug("establecida conexión con la base de datos")

	return nil
}

// ActualizarCodigosMaterias sincroniza los códigos de las materias en la
// base de datos con sus códigos correctos obtenidos del SIU.
func ActualizarCodigosMaterias(ofertas []oferta) error {
	logger := log.Default().WithPrefix("🛢️")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var n int

	err := db.QueryRow(ctx, `
SELECT count(*) FROM materia WHERE codigo LIKE 'COD%'
		`).Scan(&n)

	if err != nil {
		log.Error(err)
		return errors.New("error determinando la cantidad de materias sin actualizar")
	}

	if n == 0 {
		logger.Info("no se encontraron materias con códigos sin actualizar")
		return nil
	}

	logger.Debugf("encontradas %v materias con códigos sin actualizar", n)

	codigosMaterias, err := getCodigosDeMaterias()
	if err != nil {
		logger.Fatal(err)
		return errors.New("error obteniendo los códigos de las materias de la base de datos")
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	tx, err := db.Begin(ctx)
	if err != nil {
		logger.Error(err)
		return errors.New("error iniciando transacción SQL de actualización de códigos")
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

	logger.Debugf("obtenidos los códigos de %v materias desde el SIU", len(rows))

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err = tx.CopyFrom(
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

	logger.Debug("actualizando los códigos de las materias")

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
		return errors.New("error commiteando transacción SQL de actualización de códigos")
	}

	logger.Infof("actualizados los códigos de %v materias", updateRes.RowsAffected())

	return nil
}

// getCodigosDeMaterias retorna un hashmap donde la llave es el nombre de una
// materia sin diacríticos y en minúscula, y el valor es el código que tiene
// esa materia en la base de datos.
func getCodigosDeMaterias() (map[string]string, error) {
	logger := log.Default().WithPrefix("🛢️")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logger.Debug("obteniendo los códigos de las materias")

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
		return nil, errors.New("error obteniendo los códigos de las materias")
	}

	codigos := make(map[string]string)

	for rows.Next() {
		var cod, nombre string

		err := rows.Scan(&cod, &nombre)
		if err != nil {
			logger.Error(err, "codigo", cod, "nombre", nombre)
			return nil, errors.New("error serializando las materias")
		}

		codigos[nombre] = cod
	}

	logger.Infof("encontrado los códigos de %v materias", len(codigos))

	return codigos, nil
}

// getMateriasNoActualizadasEnCuatriActual retorna un hashset con los códigos
// de las materias que no han sido actualizadas en el cuatrimestre actual. Esto
// no implica que la materia esté desactualizada realmente, ya que esto depende
// de que haya una nueva oferta de comisiones disponible en el SIU, pero
// sugiere que la materia podría esta desactualizada si dicha oferta nueva está
// disponible.
func getMateriasNoActualizadasEnCuatriActual() (map[string]bool, error) {
	logger := log.Default().WithPrefix("🛢️")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := db.Query(ctx, `
SELECT m.codigo
FROM materia m
WHERE EXISTS (
    SELECT 1
    FROM plan_materia pm
    JOIN plan p ON pm.codigo_plan = p.codigo
    WHERE pm.codigo_materia = m.codigo
    AND p.esta_vigente = TRUE
)
AND m.codigo NOT LIKE 'COD%'
AND NOT EXISTS (
    SELECT 1
    FROM actualizacion_catedras ac
    WHERE ac.codigo_materia = m.codigo
    AND ac.codigo_cuatrimestre = (SELECT MAX(codigo) FROM cuatrimestre)
)
		`)

	if err != nil {
		logger.Error(err)
		return nil, errors.New("error obteniendo los códigos de mas materias con cátedras posiblemente desactualizadas")
	}

	codigos := []string{}

	for rows.Next() {
		var cod string

		err := rows.Scan(&cod)
		if err != nil {
			return nil, errors.New("error serializando las materias")
		}

		codigos = append(codigos, cod)
	}

	logger.Infof("encontradas %v materias que pueden requerir actualización", len(codigos))

	codigosSet := make(map[string]bool, len(codigos))
	for _, c := range codigos {
		codigosSet[c] = true
	}

	return codigosSet, nil
}

func prepActualizacionesAUltimaOfertaDeComisiones(ofertas []oferta) error {
	logger := log.Default().WithPrefix("🔄")

	materiasNoActualizadas, err := getMateriasNoActualizadasEnCuatriActual()
	if err != nil {
		return err
	}

	ultimasComisiones := filtrarUltimasComisiones(ofertas)

	var eg errgroup.Group
	eg.SetLimit(int(db.Config().MaxConns))

	for _, uc := range ultimasComisiones {
		if _, ok := materiasNoActualizadas[uc.materia.Codigo]; ok {
			if actualizacionDisponible, err := hayActualizacionDisponible(uc); err != nil {
				logger.Warn("salteando actualización de materia", "codigo", uc.materia.Codigo)
				continue
			} else if !actualizacionDisponible {
				logger.Debug("materia no requiere actualización", "codigo", uc.materia.Codigo)
				continue
			}

			logger.Debug("actualización disponible para materia", "codigo", uc.materia.Codigo)

			eg.Go(func() error {
				return prepActualizacionMateria(uc)
			})
		}
	}

	if err := eg.Wait(); err != nil {
		logger.Error(err)
		return errors.New("error actualizando las materias a las últimas comisiones")
	}

	return nil
}

// hayActualizacionDisponible retorna true si la última actualización de la
// oferta de comisiones de la materia no coincide con el cuatrimestre de la
// última oferta disponible en el SIU. Es decir, si hay una oferta de
// comisiones más reciente para la materia.
func hayActualizacionDisponible(uc ultimaComision) (bool, error) {
	logger := log.Default().WithPrefix("🛢️")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	materiaYaActualizada := true

	err := db.QueryRow(ctx, `
SELECT EXISTS (
    SELECT 1
    FROM actualizacion_catedras ac
    JOIN materia m ON ac.codigo_materia = m.codigo
    JOIN cuatrimestre c ON ac.codigo_cuatrimestre = c.codigo
    WHERE m.codigo = $1
    AND c.numero = $2
    AND c.anio = $3
);
		`, uc.materia.Codigo, uc.cuatri.anio, uc.cuatri.numero).
		Scan(&materiaYaActualizada)

	if err != nil {
		logger.Warn(err)
	}

	return !materiaYaActualizada, err
}

func prepActualizacionMateria(uc ultimaComision) error {
	logger := log.Default().WithPrefix("🔄").With("codigo", uc.materia.Codigo)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conn, err := db.Acquire(ctx)
	if err != nil {
		logger.Error(err)
		return errors.New("error obteniendo conexión de la pool")
	}

	defer conn.Release()

	err = prepActualizacionDocentesDeMateria(conn, uc)
	if err != nil {
		logger.Error(err)
		return errors.New("error actualizando los docentes de la materia")
	}

	err = prepActualizacionCatedrasDeMateria(conn, uc)
	if err != nil {
		logger.Error(err)
		return errors.New("error actualizando las cátedras de la materia")
	}

	return nil
}

func prepActualizacionDocentesDeMateria(conn *pgxpool.Conn, uc ultimaComision) error {
	return nil
}

func prepActualizacionCatedrasDeMateria(conn *pgxpool.Conn, uc ultimaComision) error {
	return nil
}
