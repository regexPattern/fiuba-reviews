package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	scraper_siu "github.com/regexPattern/fiuba-reviews/scraper-siu"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const NOMBRE_BUCKET string = "fiuba-reviews-scraper-siu"

func HandlerScraper(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error(err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte("Error interno leyendo el contenido de la solicitud."))
		if err != nil {
			slog.Error(err.Error())
		}

		return
	}

	contenidoSiu := string(data)

	meta, err := scraper_siu.ObtenerMetaData(contenidoSiu)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.Error()))

		if err != nil {
			slog.Error(err.Error())
		}

		return
	}

	materias := scraper_siu.ObtenerMaterias(meta.Cuatri.Contenido)

	data, err = json.Marshal(materias)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Error interno serializando información scrapeada."))

		if err != nil {
			slog.Error(err.Error())
		}

		return
	}

	cfg, _ := config.LoadDefaultConfig(ctx)
	s3Client := s3.NewFromConfig(cfg)

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	carrera, _, _ := transform.String(t, meta.Carrera)
	carrera = strings.ToLower(strings.ReplaceAll(carrera, " ", "-"))

	uri := fmt.Sprintf("%v-%vC-%v.json", carrera, meta.Cuatri.Numero, meta.Cuatri.Anio)

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:          aws.String(NOMBRE_BUCKET),
		Key:             aws.String(uri),
		Body:            bytes.NewReader(data),
		ContentLanguage: aws.String("es"),
		Metadata: map[string]string{
			"carrera":      meta.Carrera,
			"anio":         strconv.Itoa(meta.Cuatri.Anio),
			"cuatrimestre": strconv.Itoa(meta.Cuatri.Numero),
		},
	})

	if err != nil {
		slog.Error(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Error interno almacenando la información."))

		if err != nil {
			slog.Error(err.Error())
		}

		return
	}

	slog.Info(fmt.Sprintf("Escrito archivo `%v` con éxito.", uri))

	w.WriteHeader(http.StatusCreated)
}
