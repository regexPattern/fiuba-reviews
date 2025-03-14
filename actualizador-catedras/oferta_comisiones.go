package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/log"
	"golang.org/x/sync/errgroup"
)

type ofertaComisiones struct {
	carrera  string
	cuatri   cuatri
	materias []materia
}

type cuatri struct {
	numero int
	anio   int
}

func (c cuatri) esDespuesDe(otro cuatri) bool {
	return (otro.anio < c.anio) ||
		((otro.anio == c.anio) && (otro.numero < c.numero))
}

// getOfertasComisiones retorna todas las ofertas de comisiones almacenadas en
// el bucket. Hay, a lo sumo, una oferta por carrera, que siempre es la variante
// más reciente disponible de dicha oferta.
func getOfertasComisiones() ([]ofertaComisiones, error) {
	logger := log.Default().WithPrefix("🪣")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logger.Debug("listando ofertas de comisiones", "bucket", *bucketKey)

	bucket, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: bucketKey,
	})

	if err != nil {
		log.Error(err)
		return nil, errors.New("error listando archivos de la cubeta")
	}

	logger.Info(fmt.Sprintf("encontradas %v ofertas de comisiones", len(bucket.Contents)), "bucket", *bucketKey)

	var eg errgroup.Group
	eg.SetLimit(MAX_REQ_CONCURRENTES)

	ch := make(chan ofertaComisiones, len(bucket.Contents))

	for _, obj := range bucket.Contents {
		eg.Go(func() error {
			return newOfertaComisiones(ch, obj.Key)
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, errors.New("error obteniendo ofertas de comisiones")
	}

	close(ch)

	ofertas := make([]ofertaComisiones, 0, len(ch))
	for o := range ch {
		ofertas = append(ofertas, o)
	}

	return ofertas, nil
}

// newOfertaComisiones serializa una oferta de comisión a partir de un archivo
// almacenado en el bucket.
func newOfertaComisiones(ch chan ofertaComisiones, objKey *string) error {
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
			logger.Warn("materia sin cátedras", "materia", m.Codigo)
		} else {
			materiasConCatedras = append(materiasConCatedras, m)
		}
	}

	ch <- ofertaComisiones{
		carrera: objHead.Metadata["carrera"],
		cuatri: cuatri{
			numero: numero,
			anio:   anio,
		},
		materias: materiasConCatedras,
	}

	return nil
}
