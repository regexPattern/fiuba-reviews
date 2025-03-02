package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

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

func getMateriasOfertasComisiones() ([]materia, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	bucket := aws.String(os.Getenv("AWS_S3_BUCKET"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	output, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: bucket,
	})

	logger := slog.Default().With(slog.String("bucket", *bucket))

	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	planes := make(chan []materia)
	sem := make(chan struct{}, MAX_REQ_CONCURRENTES)

	for _, obj := range output.Contents {
		logger := logger.With(slog.String("key", *obj.Key))

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
				logger.Error("Error obteniendo el archivo del plan", slog.String("error", err.Error()))
				return
			}

			defer obj.Body.Close()
			_, err = io.ReadAll(obj.Body)

			if err != nil {
				logger.Error("Error leyendo el contenido del plan", slog.String("error", err.Error()))
				return
			}

			planes <- []materia{{Nombre: ""}}
			<-sem
		}(obj.Key)
	}

	go func() {
		wg.Wait()
		fmt.Println("closing")
		close(planes)
	}()

	for plan := range planes {
		fmt.Println(plan)
	}

	return nil, nil
}
