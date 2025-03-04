package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/log"
	"golang.org/x/sync/errgroup"
)

const BUCKET_NAME_ENV string = ""
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

func (c cuatri) esDespuesDe(otro cuatri) bool {
	return (otro.anio < c.anio) ||
		((otro.anio == c.anio) && (otro.numero < c.numero))
}

type materia struct {
	Codigo   string    `json:"codigo"`
	Nombre   string    `json:"nombre"`
	Catedras []catedra `json:"catedras"`
}

func (m *materia) Equal(otra *materia) bool {
	return m.Codigo == otra.Codigo && m.Nombre == otra.Nombre
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

func getUltimosPlanesDeEstudio() ([]plan, error) {
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

	planes := make(chan plan, len(output.Contents))

	for _, obj := range output.Contents {
		logger := logger.With("objKey", *obj.Key)

		eg.Go(func() error {
			plan, err := getPlanDeEstudio(logger, bucket, obj.Key)
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
	ultimosPlanes := make(map[string]plan, len(planes))

	for p1 := range planes {
		p2, ok := ultimosPlanes[p1.carrera]
		if !ok || p1.cuatri.esDespuesDe(p2.cuatri) {
			ultimosPlanes[p1.carrera] = p1
		}
	}

	return slices.Collect(maps.Values(ultimosPlanes)), nil
}

func getPlanDeEstudio(logger *log.Logger, bucket, objKey *string) (plan, error) {
	var plan plan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logger.Debug("Obteniendo metadata del plan")

	head, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: bucket,
		Key:    objKey,
	})

	if err != nil {
		return plan, nil
	}

	numero, err := strconv.Atoi(head.Metadata["cuatri-numero"])
	if err != nil {
		return plan, err
	}

	anio, err := strconv.Atoi(head.Metadata["cuatri-anio"])
	if err != nil {
		return plan, err
	}

	plan.carrera = head.Metadata["carrera"]
	plan.cuatri = cuatri{
		numero: numero,
		anio:   anio,
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logger.Debug("Obteniendo contenido del plan")

	obj, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: bucket,
		Key:    objKey,
	})

	if err != nil {
		return plan, err
	}

	defer obj.Body.Close()

	data, err := io.ReadAll(obj.Body)
	if err != nil {
		return plan, err
	}

	err = json.Unmarshal(data, &plan.materias)
	if err != nil {
		return plan, err
	}

	logger.Info(fmt.Sprintf("Obtenidas %v materias en el plan", len(plan.materias)))

	return plan, nil
}

func filtrarMateriasMasRecientes(planes []plan) []materia {
	maxMaterias := 0
	for _, plan := range planes {
		maxMaterias += len(plan.materias)
	}

	cuatris := make(map[string]cuatri, maxMaterias)
	materias := make(map[string]materia, maxMaterias)

	for _, plan := range planes {
		for _, materia := range plan.materias {
			cuatriUltimoCambio, ok := cuatris[materia.Nombre]

			if !ok || plan.cuatri.esDespuesDe(cuatriUltimoCambio) {
				cuatris[materia.Nombre] = plan.cuatri
				materias[materia.Nombre] = materia
			}
		}
	}

	return slices.Collect(maps.Values(materias))
}
