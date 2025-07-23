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

func TestMergeSinOfertas(t *testing.T) {
	oMaterias := mergeUltimasOfertasMateriasSiu([]*ofertaCarrera{})

	assert.Empty(t, oMaterias)
}

func TestMergeConOfertasDisjuntas(t *testing.T) {
	m0 := materiaSiu{Nombre: "Análisis Matemático II"}
	m1 := materiaSiu{Nombre: "Álgebra Lineal"}

	ofertas := []*ofertaCarrera{
		{cuatri: cuatri{1, 2025}, materias: []materiaSiu{m0}},
		{cuatri: cuatri{1, 2025}, materias: []materiaSiu{m1}},
	}

	merged := mergeUltimasOfertasMateriasSiu(ofertas)

	assert.Len(t, merged, 2)
	assert.Contains(t, merged, ofertaMateriaSiu{cuatri: cuatri{1, 2025}, materia: m0})
	assert.Contains(t, merged, ofertaMateriaSiu{cuatri: cuatri{1, 2025}, materia: m1})
}

func TestMergeConOfertasNoDisjuntas(t *testing.T) {
	m := materiaSiu{Nombre: "Análisis Matemático II"}

	ofertas := []*ofertaCarrera{
		{cuatri: cuatri{1, 2025}, materias: []materiaSiu{m}},
		{cuatri: cuatri{2, 2024}, materias: []materiaSiu{m}},
		{cuatri: cuatri{1, 2023}, materias: []materiaSiu{m}},
	}

	merged := mergeUltimasOfertasMateriasSiu(ofertas)

	assert.Len(t, merged, 1)
	assert.Equal(t, merged[0].cuatri, cuatri{1, 2025})
}

func TestMergeConOfertasIguales(t *testing.T) {
	m := materiaSiu{Nombre: "Análisis Matemático II"}

	ofertas := []*ofertaCarrera{
		{cuatri: cuatri{1, 2025}, materias: []materiaSiu{m}},
		{cuatri: cuatri{1, 2025}, materias: []materiaSiu{m}},
	}

	merged := mergeUltimasOfertasMateriasSiu(ofertas)

	assert.Len(t, merged, 1)
}
