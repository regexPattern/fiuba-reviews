package main

import (
	"slices"
	"testing"
)

func TestSeFiltranMateriasDeOfetasMasRecientes(t *testing.T) {
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
			Nombre:   "Análisis Matemático 2",
			Catedras: initDummyCatedras(codigos...),
		}
	}

	codsCatedrasMasRecientes := []int{3, 5, 8}

	p1 := plan{
		carrera:  carrera,
		cuatri:   cuatri{numero: 1, anio: 2024},
		materias: []materia{initDummyMateria(4, 7, 11)},
	}

	p2 := plan{ // plan más reciente
		carrera:  carrera,
		cuatri:   cuatri{numero: 1, anio: 2025},
		materias: []materia{initDummyMateria(codsCatedrasMasRecientes...)},
	}

	p3 := plan{
		carrera:  carrera,
		cuatri:   cuatri{numero: 2, anio: 2023},
		materias: []materia{initDummyMateria(1, 2, 6)},
	}

	materias := filtrarMateriasMasRecientes([]plan{p1, p2, p3})

	if len(materias) != 1 {
		t.Fatalf("Se agregó una mimas materias más de una vez.")
	}

	codsCatedrasFiltradas := make([]int, len(codsCatedrasMasRecientes))
	for i, c := range materias[0].Catedras {
		codsCatedrasFiltradas[i] = c.Codigo
	}

	if !slices.Equal(codsCatedrasMasRecientes, codsCatedrasFiltradas) {
		t.Fatalf("Se filtraron las cátedras de una oferta que no es la más reciente. Esperados: %v. Obtenidos: %v.",
			codsCatedrasFiltradas, codsCatedrasFiltradas)
	}
}

func TestSeDistinguenDosMateriasComoIgualPorNombre(t *testing.T) {
}
