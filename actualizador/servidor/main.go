package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
)

type OfertaComisionesCarrera struct {
	CodigoCarrera string                    `db:"codigo_carrera"`
	Cuatrimestre  Cuatrimestre              `db:"cuatrimestre"`
	Materias      []OfertaComisionesMateria `db:"contenido"`
}

type Cuatrimestre struct {
	Numero int `json:"numero"`
	Anio   int `json:"anio"`
}

type OfertaComisionesMateria struct {
	Codigo     string     `json:"codigo"`
	Nombre     string     `json:"nombre"`
	Comisiones []Comision `json:"catedras"`
}

type Comision struct {
	Codigo   int       `json:"codigo"`
	Docentes []Docente `json:"docentes"`
}

type Docente struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

func main() {
	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)
	slog.SetDefault(slog.New(logger))

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		os.Exit(1)
	}

	slog.Info("conexión establecida con la base de datos")

	rows, _ := conn.Query(context.Background(), `
SELECT
    oc.codigo_carrera,
    json_build_object('numero', cuat.numero, 'anio', cuat.anio) AS cuatrimestre,
    oc.contenido
FROM
    oferta_comisiones oc
    INNER JOIN cuatrimestre cuat ON cuat.codigo = oc.codigo_cuatrimestre
ORDER BY
    cuat.codigo DESC;
		`)

	ofertasCarreras, err := pgx.CollectRows(rows, pgx.RowToStructByName[OfertaComisionesCarrera])
	if err != nil {
		fmt.Println(err)
	}

	slog.Debug(fmt.Sprintf("encontradas %v ofertas de comisiones de carreras", len(ofertasCarreras)))

	type UltimaOfertaComisionesMateria struct {
		OfertaComisionesMateria
		Cuatrimestre
	}

	ofertasMaterias := make(map[string]UltimaOfertaComisionesMateria)
	materiasPorCuatrimestre := make(map[Cuatrimestre]int)

	for _, oc := range ofertasCarreras {
		for _, om := range oc.Materias {
			if _, ok := ofertasMaterias[om.Nombre]; !ok {
				ofertasMaterias[om.Nombre] = UltimaOfertaComisionesMateria{
					OfertaComisionesMateria: om,
					Cuatrimestre:            oc.Cuatrimestre,
				}

				materiasPorCuatrimestre[oc.Cuatrimestre]++
			}
		}
	}

	for c, n := range materiasPorCuatrimestre {
		slog.Debug(fmt.Sprintf("encontradas %v ofertas de comisiones de materias de cuatrimestre %vQ%v", n, c.Numero, c.Anio))
	}

	slog.Debug(fmt.Sprintf("encontradas %v ofertas de comisiones de materias en total", len(ofertasMaterias)))

	// ---
	// --- SINCRONIZACION
	// ---

	tx, _ := conn.Begin(context.Background())
	defer tx.Rollback(context.Background())

	nombres := make([]string, 0, len(ofertasMaterias))
	codigos := make([]string, 0, len(ofertasMaterias))
	for nombre, oferta := range ofertasMaterias {
		nombres = append(nombres, nombre)
		codigos = append(codigos, oferta.Codigo)
	}

	type MateriaSincronizada struct {
		Codigo                 string   `db:"codigo"`
		Nombre                 string   `db:"nombre"`
		DocentesMigrados       int      `db:"docentes_migrados"`
		ComentariosMigrados    int      `db:"comentarios_migrados"`
		CalificacionesMigradas int      `db:"calificaciones_migradas"`
		CodigosEquivalencias   []string `db:"codigos_equivalencias"`
	}

	rows, err = tx.Query(context.Background(), `
WITH materias_a_actualizar AS (
    SELECT
        mat.codigo AS codigo_antiguo,
        siu.codigo_siu AS codigo_nuevo,
        mat.nombre
    FROM
        materia mat
        JOIN (
            SELECT
                *
            FROM
                unnest($1::text[], $2::text[]) AS t (nombre_siu,
                codigo_siu)) siu ON lower(unaccent (mat.nombre)) = lower(unaccent (siu.nombre_siu))
    WHERE
        mat.codigo IS DISTINCT FROM siu.codigo_siu
        AND EXISTS (
            SELECT
                1
            FROM
                plan_materia pm
                INNER JOIN plan ON plan.codigo = pm.codigo_plan
            WHERE
                pm.codigo_materia = mat.codigo
                AND plan.esta_vigente)
),
materias_actualizadas AS (
    UPDATE
        materia mat
    SET
        codigo = maa.codigo_nuevo
    FROM
        materias_a_actualizar maa
    WHERE
        mat.codigo = maa.codigo_antiguo
    RETURNING
        mat.codigo AS codigo_nuevo,
        maa.codigo_antiguo,
        mat.nombre
),
equivalencias_por_materia AS (
    SELECT
        ma.codigo_nuevo,
        array_agg(e.codigo_materia_plan_anterior) AS codigos_equivalencias
FROM
    materias_actualizadas ma
    JOIN equivalencia e ON e.codigo_materia_plan_vigente = ma.codigo_antiguo
GROUP BY
    ma.codigo_nuevo
),
docentes_con_calificaciones AS (
    SELECT
        ma.codigo_nuevo,
        d.codigo AS codigo_docente_antiguo,
        d.nombre,
        d.resumen_comentarios,
        d.comentarios_ultimo_resumen,
        (
            SELECT
                count(*)
            FROM
                calificacion_dolly c
            WHERE
                c.codigo_docente = d.codigo) AS num_calificaciones
        FROM
            materias_actualizadas ma
            JOIN equivalencia e ON e.codigo_materia_plan_vigente = ma.codigo_antiguo
            JOIN docente d ON d.codigo_materia = e.codigo_materia_plan_anterior
),
docentes_a_migrar AS (
    SELECT DISTINCT ON (codigo_nuevo,
        nombre)
        codigo_nuevo,
        codigo_docente_antiguo,
        nombre,
        resumen_comentarios,
        comentarios_ultimo_resumen
    FROM
        docentes_con_calificaciones
    ORDER BY
        codigo_nuevo,
        nombre,
        num_calificaciones DESC
),
docentes_insertados AS (
INSERT INTO docente (nombre, codigo_materia, resumen_comentarios, comentarios_ultimo_resumen)
    SELECT
        nombre,
        codigo_nuevo,
        resumen_comentarios,
        comentarios_ultimo_resumen
    FROM
        docentes_a_migrar
    RETURNING
        codigo AS codigo_docente_nuevo,
        nombre,
        codigo_materia
),
mapeo_docentes AS (
    SELECT
        di.codigo_docente_nuevo,
        dm.codigo_docente_antiguo,
        di.codigo_materia
    FROM
        docentes_insertados di
        JOIN docentes_a_migrar dm ON di.nombre = dm.nombre
            AND di.codigo_materia = dm.codigo_nuevo
),
calificaciones_copiadas AS (
INSERT INTO calificacion_dolly (codigo_docente, acepta_critica, asistencia, buen_trato, claridad, clase_organizada, cumple_horarios, fomenta_participacion, panorama_amplio, responde_mails)
    SELECT
        m.codigo_docente_nuevo,
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
    JOIN mapeo_docentes m ON c.codigo_docente = m.codigo_docente_antiguo
RETURNING
    codigo_docente
),
comentarios_copiados AS (
INSERT INTO comentario (codigo_docente, codigo_cuatrimestre, contenido, es_de_dolly, fecha_creacion)
    SELECT
        m.codigo_docente_nuevo,
        cm.codigo_cuatrimestre,
        cm.contenido,
        cm.es_de_dolly,
        cm.fecha_creacion
    FROM
        comentario cm
        JOIN mapeo_docentes m ON cm.codigo_docente = m.codigo_docente_antiguo
    RETURNING
        codigo_docente
),
conteo_docentes AS (
    SELECT
        codigo_materia,
        count(*) AS docentes_migrados
    FROM
        docentes_insertados
    GROUP BY
        codigo_materia
),
conteo_calificaciones AS (
    SELECT
        md.codigo_materia,
        count(*) AS calificaciones_migradas
    FROM
        calificaciones_copiadas cc
        JOIN mapeo_docentes md ON cc.codigo_docente = md.codigo_docente_nuevo
    GROUP BY
        md.codigo_materia
),
conteo_comentarios AS (
    SELECT
        md.codigo_materia,
        count(*) AS comentarios_migrados
    FROM
        comentarios_copiados cmc
        JOIN mapeo_docentes md ON cmc.codigo_docente = md.codigo_docente_nuevo
    GROUP BY
        md.codigo_materia
)
SELECT
    ma.codigo_nuevo AS codigo,
    lower(unaccent (ma.nombre)) AS nombre,
    COALESCE(cd.docentes_migrados, 0)::int AS docentes_migrados,
    COALESCE(ccm.comentarios_migrados, 0)::int AS comentarios_migrados,
    COALESCE(ccal.calificaciones_migradas, 0)::int AS calificaciones_migradas,
    COALESCE(eq.codigos_equivalencias, ARRAY[]::text[]) AS codigos_equivalencias
FROM
    materias_actualizadas ma
    LEFT JOIN equivalencias_por_materia eq ON eq.codigo_nuevo = ma.codigo_nuevo
    LEFT JOIN conteo_docentes cd ON cd.codigo_materia = ma.codigo_nuevo
    LEFT JOIN conteo_comentarios ccm ON ccm.codigo_materia = ma.codigo_nuevo
    LEFT JOIN conteo_calificaciones ccal ON ccal.codigo_materia = ma.codigo_nuevo;
`, nombres, codigos)
	if err != nil {
		os.Exit(1)
	}

	materiasSincronizadas, err := pgx.CollectRows(rows, pgx.RowToStructByName[MateriaSincronizada])
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if len(materiasSincronizadas) == 0 {
		slog.Debug("no hay materias por sincronizar")
	} else {
		for _, m := range materiasSincronizadas {
			slog.Debug(
				fmt.Sprintf("sincronizado materia %s %s", m.Codigo, m.Nombre),
				"docentes", m.DocentesMigrados,
				"calificaciones", m.CalificacionesMigradas,
				"comentarios", m.ComentariosMigrados,
				"equivalencias", m.CodigosEquivalencias,
			)
		}
		slog.Debug(fmt.Sprintf("total de %d materias sincronizadas", len(materiasSincronizadas)))
	}

	if err := tx.Commit(context.Background()); err != nil {
		slog.Error(fmt.Sprintf("error al hacer commit de la transacción: %v", err))
		os.Exit(1)
	}

	// ---
	// ---
	// ---

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ofertasCarreras)
	})

	// log.Fatal(http.ListenAndServe(":8080", nil))
}
