package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5"
)

func iniciarServidor(conn *pgx.Conn, addr string, patches map[string]patchMateria) error {
	http.HandleFunc("GET /patches", func(w http.ResponseWriter, _ *http.Request) {
		handleGetAllPatches(w, patches)
	})
	http.HandleFunc("GET /patches/{codigoMateria}", func(w http.ResponseWriter, r *http.Request) {
		handleGetPatchMateria(w, r, conn, patches)
	})
	http.HandleFunc("PATCH /patches/{codigoMateria}", func(w http.ResponseWriter, r *http.Request) {
		handleResolverMateria(w, r, conn, patches)
	})

	slog.Info(fmt.Sprintf("servidor escuchando peticiones en dirección %v", addr))

	return http.ListenAndServe(addr, nil)
}

func handleGetAllPatches(w http.ResponseWriter, patches map[string]patchMateria) {
	type patchRes struct {
		Codigo   string `json:"codigo"`
		Nombre   string `json:"nombre"`
		Docentes int    `json:"docentes"`
	}

	patchesRes := make([]patchRes, 0)
	for cod, pat := range patches {
		patchesRes = append(
			patchesRes,
			patchRes{Codigo: cod, Nombre: pat.Nombre, Docentes: len(pat.Docentes)},
		)
	}

	slices.SortFunc(patchesRes, func(a, b patchRes) int {
		if a.Docentes != b.Docentes {
			return b.Docentes - a.Docentes
		}
		if a.Codigo < b.Codigo {
			return -1
		}
		if a.Codigo > b.Codigo {
			return 1
		}
		return 0
	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(patchesRes); err != nil {
		slog.Error(
			fmt.Sprintf("error serializando patches de materias: %v", err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleGetPatchMateria(
	w http.ResponseWriter,
	r *http.Request,
	conn *pgx.Conn,
	patches map[string]patchMateria,
) {
	codigoMateria := r.PathValue("codigoMateria")
	patch, ok := patches[codigoMateria]

	if !ok {
		http.Error(
			w,
			fmt.Sprintf(
				"Patch de actualización para materia %v no encontrado",
				codigoMateria,
			),
			http.StatusNotFound,
		)
		return
	}

	docentesResueltos, err := getDocentesResueltos(conn, codigoMateria, patch.Catedras)
	if err != nil {
		slog.Error(
			fmt.Sprintf(
				"error obteniendo estado de resolución de docentes de materia %v: %v",
				codigoMateria,
				err,
			),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	catedrasResponse := make([]catedraResponse, 0, len(patch.Catedras))
	for _, cat := range patch.Catedras {
		docentesResponse := make([]docenteCatedraResponse, 0, len(cat.Docentes))

		for _, doc := range cat.Docentes {
			codigoDocente := (*string)(nil)
			if docentesPorCatedra, ok := docentesResueltos[cat.Codigo]; ok {
				if codigo, ok := docentesPorCatedra[doc.Nombre]; ok {
					codigoDocente = codigo
				}
			}

			docentesResponse = append(docentesResponse, docenteCatedraResponse{
				Nombre:           doc.Nombre,
				CodigoYaResuelto: codigoDocente,
			})
		}

		slices.SortFunc(docentesResponse, func(a, b docenteCatedraResponse) int {
			return strings.Compare(a.Nombre, b.Nombre)
		})

		catedrasResponse = append(catedrasResponse, catedraResponse{
			Docentes: docentesResponse,
			Resuelta: cat.Resuelta,
		})
	}

	type patchMateriaRes struct {
		materia
		DocentesSinResolver []patchDocente    `json:"docentes_sin_resolver"`
		Catedras            []catedraResponse `json:"catedras"`
		cuatrimestre        `                  json:"cuatrimestre"`
	}

	res := patchMateriaRes{
		materia:             patch.materia,
		DocentesSinResolver: patch.Docentes,
		Catedras:            catedrasResponse,
		cuatrimestre:        patch.cuatrimestre,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		slog.Error(
			fmt.Sprintf(
				"error serializando patches de materias %v (%v): %v",
				patch.Codigo,
				patch.Nombre,
				err,
			),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type resolucionMateria struct {
	CodigosYaResueltos   []string                     `json:"codigos_ya_resueltos"`
	ResolucionesDocentes map[string]resolucionDocente `json:"resoluciones_actuales"`
}

type resolucionDocente struct {
	NombreDb    string  `json:"nombre_db"`
	CodigoMatch *string `json:"codigo_match"`
}

func handleResolverMateria(
	w http.ResponseWriter,
	r *http.Request,
	conn *pgx.Conn,
	patches map[string]patchMateria,
) {
	var res resolucionMateria

	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		slog.Error(fmt.Sprintf("error deserializando JSON de docentes: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	codigoMateria := r.PathValue("codigoMateria")
	patch := patches[codigoMateria]

	if err := resolverMateria(conn, patch, res); err != nil {
		slog.Error(fmt.Sprintf("error aplicando resolución de materia %v: %v", codigoMateria, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	delete(patches, codigoMateria)

	w.WriteHeader(http.StatusNoContent)
}
