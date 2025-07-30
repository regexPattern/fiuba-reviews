package resolvedor

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
)

type materiaSeleccionadaMsg struct {
	patch *patcher.Patch
	info  *patcher.InfoActualMateria
}

func materiaSeleccionadaCmd(patch *patcher.Patch) tea.Cmd {
	return tea.Batch(tea.SetWindowTitle(patch.Materia.Nombre), func() tea.Msg {
		return materiaSeleccionadaMsg{patch: patch}
	})
}
