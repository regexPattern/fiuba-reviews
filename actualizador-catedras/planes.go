package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/log"
	"golang.org/x/sync/errgroup"
)

const MAX_REQ_CONCURRENTES int = 5

type plan struct {
	carrera  string
	cuatri   cuatri
	materias []materia
}

type cuatri struct {
	numero int
	anio   int
}

func fetchPlanesDeEstudio() ([]plan, error) {
	bucketName := aws.String(os.Getenv("AWS_S3_BUCKET"))
	logger := log.Default().WithPrefix("S3 🪣").With("bucket", *bucketName)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	bucket, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: bucketName,
	})

	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Obtenidos %v planes", len(bucket.Contents)))

	var eg errgroup.Group
	eg.SetLimit(MAX_REQ_CONCURRENTES)

	ch := make(chan plan, len(bucket.Contents))

	for _, obj := range bucket.Contents {
		logger := logger.With("objKey", *obj.Key)

		eg.Go(func() error {
			plan, err := serPlanDeEstudio(logger, bucketName, obj.Key)
			if err == nil {
				ch <- plan
			}

			return err
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	close(ch)

	planes := make([]plan, 0, len(ch))
	for p := range ch {
		planes = append(planes, p)
	}

	return planes, nil
}

func serPlanDeEstudio(logger *log.Logger, bucket, objKey *string) (plan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logger.Debug("Obteniendo metadata del plan")

	objHead, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: bucket,
		Key:    objKey,
	})

	if err != nil {
		return plan{}, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logger.Debug("Obteniendo contenido del plan")

	obj, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: bucket,
		Key:    objKey,
	})

	if err != nil {
		return plan{}, err
	}

	defer obj.Body.Close()

	data, err := io.ReadAll(obj.Body)
	if err != nil {
		return plan{}, err
	}

	materias := []materia{}

	err = json.Unmarshal(data, &materias)
	if err != nil {
		return plan{}, err
	}

	logger.Info(fmt.Sprintf("Obtenidas %v materias en el plan", len(materias)))

	// En este caso se sabe de antemano que si la información fue registrada y
	// leída correctamente, los datos son serializables.
	numero, _ := strconv.Atoi(objHead.Metadata["cuatri-numero"])
	anio, _ := strconv.Atoi(objHead.Metadata["cuatri-anio"])

	plan := plan{
		carrera: objHead.Metadata["carrera"],
		cuatri: cuatri{
			numero: numero,
			anio:   anio,
		},
		materias: materias,
	}

	return plan, nil
}

func (c cuatri) esDespuesDe(otro cuatri) bool {
	return (otro.anio < c.anio) ||
		((otro.anio == c.anio) && (otro.numero < c.numero))
}
