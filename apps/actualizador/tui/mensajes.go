package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patch"
)

type materiaSeleccionadaMsg struct {
	patch *patch.Patch
	info  *patch.InfoActualMateria
}

func materiaSeleccionadaCmd(
	p *patch.Patch,
) tea.Cmd {
	return tea.Batch(tea.SetWindowTitle(p.Nombre), func() tea.Msg {
		info, _ := patch.GetInfoMateria(p.CodigoSiu)
		return materiaSeleccionadaMsg{p, info}
	})
}
