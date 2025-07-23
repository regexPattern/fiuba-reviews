package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
)

func ResolvePatches(patches []patcher.PatchGenerado) []patcher.PatchResolution {
	p := tea.NewProgram(newModel(patches))

	_, _ = p.Run()

	return nil
}
