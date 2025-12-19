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

//go:embed queries/patch/SELECT-materias-con-posible-actualizacion.sql
var materiasConPosibleActualizacionQuery string

//go:embed queries/patch/SELECT-docentes-no-resueltos-de-materia.sql
var docentesNoResueltosDeMateriaQuery string

//go:embed queries/patch/SELECT-catedras-no-resueltas-de-materia.sql
var catedrasNoResueltasDeMateriaQuery string

type patchMateria struct {
	materia
	Docentes     []patchDocente `json:"docentes"`
	Catedras     []patchCatedra `json:"catedras"`
	cuatrimestre `               json:"cuatrimestre"`
}

type patchDocente struct {
	docente
	Matches []matchDocente `json:"matches"`
}

type matchDocente struct {
	Codigo   *string  `json:"codigo"`
	NombreDb *string  `json:"nombre"`
	Score    *float64 `json:"score"`
}

type patchCatedra struct {
	catedra
}

func buildPatchesMaterias(
	conn *pgx.Conn,
	codMaterias []string,
	ofertas map[string]ultimaOfertaMateria,
) (map[string]patchMateria, error) {
	rows, err := conn.Query(
		context.TODO(),
		materiasConPosibleActualizacionQuery,
		codMaterias,
	)
	if err != nil {
		return nil, fmt.Errorf("error consultando materias candidatas a actualizarse: %w", err)
	}

	materiasCandidatas, err := pgx.CollectRows(rows, pgx.RowToStructByName[materia])
	if err != nil {
		return nil, fmt.Errorf("error deserializando materias candidatas a actualizarse: %v", err)
	}

	slog.Info(
		fmt.Sprintf(
			"encontradas %v materias con posible actualización",
			len(materiasCandidatas),
		),
	)

	var totalDocentes, docentesNuevos, totalCatedras, catedrasNuevas int
	patches := make(map[string]patchMateria, len(materiasCandidatas))

	for _, mat := range materiasCandidatas {
		oferta, ok := ofertas[mat.Codigo]
		if !ok {
			continue
		}

		if p, err := newPatchMateria(conn, oferta); err != nil {
			return nil, fmt.Errorf(
				"error determinando si oferta de materia %v (%v) tiene actualización disponible: %w",
				mat.Codigo,
				mat.Nombre,
				err,
			)
		} else if p != nil {
			totalDocentes += len(p.Docentes)
			totalCatedras += len(p.Catedras)
			for _, d := range p.Docentes {
				if len(d.Matches) == 0 {
					docentesNuevos++
				}
			}
			catedrasNuevas += len(p.Catedras)
			patches[p.Codigo] = *p
		}
	}

	slog.Info(
		fmt.Sprintf(
			"encontradas %v materias con actualización disponible",
			len(patches),
		),
	)

	return patches, nil
}

func newPatchMateria(
	conn *pgx.Conn,
	oferta ultimaOfertaMateria,
) (*patchMateria, error) {
	patchesDocentes, err := newPatchesDocentes(conn, oferta)
	if err != nil {
		return nil, fmt.Errorf(
			"error generando patch de actualización de docentes de materia %v (%v): %w",
			oferta.Codigo,
			oferta.Nombre,
			err,
		)
	}

	patchesCatedras, err := newPatchesCatedras(conn, oferta)
	if err != nil {
		return nil, fmt.Errorf(
			"error generando patch de actualización de cátedras de materia %v (%v): %w",
			oferta.Codigo,
			oferta.Nombre,
			err,
		)
	}

	if len(patchesDocentes) == 0 && len(patchesCatedras) == 0 {
		slog.Info(
			fmt.Sprintf(
				"materia %v (%v) no tiene cambios disponibles",
				oferta.Codigo,
				oferta.Nombre,
			),
		)
		return nil, nil
	}

	var docentesConMatches, docentesSinMatches int

	for _, patch := range patchesDocentes {
		if len(patch.Matches) > 0 {
			docentesConMatches++
		} else {
			docentesSinMatches++
		}
	}

	docentesUnicos := make(map[string]docente)
	for _, cat := range oferta.Catedras {
		for _, doc := range cat.Docentes {
			docentesUnicos[doc.Nombre] = doc
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
		materia:      oferta.materia,
		Docentes:     patchesDocentes,
		Catedras:     patchesCatedras,
		cuatrimestre: oferta.cuatrimestre,
	}, nil
}

func newPatchesDocentes(conn *pgx.Conn, oferta ultimaOfertaMateria) ([]patchDocente, error) {
	docentesUnicos := make(map[string]docente)
	for _, cat := range oferta.Catedras {
		for _, doc := range cat.Docentes {
			docentesUnicos[doc.Nombre] = doc
		}
	}

	nombresDocentes := slices.Collect(maps.Keys(docentesUnicos))

	rows, err := conn.Query(
		context.TODO(),
		docentesNoResueltosDeMateriaQuery,
		oferta.Codigo,
		nombresDocentes,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error consultando docentes no vinculados al siu de materia: %w",
			err,
		)
	}
	defer rows.Close()

	type docenteNoResueltosRow struct {
		NombreSiu string `db:"nombre_siu"`
		matchDocente
	}

	docentesSinMatch, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[docenteNoResueltosRow],
	)
	if err != nil {
		return nil, fmt.Errorf("error serializando docentes sin matches de materia: %w", err)
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
			docente: docentesUnicos[doc],
			Matches: matches,
		})
	}

	return patches, nil
}

func newPatchesCatedras(conn *pgx.Conn, oferta ultimaOfertaMateria) ([]patchCatedra, error) {
	catedrasJson, err := json.Marshal(oferta.Catedras)
	if err != nil {
		return nil, fmt.Errorf("error serializando cátedras de materia: %w", err)
	}

	rows, err := conn.Query(
		context.TODO(),
		catedrasNoResueltasDeMateriaQuery,
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

	catedrasNoRegistradas, err := pgx.CollectRows(rows, pgx.RowTo[int])
	if err != nil {
		return nil, fmt.Errorf(
			"error procesando cátedras no registradas de materia %v: %w",
			oferta.Codigo,
			err,
		)
	}

	patches := make([]patchCatedra, 0, len(catedrasNoRegistradas))
	codigosCatedrasNuevas := make(map[int]bool)
	for _, cod := range catedrasNoRegistradas {
		codigosCatedrasNuevas[cod] = true
	}

	for _, cat := range oferta.Catedras {
		if codigosCatedrasNuevas[cat.Codigo] {
			patches = append(patches, patchCatedra{
				catedra: cat,
			})
		}
	}

	return patches, nil
}

func aplicarPatchMateria(conn *pgx.Conn, patch patchMateria, docentesResueltos map[string]string) {
	for nombreSiu, codMatch := range docentesResueltos {
		// crear el docente en la base de datos usando sql
		// asociar el nombre siu del docente
		// TODO: por el momento poner el campo nombre del docente igual al nombre del siu
		_, _ = nombreSiu, codMatch
	}

	fmt.Println(patch.materia, docentesResueltos)
}
