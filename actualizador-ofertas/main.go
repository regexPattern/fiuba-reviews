package main

import (
	"fmt"

	"github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"
	"github.com/regexPattern/fiuba-reviews/actualizador-ofertas/actualizador"
	"github.com/regexPattern/fiuba-reviews/actualizador-ofertas/config"
)

func main() {
	log.SetDefault(config.NewLogger())

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	ps, err := actualizador.GetPatches(cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ps)

	// if tui.Run(ps) != nil {
	// 	log.Fatal(err)
	// }
}
