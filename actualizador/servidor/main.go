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

	nombres := make([]string, 0, len(ofertasMaterias))
	codigos := make([]string, 0, len(ofertasMaterias))
	for nombre, oferta := range ofertasMaterias {
		nombres = append(nombres, nombre)
		codigos = append(codigos, oferta.Codigo)
	}

	tag, err := conn.Exec(context.Background(), `
UPDATE materia mat
SET 
    codigo = siu.codigo_siu
FROM (
    SELECT * FROM unnest($1::text[], $2::text[]) AS t(nombre_siu, codigo_siu)
) siu
WHERE 
    lower(unaccent(mat.nombre)) = lower(unaccent(siu.nombre_siu))
    AND mat.codigo IS DISTINCT FROM siu.codigo_siu
    AND EXISTS (
        SELECT 1 FROM plan_materia pm
        INNER JOIN plan ON plan.codigo = pm.codigo_plan
        WHERE pm.codigo_materia = mat.codigo AND plan.esta_vigente
    );
`, nombres, codigos)
	if err != nil {
		slog.Error(err.Error())
	}

	slog.Info(fmt.Sprintf("actualizadas %v materias", tag.RowsAffected()))

	// ---
	// ---
	// ---

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ofertasCarreras)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
