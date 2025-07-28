package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/actualizador"
)

type vistaMateriaModel struct {
	patch actualizador.PatchActualizacionMateria
	info  actualizador.InfoMateria
	err   error
}

func newVistaMateria() vistaMateriaModel {
	return vistaMateriaModel{}
}

func (m vistaMateriaModel) Init() tea.Cmd {
	return nil
}

func (m vistaMateriaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m vistaMateriaModel) View() string {
	var s strings.Builder

	if m.err != nil {
		return lipgloss.NewStyle().
			Underline(true).
			Foreground(lipgloss.Color("#FF0000")).
			Render(m.err.Error())
	}

	s.WriteString(fmt.Sprintf("%s - %s\n", m.patch.CodigoSiu, m.info.Nombre))

	s.WriteString(fmt.Sprintln(m.err))

	return s.String()
}

func (m *vistaMateriaModel) setInfoMateria(p actualizador.PatchActualizacionMateria, i actualizador.InfoMateria) {
	m.patch = p
	m.info = i
	m.err = nil
}

func (m *vistaMateriaModel) setError(err error) {
	m.err = err
}
