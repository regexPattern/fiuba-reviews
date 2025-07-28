package actualizador

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
	patches := filtrarOfertasMaterias([]*oferta{})

	assert.Empty(t, patches)
}

func TestFiltrarConOfertasDisjuntas(t *testing.T) {
	m0 := materiaSiu{Nombre: "Análisis Matemático II"}
	m1 := materiaSiu{Nombre: "Álgebra Lineal"}

	ofertas := []*oferta{
		{cuatri: cuatri{1, 2025}, Materias: []materiaSiu{m0}},
		{cuatri: cuatri{1, 2025}, Materias: []materiaSiu{m1}},
	}

	patches := filtrarOfertasMaterias(ofertas)

	assert.Len(t, patches, 2)
	assert.Contains(
		t,
		patches,
		PatchActualizacionMateria{
			CodigoSiu: m0.Codigo,
			Nombre:    m0.Nombre,
			Catedras:  m0.Catedras,
			cuatri:    cuatri{1, 2025},
		},
	)
	assert.Contains(
		t,
		patches,
		PatchActualizacionMateria{
			CodigoSiu: m1.Codigo,
			Nombre:    m1.Nombre,
			Catedras:  m1.Catedras,
			cuatri:    cuatri{1, 2025},
		},
	)
}

func TestFiltrarConOfertasNoDisjuntas(t *testing.T) {
	m := materiaSiu{Nombre: "Análisis Matemático II"}

	ofertas := []*oferta{
		{cuatri: cuatri{1, 2025}, Materias: []materiaSiu{m}},
		{cuatri: cuatri{2, 2024}, Materias: []materiaSiu{m}},
		{cuatri: cuatri{1, 2023}, Materias: []materiaSiu{m}},
	}

	patches := filtrarOfertasMaterias(ofertas)

	assert.Len(t, patches, 1)
	assert.Equal(t, patches[0].cuatri, cuatri{1, 2025})
}

func TestFiltrarConOfertasIguales(t *testing.T) {
	m := materiaSiu{Nombre: "Análisis Matemático II"}

	ofertas := []*oferta{
		{cuatri: cuatri{1, 2025}, Materias: []materiaSiu{m}},
		{cuatri: cuatri{1, 2025}, Materias: []materiaSiu{m}},
	}

	patches := filtrarOfertasMaterias(ofertas)

	assert.Len(t, patches, 1)
}

func TestFiltrarConOfertasConflictivas(t *testing.T) {
	m := materiaSiu{Nombre: "Análisis Matemático II"}

	ofertas := []*oferta{
		{cuatri: cuatri{1, 2025}, Materias: []materiaSiu{m}},
		{cuatri: cuatri{1, 2025}, Materias: []materiaSiu{m}},
	}

	// Una misma materia está presente en dos ofertas de dos carreras
	// diferentes pero del mismo cuatrimestre, y las cátedras de las ofertas no
	// son idénticas entre si, sino que solo se intersectan en la cátedra con
	// código 2.

	ofertas[0].Materias[0].Catedras = []catedraSiu{{Codigo: 1}, {Codigo: 2}}
	ofertas[1].Materias[0].Catedras = []catedraSiu{{Codigo: 2}, {Codigo: 3}}

	patches := filtrarOfertasMaterias(ofertas)

	assert.Len(t, patches[0].Catedras, 3)
	assert.Contains(t, patches[0].Catedras, catedraSiu{Codigo: 1})
	assert.Contains(t, patches[0].Catedras, catedraSiu{Codigo: 2})
	assert.Contains(t, patches[0].Catedras, catedraSiu{Codigo: 3})
}

func TestMapNombreCodigoMateriasDb(t *testing.T) {
	materias := []materiaDb{
		{Codigo: "COD001", Nombre: "Análisis Matemático II"},
		{Codigo: "COD002", Nombre: "Álgebra Lineal"},
	}

	codigos := mapNombreCodigo(materias)

	assert.Len(t, codigos, 2)
	assert.Equal(t, codigos["analisis matematico ii"], "COD001")
	assert.Equal(t, codigos["algebra lineal"], "COD002")
}
