package patcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var s3Client *s3.Client

type GeneradorPatches struct {
	DbUrl         string
	DbTimeout     time.Duration
	S3BucketName  string
	S3InitTimeout time.Duration
	S3Timeout     time.Duration
}

type PatchGenerado struct {
	CodigoDb  string
	CodigoSiu string
	Nombre    string
	Catedras  []CatedraSiu
	Cuatri
}

func (g *GeneradorPatches) GenerarPatches(ctx context.Context) ([]PatchGenerado, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	var materiasDb []materiaDb
	var oMateriasSiu []ofertaMateriaSiu

	eg.Go(func() error {
		dbCtx, dbCancel := context.WithTimeout(egCtx, g.DbTimeout)
		defer dbCancel()

		var err error
		materiasDb, err = g.obtenerMateriasDb(dbCtx)
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

		oMateriasSiu = unificarOfertasMateriasSiu(oCarreras)
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	patches := vincularMateriasSiuConDb(oMateriasSiu, materiasDb)
	return patches, nil
}

type materiaDb struct {
	Codigo string `db:"codigo"`
	Nombre string `db:"nombre"`
}

func (g *GeneradorPatches) obtenerMateriasDb(ctx context.Context) ([]materiaDb, error) {
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
WHERE (
	codigo_cuatrimestre_actualizacion IS NULL OR
	codigo_cuatrimestre_actualizacion < (SELECT MAX(codigo) FROM cuatrimestre)
) AND EXISTS (
	SELECT 1
	FROM plan_materia
	JOIN plan ON plan.codigo = plan_materia.codigo_plan
	WHERE plan_materia.codigo_materia = materia.codigo AND plan.esta_vigente
);
		`)

	materias, err := pgx.CollectRows(rows, pgx.RowToStructByName[materiaDb])
	if err != nil {
		slog.Error("no se pudieron obtener las materias desactualizadas de la base de datos", "error", err)
		return nil, err
	}

	slog.Info(fmt.Sprintf("obtenidas %v materias desactualizadas de la base de datos", len(materias)))

	return materias, nil
}

func (g *GeneradorPatches) initS3Client(ctx context.Context) error {
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
	Cuatri
	materias []materiaSiu
}

type materiaSiu struct {
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

func tryNewCuatri(metadata map[string]string) (Cuatri, error) {
	var c Cuatri

	numero, err := strconv.Atoi(metadata["cuatri-numero"])
	if err != nil {
		return c, err
	}

	anio, err := strconv.Atoi(metadata["cuatri-anio"])
	if err != nil {
		return c, err
	}

	c.Numero = numero
	c.Anio = anio

	return c, nil
}

func (c Cuatri) despuesDe(otro Cuatri) bool {
	if c.Anio == otro.Anio {
		return c.Numero > otro.Numero
	} else {
		return c.Anio > otro.Anio
	}
}

func (g *GeneradorPatches) obtenerOfertasCarrerasSiu(ctx context.Context) ([]*ofertaCarrera, error) {
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

func (g *GeneradorPatches) descargarOfertasBucket(ctx context.Context) ([]s3types.Object, error) {
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

func (g *GeneradorPatches) obtenerMateriasOferta(ctx context.Context, obj s3types.Object) (*ofertaCarrera, error) {
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

	logger = slog.Default().With("carrera", carrera, "cuatri", cuatri.Numero, "anio", cuatri.Anio)

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

	o := &ofertaCarrera{
		Cuatri:   cuatri,
		materias: materias,
	}

	logger.Info("oferta procesada exitosamente")

	return o, nil
}

type ofertaMateriaSiu struct {
	Cuatri
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
			if !ok || oc.Cuatri.despuesDe(oMasReciente.Cuatri) {
				oMaterias[m.Nombre] = ofertaMateriaSiu{
					Cuatri:  oc.Cuatri,
					materia: m,
				}
			}
		}
	}

	slog.Debug("unificado las últimas ofertas de materias", "n_inicial", totalMaterias, "n_final", len(oMaterias))

	return slices.Collect(maps.Values(oMaterias))
}

func vincularMateriasSiuConDb(materiasSiu []ofertaMateriaSiu, materiasDb []materiaDb) []PatchGenerado {
	materiasDbMap := make(map[string]materiaDb, len(materiasDb))
	for _, mDb := range materiasDb {
		materiasDbMap[normalize(mDb.Nombre)] = mDb
	}

	patches := make([]PatchGenerado, 0, len(materiasSiu))
	for _, mSiu := range materiasSiu {
		if mDb, ok := materiasDbMap[normalize(mSiu.materia.Nombre)]; ok {
			patches = append(patches, PatchGenerado{
				CodigoDb:  mDb.Codigo,
				CodigoSiu: mSiu.materia.Codigo,
				Nombre:    mDb.Nombre,
				Catedras:  mSiu.materia.Catedras,
				Cuatri:    mSiu.Cuatri,
			})
		} else {
			slog.Warn(
				"materia del SIU no encontrada en la base de datos",
				"nombre", mSiu.materia.Nombre,
				"codigo", mSiu.materia.Codigo,
			)
		}
	}

	return patches
}

func normalize(s string) string {
	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)

	result, _, _ := transform.String(t, strings.ToLower(s))
	return result
}
