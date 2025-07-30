package patcher

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync/atomic"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

var pool *pgxpool.Pool

type MateriaBD struct {
	Codigo string `db:"codigo"`
	Nombre string `db:"nombre"`
}

func (i *Indexador) configPoolBD(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, i.DbInitTimeout)
	defer cancel()

	var err error
	pool, err = pgxpool.New(ctx, i.DbUrl)
	if err != nil {
		slog.Error("error configurando pool de conexiones con la base de datos", "error", err)
		return err
	}

	slog.Info("pool de conexiones con la base de datos configurado exitosamente")

	return nil
}

func (i *Indexador) syncMateriasSiuConBD(ctx context.Context, ofertas []OfertaMateriaSiu) error {
	if err := i.asociarMaterias(ctx, ofertas); err != nil {
		return err
	} else {
		return i.migrarMaterias(ctx, ofertas)
	}
}

func (i *Indexador) asociarMaterias(ctx context.Context, ofertas []OfertaMateriaSiu) error {
	sinAsociar, err := i.getMateriasSinAsociar(ctx)
	if err != nil {
		return err
	}

	bdCtx, bdCancel := context.WithTimeout(ctx, i.DbOpTimeout)
	defer bdCancel()

	var completas int64
	g, gCtx := errgroup.WithContext(bdCtx)

	for _, o := range ofertas {
		nombre := strings.ToLower(normalize(o.Materia.Nombre))
		if codigoDb, ok := sinAsociar[nombre]; ok {
			g.Go(func() error {
				asociada, err := asociarMateria(gCtx, pool, o.Materia, codigoDb)
				if asociada && err == nil {
					atomic.AddInt64(&completas, 1)
				}
				return err
			})
		}
	}

	completasVal := atomic.LoadInt64(&completas)

	if err := g.Wait(); err != nil {
		slog.Error("error asociando códigos de materias",
			"completas", completasVal,
			"incompletas", len(sinAsociar)-int(completasVal),
		)
		return err
	}

	if completasVal == 0 {
		slog.Info("no se han asociado códigos de materias")
	} else {
		slog.Info(fmt.Sprintf("asociado códigos de %v materias exitosamente", completasVal))
	}

	return nil
}

func (i *Indexador) getMateriasSinAsociar(ctx context.Context) (map[string]string, error) {
	bdCtx, bdCancel := context.WithTimeout(ctx, i.DbOpTimeout)
	defer bdCancel()

	rows, _ := pool.Query(bdCtx, `
SELECT
    codigo,
    nombre
FROM
    materia
WHERE
    codigo LIKE 'COD%'
		`)

	materias, err := pgx.CollectRows(rows, pgx.RowToStructByName[MateriaBD])
	if err != nil {
		slog.Error("error obteniendo materias con códigos desactualizados", "error", err)
		return nil, err
	}

	slog.Debug(fmt.Sprintf("encontradas %v materias con códigos desactualizados", len(materias)))

	return mapNombreCodigo(materias), nil
}

func mapNombreCodigo(materias []MateriaBD) map[string]string {
	codigos := make(map[string]string, len(materias))
	for _, m := range materias {
		codigos[strings.ToLower(normalize(m.Nombre))] = m.Codigo
	}
	return codigos
}

func asociarMateria(
	ctx context.Context,
	pool *pgxpool.Pool,
	materia MateriaSiu,
	codigoDb string,
) (bool, error) {
	l := slog.Default().With(
		"codigo_siu", materia.Codigo,
		"codigo_db", codigoDb,
		"nombre", materia.Nombre)

	res, err := pool.Exec(ctx, `
UPDATE
    materia
SET
    codigo = $1
WHERE
    codigo = $2
		`, materia.Codigo, codigoDb)

	if err != nil {
		l.Error("error asociando código de materia", "error", err)
		return false, err
	} else if res.RowsAffected() > 0 {
		l.Debug("código de materia asociado exitosamente")
		return true, nil
	} else {
		return false, nil
	}
}

func (i *Indexador) migrarMaterias(ctx context.Context, ofertas []OfertaMateriaSiu) error {
	sinMigrar, err := i.getMateriasSinMigrar(ctx)
	if err != nil {
		return err
	}

	bdCtx, bdCancel := context.WithTimeout(ctx, i.DbOpTimeout)
	defer bdCancel()

	var completas int64
	g, gCtx := errgroup.WithContext(bdCtx)

	for _, o := range ofertas {
		if _, ok := sinMigrar[o.Materia.Nombre]; ok {
			g.Go(func() error {
				err := migrarMateria(gCtx, o.Materia)
				if err == nil {
					atomic.AddInt64(&completas, 1)
				}
				return err
			})
		}
	}

	completasVal := atomic.LoadInt64(&completas)

	if err := g.Wait(); err != nil {
		slog.Error("error migrando docentes de materia desde equivalencias",
			"completas", completasVal,
			"incompletas", len(sinMigrar)-int(completasVal),
		)
		return err
	}

	if completasVal == 0 {
		slog.Info("no se han migrado materias desde equivalencias")
	} else {
		slog.Info(fmt.Sprintf("migradas %v materias desde equivalencias exitosamente", completasVal))
	}

	return nil
}

func (i *Indexador) getMateriasSinMigrar(ctx context.Context) (map[string]string, error) {
	bdCtx, bdCancel := context.WithTimeout(ctx, i.DbOpTimeout)
	defer bdCancel()

	rows, _ := pool.Query(bdCtx, `
SELECT
    m.codigo,
    m.nombre
FROM
    materia m
    JOIN plan_materia pm ON m.codigo = pm.codigo_materia
    JOIN plan p ON pm.codigo_plan = p.codigo
WHERE
    m.docentes_migrados_de_equivalencia = FALSE
    AND p.esta_vigente = TRUE
		`)

	materias, err := pgx.CollectRows(rows, pgx.RowToStructByName[MateriaBD])
	if err != nil {
		slog.Error("error obteniendo materias con equivalencias sin migrar", "error", err)
		return nil, err
	}

	slog.Debug(fmt.Sprintf("encontradas %v materias con equivalencias sin migrar", len(materias)))

	return mapNombreCodigo(materias), nil
}

func migrarMateria(ctx context.Context, materia MateriaSiu) error {
	l := slog.Default().With("codigo", materia.Codigo, "nombre", materia.Nombre)

	res, err := pool.Exec(ctx, `
WITH materias_equivalentes AS (
    SELECT
        e.codigo_materia_plan_anterior AS codigo_materia_equivalente
    FROM
        equivalencia e
    WHERE
        e.codigo_materia_plan_vigente = $1
),
docentes_equivalencias AS (
    SELECT
        d.codigo AS codigo_antiguo,
        gen_random_uuid () AS codigo_nuevo,
    d.nombre,
    d.resumen_comentarios,
    d.comentarios_ultimo_resumen
FROM
    docente d
    JOIN materias_equivalentes me ON d.codigo_materia = me.codigo_materia_equivalente
),
docentes_copiados AS (
INSERT INTO docente (codigo, nombre, codigo_materia, resumen_comentarios, comentarios_ultimo_resumen)
    SELECT
        de.codigo_nuevo,
        de.nombre,
        $1,
        de.resumen_comentarios,
        de.comentarios_ultimo_resumen
    FROM
        docentes_equivalencias de
    ON CONFLICT (nombre,
        codigo_materia)
        DO NOTHING
    RETURNING
        codigo,
        nombre
),
mapeo_codigos_docentes AS (
    SELECT
        de.codigo_antiguo,
        de.codigo_nuevo
    FROM
        docentes_equivalencias de
    WHERE
        EXISTS (
            SELECT
                1
            FROM
                docente d
            WHERE
                d.codigo = de.codigo_nuevo
                AND d.nombre = de.nombre
                AND d.codigo_materia = $1)
),
calificaciones_dolly_copiadas AS (
INSERT INTO calificacion_dolly (codigo_docente, acepta_critica, asistencia, buen_trato, claridad, clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails)
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
    FROM
        calificacion_dolly c
        JOIN mapeo_codigos_docentes m ON c.codigo_docente = m.codigo_antiguo
),
comentarios_copiados AS (
INSERT INTO comentario (codigo_docente, codigo_cuatrimestre, contenido, es_de_dolly, fecha_creacion)
    SELECT
        m.codigo_nuevo,
        cm.codigo_cuatrimestre,
        cm.contenido,
        cm.es_de_dolly,
        cm.fecha_creacion
    FROM
        comentario cm
        JOIN mapeo_codigos_docentes m ON cm.codigo_docente = m.codigo_antiguo
),
materia_actualizada AS (
    UPDATE
        materia
    SET
        docentes_migrados_de_equivalencia = TRUE
    WHERE
        codigo = $1
)
SELECT
    count(*)
FROM
    mapeo_codigos_docentes
		`, materia.Codigo)

	if err != nil {
		l.Error("error migrando equivalencias de materia", "error", err)
	} else if res.RowsAffected() > 0 {
		l.Debug("equivalencias de materia migradas exitosamente")
	}

	return nil
}

type InfoActualMateria struct {
	Nombre   string
	Docentes []DocenteDb
}

type DocenteDb struct {
	Codigo    string  `db:"codigo"`
	Nombre    string  `db:"nombre"`
	NombreSiu *string `db:"nombre_siu"`
}

func GetInfoMateria(codigo string) (*InfoActualMateria, error) {
	var nombre string

	err := pool.QueryRow(context.Background(), `
SELECT
		nombre
FROM
		materia
WHERE
		codigo = $1
		`, codigo).Scan(&nombre)

	if err != nil {
		return nil, err
	}

	rows, _ := pool.Query(context.Background(), `
SELECT
    codigo,
    nombre,
    nombre_siu
FROM
    docente
WHERE
    codigo_materia = $1
		`, codigo)

	docentes, err := pgx.CollectRows(rows, pgx.RowToStructByName[DocenteDb])
	if err != nil {
		return nil, err
	}

	info := &InfoActualMateria{
		Nombre:   nombre,
		Docentes: docentes,
	}

	return info, nil
}

type ContextoMateriaBD struct {
	Nombre            string
	CodigosDocentes   map[string]string
	ResumenesDocentes map[string]string
}

func getContextoMateriaBD(
	ctx context.Context,
	materia MateriaSiu,
	nombreBD string,
) (ContextoMateriaBD, error) {
	var ctxMateria ContextoMateriaBD

	logger := slog.Default().With("codigo", materia.Codigo, "nombre", materia.Nombre)

	rows, err := pool.Query(ctx, `
SELECT
    codigo,
    nombre,
    resumen_comentarios
FROM
    docente
WHERE
    codigo_materia = $1
		`, materia.Codigo)

	if err != nil {
		logger.Error("error obteniendo docentes de materia", "error", err)
		return ctxMateria, err
	}

	type docenteDb struct {
		Codigo             string  `db:"codigo"`
		Nombre             string  `db:"nombre"`
		ResumenComentarios *string `db:"resumen_comentarios"`
	}

	docentes, err := pgx.CollectRows(rows, pgx.RowToStructByName[docenteDb])
	if err != nil {
		logger.Error("error encontrando docentes para contexto de materia", "error", err)
		return ctxMateria, err
	} else {
		logger.Debug(fmt.Sprintf("encontrados %v docentes materia", len(docentes)))
	}

	ctxMateria.Nombre = nombreBD
	ctxMateria.CodigosDocentes = make(map[string]string, len(docentes))
	ctxMateria.ResumenesDocentes = make(map[string]string)

	for _, d := range docentes {
		nombre := strings.ToLower(normalize(d.Nombre))
		ctxMateria.CodigosDocentes[nombre] = d.Codigo
		if d.ResumenComentarios != nil {
			ctxMateria.ResumenesDocentes[nombre] = *d.ResumenComentarios
		}
	}

	return ctxMateria, nil
}

func (i *Indexador) getNombresMateriasBD(
	ctx context.Context,
	codigos []string,
) (map[string]string, error) {
	bdCtx, bdCancel := context.WithTimeout(ctx, i.DbOpTimeout)
	defer bdCancel()

	rows, err := pool.Query(bdCtx, `
SELECT
    codigo,
    nombre
FROM
    materia
WHERE
    codigo = ANY($1)
		`, codigos)

	if err != nil {
		slog.Error("error obteniendo nombres de materias", "error", err)
		return nil, err
	}

	materias, err := pgx.CollectRows(rows, pgx.RowToStructByName[MateriaBD])
	if err != nil {
		slog.Error("error procesando nombres de materias", "error", err)
		return nil, err
	}

	nombres := make(map[string]string, len(materias))
	for _, m := range materias {
		nombres[m.Codigo] = m.Nombre
	}

	return nombres, nil
}
