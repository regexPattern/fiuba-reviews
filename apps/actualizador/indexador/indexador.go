package indexador

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

type Indexador struct {
	DbUrl         string
	DbInitTimeout time.Duration
	DbOpTimeout   time.Duration
	DbTxTimeout   time.Duration
	S3BucketName  string
	S3InitTimeout time.Duration
	S3OpTimeout   time.Duration
}

func (i *Indexador) ObtenerMaterias(ctx context.Context) ([]Materia, error) {
	var err error

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error { return i.initClienteS3(gCtx) })
	g.Go(func() error { return i.initPoolDb(gCtx) })

	if err = g.Wait(); err != nil {
		return nil, err
	}

	var ofertas []OfertaMateriaSiu
	if ofertas, err = i.obtenerOfertasSiu(ctx); err != nil {
		return nil, err
	}

	var materias []Materia
	if materias, err = i.sincronizarConDb(ctx, ofertas); err != nil {
		return nil, err
	}

	return materias, nil
}
