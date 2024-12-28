package main

import (
	"fmt"
	"io"
	"os"

	"github.com/regexPattern/fiuba-reviews/pkg/scraper_siu"
)

func main() {
	bytes, _ := io.ReadAll(os.Stdin)
	materias := scraper_siu.ScrapearSiu(string(bytes))

	for _, mat := range materias {
		fmt.Println("MATERIA:", mat.Nombre)
		for _, cat := range mat.Catedras {
			fmt.Println("CATEDRA:", cat.Codigo)
			for _, doc := range cat.Docentes {
				fmt.Println(doc.Nombre, doc.Rol)
			}
			fmt.Println()
		}
	}
}
