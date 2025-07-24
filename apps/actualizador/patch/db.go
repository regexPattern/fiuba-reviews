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

func (g *GeneradorPatches) actualizarCodigosMaterias(om []ofertaMateria) error {
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

	for _, o := range om {
		codDb, ok := cm[normalize(o.materia.Nombre)]
		if !ok {
			continue
		}

		res, err := tx.Exec(ctx, `
			UPDATE materia 
			SET codigo = $1 
			WHERE codigo = $2
			`, o.materia.Codigo, codDb)

		if err != nil {
			slog.Error("error actualizando código de materia",
				"codigo", o.materia.Codigo,
				"nombre", o.materia.Nombre,
				"error", err)
			return err
		}

		if res.RowsAffected() > 0 {
			slog.Debug("actualizado código de materia exitosamente",
				"codigo", o.materia.Codigo,
				"nombre", o.materia.Nombre)
		}

	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("error confirmando transacción de actualización de códigos de materias", "error", err)
		return err
	}

	return nil
}

func obtenerCodigosMateriasDesactualizadas(ctx context.Context, conn *pgx.Conn) (map[string]string, error) {
	// Cuando se agregaron las materias de los nuevos planes a la base de datos
	// se les puso un placeholder como código ya que los PDFs de los planes no
	// tenían información sobre el código oficial de cada materia, por lo que
	// este se tiene que actualizar una vez se tienen la información del SIU.
	// Este placeholder es un valor número antecedido por el prefijo 'COD'. Es
	// decir, si una materia tiene un código que inicia con este prefijo, su
	// código no ha sido cambiado a su código real según el SIU.

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
		slog.Error("error obteniendo materias con códigos desactualizados de la base de datos", "error", err)
		return nil, err
	}

	slog.Debug(fmt.Sprintf("encontradas %v materias con códigos desactualizados en la base de datos", len(desact)))

	cm := make(map[string]string, len(desact))
	for _, m := range desact {
		cm[normalize(m.Nombre)] = m.Codigo
	}

	return cm, nil
}

func normalize(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return strings.ToLower(strings.TrimSpace(result))
}
