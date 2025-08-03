package resolver

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
)

type state uint

const (
	listaMateriasFocused state = iota
	listaDocentesFocused
	resolverDocenteFocused
)

type model struct {
	state
	listaMaterias   listaMateriasModel
	listaDocentes   listaDocentesModel
	resolverDocente resolverDocenteModel
	dimensiones     tea.WindowSizeMsg
}

func newModel(patches []patcher.Patch) model {
	return model{
		listaMaterias:   newListaMaterias(patches),
		listaDocentes:   newListaDocentes(),
		resolverDocente: newVistaMateria(),
	}
}

func (m model) Init() tea.Cmd {
	return m.listaMaterias.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case setMateriaMsg:
		cmd := m.listaDocentes.setDocentes(msg)
		m.resolverDocente.setMateria(msg)
		return m, cmd

	case setDocenteMsg:
		m.resolverDocente.setDocente(msg)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			switch m.state {
			case listaMateriasFocused:
				m.state = listaDocentesFocused
			case listaDocentesFocused:
				m.state = resolverDocenteFocused
			case resolverDocenteFocused:
				m.state = resolverDocenteFocused
			}
			return m, nil
		case "shift+tab":
			switch m.state {
			case listaMateriasFocused:
				m.state = listaMateriasFocused
			case listaDocentesFocused:
				m.state = listaMateriasFocused
			case resolverDocenteFocused:
				m.state = listaDocentesFocused
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.dimensiones = msg
		return m, nil
	}

	var cmd tea.Cmd

	switch m.state {
	case listaMateriasFocused:
		m.listaMaterias, cmd = m.listaMaterias.Update(msg)
	case listaDocentesFocused:
		m.listaDocentes, cmd = m.listaDocentes.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	var panel0, panel1, panel2 string

	var style0 lipgloss.Style
	if m.state == listaMateriasFocused {
		style0 = estiloPanelActivo
	} else {
		style0 = estiloPanelInactivo
	}
	panel0 = style0.Render(m.listaMaterias.View())

	width0, height0 := lipgloss.Size(panel0)

	var style1 lipgloss.Style
	if m.state == listaDocentesFocused {
		style1 = estiloPanelActivo
	} else {
		style1 = estiloPanelInactivo
	}
	panel1 = style1.Width(width0 - style0.GetBorderLeftSize() - style0.GetBorderRightSize()).
		Height(height0 - style0.GetBorderTopSize() - style0.GetBorderBottomSize()).
		Render(m.listaDocentes.View())

	width1 := lipgloss.Width(panel1)
	width2 := m.dimensiones.Width - width0 - width1

	var style2 lipgloss.Style
	if m.state == resolverDocenteFocused {
		style2 = estiloPanelActivo
	} else {
		style2 = estiloPanelInactivo
	}
	panel2 = style2.Width(width2 - style0.GetBorderLeftSize() - style0.GetBorderRightSize()).
		Height(height0 - style0.GetBorderTopSize() - style0.GetBorderBottomSize()).
		Render(m.resolverDocente.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, panel0, panel1, panel2) + "\n"
}
