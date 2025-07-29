package patch

import "strconv"

type Oferta struct {
	Carrera  string
	Materias []MateriaSiu
	Cuatri
}

type MateriaSiu struct {
	Codigo   string       `json:"codigo"`
	Nombre   string       `json:"nombre"`
	Catedras []CatedraSiu `json:"catedras"`
}

type CatedraSiu struct {
	Codigo   int          `json:"codigo"`
	Docentes []DocenteSiu `json:"docentes"`
}

type DocenteSiu struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

type Cuatri struct {
	Numero int
	Anio   int
}

func newCuatri(numero, anio string) (Cuatri, error) {
	var c Cuatri
	var err error
	c.Numero, err = strconv.Atoi(numero)
	if err != nil {
		return c, err
	}
	c.Anio, err = strconv.Atoi(anio)
	if err != nil {
		return c, err
	}
	return c, nil
}

func (c Cuatri) despuesDe(otro Cuatri) bool {
	if c.Anio == otro.Anio {
		return c.Numero > otro.Numero
	} else {
		return c.Anio > otro.Anio
	}
}
