package patch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFiltrarSinOfertas(t *testing.T) {
	patches := filtrarOfertasMaterias([]*Oferta{})

	assert.Empty(t, patches)
}

func TestFiltrarConOfertasDisjuntas(t *testing.T) {
	m0 := MateriaSiu{Nombre: "Análisis Matemático II"}
	m1 := MateriaSiu{Nombre: "Álgebra Lineal"}

	ofertas := []*Oferta{
		{Cuatri: Cuatri{1, 2025}, Materias: []MateriaSiu{m0}},
		{Cuatri: Cuatri{1, 2025}, Materias: []MateriaSiu{m1}},
	}

	patches := filtrarOfertasMaterias(ofertas)

	assert.Len(t, patches, 2)
	assert.Contains(
		t,
		patches,
		Patch{
			CodigoSiu: m0.Codigo,
			Nombre:    m0.Nombre,
			Catedras:  m0.Catedras,
			Cuatri:    Cuatri{1, 2025},
		},
	)
	assert.Contains(
		t,
		patches,
		Patch{
			CodigoSiu: m1.Codigo,
			Nombre:    m1.Nombre,
			Catedras:  m1.Catedras,
			Cuatri:    Cuatri{1, 2025},
		},
	)
}

func TestFiltrarConOfertasNoDisjuntas(t *testing.T) {
	m := MateriaSiu{Nombre: "Análisis Matemático II"}

	ofertas := []*Oferta{
		{Cuatri: Cuatri{1, 2025}, Materias: []MateriaSiu{m}},
		{Cuatri: Cuatri{2, 2024}, Materias: []MateriaSiu{m}},
		{Cuatri: Cuatri{1, 2023}, Materias: []MateriaSiu{m}},
	}

	patches := filtrarOfertasMaterias(ofertas)

	assert.Len(t, patches, 1)
	assert.Equal(t, patches[0].Cuatri, Cuatri{1, 2025})
}

func TestFiltrarConOfertasIguales(t *testing.T) {
	m := MateriaSiu{Nombre: "Análisis Matemático II"}

	ofertas := []*Oferta{
		{Cuatri: Cuatri{1, 2025}, Materias: []MateriaSiu{m}},
		{Cuatri: Cuatri{1, 2025}, Materias: []MateriaSiu{m}},
	}

	patches := filtrarOfertasMaterias(ofertas)

	assert.Len(t, patches, 1)
}

func TestFiltrarConOfertasConflictivas(t *testing.T) {
	m := MateriaSiu{Nombre: "Análisis Matemático II"}

	ofertas := []*Oferta{
		{Cuatri: Cuatri{1, 2025}, Materias: []MateriaSiu{m}},
		{Cuatri: Cuatri{1, 2025}, Materias: []MateriaSiu{m}},
	}

	// Una misma materia está presente en dos ofertas de dos carreras
	// diferentes pero del mismo cuatrimestre, y las cátedras de las ofertas no
	// son idénticas entre si, sino que solo se intersectan en la cátedra con
	// código 2.

	ofertas[0].Materias[0].Catedras = []CatedraSiu{{Codigo: 1}, {Codigo: 2}}
	ofertas[1].Materias[0].Catedras = []CatedraSiu{{Codigo: 2}, {Codigo: 3}}

	patches := filtrarOfertasMaterias(ofertas)

	assert.Len(t, patches[0].Catedras, 3)
	assert.Contains(t, patches[0].Catedras, CatedraSiu{Codigo: 1})
	assert.Contains(t, patches[0].Catedras, CatedraSiu{Codigo: 2})
	assert.Contains(t, patches[0].Catedras, CatedraSiu{Codigo: 3})
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
