package resolver

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patch"
)

func ResolvePatches(props *patch.PatchProposal) []patch.PatchResolution {
	p := tea.NewProgram(model{})
	_, _ = p.Run()

	return nil
}
