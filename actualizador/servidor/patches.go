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
	Carrera      string `json:"carrera"`
	cuatrimestre `               json:"cuatrimestre"`
	Docentes     []patchDocente `json:"docentes"`
	Catedras     []patchCatedra `json:"catedras"`
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
	YaExistente bool `json:"ya_existente"`
}

// getPatchesMaterias descarga las ofertas de comisiones del SIU disponibles, sincroniza las
// materias en la base de datos con los datos del SIU y retorna un hashmap donde la clave es el
// código de una materia y el valor es el patch de actualización de la misma. Solo se incluyen las
// materias que tienen actualización disponible.
func getPatchesMaterias(conn *pgx.Conn) (map[string]*patchMateria, error) {
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

	// Se tienen que sincronizar las materias antes de generar los patches de actualización para
	// armar los patches ya con los códigos oficiales.

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

// newPatchesMaterias retorna un hashmap donde la clave es el código de una materia y el valor es
// el patch de actualización de la misma. Solo se incluyen las materias que tienen actualización
// disponible.
func newPatchesMaterias(
	conn *pgx.Conn,
	codigosMaterias []string,
	ofertas map[string]ofertaMateriaMasReciente,
) (map[string]*patchMateria, error) {
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
	patches := make(map[string]*patchMateria, len(materiasCandidatas))

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
		} else if pat == nil {
			if err := marcarMateriaSinCambios(conn, oferta); err != nil {
				return nil, fmt.Errorf("error marcando materia sin cambios: %w", err)
			}
		} else {
			patches[pat.Codigo] = pat

			// Estadísticas
			totalDocentes += len(pat.Docentes)
			totalCatedras += len(pat.Catedras)
			for _, doc := range pat.Docentes {
				if len(doc.Matches) == 0 {
					docentesNuevos++
				}
			}
			catedrasNuevas += len(pat.Catedras)
		}
	}

	slog.Info(
		"materias_actualizacion_disponible",
		"con_cambios",
		len(patches),
		"sin_cambios",
		len(materiasCandidatas)-len(patches),
	)

	return patches, nil
}

// newPatchMateria retorna un puntero al patch de actualización de una materia o nil en caso de que
// no haya cambios nuevos que hacer. Una materia tiene cambios disponibles si hay docentes del SIU
// que no están registrados en la base de datos o si hay cátedras nuevas. TODO
func newPatchMateria(
	conn *pgx.Conn,
	oferta ofertaMateriaMasReciente,
) (*patchMateria, error) {
	catedrasFiltradas := make([]catedra, 0, len(oferta.Catedras))
	var catedrasDescartadas int

	// Se filtran las cátedras que tienen docentes con nombres vacios. Esto es producto de
	// errores en el scraper.

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

	oferta.Catedras = catedrasFiltradas

	patchesDocentes, err := newPatchesDocentes(conn, oferta)
	if err != nil {
		return nil, fmt.Errorf(
			"error generando patches de actualización de docentes de materia %v: %w",
			oferta.Codigo,
			err,
		)
	}

	patchesCatedras, err := newPatchesCatedras(conn, oferta)
	if err != nil {
		return nil, fmt.Errorf(
			"error generando patches de actualización de cátedras de materia %v: %w",
			oferta.Codigo,
			err,
		)
	}

	if len(patchesDocentes) == 0 && len(patchesCatedras) == 0 {
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

	var catedrasNuevas, catedrasExistentes int
	for _, pat := range patchesCatedras {
		if pat.YaExistente {
			catedrasExistentes++
		} else {
			catedrasNuevas++
		}
	}

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
		materia:      oferta.materia,
		Carrera:      oferta.NombreCarrera,
		cuatrimestre: oferta.cuatrimestre,
		Docentes:     patchesDocentes,
		Catedras:     patchesCatedras,
	}, nil
}

// newPatchesDocentes retorna un arreglo de patches de actualización para los docentes de la
// materia. En caso de que este arreglo esté vacio, significa que no hay docentes nuevos del SIU
// que deban ser registrados en la base de datos.
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

// newPatchesCatedras retorna un arreglo de patches de actualización para lás cátedras de la
// materia. En caso de que este arreglo esté vacio, significa que no hay cátedras nuevas del SIU
// que deban ser registradas en la base de datos.
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

	var codigo int
	var yaExistente bool

	catedrasConEstado := make(map[int]bool, len(oferta.Catedras))

	_, err = pgx.ForEachRow(rows, []any{&codigo, &yaExistente}, func() error {
		catedrasConEstado[codigo] = yaExistente
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf(
			"error serializando estado de cátedras de materia %v: %w",
			oferta.Codigo,
			err,
		)
	}

	todasExisten := true
	for _, yaExistente := range catedrasConEstado {
		if !yaExistente {
			todasExisten = false
			break
		}
	}

	if todasExisten {
		slog.Debug(
			"sin_catedras_nuevas",
			"count",
			len(catedrasConEstado),
			"codigo_materia",
			oferta.Codigo,
		)
		return nil, nil
	}

	patches := make([]patchCatedra, 0, len(oferta.Catedras))
	for _, cat := range oferta.Catedras {
		patches = append(patches, patchCatedra{
			catedra:     cat,
			YaExistente: catedrasConEstado[cat.Codigo],
		})
	}

	return patches, nil
}

// marcarMateriaSinCambios toma la oferta de una materia sin cambios (con patch de actualización
// nil) y actualiza el cuatrimestre de última actualización en la base de datos para indicar que
// aunque no haya cambios, esta información si corresponde al cuatrimestre de actualización.
//
// Por ejemplo, si una materia fue actualizada por última vez en 1C2025, y existe una oferta más
// reciente de 2C2025, pero sin cambios, igualmente se considera que la materia fue actualizada por
// última vez durante 2C2025, por lo tanto, se tiene que actualizar este valor.
func marcarMateriaSinCambios(conn *pgx.Conn, oferta ofertaMateriaMasReciente) error {
	_, err := conn.Exec(
		context.TODO(),
		queries.MarcarMateriaSinCambios,
		oferta.Codigo,
		oferta.Numero,
		oferta.Anio,
	)
	if err != nil {
		return fmt.Errorf("error actualizando cuatrimestre de última actualización: %w", err)
	}

	slog.Debug(
		"materia_sin_cambios",
		"codigo_materia",
		oferta.Codigo,
		"cuatrimestre",
		oferta.cuatrimestre,
	)

	return nil
}
