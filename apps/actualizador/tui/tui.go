package tui

import (
	"log/slog"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patch"
)

type state uint

const (
	enListaMaterias state = iota
	enListaDocentes
	enInfoDocente
)

var (
	stylePanelBase     = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
	stylePanelActivo   = stylePanelBase.BorderForeground(fiubaColor)
	stylePanelInactivo = stylePanelBase.BorderForeground(lipgloss.Color("240"))
)

type model struct {
	state         state
	listaMaterias listaMateriasModel
	listaDocentes listaDocentesModel
	infoDocente   infoDocenteModel
}

func newModel(patches []patch.Patch) model {
	// Ordenamos las materias según cantidad de docentes de mayor a menor, para así dar prioridad
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

	return model{
		listaMaterias: newListaMaterias(patches),
		listaDocentes: newListaDocentes(),
		infoDocente:   newInfoDocente(),
	}
}

func (m model) Init() tea.Cmd {
	slog.Info("iniciando resolver de patches gráfico")

	return tea.Batch(tea.SetWindowTitle("FIUBA Reviews"))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+f":
			switch m.state {
			case enListaMaterias:
				m.state = enListaDocentes
			case enListaDocentes:
				m.state = enInfoDocente
			case enInfoDocente:
				m.state = enInfoDocente
			}
			return m, nil
		case "ctrl+b":
			switch m.state {
			case enListaMaterias:
				m.state = enListaMaterias
			case enListaDocentes:
				m.state = enListaMaterias
			case enInfoDocente:
				m.state = enListaDocentes
			}
			return m, nil
		}
	}

	var cmd tea.Cmd

	switch m.state {
	case enListaMaterias:
		prevPatch := m.listaMaterias.GetSelectedPatch()
		m.listaMaterias, cmd = m.listaMaterias.Update(msg)
		currPatch := m.listaMaterias.GetSelectedPatch()

		if prevPatch.Nombre != currPatch.Nombre {
			m.listaDocentes.SetPatch(currPatch)
		}
	case enListaDocentes:
		m.listaDocentes, cmd = m.listaDocentes.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	var lmView, ldView, idView string

	if m.state == enListaMaterias {
		lmView = stylePanelActivo.Render(m.listaMaterias.View())
	} else {
		lmView = stylePanelInactivo.Render(m.listaMaterias.View())
	}

	if m.state == enListaDocentes {
		ldView = stylePanelActivo.Render(m.listaDocentes.View())
	} else {
		ldView = stylePanelInactivo.Render(m.listaDocentes.View())
	}

	// if m.state == enInfoDocente {
	// 	singleDocenteStr = stylePanelActivo.Render(m.infoDocente.View())
	// } else {
	// 	singleDocenteStr = stylePanelInactivo.Render(m.infoDocente.View())
	// }

	return lipgloss.JoinHorizontal(lipgloss.Top, lmView, ldView, idView)
}

func ResolvePatches(patches []patch.Patch) {
	p := tea.NewProgram(newModel(patches))
	_, _ = p.Run()
}
