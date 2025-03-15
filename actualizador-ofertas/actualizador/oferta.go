package actualizador

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/charmbracelet/log"
	"github.com/regexPattern/fiuba-reviews/actualizador-ofertas/config"
	"golang.org/x/sync/errgroup"
)

type oferta struct {
	ofertaMetadata
	materias []materia
}

type ofertaMetadata struct {
	carrera string
	cuatri  cuatri
}

type cuatri struct {
	numero int
	anio   int
}

func getOfertas(logger *log.Logger, cfg *config.S3Config) ([]oferta, error) {
	logger = logger.With("bucketName", *cfg.BucketName)

	objs, err := getBucketObjects(logger, cfg)
	if err != nil {
		return nil, err
	}

	var eg errgroup.Group
	eg.SetLimit(config.BucketMaxRequests)
	ofertach := make(chan oferta, len(objs))

	for _, obj := range objs {
		eg.Go(func() error {
			return newOfertaFromObject(logger, cfg, ofertach, obj)
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	close(ofertach)
	ofs := make([]oferta, 0, len(ofertach))
	for o := range ofertach {
		ofs = append(ofs, o)
	}

	return ofs, nil
}

func getBucketObjects(logger *log.Logger, cfg *config.S3Config) ([]types.Object, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	b, err := cfg.Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: cfg.BucketName,
	})
	if err != nil {
		msg := "error obteniendo listado de objetos del bucket"
		logger.Error(msg, "err", err)
		return nil, errors.New(msg)
	}

	return b.Contents, nil
}

func newOfertaFromObject(logger *log.Logger, cfg *config.S3Config, ch chan oferta, obj types.Object) error {
	logger = logger.With("objKey", *obj.Key)

	objMd, err := getObjMetadata(logger, cfg, obj.Key)
	if err != nil {
		return err
	}
	body, err := getObjBody(logger, cfg, obj.Key)
	if err != nil {
		return err
	}
	md, err := newPlanMetadata(objMd)
	if err != nil {
		return err
	}

	logger = log.Default().With("carrera",
		md.carrera, "cuatri", md.cuatri.numero, "anio", md.cuatri.anio)

	of, err := newOferta(logger, md, body)
	if err != nil {
		msg := "error serializando oferta"
		return wrapErrorMsg(logger, msg, err)
	}

	ch <- of
	logger.Info("oferta obtenida exitosamente")

	return nil
}

func getObjMetadata(logger *log.Logger, cfg *config.S3Config, key *string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	h, err := cfg.Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: cfg.BucketName,
		Key:    key,
	})
	if err != nil {
		msg := "error obteniendo metadata de objeto"
		return nil, wrapErrorMsg(logger, msg, err)
	}

	return h.Metadata, nil
}

func getObjBody(logger *log.Logger, cfg *config.S3Config, key *string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	obj, err := cfg.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: cfg.BucketName,
		Key:    key,
	})
	if err != nil {
		msg := "error obteniendo body del objeto"
		return nil, wrapErrorMsg(logger, msg, err)
	}

	defer obj.Body.Close()
	bs, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func newPlanMetadata(objMd map[string]string) (ofertaMetadata, error) {
	var md ofertaMetadata

	getMd := func(key string) (string, error) {
		if val, ok := objMd[key]; !ok {
			return "", fmt.Errorf("campo '%v' no está presente en metadata", key)
		} else {
			return val, nil
		}
	}
	errMd := func(err error) error {
		return fmt.Errorf("error procesando metadata del objeto: %w", err)
	}

	carr, err := getMd("carrera")
	if err != nil {
		return md, err
	}

	getValFromMd := func(key string) (int, error) {
		str, err := getMd(key)
		if err != nil {
			return 0, errMd(fmt.Errorf("campo `%v` no está presente en metadata", key))
		}
		val, err := strconv.Atoi(str)
		if err != nil {
			return 0, errMd(fmt.Errorf("campo `%v` no serializable como entero: %w", key, err))
		}
		return val, nil
	}

	num, err := getValFromMd("cuatri-numero")
	if err != nil {
		return md, err
	}
	anio, err := getValFromMd("cuatri-anio")
	if err != nil {
		return md, err
	}

	md = ofertaMetadata{
		carrera: carr,
		cuatri: cuatri{
			anio:   anio,
			numero: num,
		},
	}

	return md, nil
}

func newOferta(logger *log.Logger, md ofertaMetadata, body []byte) (oferta, error) {
	var of oferta

	mats := []materia{}
	if err := json.Unmarshal(body, &mats); err != nil {
		msg := "error serializando body del objeto"
		return of, wrapErrorMsg(logger, msg, err)
	}

	matsConCatedras := make([]materia, 0, len(mats))
	for _, m := range mats {
		if len(m.Catedras) == 0 {
			logger.Warn("materia sin cátedras",
				"codigoMateria", m.Codigo, "nombreMateria", m.Nombre)
		} else {
			matsConCatedras = append(matsConCatedras, m)
		}
	}

	of = oferta{
		ofertaMetadata: md,
		materias:       matsConCatedras,
	}

	return of, nil
}

func (c cuatri) esDespuesDe(otro cuatri) bool {
	return (otro.anio < c.anio) ||
		((otro.anio == c.anio) && (otro.numero < c.numero))
}

// newOfertaComisiones serializa una oferta de comisión a partir de un archivo
// almacenado en el bucket.
// func newOfertaComisiones(cfg config.S3Config, ch chan oferta, objKey *string) error {
// 	logger := log.Default().WithPrefix("📄").With("objKey", *objKey)
//
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
// 	defer cancel()
//
// 	logger.Debug("obteniendo metadata")
//
// 	objHead, err := client.HeadObject(ctx, &s3.HeadObjectInput{
// 		Bucket: bucketKey,
// 		Key:    objKey,
// 	})
//
// 	if err != nil {
// 		logger.Error("error obteniendo metadata", "error", err)
// 		return err
// 	}
//
// 	// En este caso se sabe de antemano que si la información fue registrada y
// 	// leída correctamente, los datos son serializables.
// 	numero, _ := strconv.Atoi(objHead.Metadata["cuatri-numero"])
// 	anio, _ := strconv.Atoi(objHead.Metadata["cuatri-anio"])
//
// 	logger = logger.With("cuatri", numero, "anio", anio)
//
// 	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
// 	defer cancel()
//
// 	logger.Debug("obteniendo contenido")
//
// 	obj, err := client.getobject(ctx, &s3.getobjectinput{
// 		bucket: bucketkey,
// 		key:    objkey,
// 	})
//
// 	if err != nil {
// 		logger.error("error obteniendo contenido", "error", err)
// 		return err
// 	}
//
// 	defer obj.body.close()
//
// 	data, err := io.readall(obj.body)
// 	if err != nil {
// 		logger.error("error leyendo contenido", "error", err)
// 		return err
// 	}
//
// 	materias := []materia{}
//
// 	err = json.unmarshal(data, &materias)
// 	if err != nil {
// 		logger.error("error serializando oferta de comisiones", "error", err)
// 		return err
// 	}
//
// 	logger.infof("encontradas %v materias en oferta de comisiones", len(materias))
//
// 	materiasconcatedras := make([]materia, 0, len(materias))
// 	for _, m := range materias {
// 		if len(m.catedras) == 0 {
// 			logger.warn("materia sin cátedras", "materia", m.codigo)
// 		} else {
// 			materiasconcatedras = append(materiasconcatedras, m)
// 		}
// 	}
//
// 	ch <- oferta{
// 		carrera: objHead.Metadata["carrera"],
// 		cuatri: cuatri{
// 			numero: numero,
// 			anio:   anio,
// 		},
// 		materias: materiasConCatedras,
// 	}
//
// 	return nil
// }
