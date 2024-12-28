package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeObtienenLosCuatrisCorrectamente(t *testing.T) {
	assert := assert.New(t)

	cuatris := obtenerCuatris(`
Período lectivo: 2024 - 1er Cuatrimestre
Período lectivo: 2024 - 2do Cuatrimestre
	`)

	assert.Len(cuatris, 2)

	assert.Equal(2024, cuatris[0].anio)
	assert.Equal(1, cuatris[0].numero)

	assert.Equal(2024, cuatris[1].anio)
	assert.Equal(2, cuatris[1].numero)
}

func TestSeIgnorarLosCuatrisConAnioInvalido(t *testing.T) {
	cuatris := obtenerCuatris(`
Período lectivo: @$^! - 1er Cuatrimestre
	`)

	assert.Empty(t, cuatris)
}

func TestSeIgnoranLosCursosDeVerano(t *testing.T) {
	cuatris := obtenerCuatris(`
Período lectivo: 2024 - Curso de Verano 2024/2025
	`)

	assert.Empty(t, cuatris)
}

func TestSeObtieneElContenidoDeLosCuatrisCorrectamente(t *testing.T) {
	assert := assert.New(t)

	cuatris := obtenerCuatris(`
Período lectivo: 2024 - 1er Cuatrimestre
CONTENIDO 1ER CUATRIMESTRE 2024
MÁS CONTENIDO 1ER CUATRIMESTRE 2024

Período lectivo: 2024 - 2do Cuatrimestre
CONTENIDO 2DO CUATRIMESTRE 2024
MÁS CONTENIDO 2DO CUATRIMESTRE 2024
	`)

	assert.Equal(
		`CONTENIDO 1ER CUATRIMESTRE 2024
MÁS CONTENIDO 1ER CUATRIMESTRE 2024`,
		strings.TrimSpace(cuatris[0].contenido),
	)
	assert.Equal(
		`CONTENIDO 2DO CUATRIMESTRE 2024
MÁS CONTENIDO 2DO CUATRIMESTRE 2024`,
		strings.TrimSpace(cuatris[1].contenido),
	)
}

func TestSeObtienenLasMateriasDeUnCuatriCorrectamente(t *testing.T) {
	assert := assert.New(t)

	materias := obtenerMateriasDeCuatri(`
Actividad: ÁLGEBRA LINEAL (CB002)
Actividad: ALGORITMOS Y ESTRUCTURAS DE DATOS (CB100)
Actividad: ANÁLISIS MATEMÁTICO II (CB001)
`)

	assert.Len(materias, 3)

	assert.Equal("ÁLGEBRA LINEAL", materias[0].Nombre)
	assert.Equal("CB002", materias[0].Codigo)

	assert.Equal("ALGORITMOS Y ESTRUCTURAS DE DATOS", materias[1].Nombre)
	assert.Equal("CB100", materias[1].Codigo)

	assert.Equal("ANÁLISIS MATEMÁTICO II", materias[2].Nombre)
	assert.Equal("CB001", materias[2].Codigo)
}

func TestSeIgnoranLasMateriasDeTrabajoProfesional(t *testing.T) {
	materias := obtenerMateriasDeCuatri(`
Actividad: TRABAJO PROFESIONAL DE INGENIERÍA INFORMÁTICA (TA053)
Actividad: TRABAJO PROFESIONAL DE INGENIERÍA QUÍMICA (TA170)
`)

	assert.Empty(t, materias)
}

func TestSeObtieneLasCatedrasDeUnaMateriaCorrectamente(t *testing.T) {
	// Por alguna razón el SIU tiene una amplia variedad de formatos para
	// nombrar cátedras.

	assert := assert.New(t)

	catedras := obtenerCatedrasDeMateria(`
Comisión: CURSO: 1
Comisión: CURSO: 02
Comisión: 03
Comisión: 04-Fontdevila
Comisión: CURSO: 05-Benitez
Comisión: CURSO: 06-Alvarez Hamelin
Comisión: CURSO: 07-Mendez/Pandolfo
Comisión: CURSO: 08- Ramos
Comisión: CURSO: 9-
Comisión: CURSO:10
Comisión: CURSO:11A
Comisión: CURSO 12
Comisión: Curso: 13
Comisión: Curso 14
		`)

	assert.Len(catedras, 14)

	codigosExpected := make([]int, 0, len(catedras))
	codigosActual := make([]int, 0, len(catedras))

	for i, cat := range catedras {
		codigosExpected = append(codigosExpected, i+1)
		codigosActual = append(codigosActual, cat.Codigo)
	}

	assert.ElementsMatch(codigosExpected, codigosActual)
}

func TestSeAsignaElCodigo1ALasCatedrasUnicasSinCodigo(t *testing.T) {
	assert := assert.New(t)

	catedras := obtenerCatedrasDeMateria(`
Comisión: CURSO: Caram
		`)

	assert.Len(catedras, 1)
	assert.Equal(catedras[0].Codigo, 1)
}

func TestSeIgnoranLasCatedrasSinCodigoONombre(t *testing.T) {
	catedras := obtenerCatedrasDeMateria(`
Comisión: CURSO:
Comisión:
		`)

	assert.Empty(t, catedras)
}

func TestSeIgnoranLasCatedrasParaCondicionales(t *testing.T) {
	catedras := obtenerCatedrasDeMateria(`
Comisión: CONDICIONALES
		`)

	assert.Empty(t, catedras)
}

func TestSeObtienenLosNombresDeLosDocentesCorrectamente(t *testing.T) {
	assert := assert.New(t)

	docentes := obtenerDocentesDeVariante(`
Docentes: BUCHWALD MARTÍN EZEQUIEL (Profesor/a Adjunto/a), PODBEREZSKI VICTOR DANIEL (Profesor/a Adjunto/a), GENENDER PEÑA EZEQUIEL DAVID (Jefe/a Trabajos Practicos)
		`)

	assert.Len(docentes, 3)

	assert.Contains(docentes, docente{"BUCHWALD MARTÍN EZEQUIEL", "PROFESOR ADJUNTO"})
	assert.Contains(docentes, docente{"PODBEREZSKI VICTOR DANIEL", "PROFESOR ADJUNTO"})
	assert.Contains(docentes, docente{"GENENDER PEÑA EZEQUIEL DAVID", "JEFE TRABAJOS PRACTICOS"})
}

func TestSeRetornaNilCuandoSeEncuentraUnaCatedraSinDocentes(t *testing.T) {
	docentes := obtenerDocentesDeVariante(`Docentes: Sin docentes`)

	assert.Nil(t, docentes)
}

func TestSeIgnoranLasCatedrasSinDocentes(t *testing.T) {
	assert := assert.New(t)

	catedras := obtenerCatedrasConDocentes(`
Comisión: CURSO: 1
Docentes: BUCHWALD MARTÍN EZEQUIEL (Profesor/a Adjunto/a), PODBEREZSKI VICTOR DANIEL (Profesor/a Adjunto/a), GENENDER PEÑA EZEQUIEL DAVID (Jefe/a Trabajos Practicos)

Comisión: CURSO: 2
Docentes: Sin docentes
		`)

	assert.Len(catedras, 1)
	assert.Equal(1, catedras[0].Codigo)
}

func TestSeAceptaCasingVariadoParaLosNombresDeLasMaterias(t *testing.T) {
	assert := assert.New(t)

	materias := obtenerMateriasDeCuatri(`
Actividad: álgebra lineal (CB002)
		`)

	assert.Len(materias, 1)
	assert.Equal("ÁLGEBRA LINEAL", materias[0].Nombre)
}

func TestSeAceptaCasingVariadoParaLosNombresDeLosDocentes(t *testing.T) {
	assert := assert.New(t)

	docentes := obtenerDocentesDeVariante(`
Docentes: BUCHWALD martín ezequiel (Profesor/a Adjunto/a)
		`)

	assert.Len(docentes, 1)
	assert.Contains(docentes, docente{"BUCHWALD MARTÍN EZEQUIEL", "PROFESOR ADJUNTO"})
}

func TestSeAgregaCadaDocenteUnaSolaVezAlUnificarVariantesDeCatedras(t *testing.T) {
	assert := assert.New(t)

	catedras := obtenerCatedrasDeMateria(`
Comisión: CURSO: 02A
Docentes: SARRIS CLAUDIA MONICA (Profesor/a Adjunto/a), FAGES LUCIANO RODOLFO (Ayudante 1ro/a)

Comisión: CURSO: 02B
Docentes: SARRIS CLAUDIA MONICA (Profesor/a Adjunto/a), GOMEZ CIAPPONI LAUTARO (Ayudante 1ro/a)
		`)

	assert.Len(catedras[0].Docentes, 3)

	assert.Contains(catedras[0].Docentes, docente{"SARRIS CLAUDIA MONICA", "PROFESOR ADJUNTO"})
	assert.Contains(catedras[0].Docentes, docente{"FAGES LUCIANO RODOLFO", "AYUDANTE 1RO"})
	assert.Contains(catedras[0].Docentes, docente{"GOMEZ CIAPPONI LAUTARO", "AYUDANTE 1RO"})
}
