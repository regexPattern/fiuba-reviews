package tui

import (
	"maps"
	"slices"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patch"
)

type listaDocentesModel struct {
	docentesMaterias    map[string][]patch.DocenteSiu
	nombreMateriaActual string
	widgetLista         list.Model
}

type docenteItem patch.DocenteSiu

func (i docenteItem) Title() string {
	return i.Nombre
}

func (i docenteItem) Description() string {
	return ""
}

func (i docenteItem) FilterValue() string {
	return i.Nombre
}

func newSelectorDocentes() listaDocentesModel {
	l := newDefaultList()
	l.Title = "Docentes"

	return listaDocentesModel{
		docentesMaterias: make(map[string][]patch.DocenteSiu),
		widgetLista:      l,
	}
}

func (m listaDocentesModel) Init() tea.Cmd {
	return nil
}

func (m listaDocentesModel) Update(msg tea.Msg) (listaDocentesModel, tea.Cmd) {
	var cmd tea.Cmd
	m.widgetLista, cmd = m.widgetLista.Update(msg)
	return m, cmd
}

func (m listaDocentesModel) View() string {
	return m.widgetLista.View()
}

func (m *listaDocentesModel) setDocentes(p *patch.Patch) {
	m.nombreMateriaActual = p.Nombre
	if _, ok := m.docentesMaterias[p.Nombre]; !ok {
		docentesUnicos := make(map[patch.DocenteSiu]bool)
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
	items := make([]list.Item, len(m.docentesMaterias[p.Nombre]))
	for i, d := range m.docentesMaterias[p.Nombre] {
		items[i] = docenteItem(d)
	}
	m.widgetLista.SetItems(items)
}
