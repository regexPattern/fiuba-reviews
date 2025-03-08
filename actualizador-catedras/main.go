package main

import (
	"github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"
)

func init() {
	initLogger()

	logger := log.Default().WithPrefix("⚙️")

	if err := InitDBPool(logger); err != nil {
		logger.Fatal(err)
	}

	if err := initS3Client(logger); err != nil {
		logger.Fatal(err)
	}
}

func main() {
	ofertas, err := fetchOfertasDeComisiones()
	if err != nil {
		log.Fatal(err)
	}

	err = ActualizarCodigosMaterias(ofertas)
	if err != nil {
		log.Fatal(err)
	}

	err = prepActualizacionesAUltimaOfertaDeComisiones(ofertas)
	if err != nil {
		log.Fatal(err)
	}
}
