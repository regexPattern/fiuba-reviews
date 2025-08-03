package resolver

import (
	"cmp"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
)

func ResolverPatches(patches []patcher.Patch) {
	priorizarPatches(patches)

	p := tea.NewProgram(newModel(patches[:10]))
	_, _ = p.Run()
}

func priorizarPatches(patches []patcher.Patch) {
	nDocentes := make(map[string]int, len(patches))
	for _, p := range patches {
		docentesUnicos := make(map[string]bool)
		for _, c := range p.Materia.Catedras {
			for _, d := range c.Docentes {
				docentesUnicos[d.Nombre] = true
			}
		}
		nDocentes[p.Materia.Codigo] = len(docentesUnicos)
	}

	slices.SortFunc(patches, func(a, b patcher.Patch) int {
		return cmp.Compare(nDocentes[b.Materia.Codigo], nDocentes[a.Materia.Codigo])
	})
}
