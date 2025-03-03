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
	_, err := getCodigosMaterias()
	if err != nil {
		log.Fatal(err)
	}

	_, err = getPlanesDeEstudio()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(cods)
	// fmt.Println(mats)
}
