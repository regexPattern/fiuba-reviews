package resolver

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
)

type listaMateriasModel struct {
	patches []patcher.Patch
	lista   list.Model
}

type patchItem patcher.Patch

func (i patchItem) Title() string {
	return i.Materia.Nombre
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
		patches: patches,
		lista:   l,
	}
}

func (m listaMateriasModel) Init() tea.Cmd {
	return setMateriaCmd(&m.patches[0])
}

func (m listaMateriasModel) Update(msg tea.Msg) (listaMateriasModel, tea.Cmd) {
	iAnterior := m.lista.GlobalIndex()

	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)

	if iActual := m.lista.GlobalIndex(); iActual != iAnterior {
		return m, tea.Batch(cmd, setMateriaCmd(&m.patches[iActual]))
	}

	return m, cmd
}

func (m listaMateriasModel) View() string {
	return m.lista.View()
}

type setMateriaMsg *patcher.Patch

func setMateriaCmd(patch *patcher.Patch) tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle(fmt.Sprintf(
			"fiuba-reviews • %s • %s",
			patch.Materia.Codigo,
			patch.Materia.Nombre,
		)),
		func() tea.Msg {
			return setMateriaMsg(patch)
		},
	)
}
