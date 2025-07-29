package tui

import (
	"log/slog"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patch"
)

type indiceVista uint

const (
	enListaMaterias indiceVista = iota
	enListaDocentes
	enVistaDocente
)

type AppModelo struct {
	indiceVista
	selectorMateria  listaMateriasModel
	selectorDocentes listaDocentesModel
	vistaDocente     vistaDocenteModel
	windowSize       tea.WindowSizeMsg
}

func newApp(patches []patch.Patch) AppModelo {
	// Ordenamos las materias según cantidad de docentes, de mayor a menor, para así dar prioridad
	// (al menos visual) a las materias que tengan más docentes.
	nDocentes := make(map[string]int, len(patches))
	for _, p := range patches {
		docentes := make(map[string]bool)
		for _, c := range p.Catedras {
			for _, d := range c.Docentes {
				docentes[d.Nombre] = true
			}
		}
		nDocentes[p.Nombre] = len(docentes)
	}

	sort.Slice(patches, func(i, j int) bool {
		return nDocentes[patches[i].Nombre] > nDocentes[patches[j].Nombre]
	})

	return AppModelo{
		selectorMateria:  newSelectorMateria(patches),
		selectorDocentes: newSelectorDocentes(),
		vistaDocente:     newVistaMateria(),
	}
}

func (m AppModelo) Init() tea.Cmd {
	slog.Info("iniciando resolvedor de patches gráfico")
	return m.selectorMateria.Init()
}

func (m AppModelo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case materiaSeleccionadaMsg:
		m.selectorDocentes.setDocentes(msg.patch)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			switch m.indiceVista {
			case enListaMaterias:
				m.indiceVista = enListaDocentes
			case enListaDocentes:
				m.indiceVista = enVistaDocente
			case enVistaDocente:
				m.indiceVista = enVistaDocente
			}
			return m, nil
		case "shift+tab":
			switch m.indiceVista {
			case enListaMaterias:
				m.indiceVista = enListaMaterias
			case enListaDocentes:
				m.indiceVista = enListaMaterias
			case enVistaDocente:
				m.indiceVista = enListaDocentes
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.windowSize = msg
	}

	var cmd tea.Cmd

	switch m.indiceVista {
	case enListaMaterias:
		m.selectorMateria, cmd = m.selectorMateria.Update(msg)
	case enListaDocentes:
		m.selectorDocentes, cmd = m.selectorDocentes.Update(msg)
	}

	return m, cmd
}

func (m AppModelo) View() string {
	var panel0, panel1, panel2 string

	if m.indiceVista == enListaMaterias {
		panel0 = estiloPanelActivo.Render(m.selectorMateria.View())
	} else {
		panel0 = estiloPanelInactivo.Render(m.selectorMateria.View())
	}

	if m.indiceVista == enListaDocentes {
		panel1 = estiloPanelActivo.Render(m.selectorDocentes.View())
	} else {
		panel1 = estiloPanelInactivo.Render(m.selectorDocentes.View())
	}

	anchoPanel0 := lipgloss.Width(panel0)
	anchoPanel1 := lipgloss.Width(panel1)
	anchoPanel2 := m.windowSize.Width - anchoPanel0 - anchoPanel1 - 2

	estiloPanel2Activo := estiloPanelActivo.Width(anchoPanel2)
	estiloPanel2Inactivo := estiloPanelInactivo.Width(anchoPanel2)

	if m.indiceVista == enVistaDocente {
		panel2 = estiloPanel2Activo.Render(m.vistaDocente.View())
	} else {
		panel2 = estiloPanel2Inactivo.Render(m.vistaDocente.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, panel0, panel1, panel2)
}

func ResolvePatches(patches []patch.Patch) {
	p := tea.NewProgram(newApp(patches))
	_, _ = p.Run()
}
