package tui

import (
	"log/slog"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
)

type viewState uint

const (
	materiasView viewState = iota
	docentesView
	singleDocenteView
)

type model struct {
	state    viewState
	materias materiasModel
}

func newModel(patches []patcher.PatchGenerado) model {
	return model{
		materias: newMateriasModel(patches),
	}
}

func (m model) Init() tea.Cmd {
	slog.Info("iniciando resolver de patches gráfico")
	return tea.Batch(tea.SetWindowTitle("FIUBA Reviews"))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.materias, cmd = m.materias.Update(msg)

	return m, cmd
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString(m.materias.View())

	return s.String()
}
