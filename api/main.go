package main

func main() {
}

// package main
//
// import (
// 	"fmt"
// 	"io"
// 	"log/slog"
// 	"net/http"
//
// 	"github.com/regexPattern/fiuba-reviews/scraper_siu"
// )
//
// func main() {
// 	mux := http.NewServeMux()
//
// 	mux.HandleFunc("POST /", handler)
//
// 	s := http.Server{
// 		Addr:    ":8080",
// 		Handler: mux,
// 	}
//
// 	if err := s.ListenAndServe(); err != nil {
// 		slog.Error(err.Error())
// 	}
// }
//
// func handler(w http.ResponseWriter, r *http.Request) {
// 	data, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		slog.Error(err.Error())
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	slog.Info(fmt.Sprintf("Contenido de %v bytes recibido", len(data)))
//
// 	contenidoSiu := string(data)
//
// 	metaData, err := scraper_siu.ObtenerMetaData(contenidoSiu)
// 	if err != nil {
// 		slog.Error(err.Error())
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	slog.Info(fmt.Sprintf(
// 		"Recibida la información de %v en %vC%v",
// 		metaData.Carrera, metaData.Cuatri.Numero, metaData.Cuatri.Anio,
// 	))
//
// 	// TODO: Verificar que necesito agregar esta informacion.
//
// 	materias := scraper_siu.ObtenerMaterias(metaData.Cuatri.Contenido)
//
// 	slog.Info(fmt.Sprintf("Obtenida la información de %v materias", len(materias)))
// }
//
