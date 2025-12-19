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

//go:embed queries/patch/UPDATE-desactivar-catedras-materia.sql
var desactivarCatedrasMateriaQuery string

//go:embed queries/patch/UPSERT-catedras-resueltas.sql
var upsertCatedrasResueltasQuery string

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

type docenteCatedraResponse struct {
	Nombre           string  `json:"nombre"`
	CodigoYaResuelto *string `json:"codigo_ya_resuelto"`
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
) (map[int]map[string]*string, error) {
	catedrasJson, err := json.Marshal(catedras)
	if err != nil {
		return nil, fmt.Errorf("error serializando cátedras de materia: %w", err)
	}

	type docenteResueltoRow struct {
		CodigoCatedra int     `db:"codigo_catedra_siu"`
		NombreDocente string  `db:"nombre_docente_siu"`
		CodigoDocente *string `db:"codigo_docente"`
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

	resueltos := make(map[int]map[string]*string)
	for _, doc := range docentesResueltos {
		if _, ok := resueltos[doc.CodigoCatedra]; !ok {
			resueltos[doc.CodigoCatedra] = make(map[string]*string)
		}
		resueltos[doc.CodigoCatedra][doc.NombreDocente] = doc.CodigoDocente
	}

	return resueltos, nil
}

func aplicarPatchMateria(
	conn *pgx.Conn,
	patch patchMateria,
	resolucion resolucionMateria,
) error {
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		return fmt.Errorf("error iniciando transacción: %w", err)
	}
	defer tx.Rollback(context.TODO())

	// Recolectar todos los códigos de docentes resueltos
	codigosResueltos := make([]string, 0, len(resolucion.CodigosYaResueltos)+len(resolucion.ResolucionesDocentes))

	// Agregar los que ya estaban resueltos
	codigosResueltos = append(codigosResueltos, resolucion.CodigosYaResueltos...)

	// Procesar resoluciones de docentes y recolectar códigos
	for nombreSiu, res := range resolucion.ResolucionesDocentes {
		var codigo string
		var err error

		if res.CodigoMatch == nil {
			codigo, err = crearNuevoDocente(tx, patch.Codigo, nombreSiu, res.NombreDb)
		} else {
			codigo, err = asociarDocenteExistente(tx, patch.Codigo, *res.CodigoMatch, nombreSiu, res.NombreDb)
		}

		if err != nil {
			return err
		}

		codigosResueltos = append(codigosResueltos, codigo)
	}

	// Sincronizar cátedras
	if err := sincronizarCatedras(tx, patch.Codigo, patch.Catedras, codigosResueltos); err != nil {
		return err
	}

	if err := tx.Commit(context.TODO()); err != nil {
		return fmt.Errorf("error confirmando transacción: %w", err)
	}

	return nil
}

func crearNuevoDocente(
	tx pgx.Tx,
	codigoMateria string,
	nombreSiu string,
	nombreDb string,
) (string, error) {
	var codigo string
	err := tx.QueryRow(
		context.TODO(),
		crearNuevoDocenteQuery,
		nombreDb,
		codigoMateria,
		nombreSiu,
	).Scan(&codigo)

	if err != nil {
		return "", fmt.Errorf("error creando nuevo docente: %w", err)
	}

	slog.Debug(
		fmt.Sprintf(
			"creado nuevo docente %v (%v) de materia %v con código %v",
			nombreSiu,
			nombreDb,
			codigoMateria,
			codigo,
		),
	)

	return codigo, nil
}

func asociarDocenteExistente(
	tx pgx.Tx,
	codigoMateria string,
	codigoDocente string,
	nombreSiu string,
	nombreDb string,
) (string, error) {
	_, err := tx.Exec(
		context.TODO(),
		asociarDocenteExistenteQuery,
		nombreDb,
		nombreSiu,
		codigoDocente,
	)

	if err != nil {
		return "", fmt.Errorf("error asociando docente existente: %w", err)
	}

	slog.Debug(
		fmt.Sprintf(
			"resuelto docente existente %v (%v) de materia %v",
			nombreSiu,
			nombreDb,
			codigoMateria,
		),
	)

	return codigoDocente, nil
}

func sincronizarCatedras(
	tx pgx.Tx,
	codigoMateria string,
	catedras []patchCatedra,
	codigosResueltos []string,
) error {
	// Desactivar todas las cátedras de la materia
	_, err := tx.Exec(context.TODO(), desactivarCatedrasMateriaQuery, codigoMateria)
	if err != nil {
		return fmt.Errorf("error desactivando cátedras de materia %v: %w", codigoMateria, err)
	}

	// Serializar cátedras a JSON
	catedrasJson, err := json.Marshal(catedras)
	if err != nil {
		return fmt.Errorf("error serializando cátedras: %w", err)
	}

	// Serializar códigos resueltos a JSON
	codigosJson, err := json.Marshal(codigosResueltos)
	if err != nil {
		return fmt.Errorf("error serializando códigos resueltos: %w", err)
	}

	// Ejecutar upsert de cátedras
	var catedrasActivadas, catedrasCreadas int
	err = tx.QueryRow(
		context.TODO(),
		upsertCatedrasResueltasQuery,
		codigoMateria,
		string(catedrasJson),
		string(codigosJson),
	).Scan(&catedrasActivadas, &catedrasCreadas)

	if err != nil {
		return fmt.Errorf("error sincronizando cátedras de materia %v: %w", codigoMateria, err)
	}

	slog.Debug(
		fmt.Sprintf(
			"sincronizadas cátedras de materia %v: %v activadas, %v creadas",
			codigoMateria,
			catedrasActivadas,
			catedrasCreadas,
		),
	)

	return nil
}
