package scraper_siu

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSeObtienenLosCuatrisCorrectamente(t *testing.T) {
	require := require.New(t)

	cuatris := obtenerCuatris(`
Período lectivo: 2023 - 1er Cuatrimestre
Período lectivo: 2023 - 2do Cuatrimestre
Período lectivo: 2024 - 1er Cuatrimestre
Período lectivo: 2024 - 2do Cuatrimestre
	`)

	require.Len(cuatris, 4)

	// Los cuatrimestres se retornan ordenados cronológicamente, con el más
	// reciente al final del listado.

	require.Equal(2023, cuatris[0].anio)
	require.Equal(1, cuatris[0].numero)

	require.Equal(2023, cuatris[1].anio)
	require.Equal(2, cuatris[1].numero)

	require.Equal(2024, cuatris[2].anio)
	require.Equal(1, cuatris[2].numero)

	require.Equal(2024, cuatris[3].anio)
	require.Equal(2, cuatris[3].numero)
}

func TestSeIgnorarLosCuatrisConAnioInvalido(t *testing.T) {
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
	require := require.New(t)

	cuatris := obtenerCuatris(`
Período lectivo: 2024 - 1er Cuatrimestre
CONTENIDO 1ER CUATRIMESTRE 2024
MÁS CONTENIDO 1ER CUATRIMESTRE 2024

Período lectivo: 2024 - 2do Cuatrimestre
CONTENIDO 2DO CUATRIMESTRE 2024
MÁS CONTENIDO 2DO CUATRIMESTRE 2024
	`)

	require.Equal(
		`CONTENIDO 1ER CUATRIMESTRE 2024
MÁS CONTENIDO 1ER CUATRIMESTRE 2024`,
		strings.TrimSpace(cuatris[0].contenido),
	)
	require.Equal(
		`CONTENIDO 2DO CUATRIMESTRE 2024
MÁS CONTENIDO 2DO CUATRIMESTRE 2024`,
		strings.TrimSpace(cuatris[1].contenido),
	)
}

func TestSeObtienenLasMateriasDeUnCuatriCorrectamente(t *testing.T) {
	require := require.New(t)

	materias := obtenerMateriasDeCuatri(`
Actividad: ÁLGEBRA LINEAL (CB002)
Actividad: ALGORITMOS Y ESTRUCTURAS DE DATOS (CB100)
Actividad: ANÁLISIS MATEMÁTICO II (CB001)
`)

	require.Len(materias, 3)

	require.Equal("ÁLGEBRA LINEAL", materias[0].Nombre)
	require.Equal("CB002", materias[0].Codigo)

	require.Equal("ALGORITMOS Y ESTRUCTURAS DE DATOS", materias[1].Nombre)
	require.Equal("CB100", materias[1].Codigo)

	require.Equal("ANÁLISIS MATEMÁTICO II", materias[2].Nombre)
	require.Equal("CB001", materias[2].Codigo)
}

func TestSeIgnoranLasMateriasDeTrabajoProfesional(t *testing.T) {
	materias := obtenerMateriasDeCuatri(`
Actividad: TRABAJO PROFESIONAL DE INGENIERÍA INFORMÁTICA (TA053)
Actividad: TRABAJO PROFESIONAL DE INGENIERÍA QUÍMICA (TA170)
`)

	require.Empty(t, materias)
}

func TestSeObtieneLasCatedrasDeUnaMateriaCorrectamente(t *testing.T) {
	require := require.New(t)

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

	require.Len(catedras, 14)

	codigosExpected := make([]int, 0, len(catedras))
	codigosActual := make([]int, 0, len(catedras))

	for i, cat := range catedras {
		codigosExpected = append(codigosExpected, i+1)
		codigosActual = append(codigosActual, cat.Codigo)
	}

	require.ElementsMatch(codigosExpected, codigosActual)
}

func TestSeAsignanCodigosSecuencialesALasCatedrasSinCodigo(t *testing.T) {
	require := require.New(t)

	catedras := obtenerCatedrasDeMateria(`
Comisión: CURSO: 01
Comisión: CURSO: 07
Comisión: CURSO: Caram
Comisión: CURSO VIRTUAL: ESPECIAL PARA RECURSANTES
		`)

	require.Len(catedras, 4)

	codsCats := []int{0, 0, 0, 0}

	for i, cat := range catedras {
		codsCats[i] = cat.Codigo
	}

	require.Contains(codsCats, 8)
	require.Contains(codsCats, 9)
}

func TestSeIgnoranLasCatedrasSinCodigoONombre(t *testing.T) {
	catedras := obtenerCatedrasDeMateria(`
Comisión: CURSO:
Comisión:
		`)

	require.Empty(t, catedras)
}

func TestSeIgnoranLasCatedrasParaCondicionales(t *testing.T) {
	catedras := obtenerCatedrasDeMateria(`
Comisión: CURSO: CONDICIONALES
		`)

	require.Empty(t, catedras)
}

func TestSeObtienenLosNombresDeLosDocentesCorrectamente(t *testing.T) {
	require := require.New(t)

	docentes := obtenerDocentesDeVariante(`
Docentes: BUCHWALD MARTÍN EZEQUIEL (Profesor/a Adjunto/a), PODBEREZSKI VICTOR DANIEL (Profesor/a Adjunto/a), GENENDER PEÑA EZEQUIEL DAVID (Jefe/a Trabajos Practicos)
		`)

	require.Len(docentes, 3)

	require.Contains(docentes, Docente{"BUCHWALD MARTÍN EZEQUIEL", "PROFESOR ADJUNTO"})
	require.Contains(docentes, Docente{"PODBEREZSKI VICTOR DANIEL", "PROFESOR ADJUNTO"})
	require.Contains(docentes, Docente{"GENENDER PEÑA EZEQUIEL DAVID", "JEFE TRABAJOS PRACTICOS"})
}

func TestSeRetornaNilCuandoSeEncuentraUnaCatedraSinDocentes(t *testing.T) {
	docentes := obtenerDocentesDeVariante(`Docentes: Sin docentes`)

	require.Nil(t, docentes)
}

func TestSeIgnoranLasCatedrasSinDocentes(t *testing.T) {
	require := require.New(t)

	catedras := obtenerCatedrasConDocentes(`
Comisión: CURSO: 1
Docentes: BUCHWALD MARTÍN EZEQUIEL (Profesor/a Adjunto/a), PODBEREZSKI VICTOR DANIEL (Profesor/a Adjunto/a), GENENDER PEÑA EZEQUIEL DAVID (Jefe/a Trabajos Practicos)

Comisión: CURSO: 2
Docentes: Sin docentes
		`)

	require.Len(catedras, 1)
	require.Equal(1, catedras[0].Codigo)
}

func TestSeAceptaCasingVariadoParaLosNombresDeLasMaterias(t *testing.T) {
	require := require.New(t)

	materias := obtenerMateriasDeCuatri(`
Actividad: álgebra lineal (CB002)
		`)

	require.Len(materias, 1)
	require.Equal("ÁLGEBRA LINEAL", materias[0].Nombre)
}

func TestSeAceptaCasingVariadoParaLosNombresDeLosDocentes(t *testing.T) {
	require := require.New(t)

	docentes := obtenerDocentesDeVariante(`
Docentes: BUCHWALD martín ezequiel (Profesor/a Adjunto/a)
		`)

	require.Len(docentes, 1)
	require.Contains(docentes, Docente{"BUCHWALD MARTÍN EZEQUIEL", "PROFESOR ADJUNTO"})
}

func TestSeAgregaCadaDocenteUnaSolaVezAlUnificarVariantesDeCatedras(t *testing.T) {
	require := require.New(t)

	catedras := obtenerCatedrasDeMateria(`
Comisión: CURSO: 02A
Docentes: SARRIS CLAUDIA MONICA (Profesor/a Adjunto/a), FAGES LUCIANO RODOLFO (Ayudante 1ro/a)

Comisión: CURSO: 02B
Docentes: SARRIS CLAUDIA MONICA (Profesor/a Adjunto/a), GOMEZ CIAPPONI LAUTARO (Ayudante 1ro/a)
		`)

	require.Len(catedras[0].Docentes, 3)

	require.Contains(catedras[0].Docentes, Docente{"SARRIS CLAUDIA MONICA", "PROFESOR ADJUNTO"})
	require.Contains(catedras[0].Docentes, Docente{"FAGES LUCIANO RODOLFO", "AYUDANTE 1RO"})
	require.Contains(catedras[0].Docentes, Docente{"GOMEZ CIAPPONI LAUTARO", "AYUDANTE 1RO"})
}

func leerArchivoTestOfertaDeComisiones(filename string) string {
	_, testFilename, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(testFilename)
	fullFilepath := filepath.Join(testDir, "testdata", filename)
	contenidoSiu, _ := os.ReadFile(fullFilepath)
	return string(contenidoSiu)
}

func TestOfertaDeComisionesInformatica2C2024(t *testing.T) {
	require := require.New(t)

	contenidoSiu := leerArchivoTestOfertaDeComisiones("informatica-28-12-2024.txt")
	materias := ScrapearSiu(string(contenidoSiu))

	require.Len(materias, 30)

	matsCantCats := make(map[string]int, len(materias))

	for _, mat := range materias {
		matsCantCats[mat.Nombre] = matsCantCats[mat.Nombre] + len(mat.Catedras)
	}

	require.Equal(14, matsCantCats["ÁLGEBRA LINEAL"])
	require.Equal(5, matsCantCats["ALGORITMOS Y ESTRUCTURAS DE DATOS"])
	require.Equal(16, matsCantCats["ANÁLISIS MATEMÁTICO II"])
	require.Equal(3, matsCantCats["ANÁLISIS MATEMÁTICO III"])
	require.Equal(1, matsCantCats["APRENDIZAJE AUTOMÁTICO"])
	require.Equal(1, matsCantCats["APRENDIZAJE PROFUNDO"])
	require.Equal(1, matsCantCats["ARQUITECTURA DE SOFTWARE"])
	require.Equal(3, matsCantCats["BASE DE DATOS"])
	require.Equal(2, matsCantCats["CIENCIA DE DATOS"])
	require.Equal(1, matsCantCats["COMPUTACIÓN CUÁNTICA"])
	require.Equal(1, matsCantCats["EMPRESAS DE BASE TECNOLÓGICA I"])
	require.Equal(1, matsCantCats["EMPRESAS DE BASE TECNOLÓGICA II"])
	require.Equal(1, matsCantCats["FÍSICA PARA INFORMÁTICA"])
	require.Equal(4, matsCantCats["FUNDAMENTOS DE PROGRAMACIÓN"])
	require.Equal(3, matsCantCats["GESTIÓN DEL DESARROLLO DE SISTEMAS INFORMÁTICOS"])
	require.Equal(3, matsCantCats["INGENIERÍA DE SOFTWARE I"])
	require.Equal(2, matsCantCats["INGENIERÍA DE SOFTWARE II"])
	require.Equal(3, matsCantCats["INTRODUCCIÓN AL DESARROLLO DE SOFTWARE"])
	require.Equal(8, matsCantCats["MODELACIÓN NUMÉRICA"])
	require.Equal(9, matsCantCats["ORGANIZACIÓN DEL COMPUTADOR"])
	require.Equal(3, matsCantCats["PARADIGMAS DE PROGRAMACIÓN"])
	require.Equal(7, matsCantCats["PROBABILIDAD Y ESTADÍSTICA"])
	require.Equal(1, matsCantCats["PROGRAMACIÓN CONCURRENTE"])
	require.Equal(2, matsCantCats["REDES"])
	require.Equal(1, matsCantCats["SIMULACIÓN"])
	require.Equal(1, matsCantCats["SISTEMAS DISTRIBUIDOS I"])
	require.Equal(2, matsCantCats["SISTEMAS OPERATIVOS"])
	require.Equal(2, matsCantCats["TALLER DE PROGRAMACIÓN"])
	require.Equal(1, matsCantCats["TALLER DE SEGURIDAD INFORMÁTICA"])
	require.Equal(5, matsCantCats["TEORÍA DE ALGORITMOS"])
}

func TestOfertaDeComisionesQuimica2C2024(t *testing.T) {
	require := require.New(t)

	contenidoSiu := leerArchivoTestOfertaDeComisiones("quimica-28-12-2024.txt")
	materias := ScrapearSiu(string(contenidoSiu))

	require.Len(materias, 37)

	matsCantCats := make(map[string]int, len(materias))

	for _, mat := range materias {
		matsCantCats[mat.Nombre] = len(mat.Catedras)
	}

	require.Equal(14, matsCantCats["ÁLGEBRA LINEAL"])
	require.Equal(16, matsCantCats["ANÁLISIS MATEMÁTICO II"])
	require.Equal(1, matsCantCats["BIOPOLÍMEROS"])
	require.Equal(3, matsCantCats["CONOCIMIENTO DE MATERIALES METÁLICOS"])
	require.Equal(1, matsCantCats["CONTROL ESTADÍSTICO DE PROCESOS"])
	require.Equal(1, matsCantCats["DINÁMICA Y CONTROL DE PROCESOS"])
	require.Equal(1, matsCantCats["DISEÑO DE PROCESOS"])
	require.Equal(1, matsCantCats["DISEÑO DE REACTORES"])
	require.Equal(10, matsCantCats["ELECTRICIDAD Y MAGNETISMO"])
	require.Equal(1, matsCantCats["ELECTROQUÍMICA"])
	require.Equal(1, matsCantCats["EMISIONES DE CONTAMINANTES QUÍMICOS Y BIOLÓGICOS"])
	require.Equal(1, matsCantCats["ENERGÍAS RENOVABLES"])
	require.Equal(1, matsCantCats["EVALUACIÓN DE PROYECTOS DE PLANTAS QUÍMICAS"])
	require.Equal(1, matsCantCats["FENÓMENOS DE TRANSPORTE"])
	require.Equal(13, matsCantCats["FÍSICA DE LOS SISTEMAS DE PARTÍCULAS"])
	require.Equal(1, matsCantCats["FUNDAMENTOS DE PROCESOS QUÍMICOS"])
	require.Equal(1, matsCantCats["GESTIÓN DE RECURSOS"])
	require.Equal(1, matsCantCats["INDUSTRIAS QUÍMICAS Y PETROQUÍMICA"])
	require.Equal(1, matsCantCats["INGENIERÍA DE BIOPROCESOS"])
	require.Equal(1, matsCantCats["INSTALACIONES DE PLANTAS DE PROCESOS"])
	require.Equal(1, matsCantCats["INTRODUCCIÓN A INGENIERÍA QUÍMICA"])
	require.Equal(4, matsCantCats["INTRODUCCIÓN A LA CIENCIA DE DATOS"])
	require.Equal(1, matsCantCats["LABORATORIO DE OPERACIONES Y PROCESOS"])
	require.Equal(4, matsCantCats["LEGISLACIÓN Y EJERCICIO PROFESIONAL"])
	require.Equal(8, matsCantCats["MODELACIÓN NUMÉRICA"])
	require.Equal(1, matsCantCats["OPERACIONES UNITARIAS DE TRANSFERENCIA DE CANTIDAD DE MOVIMIENTO Y ENERGÍA"])
	require.Equal(1, matsCantCats["OPERACIONES UNITARIAS DE TRANSFERENCIA DE MATERIA"])
	require.Equal(1, matsCantCats["ÓPTICA"])
	require.Equal(7, matsCantCats["PROBABILIDAD Y ESTADÍSTICA"])
	require.Equal(2, matsCantCats["QUÍMICA ANALÍTICA INSTRUMENTAL"])
	require.Equal(1, matsCantCats["QUÍMICA FÍSICA"])
	require.Equal(4, matsCantCats["QUÍMICA GENERAL"])
	require.Equal(3, matsCantCats["QUÍMICA INORGÁNICA"])
	require.Equal(3, matsCantCats["QUÍMICA ORGÁNICA"])
	require.Equal(1, matsCantCats["TERMODINÁMICA DE LOS PROCESOS"])
	require.Equal(1, matsCantCats["USO EFICIENTE DE LA ENERGÍA"])
	require.Equal(1, matsCantCats["DISEÑO DE PROCESOS"])
}
