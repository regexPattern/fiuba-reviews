package scraper_siu

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// https://github.com/FdelMazo/FIUBA-Plan/blob/master/src/siuparser.js
var reCuatri *regexp.Regexp = regexp.MustCompile(`Período lectivo: (\d{4}) - ([^\n]+)\n`)
var reMateria *regexp.Regexp = regexp.MustCompile(`Actividad:\s*([^\s\(]+(?:\s[^\s\(]+)*)\s?\(([^\)]+)\)`)
var reCatedra *regexp.Regexp = regexp.MustCompile(`Comisión:\s*(?:CURSO:\s*)?(?:(\d{1,2})([a-cA-C])?|([a-záéíóúñA-ZÁÉÍÓÚÑ/ ]+))`)
var reDocente *regexp.Regexp = regexp.MustCompile(`([A-ZÁÉÍÓÚÑ ]+)\s*\(([^\)]+)\)`)

type cuatri struct {
	anio   int
	numero int
	data   string
}

// TODO: Voy a necesitar una forma de asegurarme que pueda asociar estas
// materias a las materias que tengo en la base de datos. En el SIU tengo los
// códigos de las materias, en el listado de planes de FIUBA Map no (no tengo
// los oficiales).
//
// Puedo ir recuperando los códigos finales de las materias a medida me vayan
// llegando los listados de materias, para así mapearlos con los de FIUBA Map.
// Inicialmente me puedo quedar con las cátedras que ya tenía en Dolly. Este
// paquete no se encarga de esto, solo de la funcionalidad de extracción de
// esta información desde el SIU.
type Materia struct {
	Nombre string
	Codigo string
	data   string
}

type Catedra struct {
	Codigo   int
	Docentes []Docente
}

type Docente struct {
	Nombre string
	Rol    string
}

func ObtenerCatedras(contenidoSiu string) []Catedra {
	cuatris := obtenerCuatris(contenidoSiu)

	for _, cuatri := range cuatris {
		materias := obtenerMaterias(cuatri.data)
		fmt.Println(materias)
	}

	return nil
}

func obtenerCuatris(contenidoSiu string) []cuatri {
	locs := reCuatri.FindAllStringSubmatchIndex(contenidoSiu, -1)
	cuatris := make([]cuatri, 0, 2)

	for i := 0; i < len(locs); i++ {
		loc := locs[i]

		inicio := loc[1]
		var fin int

		if i+1 < len(locs) {
			fin = locs[i+1][0]
		} else {
			fin = len(contenidoSiu)
		}

		nombreStr := contenidoSiu[loc[4]:loc[5]]
		if !strings.HasSuffix(nombreStr, "Cuatrimestre") {
			// Estamos procesando un curso de verano.
			continue
		}

		// Se asume que en si ya estamos en el formato correcto, el valor del
		// cuatrimestre va a ser correcto y va tener el formato correcto (valor
		// ASCII para '1' o '2').
		numero := int(nombreStr[0]) - '0'

		anioStr := contenidoSiu[loc[2]:loc[3]]
		anio, err := strconv.Atoi(anioStr)
		if err != nil {
			continue
		}

		data := contenidoSiu[inicio:fin]

		cuatris = append(cuatris, cuatri{anio, numero, data})
	}

	return cuatris
}

func obtenerMaterias(dataCuatri string) []Materia {
	locs := reMateria.FindAllStringSubmatchIndex(dataCuatri, -1)
	materias := make([]Materia, 0, len(locs))

	for i := 0; i < len(locs); i++ {
		loc := locs[i]

		inicio := loc[1]
		var fin int

		if i+1 < len(locs) {
			fin = locs[i+1][0]
		} else {
			fin = len(dataCuatri)
		}

		nombre := dataCuatri[loc[2]:loc[3]]
		if strings.HasPrefix(nombre, "TRABAJO PROFESIONAL") {
			continue
		}

		codigo := dataCuatri[loc[4]:loc[5]]
		data := dataCuatri[inicio:fin]

		materias = append(materias, Materia{nombre, codigo, data})
	}

	return materias
}

func obtenerCatedras(dataMateria string) []Catedra {
	locs := reCatedra.FindAllStringSubmatchIndex(dataMateria, -1)
	catedrasMap := make(map[int]*Catedra, 0)

	for i := 0; i < len(locs); i++ {
		loc := locs[i]

		inicio := loc[1]
		var fin int

		if i+1 < len(locs) {
			fin = locs[i+1][0]
		} else {
			fin = len(dataMateria)
		}

		// El regex para obtener el nombre de la cátedra tiene cuatro grupos:
		// 	• 1: La captura completa.
		// 	• 2: El número de curso. Solo el número, no la letra de la
		// 		 variante. En el caso de que una cátedra sea única y no tenga
		// 		 código, este grupo no matchea. En este caso las cátedras
		// 		 tienen nombre en vez de número.
		// 	• 3: La letra que representa a la variante del curso. En el caso de
		// 		 que una cátedra sea única este grupo no matchea. Lo mismo para
		// 		 cátedras sin variantes.
		// 	• 4: Solo matchea para las cátedras únicas que tienen nombre en vez
		// 		 de número. Contiene el nombre matcheado.
		// En el caso de cátedras sin código en el SIU, fiuba-reviews les
		// asigna el número 1.

		var cod int

		if loc[6] != -1 {
			if dataMateria[loc[6]:loc[7]] == "CONDICIONALES" {
				continue
			}

			cod = 1
		} else {
			// Cualquier formato no numérico directamente no es considerado
			// como un código, por lo que matchea el regex con el grupo 3. Es
			// decir, esta serialización no puede fallar.
			cod, _ = strconv.Atoi(dataMateria[loc[2]:loc[3]])
		}

		docentes := obtenerDocentes(dataMateria[inicio:fin])

		if cat, ok := catedrasMap[cod]; ok {
			cat.Docentes = append(cat.Docentes, docentes...)
		} else {
			catedrasMap[cod] = &Catedra{cod, docentes}
		}
	}

	catedras := make([]Catedra, 0, len(catedrasMap))

	for _, cat := range catedrasMap {
		catedras = append(catedras, *cat)
	}

	return catedras
}

func obtenerDocentes(dataCatedra string) []Docente {
	matches := reDocente.FindAllStringSubmatch(dataCatedra, -1)
	docentes := make([]Docente, 0, len(matches))

	for i := 0; i < len(matches); i++ {
		nombre, rol := matches[i][1], matches[i][2]

		nombre = strings.TrimSpace(nombre)
		if nombre == "A DESIGNAR A DESIGNAR" {
			continue
		}

		rol = strings.TrimSpace(rol)

		docentes = append(docentes, Docente{nombre, rol})
	}

	return docentes
}
