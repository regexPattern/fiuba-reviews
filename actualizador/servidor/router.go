package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
)

func handleGetAllPatches(w http.ResponseWriter, patches map[string]patchMateria) {
	type PatchRes struct {
		Codigo   string `json:"codigo"`
		Nombre   string `json:"nombre"`
		Docentes int    `json:"docentes"`
	}

	patchesRes := make([]PatchRes, 0)
	for cod, pat := range patches {
		patchesRes = append(
			patchesRes,
			PatchRes{Codigo: cod, Nombre: pat.Nombre, Docentes: len(pat.Docentes)},
		)
	}

	slices.SortFunc(patchesRes, func(a, b PatchRes) int {
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
	codMateria := r.PathValue("codigoMateria")
	patchRes := patches[codMateria]

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

func handleApplyPatchMateria(
	w http.ResponseWriter,
	r *http.Request,
	patches map[string]patchMateria,
) {
	codMateria := r.PathValue("codigoMateria")
	if pat, ok := patches[codMateria]; ok {
		slog.Info(
			fmt.Sprintf("eliminado patch de materia %v (%v) del registro", pat.Codigo, pat.Nombre),
		)
	}
	w.WriteHeader(http.StatusNoContent)
}
