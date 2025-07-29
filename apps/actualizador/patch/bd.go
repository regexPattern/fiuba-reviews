package patch

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"strings"
	"sync/atomic"
	"unicode"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var pool *pgxpool.Pool

type materiaDb struct {
	Codigo string `db:"codigo"`
	Nombre string `db:"nombre"`
}

func (i *Indexador) configPoolBD(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, i.DbInitTimeout)
	defer cancel()

	var err error
	pool, err = pgxpool.New(ctx, i.DbUrl)
	if err != nil {
		slog.Error("error inicializado pool de conexiones con la base de datos", "error", err)
		return err
	}

	slog.Info("pool de conexiones con la base de datos inicializado exitosamente")

	return nil
}

func (i *Indexador) asociarMaterias(ctx context.Context, patches []Patch) error {
	ctx, cancel := context.WithTimeout(ctx, i.DbOpsTimeout)
	defer cancel()

	sinAsociar, err := getMateriasSinAsociar(ctx)
	if err != nil {
		return err
	}

	var completas int64
	g, gctx := errgroup.WithContext(ctx)

	for _, p := range patches {
		if codigoDb, ok := sinAsociar[normalize(p.Nombre)]; ok {
			g.Go(func() error {
				err := asociarMateria(gctx, pool, p, codigoDb)
				if err == nil {
					atomic.AddInt64(&completas, 1)
				}
				return err
			})
		}
	}

	if err := g.Wait(); err != nil {
		completasVal := atomic.LoadInt64(&completas)
		slog.Error("error migrando docentes de materia desde equivalencias",
			"completas", completasVal,
			"incompletas", len(sinAsociar)-int(completasVal),
		)
		return err
	}

	return nil
}

func getMateriasSinAsociar(
	ctx context.Context,
) (map[string]string, error) {
	rows, _ := pool.Query(ctx, `
SELECT
    codigo,
    nombre
FROM
    materia
WHERE
    codigo LIKE 'COD%'
		`)

	materias, err := pgx.CollectRows(rows, pgx.RowToStructByName[materiaDb])
	if err != nil {
		slog.Error("error obteniendo materias con códigos desactualizados", "error", err)
		return nil, err
	}

	slog.Debug(fmt.Sprintf("encontradas %v materias con códigos desactualizados", len(materias)))

	return mapNombreCodigo(materias), nil
}

func mapNombreCodigo(materias []materiaDb) map[string]string {
	codigos := make(map[string]string, len(materias))
	for _, m := range materias {
		codigos[normalize(m.Nombre)] = m.Codigo
	}
	return codigos
}

func asociarMateria(
	ctx context.Context,
	pool *pgxpool.Pool,
	p Patch,
	codigoDb string,
) error {
	l := slog.Default().With(
		"codigo_siu", p.CodigoSiu,
		"codigo_db", codigoDb,
		"nombre", p.Nombre)

	res, err := pool.Exec(ctx, `
UPDATE
    materia
SET
    codigo = $1
WHERE
    codigo = $2
		`, p.CodigoSiu, codigoDb)

	if err != nil {
		l.Error("error asociando código de materia", "error", err)
	} else if res.RowsAffected() > 0 {
		l.Debug("código de materia asociado exitosamente")
	}

	return err
}

func normalize(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return strings.ToLower(strings.TrimSpace(result))
}

func (i *Indexador) migrarMaterias(ctx context.Context, patches []Patch) error {
	ctx, cancel := context.WithTimeout(ctx, i.DbOpsTimeout)
	defer cancel()

	sinMigrar, err := getMateriasSinMigrar(ctx)
	if err != nil {
		return err
	}

	var completas int64

	g, gCtx := errgroup.WithContext(ctx)

	for _, p := range patches {
		if _, ok := sinMigrar[p.Nombre]; ok {
			g.Go(func() error {
				err := migrarMateria(gCtx, p)
				if err == nil {
					atomic.AddInt64(&completas, 1)
				}
				return err
			})
		}
	}

	if err := g.Wait(); err != nil {
		completasVal := atomic.LoadInt64(&completas)
		slog.Error("error migrando docentes de materia desde equivalencias",
			"completas", completasVal,
			"incompletas", len(sinMigrar)-int(completasVal),
		)
		return err
	}

	return nil
}

func getMateriasSinMigrar(ctx context.Context) (map[string]string, error) {
	rows, _ := pool.Query(ctx, `
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

	materias, err := pgx.CollectRows(rows, pgx.RowToStructByName[materiaDb])
	if err != nil {
		slog.Error("error obteniendo materias con equivalencias sin migrar", "error", err)
		return nil, err
	}

	slog.Debug(fmt.Sprintf("encontradas %v materias con equivalencias sin migrar", len(materias)))

	return mapNombreCodigo(materias), nil
}

func migrarMateria(ctx context.Context, p Patch) error {
	l := slog.Default().With("codigo", p.CodigoSiu, "nombre", p.Nombre)

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
		`, p.CodigoSiu)

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
