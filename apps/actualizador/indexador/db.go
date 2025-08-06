package indexador

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
)

type Materia struct {
	MateriaSiu
	MateriaDb
	Cuatri
}

type MateriaDb struct {
	Codigo  string `db:"codigo"`
	Nombre  string `db:"nombre"`
	Migrada bool   `db:"docentes_migrados_de_equivalencia"`
}

func (i *Indexador) sincronizarConDb(
	ctx context.Context,
	materias []OfertaMateriaSiu,
) ([]Materia, error) {
	var encontradas []Materia
	var err error
	if encontradas, err = i.asociarMaterias(ctx, materias); err != nil {
		return nil, err
	}
	if err = i.actualizarMaterias(ctx, encontradas); err != nil {
		return nil, err
	}
	return encontradas, nil
}

func (i *Indexador) asociarMaterias(
	ctx context.Context,
	materiasSiu []OfertaMateriaSiu,
) ([]Materia, error) {
	rows, _ := i.DbConn.Query(ctx, `
SELECT
    m.codigo,
    m.nombre,
	m.docentes_migrados_de_equivalencia
FROM
    materia m
    INNER JOIN plan_materia pm ON m.codigo = pm.codigo_materia
    INNER JOIN plan p ON pm.codigo_plan = p.codigo
WHERE
    p.esta_vigente = TRUE
		`)

	materiasDb, err := pgx.CollectRows(rows, pgx.RowToStructByName[MateriaDb])
	if err != nil {
		slog.Error("error obteniendo materias registradas", "error", err)
		return nil, nil
	}

	codigosDb := make(map[string]MateriaDb, len(materiasDb))
	for _, m := range materiasDb {
		codigosDb[normalize(m.Nombre)] = m
	}

	encontradas := make([]Materia, 0, len(materiasSiu))
	for _, mSiu := range materiasSiu {
		if mDb, ok := codigosDb[normalize(mSiu.Nombre)]; ok {
			encontradas = append(encontradas, Materia{
				MateriaSiu: mSiu.MateriaSiu,
				MateriaDb:  mDb,
				Cuatri:     mSiu.Cuatri,
			})
		} else {
			slog.Warn("materia no encontrada en la base de datos", "codigo_siu", mSiu.Codigo, "nombre", mSiu.Nombre)
		}
	}

	slog.Debug(
		fmt.Sprintf(
			"encontradas %v materias en la base de datos con ofetas del siu",
			len(encontradas),
		),
	)

	return encontradas, nil
}

func (i *Indexador) actualizarMaterias(
	ctx context.Context,
	materias []Materia,
) error {
	txTimeout := i.DbTxTimeout * time.Duration(len(materias))
	txctx, txcancel := context.WithTimeout(ctx, txTimeout)
	defer txcancel()

	tx, err := i.DbConn.Begin(txctx)
	if err != nil {
		slog.Error(
			"error iniciando transacción de actualización de materias",
			"error",
			err,
		)
		return err
	}

	defer func() { _ = tx.Rollback(txctx) }()

	for _, m := range materias {
		if m.MateriaSiu.Codigo != m.MateriaDb.Codigo {
			if err := i.actualizarCodigo(ctx, tx, &m); err != nil {
				return err
			}
		}
		if !m.Migrada {
			if err := i.migrarEquivalencias(ctx, tx, m); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *Indexador) actualizarCodigo(
	ctx context.Context,
	tx pgx.Tx,
	m *Materia,
) error {
	l := slog.Default().
		With("codigo_siu", m.MateriaSiu.Codigo,
			"codigo_db", m.MateriaDb.Codigo, "nombre", m.MateriaDb.Nombre)

	opctx, opcancel := context.WithTimeout(ctx, i.DbOpTimeout)
	defer opcancel()

	res, err := tx.Exec(opctx, `
UPDATE
    materia
SET
    codigo = $1
WHERE
    codigo = $2
		`, m.MateriaSiu.Codigo, m.MateriaDb.Codigo)

	if err != nil {
		l.Error("error actualizando código de materia", "error", err)
	} else if res.RowsAffected() > 0 {
		l.Debug("código de materia actualizado exitosamente")
	}

	// Actualizamos también en la copia in-memory del código de la base de datos para reflejar que
	// ahora ambos códigos están en línea con el valor obtenido del SIU. Esto por si en otra
	// parte del código llegamos a acceder el código de la base de datos.
	m.MateriaDb.Codigo = m.MateriaSiu.Codigo

	return err
}

func (i *Indexador) migrarEquivalencias(
	ctx context.Context,
	tx pgx.Tx,
	m Materia,
) error {
	l := slog.Default().
		With("codigo", m.MateriaDb.Codigo, "nombre", m.MateriaSiu.Nombre)

	opctx, opcancel := context.WithTimeout(ctx, i.DbOpTimeout)
	defer opcancel()

	res, err := tx.Exec(opctx, `
WITH materias_equivalentes AS (
    SELECT
        e.codigo_materia_plan_anterior AS codigo_materia_equivalente
    FROM
        equivalencia e
    WHERE
        e.codigo_materia_plan_vigente = $1
),
docentes_equivalencias AS (
    SELECT
        d.codigo AS codigo_antiguo,
        gen_random_uuid () AS codigo_nuevo,
    d.nombre,
    d.resumen_comentarios,
    d.comentarios_ultimo_resumen
FROM
    docente d
    JOIN materias_equivalentes me ON d.codigo_materia = me.codigo_materia_equivalente
),
docentes_copiados AS (
INSERT INTO docente (codigo, nombre, codigo_materia, resumen_comentarios, comentarios_ultimo_resumen)
    SELECT
        de.codigo_nuevo,
        de.nombre,
        $1,
        de.resumen_comentarios,
        de.comentarios_ultimo_resumen
    FROM
        docentes_equivalencias de
    ON CONFLICT (nombre,
        codigo_materia)
        DO NOTHING
    RETURNING
        codigo,
        nombre
),
mapeo_codigos_docentes AS (
    SELECT
        de.codigo_antiguo,
        de.codigo_nuevo
    FROM
        docentes_equivalencias de
    WHERE
        EXISTS (
            SELECT
                1
            FROM
                docente d
            WHERE
                d.codigo = de.codigo_nuevo
                AND d.nombre = de.nombre
                AND d.codigo_materia = $1)
),
calificaciones_dolly_copiadas AS (
INSERT INTO calificacion_dolly (codigo_docente, acepta_critica, asistencia, buen_trato, claridad, clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails)
    SELECT
        m.codigo_nuevo,
        c.acepta_critica,
        c.asistencia,
        c.buen_trato,
        c.claridad,
        c.clase_organizada,
        c.cumple_horarios,
        c.fomenta_participacion,
        c.panorama_amplio,
        c.responde_mails
    FROM
        calificacion_dolly c
        JOIN mapeo_codigos_docentes m ON c.codigo_docente = m.codigo_antiguo
),
comentarios_copiados AS (
INSERT INTO comentario (codigo_docente, codigo_cuatrimestre, contenido, es_de_dolly, fecha_creacion)
    SELECT
        m.codigo_nuevo,
        cm.codigo_cuatrimestre,
        cm.contenido,
        cm.es_de_dolly,
        cm.fecha_creacion
    FROM
        comentario cm
        JOIN mapeo_codigos_docentes m ON cm.codigo_docente = m.codigo_antiguo
),
materia_actualizada AS (
    UPDATE
        materia
    SET
        docentes_migrados_de_equivalencia = TRUE
    WHERE
        codigo = $1
)
SELECT
    count(*)
FROM
    mapeo_codigos_docentes
		`, m.MateriaDb.Codigo)

	if err != nil {
		l.Error("error migrando equivalencias de materia", "error", err)
	} else if res.RowsAffected() > 0 {
		l.Debug("equivalencias de materia migradas exitosamente")
	}

	return err
}
