package scraper_siu

import (
	"regexp"
	"strconv"
	"strings"
)

// https://github.com/FdelMazo/FIUBA-Plan/blob/master/src/siuparser.js
var reCuatri *regexp.Regexp = regexp.MustCompile(`Período lectivo: (\d{4}) - (\d).*`)
var reMateria *regexp.Regexp = regexp.MustCompile(`Actividad: ([^\s\(]+(?:\s[^\s\(]+)*)\s?\(([^\)]+)\)`)
var reCatedra *regexp.Regexp = regexp.MustCompile(`(?i)Comisión: (?:(?:CURSO:? ?)?(\d{1,2})([a-cA-C])?|CURSO:? ?([a-záéíóúñA-ZÁÉÍÓÚÑ ]+))`)
var reDocente *regexp.Regexp = regexp.MustCompile(`([a-záéíóúñA-ZÁÉÍÓÚÑ ]+)\s*\(([^\)]+)\)`)

// Los cuatrimestres no contienen las materias parseadas, ya que esta tarea
// realmente solo se quiere hacer con el último cuatrimestre, para actualizar
// los registros.
type cuatri struct {
	anio      int
	numero    int
	contenido string
}

type Materia struct {
	Nombre   string
	Codigo   string
	Catedras []Catedra
}

type Catedra struct {
	Codigo   int
	Docentes []Docente
}

type Docente struct {
	Nombre string
	Rol    string
}

func ScrapearSiu(contenidoSiu string) []Materia {
	cuatris := obtenerCuatris(contenidoSiu)
	ultimoCuat := cuatris[len(cuatris)-1]
	return obtenerMaterias(ultimoCuat.contenido)
}

func obtenerCuatris(contenidoSiu string) []cuatri {
	locs := reCuatri.FindAllStringSubmatchIndex(contenidoSiu, -1)

	// En el SIU nunca hay más de dos cuatrimestres listados al mismo tiempo.
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

		anioStr := contenidoSiu[loc[2]:loc[3]]
		numeroStr := contenidoSiu[loc[4]:loc[5]]

		anio, _ := strconv.Atoi(anioStr)
		// El regex matchea un solo dígito.
		numero := int(numeroStr[0]) - '0'

		cuatris = append(cuatris, cuatri{anio, numero, contenidoSiu[inicio:fin]})
	}

	return cuatris
}

func obtenerMaterias(contenidoCuatri string) []Materia {
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
		catedras := obtenerCatedras(contenidoCuatri[inicio:fin])

		materias = append(materias, Materia{nombre, codigo, catedras})
	}

	return materias
}

func obtenerCatedrasConDocentes(contenidoMateria string) []Catedra {
	catedras := obtenerCatedras(contenidoMateria)
	catedrasConDocentes := make([]Catedra, 0, len(catedras))

	for _, cat := range catedras {
		if len(cat.Docentes) != 0 {
			catedrasConDocentes = append(catedrasConDocentes, cat)
		}
	}

	return catedrasConDocentes
}

func obtenerCatedras(contenidoMateria string) []Catedra {
	locs := reCatedra.FindAllStringSubmatchIndex(contenidoMateria, -1)
	catedrasMap := make(map[int]*Catedra, 0)

	for i := 0; i < len(locs); i++ {
		loc := locs[i]

		inicio := loc[1]
		var fin int

		if i+1 < len(locs) {
			fin = locs[i+1][0]
		} else {
			fin = len(contenidoMateria)
		}

		// El regex para obtener el nombre de la cátedra tiene tres grupos (sin
		// contar la captura completa):
		// 1: El número de curso. Solo el número, no la letra de la
		//    variante. En el caso de que una cátedra sea única y no tenga
		//    código, este grupo no matchea. En este caso las cátedras
		//    tienen nombre en vez de número.
		// 2: La letra que representa a la variante del curso. En el caso de
		//    que una cátedra sea única este grupo no matchea. Lo mismo para
		//    cátedras sin variantes.
		// 3: Solo matchea para las cátedras únicas que tienen nombre en vez
		//    de número. Contiene el nombre matcheado.
		// En el caso de cátedras sin código en el SIU, FIUBA Reviews les
		// asigna el número 1.

		var cod int

		if loc[6] != -1 {
			if contenidoMateria[loc[6]:loc[7]] == "CONDICIONALES" {
				continue
			}

			cod = 1
		} else {
			// Cualquier formato no numérico directamente no es considerado
			// como un código, por lo que matchea el regex con el grupo 3. Es
			// decir, esta serialización no puede fallar.
			cod, _ = strconv.Atoi(contenidoMateria[loc[2]:loc[3]])
		}

		// Acá docentes puede retornar `nil`, sin embargo, no vamos a hacer el
		// filtrado de las cátedras ahora, sino que simplemente guardamos la
		// cátedra o extendemos los docentes, y luego podemos hacer un filtrado
		// de cátedras sin docentes.
		docentes := obtenerDocentes(contenidoMateria[inicio:fin])

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

func obtenerDocentes(contenidoCatedra string) []Docente {
	matches := reDocente.FindAllStringSubmatch(contenidoCatedra, -1)
	if len(matches) == 0 {
		// La cátedra no tiene docentes (es un formato válido, no un error).
		return nil
	}

	docentes := make([]Docente, 0, len(matches))

	for i := 0; i < len(matches); i++ {
		nombre, rol := matches[i][1], matches[i][2]

		nombre = strings.ToUpper(strings.TrimSpace(nombre))
		if nombre == "A DESIGNAR A DESIGNAR" {
			continue
		}

		rol = strings.ReplaceAll(strings.TrimSpace(rol), "/a", "")

		docentes = append(docentes, Docente{nombre, rol})
	}

	return docentes
}
