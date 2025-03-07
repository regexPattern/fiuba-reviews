package main

import (
	"maps"
	"slices"
)

type ofertaComisiones struct {
	carrera  string
	cuatri   cuatri
	materias []materia
}

type cuatri struct {
	numero int
	anio   int
}

func (c cuatri) esDespuesDe(otro cuatri) bool {
	return (otro.anio < c.anio) ||
		((otro.anio == c.anio) && (otro.numero < c.numero))
}

type materia struct {
	Codigo   string    `json:"codigo"`
	Nombre   string    `json:"nombre"`
	Catedras []catedra `json:"catedras"`
}

type catedra struct {
	Codigo   int       `json:"codigo"`
	Docentes []docente `json:"docentes"`
}

type docente struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

func filtrarMateriasMasRecientes(ofertas []ofertaComisiones) []materia {
	max := 0
	for _, o := range ofertas {
		max += len(o.materias)
	}

	cuatris := make(map[string]cuatri, max)
	materias := make(map[string]materia, max)

	for _, o := range ofertas {
		for _, m := range o.materias {
			cuatriUltimoCambio, ok := cuatris[m.Nombre]

			if !ok || o.cuatri.esDespuesDe(cuatriUltimoCambio) {
				cuatris[m.Nombre] = o.cuatri
				materias[m.Nombre] = m
			}
		}
	}

	return slices.Collect(maps.Values(materias))
}
