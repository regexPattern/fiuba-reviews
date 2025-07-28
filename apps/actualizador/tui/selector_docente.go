package tui

import (
	"maps"
	"slices"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/actualizador"
)

type selectorDocentes struct {
	docentesMaterias    map[string][]actualizador.DocenteSiu
	nombreMateriaActual string
	lista               listaModel
}

func newSelectorDocentes() selectorDocentes {
	l := NewLista("Docentes")

	return selectorDocentes{
		docentesMaterias: make(map[string][]actualizador.DocenteSiu),
		lista:            l,
	}
}

func (m selectorDocentes) Init() tea.Cmd {
	return nil
}

func (m selectorDocentes) Update(msg tea.Msg) (selectorDocentes, tea.Cmd) {
	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

func (m selectorDocentes) View() string {
	nombresDocentes := make([]string, 0)

	for i, d := range m.docentesMaterias[m.nombreMateriaActual] {
		if i < 15 {
			nombresDocentes = append(nombresDocentes, d.Nombre)
		}
	}

	return strings.Join(nombresDocentes, "\n")
}

func (m *selectorDocentes) setDocentes(p *actualizador.PatchActualizacionMateria) {
	m.nombreMateriaActual = p.Nombre
	if _, ok := m.docentesMaterias[p.Nombre]; !ok {
		docentesUnicos := make(map[actualizador.DocenteSiu]bool)
		for _, c := range p.Catedras {
			for _, d := range c.Docentes {
				docentesUnicos[d] = true
			}
		}

		docentesOrdenados := slices.Collect(maps.Keys(docentesUnicos))
		sort.Slice(docentesOrdenados, func(i, j int) bool {
			return docentesOrdenados[i].Nombre < docentesOrdenados[j].Nombre
		})

		m.docentesMaterias[p.Nombre] = docentesOrdenados
	}
}
