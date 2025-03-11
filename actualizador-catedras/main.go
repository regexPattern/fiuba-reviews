package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"
)

func init() {
	initLogger()

	logger := log.Default().WithPrefix("⚙️")

	if err := initDbPool(logger); err != nil {
		logger.Fatal(err)
	}

	if err := initS3Client(logger); err != nil {
		logger.Fatal(err)
	}
}

func main() {
	ofertas, err := getOfertasComisiones()
	if err != nil {
		log.Fatal(err)
	}

	err = updateCodigosMaterias(ofertas)
	if err != nil {
		log.Fatal(err)
	}

	patches, err := getPatchesMaterias(ofertas)
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(newModeloApp(patches))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
