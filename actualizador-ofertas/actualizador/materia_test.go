package actualizador

import (
	"slices"
	"testing"
)

func TestSeFiltranLasMateriasDeLasOfetasMasRecientes(t *testing.T) {
	carrera := "Ingeniería en Informática"

	initDummyCatedras := func(codigos ...int) []catedra {
		catedras := make([]catedra, len(codigos))
		for i, c := range codigos {
			catedras[i].Codigo = c
		}
		return catedras
	}

	initDummyMateria := func(codigos ...int) materia {
		return materia{
			Codigo:   "AM2",
			Nombre:   "Análisis Matemático II",
			Catedras: initDummyCatedras(codigos...),
		}
	}

	codsCatedrasEsperadas := []int{3, 5, 8}

	o1 := oferta{
		ofertaMetadata: ofertaMetadata{
			carrera: carrera,
			cuatri:  cuatri{numero: 1, anio: 2024},
		},
		materias: []materia{initDummyMateria(4, 7, 11)},
	}
	masReciente := oferta{
		ofertaMetadata: ofertaMetadata{
			carrera: carrera,
			cuatri:  cuatri{numero: 1, anio: 2025},
		},
		materias: []materia{initDummyMateria(codsCatedrasEsperadas...)},
	}
	o3 := oferta{
		ofertaMetadata: ofertaMetadata{
			carrera: carrera,
			cuatri:  cuatri{numero: 2, anio: 2023},
		},
		materias: []materia{initDummyMateria(1, 2, 6)},
	}

	uofs := filtrarUltimasOfertas([]oferta{o1, masReciente, o3})

	mats := make([]materia, 0, len(uofs))
	for _, uc := range uofs {
		mats = append(mats, uc.materia)
	}

	if len(mats) != 1 {
		t.Fail()
	}

	codsCatedrasFiltradas := make([]int, len(codsCatedrasEsperadas))
	for i, c := range mats[0].Catedras {
		codsCatedrasFiltradas[i] = c.Codigo
	}

	if !slices.Equal(codsCatedrasEsperadas, codsCatedrasFiltradas) {
		t.Fail()
	}
}

func TestSeDistinguenDosMateriasComoIgualesPorSuNombre(t *testing.T) {
	o1 := oferta{
		ofertaMetadata: ofertaMetadata{
			cuatri: cuatri{numero: 1, anio: 2025},
		},
		materias: []materia{{
			Codigo:   "AM2",
			Nombre:   "Análisis Matemático II",
			Catedras: []catedra{{Codigo: 7}},
		}},
	}
	o2 := oferta{
		ofertaMetadata: ofertaMetadata{
			cuatri: cuatri{numero: 2, anio: 2021},
		},
		materias: []materia{{
			Codigo:   "AM2",
			Nombre:   "Física de los Sistemas de Partículas",
			Catedras: []catedra{{Codigo: 1}},
		}},
	}

	mats := filtrarUltimasOfertas([]oferta{o1, o2})
	if len(mats) != 2 {
		t.Fail()
	}
}

func TestSeConservanLasMateriasSinActualizacion(t *testing.T) {
	masReciente := oferta{
		ofertaMetadata: ofertaMetadata{
			carrera: "Ingeniería Civil",
			cuatri:  cuatri{numero: 1, anio: 2025},
		},
		materias: []materia{{
			Codigo:   "AM2",
			Nombre:   "Análisis Matemático II",
			Catedras: []catedra{{Codigo: 7}},
		}},
	}
	o2 := oferta{
		ofertaMetadata: ofertaMetadata{
			carrera: "Ingeniería en Informática",
			cuatri:  cuatri{numero: 2, anio: 2021},
		},
		materias: []materia{{
			Codigo:   "AM2",
			Nombre:   "Análisis Matemático II",
			Catedras: []catedra{{Codigo: 1}},
		}, {
			Codigo:   "FIS",
			Nombre:   "Física de los Sistemas de Partículas",
			Catedras: []catedra{{Codigo: 5}},
		}},
	}

	uofs := filtrarUltimasOfertas([]oferta{masReciente, o2})

	mats := make([]materia, 0, len(uofs))
	for _, uc := range uofs {
		mats = append(mats, uc.materia)
	}

	codsMateriasEsperadas := []string{"AM2", "FIS"}

	if len(mats) != 2 {
		t.Fail()
	}

	codsMateriasFiltradas := []string{mats[0].Codigo, mats[1].Codigo}
	slices.Sort(codsMateriasFiltradas)

	if !slices.Equal(codsMateriasFiltradas, codsMateriasEsperadas) {
		t.Fail()
	}
}

func TestSeFiltranLasCatedrasMasRecientesSinImportarLaCarrera(t *testing.T) {
	masReciente := oferta{
		ofertaMetadata: ofertaMetadata{
			carrera: "Ingeniería Civil",
			cuatri:  cuatri{numero: 1, anio: 2025},
		},
		materias: []materia{{
			Codigo:   "AM2",
			Nombre:   "Análisis Matemático II",
			Catedras: []catedra{{Codigo: 7}},
		}},
	}
	o2 := oferta{
		ofertaMetadata: ofertaMetadata{
			carrera: "Ingeniería en Informática",
			cuatri:  cuatri{numero: 2, anio: 2021},
		},
		materias: []materia{{
			Codigo:   "AM2",
			Nombre:   "Análisis Matemático II",
			Catedras: []catedra{{Codigo: 1}},
		}},
	}

	uofs := filtrarUltimasOfertas([]oferta{masReciente, o2})

	mats := make([]materia, 0, len(uofs))
	for _, uc := range uofs {
		mats = append(mats, uc.materia)
	}

	if mats[0].Catedras[0].Codigo != 7 {
		t.Fail()
	}
}
