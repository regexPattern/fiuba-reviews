package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/regexPattern/fiuba-reviews/scraper"
)

func HandlerScraper(w http.ResponseWriter, r *http.Request) {
	initClientS3()

	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(500)
		return
	}

	contenidoSiu := string(data)

	meta, err := scraper.ObtenerMetaData(contenidoSiu)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	materias := scraper.ObtenerMaterias(meta.Cuatri.Contenido)

	data, err = json.Marshal(materias)
	if err != nil {
		return
	}

	// TODO: realmente lo que quiero es guardarlas en el bucket de S3.
	_, err = w.Write(data)
	if err != nil {
		return
	}

	w.WriteHeader(200)
}

func initClientS3() *s3.Client {
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// return s3.NewFromConfig(cfg)
	return nil
}
