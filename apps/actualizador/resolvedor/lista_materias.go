package resolvedor

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
)

type listaMateriasModel struct {
	patches     []patcher.Patch
	widgetLista list.Model
}

type patchItem patcher.Patch

func (i patchItem) Title() string {
	return i.ContextoMateriaBD.Nombre
}

func (i patchItem) Description() string {
	return ""
}

func (i patchItem) FilterValue() string {
	return i.Materia.Nombre
}

func newListaMaterias(patches []patcher.Patch) listaMateriasModel {
	l := newDefaultList()
	l.Title = "Materias"

	items := make([]list.Item, len(patches))
	for i, p := range patches {
		items[i] = patchItem(p)
	}
	l.SetItems(items)

	return listaMateriasModel{
		patches:     patches,
		widgetLista: l,
	}
}

func (m listaMateriasModel) Init() tea.Cmd {
	if len(m.patches) > 0 {
		m := &m.patches[0]
		return materiaSeleccionadaCmd(m)
	} else {
		return nil
	}
}

func (m listaMateriasModel) Update(msg tea.Msg) (listaMateriasModel, tea.Cmd) {
	iAnterior := m.widgetLista.GlobalIndex()

	var cmd tea.Cmd
	m.widgetLista, cmd = m.widgetLista.Update(msg)

	if iActual := m.widgetLista.GlobalIndex(); iActual != iAnterior {
		return m, tea.Batch(cmd, materiaSeleccionadaCmd(&m.patches[iActual]))
	}

	return m, cmd
}

func (m listaMateriasModel) View() string {
	return m.widgetLista.View()
}
