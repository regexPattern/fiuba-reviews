package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
)

func main() {
	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)
	slog.SetDefault(slog.New(logger))

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("")
		os.Exit(1)
	}

	slog.Info("conexión establecida con la base de datos")

	ofertasMaterias, err := GetOfertasMaterias(conn)
	if err != nil {
		slog.Error(fmt.Sprintf("error obteniendo ofertas de comisiones de materias: %v", err))
		os.Exit(1)
	}

	codigos := make([]string, 0, len(ofertasMaterias))
	nombres := make([]string, 0, len(ofertasMaterias))
	for nombre, oferta := range ofertasMaterias {
		codigos = append(codigos, oferta.Codigo)
		nombres = append(nombres, nombre)
	}

	if err := SyncMateriasDb(conn, codigos, nombres); err != nil {
		slog.Error(
			fmt.Sprintf("error sincronizando materias de la base de datos con el siu: %v", err),
		)
		os.Exit(1)
	}

	materiasAActualizar, err := GetMateriasAActualizar(conn, codigos)
	if err != nil {
		slog.Error(fmt.Sprintf("error obteniendo materias a actualizar: %v", err))
		os.Exit(1)
	}

	fmt.Println(materiasAActualizar)
}

func GetOfertasMaterias(conn *pgx.Conn) (map[string]UltimaOfertaMateria, error) {
	rows, err := conn.Query(context.Background(), `
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
	if err != nil {
		return nil, fmt.Errorf("error consultando ofertas de comisiones de carreras: %w", err)
	}

	ofertasCarreras, err := pgx.CollectRows(rows, pgx.RowToStructByName[OfertaCarrera])
	if err != nil {
		return nil, fmt.Errorf("error procesando ofertas de comisiones de carreras")
	}

	slog.Debug(
		fmt.Sprintf("encontradas %v ofertas de comisiones de carreras", len(ofertasCarreras)),
	)

	ofertasMaterias := make(map[string]UltimaOfertaMateria)
	materiasPorCuatrimestre := make(map[Cuatrimestre]int)

	for _, oc := range ofertasCarreras {
		for _, om := range oc.Materias {
			if _, ok := ofertasMaterias[om.Nombre]; !ok {
				ofertasMaterias[om.Nombre] = UltimaOfertaMateria{
					OfertaMateria: om,
					Cuatrimestre:  oc.Cuatrimestre,
				}

				materiasPorCuatrimestre[oc.Cuatrimestre]++
			}
		}
	}

	for c, n := range materiasPorCuatrimestre {
		slog.Debug(
			fmt.Sprintf(
				"encontradas %v ofertas de comisiones de materias de cuatrimestre %vQ%v",
				n,
				c.Numero,
				c.Anio,
			),
		)
	}

	slog.Debug(
		fmt.Sprintf(
			"encontradas %v ofertas de comisiones de materias en total",
			len(ofertasMaterias),
		),
	)

	return ofertasMaterias, nil
}

func SyncMateriasDb(conn *pgx.Conn, codigos, nombres []string) error {
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("error iniciando transacción de sincronización de materias: %w", err)
	}

	defer tx.Rollback(context.Background())

	type MateriaSincronizada struct {
		Codigo                 string   `db:"codigo"`
		Nombre                 string   `db:"nombre"`
		DocentesMigrados       int      `db:"docentes_migrados"`
		ComentariosMigrados    int      `db:"comentarios_migrados"`
		CalificacionesMigradas int      `db:"calificaciones_migradas"`
		CodigosEquivalencias   []string `db:"codigos_equivalencias"`
	}

	rows, err := tx.Query(context.Background(), `
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
		return fmt.Errorf("error ejecutando consulta de sincronización de materias: %w", err)
	}

	materiasSincronizadas, err := pgx.CollectRows(rows, pgx.RowToStructByName[MateriaSincronizada])
	if err != nil {
		return fmt.Errorf("error procesando materias sincronizadas: %w", err)
	}

	for _, m := range materiasSincronizadas {
		slog.Debug(
			fmt.Sprintf("sincronizado materia %s %s", m.Codigo, m.Nombre),
			"docentes", m.DocentesMigrados,
			"calificaciones", m.CalificacionesMigradas,
			"comentarios", m.ComentariosMigrados,
			"equivalencias", m.CodigosEquivalencias,
		)
	}

	slog.Debug(fmt.Sprintf("sincronizadas %d materias en total", len(materiasSincronizadas)))

	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf(
			"error haciendo commit de la transacción de sincronización de materias: %w",
			err,
		)
	}

	if err := WarnMateriasNoRegistradas(conn, codigos, nombres); err != nil {
		return fmt.Errorf("error verificando materias no registradas en la base de datos: %w", err)
	}

	return nil
}

func WarnMateriasNoRegistradas(conn *pgx.Conn, codigos, nombres []string) error {
	rows, err := conn.Query(context.Background(), `
		SELECT 
			siu.codigo_siu AS codigo,
			siu.nombre_siu AS nombre
		FROM unnest($1::text[], $2::text[]) AS siu(nombre_siu, codigo_siu)
		WHERE NOT EXISTS (
			SELECT 1 
			FROM materia mat
			WHERE lower(unaccent(mat.nombre)) = lower(unaccent(siu.nombre_siu))
		);
	`, nombres, codigos)
	if err != nil {
		return fmt.Errorf("error consultando materias no registradas: %w", err)
	}

	materiasNoRegistradas, err := pgx.CollectRows(rows, pgx.RowToStructByName[Materia])
	if err != nil {
		return fmt.Errorf("error procesando materias no registradas: %v", err)
	}

	for _, m := range materiasNoRegistradas {
		slog.Warn(
			fmt.Sprintf("materia %s %s no registrada en la base de datos", m.Codigo, m.Nombre),
		)
	}

	return nil
}

func GetMateriasAActualizar(conn *pgx.Conn, codigos []string) ([]Materia, error) {
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
	`, codigos)
	if err != nil {
		return nil, fmt.Errorf("error consultando materias pendientes de actualización: %w", err)
	}

	materias, err := pgx.CollectRows(rows, pgx.RowToStructByName[Materia])
	if err != nil {
		return nil, fmt.Errorf("error procesando materias pendientes de actualización: %v", err)
	}

	// TODO: Quedarme solo con las materias que realmente tienen cambios

	// Las materias a actualizar que tenemos en la variable anterior `materias` incluyen todas las
	// materias que no han sido actualizadas al ultimo cuatrimestre. Pero realmente podrian ya estar
	// totalmente al dia en cuanto a su oferta de comisiones y planilla de docentes, que son las
	// cosas que a mi me interesa actualizar en mi app.

	// La oferta de comisiones de cada materia proveniente del siu posee tanto un listado de
	// comisiones/catedras como de los docentes que las componen. Actualizar la oferta de comisiones
	// en la base de datos implica tener armadas las catedras que estan presentes en esta ultima
	// version de la oferta de comisiones que se tiene para cada materia disponible. Como una
	// catedra no es mucho mas que una agrupacion de docentes, esto implica que deben estar
	// registrados los docentes listados en el siu en la base de datos. Contamos con el listado de
	// docentes que necesitamos encontrar en la base de datos para armar las catedras de una
	// materia, son justamente los que traemos del siu.

	// # Matcheo de docentes
	//
	// Para cada docente de esta lista se pueden presentar 3 situaciones:
	// 1. Hay un docente en la base de datos cuyo `nombre_siu` coincide con el nombre del docente
	// encontrado en el siu. En este caso podemos considerar al docente como resuelto.
	// 2. No hay un docente en la base de datos cuyo `nombre_siu` coincida, pero hay varios docentes
	// cuyo campo `nombre` coincide con el nombre del docente encontrado en el siu. Estos otros
	// docentes no deben tener un `nombre_siu` asignado, ya que esto significa que ya estan
	// "asociados" con otro docente del siu. Para esto se utilizara fuzzy matching.
	// 3. No hay ningun docente que matchee en nombre con el docente del siu. En este caso significa
	// que estamos ante un docente totalmente nuevo que hay que registrar.

	// # Armado de catedras/comisiones
	//
	// Una catedra no es mas que una agrupacion de docentes. Si bien en la base de datos estan
	// identificadas con un codigo, lo verdaderamente representativo de las catedras es su nombre,
	// que se puede obtener concatenando los nombres de los docentes que la componen. Asi, por
	// ejemplo, la catedra de los docentes "Carlos Castillo", "Irene Cardona" y "Lionel Messi",
	// seria la catedra "Cardona-Castillo-Messi", ordenada en orden alfabetico.
	//
	// NOTE: Estaria bueno agregar el campo `apellido` a los docentes en la base de datos, para que
	// tengan un nombre que los represente de manera corta en vez de tener que usar siempre su
	// nombre completo real. Quiza incluso podamos ocupar el campo `nombre` que ya existe.
	//
	// De un cuatrimestre a otro, las catedras pueden continuar existiendo. Es decir, si la catedra
	// "Cardona-Castillo-Messi" ya existe, no hay necesidad de eliminar este registro y crear uno
	// nuevo con los mismos docentes. Si se diera el caso en el que se agregara un docente
	// adicional, para simplificar, ahi si deberiamos crear una nueva catedra. Lo mismo para si se
	// retira un docente.
	//
	// A medida vayan pasando los cuatrimestres de actualizacion, van a ir quedando catedras
	// registradas en la base de datos que no existan. Para esto tenemos dos opciones: borramos
	// siempre todas las catedras que no se mantienen en cada cuatrimestre, o usamos algun tipo de
	// flag para saber si estan activas o son de cuatrimestres anteriores. Para elegir alguna de
	// estas dos opciones, hay que tener en consideracion que pasa cuando una materia no tiene
	// actualizaciones para este cuatrimestre.

	// # Problema principal
	//
	// Justo ahi recae el problema. Quiero filtrar las materias que si tengan actualizaciones
	// pendientes, ya sea para crear registrar docentes que no estan registrados, y/o para registrar
	// catedras (reagrupaciones de estos docentes). Quiero en un arreglo tener solo estas materias.
	// La informacion que debe quedar en el arreglo final es:
	// - codigo y nombre (normalizado con lower(unaccent())) de la materia
	// - listado de docentes que no estan resueltos, y sus respectivos matches (tanto el nombre
	// fuzzy matcheado como el codigo del match en la db).
	// - listado de catedras a crear (nombre concatenado de los docentes).
	// Pero, nuevamente, el tema es que solo quiero esta informacion de las materias que tienen
	// actualizaciones.

	return materias, nil
}
