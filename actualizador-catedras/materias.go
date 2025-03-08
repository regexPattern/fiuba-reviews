package main

import (
	"maps"
	"slices"

	"github.com/charmbracelet/log"
)

type oferta struct {
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

type ultimaComision struct {
	materia materia
	cuatri  cuatri
}

// entre las ofertas del SIU disponible, nos quedamos con la mas reciente de cada materia.
// TODO: buscar mejores nombres definitivamente
func filtrarUltimasComisiones(ofertas []oferta) []ultimaComision {
	logger := log.Default().WithPrefix("🧹")

	max := 0
	for _, o := range ofertas {
		max += len(o.materias)
	}

	cuatris := make(map[string]cuatri, max)
	materias := make(map[string]ultimaComision, max)

	logger.Info("filtrando solo las ofertas de comisiones más recientes")

	for _, o := range ofertas {
		for _, m := range o.materias {
			cuatriUltimaActualizacion, ok := cuatris[m.Nombre]

			if !ok || o.cuatri.esDespuesDe(cuatriUltimaActualizacion) {
				cuatris[m.Nombre] = o.cuatri
				materias[m.Nombre] = ultimaComision{
					materia: m,
					cuatri:  o.cuatri,
				}
			}
		}
	}

	return slices.Collect(maps.Values(materias))
}
