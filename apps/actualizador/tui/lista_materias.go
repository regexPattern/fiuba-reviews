package tui

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
	return seleccionarMateriaCmd(&m.patches[0])
}

func (m listaMateriasModel) Update(msg tea.Msg) (listaMateriasModel, tea.Cmd) {
	iAnterior := m.lista.GlobalIndex()

	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)

	if iActual := m.lista.GlobalIndex(); iActual != iAnterior {
		return m, tea.Batch(cmd, seleccionarMateriaCmd(&m.patches[iActual]))
	}

	return m, cmd
}

func (m listaMateriasModel) View() string {
	return m.lista.View()
}

type materiaSeleccionadaMsg *patcher.Patch

func seleccionarMateriaCmd(patch *patcher.Patch) tea.Cmd {
	titulo := fmt.Sprintf(
		"fiuba-reviews • %s • %s",
		patch.Materia.Codigo,
		patch.Materia.Nombre,
	)
	return tea.Batch(
		tea.SetWindowTitle(titulo),
		func() tea.Msg {
			return materiaSeleccionadaMsg(patch)
		},
	)
}
