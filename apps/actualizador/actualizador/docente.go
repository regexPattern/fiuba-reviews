package actualizador

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/regexPattern/fiuba-reviews/actualizador-ofertas/config"
)

// Patch de actualización con los docentes registrados en la base de datos que
// han aún no han sido vinculados con un docente del SIU, y los docentes del
// SIU que tampoco se han registrado en la base de datos.
type PatchDocentesOutput struct {
	Registrados map[string]string
	Nuevos      map[string]string
}

// newPatchDocentes retorna un patch de actualización de los docentes de la
// materia a partir de su última oferta de comisiones.
func newPatchDocentes(logger *log.Logger, conn *pgxpool.Conn, uoferta ultimaOfertaMateria) (PatchDocentesOutput, error) {
	var p PatchDocentesOutput

	// Si la materia está siendo actualizada por primera vez, se migran los
	// docentes de sus equivalencias (con sus respectivas calificaciones y
	// comentarios) y se los asigna a la materia en cuestión.

	if yes, err := checkDocentesYaMigrados(logger, conn, uoferta); err != nil {
		return p, err
	} else if !yes {
		n, err := migrarDocentesEquivalencias(logger, conn, uoferta.Codigo)
		if err != nil {
			return p, err
		}
		logger.Log(config.DebugIndividualOps,
			fmt.Sprintf("migrado %v docentes", n))
	} else {
		logger.Log(config.DebugIndividualOps, "docentes ya migrados")
	}

	dsDb, err := getDocentesDb(logger, conn, uoferta.Codigo)
	if err != nil {
		return p, err
	}

	dsSiu := make(map[string]string)
	for _, c := range uoferta.Catedras {
		for _, d := range c.Docentes {
			dsSiu[d.Nombre] = d.Rol
		}
	}

	dsRegist := make(map[string]string, len(dsDb))
	dsNuevos := make(map[string]string, len(dsSiu))

	// Si un docente de la base de datos y otro del SIU tienen el mismo nombre
	// entonces ya fueron vinculados previamente. Solo nos interesan los
	// patches de actualizaciones de los que aún no han sido vinculados.

	for nom, cod := range dsDb {
		if _, ok := dsSiu[nom]; !ok {
			dsRegist[nom] = cod
		}
	}
	for nom, rol := range dsSiu {
		if _, ok := dsDb[nom]; !ok {
			dsNuevos[nom] = rol
		}
	}

	p = PatchDocentesOutput{
		Registrados: dsRegist,
		Nuevos:      dsNuevos,
	}

	return p, nil
}

// checkDocentesYaMigrados verifica si los docentes de una materia ya han sido
// migrados desde sus equivalencias.
func checkDocentesYaMigrados(logger *log.Logger, conn *pgxpool.Conn, uoferta ultimaOfertaMateria) (bool, error) {
	var done bool
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	err := conn.QueryRow(ctx, `
	SELECT docentes_migrados_de_equivalencia
	FROM materia
	WHERE codigo = $1
		`, uoferta.materia.Codigo).Scan(&done)
	if err != nil {
		msg := "error determinando si docentes de equivalencia ya fueron migrados"
		return false, logErrRetMsg(logger, msg, err)
	}

	return done, nil
}

// getDocentesDb obtiene los docentes de la materia de la base de datos. La
// obtención de los docentes es directa, es decir, no se obtienen los docentes
// de las equivalencias, sino que directamente los de la materia en cuestión.
func getDocentesDb(logger *log.Logger, conn *pgxpool.Conn, codMateria string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := conn.Query(ctx, `
	SELECT codigo, lower(unaccent(nombre))
	FROM docente
	WHERE codigo_materia = $1
		`, codMateria)
	if err != nil {
		msg := "error obteniendo códigos de docentes"
		return nil, logErrRetMsg(logger, msg, err)
	}

	docs := make(map[string]string)
	for rows.Next() {
		var cod, nombre string
		err := rows.Scan(&cod, &nombre)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		docs[nombre] = cod
	}

	logger.Log(config.DebugIndividualOps,
		fmt.Sprintf("encontrados %v docentes en la base de datos", len(docs)))

	return docs, nil
}

// migrarDocentesEquivalencias copia los docentes de las equivalencias de una
// materia y los asigna a la materia en cuestión, como si fueran suyos. Esto
// solo cobra importancia en el caso en el que no haya ningún docente en una
// materia del nuevo plan aún, para generar un estado inicial con los datos de
// Dolly.
func migrarDocentesEquivalencias(logger *log.Logger, conn *pgxpool.Conn, codigoMateria string) (int, error) {
	var n int

	tx, err := beginTxMigracion(logger, conn)
	if err != nil {
		return 0, err
	}
	n, err = applyModsTxMigracion(logger, tx, codigoMateria)
	if err != nil {
		return 0, err
	}
	if err := commitTxMigracion(logger, tx); err != nil {
		return 0, err
	}

	return n, nil
}

// beginTxMigracion inicia la transacción SQL para migrar los docentes de las
// equivalencias de una materia.
func beginTxMigracion(logger *log.Logger, conn *pgxpool.Conn) (pgx.Tx, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	tx, err := conn.Begin(ctx)
	if err != nil {
		msg := "error iniciando transacción SQL de migración de docentes de equivalencias"
		return nil, logErrRetMsg(logger, msg, err)
	}

	return tx, nil
}

// applyModsTxMigracion efectúa la migración los docentes de las equivalencias
// de una materia y marca la flag de migración de docentes de la materia.
func applyModsTxMigracion(logger *log.Logger, tx pgx.Tx, codMateria string) (int, error) {
	var n int
	ctx, cancel1 := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel1()

	// La query la hice con IA pero creeme que anda y hace justo lo que
	// debería hacer.

	err := tx.QueryRow(ctx, `
	WITH materias_equivalentes AS (
	    SELECT e.codigo_materia_plan_anterior AS codigo_materia_equivalente
	    FROM equivalencia e
	    WHERE e.codigo_materia_plan_vigente = $1
	),
	docentes_equivalencias AS (
	    SELECT
	        d.codigo AS codigo_antiguo,
	        gen_random_uuid() AS codigo_nuevo,
	        d.nombre,
	        d.resumen_comentarios,
	        d.comentarios_ultimo_resumen
	    FROM docente d
	    JOIN materias_equivalentes me
		ON d.codigo_materia = me.codigo_materia_equivalente
	),
	docentes_copiados AS (
	    INSERT INTO docente (codigo, nombre, codigo_materia, resumen_comentarios, comentarios_ultimo_resumen)
	    SELECT
	        de.codigo_nuevo,
	        de.nombre,
	        $1,
	        de.resumen_comentarios,
	        de.comentarios_ultimo_resumen
	    FROM docentes_equivalencias de
	),
	mapeo_codigos_docentes AS (
	    SELECT de.codigo_antiguo, de.codigo_nuevo
	    FROM docentes_equivalencias de
	),
	calificaciones_dolly_copiadas AS (
	    INSERT INTO
	        calificacion_dolly (
	            codigo_docente,
	            acepta_critica,
	            asistencia,
	            buen_trato,
	            claridad,
	            clase_organizada,
	            cumple_horarios,
	            fomenta_participacion,
	            panorama_amplio,
	            responde_mails
	        )
	    SELECT
	        m.codigo_nuevo,
	        c.acepta_critica,
	        c.asistencia,
	        c.buen_trato,
	        c.claridad,
	        c.clase_organizada,
	        c.cumple_horarios,
	        c.fomenta_participacion,
	        c.panorama_amplio,
	        c.responde_mails
	    FROM calificacion_dolly c
	    JOIN mapeo_codigos_docentes m
		ON c.codigo_docente = m.codigo_antiguo
	),
	comentarios_copiados AS (
	    INSERT INTO
	        comentario (
	            codigo_docente,
	            codigo_cuatrimestre,
	            contenido,
	            es_de_dolly,
	            fecha_creacion
	        )
	    SELECT
	        m.codigo_nuevo,
	        cm.codigo_cuatrimestre,
	        cm.contenido,
	        cm.es_de_dolly,
	        cm.fecha_creacion
	    FROM comentario cm
	    JOIN mapeo_codigos_docentes m
		ON cm.codigo_docente = m.codigo_antiguo
	)
	SELECT count(*)
	FROM mapeo_codigos_docentes
			`, codMateria).Scan(&n)
	if err != nil {
		msg := "error migrando docentes"
		return 0, logErrRetMsg(logger, msg, err)
	}

	ctx, cancel2 := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel2()

	_, err = tx.Exec(ctx, `
	UPDATE materia
	SET docentes_migrados_de_equivalencia = TRUE
	WHERE codigo = $1
			`, codMateria)
	if err != nil {
		msg := "error marcando que docentes de equivalencias fueron migrados"
		return 0, logErrRetMsg(logger, msg, err)
	}

	return n, nil
}

// commitTxMigracion commitea la transacción SQL para migrar los docentes de
// las equivalencias de una materia.
func commitTxMigracion(logger *log.Logger, tx pgx.Tx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	if err := tx.Commit(ctx); err != nil {
		msg := "error commiteando transacción SQL de migración de docentes de equivalencias"
		return logErrRetMsg(logger, msg, err)
	}

	return nil
}
