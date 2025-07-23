package patch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"slices"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"
)

var s3Client *s3.Client

type Generador struct {
	DbUrl         string
	DbTimeout     time.Duration
	S3BucketName  string
	S3InitTimeout time.Duration
	S3Timeout     time.Duration
}

func (g *Generador) GenerarPatches(ctx context.Context) (*PatchProposal, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	var materiasDb []materiaDb
	var oMaterias []ofertaMateriaSiu

	eg.Go(func() error {
		dbCtx, dbCancel := context.WithTimeout(egCtx, g.DbTimeout)
		defer dbCancel()

		var err error
		materiasDb, err = g.obtenerMateriasDB(dbCtx)

		return err
	})

	eg.Go(func() error {
		s3ctx, s3cancel := context.WithTimeout(egCtx, g.S3Timeout)
		defer s3cancel()

		if err := g.initS3Client(s3ctx); err != nil {
			return err
		}

		oCarreras, err := g.obtenerOfertasCarrerasSiu(s3ctx)
		if err != nil {
			return err
		}

		oMaterias = unificarOfertasMateriasSiu(oCarreras)

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	fmt.Println(materiasDb)
	fmt.Println()
	fmt.Println(oMaterias)

	return nil, nil
}

type materiaDb struct {
	Codigo string `db:"codigo"`
	Nombre string `db:"nombre"`
}

func (g *Generador) obtenerMateriasDB(ctx context.Context) ([]materiaDb, error) {
	conn, err := pgx.Connect(ctx, g.DbUrl)
	if err != nil {
		slog.Error("no se pudo conectar con la base de datos", "error", err)
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

	materias, err := pgx.CollectRows(rows, pgx.RowToStructByName[materiaDb])
	if err != nil {
		slog.Error("no se pudieron obtener las materias desactualizadas de la base de datos", "error", err)
		return nil, err
	}

	slog.Info(fmt.Sprintf("obtenidas %v materias desactualizadas de la base de datos", len(materias)))

	return materias, nil
}

func (g *Generador) initS3Client(ctx context.Context) error {
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
	cuatri
	materias []materiaSiu
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

func tryNewCuatri(metadata map[string]string) (cuatri, error) {
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

func (g *Generador) obtenerOfertasCarrerasSiu(ctx context.Context) ([]*ofertaCarrera, error) {
	objs, err := g.descargarOfertasBucket(ctx)
	if err != nil {
		return nil, err
	}

	ofertas := make([]*ofertaCarrera, 0, len(objs))

	for _, obj := range objs {
		if o, err := g.obtenerMateriasOferta(ctx, obj); err != nil {
			slog.Warn("omitiendo indexado de oferta", "key", *obj.Key)
		} else {
			ofertas = append(ofertas, o)
		}
	}

	return ofertas, nil
}

func (g *Generador) descargarOfertasBucket(ctx context.Context) ([]s3types.Object, error) {
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

func (g *Generador) obtenerMateriasOferta(ctx context.Context, obj s3types.Object) (*ofertaCarrera, error) {
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

	cuatri, err := tryNewCuatri(content.Metadata)
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

	logger.Debug(fmt.Sprintf("obtenidos %v bytes de contenido de oferta", len(bytes)))

	var materias []materiaSiu
	if err := json.Unmarshal(bytes, &materias); err != nil {
		logger.Error("error serializando contenido de oferta", "error", err)
		return nil, err
	}

	logger.Info("oferta procesada correctamente")

	o := &ofertaCarrera{
		cuatri:   cuatri,
		materias: materias,
	}

	return o, nil
}

type ofertaMateriaSiu struct {
	cuatri
	materia materiaSiu
}

func unificarOfertasMateriasSiu(oCarreras []*ofertaCarrera) []ofertaMateriaSiu {
	totalMaterias := 0

	for _, oc := range oCarreras {
		totalMaterias += len(oc.materias)
	}

	oMaterias := make(map[string]ofertaMateriaSiu, totalMaterias)

	for _, oc := range oCarreras {
		for _, m := range oc.materias {
			oMasReciente, ok := oMaterias[m.Nombre]
			if !ok || oc.cuatri.despuesDe(oMasReciente.cuatri) {
				oMaterias[m.Nombre] = ofertaMateriaSiu{
					cuatri:  oc.cuatri,
					materia: m,
				}
			}
		}
	}

	slog.Debug("unificado últimas ofertas de materias", "n_inicial", totalMaterias, "n_final", len(oMaterias))

	return slices.Collect(maps.Values(oMaterias))
}
