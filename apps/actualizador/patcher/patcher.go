package patcher

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

type Patch struct {
	OfertaMateriaSiu
	ContextoMateriaDb
}

type Indexador struct {
	DbUrl         string
	DbInitTimeout time.Duration
	DbOpTimeout   time.Duration
	S3BucketName  string
	S3InitTimeout time.Duration
	S3OpTimeout   time.Duration
}

func (i *Indexador) GenerarPatches(ctx context.Context) ([]Patch, error) {
	if err := i.configClientesDatos(ctx); err != nil {
		return nil, err
	}

	var ofertas []OfertaMateriaSiu
	var err error

	if ofertas, err = i.getOfertasMateriasSiu(ctx); err != nil {
		return nil, err
	}

	if err := i.syncMateriasSiuConDb(ctx, ofertas); err != nil {
		return nil, err
	}

	if patches, err := i.completarPatches(ctx, ofertas); err != nil {
		return nil, err
	} else {
		return patches, nil
	}
}

func (i *Indexador) configClientesDatos(ctx context.Context) error {
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return i.configClienteS3(gCtx)
	})
	g.Go(func() error {
		return i.configPoolDb(gCtx)
	})
	return g.Wait()
}

func (i *Indexador) completarPatches(
	ctx context.Context,
	ofertas []OfertaMateriaSiu,
) ([]Patch, error) {
	bdCtx, bdCancel := context.WithTimeout(ctx, i.DbOpTimeout)
	defer bdCancel()

	patchesCh := make(chan Patch, len(ofertas))
	g, gCtx := errgroup.WithContext(bdCtx)

	for _, o := range ofertas {
		g.Go(func() error {
			if c, err := getContextoMateriaDb(gCtx, o.Materia); err != nil {
				return err
			} else if c != nil {
				patchesCh <- Patch{
					OfertaMateriaSiu:  o,
					ContextoMateriaDb: *c,
				}
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	close(patchesCh)

	patches := make([]Patch, 0, len(ofertas))
	for p := range patchesCh {
		patches = append(patches, p)
	}

	return patches, nil
}
