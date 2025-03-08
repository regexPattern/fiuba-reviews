package main

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/charmbracelet/log"
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
