package resolvedor

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
)

type vistaDocenteModel struct {
	info patcher.InfoActualMateria
	err  error
}

func newVistaMateria() vistaDocenteModel {
	return vistaDocenteModel{}
}

func (m vistaDocenteModel) Init() tea.Cmd {
	return nil
}

func (m vistaDocenteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m vistaDocenteModel) View() string {
	var s strings.Builder

	if m.err != nil {
		return lipgloss.NewStyle().
			Underline(true).
			Foreground(lipgloss.Color("#FF0000")).
			Render(m.err.Error())
	}

	// s.WriteString(fmt.Sprintf("%s - %s\n", m.patch.CodigoSiu, m.info.Nombre))

	s.WriteString(fmt.Sprintln(m.err))

	return s.String()
}

func (m *vistaDocenteModel) setInfoMateria(i materiaSeleccionadaMsg) {
}

func (m *vistaDocenteModel) setError(err error) {
	m.err = err
}
