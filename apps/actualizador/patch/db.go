package patch

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type materiaDb struct {
	Codigo string `db:"codigo"`
	Nombre string `db:"nombre"`
}

// actualizarCodigosMaterias actualiza los códigos de las materias de la base de datos con los
// códigos de las materias obtenidos desde el SIU.
func (g *GeneradorPatches) actualizarCodigosMaterias(patches []Patch) error {
	// Cuando se agregaron las materias de los nuevos planes a la base de datos se les puso un
	// placeholder como código ya que los PDFs de los planes no tenían información sobre el código
	// oficial de cada materia, por lo que este se tiene que actualizar una vez se tienen la
	// información del SIU.

	ctx, cancel := context.WithTimeout(context.Background(), g.DbTimeout)
	defer cancel()

	conn, err := pgx.Connect(ctx, g.DbUrl)
	if err != nil {
		slog.Error("error estableciendo conexión con la base de datos", "error", err)
		return err
	}

	defer conn.Close(ctx)

	cm, err := obtenerCodigosMateriasDesactualizadas(ctx, conn)
	if err != nil {
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	for _, p := range patches {
		codDb, ok := cm[normalize(p.Nombre)]
		if !ok {
			continue
		}

		res, err := tx.Exec(ctx, `
UPDATE
    materia
SET
    codigo = $1
WHERE
    codigo = $2
			`, p.Codigo, codDb)

		if err != nil {
			slog.Error("error actualizando código de materia",
				"codigo", p.Codigo,
				"nombre", p.Nombre,
				"error", err)
			return err
		}

		if res.RowsAffected() > 0 {
			slog.Debug("código de materia actualizado exitosamente",
				"codigo", p.Codigo,
				"nombre", p.Nombre)
		}

	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error(
			"error confirmando transacción de actualización de códigos de materias",
			"error",
			err,
		)
		return err
	}

	return nil
}

// obtenerCodigosMateriasDesactualizadas obtiene los códigos de las materias con códigos
// desactualizados de la base de datos. Revisar actualizarCodigosMaterias para mayor información.
func obtenerCodigosMateriasDesactualizadas(
	ctx context.Context,
	conn *pgx.Conn,
) (map[string]string, error) {
	// Las materia con código desactualizado son aquellas cuyo código actual todavía tiene el
	// prefijo 'COD'. Este fue el placeholder elegido cuando se cargaron las materias de los nuevos
	// planes.

	rows, _ := conn.Query(ctx, `
SELECT
    codigo,
    nombre
FROM
    materia
WHERE
    codigo LIKE 'COD%'
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

	cm := make(map[string]string, len(desact))
	for _, m := range desact {
		cm[normalize(m.Nombre)] = m.Codigo
	}

	return cm, nil
}

// normalize normaliza una string haciendola lowercase y eliminando los acentos.
func normalize(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return strings.ToLower(strings.TrimSpace(result))
}
