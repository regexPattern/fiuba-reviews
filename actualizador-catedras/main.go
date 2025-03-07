package main

import (
	"github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"
)

func init() {
	initLogger()

	logger := log.Default().WithPrefix("⚙️")

	if err := initDBConn(logger); err != nil {
		logger.Fatal(err)
	}

	if err := initS3Client(logger); err != nil {
		logger.Fatal(err)
	}
}

func main() {
	codigosMaterias, err := selectNombresACodigosMaterias()
	if err != nil {
		log.Fatal(err)
	}

	ofertas, err := fetchOfertasDeComisiones()
	if err != nil {
		log.Fatal(err)
	}

	err = syncCodigosMaterias(codigosMaterias, ofertas)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Filtrando solo ofertas de comisiones más recientes")

	filtrarMateriasMasRecientes(ofertas)

	// Lo que quiero hacer ahora es actualizar las catedras de las materias que
	// tienen actualizacion.
	// 1. Las materias que tienen actualizacion son aquellas que no tienen
	// registro en `actualizacion_catedras` o cuyo registro fue actualizado en
	// el cuatrimestre anterior. SOLO MATERIAS DE LOS NUEVOS PLANES.

	_ = updateCatedrasMaterias()
}
