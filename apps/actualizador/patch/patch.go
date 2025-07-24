package patch

import (
	"context"
	"log/slog"
	"maps"
	"slices"
	"time"
)

type Patch struct {
	CodigoDb  string
	CodigoSiu string
	Nombre    string
	Catedras  []catedraSiu
	cuatri
}

type GeneradorPatches struct {
	DbUrl         string
	DbTimeout     time.Duration
	S3BucketName  string
	S3InitTimeout time.Duration
	S3Timeout     time.Duration
}

func (g *GeneradorPatches) GenerarPatches(ctx context.Context) ([]Patch, error) {
	if err := g.initClienteS3(ctx); err != nil {
		return nil, err
	}

	var oc []*ofertaCarrera
	var err error
	if oc, err = g.obtenerOfertasCarreras(ctx); err != nil {
		return nil, err
	}

	om := filtrarOfertasMaterias(oc)
	if err := g.actualizarCodigosMaterias(om); err != nil {
		return nil, err
	}

	return []Patch{}, nil
}

type ofertaMateria struct {
	materia materiaSiu
	cuatri
	carrera string
}

func filtrarOfertasMaterias(oc []*ofertaCarrera) []ofertaMateria {
	nMaterias := 0
	for _, o := range oc {
		nMaterias += len(o.Materias)
	}

	om := make(map[string]ofertaMateria, nMaterias)
	for _, o := range oc {
		for _, m := range o.Materias {
			oMasReciente, ok := om[m.Nombre]
			if !ok || o.cuatri.despuesDe(oMasReciente.cuatri) {
				om[m.Nombre] = ofertaMateria{
					materia: m,
					cuatri:  o.cuatri,
					carrera: o.carrera,
				}
			} else if ok && oMasReciente.cuatri == o.cuatri {
				// Si tenemos dos ofertas para una misma materia, y ambas
				// ofertas corresponden al mismo cuatrimestre, entonces
				// agregamos todas las cátedras de ambas ofertas al total de
				// la materia.

				c := make(map[int]catedraSiu)
				for _, cat := range oMasReciente.materia.Catedras {
					c[cat.Codigo] = cat
				}
				for _, cat := range m.Catedras {
					c[cat.Codigo] = cat
				}

				materiaUnificada := oMasReciente.materia
				materiaUnificada.Catedras = slices.Collect(maps.Values(c))

				om[m.Nombre] = ofertaMateria{
					materia: materiaUnificada,
					cuatri:  o.cuatri,
					carrera: o.carrera,
				}
			}
		}
	}

	slog.Debug("unificado las últimas ofertas de materias", "n_inicial", nMaterias, "n_final", len(om))

	return slices.Collect(maps.Values(om))
}
