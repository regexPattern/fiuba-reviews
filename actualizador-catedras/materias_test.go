package main

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

	p1 := oferta{
		carrera:  carrera,
		cuatri:   cuatri{numero: 1, anio: 2024},
		materias: []materia{initDummyMateria(4, 7, 11)},
	}

	p2 := oferta{ // plan más reciente
		carrera:  carrera,
		cuatri:   cuatri{numero: 1, anio: 2025},
		materias: []materia{initDummyMateria(codsCatedrasEsperadas...)},
	}

	p3 := oferta{
		carrera:  carrera,
		cuatri:   cuatri{numero: 2, anio: 2023},
		materias: []materia{initDummyMateria(1, 2, 6)},
	}

	ultimasComisiones := FiltrarUltimasComisiones([]oferta{p1, p2, p3})

	materias := make([]materia, 0, len(ultimasComisiones))
	for _, uc := range ultimasComisiones {
		materias = append(materias, uc.materia)
	}

	if len(materias) != 1 {
		t.Fail()
	}

	codsCatedrasFiltradas := make([]int, len(codsCatedrasEsperadas))
	for i, c := range materias[0].Catedras {
		codsCatedrasFiltradas[i] = c.Codigo
	}

	if !slices.Equal(codsCatedrasEsperadas, codsCatedrasFiltradas) {
		t.Fail()
	}
}

func TestSeDistinguenDosMateriasComoIgualesPorSuNombre(t *testing.T) {
	p1 := oferta{
		cuatri: cuatri{numero: 1, anio: 2025},
		materias: []materia{{
			Codigo:   "AM2",
			Nombre:   "Análisis Matemático II",
			Catedras: []catedra{{Codigo: 7}},
		}},
	}

	p2 := oferta{
		cuatri: cuatri{numero: 2, anio: 2021},
		materias: []materia{{
			Codigo:   "AM2",
			Nombre:   "Física de los Sistemas de Partículas",
			Catedras: []catedra{{Codigo: 1}},
		}},
	}

	materias := FiltrarUltimasComisiones([]oferta{p1, p2})

	if len(materias) != 2 {
		t.Fail()
	}
}

func TestSeConservanLasMateriasSinActualizacion(t *testing.T) {
	p1 := oferta{ // plan más reciente
		carrera: "Ingeniería Civil",
		cuatri:  cuatri{numero: 1, anio: 2025},
		materias: []materia{{
			Codigo:   "AM2",
			Nombre:   "Análisis Matemático II",
			Catedras: []catedra{{Codigo: 7}},
		}},
	}

	p2 := oferta{
		carrera: "Ingeniería en Informática",
		cuatri:  cuatri{numero: 2, anio: 2021},
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

	ultimasComisiones := FiltrarUltimasComisiones([]oferta{p1, p2})

	materias := make([]materia, 0, len(ultimasComisiones))
	for _, uc := range ultimasComisiones {
		materias = append(materias, uc.materia)
	}

	codsMateriasEsperadas := []string{"AM2", "FIS"}

	if len(materias) != 2 {
		t.Fail()
	}

	codsMateriasFiltradas := []string{materias[0].Codigo, materias[1].Codigo}
	slices.Sort(codsMateriasFiltradas)

	if !slices.Equal(codsMateriasFiltradas, codsMateriasEsperadas) {
		t.Fail()
	}
}

func TestSeFiltranLasCatedrasMasRecientesSinImportarLaCarrera(t *testing.T) {
	p1 := oferta{ // plan más reciente
		carrera: "Ingeniería Civil",
		cuatri:  cuatri{numero: 1, anio: 2025},
		materias: []materia{{
			Codigo:   "AM2",
			Nombre:   "Análisis Matemático II",
			Catedras: []catedra{{Codigo: 7}},
		}},
	}

	p2 := oferta{
		carrera: "Ingeniería en Informática",
		cuatri:  cuatri{numero: 2, anio: 2021},
		materias: []materia{{
			Codigo:   "AM2",
			Nombre:   "Análisis Matemático II",
			Catedras: []catedra{{Codigo: 1}},
		}},
	}

	ultimasComisiones := FiltrarUltimasComisiones([]oferta{p1, p2})

	materias := make([]materia, 0, len(ultimasComisiones))
	for _, uc := range ultimasComisiones {
		materias = append(materias, uc.materia)
	}

	if materias[0].Catedras[0].Codigo != 7 {
		t.Fail()
	}
}
