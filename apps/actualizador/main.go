package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/log"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/indexador"
)

func main() {
	setupLogger()

	var client *s3.Client
	if c, err := setupS3Client(); err != nil {
		slog.Error("no se pudo inicializar el cliente de S3", "error", err)
		os.Exit(1)
	} else {
		client = c
	}

	bucketName := os.Getenv("AWS_S3_BUCKET")
	if bucketName == "" {
		slog.Error("la variable de entorno AWS_S3_BUCKET no está definida")
		os.Exit(1)
	}

	var ofertas []*indexador.Oferta
	idx := indexador.New(client, &bucketName)

	if o, err := idx.IndexarOfertasComisiones(); err != nil {
		slog.Error("no se pudieron indexar las ofertas", "error", err)
		os.Exit(1)
	} else {
		ofertas = o
	}

	fmt.Println(ofertas)
}

func setupLogger() {
	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.RFC3339,
		Level:           log.InfoLevel,
	})

	slog.SetDefault(slog.New(logger))
}

func setupS3Client() (*s3.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}
