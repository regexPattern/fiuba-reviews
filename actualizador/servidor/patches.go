package main

import (
	"context"
	_ "embed"
	"encoding/json"
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

//go:embed queries/6-SELECT-catedras-no-registradas-de-materia.sql
var selectCatedrasNoRegistradasDeMateriaQuery string

type patchMateria struct {
	Materia
	Docentes []patchDocente
	Catedras []patchCatedra
}

type patchDocente struct {
	Docente
	Matches []matchDocente
}

type matchDocente struct {
	Codigo    *string
	NombreDb  *string
	Similitud *float64
}

type patchCatedra struct {
	Catedra
}

func buildPatchesMaterias(
	conn *pgx.Conn,
	codigos []string,
	ofertas map[string]UltimaOfertaMateria,
) ([]patchMateria, error) {
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

	var totalDocentes, docentesNuevos, totalCatedras, catedrasNuevas int
	patches := make([]patchMateria, 0, len(candidatas))

	for _, mat := range candidatas {
		oferta, ok := ofertas[mat.Codigo]
		if !ok {
			continue
		}

		if p, err := newPatchMateria(conn, oferta); err != nil {
			return nil, fmt.Errorf(
				"error determinando si oferta de materia %v (%v) tiene actualización de oferta disponible: %w",
				mat.Codigo,
				mat.Nombre,
				err,
			)
		} else if p != nil {
			totalDocentes += len(p.Docentes)
			totalCatedras += len(p.Catedras)
			for _, pd := range p.Docentes {
				if len(pd.Matches) == 0 {
					docentesNuevos++
				}
			}
			catedrasNuevas += len(p.Catedras)
			patches = append(patches, *p)
		}
	}

	slog.Info(
		fmt.Sprintf(
			"encontradas %v materias con actualizaciones de oferta disponible",
			len(patches),
		),
	)

	return patches, nil
}

func newPatchMateria(
	conn *pgx.Conn,
	oferta UltimaOfertaMateria,
) (*patchMateria, error) {
	patchesDocentes, err := newPatchesDocentes(conn, oferta)
	if err != nil {
		return nil, nil
	}

	patchesCatedras, err := newPatchesCatedras(conn, oferta)
	if err != nil {
		return nil, nil
	}

	if len(patchesDocentes) == 0 && len(patchesCatedras) == 0 {
		return nil, nil
	}

	docentesConMatches := 0
	docentesSinMatches := 0
	for _, patch := range patchesDocentes {
		if len(patch.Matches) > 0 {
			docentesConMatches++
		} else {
			docentesSinMatches++
		}
	}

	docentesUnicos := make(map[string]Docente)
	for _, c := range oferta.Catedras {
		for _, d := range c.Docentes {
			docentesUnicos[d.Nombre] = d
		}
	}
	docentesExistentes := len(docentesUnicos) - len(patchesDocentes)

	catedrasNuevas := len(patchesCatedras)
	catedrasExistentes := len(oferta.Catedras) - catedrasNuevas

	slog.Debug(
		fmt.Sprintf(
			"generado patch de actualización para materia %v (%v)",
			oferta.Codigo,
			oferta.Nombre,
		),
		slog.Group("docentes",
			"sin_matches", docentesSinMatches,
			"con_matches", docentesConMatches,
			"existentes", docentesExistentes,
		),
		slog.Group("catedras",
			"nuevas", catedrasNuevas,
			"existentes", catedrasExistentes,
		),
	)

	return &patchMateria{
		Materia:  oferta.Materia,
		Docentes: patchesDocentes,
		Catedras: patchesCatedras,
	}, nil
}

func newPatchesDocentes(conn *pgx.Conn, oferta UltimaOfertaMateria) ([]patchDocente, error) {
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

	docentesSinMatch, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[docenteSinMatchRow],
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error procesando docentes sin match de materia %v: %w",
			oferta.Codigo,
			err,
		)
	}

	matchesPorDocente := make(map[string][]matchDocente)
	for _, row := range docentesSinMatch {
		if row.Codigo != nil {
			matchesPorDocente[row.NombreSiu] = append(
				matchesPorDocente[row.NombreSiu],
				row.matchDocente,
			)
		} else {
			if _, ok := matchesPorDocente[row.NombreSiu]; !ok {
				matchesPorDocente[row.NombreSiu] = make([]matchDocente, 0)
			}
		}
	}

	patches := make([]patchDocente, 0, len(matchesPorDocente))
	for doc, matches := range matchesPorDocente {
		patches = append(patches, patchDocente{
			Docente: docentesUnicos[doc],
			Matches: matches,
		})
	}

	return patches, nil
}

func newPatchesCatedras(conn *pgx.Conn, oferta UltimaOfertaMateria) ([]patchCatedra, error) {
	catedrasJson, err := json.Marshal(oferta.Catedras)
	if err != nil {
		return nil, fmt.Errorf(
			"error serializando cátedras de materia %v a json: %w",
			oferta.Codigo,
			err,
		)
	}

	rows, err := conn.Query(
		context.Background(),
		selectCatedrasNoRegistradasDeMateriaQuery,
		oferta.Codigo,
		string(catedrasJson),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error consultando cátedras no registradas de materia %v: %w",
			oferta.Codigo,
			err,
		)
	}
	defer rows.Close()

	type catedraNoRegistradaRow struct {
		CodigoSiu int `db:"codigo_siu"`
	}

	catedrasNoRegistradas, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[catedraNoRegistradaRow],
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error procesando cátedras no registradas de materia %v: %w",
			oferta.Codigo,
			err,
		)
	}

	patches := make([]patchCatedra, 0, len(catedrasNoRegistradas))
	codigosCatedrasNuevas := make(map[int]bool)
	for _, row := range catedrasNoRegistradas {
		codigosCatedrasNuevas[row.CodigoSiu] = true
	}

	for _, cat := range oferta.Catedras {
		if codigosCatedrasNuevas[cat.Codigo] {
			patches = append(patches, patchCatedra{
				Catedra: cat,
			})
		}
	}

	return patches, nil
}
