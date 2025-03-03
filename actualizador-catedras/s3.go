package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/log"
)

const BUCKET_NAME_ENV string = ""
const MAX_REQ_CONCURRENTES int = 5

type materia struct {
	Codigo   string    `json:"codigo"`
	Nombre   string    `json:"nombre"`
	Catedras []catedra `json:"catedras"`
}

type catedra struct {
	Codigo   int       `json:"codigo"`
	Docentes []docente `json:"docentes"`
}

type docente struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

var client *s3.Client

func newS3Client() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client = s3.NewFromConfig(cfg)
	return nil
}

func getPlanesDeEstudio() ([]materia, error) {
	bucket := aws.String(os.Getenv("AWS_S3_BUCKET"))
	logger := log.Default().WithPrefix("S3 🪣").With("bucket", *bucket)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	output, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: bucket,
	})

	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Obtenidos %v planes", len(output.Contents)))

	var wg sync.WaitGroup
	planes := make(chan []materia)
	sem := make(chan struct{}, MAX_REQ_CONCURRENTES)

	for _, obj := range output.Contents {
		logger := logger.With("obj", *obj.Key)

		wg.Add(1)

		go func(key *string) {
			sem <- struct{}{}

			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			obj, err := client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: bucket,
				Key:    key,
			})

			if err != nil {
				// TODO: ocupar errgroup para manejar esto
				logger.Error("Error obteniendo el archivo del plan", slog.String("error", err.Error()))
				return
			}

			defer obj.Body.Close()
			data, err := io.ReadAll(obj.Body)

			if err != nil {
				// TODO: ocupar errgroup para manejar esto
				logger.Error("Error leyendo el contenido del plan", slog.String("error", err.Error()))
				return
			}

			var plan []materia
			_ = json.Unmarshal(data, &plan)

			// TODO: obtener la metadata del plan (ya la tengo seguro) para imprimir el nombre del plan bonito
			logger.Info(fmt.Sprintf("Obtenidas %v materias en el plan", len(plan)))

			planes <- plan
			<-sem
		}(obj.Key)
	}

	go func() {
		wg.Wait()
		close(planes)
	}()

	// TODO: lo que necesito es tirar a todas las materias en un hashset
	plan := <-planes
	fmt.Println(plan)

	return nil, nil
}
