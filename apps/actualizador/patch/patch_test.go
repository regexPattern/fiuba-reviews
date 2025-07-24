package patch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparacionCuatrimestres(t *testing.T) {
	assert.True(t, cuatri{1, 2025}.despuesDe(cuatri{2, 2023}))
	assert.False(t, cuatri{2, 2023}.despuesDe(cuatri{1, 2025}))
	assert.False(t, cuatri{1, 2025}.despuesDe(cuatri{1, 2025}))
}

func TestFiltrarSinOfertas(t *testing.T) {
	om := filtrarOfertasMaterias([]*ofertaCarrera{})

	assert.Empty(t, om)
}

func TestFiltrarConOfertasDisjuntas(t *testing.T) {
	m0 := materiaSiu{Nombre: "Análisis Matemático II"}
	m1 := materiaSiu{Nombre: "Álgebra Lineal"}

	oc := []*ofertaCarrera{
		{cuatri: cuatri{1, 2025}, Materias: []materiaSiu{m0}},
		{cuatri: cuatri{1, 2025}, Materias: []materiaSiu{m1}},
	}

	om := filtrarOfertasMaterias(oc)

	assert.Len(t, om, 2)
	assert.Contains(
		t,
		om,
		Patch{Codigo: m0.Codigo, Nombre: m0.Nombre, Catedras: m0.Catedras, cuatri: cuatri{1, 2025}},
	)
	assert.Contains(
		t,
		om,
		Patch{Codigo: m1.Codigo, Nombre: m1.Nombre, Catedras: m1.Catedras, cuatri: cuatri{1, 2025}},
	)
}

func TestFiltrarConOfertasNoDisjuntas(t *testing.T) {
	m := materiaSiu{Nombre: "Análisis Matemático II"}

	oc := []*ofertaCarrera{
		{cuatri: cuatri{1, 2025}, Materias: []materiaSiu{m}},
		{cuatri: cuatri{2, 2024}, Materias: []materiaSiu{m}},
		{cuatri: cuatri{1, 2023}, Materias: []materiaSiu{m}},
	}

	om := filtrarOfertasMaterias(oc)

	assert.Len(t, om, 1)
	assert.Equal(t, om[0].cuatri, cuatri{1, 2025})
}

func TestFiltrarConOfertasIguales(t *testing.T) {
	m := materiaSiu{Nombre: "Análisis Matemático II"}

	oc := []*ofertaCarrera{
		{cuatri: cuatri{1, 2025}, Materias: []materiaSiu{m}},
		{cuatri: cuatri{1, 2025}, Materias: []materiaSiu{m}},
	}

	om := filtrarOfertasMaterias(oc)

	assert.Len(t, om, 1)
}

func TestFiltrarConOfertasConflictivas(t *testing.T) {
	m := materiaSiu{Nombre: "Análisis Matemático II"}

	oc := []*ofertaCarrera{
		{cuatri: cuatri{1, 2025}, Materias: []materiaSiu{m}},
		{cuatri: cuatri{1, 2025}, Materias: []materiaSiu{m}},
	}

	// Una misma materia está presente en dos ofertas de dos carreras
	// diferentes pero del mismo cuatrimestre, y las cátedras de las ofertas no
	// son idénticas entre si, sino que solo se intersectan en la cátedra con
	// código 2.

	oc[0].Materias[0].Catedras = []catedraSiu{{Codigo: 1}, {Codigo: 2}}
	oc[1].Materias[0].Catedras = []catedraSiu{{Codigo: 2}, {Codigo: 3}}

	om := filtrarOfertasMaterias(oc)

	assert.Len(t, om[0].Catedras, 3)
	assert.Contains(t, om[0].Catedras, catedraSiu{Codigo: 1})
	assert.Contains(t, om[0].Catedras, catedraSiu{Codigo: 2})
	assert.Contains(t, om[0].Catedras, catedraSiu{Codigo: 3})
}
