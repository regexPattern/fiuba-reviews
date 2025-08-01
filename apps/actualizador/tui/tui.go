package tui

import (
	"log/slog"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
)

type indiceVista uint

const (
	enListaMaterias indiceVista = iota
	enListaDocentes
	enResolverDocente
)

type Model struct {
	indiceVista
	listaMaterias   listaMateriasModel
	listaDocentes   listaDocentesModel
	resolverDocente resolverDocenteModel
	dimensiones     tea.WindowSizeMsg
}

func NewModel(patches []patcher.Patch) Model {
	// Ordenamos los patches de materias según cantidad de docentes, de mayor a menor, para así
	// dar prioridad (al menos visual) a las materias que tengan más docentes.
	nDocentes := make(map[string]int, len(patches))
	for _, p := range patches {
		docentes := make(map[string]bool)
		for _, c := range p.Materia.Catedras {
			for _, d := range c.Docentes {
				docentes[d.Nombre] = true
			}
		}
		nDocentes[p.Materia.Codigo] = len(docentes)
	}

	sort.Slice(patches, func(i, j int) bool {
		return nDocentes[patches[i].Materia.Codigo] > nDocentes[patches[j].Materia.Codigo]
	})

	return Model{
		listaMaterias:   newListaMaterias(patches),
		listaDocentes:   newListaDocentes(),
		resolverDocente: newVistaMateria(),
	}
}

func (m Model) Init() tea.Cmd {
	slog.Info("iniciando resolvedor de patches gráfico")
	return m.listaMaterias.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Manejamos el cambio de la materia y docentes seleccionado desde el modelo principal de la app
	// porque este es el que tiene acceso a las tres vistas. Esto tiene que ser asi ya no podemos
	// hacer que cada componente se actualice por su cuenta, debido al conflicto de keybinds que
	// habría entre ambas listas.
	case materiaSeleccionadaMsg:
		m.resolverDocente.setMateria(msg)
		cmd := m.listaDocentes.setDocentes(msg)
		return m, cmd

	case docenteSeleccionadoMsg:
		m.resolverDocente.setDocente(msg)
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
				m.indiceVista = enResolverDocente
			case enResolverDocente:
				m.indiceVista = enResolverDocente
			}
			return m, nil
		case "shift+tab":
			switch m.indiceVista {
			case enListaMaterias:
				m.indiceVista = enListaMaterias
			case enListaDocentes:
				m.indiceVista = enListaMaterias
			case enResolverDocente:
				m.indiceVista = enListaDocentes
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.dimensiones = msg
		return m, nil
	}

	var cmd tea.Cmd

	switch m.indiceVista {
	case enListaMaterias:
		m.listaMaterias, cmd = m.listaMaterias.Update(msg)
	case enListaDocentes:
		m.listaDocentes, cmd = m.listaDocentes.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	var panel0, panel1, panel2 string

	if m.indiceVista == enListaMaterias {
		panel0 = estiloPanelActivo.Render(m.listaMaterias.View())
	} else {
		panel0 = estiloPanelInactivo.Render(m.listaMaterias.View())
	}

	if m.indiceVista == enListaDocentes {
		panel1 = estiloPanelActivo.Render(m.listaDocentes.View())
	} else {
		panel1 = estiloPanelInactivo.Render(m.listaDocentes.View())
	}

	anchoPanel0 := lipgloss.Width(panel0)
	anchoPanel1 := lipgloss.Width(panel1)
	anchoPanel2 := m.dimensiones.Width - anchoPanel0 - anchoPanel1 - 2

	if m.indiceVista == enResolverDocente {
		panel2 = estiloPanelActivo.Width(anchoPanel2).Render(m.resolverDocente.View())
	} else {
		panel2 = estiloPanelInactivo.Width(anchoPanel2).Render(m.resolverDocente.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, panel0, panel1, panel2) + "\n"
}
