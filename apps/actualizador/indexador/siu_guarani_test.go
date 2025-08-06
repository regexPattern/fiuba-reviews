package indexador

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparacionCuatrimestres(t *testing.T) {
	assert.True(t, Cuatri{1, 2025}.despuesDe(Cuatri{2, 2023}))
	assert.False(t, Cuatri{2, 2023}.despuesDe(Cuatri{1, 2025}))
	assert.False(t, Cuatri{1, 2025}.despuesDe(Cuatri{1, 2025}))
}

func TestFiltrarSinOfertas(t *testing.T) {
	materias := unificarOfertasSiu([]OfertaCarreraSiu{})

	assert.Empty(t, materias)
}

func TestFiltrarConOfertasDisjuntas(t *testing.T) {
	m0 := MateriaSiu{Nombre: "Análisis Matemático II"}
	m1 := MateriaSiu{Nombre: "Álgebra Lineal"}

	ofertasCarreras := []OfertaCarreraSiu{
		{Cuatri: Cuatri{1, 2025}, Materias: []MateriaSiu{m0}},
		{Cuatri: Cuatri{1, 2025}, Materias: []MateriaSiu{m1}},
	}

	materias := unificarOfertasSiu(ofertasCarreras)

	assert.Len(t, materias, 2)
	assert.Contains(
		t,
		materias,
		OfertaMateriaSiu{
			MateriaSiu: m0,
			Cuatri:     Cuatri{1, 2025},
		},
	)
	assert.Contains(
		t,
		materias,
		OfertaMateriaSiu{
			MateriaSiu: m1,
			Cuatri:     Cuatri{1, 2025},
		},
	)
}

func TestFiltrarConOfertasNoDisjuntas(t *testing.T) {
	m := MateriaSiu{Nombre: "Análisis Matemático II"}

	ofertasCarreras := []OfertaCarreraSiu{
		{Cuatri: Cuatri{1, 2025}, Materias: []MateriaSiu{m}},
		{Cuatri: Cuatri{2, 2024}, Materias: []MateriaSiu{m}},
		{Cuatri: Cuatri{1, 2023}, Materias: []MateriaSiu{m}},
	}

	materias := unificarOfertasSiu(ofertasCarreras)

	assert.Len(t, materias, 1)
	assert.Equal(t, materias[0].Cuatri, Cuatri{1, 2025})
}

func TestFiltrarConOfertasIguales(t *testing.T) {
	m := MateriaSiu{Nombre: "Análisis Matemático II"}

	ofertasCarreras := []OfertaCarreraSiu{
		{Cuatri: Cuatri{1, 2025}, Materias: []MateriaSiu{m}},
		{Cuatri: Cuatri{1, 2025}, Materias: []MateriaSiu{m}},
	}

	materias := unificarOfertasSiu(ofertasCarreras)

	assert.Len(t, materias, 1)
}

func TestFiltrarConOfertasConflictivas(t *testing.T) {
	m := MateriaSiu{Nombre: "Análisis Matemático II"}

	ofertasCarreras := []OfertaCarreraSiu{
		{Cuatri: Cuatri{1, 2025}, Materias: []MateriaSiu{m}},
		{Cuatri: Cuatri{1, 2025}, Materias: []MateriaSiu{m}},
	}

	// Una misma materia está presente en dos ofertas de dos carreras
	// diferentes pero del mismo cuatrimestre, y las cátedras de las ofertas no
	// son idénticas entre si, sino que solo se intersectan en la cátedra con
	// código 2.

	ofertasCarreras[0].Materias[0].Catedras = []CatedraSiu{{Codigo: 1}, {Codigo: 2}}
	ofertasCarreras[1].Materias[0].Catedras = []CatedraSiu{{Codigo: 2}, {Codigo: 3}}

	materias := unificarOfertasSiu(ofertasCarreras)

	assert.Len(t, materias[0].Catedras, 3)
	assert.Contains(t, materias[0].Catedras, CatedraSiu{Codigo: 1})
	assert.Contains(t, materias[0].Catedras, CatedraSiu{Codigo: 2})
	assert.Contains(t, materias[0].Catedras, CatedraSiu{Codigo: 3})
}
