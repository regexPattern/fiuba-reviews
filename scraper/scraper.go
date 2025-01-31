package scraper

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Se utilizaron los regexes de FIUBA Plan como inspiración: https://github.com/FdelMazo/FIUBA-Plan/blob/master/src/siuparser.js
var reCarrera *regexp.Regexp = regexp.MustCompile(`Propuesta: ([a-záéíóúñA-ZÁÉÍÓÚÑ ]+)`) // https://regex101.com/r/cfElw2/2
var reCuatri *regexp.Regexp = regexp.MustCompile(`Período lectivo: (\d{4}) - (\d).*`)                                                       // https://regex101.com/r/b4DVgP/1
var reMateria *regexp.Regexp = regexp.MustCompile(`Actividad: ([^\s\(]+(?:\s[^\s\(]+)*)\s?\(([^\)]+)\)`)                                    // https://regex101.com/r/L7SlFt/2
var reCatedra *regexp.Regexp = regexp.MustCompile(`(?i)Comisión: (?:(?:CURSO:? ?)?(\d{1,2})([a-cA-C])?|CURSO:? ?([a-záéíóúñA-ZÁÉÍÓÚÑ ]+))`) // https://regex101.com/r/oli0zO/2
var reDocente *regexp.Regexp = regexp.MustCompile(`([a-záéíóúñA-ZÁÉÍÓÚÑ ]+)\s*\(([^\)]+)\)`)                                                // https://regex101.com/r/IfwXK0/1

type MetaData struct {
	Cuatri  Cuatri
	Carrera string
}

type Cuatri struct {
	Anio   int
	Numero int

	// Texto plano que contiene la información de las materias de cada
	// cuatrimestre. Esto se hace ya que solo nos interesa parsear las materias
	// de cuatrimestres selectos on-demand, no hay necesidad de parsear la
	// información de todos los cuatrimestres disponibles.
	Contenido string
}

type Materia struct {
	Nombre   string    `json:"nombre"`
	Codigo   string    `json:"codigo"`
	Catedras []Catedra `json:"catedras"`
}

type Catedra struct {
	Codigo   int       `json:"codigo"`
	Docentes []Docente `json:"docentes"`
}

type Docente struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

type By func(p1, p2 *Cuatri) bool

type cuatriSorter struct {
	cuatris []Cuatri
	by      By
}

func (by By) Sort(cuatris []Cuatri) {
	cs := &cuatriSorter{
		cuatris: cuatris,
		by:      by,
	}
	sort.Sort(cs)
}

func (s *cuatriSorter) Len() int {
	return len(s.cuatris)
}

func (s *cuatriSorter) Swap(i, j int) {
	s.cuatris[i], s.cuatris[j] = s.cuatris[j], s.cuatris[i]
}

func (s *cuatriSorter) Less(i, j int) bool {
	return s.by(&s.cuatris[i], &s.cuatris[j])
}

// ObtenerMetaData encuentra la carrera a la que corresponde la información del
// SIU enviada y el cuatrimestre relevante de la misma (el último
// cuatrimestre).
func ObtenerMetaData(contenidoSiu string) (MetaData, error) {
	var metaData MetaData

	carrera, err := obtenerCarrera(contenidoSiu)
	if err != nil {
		return metaData, err
	}

	cuatris := obtenerCuatris(contenidoSiu)
	if len(cuatris) == 0 {
		return metaData, fmt.Errorf("No se encontraron cuatrimestres")
	}

	metaData.Carrera = carrera
	metaData.Cuatri = cuatris[len(cuatris)-1]

	return metaData, nil
}

func obtenerCarrera(contenidoSiu string) (string, error) {
	matches := reCarrera.FindStringSubmatch(contenidoSiu)

	if len(matches) == 0 || matches[1] == "" {
		return "", fmt.Errorf("No se encontró la carrera")
	}

	return strings.ToUpper(matches[1]), nil
}

func obtenerCuatris(contenidoSiu string) []Cuatri {
	locs := reCuatri.FindAllStringSubmatchIndex(contenidoSiu, -1)

	// En el SIU nunca hay más de dos cuatrimestres listados al mismo tiempo.
	// De igual forma, este es un vector dinámico.
	cuatris := make([]Cuatri, 0, 2)

	for i := 0; i < len(locs); i++ {
		loc := locs[i]

		inicio := loc[1]
		var fin int

		if i+1 < len(locs) {
			fin = locs[i+1][0]
		} else {
			fin = len(contenidoSiu)
		}

		anioStr := contenidoSiu[loc[2]:loc[3]]
		numeroStr := contenidoSiu[loc[4]:loc[5]]

		anio, _ := strconv.Atoi(anioStr)
		// El regex matchea un solo dígito directamente.
		numero := int(numeroStr[0]) - '0'

		cuatris = append(cuatris, Cuatri{anio, numero, contenidoSiu[inicio:fin]})
	}

	By(func(p1, p2 *Cuatri) bool {
		if p1.Anio < p2.Anio {
			return true
		} else if p1.Anio > p2.Anio {
			return false
		} else {
			return p1.Numero < p2.Numero
		}
	}).Sort(cuatris)

	return cuatris
}

// ObtenerMaterias scrapea el contenido del cuatrimestre y retorna la
// información de las materias del mismo.
func ObtenerMaterias(contenidoCuatri string) []Materia {
	locs := reMateria.FindAllStringSubmatchIndex(contenidoCuatri, -1)
	materias := make([]Materia, 0, len(locs))

	for i := 0; i < len(locs); i++ {
		loc := locs[i]

		inicio := loc[1]
		var fin int

		if i+1 < len(locs) {
			fin = locs[i+1][0]
		} else {
			fin = len(contenidoCuatri)
		}

		nombre := strings.ToUpper(contenidoCuatri[loc[2]:loc[3]])
		if strings.Contains(nombre, "TRABAJO PROFESIONAL") {
			continue
		}

		codigo := contenidoCuatri[loc[4]:loc[5]]
		catedras := obtenerCatedrasDeMateria(contenidoCuatri[inicio:fin])

		materias = append(materias, Materia{nombre, codigo, catedras})
	}

	return materias
}

func obtenerCatedrasConDocentes(contenidoMateria string) []Catedra {
	catedrasMat := obtenerCatedrasDeMateria(contenidoMateria)
	catedras := make([]Catedra, 0, len(catedrasMat))

	for _, cat := range catedrasMat {
		if len(cat.Docentes) != 0 {
			catedras = append(catedras, cat)
		}
	}

	return catedras
}

func obtenerCatedrasDeMateria(contenidoMateria string) []Catedra {
	locs := reCatedra.FindAllStringSubmatchIndex(contenidoMateria, -1)

	catedrasMap := make(map[int]*Catedra)
	catedraDocentesMap := make(map[int]map[Docente]bool)

	// Sirve para llevar control de la cantidad de cátedras sin código. Su
	// valor es siempre <= 0 para hacer las veces de flag para distinguir
	// aquellas cátedras que necesitan que se les autoasigne un código.
	var catsSinCodigoAccum int

	for i := 0; i < len(locs); i++ {
		loc := locs[i]

		inicio := loc[1]
		var fin int

		if i+1 < len(locs) {
			fin = locs[i+1][0]
		} else {
			fin = len(contenidoMateria)
		}

		var codigo int

		// El regex para obtener el nombre de la cátedra tiene tres grupos (sin
		// contar la captura completa):
		//
		// 1: El número de curso, sin la letra de la variante. En el caso de
		//    que una cátedra sea única y no tenga código, este grupo no
		//    matchea.
		//
		// 2: La letra de la variante del curso. En el caso de que una cátedra
		//    sea única este grupo no matchea. Lo mismo para cátedras sin
		//    variantes.
		//
		// 3: Solo matchea para las cátedras que tienen nombre en vez de
		//    número. Contiene el nombre matcheado. En este caso lo que se hace
		//    es autoasignar un código a estas cátedras.
		//
		// Podés visualizar el patrón acá: https://regex101.com/r/oli0zO/2

		if loc[6] != -1 { // Matchea el tercer grupo (con nombre y sin código).
			if contenidoMateria[loc[6]:loc[7]] == "CONDICIONALES" {
				continue
			}

			catsSinCodigoAccum--
			codigo = catsSinCodigoAccum
		} else {
			// Si no matchea el tercer grupo entonces matchea el primero y/o
			// segundo, cuyo patrón matchea dígitos numéricos, por lo que esta
			// serialización no puede fallar, ya sabemos que tenemos un dígito
			// válido.
			codigo, _ = strconv.Atoi(contenidoMateria[loc[2]:loc[3]])
		}

		docentesMap := obtenerDocentesDeVariante(contenidoMateria[inicio:fin])

		if docentesCat, ok := catedraDocentesMap[codigo]; ok {
			for doc := range docentesMap {
				docentesCat[doc] = true
			}
		} else {
			catedraDocentesMap[codigo] = docentesMap
		}

		catedrasMap[codigo] = &Catedra{Codigo: codigo}
	}

	catedras := make([]Catedra, 0, len(catedrasMap))

	// Las cátedras que matchearon el grupo 3 por no tener código, van a ser
	// asignadas códigos numéricos sequenciales que inician a partir del máximo
	// (+1) código de las cátedras de esa materia que si tienen código.
	//
	// Ver el test TestSeAsignanCodigosSecuencialesALasCatedrasSinCodigo para
	// contrastar con un ejemplo.
	var maxCodigoCat int

	for _, cat := range catedrasMap {
		if cat.Codigo > maxCodigoCat {
			maxCodigoCat = cat.Codigo
		}
	}

	for cod, cat := range catedrasMap {
		// Usamos un código negativo como flag para saber a qué cátedras se les
		// debe autoasignar un código.
		if cat.Codigo < 0 {
			maxCodigoCat++
			cat.Codigo = maxCodigoCat
		}

		docentesMap := catedraDocentesMap[cod]
		docentes := make([]Docente, 0, len(docentesMap))

		for doc := range docentesMap {
			docentes = append(docentes, doc)
		}

		cat.Docentes = docentes
		catedras = append(catedras, *cat)
	}

	return catedras
}

func obtenerDocentesDeVariante(contenidoCatedra string) map[Docente]bool {
	matches := reDocente.FindAllStringSubmatch(contenidoCatedra, -1)
	if len(matches) == 0 {
		// La cátedra no tiene docentes (es un formato válido, no un error).
		return nil
	}

	docentesMap := make(map[Docente]bool)

	for i := 0; i < len(matches); i++ {
		nombre, rol := matches[i][1], matches[i][2]

		nombre = strings.ToUpper(strings.TrimSpace(nombre))
		if nombre == "A DESIGNAR A DESIGNAR" {
			continue
		}

		rol = strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(rol), "/a", ""))
		docentesMap[Docente{nombre, rol}] = true
	}

	return docentesMap
}
