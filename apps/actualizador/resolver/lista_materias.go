package resolver

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/indexador"
)

type listaMateriasModel struct {
	patches []indexador.OfertaMateriaSiu
	lista   list.Model
}

type patchItem indexador.OfertaMateriaSiu

func (i patchItem) Title() string {
	return i.Nombre
}

func (i patchItem) Description() string {
	return ""
}

func (i patchItem) FilterValue() string {
	return i.Nombre
}

func newListaMaterias(patches []indexador.OfertaMateriaSiu) listaMateriasModel {
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
	patch := m.patches[0]
	docente := patch.Catedras[0].Docentes[0]
	return tea.Batch(setMateriaCmd(patch), setDocenteCmd(docente))
}

func (m listaMateriasModel) Update(msg tea.Msg) (listaMateriasModel, tea.Cmd) {
	iAnterior := m.lista.GlobalIndex()

	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)

	if iActual := m.lista.GlobalIndex(); iActual != iAnterior {
		return m, tea.Batch(cmd, setMateriaCmd(m.patches[iActual]))
	}

	return m, cmd
}

func (m listaMateriasModel) View() string {
	return m.lista.View()
}

type setMateriaMsg indexador.OfertaMateriaSiu

func setMateriaCmd(patch indexador.OfertaMateriaSiu) tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle(fmt.Sprintf(
			"fiuba-reviews • %s • %s",
			patch.Codigo,
			patch.Nombre,
		)),
		func() tea.Msg {
			return setMateriaMsg(patch)
		},
		setDocenteCmd(patch.Catedras[0].Docentes[0]),
	)
}
