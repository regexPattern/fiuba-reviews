package main

import (
	"github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"
)

func init() {
	newLogger()

	if err := newDbConn(); err != nil {
		log.Fatal(err)
	}

	if err := newS3Client(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Info("Obteniendo códigos de materias")

	_, err := getCodigosMaterias()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Obteniendo últimos planes de estudio")

	planes, err := getUltimosPlanesDeEstudio()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Filtrando solo ofertas de comisiones más recientes")

	filtrarMateriasMasRecientes(planes)
}
