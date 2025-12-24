package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/regexPattern/fiuba-reviews/actualizador/queries"
)

type resolucion struct {
	NombreSiu   string  `json:"nombre_siu"`
	NombreDb    string  `json:"nombre_db"`
	CodigoMatch *string `json:"codigo_match"`
}

func resolverMateria(
	conn *pgx.Conn,
	patch patchMateria,
	resoluciones []resolucion,
) error {
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		return fmt.Errorf("error iniciando transacción de resolución de materia: %w", err)
	}
	defer func() { _ = tx.Rollback(context.TODO()) }()

	var codigosUpdate, nombresSiuUpdate, nombresDbUpdate []string
	var nombresSiuInsert, nombresDbInsert []string

	for _, res := range resoluciones {
		if res.CodigoMatch != nil {
			codigosUpdate = append(codigosUpdate, *res.CodigoMatch)
			nombresSiuUpdate = append(nombresSiuUpdate, res.NombreSiu)
			nombresDbUpdate = append(nombresDbUpdate, res.NombreDb)
		} else {
			nombresSiuInsert = append(nombresSiuInsert, res.NombreSiu)
			nombresDbInsert = append(nombresDbInsert, res.NombreDb)
		}
	}

	fmt.Println("codigo_update", codigosUpdate)
	fmt.Println("nombre_update_(siu)_(db)", nombresSiuUpdate, nombresDbUpdate)
	fmt.Println("nombres_insert_(siu)_(db)", nombresSiuInsert, nombresDbInsert)

	if len(codigosUpdate) > 0 {
		_, err := tx.Exec(
			context.TODO(),
			queries.UpdateDocentes,
			codigosUpdate,
			nombresSiuUpdate,
			nombresDbUpdate,
		)
		if err != nil {
			return fmt.Errorf("error actualizando docentes existentes: %w", err)
		}
	}

	if len(nombresSiuInsert) > 0 {
		_, err := tx.Exec(
			context.TODO(),
			queries.InsertDocentes,
			patch.Codigo,
			nombresSiuInsert,
			nombresDbInsert,
		)
		if err != nil {
			return fmt.Errorf("error insertando docentes nuevos: %w", err)
		}
	}

	slog.Debug(
		"materia_resuelta",
		slog.Group(
			"docentes",
			"actualizados",
			len(codigosUpdate),
			"creados",
			len(nombresSiuInsert),
		),
	)

	catedrasJson, err := json.Marshal(patch.Catedras)
	_ = catedrasJson
	if err != nil {
		return fmt.Errorf("error serializando cátedras: %w", err)
	}

	// row := tx.QueryRow(context.TODO(), queries.UpsertCatedras, patch.Codigo,
	// string(catedrasJson))

	var catedrasActivadas, catedrasCreadas int
	// if err := row.Scan(&catedrasActivadas, &catedrasCreadas); err != nil {
	// 	return fmt.Errorf("error sincronizando cátedras: %w", err)
	// }

	if catedrasActivadas > 0 || catedrasCreadas > 0 {
		slog.Info("catedras_actualizadas",
			"codigo_materia", patch.Codigo,
			"activadas", catedrasActivadas,
			"creadas", catedrasCreadas)
	}

	if err := tx.Commit(context.TODO()); err != nil {
		return fmt.Errorf("error confirmando transacción: %w", err)
	}

	return nil
}

func getDocentesConEstadoPorCatedra(
	conn *pgx.Conn,
	codigoMateria string,
	catedras []patchCatedra,
) (map[int]map[string]*string, error) {
	catedrasJson, err := json.Marshal(catedras)
	if err != nil {
		return nil, fmt.Errorf("error serializando cátedras de materia: %w", err)
	}

	rows, err := conn.Query(
		context.TODO(),
		queries.DocentesConEstado,
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

	type docenteConEstadoRow struct {
		CodigoCatedra int     `db:"codigo_catedra_siu"`
		NombreDocente string  `db:"nombre_docente_siu"`
		CodigoDocente *string `db:"codigo_docente"`
	}

	docentesConEstado, err := pgx.CollectRows(rows, pgx.RowToStructByName[docenteConEstadoRow])
	if err != nil {
		return nil, fmt.Errorf(
			"error serializando docentes resueltos de cátedras de materia %v: %w",
			codigoMateria,
			err,
		)
	}

	docentesPorCatedra := make(map[int]map[string]*string)
	for _, doc := range docentesConEstado {
		if _, ok := docentesPorCatedra[doc.CodigoCatedra]; !ok {
			docentesPorCatedra[doc.CodigoCatedra] = make(map[string]*string)
		}
		docentesPorCatedra[doc.CodigoCatedra][doc.NombreDocente] = doc.CodigoDocente
	}

	return docentesPorCatedra, nil
}
