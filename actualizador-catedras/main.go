package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"
)

func init() {
	initLogger()
}

func main() {
	logger := log.Default().WithPrefix("⚙️")

	if err := conectarDb(logger); err != nil {
		logger.Fatal(err)
	}

	if err := conectarS3(logger); err != nil {
		logger.Fatal(err)
	}

	ofertas, err := getOfertasComisiones()
	if err != nil {
		log.Fatal(err)
	}

	err = updateCodigosMaterias(ofertas)
	if err != nil {
		log.Fatal(err)
	}

	patches, err := getPatchesActualizacion(ofertas)
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(newModeloApp(patches))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
