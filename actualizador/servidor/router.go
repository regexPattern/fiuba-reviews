package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func handleGetAllPatches(w http.ResponseWriter, patches map[string]patchMateria) {
	type PatchRes struct {
		Codigo string `json:"codigo"`
		Nombre string `json:"nombre"`
	}

	patchesRes := make([]PatchRes, 0)
	for cod, pat := range patches {
		patchesRes = append(patchesRes, PatchRes{Codigo: cod, Nombre: pat.Nombre})
	}

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
	patchRes, ok := patches[codMateria]
	if !ok {
		http.Error(
			w,
			fmt.Sprintf("patch de materia %v no encontrado", codMateria),
			http.StatusNotFound,
		)
		return
	}

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
		// delete(patches, codMateria)
	}
	w.WriteHeader(http.StatusNoContent)
}
