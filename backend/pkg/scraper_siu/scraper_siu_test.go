package scraper_siu

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObtenerCuatris(t *testing.T) {
	contenidoSiu := `
Período lectivo: 2024 - 1er Cuatrimestre
Materias del primer cuatrimestre aquí...
Cátedras del primer cuatrimestre aquí...

Período lectivo: 2024 - 2do Cuatrimestre
Materias del segundo cuatrimestre aquí...
Cátedras del segundo cuatrimestre aquí...
	`

	cuatris := obtenerCuatris(contenidoSiu)

	assert.Len(t, cuatris, 2)

	assert.Equal(t, 2024, cuatris[0].anio)
	assert.Equal(t, 1, cuatris[0].numero)
	assert.Equal(
		t,
		`Materias del primer cuatrimestre aquí...
Cátedras del primer cuatrimestre aquí...`,
		strings.TrimSpace(cuatris[0].data),
	)

	assert.Equal(t, 2024, cuatris[1].anio)
	assert.Equal(t, 2, cuatris[1].numero)
	assert.Equal(
		t,
		`Materias del segundo cuatrimestre aquí...
Cátedras del segundo cuatrimestre aquí...`,
		strings.TrimSpace(cuatris[1].data),
	)
}

func TestNoSeConsideranLosCursoDeVerano(t *testing.T) {
	contenidoSiu := `
Período lectivo: 2025 - Curso de verano
Materias del curso de verano aquí...
Cátedras del curso de verano aquí...
	`

	cuatris := obtenerCuatris(contenidoSiu)

	assert.Empty(t, cuatris)
}

func TestObtenerMaterias(t *testing.T) {
	data := `
Actividad: INTRODUCCIÓN AL DESARROLLO DE SOFTWARE (TB022)
Cátedras de la materia TB022 aquí...

Actividad: MODELACIÓN NUMÉRICA(CB051)
Cátedras de la materia CB051 aquí...

Actividad: ORGANIZACIÓN DEL COMPUTADOR (TB023)
Cátedras de la materia TB023 aquí...
	`

	materias := obtenerMaterias(data)

	assert.Len(t, materias, 3)

	assert.Equal(t, "INTRODUCCIÓN AL DESARROLLO DE SOFTWARE", materias[0].Nombre)
	assert.Equal(t, "TB022", materias[0].Codigo)
	assert.Equal(
		t,
		`Cátedras de la materia TB022 aquí...`,
		strings.TrimSpace(materias[0].data),
	)

	assert.Equal(t, "MODELACIÓN NUMÉRICA", materias[1].Nombre)
	assert.Equal(t, "CB051", materias[1].Codigo)
	assert.Equal(
		t,
		`Cátedras de la materia CB051 aquí...`,
		strings.TrimSpace(materias[1].data),
	)

	assert.Equal(t, "ORGANIZACIÓN DEL COMPUTADOR", materias[2].Nombre)
	assert.Equal(t, "TB023", materias[2].Codigo)
	assert.Equal(
		t,
		`Cátedras de la materia TB023 aquí...`,
		strings.TrimSpace(materias[2].data),
	)
}

func TestObtenerCatedras(t *testing.T) {
	data := `
Comisión: CURSO: 05
Docentes: BUCHWALD MARTÍN EZEQUIEL (Profesor Adjunto), PODBEREZSKI VICTOR DANIEL (Profesor Adjunto), GENENDER PEÑA EZEQUIEL DAVID (Jefe Trabajos Practicos)
Horarios de la cátedra 05 aquí...

Comisión: CURSO: 07
Docentes: BUCHWALD MARTÍN EZEQUIEL (Profesor Adjunto), GENENDER PEÑA EZEQUIEL DAVID (Jefe Trabajos Practicos)
Horarios de la cátedra 07 aquí...
	`

	catedras := obtenerCatedras(data)

	assert.Len(t, catedras, 2)

	assert.Equal(t, 5, catedras[0].Codigo)
	assert.Contains(t, catedras[0].Docentes, Docente{"BUCHWALD MARTÍN EZEQUIEL", "Profesor Adjunto"})
	assert.Contains(t, catedras[0].Docentes, Docente{"PODBEREZSKI VICTOR DANIEL", "Profesor Adjunto"})
	assert.Contains(t, catedras[0].Docentes, Docente{"GENENDER PEÑA EZEQUIEL DAVID", "Jefe Trabajos Practicos"})

	assert.Equal(t, 7, catedras[1].Codigo)
	assert.Contains(t, catedras[0].Docentes, Docente{"BUCHWALD MARTÍN EZEQUIEL", "Profesor Adjunto"})
	assert.Contains(t, catedras[0].Docentes, Docente{"GENENDER PEÑA EZEQUIEL DAVID", "Jefe Trabajos Practicos"})
}

func TestFormatosCodigosCatedras(t *testing.T) {
	// Existen más de 10 formatos diferentes en los que el SIU muestra los
	// nombre de las cátedras. ¿Por qué?. No lo sé.

	data := `
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
	`

	catedras := obtenerCatedras(data)

	assert.Len(t, catedras, 11)
}

func TestSeAsignaUnCodigoALasCatedrasSinCodigo(t *testing.T) {
	// Uno de los formatos de cátedra es aquel en el que solo se detalle el
	// nombre de la misma, es decir, la cátedra no tiene ni código ni mucho
	// menos variante. Esto sucede para algunas cátedras únicas. En este caso
	// se le asigna el código 1.

	data := "Comisión: CURSO: Caram"
	catedras := obtenerCatedras(data)

	assert.Equal(t, catedras[0].Codigo, 1)
}

func TestFormatosNombresDocentes(t *testing.T) {
	// La string que se pasa es por cátedra, así que tenemos que hacer el test
	// 4 veces con argumentos diferentes para probar todos los formatos.

	data := "Docentes: RAMOS SILVIA ADRIANA (Profesor Adjunto)"
	docentes := obtenerDocentes(data)

	assert.Len(t, docentes, 1)
	assert.Equal(t, docentes[0].Nombre, "RAMOS SILVIA ADRIANA")
	assert.Equal(t, docentes[0].Rol, "Profesor Adjunto")

	data = "Docentes: BUCHWALD MARTÍN EZEQUIEL (Profesor Adjunto), PODBEREZSKI VICTOR DANIEL (Profesor Adjunto), GENENDER PEÑA EZEQUIEL DAVID (Jefe Trabajos Practicos)"
	docentes = obtenerDocentes(data)

	assert.Len(t, docentes, 3)

	assert.Contains(t, docentes, Docente{"BUCHWALD MARTÍN EZEQUIEL", "Profesor Adjunto"})
	assert.Contains(t, docentes, Docente{"PODBEREZSKI VICTOR DANIEL", "Profesor Adjunto"})
	assert.Contains(t, docentes, Docente{"GENENDER PEÑA EZEQUIEL DAVID", "Jefe Trabajos Practicos"})

	data = "Docentes: Sin docentes"
	docentes = obtenerDocentes(data)

	assert.Empty(t, docentes)
}

func TestNoSeConsideranLosDocentesPorDesignar(t *testing.T) {
	data := "Docentes: A DESIGNAR A DESIGNAR (Profesor Adjunto), BOGGI SILVINA (Profesor Adjunto), VENTURIELLO VERONICA LAURA (Ayudante 1ro)"
	docentes := obtenerDocentes(data)

	assert.Len(t, docentes, 2)
	assert.NotEqual(t, docentes[0], "A DESIGNAR A DESIGNAR")
	assert.NotEqual(t, docentes[1], "A DESIGNAR A DESIGNAR")
}

func TestSeUnificanLasCatedrasConVariantes(t *testing.T) {
	// Algunas cátedras tienen horarios diferentes para diferentes grupos, por
	// ejemplo, cuando yo cursé 'Fisica I', la mitad de los alumnos de mi
	// cátedra iban al laboratorio mientras los demás se quedaban en el aula
	// recibiendo la clase teórica. El SIU trata estos horarios como cátedras
	// diferentes con un mismo código, pero les agrega un sufijo alfabético
	// diferente a cada una. FIUBA Review los trata como una sola cátedra.

	data := `
Comisión: CURSO: 23A
Docentes: BUCHWALD MARTÍN EZEQUIEL (Profesor Adjunto)

Comisión: CURSO: 23B
Docentes: PODBEREZSKI VICTOR DANIEL (Profesor Adjunto)

Comisión: CURSO: 23C
Docentes: GENENDER PEÑA EZEQUIEL DAVID (Jefe Trabajos Practicos)
	`

	catedras := obtenerCatedras(data)

	assert.Len(t, catedras, 1)
	assert.Equal(t, 23, catedras[0].Codigo)

	assert.Contains(t, catedras[0].Docentes, Docente{"BUCHWALD MARTÍN EZEQUIEL", "Profesor Adjunto"})
	assert.Contains(t, catedras[0].Docentes, Docente{"PODBEREZSKI VICTOR DANIEL", "Profesor Adjunto"})
	assert.Contains(t, catedras[0].Docentes, Docente{"GENENDER PEÑA EZEQUIEL DAVID", "Jefe Trabajos Practicos"})
}

func TestNoSeConsideranLasCatedraParaCondicionales(t *testing.T) {
	data := "Comisión: CONDICIONALES"

	catedras := obtenerCatedras(data)

	assert.Len(t, catedras, 0)
}

func TestNoSeConsideraElTrabajoProfesional(t *testing.T) {
	data := `
Actividad: TRABAJO PROFESIONAL DE INGENIERÍA INFORMÁTICA (TA053)
Cátedras de la materia TA053 aquí...

Actividad: TRABAJO PROFESIONAL DE INGENIERÍA QUÍMICA (TA170)
Cátedras de la materia TA170 aquí...
	`

	materias := obtenerMaterias(data)

	assert.Empty(t, materias)
}
