package main

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"maps"
	"slices"

	"github.com/jackc/pgx/v5"
)

//go:embed queries/4-SELECT-materias-con-posible-actualizacion.sql
var selectMateriasConPosibleActualizacionQuery string

//go:embed queries/5-SELECT-docentes-no-vinculados-de-materia.sql
var selectDocentesNoVinculadosDeMateriaQuery string

type patchActualizacionMateria struct {
	Materia
	Docentes []patchDocente
	Catedras []patchCatedra
}

type patchDocente struct {
	Docente
	Matches []matchDocente
}

type matchDocente struct {
	Codigo    string
	NombreDb  string
	Similitud float64
}

type patchCatedra struct{}

func buildPatchesActualizacionMaterias(
	conn *pgx.Conn,
	codigos []string,
	ofertas map[string]UltimaOfertaMateria,
) ([]patchActualizacionMateria, error) {
	rows, err := conn.Query(
		context.Background(),
		selectMateriasConPosibleActualizacionQuery,
		codigos,
	)
	if err != nil {
		return nil, fmt.Errorf("error consultando materias candidatas a actualizarse: %w", err)
	}

	candidatas, err := pgx.CollectRows(rows, pgx.RowToStructByName[Materia])
	if err != nil {
		return nil, fmt.Errorf("error procesando materias candidatas a actualizarse: %v", err)
	}

	slog.Debug(
		fmt.Sprintf(
			"encontradas %v materias con posible actualización de oferta",
			len(candidatas),
		),
	)

	patches := make([]patchActualizacionMateria, 0, len(candidatas))
	totalDocentes := 0
	docentesSinMatches := 0

	for _, mat := range candidatas {
		oferta, ok := ofertas[mat.Codigo]
		if !ok {
			continue
		}

		if p, err := newPatchActualizacionMateria(conn, oferta); err != nil {
			return nil, fmt.Errorf(
				"error determinando si oferta de materia %v (%v) tiene cambios: %w",
				mat.Codigo,
				mat.Nombre,
				err,
			)
		} else if p != nil {
			totalDocentes += len(p.Docentes)
			for _, d := range p.Docentes {
				if len(d.Matches) == 0 {
					docentesSinMatches++
				}
			}
			patches = append(patches, *p)
		}
	}

	slog.Debug(fmt.Sprintf("encontradas %v materias con actualizaciones pendientes", len(patches)))
	slog.Debug(
		fmt.Sprintf(
			"encontradas %v docentes con matches en la base de datos",
			totalDocentes-docentesSinMatches,
		),
	)
	slog.Debug(
		fmt.Sprintf("encontradas %v docentes sin matches en la base de datos", docentesSinMatches),
	)

	return patches, nil
}

func newPatchActualizacionMateria(
	conn *pgx.Conn,
	oferta UltimaOfertaMateria,
) (*patchActualizacionMateria, error) {
	docentesUnicos := make(map[string]Docente)

	for _, c := range oferta.Catedras {
		for _, d := range c.Docentes {
			docentesUnicos[d.Nombre] = d
		}
	}

	nombresDocentes := slices.Collect(maps.Keys(docentesUnicos))

	rows, err := conn.Query(
		context.Background(),
		selectDocentesNoVinculadosDeMateriaQuery,
		oferta.Codigo,
		nombresDocentes,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error consultando docentes no vinculados al siu de materia %v: %w",
			oferta.Codigo,
			err,
		)
	}
	defer rows.Close()

	type docenteSinMatchRow struct {
		NombreSiu string `db:"nombre_siu"`
		matchDocente
	}

	docentesSinMatch, err := pgx.CollectRows(rows, pgx.RowToStructByName[docenteSinMatchRow])
	if err != nil {
		return nil, fmt.Errorf(
			"error procesando docentes sin match de materia %v: %w",
			oferta.Codigo,
			err,
		)
	}

	matchesPorDocente := make(map[string][]matchDocente)
	for _, row := range docentesSinMatch {
		matchesPorDocente[row.NombreSiu] = append(
			matchesPorDocente[row.NombreSiu],
			row.matchDocente,
		)
	}

	patchDocentes := make([]patchDocente, 0, len(matchesPorDocente))
	for doc, matches := range matchesPorDocente {
		patchDocentes = append(patchDocentes, patchDocente{
			Docente: docentesUnicos[doc],
			Matches: matches,
		})
	}

	// Por ahora retornamos solo los docentes, las cátedras se implementarán después
	if len(patchDocentes) == 0 {
		return nil, nil
	}

	return &patchActualizacionMateria{
		Materia:  oferta.Materia,
		Docentes: patchDocentes,
		Catedras: nil,
	}, nil
}
