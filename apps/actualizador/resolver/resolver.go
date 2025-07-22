package resolver

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
)

func ResolvePatches(props *patcher.PatchProposal) []patcher.PatchResolution {
	p := tea.NewProgram(model{})
	_, _ = p.Run()

	return nil
}
