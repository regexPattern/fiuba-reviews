package main

import (
	"context"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type actualizacion struct{}

func GetActualizacionesMaterias(ofertas []oferta) ([]actualizacion, error) {
	logger := log.Default().WithPrefix("🔄")

	materiasNoActualizadas, err := getMateriasNoActualizadasEnCuatriActual()
	if err != nil {
		return nil, err
	}

	ultimasComisiones := filtrarUltimasComisiones(ofertas)

	var wg sync.WaitGroup
	sem := make(chan struct{}, int(db.Config().MaxConns))

	ch := make(chan actualizacion, len(ultimasComisiones))

	for _, uc := range ultimasComisiones {
		if _, ok := materiasNoActualizadas[uc.materia.Codigo]; ok {
			logger := logger.With("codigo", uc.materia.Codigo)
			wg.Add(1)
			go func() {
				sem <- struct{}{}

				if actualizacion, yaActualizada, err := prepararActualizacion(logger, uc); err != nil {
					logger.Warn("saltando actualización de materia")
				} else if yaActualizada {
					ch <- actualizacion
				}

				<-sem
				wg.Done()
			}()
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	actualizaciones := make([]actualizacion, len(ch))
	for a := range ch {
		actualizaciones = append(actualizaciones, a)
	}

	return actualizaciones, nil
}

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
		return nil, err
	}

	codigosMaterias := make(map[string]bool)

	for rows.Next() {
		var cod string

		err := rows.Scan(&cod)
		if err != nil {
			logger.Error("error serializando las materias",
				"error", err, "codigo", cod)
			return nil, err
		}

		codigosMaterias[cod] = true
	}

	logger.Infof("encontradas %v materias que pueden requerir actualización", len(codigosMaterias))

	return codigosMaterias, nil
}

func prepararActualizacion(logger *log.Logger, uc ultimaComision) (actualizacion, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conn, err := db.Acquire(ctx)
	if err != nil {
		logger.Error(err)
		return actualizacion{}, false, err
	}

	defer conn.Release()

	var yaActualizada bool

	err = conn.QueryRow(ctx, `
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
		Scan(&yaActualizada)

	if err != nil {
		logger.Error(err)
		return actualizacion{}, false, err
	}

	err = prepararActualizacionDocentes(logger, conn, uc)
	if err != nil {
		return actualizacion{}, false, err
	}

	err = prepararActualizacionCatedras(logger, conn, uc)
	if err != nil {
		return actualizacion{}, false, err
	}

	return actualizacion{}, yaActualizada, nil
}

func prepararActualizacionDocentes(logger *log.Logger, conn *pgxpool.Conn, uc ultimaComision) error {
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
		return err
	}

	if primeraActualizacion {
		if err := migrarDocentesEquivalencias(conn, uc.materia.Codigo); err != nil {
			return err
		}
	}

	// A PARTIR DE AQUI SI INICIA EL TEMA DE LAS COMISIONES DEL SIU

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := conn.Query(ctx, `
SELECT codigo, lower(unaccent(nombre))
FROM docente
WHERE codigo_materia = $1
		`, uc.materia.Codigo)

	if err != nil {
		logger.Error(err)
		return err
	}

	docentesMateria := make(map[string]string)

	for rows.Next() {
		var cod, nombre string

		err := rows.Scan(&cod, &nombre)
		if err != nil {
			logger.Error("error serializando los docentes",
				"error", err, "codigo", cod, "nombre", nombre)
			return err
		}

		docentesMateria[cod] = nombre
	}

	if err != nil {

	}

	// 2. traigo los docentes de la materia
	//
	// 3. por cada uno de los docentes de la oferta de ultimas comisiones,
	// busco entre los docente que me traje de la db cual es su match. Esta
	// busqueda es por nombre. Un match perfecto seria que el docente tenga el
	// mismo nombre en la DB que en el SIU. En caso de que no haya match
	// perfecto para cada docente, se debe hacer una asociacion para cada
	// docente. Se buscan los nombres parecidos usando fuzzy match.
	//
	// 4.
	// 4.1. en caso de match perfecto solo tenemos una opcion para hacer el
	// link

	return nil
}

func migrarDocentesEquivalencias(conn *pgxpool.Conn, codigoMateria string) error {
	logger := log.Default().WithPrefix("♻️").With("codigo", codigoMateria)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	tx, err := conn.Begin(ctx)
	if err != nil {
		logger.Error(err)
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var n int

	err = tx.QueryRow(ctx, `
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
		`, codigoMateria).Scan(&n)

	if err != nil {
		logger.Error(err)
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err = tx.Exec(ctx, `
UPDATE materia
SET docentes_migrados_de_equivalencia = TRUE
WHERE codigo = $1
		`, codigoMateria)

	if err != nil {
		logger.Error(err)
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err = tx.Commit(ctx)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Debugf("copiados %v docentes de las equivalencias", n)

	return nil
}

func prepararActualizacionCatedras(logger *log.Logger, conn *pgxpool.Conn, uc ultimaComision) error {
	// logger = logger.With("codigo", uc.materia.Codigo)

	return nil
}
