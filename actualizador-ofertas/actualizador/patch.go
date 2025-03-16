package actualizador

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PatchMateriaOutput struct {
	Codigo   string
	Nombre   string
	Docentes PatchDocentesOutput
	Catedras PatchCatedrasOutput
}

func getPatchesMateriaOutput(logger *log.Logger, db *pgxpool.Pool, ofertas []oferta) ([]PatchMateriaOutput, error) {
	// Para generar los patches de las materias que requieren actualización,
	// primero nos aseguramos de tener los códigos correctamente sincronizados
	// para generar los patches con la información más reciente. En caso de que
	// falle la generación de patches, igualmente ya quedan actualizados los
	// códigos para la siguiente ejecución.

	if err := syncCodigosMaterias(logger, db, ofertas); err != nil {
		return nil, err
	}

	mats, err := getMateriasConPosibleActualizacion(logger, db)
	if err != nil {
		return nil, err
	}

	// Para actualizar las comisiones de las materias tomamos
	// únicamente la última versión disponible de las ofertas de
	// cada materia.

	uofs := filtrarUltimasOfertas(ofertas)

	var wg sync.WaitGroup
	semch := make(chan struct{}, int(db.Config().MaxConns))
	patchch := make(chan PatchMateriaOutput, len(uofs))

	for _, of := range uofs {
		if _, ok := mats[of.materia.Codigo]; ok {
			logger := logger.With("codigoMateria", of.materia.Codigo)

			wg.Add(1)
			go func() {
				defer wg.Done()
				semch <- struct{}{}
				if p, done, err := newPatchMateriaOutput(logger, db, of); err != nil {
					logger.Warn("saltando actualización de materia")
				} else if !done {
					patchch <- p
				}
				<-semch
			}()
		}
	}

	go func() {
		wg.Wait()
		close(patchch)
	}()

	ps := make([]PatchMateriaOutput, len(patchch))
	for p := range patchch {
		ps = append(ps, p)
	}

	return ps, nil
}

// getMateriasConPosibleActualizacion retorna un hashset con los códigos de las
// materias cuyas cátedras no han sido actualizadas al último cuatrimestre.
// Esto no significa que la materia esté desactualizada, sino que podría haber
// alguna oferta de comisiones del SIU que sea más reciente, ya que la materia
// no ha sido actualizada este cuatrimestre.
func getMateriasConPosibleActualizacion(logger *log.Logger, db *pgxpool.Pool) (map[string]bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	// Solo las materias de los nuevos planes cuyo código ya ha sido
	// actualizado, ya que si el código no ha sido actualizado, es porque no
	// estaba incluido en ninguna de las ofertas descargadas del SIU
	// anteriormente, por lo tanto no tiene actualización posible.
	//
	// Hay más información sobre porqué el prefijo 'COD' determina si el código
	// ya fue actualizado o no en la función [getCantidadMateriasNoSync].

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
		return nil, err
	}

	cods := make(map[string]bool)
	for rows.Next() {
		var cod string
		err := rows.Scan(&cod)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		cods[cod] = true
	}

	logger.Debug(fmt.Sprintf("encontradas %v materias con ofertas posiblemente desactualizadas", len(cods)))

	return cods, nil
}

// newPatchMateriaOutput retorna un patch de materia con los patches de
// actualización de los docentes y las cátedras de la materia, pero solo si la
// materia no está actualizada ya. En este último caso (si ya está actualizada)
// retorna true como segundo valor, por lo que el patch no va a ser utilizado.
// En caso de que el patch si deba ser utilizado (si la materia ya está
// desactualizada) retorna false.
func newPatchMateriaOutput(logger *log.Logger, db *pgxpool.Pool, uoferta ultimaOfertaMateria) (PatchMateriaOutput, bool, error) {
	var p PatchMateriaOutput

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	conn, err := db.Acquire(ctx)
	if err != nil {
		msg := "error adquiriendo conexión de la pool de la base de datos"
		return p, true, logErrRetMsg(logger, msg, err)
	}
	defer conn.Release()

	if yes, err := checkMateriaYaActualizada(logger, conn, uoferta); err != nil {
		return p, true, err
	} else if yes {
		logger.Debug("comisiones de materia ya actualizadas")
		return p, true, nil
	}

	pds, err := newPatchDocentes(logger, conn, uoferta)
	if err != nil {
		return p, true, err
	}
	pcs, err := newPatchCatedras(logger, conn, uoferta)
	if err != nil {
		return p, true, err
	}

	p = PatchMateriaOutput{
		Codigo:   uoferta.Codigo,
		Nombre:   uoferta.Nombre,
		Docentes: pds,
		Catedras: pcs,
	}

	return p, false, nil
}

// checkMateriaYaActualizada verifica si la materia ya fue actualizada.
func checkMateriaYaActualizada(logger *log.Logger, conn *pgxpool.Conn, uoferta ultimaOfertaMateria) (bool, error) {
	var done bool
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	err := conn.QueryRow(ctx, `
	SELECT EXISTS (
	    SELECT 1
	    FROM actualizacion_catedras ac
	    JOIN materia m ON ac.codigo_materia = m.codigo
	    JOIN cuatrimestre c ON ac.codigo_cuatrimestre = c.codigo
	    WHERE m.codigo = $1
	    AND c.numero = $2
	    AND c.anio = $3
	);
			`, uoferta.Codigo, uoferta.cuatri.anio, uoferta.cuatri.numero).
		Scan(&done)
	if err != nil {
		msg := "error determinando si oferta de materia ya está actualizada"
		return false, logErrRetMsg(logger, msg, err)
	}

	return done, nil
}
