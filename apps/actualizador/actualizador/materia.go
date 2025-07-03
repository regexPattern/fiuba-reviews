package actualizador

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type materia struct {
	Codigo   string    `json:"codigo"`
	Nombre   string    `json:"nombre"`
	Catedras []catedra `json:"catedras"`
}

type catedra struct {
	Codigo   int       `json:"codigo"`
	Docentes []docente `json:"docentes"`
}

type docente struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

type ultimaOfertaMateria struct {
	materia
	cuatri cuatri
}

func syncCodigosMaterias(logger *log.Logger, db *pgxpool.Pool, ofertas []oferta) error {
	if n, err := getCantidadMateriasNoSync(logger, db); err != nil {
		return err
	} else if n == 0 {
		logger.Info("no se encontraron materias con códigos sin sincronizar")
		return nil
	} else {
		logger.Info(fmt.Sprintf("encontradas %v materias con códigos sin sincronizar", n))
	}

	cods, err := getCodigosActualesMaterias(logger, db)
	if err != nil {
		return err
	}

	logger.Debug(fmt.Sprintf("encontrados %v códigos de materias desde el SIU", len(cods)))
	rows := getCodigosActualesASiuRows(logger, cods, ofertas)

	tx, err := beginTxSync(logger, db)
	if err != nil {
		return err
	}
	n, err := applyModsTxSync(logger, tx, rows)
	if err != nil {
		return err
	}
	if err := commitTxSync(logger, tx); err != nil {
		return err
	}

	// Que no se hayan sincronizado los códigos de ninguna materia de las que
	// estaban pendientes no es necesariamente un error, sino que a veces hay
	// cuatrimestres en los que no hay cátedras para algunas materias, por lo
	// que ni siquiera aparecen en el SIU.

	logger.Info(fmt.Sprintf("sincronizado códigos de %v materias", n))
	return nil
}

func getCantidadMateriasNoSync(logger *log.Logger, db *pgxpool.Pool) (int, error) {
	var n int
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// Cuando se crearon las materias de los nuevos planes en FIUBA Reviews,
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
		msg := "error obteniendo cantidad de materias con códigos no sincronizados"
		return 0, logErrRetMsg(logger, msg, err)
	}

	return n, nil
}

func getCodigosActualesMaterias(logger *log.Logger, db *pgxpool.Pool) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

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
		msg := "error obteniendo códigos actuales de materias"
		return nil, logErrRetMsg(logger, msg, err)
	}

	cods := make(map[string]string)
	for rows.Next() {
		var cod, nombre string
		err := rows.Scan(&cod, &nombre)
		if err != nil {
			msg := "error serializando materias"
			return nil, logErrRetMsg(logger, msg, err)
		}
		cods[nombre] = cod
	}

	return cods, nil
}

func beginTxSync(logger *log.Logger, db *pgxpool.Pool) (pgx.Tx, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	tx, err := db.Begin(ctx)
	if err != nil {
		msg := "error iniciando transacción SQL de actualización de códigos"
		return nil, logErrRetMsg(logger, msg, err)
	}

	return tx, nil
}

func applyModsTxSync(logger *log.Logger, tx pgx.Tx, rows map[string][]any) (int, error) {
	if err := createTblCodigos(logger, tx); err != nil {
		return 0, err
	}
	if err := insertCodigos(logger, tx, rows); err != nil {
		return 0, err
	}
	n, err := updateCodigos(logger, tx)
	if err != nil {
		return 0, err
	}
	if err := dropTblCodigos(logger, tx); err != nil {
		return 0, err
	}

	return n, nil
}

func createTblCodigos(logger *log.Logger, tx pgx.Tx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	_, err := tx.Exec(ctx, `
	CREATE TABLE tmp_codigos_materias (
		nombre_materia TEXT PRIMARY KEY,
		codigo_materia_actual TEXT NOT NULL,
		codigo_materia_siu TEXT NOT NULL
	)
			`)
	if err != nil {
		msg := "error creando tabla de asociación de códigos de materias"
		return logErrRetMsg(logger, msg, err)
	}

	return nil
}

func getCodigosActualesASiuRows(logger *log.Logger, cods map[string]string, ofertas []oferta) map[string][]any {
	matches := make(map[string][]any, len(cods))
	for _, of := range ofertas {
		for _, m := range of.materias {
			if codActual, ok := cods[m.Nombre]; ok {
				if _, ok := matches[m.Nombre]; !ok {
					matches[m.Nombre] = []any{m.Nombre, codActual, m.Codigo}
				}
			} else {
				logger.Warn("materia no está en la base de datos",
					"codigoMateria", m.Codigo, "nombreMateria", m.Nombre)
			}
		}
	}

	return matches
}

func insertCodigos(logger *log.Logger, tx pgx.Tx, rows map[string][]any) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	_, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"tmp_codigos_materias"},
		[]string{"nombre_materia", "codigo_materia_actual", "codigo_materia_siu"},
		pgx.CopyFromRows(slices.Collect(maps.Values(rows))),
	)
	if err != nil {
		msg := "error copiando filas de asociación de códigos de materias"
		return logErrRetMsg(logger, msg, err)
	}

	return nil
}

func updateCodigos(logger *log.Logger, tx pgx.Tx) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

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
		msg := "error actualizando códigos de materias"
		return 0, logErrRetMsg(logger, msg, err)
	}

	return int(rows.RowsAffected()), nil
}

func dropTblCodigos(logger *log.Logger, tx pgx.Tx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	if _, err := tx.Exec(ctx, `DROP TABLE tmp_codigos_materias`); err != nil {
		msg := "error eliminando tabla de asociación de códigos de materias"
		return logErrRetMsg(logger, msg, err)
	}

	return nil
}

func commitTxSync(logger *log.Logger, tx pgx.Tx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	if err := tx.Commit(ctx); err != nil {
		msg := "error commiteando transacción SQL de actualización de códigos"
		return logErrRetMsg(logger, msg, err)
	}

	return nil
}

func filtrarUltimasOfertas(ofertas []oferta) []ultimaOfertaMateria {
	max := 0
	for _, of := range ofertas {
		max += len(of.materias)
	}

	cuatris := make(map[string]cuatri, max)
	mats := make(map[string]ultimaOfertaMateria, max)

	for _, of := range ofertas {
		for _, m := range of.materias {
			c, ok := cuatris[m.Nombre]

			if !ok || of.cuatri.esDespuesDe(c) {
				cuatris[m.Nombre] = of.cuatri
				mats[m.Nombre] = ultimaOfertaMateria{
					materia: m,
					cuatri:  of.cuatri,
				}
			}
		}
	}

	return slices.Collect(maps.Values(mats))
}
