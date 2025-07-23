package patcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparacionCuatrimestres(t *testing.T) {
	assert.True(t, Cuatri{1, 2025}.despuesDe(Cuatri{2, 2023}))
	assert.False(t, Cuatri{2, 2023}.despuesDe(Cuatri{1, 2025}))
	assert.False(t, Cuatri{1, 2025}.despuesDe(Cuatri{1, 2025}))
}

func TestUnificarSinOfertas(t *testing.T) {
	oMateriasSiu := unificarOfertasMateriasSiu([]*ofertaCarrera{})

	assert.Empty(t, oMateriasSiu)
}

func TestUnificarConOfertasDisjuntas(t *testing.T) {
	m0 := materiaSiu{Nombre: "Análisis Matemático II"}
	m1 := materiaSiu{Nombre: "Álgebra Lineal"}

	oCarreras := []*ofertaCarrera{
		{Cuatri: Cuatri{1, 2025}, materias: []materiaSiu{m0}},
		{Cuatri: Cuatri{1, 2025}, materias: []materiaSiu{m1}},
	}

	oMateriasSiu := unificarOfertasMateriasSiu(oCarreras)

	assert.Len(t, oMateriasSiu, 2)
	assert.Contains(t, oMateriasSiu, ofertaMateriaSiu{Cuatri: Cuatri{1, 2025}, materia: m0})
	assert.Contains(t, oMateriasSiu, ofertaMateriaSiu{Cuatri: Cuatri{1, 2025}, materia: m1})
}

func TestUnificarConOfertasNoDisjuntas(t *testing.T) {
	m := materiaSiu{Nombre: "Análisis Matemático II"}

	oCarreras := []*ofertaCarrera{
		{Cuatri: Cuatri{1, 2025}, materias: []materiaSiu{m}},
		{Cuatri: Cuatri{2, 2024}, materias: []materiaSiu{m}},
		{Cuatri: Cuatri{1, 2023}, materias: []materiaSiu{m}},
	}

	oMateriasSiu := unificarOfertasMateriasSiu(oCarreras)

	assert.Len(t, oMateriasSiu, 1)
	assert.Equal(t, oMateriasSiu[0].Cuatri, Cuatri{1, 2025})
}

func TestUnificarConOfertasIguales(t *testing.T) {
	m := materiaSiu{Nombre: "Análisis Matemático II"}

	oCarreras := []*ofertaCarrera{
		{Cuatri: Cuatri{1, 2025}, materias: []materiaSiu{m}},
		{Cuatri: Cuatri{1, 2025}, materias: []materiaSiu{m}},
	}

	oMaterias := unificarOfertasMateriasSiu(oCarreras)

	assert.Len(t, oMaterias, 1)
}

func TestVincularMateriasSiuConDb(t *testing.T) {
	m0 := materiaSiu{
		Codigo:   "CB001",
		Nombre:   "ANÁLISIS matematico ii",
		Catedras: []CatedraSiu{{Codigo: 1}},
	}

	m1 := materiaSiu{
		Codigo:   "CB002",
		Nombre:   "algebra lineal",
		Catedras: []CatedraSiu{{Codigo: 2}},
	}

	m2 := materiaSiu{
		Codigo:   "CB003",
		Nombre:   "probabilidad y estadistica",
		Catedras: []CatedraSiu{{Codigo: 3}, {Codigo: 4}},
	}

	mSiu := []ofertaMateriaSiu{
		{Cuatri: Cuatri{1, 2025}, materia: m0},
		{Cuatri: Cuatri{2, 2024}, materia: m1},
		{Cuatri: Cuatri{1, 2025}, materia: m2},
	}

	mDb := []materiaDb{
		{Codigo: "CODXX1", Nombre: "Análisis Matemático II"},
		{Codigo: "CODXX2", Nombre: "Álgebra Lineal"},
		{Codigo: "CODXX3", Nombre: "Probabilidad y Estadística"},
	}

	patches := vincularMateriasSiuConDb(mSiu, mDb)

	assert.Len(t, patches, 3)

	assert.Contains(t, patches, PatchGenerado{
		CodigoDb:  "CODXX1",
		CodigoSiu: "CB001",
		Nombre:    "Análisis Matemático II",
		Catedras:  m0.Catedras,
		Cuatri:    Cuatri{1, 2025},
	})

	assert.Contains(t, patches, PatchGenerado{
		CodigoDb:  "CODXX2",
		CodigoSiu: "CB002",
		Nombre:    "Álgebra Lineal",
		Catedras:  m1.Catedras,
		Cuatri:    Cuatri{2, 2024},
	})

	assert.Contains(t, patches, PatchGenerado{
		CodigoDb:  "CODXX3",
		CodigoSiu: "CB003",
		Nombre:    "Probabilidad y Estadística",
		Catedras:  m2.Catedras,
		Cuatri:    Cuatri{1, 2025},
	})
}
