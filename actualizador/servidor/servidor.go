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
	http.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		slog.Info("get_patches_pendientes", "method", "GET", "path", "/")
		handleGetPatchesPendientes(w, patches)
	})
	http.HandleFunc("GET /{codigoMateria}", func(w http.ResponseWriter, r *http.Request) {
		slog.Info(
			"get_patch_materia",
			"method",
			"GET",
			"path",
			"/{codigoMateria}",
			"codigo_materia",
			r.PathValue("codigoMateria"),
		)
		handleGetPatchMateria(w, r, conn, patches)
	})
	http.HandleFunc("PATCH /{codigoMateria}", func(w http.ResponseWriter, r *http.Request) {
		slog.Info(
			"patch_resolver_materia",
			"method",
			"PATCH",
			"path",
			"/{codigoMateria}",
			"codigo_materia",
			r.PathValue("codigoMateria"),
		)
		handleResolverMateria(w, r, conn, patches)
	})

	slog.Info("servidor_iniciado", "addr", addr)

	return http.ListenAndServe(addr, nil)
}

func handleGetPatchesPendientes(w http.ResponseWriter, patches map[string]patchMateria) {
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
		slog.Error("encode_patches_failed", "error", err)
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
				"patch de actualización para materia %v no encontrado",
				codigoMateria,
			),
			http.StatusNotFound,
		)
		return
	}

	docentesPorCatedra, err := getDocentesConEstadoPorCatedra(conn, codigoMateria, patch.Catedras)
	if err != nil {
		slog.Error("get_docentes_estado_failed", "codigo_materia", codigoMateria, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type docenteCatedraRes struct {
		Nombre string  `json:"nombre"`
		Codigo *string `json:"codigo"`
	}

	catedras := make([][]docenteCatedraRes, 0, len(patch.Catedras))
	for _, cat := range patch.Catedras {
		docentesCatedra := make([]docenteCatedraRes, 0, len(cat.Docentes))

		for _, doc := range cat.Docentes {
			docentesCatedra = append(docentesCatedra, docenteCatedraRes{
				Nombre: doc.Nombre,
				Codigo: docentesPorCatedra[cat.Codigo][doc.Nombre],
			})
		}

		slices.SortFunc(docentesCatedra, func(a, b docenteCatedraRes) int {
			return strings.Compare(a.Nombre, b.Nombre)
		})

		catedras = append(catedras, docentesCatedra)
	}

	type patchMateriaRes struct {
		materia
		cuatrimestre       `                      json:"cuatrimestre"`
		DocentesPendientes []patchDocente        `json:"docentes_pendientes"`
		DocentesPorCatedra [][]docenteCatedraRes `json:"docentes_por_catedra"`
	}

	res := patchMateriaRes{
		materia:            patch.materia,
		cuatrimestre:       patch.cuatrimestre,
		DocentesPendientes: patch.Docentes,
		DocentesPorCatedra: catedras,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		slog.Error(
			"encode_patch_materia_failed",
			"codigo_materia",
			patch.Codigo,
			"nombre",
			patch.Nombre,
			"error",
			err,
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleResolverMateria(
	w http.ResponseWriter,
	r *http.Request,
	conn *pgx.Conn,
	patches map[string]patchMateria,
) {
	codigoMateria := r.PathValue("codigoMateria")
	patch, ok := patches[codigoMateria]
	if !ok {
		slog.Warn("patch_not_found", "codigo_materia", codigoMateria)
		http.Error(
			w,
			fmt.Sprintf("materia %v no tiene actualización disponible", codigoMateria),
			http.StatusNotFound,
		)
		return
	}

	var res resolucionMateria
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		slog.Error("decode_resolucion_failed", "codigo_materia", codigoMateria, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := resolverMateria(conn, patch, res); err != nil {
		slog.Error("resolver_materia_failed", "codigo_materia", codigoMateria, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	delete(patches, codigoMateria)

	w.WriteHeader(http.StatusNoContent)
}
