package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/actualizador"
)

type materiaSeleccionadaMsg struct {
	patch *actualizador.PatchActualizacionMateria
}

func materiaSeleccionadaCmd(
	p *actualizador.PatchActualizacionMateria,
) tea.Cmd {
	return tea.Batch(tea.SetWindowTitle(p.Nombre), func() tea.Msg {
		return materiaSeleccionadaMsg{p}
	})
}
