package patcher

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

type Indexador struct {
	DbUrl         string
	DbInitTimeout time.Duration
	DbOpTimeout   time.Duration
	S3BucketName  string
	S3InitTimeout time.Duration
	S3OpTimeout   time.Duration
}

type Patch struct {
	OfertaMateriaSiu
	ContextoMateriaBD
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

	if err := i.syncMateriasSiuConBD(ctx, ofertas); err != nil {
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
		return i.configPoolBD(gCtx)
	})
	return g.Wait()
}

func (i *Indexador) completarPatches(
	ctx context.Context,
	ofertas []OfertaMateriaSiu,
) ([]Patch, error) {
	codigos := make([]string, len(ofertas))
	for i, o := range ofertas {
		codigos[i] = o.Materia.Codigo
	}

	nombres, err := i.getNombresMateriasBD(ctx, codigos)
	if err != nil {
		return nil, err
	}

	bdCtx, bdCancel := context.WithTimeout(ctx, i.DbOpTimeout)
	defer bdCancel()

	patchesCh := make(chan Patch, len(ofertas))
	g, gCtx := errgroup.WithContext(bdCtx)

	for _, o := range ofertas {
		g.Go(func() error {
			nombreBD := nombres[o.Materia.Codigo]
			if c, err := getContextoMateriaBD(gCtx, o.Materia, nombreBD); err != nil {
				return err
			} else {
				patchesCh <- Patch{
					OfertaMateriaSiu:  o,
					ContextoMateriaBD: c,
				}
				return nil
			}
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
