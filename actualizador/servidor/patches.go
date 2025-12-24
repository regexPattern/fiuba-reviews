package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/regexPattern/fiuba-reviews/actualizador/queries"
)

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

type patchCatedra catedra

func getPatchesMaterias(conn *pgx.Conn) (map[string]patchMateria, error) {
	ofertas, err := newOfertasMaterias(conn)
	if err != nil {
		return nil, fmt.Errorf(
			"error obteniendo ofertas de comisiones de materias: %w",
			err,
		)
	}

	codigosMaterias := make([]string, 0, len(ofertas))
	nombresMaterias := make([]string, 0, len(ofertas))

	for codMat, ofMat := range ofertas {
		codigosMaterias = append(codigosMaterias, codMat)
		nombresMaterias = append(nombresMaterias, ofMat.Nombre)
	}

	if err := sincronizarMaterias(conn, codigosMaterias, nombresMaterias); err != nil {
		return nil, fmt.Errorf(
			"error sincronizando materias de la base de datos con el siu: %w",
			err,
		)
	}

	patches, err := newPatchesMaterias(conn, codigosMaterias, ofertas)
	if err != nil {
		return nil, fmt.Errorf(
			"error construyendo patches de actualización de materias: %w",
			err,
		)
	}

	return patches, nil
}

func newPatchesMaterias(
	conn *pgx.Conn,
	codigosMaterias []string,
	ofertas map[string]ofertaMateriaMasReciente,
) (map[string]patchMateria, error) {
	rows, err := conn.Query(context.TODO(), queries.MateriasCandidatas, codigosMaterias)
	if err != nil {
		return nil, fmt.Errorf("error consultando materias candidatas a actualizarse: %w", err)
	}

	materiasCandidatas, err := pgx.CollectRows(rows, pgx.RowToStructByName[materia])
	if err != nil {
		return nil, fmt.Errorf("error deserializando materias candidatas a actualizarse: %v", err)
	}

	slog.Info("materias_actualizacion_pendiente", "count", len(materiasCandidatas))

	var totalDocentes, docentesNuevos, totalCatedras, catedrasNuevas int
	patches := make(map[string]patchMateria, len(materiasCandidatas))

	for _, mat := range materiasCandidatas {
		oferta, ok := ofertas[mat.Codigo]
		if !ok {
			slog.Debug("materia_sin_oferta", "codigo_materia", mat.Codigo)
			continue
		}

		if pat, err := newPatchMateria(conn, oferta); err != nil {
			return nil, fmt.Errorf(
				"error determinando si oferta de materia %v tiene actualización disponible: %w",
				mat.Codigo,
				err,
			)
		} else if pat != nil {
			totalDocentes += len(pat.Docentes)
			totalCatedras += len(pat.Catedras)
			for _, doc := range pat.Docentes {
				if len(doc.Matches) == 0 {
					docentesNuevos++
				}
			}
			catedrasNuevas += len(pat.Catedras)
			patches[pat.Codigo] = *pat
		}
	}

	slog.Info("materias_actualizacion_disponible", "count", len(patches))

	return patches, nil
}

func newPatchMateria(
	conn *pgx.Conn,
	oferta ofertaMateriaMasReciente,
) (*patchMateria, error) {
	catedrasFiltradas := make([]catedra, 0, len(oferta.Catedras))
	var catedrasDescartadas int

	for _, cat := range oferta.Catedras {
		tieneDocenteVacio := false
		for _, doc := range cat.Docentes {
			if doc.Nombre == "" {
				tieneDocenteVacio = true
				break
			}
		}
		if !tieneDocenteVacio {
			catedrasFiltradas = append(catedrasFiltradas, cat)
		} else {
			catedrasDescartadas++
		}
	}

	if catedrasDescartadas > 0 {
		slog.Warn(
			"catedras_descartadas",
			"count",
			catedrasDescartadas,
			"codigo_materia",
			oferta.Codigo,
			"motivo",
			"docentes_vacios",
		)
	}

	ofertaFiltrada := oferta
	ofertaFiltrada.Catedras = catedrasFiltradas

	patchesDocentes, err := newPatchesDocentes(conn, ofertaFiltrada)
	if err != nil {
		return nil, fmt.Errorf(
			"error generando patches de actualización de docentes de materia %v: %w",
			oferta.Codigo,
			err,
		)
	}

	patchesCatedras, err := newPatchesCatedras(conn, ofertaFiltrada)
	if err != nil {
		return nil, fmt.Errorf(
			"error generando patches de actualización de cátedras de materia %v: %w",
			oferta.Codigo,
			err,
		)
	}

	if len(patchesDocentes) == 0 && len(patchesCatedras) == 0 {
		slog.Debug("materia_sin_cambios", "codigo_materia", oferta.Codigo)
		return nil, nil
	}

	var docentesConMatches, docentesSinMatches int

	for _, pat := range patchesDocentes {
		if len(pat.Matches) > 0 {
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

	docentesDeOfertasYaResueltos := len(docentesUnicos) - len(patchesDocentes)
	catedrasNuevas := len(patchesCatedras)
	catedrasExistentes := len(oferta.Catedras) - catedrasNuevas

	slog.Debug("patch_materia_generado", "codigo_materia", oferta.Codigo,
		slog.Group("docentes",
			"sin_matches", docentesSinMatches,
			"con_matches", docentesConMatches,
			"ya_resueltos", docentesDeOfertasYaResueltos,
		),
		slog.Group("catedras",
			"nuevas", catedrasNuevas,
			"existentes", catedrasExistentes,
		),
	)

	return &patchMateria{
		materia:      ofertaFiltrada.materia,
		Docentes:     patchesDocentes,
		Catedras:     patchesCatedras,
		cuatrimestre: ofertaFiltrada.cuatrimestre,
	}, nil
}

func newPatchesDocentes(conn *pgx.Conn, oferta ofertaMateriaMasReciente) ([]patchDocente, error) {
	docentesUnicos := make(map[string]docente)
	for _, cat := range oferta.Catedras {
		for _, doc := range cat.Docentes {
			docentesUnicos[doc.Nombre] = doc
		}
	}

	nombresDocentes := slices.Collect(maps.Keys(docentesUnicos))

	rows, err := conn.Query(
		context.TODO(),
		queries.DocentesPendientes,
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

	type docentePendienteRow struct {
		NombreSiu string `db:"nombre_siu"`
		matchDocente
	}

	docentesPendientes, err := pgx.CollectRows(rows, pgx.RowToStructByName[docentePendienteRow])
	if err != nil {
		return nil, fmt.Errorf(
			"error serializando docentes pendientes materia %v: %w",
			oferta.Codigo,
			err,
		)
	}

	matchesPorDocente := make(map[string][]matchDocente)
	for _, doc := range docentesPendientes {
		if doc.Codigo != nil {
			matchesPorDocente[doc.NombreSiu] = append(
				matchesPorDocente[doc.NombreSiu],
				doc.matchDocente,
			)
		} else {
			if _, ok := matchesPorDocente[doc.NombreSiu]; !ok {
				matchesPorDocente[doc.NombreSiu] = make([]matchDocente, 0)
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

func newPatchesCatedras(conn *pgx.Conn, oferta ofertaMateriaMasReciente) ([]patchCatedra, error) {
	catedrasJson, err := json.Marshal(oferta.Catedras)
	if err != nil {
		return nil, fmt.Errorf("error serializando cátedras de materia %v: %w", oferta.Codigo, err)
	}

	rows, err := conn.Query(
		context.TODO(),
		queries.CatedrasConEstado,
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

	type catedraConEstadoRow struct {
		Codigo   int  `db:"codigo_siu"`
		Resuelta bool `db:"resuelta"`
	}

	catedrasConEstado, err := pgx.CollectRows(rows, pgx.RowToStructByName[catedraConEstadoRow])
	if err != nil {
		return nil, fmt.Errorf(
			"error serializando cátedras de materia %v: %w",
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
		patches = append(patches, patchCatedra(cat))
	}

	return patches, nil
}
