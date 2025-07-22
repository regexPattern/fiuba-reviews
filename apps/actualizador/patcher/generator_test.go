package patcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparacionCuatrimestres(t *testing.T) {
	assert.Equal(t, 1, Cuatri{1, 2025}.Compare(Cuatri{2, 2024}))
	assert.Equal(t, -1, Cuatri{1, 2024}.Compare(Cuatri{2, 2024}))
	assert.Equal(t, 0, Cuatri{1, 2025}.Compare(Cuatri{1, 2025}))
}

func TestMergeSinOfertas(t *testing.T) {
	materias := mergeLatestOfertas([]*Oferta{})

	assert.Empty(t, materias)
}

func TestMergeConOfertasDisjuntas(t *testing.T) {
	m0 := Materia{Nombre: "Análisis Matemático II"}
	m1 := Materia{Nombre: "Álgebra Lineal"}

	ofertas := []*Oferta{
		{Materias: []Materia{m0}, Cuatri: Cuatri{1, 2025}},
		{Materias: []Materia{m1}, Cuatri: Cuatri{1, 2025}},
	}

	merged := mergeLatestOfertas(ofertas)

	assert.Len(t, merged, 2)
	assert.Contains(t, merged, LatestMateria{Materia: m0, Cuatri: Cuatri{1, 2025}})
	assert.Contains(t, merged, LatestMateria{Materia: m1, Cuatri: Cuatri{1, 2025}})
}

func TestMergeConOfertasNoDisjuntas(t *testing.T) {
	m := Materia{Nombre: "Análisis Matemático II"}

	ofertas := []*Oferta{
		{Materias: []Materia{m}, Cuatri: Cuatri{1, 2025}},
		{Materias: []Materia{m}, Cuatri: Cuatri{2, 2024}},
		{Materias: []Materia{m}, Cuatri: Cuatri{1, 2023}},
	}

	merged := mergeLatestOfertas(ofertas)

	assert.Len(t, merged, 1)
	assert.Equal(t, merged[0].Cuatri, Cuatri{1, 2025})
}

func TestMergeConOfertasIguales(t *testing.T) {
	m := Materia{Nombre: "Análisis Matemático II"}

	ofertas := []*Oferta{
		{Materias: []Materia{m}, Cuatri: Cuatri{1, 2025}},
		{Materias: []Materia{m}, Cuatri: Cuatri{1, 2025}},
	}

	merged := mergeLatestOfertas(ofertas)

	assert.Len(t, merged, 1)
}
