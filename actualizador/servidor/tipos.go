package main

type OfertaCarrera struct {
	CodigoCarrera string          `db:"codigo_carrera"`
	Cuatrimestre  Cuatrimestre    `db:"cuatrimestre"`
	Materias      []OfertaMateria `db:"contenido"`
}

type Cuatrimestre struct {
	Numero int `json:"numero"`
	Anio   int `json:"anio"`
}

type Materia struct {
	Codigo string `db:"codigo" json:"codigo"`
	Nombre string `db:"nombre" json:"nombre"`
}

type OfertaMateria struct {
	Materia
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

type UltimaOfertaMateria struct {
	OfertaMateria
	Cuatrimestre
}
