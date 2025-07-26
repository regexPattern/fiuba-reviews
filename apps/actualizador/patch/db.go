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

//go:embed queries/materias_sin_asociar.sql
var queryMateriasSinAsociar string

//go:embed queries/asociar_materia.sql
var queryAsociarMateria string

//go:embed queries/materias_sin_migrar.sql
var queryMateriasSinMigrar string

//go:embed queries/migrar_materia.sql
var queryMigrarMateria string

var pool *pgxpool.Pool

type materiaDb struct {
	Codigo string `db:"codigo"`
	Nombre string `db:"nombre"`
}

func (g *GeneradorPatches) initDbPool(dbCtx context.Context) error {
	dbCtx, cancel := context.WithTimeout(dbCtx, g.DbInitTimeout)
	defer cancel()

	var err error
	pool, err = pgxpool.New(dbCtx, g.DbUrl)

	if err != nil {
		slog.Error("error inicializado pool de conexiones con la base de datos", "error", err)
		return err
	}

	slog.Info("pool de conexiones con la base de datos inicializado exitosamente")

	return nil
}

func (g *GeneradorPatches) asociarMaterias(ctx context.Context, patches []PatchMateria) error {
	ctx, cancel := context.WithTimeout(ctx, g.DbOpsTimeout)
	defer cancel()

	sinAsociar, err := getMateriasSinAsociar(ctx)
	if err != nil {
		return err
	}

	var completas int64

	eg, egCtx := errgroup.WithContext(ctx)

	for _, p := range patches {
		if codigoDb, ok := sinAsociar[normalize(p.Nombre)]; ok {
			eg.Go(func() error {
				err := asociarMateria(egCtx, pool, p, codigoDb)
				if err == nil {
					atomic.AddInt64(&completas, 1)
				}
				return err
			})
		}
	}

	if err := eg.Wait(); err != nil {
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
	rows, _ := pool.Query(ctx, queryMateriasSinAsociar)

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
	p PatchMateria,
	codigoDb string,
) error {
	l := slog.Default().With(
		"codigo_siu", p.CodigoSiu,
		"codigo_db", codigoDb,
		"nombre", p.Nombre)

	res, err := pool.Exec(ctx, queryAsociarMateria, p.CodigoSiu, codigoDb)

	if err != nil {
		l.Error("error asociando código de materia", "error", err)
	} else if res.RowsAffected() > 0 {
		l.Debug("código de materia asociado exitosamente")
	}

	return err
}

// normalize normaliza una string haciendola lowercase y eliminando los acentos.
func normalize(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return strings.ToLower(strings.TrimSpace(result))
}

func (g *GeneradorPatches) migrarMaterias(
	ctx context.Context,
	patches []PatchMateria,
) error {
	ctx, cancel := context.WithTimeout(ctx, g.DbOpsTimeout)
	defer cancel()

	sinMigrar, err := getMateriasSinMigrar(ctx)
	if err != nil {
		return err
	}

	var completas int64

	eg, egCtx := errgroup.WithContext(ctx)

	for _, p := range patches {
		if _, ok := sinMigrar[p.Nombre]; ok {
			eg.Go(func() error {
				err := migrarMateria(egCtx, p)
				if err == nil {
					atomic.AddInt64(&completas, 1)
				}
				return err
			})
		}
	}

	if err := eg.Wait(); err != nil {
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
	rows, _ := pool.Query(ctx, queryMateriasSinMigrar)

	materias, err := pgx.CollectRows(rows, pgx.RowToStructByName[materiaDb])
	if err != nil {
		slog.Error("error obteniendo materias con equivalencias sin migrar", "error", err)
		return nil, err
	}

	slog.Debug(fmt.Sprintf("encontradas %v materias con equivalencias sin migrar", len(materias)))

	return mapNombreCodigo(materias), nil
}

func migrarMateria(ctx context.Context, p PatchMateria) error {
	l := slog.Default().With("codigo", p.CodigoSiu, "nombre", p.Nombre)

	res, err := pool.Exec(ctx, queryMigrarMateria, p.CodigoSiu)

	if err != nil {
		l.Error("error migrando equivalencias de materia", "error", err)
	} else if res.RowsAffected() > 0 {
		l.Debug("equivalencias de materia migradas exitosamente")
	}

	return nil
}

type DocenteDb struct {
	Codigo string `db:"codigo"`
	Nombre string `db:"nombre"`
}

func ObtenerDocentesMateria(codigoMateria string) ([]*DocenteDb, error) {
	l := slog.Default().With("codigo", codigoMateria)

	rows, _ := pool.Query(context.Background(), `
SELECT codigo, nombre
FROM docente
WHERE codigo_materia = $1
		`, codigoMateria)

	d, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[DocenteDb])
	if err != nil {
		return nil, err
	}

	l.Debug(
		fmt.Sprintf(
			"encontrados %v docentes de materia",
			len(d),
		),
	)

	return d, nil
}
