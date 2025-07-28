package tui

import (
	"log/slog"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/actualizador"
)

var (
	colorFiuba = lipgloss.Color("#4EACD4")
)

type indiceVista uint

const (
	enListaMaterias indiceVista = iota
	enListaDocentes
	enVistaDocente
)

var (
	stylePanelBase     = lipgloss.NewStyle().Border(lipgloss.ThickBorder())
	stylePanelActivo   = stylePanelBase.BorderForeground(colorFiuba)
	stylePanelInactivo = stylePanelBase.BorderForeground(lipgloss.Color("240"))
)

type App struct {
	indiceVista
	selectorMateria  selectorMateriaModel
	selectorDocentes selectorDocentes
	vistaDocente     vistaDocenteModel
}

func newApp(patches []actualizador.PatchActualizacionMateria) App {
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

	return App{
		selectorMateria:  newSelectorMateria(patches),
		selectorDocentes: newSelectorDocentes(),
		vistaDocente:     newVistaMateria(),
	}
}

func (m App) Init() tea.Cmd {
	slog.Info("iniciando resolvedor de patches gráfico")
	return m.selectorMateria.Init()
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m App) View() string {
	var lm, ld, id string

	if m.indiceVista == enListaMaterias {
		lm = stylePanelActivo.Width(30 + 4).Height(21).Render(m.selectorMateria.View())
	} else {
		lm = stylePanelInactivo.Width(30 + 4).Height(21).Render(m.selectorMateria.View())
	}

	if m.indiceVista == enListaDocentes {
		ld = stylePanelActivo.Width(30 + 4).Height(21).Render(m.selectorDocentes.View())
	} else {
		ld = stylePanelInactivo.Width(30 + 4).Height(21).Render(m.selectorDocentes.View())
	}

	if m.indiceVista == enVistaDocente {
		id = stylePanelActivo.Height(21).Render(m.vistaDocente.View())
	} else {
		id = stylePanelInactivo.Height(21).Render(m.vistaDocente.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, lm, ld, id)
}

func ResolvePatches(patches []actualizador.PatchActualizacionMateria) {
	p := tea.NewProgram(newApp(patches))
	_, _ = p.Run()
}
