package api

import (
	"fmt"
	"net/http"

	"github.com/regexPattern/fiuba-reviews/scraper"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	scraper.ObtenerMaterias("")
	fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
}
