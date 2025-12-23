package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/regexPattern/fiuba-reviews/actualizador/queries"
)

type resolucionMateria struct {
	CodigosYaResueltos   []string                     `json:"codigos_ya_resueltos"`
	ResolucionesDocentes map[string]resolucionDocente `json:"resoluciones"`
}

type resolucionDocente struct {
	NombreDb    string  `json:"nombre_db"`
	CodigoMatch *string `json:"codigo_match"`
}

func resolverMateria(
	conn *pgx.Conn,
	patch patchMateria,
	resolucion resolucionMateria,
) error {
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		return fmt.Errorf("error iniciando transacción de resolución de materia: %w", err)
	}
	defer func() { _ = tx.Rollback(context.TODO()) }()

	codigosResueltos := make(
		[]string,
		0,
		len(resolucion.CodigosYaResueltos)+len(resolucion.ResolucionesDocentes),
	)

	codigosResueltos = append(codigosResueltos, resolucion.CodigosYaResueltos...)

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

	if err := sincronizarCatedras(tx, patch.Codigo, patch.Catedras, codigosResueltos); err != nil {
		return err
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
		"docente_nuevo_creado",
		"nombre_siu",
		nombreSiu,
		"nombre_db",
		nombreDb,
		"codigo_materia",
		codigoMateria,
		"codigo_docente",
		codigo,
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
		"docente_existente_asociado",
		"nombre_siu",
		nombreSiu,
		"nombre_db",
		nombreDb,
		"codigo_materia",
		codigoMateria,
	)

	return codigoDocente, nil
}

func sincronizarCatedras(
	tx pgx.Tx,
	codigoMateria string,
	catedras []patchCatedra,
	codigosResueltos []string,
) error {
	_, err := tx.Exec(context.TODO(), desactivarCatedrasMateriaQuery, codigoMateria)
	if err != nil {
		return fmt.Errorf("error desactivando cátedras de materia %v: %w", codigoMateria, err)
	}

	catedrasJson, err := json.Marshal(catedras)
	if err != nil {
		return fmt.Errorf("error serializando cátedras: %w", err)
	}

	codigosJson, err := json.Marshal(codigosResueltos)
	if err != nil {
		return fmt.Errorf("error serializando códigos resueltos: %w", err)
	}

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
		"catedras_sincronizadas",
		"codigo_materia",
		codigoMateria,
		"activadas",
		catedrasActivadas,
		"creadas",
		catedrasCreadas,
	)

	return nil
}
