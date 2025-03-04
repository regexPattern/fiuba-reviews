package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/log"
	"golang.org/x/sync/errgroup"
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

	var eg errgroup.Group
	eg.SetLimit(MAX_REQ_CONCURRENTES)

	planes := make(chan []materia, len(output.Contents))

	for _, obj := range output.Contents {
		logger := logger.With("objKey", *obj.Key)

		eg.Go(func() error {
			plan, err := procPlanDeEstudio(logger, &s3.GetObjectInput{
				Bucket: bucket,
				Key:    obj.Key,
			})

			if err == nil {
				planes <- plan
			}

			return err
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	close(planes)

	// TODO: agregar las materias a un hashset segun su codigo
	for plan := range planes {
		fmt.Println(len(plan))
	}

	return nil, nil
}

func procPlanDeEstudio(logger *log.Logger, objInputOpts *s3.GetObjectInput) ([]materia, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logger.Debug("Obteniendo contenido del plan")

	obj, err := client.GetObject(ctx, objInputOpts)
	if err != nil {
		return nil, err
	}

	defer obj.Body.Close()
	data, err := io.ReadAll(obj.Body)

	if err != nil {
		return nil, err
	}

	var plan []materia
	_ = json.Unmarshal(data, &plan)

	// TODO: obtener la metadata del plan (ya la tengo seguro) para imprimir el nombre del plan bonito
	logger.Info(fmt.Sprintf("Obtenidas %v materias en el plan", len(plan)))

	return plan, nil
}
