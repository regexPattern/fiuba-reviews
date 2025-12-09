package main

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/jackc/pgx/v5"
)

func getMateriasPendientes(
	conn *pgx.Conn,
	codigosMaterias []string,
	ofertasMaterias map[string]UltimaOfertaMateria,
) ([]MateriaConActualizaciones, error) {
	rows, err := conn.Query(context.Background(), `
		SELECT DISTINCT
			mat.codigo,
			mat.nombre
		FROM materia mat
		INNER JOIN plan_materia pm ON pm.codigo_materia = mat.codigo
		INNER JOIN plan ON plan.codigo = pm.codigo_plan
		WHERE plan.esta_vigente
		  AND mat.codigo = ANY($1::text[])
		  AND mat.cuatrimestre_ultima_actualizacion IS DISTINCT FROM (SELECT max(codigo) FROM cuatrimestre);
	`, codigosMaterias)
	if err != nil {
		return nil, fmt.Errorf("error consultando materias candidatas a actualizarse: %w", err)
	}

	materiasCandidatas, err := pgx.CollectRows(rows, pgx.RowToStructByName[Materia])
	if err != nil {
		return nil, fmt.Errorf("error procesando materias candidatas a actualizarse: %v", err)
	}

	materiasConActualizaciones := make([]MateriaConActualizaciones, 0, len(materiasCandidatas))

	for _, m := range materiasCandidatas {
		oferta, ok := ofertasMaterias[m.Codigo]
		if !ok {
			continue
		}

		if docentesPendientes, err := ofertaTieneCambios(conn, oferta); err != nil {
			return nil, fmt.Errorf(
				"error determinando si oferta de materia %v tiene cambios: %w",
				m.Codigo,
				err,
			)
		} else if len(docentesPendientes) > 0 {
			materiaConActualizaciones := MateriaConActualizaciones{
				Codigo:             m.Codigo,
				Nombre:             m.Nombre,
				DocentesPendientes: docentesPendientes,
				CatedrasNuevas:     []CatedraNueva{}, // TODO: Implementar lógica para catedras nuevas
			}
			materiasConActualizaciones = append(materiasConActualizaciones, materiaConActualizaciones)
		}
	}

	return materiasConActualizaciones, nil
}

type resultadoMatchDocente struct {
	NombreSiu     string   `db:"nombre_siu"`
	CodigoDocente *string  `db:"codigo_docente"`
	NombreDb      *string  `db:"nombre_db"`
	Similitud     *float64 `db:"similitud"`
}

func ofertaTieneCambios(
	conn *pgx.Conn,
	ofertaMateria UltimaOfertaMateria,
) ([]DocentePendiente, error) {
	nombresUnicosDocentes := make(map[string]bool)
	for _, c := range ofertaMateria.Catedras {
		for _, d := range c.Docentes {
			nombresUnicosDocentes[d.Nombre] = true
		}
	}

	nombresDocentes := slices.Collect(maps.Keys(nombresUnicosDocentes))

	rows, err := conn.Query(context.Background(), `
		SELECT 
			ns.nombre_siu,
			d.codigo AS codigo_docente,
			d.nombre AS nombre_db,
			similarity(lower(unaccent(ns.nombre_siu)), lower(unaccent(d.nombre))) AS similitud
		FROM unnest($2::text[]) AS ns(nombre_siu)
		LEFT JOIN LATERAL (
			SELECT codigo, nombre
			FROM docente
			WHERE codigo_materia = $1
			  AND nombre_siu IS NULL
			  AND similarity(lower(unaccent(ns.nombre_siu)), lower(unaccent(docente.nombre))) >= 0.5
			ORDER BY similarity(lower(unaccent(ns.nombre_siu)), lower(unaccent(docente.nombre))) DESC
		) d ON true
		WHERE NOT EXISTS (
			SELECT 1 
			FROM docente 
			WHERE codigo_materia = $1 
			  AND docente.nombre_siu = ns.nombre_siu
		);
	`, ofertaMateria.Codigo, nombresDocentes)
	if err != nil {
		return nil, fmt.Errorf(
			"error consultando docentes sin match perfecto de materia %v: %w",
			ofertaMateria.Codigo,
			err,
		)
	}
	defer rows.Close()

	resultados, err := pgx.CollectRows(rows, pgx.RowToStructByPos[resultadoMatchDocente])
	if err != nil {
		return nil, fmt.Errorf("error procesando resultados de docentes: %w", err)
	}

	// Agrupar resultados por nombre_siu
	docentesPendientesMap := make(map[string]*DocentePendiente)

	for _, resultado := range resultados {
		// Si no existe el docente pendiente, crearlo
		if _, exists := docentesPendientesMap[resultado.NombreSiu]; !exists {
			docentesPendientesMap[resultado.NombreSiu] = &DocentePendiente{
				NombreSiu:       resultado.NombreSiu,
				Rol:             "", // Se podría determinar si es necesario
				PosiblesMatches: []DocenteMatch{},
			}
		}

		// Si hay un match potencial, agregarlo
		if resultado.CodigoDocente != nil && resultado.NombreDb != nil &&
			resultado.Similitud != nil {
			match := DocenteMatch{
				Codigo:    *resultado.CodigoDocente,
				NombreDb:  *resultado.NombreDb,
				Similitud: *resultado.Similitud,
			}
			docentesPendientesMap[resultado.NombreSiu].PosiblesMatches = append(
				docentesPendientesMap[resultado.NombreSiu].PosiblesMatches,
				match,
			)
		}
	}

	// Convertir mapa a slice
	docentesPendientes := make([]DocentePendiente, 0, len(docentesPendientesMap))
	for _, docente := range docentesPendientesMap {
		docentesPendientes = append(docentesPendientes, *docente)
	}

	return docentesPendientes, nil
}
