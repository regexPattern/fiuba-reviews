package patch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var clienteS3 *s3.Client

func (i *Indexador) configClienteS3(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, i.S3InitTimeout)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		slog.Error("error inicializado cliente de S3", "error", err)
		return err
	}
	clienteS3 = s3.NewFromConfig(cfg)
	slog.Info("cliente de S3 inicializado exitosamente")
	return nil
}

func (i *Indexador) obtenerOfertasCarreras(ctx context.Context) ([]*Oferta, error) {
	ctx, cancel := context.WithTimeout(ctx, i.S3OpsTimeout)
	defer cancel()

	objs, err := i.descargarObjetosBucket(ctx)
	if err != nil {
		return nil, err
	}

	ofertas := make([]*Oferta, 0, len(objs))
	for _, obj := range objs {
		if o, err := i.newOfertaCarrera(ctx, obj.Key); err != nil {
			slog.Warn("omitiendo indexado de oferta", "key", *obj.Key)
		} else {
			ofertas = append(ofertas, o)
		}
	}

	return ofertas, nil
}

func (i *Indexador) descargarObjetosBucket(ctx context.Context) ([]s3types.Object, error) {
	output, err := clienteS3.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &i.S3BucketName,
	})
	if err != nil {
		slog.Error("error enlistando archivos del bucket", "error", err)
		return nil, err
	}
	slog.Debug(fmt.Sprintf("obtenidas %v ofertas del bucket", len(output.Contents)))
	return output.Contents, nil
}

func (i *Indexador) newOfertaCarrera(ctx context.Context, objKey *string) (*Oferta, error) {
	logger := slog.Default().With("key", *objKey)

	obj, err := clienteS3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &i.S3BucketName,
		Key:    objKey,
	})

	if err != nil {
		logger.Error("error obteniendo contenido del objeto", "error", err)
		return nil, err
	}

	carrera := obj.Metadata["carrera"]

	cuatri, err := newCuatri(obj.Metadata["cuatri-numero"], obj.Metadata["cuatri-anio"])
	if err != nil {
		logger.Error("error obteniendo cuatrimestre de la oferta",
			"carrera", carrera,
			"error", err,
		)
		return nil, err
	}

	logger = slog.Default().With(
		"carrera", carrera,
		"cuatri", cuatri.Numero,
		"anio", cuatri.Anio,
	)

	defer obj.Body.Close()
	bytes, err := io.ReadAll(obj.Body)
	if err != nil {
		logger.Error("error leyendo bytes de contenido de oferta", "error", err)
		return nil, err
	}

	var materias []MateriaSiu
	if err := json.Unmarshal(bytes, &materias); err != nil {
		logger.Error("error serializando contenido de oferta", "error", err)
		return nil, err
	}

	for _, m := range materias {
		if len(m.Catedras) == 0 {
			logger.Warn("materia no tiene cátedras",
				"codigo", m.Codigo,
				"nombre", m.Nombre)
		}
	}

	logger.Info("oferta procesada exitosamente")

	o := &Oferta{
		Materias: materias,
		Cuatri:   cuatri,
		Carrera:  carrera,
	}

	return o, nil
}
