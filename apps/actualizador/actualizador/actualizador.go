package actualizador

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"time"

	"golang.org/x/sync/errgroup"
)

type PatchActualizacionMateria struct {
	CodigoSiu string
	Nombre    string
	Catedras  []CatedraSiu
	cuatri
}

type IndexadorOfertas struct {
	DbUrl         string
	DbInitTimeout time.Duration
	DbOpsTimeout  time.Duration
	S3BucketName  string
	S3InitTimeout time.Duration
	S3OpsTimeout  time.Duration
}

func (i *IndexadorOfertas) GenerarPatchesDeActualizacion(
	ctx context.Context,
) ([]PatchActualizacionMateria, error) {
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error { return i.initS3Client(gCtx) })
	g.Go(func() error { return i.initDbPool(gCtx) })

	if err := g.Wait(); err != nil {
		return nil, err
	}

	var ofertas []*Oferta
	var err error

	if ofertas, err = i.obtenerOfertasCarreras(ctx); err != nil {
		return nil, err
	}

	patches := filtrarOfertasMaterias(ofertas)
	if err := i.asociarMaterias(ctx, patches); err != nil {
		return nil, err
	}

	if err := i.migrarMaterias(ctx, patches); err != nil {
		return nil, err
	}

	return patches, nil
}

func filtrarOfertasMaterias(ofertas []*Oferta) []PatchActualizacionMateria {
	nMaterias := 0
	for _, o := range ofertas {
		nMaterias += len(o.Materias)
	}

	patches := make(map[string]PatchActualizacionMateria, nMaterias)
	for _, o := range ofertas {
		for _, m := range o.Materias {
			pActual, ok := patches[m.Nombre]
			if !ok || o.cuatri.despuesDe(pActual.cuatri) {
				patches[m.Nombre] = PatchActualizacionMateria{
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

				c := make(map[int]CatedraSiu)
				for _, cat := range pActual.Catedras {
					c[cat.Codigo] = cat
				}
				for _, cat := range m.Catedras {
					c[cat.Codigo] = cat
				}

				patches[m.Nombre] = PatchActualizacionMateria{
					CodigoSiu: pActual.CodigoSiu,
					Nombre:    pActual.Nombre,
					Catedras:  slices.Collect(maps.Values(c)),
					cuatri:    o.cuatri,
				}
			}
		}
	}

	slog.Debug(
		"unificado las últimas ofertas de materias",
		"unificadas",
		nMaterias-len(patches),
	)
	slog.Info(
		fmt.Sprintf("generados %v patches de actualización de materias", len(patches)),
	)

	return slices.Collect(maps.Values(patches))
}
