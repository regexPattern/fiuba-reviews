package scraper_siu

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeObtieneLaCarreraCorrectamente(t *testing.T) {
	carrera, _ := obtenerCarrera(`
Propuesta: Ingeniería en Informática
		`)

	assert.Equal(t, "INGENIERÍA EN INFORMÁTICA", carrera)
}

func TestSeRetornaUnErrorCuandoNoHayCarrera(t *testing.T) {
	_, err := obtenerCarrera(`
Propuesta: 
		`)

	assert.EqualError(t, err, "No se encontró la carrera.")
}

func TestSeObtienenLosCuatrisCorrectamente(t *testing.T) {
	assert := assert.New(t)

	cuatris := obtenerCuatris(`
Período lectivo: 2024 - 1er Cuatrimestre
Período lectivo: 2024 - 2do Cuatrimestre
Período lectivo: 2023 - 1er Cuatrimestre
Período lectivo: 2023 - 2do Cuatrimestre
	`)

	require.Len(t, cuatris, 4)

	// Los cuatrimestres se retornan ordenados cronológicamente, con el más
	// reciente al final del listado.

	assert.Equal(2023, cuatris[0].Anio)
	assert.Equal(1, cuatris[0].Numero)

	assert.Equal(2023, cuatris[1].Anio)
	assert.Equal(2, cuatris[1].Numero)

	assert.Equal(2024, cuatris[2].Anio)
	assert.Equal(1, cuatris[2].Numero)

	assert.Equal(2024, cuatris[3].Anio)
	assert.Equal(2, cuatris[3].Numero)
}

func TestSeRetornaUnErrorCuandoNoHayCuatrimestres(t *testing.T) {
	_, err := ObtenerMetaData(`
Propuesta: Ingeniería en Informática
		`)

	assert.EqualError(t, err, "No se encontraron cuatrimestres.")
}

func TestSeIgnoranLosCuatrisConAnioInvalido(t *testing.T) {
	cuatris := obtenerCuatris(`
Período lectivo: @$^! - 1er Cuatrimestre
	`)

	require.Empty(t, cuatris)
}

func TestSeIgnoranLosCursosDeVerano(t *testing.T) {
	cuatris := obtenerCuatris(`
Período lectivo: 2024 - Curso de Verano 2024/2025
	`)

	require.Empty(t, cuatris)
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
		strings.TrimSpace(cuatris[0].Contenido),
	)
	assert.Equal(
		`CONTENIDO 2DO CUATRIMESTRE 2024
MÁS CONTENIDO 2DO CUATRIMESTRE 2024`,
		strings.TrimSpace(cuatris[1].Contenido),
	)
}

func TestSeObtienenLasMateriasDeUnCuatriCorrectamente(t *testing.T) {
	assert := assert.New(t)

	materias := ObtenerMaterias(`
Actividad: ÁLGEBRA LINEAL (CB002)
Actividad: ALGORITMOS Y ESTRUCTURAS DE DATOS (CB100)
Actividad: ANÁLISIS MATEMÁTICO II (CB001)
`)

	require.Len(t, materias, 3)

	assert.Equal("ÁLGEBRA LINEAL", materias[0].Nombre)
	assert.Equal("CB002", materias[0].Codigo)

	assert.Equal("ALGORITMOS Y ESTRUCTURAS DE DATOS", materias[1].Nombre)
	assert.Equal("CB100", materias[1].Codigo)

	assert.Equal("ANÁLISIS MATEMÁTICO II", materias[2].Nombre)
	assert.Equal("CB001", materias[2].Codigo)
}

func TestSeIgnoranLasMateriasDeTrabajoProfesional(t *testing.T) {
	materias := ObtenerMaterias(`
Actividad: TRABAJO PROFESIONAL DE INGENIERÍA INFORMÁTICA (TA053)
Actividad: TRABAJO PROFESIONAL DE INGENIERÍA QUÍMICA (TA170)
`)

	assert.Empty(t, materias)
}

func TestSeObtieneLasCatedrasDeUnaMateriaCorrectamente(t *testing.T) {
	assert := assert.New(t)

	// Por alguna razón el SIU tiene una amplia variedad de formatos para
	// nombrar cátedras.

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

func TestSeAsignanCodigosSecuencialesALasCatedrasSinCodigo(t *testing.T) {
	assert := assert.New(t)

	catedras := obtenerCatedrasDeMateria(`
Comisión: CURSO: 01
Comisión: CURSO: 07
Comisión: CURSO: Caram
Comisión: CURSO VIRTUAL: ESPECIAL PARA RECURSANTES
		`)

	require.Len(t, catedras, 4)

	codsCats := []int{0, 0, 0, 0}

	for i, cat := range catedras {
		codsCats[i] = cat.Codigo
	}

	assert.Contains(codsCats, 8)
	assert.Contains(codsCats, 9)
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
Comisión: CURSO: CONDICIONALES
		`)

	assert.Empty(t, catedras)
}

func TestSeObtienenLosNombresDeLosDocentesCorrectamente(t *testing.T) {
	assert := assert.New(t)

	docentes := obtenerDocentesDeVariante(`
Docentes: BUCHWALD MARTÍN EZEQUIEL (Profesor/a Adjunto/a), PODBEREZSKI VICTOR DANIEL (Profesor/a Adjunto/a), GENENDER PEÑA EZEQUIEL DAVID (Jefe/a Trabajos Practicos)
		`)

	require.Len(t, docentes, 3)

	assert.Contains(docentes, Docente{"BUCHWALD MARTÍN EZEQUIEL", "PROFESOR ADJUNTO"})
	assert.Contains(docentes, Docente{"PODBEREZSKI VICTOR DANIEL", "PROFESOR ADJUNTO"})
	assert.Contains(docentes, Docente{"GENENDER PEÑA EZEQUIEL DAVID", "JEFE TRABAJOS PRACTICOS"})
}

func TestSeRetornaNilCuandoSeEncuentraUnaCatedraSinDocentes(t *testing.T) {
	docentes := obtenerDocentesDeVariante(`Docentes: Sin docentes`)

	assert.Nil(t, docentes)
}

func TestSeIgnoranLasCatedrasSinDocentes(t *testing.T) {
	catedras := obtenerCatedrasConDocentes(`
Comisión: CURSO: 1
Docentes: BUCHWALD MARTÍN EZEQUIEL (Profesor/a Adjunto/a), PODBEREZSKI VICTOR DANIEL (Profesor/a Adjunto/a), GENENDER PEÑA EZEQUIEL DAVID (Jefe/a Trabajos Practicos)

Comisión: CURSO: 2
Docentes: Sin docentes
		`)

	require.Len(t, catedras, 1)
	assert.Equal(t, 1, catedras[0].Codigo)
}

func TestSeAceptaCasingVariadoParaLosNombresDeLasMaterias(t *testing.T) {
	materias := ObtenerMaterias(`
Actividad: álgebra lineal (CB002)
		`)

	require.Len(t, materias, 1)
	require.Equal(t, "ÁLGEBRA LINEAL", materias[0].Nombre)
}

func TestSeAceptaCasingVariadoParaLosNombresDeLosDocentes(t *testing.T) {
	docentes := obtenerDocentesDeVariante(`
Docentes: BUCHWALD martín ezequiel (Profesor/a Adjunto/a)
		`)

	require.Len(t, docentes, 1)
	assert.Contains(t, docentes, Docente{"BUCHWALD MARTÍN EZEQUIEL", "PROFESOR ADJUNTO"})
}

func TestSeAgregaCadaDocenteUnaSolaVezAlUnificarVariantesDeCatedras(t *testing.T) {
	assert := assert.New(t)

	catedras := obtenerCatedrasDeMateria(`
Comisión: CURSO: 02A
Docentes: SARRIS CLAUDIA MONICA (Profesor/a Adjunto/a), FAGES LUCIANO RODOLFO (Ayudante 1ro/a)

Comisión: CURSO: 02B
Docentes: SARRIS CLAUDIA MONICA (Profesor/a Adjunto/a), GOMEZ CIAPPONI LAUTARO (Ayudante 1ro/a)
		`)

	require.Len(t, catedras[0].Docentes, 3)

	assert.Contains(catedras[0].Docentes, Docente{"SARRIS CLAUDIA MONICA", "PROFESOR ADJUNTO"})
	assert.Contains(catedras[0].Docentes, Docente{"FAGES LUCIANO RODOLFO", "AYUDANTE 1RO"})
	assert.Contains(catedras[0].Docentes, Docente{"GOMEZ CIAPPONI LAUTARO", "AYUDANTE 1RO"})
}

func leerArchivoTestOfertaDeComisiones(filename string) string {
	_, testFilename, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(testFilename)
	fullFilepath := filepath.Join(testDir, "testdata", filename)
	contenidoSiu, _ := os.ReadFile(fullFilepath)
	return string(contenidoSiu)
}

func TestOfertaDeComisionesInformatica2C2024(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	contenidoSiu := leerArchivoTestOfertaDeComisiones("informatica-28-12-2024.txt")

	meta, _ := ObtenerMetaData(string(contenidoSiu))
	materias := ObtenerMaterias(string(contenidoSiu))

	require.Equal(meta.Carrera, "INGENIERÍA EN INFORMÁTICA")
	require.Len(materias, 30)

	matsCantCats := make(map[string]int, len(materias))

	for _, mat := range materias {
		matsCantCats[mat.Nombre] = matsCantCats[mat.Nombre] + len(mat.Catedras)
	}

	assert.Equal(14, matsCantCats["ÁLGEBRA LINEAL"])
	assert.Equal(5, matsCantCats["ALGORITMOS Y ESTRUCTURAS DE DATOS"])
	assert.Equal(16, matsCantCats["ANÁLISIS MATEMÁTICO II"])
	assert.Equal(3, matsCantCats["ANÁLISIS MATEMÁTICO III"])
	assert.Equal(1, matsCantCats["APRENDIZAJE AUTOMÁTICO"])
	assert.Equal(1, matsCantCats["APRENDIZAJE PROFUNDO"])
	assert.Equal(1, matsCantCats["ARQUITECTURA DE SOFTWARE"])
	assert.Equal(3, matsCantCats["BASE DE DATOS"])
	assert.Equal(2, matsCantCats["CIENCIA DE DATOS"])
	assert.Equal(1, matsCantCats["COMPUTACIÓN CUÁNTICA"])
	assert.Equal(1, matsCantCats["EMPRESAS DE BASE TECNOLÓGICA I"])
	assert.Equal(1, matsCantCats["EMPRESAS DE BASE TECNOLÓGICA II"])
	assert.Equal(1, matsCantCats["FÍSICA PARA INFORMÁTICA"])
	assert.Equal(4, matsCantCats["FUNDAMENTOS DE PROGRAMACIÓN"])
	assert.Equal(3, matsCantCats["GESTIÓN DEL DESARROLLO DE SISTEMAS INFORMÁTICOS"])
	assert.Equal(3, matsCantCats["INGENIERÍA DE SOFTWARE I"])
	assert.Equal(2, matsCantCats["INGENIERÍA DE SOFTWARE II"])
	assert.Equal(3, matsCantCats["INTRODUCCIÓN AL DESARROLLO DE SOFTWARE"])
	assert.Equal(8, matsCantCats["MODELACIÓN NUMÉRICA"])
	assert.Equal(9, matsCantCats["ORGANIZACIÓN DEL COMPUTADOR"])
	assert.Equal(3, matsCantCats["PARADIGMAS DE PROGRAMACIÓN"])
	assert.Equal(7, matsCantCats["PROBABILIDAD Y ESTADÍSTICA"])
	assert.Equal(1, matsCantCats["PROGRAMACIÓN CONCURRENTE"])
	assert.Equal(2, matsCantCats["REDES"])
	assert.Equal(1, matsCantCats["SIMULACIÓN"])
	assert.Equal(1, matsCantCats["SISTEMAS DISTRIBUIDOS I"])
	assert.Equal(2, matsCantCats["SISTEMAS OPERATIVOS"])
	assert.Equal(2, matsCantCats["TALLER DE PROGRAMACIÓN"])
	assert.Equal(1, matsCantCats["TALLER DE SEGURIDAD INFORMÁTICA"])
	assert.Equal(5, matsCantCats["TEORÍA DE ALGORITMOS"])
}

func TestOfertaDeComisionesQuimica2C2024(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	contenidoSiu := leerArchivoTestOfertaDeComisiones("quimica-28-12-2024.txt")

	meta, _ := ObtenerMetaData(string(contenidoSiu))
	materias := ObtenerMaterias(string(contenidoSiu))

	require.Equal(meta.Carrera, "INGENIERÍA QUÍMICA")
	require.Len(materias, 37)

	matsCantCats := make(map[string]int, len(materias))

	for _, mat := range materias {
		matsCantCats[mat.Nombre] = len(mat.Catedras)
	}

	assert.Equal(14, matsCantCats["ÁLGEBRA LINEAL"])
	assert.Equal(16, matsCantCats["ANÁLISIS MATEMÁTICO II"])
	assert.Equal(1, matsCantCats["BIOPOLÍMEROS"])
	assert.Equal(3, matsCantCats["CONOCIMIENTO DE MATERIALES METÁLICOS"])
	assert.Equal(1, matsCantCats["CONTROL ESTADÍSTICO DE PROCESOS"])
	assert.Equal(1, matsCantCats["DINÁMICA Y CONTROL DE PROCESOS"])
	assert.Equal(1, matsCantCats["DISEÑO DE PROCESOS"])
	assert.Equal(1, matsCantCats["DISEÑO DE REACTORES"])
	assert.Equal(10, matsCantCats["ELECTRICIDAD Y MAGNETISMO"])
	assert.Equal(1, matsCantCats["ELECTROQUÍMICA"])
	assert.Equal(1, matsCantCats["EMISIONES DE CONTAMINANTES QUÍMICOS Y BIOLÓGICOS"])
	assert.Equal(1, matsCantCats["ENERGÍAS RENOVABLES"])
	assert.Equal(1, matsCantCats["EVALUACIÓN DE PROYECTOS DE PLANTAS QUÍMICAS"])
	assert.Equal(1, matsCantCats["FENÓMENOS DE TRANSPORTE"])
	assert.Equal(13, matsCantCats["FÍSICA DE LOS SISTEMAS DE PARTÍCULAS"])
	assert.Equal(1, matsCantCats["FUNDAMENTOS DE PROCESOS QUÍMICOS"])
	assert.Equal(1, matsCantCats["GESTIÓN DE RECURSOS"])
	assert.Equal(1, matsCantCats["INDUSTRIAS QUÍMICAS Y PETROQUÍMICA"])
	assert.Equal(1, matsCantCats["INGENIERÍA DE BIOPROCESOS"])
	assert.Equal(1, matsCantCats["INSTALACIONES DE PLANTAS DE PROCESOS"])
	assert.Equal(1, matsCantCats["INTRODUCCIÓN A INGENIERÍA QUÍMICA"])
	assert.Equal(4, matsCantCats["INTRODUCCIÓN A LA CIENCIA DE DATOS"])
	assert.Equal(1, matsCantCats["LABORATORIO DE OPERACIONES Y PROCESOS"])
	assert.Equal(4, matsCantCats["LEGISLACIÓN Y EJERCICIO PROFESIONAL"])
	assert.Equal(8, matsCantCats["MODELACIÓN NUMÉRICA"])
	assert.Equal(1, matsCantCats["OPERACIONES UNITARIAS DE TRANSFERENCIA DE CANTIDAD DE MOVIMIENTO Y ENERGÍA"])
	assert.Equal(1, matsCantCats["OPERACIONES UNITARIAS DE TRANSFERENCIA DE MATERIA"])
	assert.Equal(1, matsCantCats["ÓPTICA"])
	assert.Equal(7, matsCantCats["PROBABILIDAD Y ESTADÍSTICA"])
	assert.Equal(2, matsCantCats["QUÍMICA ANALÍTICA INSTRUMENTAL"])
	assert.Equal(1, matsCantCats["QUÍMICA FÍSICA"])
	assert.Equal(4, matsCantCats["QUÍMICA GENERAL"])
	assert.Equal(3, matsCantCats["QUÍMICA INORGÁNICA"])
	assert.Equal(3, matsCantCats["QUÍMICA ORGÁNICA"])
	assert.Equal(1, matsCantCats["TERMODINÁMICA DE LOS PROCESOS"])
	assert.Equal(1, matsCantCats["USO EFICIENTE DE LA ENERGÍA"])
	assert.Equal(1, matsCantCats["DISEÑO DE PROCESOS"])
}
