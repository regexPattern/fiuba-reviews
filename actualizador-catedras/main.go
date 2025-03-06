package main

import (
	"github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"
)

func init() {
	initLogger()

	if err := initDbConn(); err != nil {
		log.Fatal(err)
	}

	if err := initS3Client(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Info("Obteniendo códigos de materias")

	codigosMaterias, err := fetchNombresACodigosMateriasDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Obteniendo últimos planes de estudio")

	planes, err := fetchPlanesDeEstudio()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Filtrando solo ofertas de comisiones más recientes")

	materias := filtrarMateriasMasRecientes(planes)

	for _, m := range materias {
		if _, ok := codigosMaterias[m.Nombre]; !ok {
			log.Warn("Materia no está en la base de datos.", "codigo", m.Codigo, "nombre", m.Nombre)
		}
	}
}
