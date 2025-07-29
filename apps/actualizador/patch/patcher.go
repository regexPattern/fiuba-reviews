package patch

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

type Patch struct {
	CodigoSiu string
	Nombre    string
	Catedras  []CatedraSiu
	Cuatri
}

type Indexador struct {
	DbUrl         string
	DbInitTimeout time.Duration
	DbOpTimeout   time.Duration
	S3BucketName  string
	S3InitTimeout time.Duration
	S3OpTimeout   time.Duration
}

type NuevoPatch struct {
	OfertaMateriaSiu
	ContextoMateriaBD
}

func (i *Indexador) GenerarPatches(ctx context.Context) ([]NuevoPatch, error) {
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
) ([]NuevoPatch, error) {
	patchesCh := make(chan NuevoPatch, len(ofertas))

	bdCtx, bdCancel := context.WithTimeout(ctx, i.DbOpTimeout)
	defer bdCancel()

	g, gCtx := errgroup.WithContext(bdCtx)

	for _, o := range ofertas {
		g.Go(func() error {
			if c, err := getContextoMateriaBD(gCtx, o.Materia); err != nil {
				return err
			} else {
				patchesCh <- NuevoPatch{
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

	patches := make([]NuevoPatch, 0, len(ofertas))
	for p := range patchesCh {
		patches = append(patches, p)
	}

	return patches, nil
}
