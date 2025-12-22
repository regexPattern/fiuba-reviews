package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

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

func getDocentesResueltos(
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
		fmt.Sprintf(
			"sincronizadas cátedras de materia %v: %v activadas, %v creadas",
			codigoMateria,
			catedrasActivadas,
			catedrasCreadas,
		),
	)

	return nil
}

// ---

type docenteCatedraResponse struct {
	Nombre           string  `json:"nombre"`
	CodigoYaResuelto *string `json:"codigo_ya_resuelto"`
}

type catedraResponse struct {
	Docentes []docenteCatedraResponse `json:"docentes"`
	Resuelta bool                     `json:"resuelta"`
}

// ---
