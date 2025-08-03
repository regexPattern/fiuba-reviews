package resolver

import (
	tea "github.com/charmbracelet/bubbletea"
)

type resolverDocenteModel struct {
	docentesPorMateria map[string]setDocenteMsg
	materiaActual      string
}

func newVistaMateria() resolverDocenteModel {
	return resolverDocenteModel{
		docentesPorMateria: make(map[string]setDocenteMsg),
	}
}

func (m resolverDocenteModel) Init() tea.Cmd {
	return nil
}

func (m resolverDocenteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case setDocenteMsg:
		m.setDocente(msg)
	case setMateriaMsg:
		m.setMateria(msg)
	}
	return m, nil
}

func (m resolverDocenteModel) View() string {
	if docente, ok := m.docentesPorMateria[m.materiaActual]; ok {
		return docente.Nombre
	}
	return "Selecciona un docente"
}


func (m *resolverDocenteModel) setMateria(msg setMateriaMsg) {
	m.materiaActual = msg.Materia.Nombre
}

func (m *resolverDocenteModel) setDocente(msg setDocenteMsg) {
	m.docentesPorMateria[m.materiaActual] = msg
}
