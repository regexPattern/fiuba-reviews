package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"

	"github.com/jackc/pgx/v5"
)

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
	patches map[string]patchMateria,
) {
	codigoMateria := r.PathValue("codigoMateria")
	patchRes := patches[codigoMateria]

	// for i := range patchRes.Catedras {
	// 	// slices.SortFunc(patchRes.Catedras[i].DocentesResueltos, func(a, b docente) int {
	// 	// 	return strings.Compare(a.Nombre, b.Nombre)
	// 	// })
	// 	// slices.SortFunc(patchRes.Catedras[i].DocentesNoResueltos, func(a, b docente) int {
	// 	// 	return strings.Compare(a.Nombre, b.Nombre)
	// 	// })
	// }

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(patchRes); err != nil {
		slog.Error(
			fmt.Sprintf(
				"error serializando patches de materias %v (%v): %v",
				patchRes.Codigo,
				patchRes.Nombre,
				err,
			),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleAplicarPatchMateria(
	w http.ResponseWriter,
	r *http.Request,
	conn *pgx.Conn,
	patches map[string]patchMateria,
) {
	var resoluciones map[string]struct {
		NombreDb    string  `json:"nombre_db"`
		CodigoMatch *string `json:"codigo_match"`
	}

	if err := json.NewDecoder(r.Body).Decode(&resoluciones); err != nil {
		slog.Error(fmt.Sprintf("error parseando JSON de docentes: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for k, r := range resoluciones {
		fmt.Println(k, r)
	}

	codigoMateria := r.PathValue("codigoMateria")
	patch := patches[codigoMateria]

	_ = patch

	w.WriteHeader(http.StatusNoContent)
}
