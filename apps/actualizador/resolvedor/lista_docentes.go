package resolvedor

import (
	"fmt"
	"maps"
	"slices"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
)

type listaDocentesModel struct {
	docentesMaterias    map[string][]patcher.DocenteSiu
	nombreMateriaActual string
	widgetLista         list.Model
}

type docenteItem patcher.DocenteSiu

func (i docenteItem) Title() string {
	return fmt.Sprintf("%s (%s)", i.Nombre, i.Rol)
}

func (i docenteItem) Description() string {
	return ""
}

func (i docenteItem) FilterValue() string {
	return i.Nombre
}

func newListaDocentes() listaDocentesModel {
	l := newDefaultList()
	l.Title = "Docentes"

	return listaDocentesModel{
		docentesMaterias: make(map[string][]patcher.DocenteSiu),
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

func (m *listaDocentesModel) setDocentes(p *patcher.Patch) {
	m.nombreMateriaActual = p.Materia.Nombre
	if _, ok := m.docentesMaterias[p.Materia.Nombre]; !ok {
		docentesUnicos := make(map[patcher.DocenteSiu]bool)
		for _, c := range p.Materia.Catedras {
			for _, d := range c.Docentes {
				docentesUnicos[d] = true
			}
		}

		docentesOrdenados := slices.Collect(maps.Keys(docentesUnicos))
		sort.Slice(docentesOrdenados, func(i, j int) bool {
			return docentesOrdenados[i].Nombre < docentesOrdenados[j].Nombre
		})

		m.docentesMaterias[p.Materia.Nombre] = docentesOrdenados
	}
	items := make([]list.Item, len(m.docentesMaterias[p.Materia.Nombre]))
	for i, d := range m.docentesMaterias[p.Materia.Nombre] {
		items[i] = docenteItem(d)
	}
	m.widgetLista.SetItems(items)
}
