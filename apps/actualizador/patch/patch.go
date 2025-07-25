package patch

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"time"

	"golang.org/x/sync/errgroup"
)

type PatchMateria struct {
	CodigoSiu string
	Nombre    string
	Catedras  []catedraSiu
	cuatri
}

type GeneradorPatches struct {
	DbUrl         string
	DbInitTimeout time.Duration
	DbOpsTimeout  time.Duration
	S3BucketName  string
	S3InitTimeout time.Duration
	S3OpsTimeout  time.Duration
}

func (g *GeneradorPatches) GenerarPatches(ctx context.Context) ([]PatchMateria, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error { return g.initS3Client(egCtx) })
	eg.Go(func() error { return g.initDbPool(egCtx) })

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	var ofertas []*oferta
	var err error

	if ofertas, err = g.obtenerOfertasCarreras(ctx); err != nil {
		return nil, err
	}

	patches := filtrarOfertasMaterias(ofertas)
	if err := g.asociarMaterias(ctx, patches); err != nil {
		return nil, err
	}

	if err := g.migrarMaterias(ctx, patches); err != nil {
		return nil, err
	}

	return patches, nil
}

func filtrarOfertasMaterias(ofertas []*oferta) []PatchMateria {
	nMaterias := 0
	for _, o := range ofertas {
		nMaterias += len(o.Materias)
	}

	patches := make(map[string]PatchMateria, nMaterias)
	for _, o := range ofertas {
		for _, m := range o.Materias {
			pActual, ok := patches[m.Nombre]
			if !ok || o.cuatri.despuesDe(pActual.cuatri) {
				patches[m.Nombre] = PatchMateria{
					CodigoSiu: m.Codigo,
					Nombre:    m.Nombre,
					Catedras:  m.Catedras,
					cuatri:    o.cuatri,
				}
			} else if ok && pActual.cuatri == o.cuatri {
				// Si tenemos dos ofertas para una misma materia, y ambas
				// ofertas corresponden al mismo cuatrimestre, entonces
				// agregamos todas las cátedras de ambas ofertas al total de
				// la materia.

				c := make(map[int]catedraSiu)
				for _, cat := range pActual.Catedras {
					c[cat.Codigo] = cat
				}
				for _, cat := range m.Catedras {
					c[cat.Codigo] = cat
				}

				patches[m.Nombre] = PatchMateria{
					CodigoSiu: pActual.CodigoSiu,
					Nombre:    pActual.Nombre,
					Catedras:  slices.Collect(maps.Values(c)),
					cuatri:    o.cuatri,
				}
			}
		}
	}

	slog.Debug("unificado las últimas ofertas de materias", "unificadas", nMaterias-len(patches))
	slog.Info(fmt.Sprintf("generados %v patches de actualización de materias", len(patches)))

	return slices.Collect(maps.Values(patches))
}
