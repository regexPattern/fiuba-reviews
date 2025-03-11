package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/log"
)

const bucketEnvVarName string = "AWS_S3_BUCKET"
const MAX_REQ_CONCURRENTES int = 5

var client *s3.Client
var bucketKey *string

// conectarDb configura y crea el cliente de S3 utilizando. No se conecta al
// bucket, pero sí configura el identificador del mismo para usos posteriores.
func conectarS3(logger *log.Logger) error {
	bucketEnv, ok := os.LookupEnv(bucketEnvVarName)
	if !ok {
		return fmt.Errorf("variable de entorno `%v` no configurada", bucketEnvVarName)
	}

	bucketKey = aws.String(bucketEnv)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.Error(err)
		return errors.New("error cargando la configuración de AWS")
	}

	logger.Info("configuración de AWS cargada exitosamente")

	client = s3.NewFromConfig(cfg)

	logger.Info("cliente de S3 creado exitosamente")

	return nil
}
