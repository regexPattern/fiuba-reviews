package main

import (
	"context"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type patch struct {
	codigoMateria string
	nombreMateria string
	docentes      *patchDocentes
	catedras      *patchCatedras
}

type patchDocentes struct {
	db  map[string]string
	siu map[string]string
}

type patchCatedras struct{}

func getPatchesMaterias(ofertas []oferta) ([]patch, error) {
	logger := log.Default().WithPrefix("⬆️")

	logger.Info("actualizando ofertas de comisiones")

	materiasPendientes, err := getMateriasPatchPendiente()
	if err != nil {
		return nil, err
	}

	log.Default().WithPrefix("🧹").Debug("filtrando solo las ofertas de comisiones más recientes")

	ultimasComisiones := filtrarUltimasComisiones(ofertas)

	var wg sync.WaitGroup
	sem := make(chan struct{}, int(db.Config().MaxConns))

	ch := make(chan patch, len(ultimasComisiones))

	for _, uc := range ultimasComisiones {
		if _, ok := materiasPendientes[uc.materia.Codigo]; ok {
			logger := logger.With("codigoMateria", uc.materia.Codigo)

			wg.Add(1)
			go func() {
				defer wg.Done()

				sem <- struct{}{}
				defer func() { <-sem }()

				if genPatchMateria(logger, ch, uc) != nil {
					logger.Warn("saltando actualización de materia")
				}
			}()
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	patches := make([]patch, len(ch))
	for a := range ch {
		patches = append(patches, a)
	}

	logger.Info("terminada generación de patches de actualización de materias")

	return patches, nil
}

func getMateriasPatchPendiente() (map[string]bool, error) {
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
			logger.Error(err)
			return nil, err
		}

		codigosMaterias[cod] = true
	}

	logger.Infof("encontradas %v materias con ofertas posiblemente desactualizadas", len(codigosMaterias))

	return codigosMaterias, nil
}

func genPatchMateria(logger *log.Logger, ch chan patch, uc ultimaComision) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conn, err := db.Acquire(ctx)
	if err != nil {
		logger.Error(err)
		return err
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
		return err
	}

	if yaActualizada {
		return nil
	}

	patchesDocentes, err := genPatchDocentes(conn, uc)
	if err != nil {
		return err
	}

	_, err = genPatchCatedras(conn, uc)
	if err != nil {
		return err
	}

	ch <- patch{
		codigoMateria: uc.materia.Codigo,
		nombreMateria: uc.materia.Nombre,
		docentes:      patchesDocentes,
		catedras:      nil,
	}

	return nil
}

func genPatchDocentes(conn *pgxpool.Conn, uc ultimaComision) (*patchDocentes, error) {
	logger := log.Default().WithPrefix("🎓").With("codigoMateria", uc.materia.Codigo)

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

	logger.Debugf("encontrados %v docentes en la base de datos", len(docentesDb))

	docentesSiu := make(map[string]string)
	for _, c := range uc.materia.Catedras {
		for _, d := range c.Docentes {
			docentesSiu[d.Nombre] = d.Rol
		}
	}

	docentesDbPendientes := make(map[string]string)
	for nombre, cod := range docentesDb {
		if _, ok := docentesSiu[nombre]; ok {
			logger.Debug("docente encontrado en la base de datos", "nombre", nombre)
			continue
		}

		docentesDbPendientes[nombre] = cod
	}

	patch := &patchDocentes{
		db:  docentesDbPendientes,
		siu: docentesSiu,
	}

	return patch, nil
}

func migrarDocentesEquivalencias(logger *log.Logger, conn *pgxpool.Conn, codigoMateria string) error {
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

func genPatchCatedras(_ *pgxpool.Conn, uc ultimaComision) (*patchCatedras, error) {
	logger := log.Default().WithPrefix("📚").With("codigoMateria", uc.materia.Codigo)

	logger.Debugf("encontradas %v cátedras en materia", 0)

	return nil, nil
}
