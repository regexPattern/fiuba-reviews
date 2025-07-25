package tui

import (
	"fmt"
	"log/slog"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patch"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/tui/color"
)

type actualizacionDocentesMsg struct {
	patch      patch.PatchMateria
	docentesDb []*patch.DocenteDb
}

type errorDocentesMsg struct {
	patch patch.PatchMateria
	error error
}

func cargarDocentesCmd(p patch.PatchMateria) tea.Cmd {
	return func() tea.Msg {
		d, err := patch.ObtenerDocentesMateria(p.CodigoSiu)
		if err != nil {
			return errorDocentesMsg{patch: p, error: err}
		}
		return actualizacionDocentesMsg{patch: p, docentesDb: d}
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
	stylePanelActivo   = stylePanelBase.BorderForeground(color.FiubaColor)
	stylePanelInactivo = stylePanelBase.BorderForeground(lipgloss.Color("240"))
)

type model struct {
	view          view
	listaMaterias listaMateriasModel
	listaDocentes listaDocentesModel
	infoDocente   infoDocenteModel
}

func newModel(patches []patch.PatchMateria) model {
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
	slog.Info("iniciando resolvedor de patches gráfico")
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case actualizacionDocentesMsg:
		m.infoDocente.setDocentesMateria(msg.patch.CodigoSiu, msg.docentesDb)
		return m, nil
	case errorDocentesMsg:
		m.infoDocente.setErrorDocentes(msg.patch.CodigoSiu, msg.error)
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
		prevIdx := m.listaMaterias.lista.GlobalIndex()
		m.listaMaterias, cmd = m.listaMaterias.Update(msg)

		if m.listaMaterias.lista.GlobalIndex() != prevIdx {
			p := m.listaMaterias.getPatchSeleccionado()
			return m, tea.Batch(cmd, cargarDocentesCmd(p))
		}
	case enListaDocentes:
		m.listaDocentes, cmd = m.listaDocentes.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	var lm, ld, id string

	if m.view == enListaMaterias {
		lm = stylePanelActivo.Render(m.listaMaterias.View())
	} else {
		lm = stylePanelInactivo.Render(m.listaMaterias.View())
	}

	if m.view == enListaDocentes {
		ld = stylePanelActivo.Render(m.listaDocentes.View())
	} else {
		ld = stylePanelInactivo.Render(m.listaDocentes.View())
	}

	if m.view == enInfoDocente {
		id = stylePanelActivo.Render(m.infoDocente.View())
	} else {
		id = stylePanelInactivo.Render(m.infoDocente.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, lm, ld, id)
}

func ResolvePatches(patches []patch.PatchMateria) {
	// f, _ := tea.LogToFile("debug.log", "debug")
	// defer f.Close()

	// p := tea.NewProgram(newModel(patches))
	// _, _ = p.Run()
	fmt.Println()
}
