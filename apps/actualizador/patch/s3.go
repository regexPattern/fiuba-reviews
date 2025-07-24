package patch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var s3Client *s3.Client

func (g *GeneradorPatches) initClienteS3(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, g.S3InitTimeout)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		slog.Error("no se puedo configurar cliente de S3", "error", err)
		return err
	}

	s3Client = s3.NewFromConfig(cfg)
	slog.Info("cliente de S3 inicializado exitosamente")

	return nil
}

type ofertaCarrera struct {
	Materias []materiaSiu
	cuatri
	carrera string
}

type materiaSiu struct {
	Codigo   string       `json:"codigo"`
	Nombre   string       `json:"nombre"`
	Catedras []catedraSiu `json:"catedras"`
}

type catedraSiu struct {
	Codigo   int          `json:"codigo"`
	Docentes []docenteSiu `json:"docentes"`
}

type docenteSiu struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

type cuatri struct {
	numero int
	anio   int
}

func newCuatri(metadata map[string]string) (cuatri, error) {
	var c cuatri
	numero, err := strconv.Atoi(metadata["cuatri-numero"])
	if err != nil {
		return c, err
	}
	anio, err := strconv.Atoi(metadata["cuatri-anio"])
	if err != nil {
		return c, err
	}
	c.numero = numero
	c.anio = anio
	return c, nil
}

func (c cuatri) despuesDe(otro cuatri) bool {
	if c.anio == otro.anio {
		return c.numero > otro.numero
	} else {
		return c.anio > otro.anio
	}
}

func (g *GeneradorPatches) obtenerOfertasCarreras(ctx context.Context) ([]*ofertaCarrera, error) {
	ctx, cancel := context.WithTimeout(ctx, g.S3Timeout)
	defer cancel()

	objs, err := g.descargarObjetosBucket(ctx)
	if err != nil {
		return nil, err
	}

	ofertas := make([]*ofertaCarrera, 0, len(objs))
	for _, obj := range objs {
		if o, err := g.newOfertaCarrera(ctx, obj); err != nil {
			slog.Warn("omitiendo indexado de oferta", "key", *obj.Key)
		} else {
			ofertas = append(ofertas, o)
		}
	}

	return ofertas, nil
}

func (g *GeneradorPatches) descargarObjetosBucket(ctx context.Context) ([]s3types.Object, error) {
	output, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &g.S3BucketName,
	})

	if err != nil {
		slog.Error("no se pudieron enlistar los archivos del bucket", "error", err)
		return nil, err
	}

	slog.Debug(fmt.Sprintf("obtenidas %v ofertas del bucket", len(output.Contents)))

	return output.Contents, nil
}

func (g *GeneradorPatches) newOfertaCarrera(ctx context.Context, obj s3types.Object) (*ofertaCarrera, error) {
	logger := slog.Default().With("key", *obj.Key)

	content, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &g.S3BucketName,
		Key:    obj.Key,
	})

	if err != nil {
		logger.Error("error obteniendo el contenido del objeto", "error", err)
		return nil, err
	}

	carrera := content.Metadata["carrera"]
	cuatri, err := newCuatri(content.Metadata)
	if err != nil {
		logger.Error("error obteniendo el cuatrimestre de la oferta", "carrera", carrera, "error", err)
		return nil, err
	}

	logger = slog.Default().With("carrera", carrera, "cuatri", cuatri.numero, "anio", cuatri.anio)

	defer content.Body.Close()
	bytes, err := io.ReadAll(content.Body)
	if err != nil {
		logger.Error("error leyendo bytes de contenido de oferta", "error", err)
		return nil, err
	}

	var materias []materiaSiu
	if err := json.Unmarshal(bytes, &materias); err != nil {
		logger.Error("error serializando contenido de oferta", "error", err)
		return nil, err
	}

	for _, m := range materias {
		if len(m.Catedras) == 0 {
			logger.Warn("materia no tiene cátedras", "codigo", m.Codigo, "nombre", m.Nombre)
		}
	}

	o := &ofertaCarrera{
		Materias: materias,
		cuatri:   cuatri,
		carrera:  carrera,
	}

	logger.Info("oferta procesada exitosamente")

	return o, nil
}
