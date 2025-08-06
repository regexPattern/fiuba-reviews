package resolver

import (
	"cmp"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/indexador"
)

func ResolverPatches(patches []indexador.OfertaMateriaSiu) {
	sortPatchesSegunPrioridad(patches)

	p := tea.NewProgram(newModel(patches))
	_, _ = p.Run()
}

func sortPatchesSegunPrioridad(patches []indexador.OfertaMateriaSiu) {
	nDocentes := make(map[string]int, len(patches))
	for _, p := range patches {
		docentesUnicos := make(map[string]bool)
		for _, c := range p.Catedras {
			for _, d := range c.Docentes {
				docentesUnicos[d.Nombre] = true
			}
		}
		nDocentes[p.Codigo] = len(docentesUnicos)
	}

	slices.SortFunc(patches, func(a, b indexador.OfertaMateriaSiu) int {
		return cmp.Compare(nDocentes[b.Codigo], nDocentes[a.Codigo])
	})
}
