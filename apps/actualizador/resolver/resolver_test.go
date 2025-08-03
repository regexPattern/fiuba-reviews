package resolver

import (
	"testing"

	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
	"github.com/stretchr/testify/assert"
)

func TestPriorizacionDePatches(t *testing.T) {
	// Materia con 3 docentes.
	m0 := patcher.MateriaSiu{
		Codigo: "CB001",
		Nombre: "Análisis Matemático II",
		Catedras: []patcher.CatedraSiu{
			{Docentes: []patcher.DocenteSiu{
				{Nombre: "Seminara"},
				{Nombre: "Maulhardt"},
				{Nombre: "Acero"},
			}},
		},
	}

	// Materia con 2 docentes.
	m1 := patcher.MateriaSiu{
		Codigo: "CB002",
		Nombre: "Álgebra Lineal",
		Catedras: []patcher.CatedraSiu{
			{Docentes: []patcher.DocenteSiu{
				{Nombre: "Vargas"},
				{Nombre: "López"},
			}},
		},
	}

	// Materia con 1 docente.
	m2 := patcher.MateriaSiu{
		Codigo: "CB003",
		Nombre: "Probabilidad y Estadística",
		Catedras: []patcher.CatedraSiu{
			{Docentes: []patcher.DocenteSiu{
				{Nombre: "García"},
			}},
		},
	}

	p0 := patcher.Patch{OfertaMateriaSiu: patcher.OfertaMateriaSiu{Materia: m0}}
	p1 := patcher.Patch{OfertaMateriaSiu: patcher.OfertaMateriaSiu{Materia: m1}}
	p2 := patcher.Patch{OfertaMateriaSiu: patcher.OfertaMateriaSiu{Materia: m2}}

	patches := []patcher.Patch{p0, p1, p2}

	priorizarPatches(patches)

	// La materia con mayor cantidad de docentes queda primero en el orden de prioridades.
	assert.Equal(t, patches[0], p0) // m0: 3 docentes
	assert.Equal(t, patches[1], p1) // m1: 2 docentes
	assert.Equal(t, patches[2], p2) // m2: 1 docente
}
