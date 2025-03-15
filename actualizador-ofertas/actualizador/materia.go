package actualizador

import (
	"context"
	"errors"
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

type ultimaComision struct {
	materia materia
	cuatri  cuatri
}

// updateCodigosMaterias sincroniza los códigos de las materias en la
// base de datos con sus códigos correctos obtenidos del SIU.
func updateCodigosMaterias(db *pgxpool.Pool, coms []ultimaComision) error {
	lg := log.Default().WithPrefix("🔢")

	if n, err := getCantMateriasDesactualizadas(db, lg); err != nil {
		return errors.New("error determinando la cantidad de materias sin actualizar")
	} else if n == 0 {
		lg.Info("no se encontraron materias con códigos sin actualizar")
		return nil
	} else {
		lg.Debugf("encontradas %v materias con códigos sin actualizar", n)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	tx, err := db.Begin(ctx)
	if err != nil {
		lg.Error(err)
		return errors.New("error iniciando transacción SQL de actualización de códigos")
	}

	if err := createTablaCodigos(lg, tx); err != nil {
		return errors.New("error creando tabla de asociación de códigos de materias")
	}

	if err := asociarCodigosActualesSiu(db, tx, lg, coms); err != nil {
		return errors.New("error sincronizando códigos de materias")
	}

	n, err := updateCodigosActuales(tx, lg)
	if err != nil {
		return errors.New("error actualizando códigos de materia")
	}

	if err := deleteTablaCodigos(tx); err != nil {
		lg.Error("error eliminando tabla de asociación de códigos de materias", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		lg.Error(err)
		return errors.New("error commiteando transacción SQL de actualización de códigos")
	}

	// Que no se hayan actualizado los códigos de ninguna materia de las que
	// estaban pendientes no es necesariamente un error, sino que a veces hay
	// cuatrimestres en los que no hay comisiones para algunas materias, por lo
	// que ni siquiera aparecen en el SIU.

	lg.Infof("actualizado los códigos de %v materias exitosamente", n)

	return nil
}

// getCantMateriasDesactualizadas retorna la cantidad de materias cuyos códigos
// no han sido sincronizados con los códigos correctos del SIU.
func getCantMateriasDesactualizadas(db *pgxpool.Pool, lg *log.Logger) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var n int

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
		lg.Error(err)
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

// asociarCodigosActualesSiu completa la tabla de asociación de códigos actuales de las
// materias con los códigos correctos obtenidos desde el SIU.
func asociarCodigosActualesSiu(db *pgxpool.Pool, tx pgx.Tx, lg *log.Logger, coms []ultimaComision) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
// 	defer cancel()
//
// 	lg.Debug("obteniendo códigos de materias de la base de datos")
//
// 	rows, err := db.Query(ctx, `
// SELECT m.codigo, lower(unaccent(m.nombre))
// FROM materia m
// INNER JOIN plan_materia pm
// ON m.codigo = pm.codigo_materia
// INNER JOIN plan p
// ON p.codigo = pm.codigo_plan
// WHERE p.esta_vigente = true
// 		`)
// 	if err != nil {
// 		lg.Error(err)
// 		return err
// 	}
//
// 	codigosMaterias := make(map[string]string)
//
// 	for rows.Next() {
// 		var cod, nombre string
//
// 		err := rows.Scan(&cod, &nombre)
// 		if err != nil {
// 			lg.Error("error serializando las materias", "error", err)
// 			return err
// 		}
//
// 		codigosMaterias[nombre] = cod
// 	}
//
// 	lg.Debugf("encontrados los códigos de %v materias en la base de datos", len(codigosMaterias))
//
// 	materias := make(map[string][]any, len(codigosMaterias))
//
// 	for _, c := range coms {
// 		for _, m := range c.materias {
// 			if codActual, ok := codigosMaterias[m.Nombre]; ok {
// 				if _, ok := materias[m.Nombre]; !ok {
// 					materias[m.Nombre] = []any{m.Nombre, codActual, m.Codigo}
// 				}
// 			} else {
// 				lg.Warn("materia no está en la base de datos",
// 					"materia", m.Codigo, "nombre", m.Nombre)
// 			}
// 		}
// 	}
//
// 	lg.Debugf("obtenidos los códigos de %v materias desde el SIU", len(materias))
//
// 	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
// 	defer cancel()
//
// 	_, err = tx.CopyFrom(
// 		ctx,
// 		pgx.Identifier{"tmp_codigos_materias"},
// 		[]string{"nombre_materia", "codigo_materia_actual", "codigo_materia_siu"},
// 		pgx.CopyFromRows(slices.Collect(maps.Values(materias))),
// 	)
//
// 	if err != nil {
// 		lg.Error(err)
// 		return err
// 	}
//
// 	return nil
	return nil
}

// updateCodigosActuales efectúa la actualización de los códigos de las
// materias con código desactualizado. Retorna la cantidad de registros que
// fueron afectados por la query de actualización.
func updateCodigosActuales(tx pgx.Tx, lg *log.Logger) (int, error) {
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
		lg.Error(err)
		return -1, errors.New("error actualizando códigos de materias")
	}

	return int(rows.RowsAffected()), nil
}

// filtrarUltimasComisiones se queda con la oferta de comisiones más reciente
// para cada materia. Por ejemplo, si la oferta de comisiones de Ingeniería
// Química está actualizada al 1C 2025 y la de Ingeniería Informática al 2C
// 2024, para una materia en común como podría ser Álgebra Lineal, presente en
// ambas ofertas, esta función retorna solamente la oferta de comisiones de
// Álgebra Lineal del 1C 2025.
func filtrarUltimasComisiones(ofertas []*oferta) []ultimaComision {
	max := 0
	for _, o := range ofertas {
		max += len(o.materias)
	}

	cuatris := make(map[string]cuatri, max)
	mats := make(map[string]ultimaComision, max)

	for _, o := range ofertas {
		for _, m := range o.materias {
			c, ok := cuatris[m.Nombre]

			if !ok || o.cuatri.esDespuesDe(c) {
				cuatris[m.Nombre] = o.cuatri
				mats[m.Nombre] = ultimaComision{
					materia: m,
					cuatri:  o.cuatri,
				}
			}
		}
	}

	return slices.Collect(maps.Values(mats))
}
