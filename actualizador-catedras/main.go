package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/text/unicode/norm"
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

	codsMaterias, err := getCodigosMaterias()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Obteniendo últimos planes de estudio")

	planes, err := getUltimosPlanesDeEstudio()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Filtrando solo ofertas de comisiones más recientes")

	materias := filtrarMateriasMasRecientes(planes)

	log.Info(fmt.Sprintf("Encontradas %v materias", len(materias)))

	for _, m := range materias {
		// TODO: tengo que pasar esto al parser mejor
		nombre := norm.NFD.String(m.Nombre)

		var nombreSanit strings.Builder
		for _, r := range nombre {
			if !unicode.Is(unicode.Mn, r) {
				nombreSanit.WriteRune(r)
			}
		}

		nombre = strings.ToLower(nombreSanit.String())
		if _, ok := codsMaterias[nombre]; !ok {
			log.Warn(nombre)
		}
	}
}
