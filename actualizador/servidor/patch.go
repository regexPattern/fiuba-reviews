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

//go:embed queries/patch/SELECT-catedras-de-materia.sql
var catedrasDeMateriaQuery string

//go:embed queries/patch/INSERT-nuevo-docente.sql
var crearNuevoDocenteQuery string

//go:embed queries/patch/UPDATE-asociar-docente-existente.sql
var asociarDocenteExistenteQuery string

//go:embed queries/patch/SELECT-docentes-resueltos-de-catedras.sql
var docentesResueltosDeCatedrasQuery string

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
	Resuelta bool `json:"resuelta"`
}

// Tipos de respuesta específicos para el endpoint GET /patches/{codigoMateria}
type docenteCatedraResponse struct {
	Nombre   string `json:"nombre"`
	Resuelto bool   `json:"resuelto"`
}

type catedraResponse struct {
	Docentes []docenteCatedraResponse `json:"docentes"`
	Resuelta bool                     `json:"resuelta"`
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

	type catedraConEstadoRow struct {
		Codigo   int  `db:"codigo_siu"`
		Resuelta bool `db:"resuelta"`
	}

	rows, err := conn.Query(
		context.TODO(),
		catedrasDeMateriaQuery,
		oferta.Codigo,
		string(catedrasJson),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error consultando cátedras de materia %v: %w",
			oferta.Codigo,
			err,
		)
	}
	defer rows.Close()

	catedrasConEstado, err := pgx.CollectRows(rows, pgx.RowToStructByName[catedraConEstadoRow])
	if err != nil {
		return nil, fmt.Errorf(
			"error procesando cátedras de materia %v: %w",
			oferta.Codigo,
			err,
		)
	}

	estadoPorCodigo := make(map[int]bool)
	for _, cat := range catedrasConEstado {
		estadoPorCodigo[cat.Codigo] = cat.Resuelta
	}

	patches := make([]patchCatedra, 0, len(oferta.Catedras))
	for _, cat := range oferta.Catedras {
		patches = append(patches, patchCatedra{
			catedra:  cat,
			Resuelta: estadoPorCodigo[cat.Codigo],
		})
	}

	return patches, nil
}

func getDocentesResueltosDeCatedras(
	conn *pgx.Conn,
	codigoMateria string,
	catedras []patchCatedra,
) (map[int]map[string]bool, error) {
	catedrasJson, err := json.Marshal(catedras)
	if err != nil {
		return nil, fmt.Errorf("error serializando cátedras de materia: %w", err)
	}

	type docenteResueltoRow struct {
		CodigoCatedra int    `db:"codigo_catedra_siu"`
		NombreDocente string `db:"nombre_docente_siu"`
		Resuelto      bool   `db:"resuelto"`
	}

	rows, err := conn.Query(
		context.TODO(),
		docentesResueltosDeCatedrasQuery,
		codigoMateria,
		string(catedrasJson),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error consultando docentes resueltos de cátedras de materia %v: %w",
			codigoMateria,
			err,
		)
	}
	defer rows.Close()

	docentesResueltos, err := pgx.CollectRows(rows, pgx.RowToStructByName[docenteResueltoRow])
	if err != nil {
		return nil, fmt.Errorf(
			"error serializando docentes resueltos de cátedras de materia %v: %w",
			codigoMateria,
			err,
		)
	}

	resueltos := make(map[int]map[string]bool)
	for _, doc := range docentesResueltos {
		if _, ok := resueltos[doc.CodigoCatedra]; !ok {
			resueltos[doc.CodigoCatedra] = make(map[string]bool)
		}
		resueltos[doc.CodigoCatedra][doc.NombreDocente] = doc.Resuelto
	}

	return resueltos, nil
}

func aplicarPatchMateria(
	conn *pgx.Conn,
	patch patchMateria,
	resoluciones map[string]struct {
		NombreDb    string  `json:"nombre_db"`
		CodigoMatch *string `json:"codigo_match"`
	},
) error {
	tx, _ := conn.Begin(context.TODO())
	defer tx.Rollback(context.TODO())

	for nombreSiu, resolucion := range resoluciones {
		if resolucion.CodigoMatch == nil {
			_ = crearNuevoDocente(
				tx,
				patch.Codigo,
				nombreSiu,
				resolucion.NombreDb,
			)
		} else {
			_ = asociarDocenteExistente(
				tx,
				patch.Codigo,
				*resolucion.CodigoMatch,
				nombreSiu,
				resolucion.NombreDb,
			)
		}
	}

	_ = tx.Commit(context.TODO())
	return nil
}

func crearNuevoDocente(
	tx pgx.Tx,
	codigoMateria string,
	nombreSiu string,
	nombreDb string,
) error {
	_, _ = tx.Exec(
		context.TODO(),
		crearNuevoDocenteQuery,
		nombreDb,
		codigoMateria,
		nombreSiu,
	)

	slog.Debug(
		fmt.Sprintf(
			"creado nuevo docente %v (%v) de materia %v",
			nombreSiu,
			nombreDb,
			codigoMateria,
		),
	)

	return nil
}

func asociarDocenteExistente(
	tx pgx.Tx,
	codigoMateria string,
	codigoDocente string,
	nombreSiu string,
	nombreDb string,
) error {
	_, _ = tx.Exec(
		context.TODO(),
		asociarDocenteExistenteQuery,
		nombreDb,
		nombreSiu,
		codigoDocente,
	)

	slog.Debug(
		fmt.Sprintf(
			"resuelto docente existente %v (%v) de materia %v",
			nombreSiu,
			nombreDb,
			codigoMateria,
		),
	)

	return nil
}
