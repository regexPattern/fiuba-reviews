package patcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"os"
	"slices"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/jackc/pgx/v5"
)

var client *s3.Client
var bucketName *string
var dbUrl string

type Oferta struct {
	Materias []Materia
	Cuatri
}

type Cuatri struct {
	Numero int
	Anio   int
}

func (c Cuatri) Compare(other Cuatri) int {
	if c.Anio > other.Anio {
		return 1
	}
	if c.Anio < other.Anio {
		return -1
	}
	if c.Numero > other.Numero {
		return 1
	}
	if c.Numero < other.Numero {
		return -1
	}
	return 0
}

type Materia struct {
	Codigo   string    `json:"codigo" db:"codigo"`
	Nombre   string    `json:"nombre" db:"nombre"`
	Catedras []Catedra `json:"catedras"`
}

type Catedra struct {
	Codigo   int       `json:"codigo"`
	Docentes []Docente `json:"docentes"`
}

type Docente struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

type LatestMateria struct {
	Materia
	Cuatri
}

func initS3Client(ctx context.Context) error {
	bucketNameInput := "fiuba-reviews"
	bucketName = &bucketNameInput

	dbUrl = os.Getenv("DATABASE_URL")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client = s3.NewFromConfig(cfg)

	slog.Debug("cliente de S3 se ha creado exitosamente")

	return nil
}

func getOfertasFromSiu(ctx context.Context) ([]*Oferta, error) {
	objs, err := fetchObjectsFromBucket(ctx)
	if err != nil {
		return nil, err
	}

	ofertas := make([]*Oferta, 0, len(objs))

	for _, obj := range objs {
		if of, err := newOferta(ctx, obj); err != nil {
			slog.Warn("omitiendo indexado de oferta", "key", *obj.Key)
		} else {
			ofertas = append(ofertas, of)
		}
	}

	slog.Debug(fmt.Sprintf("obtenidas %v ofertas del bucket", len(ofertas)))

	return ofertas, nil
}

func fetchObjectsFromBucket(ctx context.Context) ([]s3types.Object, error) {
	output, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: bucketName,
	})
	if err != nil {
		slog.Error("no se pudieron enlistar los archivos del bucket", "error", err)
		return nil, err
	}

	slog.Info(fmt.Sprintf("obtenidos %v archivos del bucket", len(output.Contents)))

	return output.Contents, nil
}

func newOferta(ctx context.Context, obj s3types.Object) (*Oferta, error) {
	logger := slog.Default().With("key", *obj.Key)

	content, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: bucketName,
		Key:    obj.Key,
	})

	if err != nil {
		logger.Error("error obteniendo contenido de objeto", "error", err)
		return nil, err
	}

	logger = slog.Default().With(
		"carrera", content.Metadata["carrera"],
		"cuatri", content.Metadata["cuatri-numero"],
		"anio", content.Metadata["cuatri-anio"])

	defer content.Body.Close()
	bytes, err := io.ReadAll(content.Body)
	if err != nil {
		logger.Error("error leyendo bytes de contenido de oferta", "error", err)
		return nil, err
	}

	logger.Debug(fmt.Sprintf("obtenido %v bytes de contenido de oferta", len(bytes)))

	var materias []Materia
	if err := json.Unmarshal(bytes, &materias); err != nil {
		logger.Error("error serializando contenido de oferta", "error", err)
		return nil, err
	}

	o := Oferta{
		Materias: materias,
	}

	return &o, nil
}

func mergeLatestOfertas(ofertas []*Oferta) []LatestMateria {
	merged := make(map[string]LatestMateria, len(ofertas))

	for _, o := range ofertas {
		for _, m := range o.Materias {
			curr, ok := merged[m.Nombre]
			if !ok || o.Cuatri.Compare(curr.Cuatri) == 1 {
				merged[m.Nombre] = LatestMateria{
					Materia: m,
					Cuatri:  o.Cuatri,
				}
			}
		}
	}

	return slices.Collect(maps.Values(merged))
}

func getMateriasFromDb(ctx context.Context) ([]Materia, error) {
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		slog.Error("no se puedo conectar con la base de datos", "error", err)
		return nil, err
	}

	defer conn.Close(ctx)

	rows, _ := conn.Query(ctx, `
SELECT
	codigo,
	nombre
FROM
	materia
WHERE
	codigo_cuatrimestre_actualizacion IS NULL OR
	codigo_cuatrimestre_actualizacion < (SELECT MAX(codigo) FROM cuatrimestre);
		`)

	materias, err := pgx.CollectRows(rows, pgx.RowTo[Materia])
	if err != nil {
		slog.Error("no se pudieron obtener materias desactualizadas de la base de datos", "error", err)
		return nil, err
	}

	return materias, nil
}
