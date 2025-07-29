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

type Patch struct {
	CodigoSiu string
	Nombre    string
	Catedras  []CatedraSiu
	Cuatri
}

type Indexador struct {
	DbUrl         string
	DbInitTimeout time.Duration
	DbOpsTimeout  time.Duration
	S3BucketName  string
	S3InitTimeout time.Duration
	S3OpsTimeout  time.Duration
}

func (i *Indexador) GenerarPatches(ctx context.Context) ([]Patch, error) {
	if err := i.configClientesDatos(ctx); err != nil {
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

func (i *Indexador) configClientesDatos(ctx context.Context) error {
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return i.configClienteS3(gctx)
	})
	g.Go(func() error {
		return i.configPoolBD(gctx)
	})
	return g.Wait()
}

func filtrarOfertasMaterias(ofertas []*Oferta) []Patch {
	nMaterias := 0
	for _, o := range ofertas {
		nMaterias += len(o.Materias)
	}

	filtradas := make(map[string]Patch, nMaterias)
	for _, o := range ofertas {
		for _, m := range o.Materias {
			if pActual, ok := filtradas[m.Nombre]; !ok || o.Cuatri.despuesDe(pActual.Cuatri) {
				filtradas[m.Nombre] = Patch{
					CodigoSiu: m.Codigo,
					Nombre:    m.Nombre,
					Catedras:  m.Catedras,
					Cuatri:    o.Cuatri,
				}
			} else if ok && pActual.Cuatri == o.Cuatri {
				// Si tenemos dos ofertas para una misma materia, y ambas
				// ofertas corresponden al mismo cuatrimestre, entonces
				// agregamos todas las cátedras de ambas ofertas al total de
				// la materia.

				catedras := make(map[int]CatedraSiu)
				for _, c := range pActual.Catedras {
					catedras[c.Codigo] = c
				}
				for _, c := range m.Catedras {
					catedras[c.Codigo] = c
				}

				filtradas[m.Nombre] = Patch{
					CodigoSiu: pActual.CodigoSiu,
					Nombre:    pActual.Nombre,
					Catedras:  slices.Collect(maps.Values(catedras)),
					Cuatri:    o.Cuatri,
				}
			}
		}
	}

	slog.Debug(
		"unificado las últimas ofertas de materias",
		"unificadas",
		nMaterias-len(filtradas),
	)
	slog.Info(
		fmt.Sprintf("generados %v patches de actualización de materias", len(filtradas)),
	)

	return slices.Collect(maps.Values(filtradas))
}
