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

type getInfoMateriaMsg struct {
	patch actualizador.PatchActualizacionMateria
	info  actualizador.InfoMateria
	err   error
}

func getInfoMateriaCmd(p actualizador.PatchActualizacionMateria) tea.Cmd {
	return func() tea.Msg {
		info, err := actualizador.GetInfoMateria(p.CodigoSiu)
		if err != nil {
			return getInfoMateriaMsg{err: err}
		}
		return getInfoMateriaMsg{
			patch: p,
			info:  *info,
		}
	}
}

type view uint

const (
	enListaMaterias view = iota
	enListaDocentes
	enInfoDocente
)

var (
	stylePanelBase     = lipgloss.NewStyle().Border(lipgloss.ThickBorder())
	stylePanelActivo   = stylePanelBase.BorderForeground(colorFiuba)
	stylePanelInactivo = stylePanelBase.BorderForeground(lipgloss.Color("240"))
)

type Model struct {
	view          view
	listaMaterias listaMateriasModel
	listaDocentes listaDocentesModel
	vistaMateria  vistaMateriaModel
}

func newModel(patches []actualizador.PatchActualizacionMateria) Model {
	// Ordenamos las materias según cantidad de de mayor a menor, para así dar prioridad
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

	return Model{
		listaMaterias: newListaMaterias(patches),
		listaDocentes: newListaDocentes(),
		vistaMateria:  newVistaMateria(),
	}
}

func (m Model) Init() tea.Cmd {
	slog.Info("iniciando resolvedor de patches gráfico")
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case getInfoMateriaMsg:
		if msg.err != nil {
			m.vistaMateria.setError(msg.err)
		} else {
			m.vistaMateria.setInfoMateria(msg.patch, msg.info)
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			switch m.view {
			case enListaMaterias:
				m.view = enListaDocentes
			case enListaDocentes:
				m.view = enInfoDocente
			case enInfoDocente:
				m.view = enInfoDocente
			}
			return m, nil
		case "shift+tab":
			switch m.view {
			case enListaMaterias:
				m.view = enListaMaterias
			case enListaDocentes:
				m.view = enListaMaterias
			case enInfoDocente:
				m.view = enListaDocentes
			}
			return m, nil
		}
	}

	var cmd tea.Cmd

	switch m.view {
	case enListaMaterias:
		// Si luego de haber actualizado el componente de la lista de materias la materia
		// seleccionada cambió, entonces descargamos los docentes de la nueva materia seleccionada.
		// Para esto necesitamos guardarnos el indice seleccionado antes de ejecutar el comando en
		// la lista.
		prevIdx := m.listaMaterias.lista.globalIndex()
		m.listaMaterias, cmd = m.listaMaterias.Update(msg)

		if m.listaMaterias.lista.globalIndex() != prevIdx {
			p := m.listaMaterias.getPatchSeleccionado()
			return m, tea.Batch(cmd, getInfoMateriaCmd(p))
		}
	case enListaDocentes:
		m.listaDocentes, cmd = m.listaDocentes.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	var lm, ld, id string

	if m.view == enListaMaterias {
		lm = stylePanelActivo.Width(30 + 4).Height(21).Render(m.listaMaterias.View())
	} else {
		lm = stylePanelInactivo.Width(30).Height(21).Render(m.listaMaterias.View())
	}

	if m.view == enListaDocentes {
		ld = stylePanelActivo.Width(30).Height(21).Render(m.listaDocentes.View())
	} else {
		ld = stylePanelInactivo.Width(30).Height(21).Render(m.listaDocentes.View())
	}

	if m.view == enInfoDocente {
		id = stylePanelActivo.Height(21).Render(m.vistaMateria.View())
	} else {
		id = stylePanelInactivo.Height(21).Render(m.vistaMateria.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, lm, ld, id)
}

func ResolvePatches(patches []actualizador.PatchActualizacionMateria) {
	p := tea.NewProgram(newModel(patches))
	_, _ = p.Run()
}
