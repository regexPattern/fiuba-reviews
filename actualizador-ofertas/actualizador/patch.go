package actualizador

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/regexPattern/fiuba-reviews/actualizador-ofertas/config"
)

type PatchMateriaOutput struct {
	Codigo   string
	Nombre   string
	Docentes *PatchDocentesOutput
	Catedras *PatchCatedrasOutput
}

func newPatchMateriaOutput(ch chan PatchMateriaOutput, uc ultimaComision) (*PatchMateriaOutput, error) {
	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// 	defer cancel()
	//
	// 	conn, err := db.Acquire(ctx)
	// 	if err != nil {
	// 		logger.Error(err)
	// 		return err
	// 	}
	//
	// 	defer conn.Release()
	//
	// 	// Acá finalmente se verifica si la materia tiene una actualización
	// 	// posible o no. Antes solo traímos las que posiblemente tenían
	// 	// actualización porque no habían sido actualizadas en el último
	// 	// cuatrimestre. Ahora se verifica si el plan más reciente para esta
	// 	// materia es el que ya está registrado.
	//
	// 	var yaEnPlanMasRecienteDisp bool
	//
	// 	err = conn.QueryRow(ctx, `
	// SELECT EXISTS (
	//     SELECT 1
	//     FROM actualizacion_catedras ac
	//     JOIN materia m ON ac.codigo_materia = m.codigo
	//     JOIN cuatrimestre c ON ac.codigo_cuatrimestre = c.codigo
	//     WHERE m.codigo = $1
	//     AND c.numero = $2
	//     AND c.anio = $3
	// );
	// 		`, uc.materia.Codigo, uc.cuatri.anio, uc.cuatri.numero).
	// 		Scan(&yaEnPlanMasRecienteDisp)
	//
	// 	if err != nil {
	// 		logger.Error(err)
	// 		return err
	// 	}
	//
	// 	if yaEnPlanMasRecienteDisp {
	// 		return nil
	// 	}
	//
	// 	patchesDocentes, err := newPatchDocentes(logger, conn, uc)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	_, err = newPatchCatedras(logger, conn, uc)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	ch <- PatchMateriaOutput{
	// 		Codigo:   uc.materia.Codigo,
	// 		Nombre:   uc.materia.Nombre,
	// 		Docentes: patchesDocentes,
	// 		Catedras: nil, // TODO
	// 	}
	//
	// 	return nil

	return nil, nil
}

// getPatchesMateriaOutput retorna los patches de actualización para las
// materias.
func getPatchesMateriaOutput(ofertas []oferta) ([]PatchMateriaOutput, error) {
	// logger := log.Default().WithPrefix("👨‍🏫")
	//
	// matsPendientes, err := getMateriasPendientes(logger)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// logger.Debug("filtrando solo las ofertas de comisiones más recientes")
	//
	// ultimasComisiones := filtrarUltimasComisiones(ofertas)
	//
	// var wg sync.WaitGroup
	// sem := make(chan struct{}, int(db.Config().MaxConns))
	//
	// ch := make(chan PatchMateriaOutput, len(ultimasComisiones))
	//
	// for _, uc := range ultimasComisiones {
	// 	if _, ok := matsPendientes[uc.materia.Codigo]; ok {
	// 		logger := logger.With("materia", uc.materia.Codigo)
	//
	// 		wg.Add(1)
	// 		go func() {
	// 			defer wg.Done()
	//
	// 			sem <- struct{}{}
	// 			defer func() { <-sem }()
	//
	// 			if newPatchMateriaOutput(logger, ch, uc) != nil {
	// 				logger.Warn("saltando actualización de materia")
	// 			}
	// 		}()
	// 	}
	// }
	//
	// go func() {
	// 	wg.Wait()
	// 	close(ch)
	// }()
	//
	// patches := make([]PatchMateriaOutput, len(ch))
	// for p := range ch {
	// 	patches = append(patches, p)
	// }
	//
	// return patches, nil

	return nil, nil
}

type PatchMateriaInput struct{}

// getMateriasPendientes retorna un hashset con los códigos de las materias
// cuyas ofertas de comisiones no han sido actualizadas al último cuatrimestre.
// Esto no significa que la materia esté desactualizada, sino que podría haber
// alguna oferta de comisiones de las del SIU que sea más reciente, ya que la
// materia no ha sido actualizada este cuatrimestre.
func getMateriasPendientes(logger *log.Logger) (map[string]bool, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
// 	defer cancel()
//
// 	// Solo las materias del nuevo plan cuyo código ya ha sido actualizado, ya
// 	// que si el código no ha sido actualizado, es porque no estaba incluído en
// 	// ninguna de las ofertas descargadas del SIU anteriormente.
//
// 	rows, err := db.Query(ctx, `
// SELECT m.codigo
// FROM materia m
// WHERE EXISTS (
//     SELECT 1
//     FROM plan_materia pm
//     JOIN plan p ON pm.codigo_plan = p.codigo
//     WHERE pm.codigo_materia = m.codigo
//     AND p.esta_vigente = TRUE
// )
// AND m.codigo NOT LIKE 'COD%'
// AND NOT EXISTS (
//     SELECT 1
//     FROM actualizacion_catedras ac
//     WHERE ac.codigo_materia = m.codigo
//     AND ac.codigo_cuatrimestre = (SELECT MAX(codigo) FROM cuatrimestre)
// )
// 		`)
//
// 	if err != nil {
// 		logger.Error(err)
// 		return nil, err
// 	}
//
// 	codigosMaterias := make(map[string]bool)
//
// 	for rows.Next() {
// 		var cod string
//
// 		err := rows.Scan(&cod)
// 		if err != nil {
// 			logger.Error(err)
// 			return nil, err
// 		}
//
// 		codigosMaterias[cod] = true
// 	}
//
// 	logger.Logf(config.DebugIndividualOps, "encontradas %v materias con ofertas posiblemente desactualizadas", len(codigosMaterias))
//
// 	return codigosMaterias, nil
	return nil, nil
}

type PatchDocentesOutput struct {
	Registrados map[string]string
	Nuevos      map[string]string
}

// newPatchDocentes retorna un patch de docentes con los docentes ya
// registrados en la base de datos que han sido relacionados con un docente del
// SIU, y un listado de los nuevos docentes del SIU que no han sido agregados a
// la base de datos.
func newPatchDocentes(logger *log.Logger, conn *pgxpool.Conn, uc ultimaComision) (*PatchDocentesOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// Si la materia nunca ha sido actualizada entonces copio los docentes de
	// sus equivalencias (además de sus calificaciones y comentarios) y los
	// asigno a la materia en cuestión.

	var primeraActualizacion bool

	err := conn.QueryRow(ctx, `
SELECT NOT docentes_migrados_de_equivalencia
FROM materia
WHERE codigo = $1
		`, uc.materia.Codigo).Scan(&primeraActualizacion)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if primeraActualizacion {
		if err := migrarDocentesEquivalencias(logger, conn, uc.materia.Codigo); err != nil {
			return nil, err
		}
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := conn.Query(ctx, `
SELECT codigo, lower(unaccent(nombre))
FROM docente
WHERE codigo_materia = $1
		`, uc.materia.Codigo)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	docentesDb := make(map[string]string)

	for rows.Next() {
		var cod, nombre string

		err := rows.Scan(&cod, &nombre)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		docentesDb[nombre] = cod
	}

	logger.Logf(config.DebugIndividualOps, "encontrados %v docentes en la base de datos", len(docentesDb))

	docentesSiu := make(map[string]string)
	for _, c := range uc.materia.Catedras {
		for _, d := range c.Docentes {
			docentesSiu[d.Nombre] = d.Rol
		}
	}

	docentesDbPendientes := make(map[string]string, len(docentesDb))
	docentesSiuPendientes := make(map[string]string, len(docentesSiu))

	for nombre, cod := range docentesDb {
		if _, ok := docentesSiu[nombre]; ok {
			continue
		}

		docentesDbPendientes[nombre] = cod
	}

	for nombre, rol := range docentesSiu {
		if _, ok := docentesDb[nombre]; ok {
			continue
		}

		docentesSiuPendientes[nombre] = rol
	}

	patch := &PatchDocentesOutput{
		Registrados: docentesDbPendientes,
		Nuevos:      docentesSiuPendientes,
	}

	return patch, nil
}

// migrarDocentesEquivalencias copia los docentes de las equivalencias de una
// materia y los asigna a la materia en cuestión, como si fueran suyos. Esto
// solo cobra importancia en el caso en el que no haya ningún docente en una
// materia del nuevo plan aún, para generar un estado inicial con los datos de
// Dolly.
func migrarDocentesEquivalencias(logger *log.Logger, conn *pgxpool.Conn, codigoMateria string) error {
	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// 	defer cancel()
	//
	// 	tx, err := conn.Begin(ctx)
	// 	if err != nil {
	// 		logger.Error(err)
	// 		return err
	// 	}
	//
	// 	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	// 	defer cancel()
	//
	// 	var n int
	//
	// 	// La query la escribió la IA, pero anda y hace justo lo que debería.
	//
	// 	err = tx.QueryRow(ctx, `
	// WITH materias_equivalentes AS (
	//     SELECT e.codigo_materia_plan_anterior AS codigo_materia_equivalente
	//     FROM equivalencia e
	//     WHERE e.codigo_materia_plan_vigente = $1
	// ),
	// docentes_equivalencias AS (
	//     SELECT
	//         d.codigo AS codigo_antiguo,
	//         gen_random_uuid() AS codigo_nuevo,
	//         d.nombre,
	//         d.resumen_comentarios,
	//         d.comentarios_ultimo_resumen
	//     FROM docente d
	//     JOIN materias_equivalentes me
	// 	ON d.codigo_materia = me.codigo_materia_equivalente
	// ),
	// docentes_copiados AS (
	//     INSERT INTO docente (codigo, nombre, codigo_materia, resumen_comentarios, comentarios_ultimo_resumen)
	//     SELECT
	//         de.codigo_nuevo,
	//         de.nombre,
	//         $1,
	//         de.resumen_comentarios,
	//         de.comentarios_ultimo_resumen
	//     FROM docentes_equivalencias de
	// ),
	// mapeo_codigos_docentes AS (
	//     SELECT de.codigo_antiguo, de.codigo_nuevo
	//     FROM docentes_equivalencias de
	// ),
	// calificaciones_dolly_copiadas AS (
	//     INSERT INTO
	//         calificacion_dolly (
	//             codigo_docente,
	//             acepta_critica,
	//             asistencia,
	//             buen_trato,
	//             claridad,
	//             clase_organizada,
	//             cumple_horarios,
	//             fomenta_participacion,
	//             panorama_amplio,
	//             responde_mails
	//         )
	//     SELECT
	//         m.codigo_nuevo,
	//         c.acepta_critica,
	//         c.asistencia,
	//         c.buen_trato,
	//         c.claridad,
	//         c.clase_organizada,
	//         c.cumple_horarios,
	//         c.fomenta_participacion,
	//         c.panorama_amplio,
	//         c.responde_mails
	//     FROM calificacion_dolly c
	//     JOIN mapeo_codigos_docentes m
	// 	ON c.codigo_docente = m.codigo_antiguo
	// ),
	// comentarios_copiados AS (
	//     INSERT INTO
	//         comentario (
	//             codigo_docente,
	//             codigo_cuatrimestre,
	//             contenido,
	//             es_de_dolly,
	//             fecha_creacion
	//         )
	//     SELECT
	//         m.codigo_nuevo,
	//         cm.codigo_cuatrimestre,
	//         cm.contenido,
	//         cm.es_de_dolly,
	//         cm.fecha_creacion
	//     FROM comentario cm
	//     JOIN mapeo_codigos_docentes m
	// 	ON cm.codigo_docente = m.codigo_antiguo
	// )
	// SELECT count(*)
	// FROM mapeo_codigos_docentes
	// 		`, codigoMateria).Scan(&n)
	//
	// 	if err != nil {
	// 		logger.Error(err)
	// 		return err
	// 	}
	//
	// 	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	// 	defer cancel()
	//
	// 	_, err = tx.Exec(ctx, `
	// UPDATE materia
	// SET docentes_migrados_de_equivalencia = TRUE
	// WHERE codigo = $1
	// 		`, codigoMateria)
	//
	// 	if err != nil {
	// 		logger.Error(err)
	// 		return err
	// 	}
	//
	// 	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	// 	defer cancel()
	//
	// 	err = tx.Commit(ctx)
	// 	if err != nil {
	// 		logger.Error(err)
	// 		return err
	// 	}
	//
	// 	logger.Debugf("copiados %v docentes de las equivalencias", n)
	//
	// 	return nil
	return nil
}

type PatchCatedrasOutput struct{}

func newPatchCatedras(_ *pgxpool.Conn, _ ultimaComision) (*PatchCatedrasOutput, error) {
	return nil, nil
}
