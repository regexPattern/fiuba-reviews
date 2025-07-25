package patch

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
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

// actualizarCodigosMaterias actualiza los códigos de las materias de la base de datos con los
// códigos de las materias obtenidos desde el SIU.
func (g *GeneradorPatches) actualizarCodigosMaterias(ctx context.Context, patches []Patch) error {
	// Cuando se agregaron las materias de los nuevos planes a la base de datos se les puso un
	// placeholder como código ya que los PDFs de los planes no tenían información sobre el código
	// oficial de cada materia, por lo que este se tiene que actualizar una vez se tienen la
	// información del SIU.
	ctx, cancel := context.WithTimeout(ctx, g.DbTimeout)
	defer cancel()

	var err error
	pool, err = pgxpool.New(ctx, g.DbUrl)
	if err != nil {
		slog.Error("error estableciendo conexión con la base de datos", "error", err)
		return err
	}

	cods, err := obtenerCodigosMateriasDesactualizadas(ctx)
	if err != nil {
		return err
	}

	eg, egCtx := errgroup.WithContext(ctx)
	sem := make(chan struct{}, pool.Config().MaxConns)

	for _, p := range patches {
		cod, ok := cods[normalize(p.Nombre)]
		if !ok {
			continue
		}

		eg.Go(func() error {
			sem <- struct{}{}
			defer func() { <-sem }()
			return actualizarCodigoMateria(egCtx, pool, p.Codigo, p.Nombre, cod)
		})
	}

	if err := eg.Wait(); err != nil {
		slog.Error("error actualizando códigos de materias", "error", err)
		return err
	}

	return nil
}

// obtenerCodigosMateriasDesactualizadas obtiene los códigos de las materias con códigos
// desactualizados de la base de datos. Revisar actualizarCodigosMaterias para mayor información.
func obtenerCodigosMateriasDesactualizadas(
	ctx context.Context,
) (map[string]string, error) {
	// Las materia con código desactualizado son aquellas cuyo código actual todavía tiene el
	// prefijo 'COD'. Este fue el placeholder elegido cuando se cargaron las materias de los nuevos
	// planes.
	rows, _ := pool.Query(ctx, `
SELECT codigo, nombre
FROM materia
WHERE codigo LIKE 'COD%'
		`)

	desact, err := pgx.CollectRows(rows, pgx.RowToStructByName[materiaDb])
	if err != nil {
		slog.Error(
			"error obteniendo materias con códigos desactualizados de la base de datos",
			"error",
			err,
		)
		return nil, err
	}

	slog.Debug(
		fmt.Sprintf(
			"encontradas %v materias con códigos desactualizados en la base de datos",
			len(desact),
		),
	)

	cods := make(map[string]string, len(desact))
	for _, m := range desact {
		cods[normalize(m.Nombre)] = m.Codigo
	}

	return cods, nil
}

// actualizarCodigoMateria actualiza el código de una materia específica en la base de datos, con
// el código oficial del SIU.
func actualizarCodigoMateria(
	ctx context.Context,
	pool *pgxpool.Pool,
	codigoSiu, nombre, codigoDb string,
) error {
	logger := slog.Default().With(
		"codigo_siu", codigoSiu,
		"codigo_db", codigoDb,
		"nombre", nombre)

	res, err := pool.Exec(ctx, `
UPDATE materia
SET codigo = $1
WHERE codigo = $2
			`, codigoSiu, codigoDb)

	if err != nil {
		logger.Error("error actualizando código de materia", "error", err)
		return err
	}

	if res.RowsAffected() > 0 {
		logger.Debug("código de materia actualizado exitosamente")
	}

	return nil
}

// normalize normaliza una string haciendola lowercase y eliminando los acentos.
func normalize(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return strings.ToLower(strings.TrimSpace(result))
}
