package resolver

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type vistaDocenteModel struct {
	docente patchDocente
	materia string
}

func newVistaDocente() vistaDocenteModel {
	return vistaDocenteModel{}
}

func (m vistaDocenteModel) Init() tea.Cmd {
	return nil
}

func (m vistaDocenteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case setMateriaMsg:
		m.setMateria(msg)
	}
	return m, nil
}

func (m vistaDocenteModel) View() string {
	var s strings.Builder

	s.WriteString(m.materia + "\n")

	s.WriteString(m.docente.Nombre)

	return s.String()
}

func (m *vistaDocenteModel) setMateria(patch setMateriaMsg) tea.Cmd {
	m.materia = patch.Nombre
	return setDocenteCmd(patch.Catedras[0].Docentes[0])
}

func (m *vistaDocenteModel) setDocente(docente setDocenteMsg) {
	m.docente = patchDocente(docente)
}
