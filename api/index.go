package api

import (
	"fmt"
	"net/http"
)

func Another(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
}
