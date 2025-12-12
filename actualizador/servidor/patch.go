package main

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5"
)

type patchMateria struct {
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

type docenteSinMatchRow struct {
	NombreSiu string `db:"nombre_siu"`
	matchDocente
}

type patchCatedra struct{}

func getPatchesMateria(
	conn *pgx.Conn,
	codigosMaterias []string,
	ofertasMaterias map[string]UltimaOfertaMateria,
) ([]patchMateria, error) {
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

	slog.Debug(
		fmt.Sprintf(
			"encontradas %v materias con posible actualización de oferta",
			len(materiasCandidatas),
		),
	)

	patchesMaterias := make([]patchMateria, 0, len(materiasCandidatas))

	for _, m := range materiasCandidatas {
		oferta, ok := ofertasMaterias[m.Codigo]
		if !ok {
			continue
		}

		if p, err := newPatchMateria(conn, oferta); err != nil {
			return nil, fmt.Errorf(
				"error determinando si oferta de materia %v tiene cambios: %w",
				m.Codigo,
				err,
			)
		} else if p != nil {
			patchesMaterias = append(patchesMaterias, *p)
		}
	}

	return patchesMaterias, nil
}

func newPatchMateria(
	conn *pgx.Conn,
	ofertaMateria UltimaOfertaMateria,
) (*patchMateria, error) {
	docentesUnicos := make(map[string]Docente)
	tieneTitularOAdjunto := false

	for _, c := range ofertaMateria.Catedras {
		for _, d := range c.Docentes {
			if strings.Contains(d.Rol, "titular") || strings.Contains(d.Rol, "adjunto") {
				tieneTitularOAdjunto = true
			}
			docentesUnicos[d.Nombre] = d
		}
	}

	if !tieneTitularOAdjunto {
		slog.Warn(
			fmt.Sprintf(
				"materia %v no tiene docente titular ni adjunto",
				ofertaMateria.Codigo,
			),
			"carrera",
			ofertaMateria.NombreCarrera,
		)
		for _, c := range ofertaMateria.Catedras {
			fmt.Println(c.Docentes)
		}
		fmt.Println()
	}

	nombresDocentes := slices.Collect(maps.Keys(docentesUnicos))

	rows, err := conn.Query(context.Background(), `
		WITH nombres_siu AS (
			SELECT unnest($2::text[]) AS nombre
		),
		con_match_exacto AS (
			SELECT ns.nombre
			FROM nombres_siu ns
			WHERE EXISTS (
				SELECT 1 FROM docente d
				WHERE d.codigo_materia = $1
				  AND d.nombre_siu = ns.nombre
			)
		),
		sin_match_exacto AS (
			SELECT nombre FROM nombres_siu
			EXCEPT
			SELECT nombre FROM con_match_exacto
		)
		SELECT 
			sme.nombre AS nombre_siu,
			d.codigo::text AS codigo,
			d.nombre AS nombre_db,
			similarity(d.nombre, sme.nombre) AS similitud
		FROM sin_match_exacto sme
		INNER JOIN LATERAL (
			SELECT codigo, nombre
			FROM docente
			WHERE codigo_materia = $1
			  AND nombre_siu IS NULL
			  AND similarity(nombre, sme.nombre) >= 0.5
		) d ON true
		ORDER BY sme.nombre, similitud DESC
	`, ofertaMateria.Codigo, nombresDocentes)
	if err != nil {
		return nil, fmt.Errorf(
			"error consultando docentes no vinculados al siu de materia %v: %w",
			ofertaMateria.Codigo,
			err,
		)
	}
	defer rows.Close()

	docenteRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[docenteSinMatchRow])
	if err != nil {
		return nil, fmt.Errorf(
			"error procesando docentes sin match de materia %v: %w",
			ofertaMateria.Codigo,
			err,
		)
	}

	// Agrupar matches por nombre_siu
	matchesPorDocente := make(map[string][]matchDocente)
	for _, row := range docenteRows {
		matchesPorDocente[row.NombreSiu] = append(
			matchesPorDocente[row.NombreSiu],
			row.matchDocente,
		)
	}

	// Construir patchDocentes
	patchDocentes := make([]patchDocente, 0, len(matchesPorDocente))
	for nombreSiu, matches := range matchesPorDocente {
		patchDocentes = append(patchDocentes, patchDocente{
			Docente: docentesUnicos[nombreSiu],
			Matches: matches,
		})
	}

	// Por ahora retornamos solo los docentes, las cátedras se implementarán después
	if len(patchDocentes) == 0 {
		return nil, nil
	}

	return &patchMateria{
		Docentes: patchDocentes,
		Catedras: nil,
	}, nil
}
