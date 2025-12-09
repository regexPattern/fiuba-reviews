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

type MateriaConActualizaciones struct {
	Codigo             string             `db:"codigo"`
	Nombre             string             `db:"nombre"` // normalizado con lower(unaccent())
	DocentesPendientes []DocentePendiente `json:"docentes_pendientes"`
	CatedrasNuevas     []CatedraNueva     `json:"catedras_nuevas"`
}

type DocentePendiente struct {
	NombreSiu       string         `json:"nombre_siu"`
	Rol             string         `json:"rol"`
	PosiblesMatches []DocenteMatch `json:"posibles_matches"`
}

type DocenteMatch struct {
	Codigo    string  `json:"codigo"`
	NombreDb  string  `json:"nombre_db"`
	Similitud float64 `json:"similitud"`
}

type CatedraNueva struct {
	Nombre   string           `json:"nombre"`
	Docentes []DocenteCatedra `json:"docentes"`
}

type DocenteCatedra struct {
	NombreSiu     string  `json:"nombre_siu"`
	Rol           string  `json:"rol"`
	CodigoDocente *string `json:"codigo_docente"` // nil si no está resuelto
}
