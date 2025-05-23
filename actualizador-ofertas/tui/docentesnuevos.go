package tui

import tea "github.com/charmbracelet/bubbletea"

type docentesNuevosView struct {
	docentes    []docenteSiu
	currDocente int
}

func (v docentesNuevosView) Update(msg tea.Msg) (docentesNuevosView, tea.Cmd) {
	var cmd tea.Cmd

	return v, cmd
}

func (v docentesNuevosView) View(_ int) string {
	return ""
}
