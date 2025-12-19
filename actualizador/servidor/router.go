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

type patchMateriaResponse struct {
	materia
	DocentesSinResolver []patchDocente    `json:"docentes_sin_resolver"`
	Catedras            []catedraResponse `json:"catedras"`
	cuatrimestre        `json:"cuatrimestre"`
}

func handleGetPatchMateria(
	w http.ResponseWriter,
	r *http.Request,
	conn *pgx.Conn,
	patches map[string]patchMateria,
) {
	codigoMateria := r.PathValue("codigoMateria")
	patchRes := patches[codigoMateria]

	docentesResueltos, err := getDocentesResueltosDeCatedras(conn, codigoMateria, patchRes.Catedras)
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

	catedrasResponse := make([]catedraResponse, 0, len(patchRes.Catedras))
	for _, cat := range patchRes.Catedras {
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

	res := patchMateriaResponse{
		materia:             patchRes.materia,
		DocentesSinResolver: patchRes.Docentes,
		Catedras:            catedrasResponse,
		cuatrimestre:        patchRes.cuatrimestre,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
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
	var resoluciones struct {
		CodigosYaResueltos   []string `json:"codigos_ya_resueltos"`
		ResolucionesActuales map[string]struct {
			NombreDb    string  `json:"nombre_db"`
			CodigoMatch *string `json:"codigo_match"`
		} `json:"resoluciones_actuales"`
	}

	if err := json.NewDecoder(r.Body).Decode(&resoluciones); err != nil {
		slog.Error(fmt.Sprintf("error deserializando JSON de docentes: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	codigoMateria := r.PathValue("codigoMateria")
	patch := patches[codigoMateria]

	fmt.Println(resoluciones)

	// if err := aplicarPatchMateria(conn, patch, resoluciones); err != nil {
	// 	slog.Error(fmt.Sprintf("error aplicando patch de materia %v: %v", codigoMateria, err))
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	_ = patch

	w.WriteHeader(http.StatusNoContent)
}
