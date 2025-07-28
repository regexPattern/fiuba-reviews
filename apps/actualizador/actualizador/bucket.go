package actualizador

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

var clienteS3 *s3.Client

func (i *IndexadorOfertas) initS3Client(s3Ctx context.Context) error {
	s3Ctx, cancel := context.WithTimeout(s3Ctx, i.S3InitTimeout)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(s3Ctx)
	if err != nil {
		slog.Error("error inicializado cliente de S3", "error", err)
		return err
	}

	clienteS3 = s3.NewFromConfig(cfg)
	slog.Info("cliente de S3 inicializado exitosamente")

	return nil
}

// Oferta de una carrera.
type Oferta struct {
	Materias []MateriaSiu
	cuatri
	carrera string
}

// Materia de una oferta.
type MateriaSiu struct {
	Codigo   string       `json:"codigo"`
	Nombre   string       `json:"nombre"`
	Catedras []CatedraSiu `json:"catedras"`
}

// Cátedra de una materia.
type CatedraSiu struct {
	Codigo   int          `json:"codigo"`
	Docentes []DocenteSiu `json:"docentes"`
}

// Cátedra de una cátedra.
type DocenteSiu struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

// Cuatrimestre.
type cuatri struct {
	numero int
	anio   int
}

// newCuatri crea un cuatrimestre a partir de un hashmap con los campos 'cuatri-numero' y
func newCuatri(numero, anio string) (cuatri, error) {
	var c cuatri
	var err error
	c.numero, err = strconv.Atoi(numero)
	if err != nil {
		return c, err
	}
	c.anio, err = strconv.Atoi(anio)
	if err != nil {
		return c, err
	}
	return c, nil
}

// despuesDe compara si el cuatrimestre viene después del otro en orden cronológico.
func (c cuatri) despuesDe(otro cuatri) bool {
	if c.anio == otro.anio {
		return c.numero > otro.numero
	} else {
		return c.anio > otro.anio
	}
}

// obtenerOfertasCarreras obtiene las ofertas de carreras disponibles en el bucket.
func (i *IndexadorOfertas) obtenerOfertasCarreras(ctx context.Context) ([]*Oferta, error) {
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

// descargarObjetosBucket descarga los archivos del bucket.
func (i *IndexadorOfertas) descargarObjetosBucket(ctx context.Context) ([]s3types.Object, error) {
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

// newOfertaCarrera crea una nueva oferta de carrera a partir de un archivo del bucket.
func (i *IndexadorOfertas) newOfertaCarrera(
	ctx context.Context,
	objKey *string,
) (*Oferta, error) {
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
		"cuatri", cuatri.numero,
		"anio", cuatri.anio,
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

	o := &Oferta{
		Materias: materias,
		cuatri:   cuatri,
		carrera:  carrera,
	}

	logger.Info("oferta procesada exitosamente")

	return o, nil
}
