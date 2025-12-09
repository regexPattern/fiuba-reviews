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
	for codigo, oferta := range ofertasMaterias {
		codigos = append(codigos, codigo)
		nombres = append(nombres, oferta.Nombre)
	}

	if err := SyncMateriasDb(conn, codigos, nombres); err != nil {
		slog.Error(
			fmt.Sprintf("error sincronizando materias de la base de datos con el siu: %v", err),
		)
		os.Exit(1)
	}

	materiasPendientes, err := GetMateriasPendientes(conn, codigos, ofertasMaterias)
	if err != nil {
		slog.Error(fmt.Sprintf("error obteniendo materias a actualizar: %v", err))
		os.Exit(1)
	}

	slog.Debug(
		fmt.Sprintf(
			"encontradas %d materias con actualizaciones pendientes",
			len(materiasPendientes),
		),
	)

	// for i, mat := range materiasAActualizar {
	// 	slog.Debug(fmt.Sprintf("materia %d: %s - %d docentes pendientes - %d catedras nuevas",
	// 		i+1, mat.Nombre, len(mat.DocentesPendientes), len(mat.CatedrasNuevas)))
	// }

	for _, m := range materiasPendientes {
		fmt.Println(m)
	}
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

	ofertasMaterias := make(map[string]UltimaOfertaMateria) // clave: codigo de materia
	materiasPorCuatrimestre := make(map[Cuatrimestre]int)

	for _, oc := range ofertasCarreras {
		for _, om := range oc.Materias {
			if _, ok := ofertasMaterias[om.Codigo]; !ok {
				ofertasMaterias[om.Codigo] = UltimaOfertaMateria{
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

func GetMateriasPendientes(
	conn *pgx.Conn,
	codigos []string,
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
	`, codigos)
	if err != nil {
		return nil, fmt.Errorf("error consultando materias candidatas a actualizarse: %w", err)
	}

	materiasCandidatas, err := pgx.CollectRows(rows, pgx.RowToStructByName[Materia])
	if err != nil {
		return nil, fmt.Errorf("error procesando materias candidatas a actualizarse: %v", err)
	}

	materiasPendientes := make([]any, 0, len(materiasCandidatas))

	for _, m := range materiasCandidatas {
		oferta, ok := ofertasMaterias[m.Codigo]
		if !ok {
			continue
		}

		if tiene, err := OfertaTieneCambios(conn, oferta); err != nil {
			return nil, fmt.Errorf(
				"error determinando si oferta de materia %v tiene cambios: %w",
				m.Codigo,
				err,
			)
		} else if tiene {
			materiasPendientes = append(materiasPendientes, oferta)
		}
	}

	return nil, nil
}

func OfertaTieneCambios(conn *pgx.Conn, oferta UltimaOfertaMateria) (bool, error) {
	// // 2. Para cada materia candidata, analizar si tiene actualizaciones pendientes
	// for _, mat := range materiasCandidatas {
	// 	oferta, ok := ofertasMaterias[mat.Codigo]
	// 	if !ok {
	// 		continue
	// 	}
	//
	// 	// 3. Extraer todos los docentes únicos de la oferta
	// 	docentesSiu := extraerDocentesUnicos(oferta.Catedras)
	//
	// 	// 4. Buscar matches para cada docente
	// 	docentesPendientes, err := resolverDocentes(conn, mat.Codigo, docentesSiu)
	// 	if err != nil {
	// 		return nil, fmt.Errorf(
	// 			"error resolviendo docentes para materia %s: %w",
	// 			mat.Codigo,
	// 			err,
	// 		)
	// 	}
	//
	// 	// 5. Verificar qué cátedras son nuevas
	// 	catedrasNuevas, err := verificarCatedras(
	// 		conn,
	// 		mat.Codigo,
	// 		oferta.Catedras,
	// 		docentesPendientes,
	// 	)
	// 	if err != nil {
	// 		return nil, fmt.Errorf(
	// 			"error verificando catedras para materia %s: %w",
	// 			mat.Codigo,
	// 			err,
	// 		)
	// 	}
	//
	// 	// 6. Si hay algo pendiente, incluir en resultado
	// 	if len(docentesPendientes) > 0 || len(catedrasNuevas) > 0 {
	// 		materiasPendientes = append(materiasPendientes, MateriaConActualizaciones{
	// 			Codigo:             mat.Codigo,
	// 			Nombre:             mat.Nombre,
	// 			DocentesPendientes: docentesPendientes,
	// 			CatedrasNuevas:     catedrasNuevas,
	// 		})
	// 	}
	// }

	return false, nil
}

// // extraerDocentesUnicos extrae todos los docentes únicos de las catedras de una oferta
// func extraerDocentesUnicos(catedras []Catedra) map[string]string {
// 	docentes := make(map[string]string) // nombre_siu -> rol
//
// 	for _, catedra := range catedras {
// 		for _, docente := range catedra.Docentes {
// 			// Guardar el rol más "importante" si hay duplicados
// 			if rolExistente, existe := docentes[docente.Nombre]; !existe ||
// 				(docente.Rol == "profesor titular" || docente.Rol == "profesor asociado") {
// 				docentes[docente.Nombre] = docente.Rol
// 			} else if existe && (docente.Rol == "profesor titular" || docente.Rol == "profesor asociado")
// {
// 				// Reemplazar si el nuevo rol es más importante
// 				docentes[docente.Nombre] = docente.Rol
// 			} else {
// 				// Mantener el rol existente
// 				docentes[docente.Nombre] = rolExistente
// 			}
// 		}
// 	}
//
// 	return docentes
// }
//
// // resolverDocentes busca matches para docentes del SIU en la base de datos
// func resolverDocentes(
// 	conn *pgx.Conn,
// 	codigoMateria string,
// 	docentesSiu map[string]string,
// ) ([]DocentePendiente, error) {
// 	if len(docentesSiu) == 0 {
// 		return nil, nil
// 	}
//
// 	// Convertir mapa a array para la query
// 	nombresSiu := make([]string, 0, len(docentesSiu))
// 	for nombre := range docentesSiu {
// 		nombresSiu = append(nombresSiu, nombre)
// 	}
//
// 	// Query para buscar matches fuzzy
// 	rows, err := conn.Query(context.Background(), `
// 		WITH docentes_siu AS (
// 			SELECT unnest($1::text[]) AS nombre_siu
// 		),
// 		matches AS (
// 			SELECT
// 				ds.nombre_siu,
// 				d.codigo,
// 				d.nombre AS nombre_db,
// 				d.nombre_siu AS nombre_siu_actual,
// 				similarity(
// 					lower(unaccent(d.nombre)),
// 					lower(unaccent(split_part(ds.nombre_siu, ' ', 1)))
// 				) AS similitud
// 			FROM docentes_siu ds
// 			LEFT JOIN docente d ON d.codigo_materia = $2
// 			WHERE d.nombre_siu IS NULL  -- solo docentes no resueltos
// 			  AND similarity(
// 					lower(unaccent(d.nombre)),
// 					lower(unaccent(split_part(ds.nombre_siu, ' ', 1)))
// 				  ) > 0.5  -- umbral mínimo
// 		)
// 		SELECT
// 			ds.nombre_siu,
// 			COALESCE(
// 				json_agg(
// 					json_build_object(
// 						'codigo', matches.codigo,
// 						'nombre_db', matches.nombre_db,
// 						'similitud', matches.similitud
// 					) ORDER BY matches.similitud DESC
// 				) FILTER (WHERE matches.codigo IS NOT NULL),
// 				'[]'::json
// 			) AS matches
// 		FROM docentes_siu ds
// 		LEFT JOIN matches ON matches.nombre_siu = ds.nombre_siu
// 		WHERE NOT EXISTS (
// 			-- Excluir docentes ya resueltos exactamente
// 			SELECT 1 FROM docente d2
// 			WHERE d2.codigo_materia = $2
// 			  AND d2.nombre_siu = ds.nombre_siu
// 		)
// 		GROUP BY ds.nombre_siu;
// 	`, nombresSiu, codigoMateria)
// 	if err != nil {
// 		return nil, fmt.Errorf("error buscando matches de docentes: %w", err)
// 	}
// 	defer rows.Close()
//
// 	var docentesPendientes []DocentePendiente
//
// 	for rows.Next() {
// 		var nombreSiu string
// 		var matchesJSON string
//
// 		if err := rows.Scan(&nombreSiu, &matchesJSON); err != nil {
// 			return nil, fmt.Errorf("error escaneando matches de docentes: %w", err)
// 		}
//
// 		// Parsear matches JSON
// 		var matches []DocenteMatch
// 		if matchesJSON != "null" && matchesJSON != "[]" {
// 			if err := json.Unmarshal([]byte(matchesJSON), &matches); err != nil {
// 				return nil, fmt.Errorf("error parseando matches JSON: %w", err)
// 			}
// 		}
//
// 		// Solo incluir si no hay matches o hay múltiples (requiere intervención manual)
// 		if len(matches) == 0 || len(matches) > 1 ||
// 			(len(matches) == 1 && matches[0].Similitud < 0.8) {
// 			docentesPendientes = append(docentesPendientes, DocentePendiente{
// 				NombreSiu:       nombreSiu,
// 				Rol:             docentesSiu[nombreSiu],
// 				PosiblesMatches: matches,
// 			})
// 		}
// 	}
//
// 	return docentesPendientes, nil
// }
//
// // verificarCatedras identifica qué cátedras del SIU son nuevas
// func verificarCatedras(
// 	conn *pgx.Conn,
// 	codigoMateria string,
// 	catedrasSiu []Catedra,
// 	docentesPendientes []DocentePendiente,
// ) ([]CatedraNueva, error) {
// 	if len(catedrasSiu) == 0 {
// 		return nil, nil
// 	}
//
// 	// Crear mapa de docentes pendientes para referencia rápida
// 	docentesPendientesMap := make(map[string]bool)
// 	for _, dp := range docentesPendientes {
// 		docentesPendientesMap[dp.NombreSiu] = true
// 	}
//
// 	var catedrasNuevas []CatedraNueva
//
// 	for _, catedraSiu := range catedrasSiu {
// 		// Verificar si todos los docentes de esta cátedra están resueltos
// 		tieneDocentesPendientes := false
// 		for _, docente := range catedraSiu.Docentes {
// 			if docentesPendientesMap[docente.Nombre] {
// 				tieneDocentesPendientes = true
// 				break
// 			}
// 		}
//
// 		// Si hay docentes pendientes, la cátedra es nueva por definición
// 		if tieneDocentesPendientes {
// 			catedraNueva := CatedraNueva{
// 				Nombre:   generarNombreCatedra(catedraSiu.Docentes),
// 				Docentes: make([]DocenteCatedra, 0, len(catedraSiu.Docentes)),
// 			}
//
// 			for _, docente := range catedraSiu.Docentes {
// 				dc := DocenteCatedra{
// 					NombreSiu: docente.Nombre,
// 					Rol:       docente.Rol,
// 				}
// 				catedraNueva.Docentes = append(catedraNueva.Docentes, dc)
// 			}
//
// 			catedrasNuevas = append(catedrasNuevas, catedraNueva)
// 			continue
// 		}
//
// 		// Si todos los docentes están resueltos, verificar si la cátedra existe
// 		existe, err := verificarCatedraExistente(conn, codigoMateria, catedraSiu.Docentes)
// 		if err != nil {
// 			// Si hay error porque no se encontró algún docente, significa que hay docentes
// 			// pendientes
// 			// En ese caso, la cátedra se considera nueva
// 			if strings.Contains(err.Error(), "docente no encontrado") {
// 				catedraNueva := CatedraNueva{
// 					Nombre:   generarNombreCatedra(catedraSiu.Docentes),
// 					Docentes: make([]DocenteCatedra, 0, len(catedraSiu.Docentes)),
// 				}
//
// 				for _, docente := range catedraSiu.Docentes {
// 					dc := DocenteCatedra{
// 						NombreSiu: docente.Nombre,
// 						Rol:       docente.Rol,
// 					}
// 					catedraNueva.Docentes = append(catedraNueva.Docentes, dc)
// 				}
//
// 				catedrasNuevas = append(catedrasNuevas, catedraNueva)
// 				continue
// 			}
// 			return nil, fmt.Errorf("error verificando si existe catedra: %w", err)
// 		}
//
// 		if !existe {
// 			// La cátedra no existe, es nueva
// 			catedraNueva := CatedraNueva{
// 				Nombre:   generarNombreCatedra(catedraSiu.Docentes),
// 				Docentes: make([]DocenteCatedra, 0, len(catedraSiu.Docentes)),
// 			}
//
// 			// Para docentes resueltos, obtener sus códigos
// 			for _, docente := range catedraSiu.Docentes {
// 				codigoDocente, err := obtenerCodigoDocentePorNombreSiu(
// 					conn,
// 					codigoMateria,
// 					docente.Nombre,
// 				)
// 				if err != nil {
// 					return nil, fmt.Errorf("error obteniendo código de docente: %w", err)
// 				}
//
// 				dc := DocenteCatedra{
// 					NombreSiu:     docente.Nombre,
// 					Rol:           docente.Rol,
// 					CodigoDocente: &codigoDocente,
// 				}
// 				catedraNueva.Docentes = append(catedraNueva.Docentes, dc)
// 			}
//
// 			catedrasNuevas = append(catedrasNuevas, catedraNueva)
// 		}
// 	}
//
// 	return catedrasNuevas, nil
// }
//
// // generarNombreCatedra genera el nombre representativo de una cátedra
// func generarNombreCatedra(docentes []Docente) string {
// 	apellidos := make([]string, 0, len(docentes))
// 	for _, d := range docentes {
// 		// Extraer primer token (apellido) del nombre del SIU
// 		parts := strings.Split(d.Nombre, " ")
// 		if len(parts) > 0 {
// 			apellido := strings.Title(parts[0])
// 			apellidos = append(apellidos, apellido)
// 		}
// 	}
// 	sort.Strings(apellidos)
// 	return strings.Join(apellidos, "-")
// }
//
// // verificarCatedraExistente verifica si una cátedra con esos docentes ya existe
// func verificarCatedraExistente(
// 	conn *pgx.Conn,
// 	codigoMateria string,
// 	docentesSiu []Docente,
// ) (bool, error) {
// 	// Obtener códigos de docentes para esta materia
// 	codigosDocentes := make([]string, 0, len(docentesSiu))
// 	for _, docente := range docentesSiu {
// 		codigo, err := obtenerCodigoDocentePorNombreSiu(conn, codigoMateria, docente.Nombre)
// 		if err != nil {
// 			return false, err
// 		}
// 		codigosDocentes = append(codigosDocentes, codigo)
// 	}
//
// 	// Ordenar códigos para comparación
// 	sort.Strings(codigosDocentes)
//
// 	// Verificar si existe una cátedra con exactamente estos docentes
// 	var existe bool
// 	err := conn.QueryRow(context.Background(), `
// 		WITH catedras_existentes AS (
// 			SELECT
// 				c.codigo,
// 				array_agg(cd.codigo_docente ORDER BY cd.codigo_docente) AS docentes
// 			FROM catedra c
// 			JOIN catedra_docente cd ON cd.codigo_catedra = c.codigo
// 			WHERE c.codigo_materia = $1
// 			GROUP BY c.codigo
// 		)
// 		SELECT EXISTS (
// 			SELECT 1 FROM catedras_existentes
// 			WHERE docentes = $2::uuid[]
// 		);
// 	`, codigoMateria, codigosDocentes).Scan(&existe)
// 	if err != nil {
// 		// Si no hay filas, significa que no hay catedras existentes
// 		if err.Error() == "no rows in result set" {
// 			return false, nil
// 		}
// 		return false, err
// 	}
//
// 	return existe, nil
// }
//
// // obtenerCodigoDocentePorNombreSiu obtiene el código de un docente por su nombre_siu
// func obtenerCodigoDocentePorNombreSiu(
// 	conn *pgx.Conn,
// 	codigoMateria, nombreSiu string,
// ) (string, error) {
// 	var codigo string
// 	err := conn.QueryRow(context.Background(), `
// 		SELECT codigo FROM docente
// 		WHERE codigo_materia = $1 AND nombre_siu = $2
// 	`, codigoMateria, nombreSiu).Scan(&codigo)
// 	if err != nil {
// 		// Si no hay filas, significa que el docente no existe
// 		if err.Error() == "no rows in result set" {
// 			return "", fmt.Errorf("docente no encontrado: %s", nombreSiu)
// 		}
// 		return "", err
// 	}
//
// 	return codigo, nil
// }
