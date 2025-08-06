package resolver

import (
	"testing"

	"github.com/regexPattern/fiuba-reviews/apps/actualizador/indexador"
	"github.com/stretchr/testify/assert"
)

func TestPriorizacionDePatches(t *testing.T) {
	// Materia con 3 docentes.
	m0 := indexador.OfertaMateriaSiu{
		MateriaSiu: indexador.MateriaSiu{
			Codigo: "CB001",
			Nombre: "Análisis Matemático II",
			Catedras: []indexador.CatedraSiu{
				{Docentes: []indexador.DocenteSiu{
					{Nombre: "Seminara"},
					{Nombre: "Maulhardt"},
					{Nombre: "Acero"},
				}},
			},
		},
	}

	// Materia con 2 docentes.
	m1 := indexador.OfertaMateriaSiu{
		MateriaSiu: indexador.MateriaSiu{
			Codigo: "CB002",
			Nombre: "Álgebra Lineal",
			Catedras: []indexador.CatedraSiu{
				{Docentes: []indexador.DocenteSiu{
					{Nombre: "Vargas"},
					{Nombre: "López"},
				}},
			},
		},
	}

	// Materia con 1 docente.
	m2 := indexador.OfertaMateriaSiu{
		MateriaSiu: indexador.MateriaSiu{
			Codigo: "CB003",
			Nombre: "Probabilidad y Estadística",
			Catedras: []indexador.CatedraSiu{
				{Docentes: []indexador.DocenteSiu{
					{Nombre: "García"},
				}},
			},
		},
	}

	materias := []indexador.OfertaMateriaSiu{m0, m1, m2}

	sortPatchesSegunPrioridad(materias)

	// La materia con mayor cantidad de docentes queda primero en el orden de prioridades.
	assert.Equal(t, materias[0], m0) // m0: 3 docentes
	assert.Equal(t, materias[1], m1) // m1: 2 docentes
	assert.Equal(t, materias[2], m2) // m2: 1 docente
}
