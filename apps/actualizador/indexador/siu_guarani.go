package indexador

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"slices"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var clienteS3 *s3.Client

type OfertaMateriaSiu struct {
	MateriaSiu
	Cuatri
}

type OfertaCarreraSiu struct {
	Carrera  string
	Materias []MateriaSiu
	Cuatri
}

type MateriaSiu struct {
	Codigo   string       `json:"codigo"`
	Nombre   string       `json:"nombre"`
	Catedras []CatedraSiu `json:"catedras"`
}

type CatedraSiu struct {
	Codigo   int          `json:"codigo"`
	Docentes []DocenteSiu `json:"docentes"`
}

type DocenteSiu struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

type Cuatri struct {
	Numero int
	Anio   int
}

func newCuatri(numero, anio string) (Cuatri, error) {
	var c Cuatri
	var err error
	c.Numero, err = strconv.Atoi(numero)
	if err != nil {
		return c, err
	}
	c.Anio, err = strconv.Atoi(anio)
	if err != nil {
		return c, err
	}
	return c, nil
}

func (c Cuatri) despuesDe(otro Cuatri) bool {
	if c.Anio == otro.Anio {
		return c.Numero > otro.Numero
	} else {
		return c.Anio > otro.Anio
	}
}

func (i *Indexador) initClienteS3(ctx context.Context) error {
	initctx, initcancel := context.WithTimeout(ctx, i.S3InitTimeout)
	defer initcancel()

	cfg, err := config.LoadDefaultConfig(initctx)
	if err != nil {
		slog.Error("error configurando cliente de S3", "error", err)
		return err
	}

	clienteS3 = s3.NewFromConfig(cfg)

	slog.Info("cliente de S3 configurado exitosamente")

	return nil
}

func (i *Indexador) obtenerOfertasSiu(ctx context.Context) ([]OfertaMateriaSiu, error) {
	objs, err := i.fetchBucketObjects(ctx)
	if err != nil {
		return nil, err
	}

	ofertas := make([]OfertaCarreraSiu, 0, len(objs))
	for _, obj := range objs {
		if o, err := i.newOfertaCarreraSiu(ctx, obj.Key); err != nil {
			slog.Warn("omitiendo indexado de oferta", "key", *obj.Key)
		} else {
			ofertas = append(ofertas, o)
		}
	}

	return unificarOfertasSiu(ofertas), nil
}

func (i *Indexador) fetchBucketObjects(ctx context.Context) ([]s3types.Object, error) {
	opctx, opcancel := context.WithTimeout(ctx, i.S3OpTimeout)
	defer opcancel()

	output, err := clienteS3.ListObjectsV2(opctx, &s3.ListObjectsV2Input{
		Bucket: &i.S3BucketName,
	})
	if err != nil {
		slog.Error("error enlistando archivos del bucket", "error", err)
		return nil, err
	}

	slog.Debug(fmt.Sprintf("encontradas %v ofertas de carrera en el bucket", len(output.Contents)))

	return output.Contents, nil
}

func (i *Indexador) newOfertaCarreraSiu(
	ctx context.Context,
	objKey *string,
) (OfertaCarreraSiu, error) {
	var oferta OfertaCarreraSiu

	l := slog.Default().With("key", *objKey)

	opctx, opcancel := context.WithTimeout(ctx, i.S3OpTimeout)
	defer opcancel()

	obj, err := clienteS3.GetObject(opctx, &s3.GetObjectInput{
		Bucket: &i.S3BucketName,
		Key:    objKey,
	})
	if err != nil {
		l.Error("error obteniendo contenido del objeto", "error", err)
		return oferta, err
	}

	carrera := obj.Metadata["carrera"]

	cuatri, err := newCuatri(obj.Metadata["cuatri-numero"], obj.Metadata["cuatri-anio"])
	if err != nil {
		l.Error("error obteniendo cuatrimestre de la oferta",
			"carrera", carrera,
			"error", err,
		)
		return oferta, err
	}

	l = slog.Default().With(
		"carrera", carrera,
		"cuatri", cuatri.Numero,
		"anio", cuatri.Anio,
	)

	defer func() { _ = obj.Body.Close() }()
	bytes, err := io.ReadAll(obj.Body)
	if err != nil {
		l.Error("error leyendo bytes de contenido de oferta", "error", err)
		return oferta, err
	}

	var materias []MateriaSiu
	if err := json.Unmarshal(bytes, &materias); err != nil {
		l.Error("error serializando contenido de oferta", "error", err)
		return oferta, err
	}

	for _, m := range materias {
		if len(m.Catedras) == 0 {
			l.Warn("materia no tiene cátedras", "codigo_siu", m.Codigo, "nombre", m.Nombre)
		}
	}

	l.Info("oferta de carrera procesada exitosamente")

	oferta = OfertaCarreraSiu{
		Materias: materias,
		Cuatri:   cuatri,
		Carrera:  carrera,
	}

	return oferta, nil
}

func unificarOfertasSiu(ofertas []OfertaCarreraSiu) []OfertaMateriaSiu {
	filtradas := make(map[string]OfertaMateriaSiu)

	for _, o := range ofertas {
		for _, m := range o.Materias {
			if ofertaActual, ok := filtradas[m.Nombre]; !ok ||
				o.despuesDe(ofertaActual.Cuatri) {
				filtradas[m.Nombre] = OfertaMateriaSiu{
					MateriaSiu: m,
					Cuatri:     o.Cuatri,
				}
			} else if ok && ofertaActual.Cuatri == o.Cuatri {
				// Si tenemos dos ofertas para una misma materia, y ambas
				// ofertas corresponden al mismo cuatrimestre, entonces
				// agregamos todas las cátedras de ambas ofertas al total de
				// la materia.

				catedras := make(map[int]CatedraSiu)
				for _, c := range ofertaActual.Catedras {
					catedras[c.Codigo] = c
				}
				for _, c := range m.Catedras {
					catedras[c.Codigo] = c
				}

				filtradas[m.Nombre] = OfertaMateriaSiu{
					MateriaSiu: MateriaSiu{
						Codigo:   ofertaActual.Codigo,
						Nombre:   ofertaActual.Nombre,
						Catedras: slices.Collect(maps.Values(catedras)),
					},
					Cuatri: o.Cuatri,
				}

				slog.Debug("unificadas cátedras de materia",
					"codigo_siu", m.Codigo, "nombre", m.Nombre, "cuatri", o.Cuatri.Numero, "anio", o.Cuatri.Anio)
			}
		}
	}

	slog.Info(fmt.Sprintf("extraído ofertas de %v materias", len(filtradas)))

	return slices.Collect(maps.Values(filtradas))
}
