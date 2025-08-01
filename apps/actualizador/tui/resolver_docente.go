package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type resolverDocenteModel struct {
	docentesPorMateria map[string]docenteSeleccionadoMsg
	materiaActual      string
}

func newVistaMateria() resolverDocenteModel {
	return resolverDocenteModel{
		docentesPorMateria: make(map[string]docenteSeleccionadoMsg),
	}
}

func (m resolverDocenteModel) Init() tea.Cmd {
	return nil
}

func (m resolverDocenteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case docenteSeleccionadoMsg:
		m.setDocente(msg)
	case materiaSeleccionadaMsg:
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


func (m *resolverDocenteModel) setMateria(msg materiaSeleccionadaMsg) {
	m.materiaActual = msg.Materia.Nombre
}

func (m *resolverDocenteModel) setDocente(msg docenteSeleccionadoMsg) {
	m.docentesPorMateria[m.materiaActual] = msg
}
