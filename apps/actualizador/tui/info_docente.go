package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patch"
)

type infoDocenteModel struct {
	codigo     string
	docentesDb []*patch.DocenteDb
	err        error
}

func newInfoDocente() infoDocenteModel {
	return infoDocenteModel{}
}

func (m infoDocenteModel) Init() tea.Cmd {
	return nil
}

func (m infoDocenteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m infoDocenteModel) View() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintln(m.codigo))

	s.WriteString(fmt.Sprintln(len(m.docentesDb)))
	for _, d := range m.docentesDb {
		s.WriteString(fmt.Sprintln(d.Nombre))
	}

	s.WriteString(fmt.Sprintln(m.err))

	return s.String()
}

func (m *infoDocenteModel) setDocentesMateria(codigo string, d []*patch.DocenteDb) {
	m.codigo = codigo
	m.docentesDb = d
}

func (m *infoDocenteModel) setErrorDocentes(_ string, err error) {
	m.err = err
}
