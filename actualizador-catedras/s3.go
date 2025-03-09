package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/log"
	"golang.org/x/sync/errgroup"
)

const MAX_REQ_CONCURRENTES int = 5

var client *s3.Client
var bucketKey *string

func initS3Client(logger *log.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		logger.Error(err)
		return errors.New("error cargando la configuración de AWS")
	}

	logger.Debug("cargada la configuración de AWS")

	client = s3.NewFromConfig(cfg)

	bucketKeyEnv, ok := os.LookupEnv("AWS_S3_BUCKET")
	if !ok {
		return errors.New("variable de entorno `AWS_S3_BUCKET` no configurada")
	}

	bucketKey = aws.String(bucketKeyEnv)

	return nil
}

// getOfertasComisiones retorna todas las ofertas de comisiones almacenadas en
// el bucket. Hay a lo sumo una oferta por carrera, que siempre es la variante
// más reciente disponible de dicha oferta.
func getOfertasComisiones() ([]oferta, error) {
	logger := log.Default().WithPrefix("🪣")

	logger.Info("obteniendo ofertas de comisiones")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logger.Debug("listando archivos de la cubeta", "bucket", *bucketKey)

	bucket, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: bucketKey,
	})

	if err != nil {
		log.Error(err)
		return nil, errors.New("error listando archivos de la cubeta")
	}

	logger.Info(fmt.Sprintf("encontrados %v archivos de ofertas de comisiones", len(bucket.Contents)), "bucket", *bucketKey)

	var eg errgroup.Group
	eg.SetLimit(MAX_REQ_CONCURRENTES)

	ch := make(chan oferta, len(bucket.Contents))

	for _, obj := range bucket.Contents {
		eg.Go(func() error {
			return serOfertaComisiones(ch, obj.Key)
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, errors.New("error obteniendo ofertas de comisiones")
	}

	close(ch)

	ofertas := make([]oferta, 0, len(ch))
	for o := range ch {
		ofertas = append(ofertas, o)
	}

	return ofertas, nil
}

// serOfertaComisiones serializa una oferta de comisión a partir de un archivo
// almacenado en el bucket.
func serOfertaComisiones(ch chan oferta, objKey *string) error {
	logger := log.Default().WithPrefix("📄").With("objKey", *objKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logger.Debug("obteniendo metadata")

	objHead, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: bucketKey,
		Key:    objKey,
	})

	if err != nil {
		logger.Error("error obteniendo metadata", "error", err)
		return err
	}

	// En este caso se sabe de antemano que si la información fue registrada y
	// leída correctamente, los datos son serializables.
	numero, _ := strconv.Atoi(objHead.Metadata["cuatri-numero"])
	anio, _ := strconv.Atoi(objHead.Metadata["cuatri-anio"])

	logger = logger.With("cuatri", numero, "anio", anio)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logger.Debug("obteniendo contenido")

	obj, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: bucketKey,
		Key:    objKey,
	})

	if err != nil {
		logger.Error("error obteniendo contenido", "error", err)
		return err
	}

	defer obj.Body.Close()

	data, err := io.ReadAll(obj.Body)
	if err != nil {
		logger.Error("error leyendo contenido", "error", err)
		return err
	}

	materias := []materia{}

	err = json.Unmarshal(data, &materias)
	if err != nil {
		logger.Error("error serializando oferta de comisiones", "error", err)
		return err
	}

	logger.Infof("encontradas %v materias en oferta de comisiones", len(materias))

	materiasConCatedras := make([]materia, 0, len(materias))
	for _, m := range materias {
		if len(m.Catedras) == 0 {
			logger.Warn("materia sin cátedras", "codigoMateria", m.Codigo)
		} else {
			materiasConCatedras = append(materiasConCatedras, m)
		}
	}

	ch <- oferta{
		carrera: objHead.Metadata["carrera"],
		cuatri: cuatri{
			numero: numero,
			anio:   anio,
		},
		materias: materiasConCatedras,
	}

	return nil
}
