package patch

import (
	"context"
	"log/slog"
	"maps"
	"slices"
	"time"
)

// Ofeta de cátedras y docentes de una materia obtenida desde las ofertas de comisiones del SIU.
type Patch struct {
	Codigo   string
	Nombre   string
	Catedras []catedraSiu
	cuatri
}

// Genera los patches de actualización para las materias accediendo a las ofetas en S3 y las
// materias en la base de datos utilizando la configuración provista.
type GeneradorPatches struct {
	DbUrl         string
	DbTimeout     time.Duration
	S3BucketName  string
	S3InitTimeout time.Duration
	S3Timeout     time.Duration
}

// GenerarPatches genera los patches de actualización para las materias.
func (g *GeneradorPatches) GenerarPatches(ctx context.Context) ([]Patch, error) {
	if err := g.initClienteS3(ctx); err != nil {
		return nil, err
	}

	var oc []*ofertaCarrera
	var err error
	if oc, err = g.obtenerOfertasCarreras(ctx); err != nil {
		return nil, err
	}

	p := filtrarOfertasMaterias(oc)
	if err := g.actualizarCodigosMaterias(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

// filtrarOfertasMaterias unifica las materias de diferentes ofertas de carrera y se queda con la
// oferta más reciente para cada una.
func filtrarOfertasMaterias(oc []*ofertaCarrera) []Patch {
	nMaterias := 0
	for _, o := range oc {
		nMaterias += len(o.Materias)
	}

	p := make(map[string]Patch, nMaterias)
	for _, o := range oc {
		for _, m := range o.Materias {
			pActual, ok := p[m.Nombre]
			if !ok || o.cuatri.despuesDe(pActual.cuatri) {
				p[m.Nombre] = Patch{
					Codigo:   m.Codigo,
					Nombre:   m.Nombre,
					Catedras: m.Catedras,
					cuatri:   o.cuatri,
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

				p[m.Nombre] = Patch{
					Codigo:   pActual.Codigo,
					Nombre:   pActual.Nombre,
					Catedras: slices.Collect(maps.Values(c)),
					cuatri:   o.cuatri,
				}
			}
		}
	}

	slog.Debug("unificado las últimas ofertas de materias",
		"n_inicial", nMaterias,
		"n_final", len(p))

	return slices.Collect(maps.Values(p))
}
