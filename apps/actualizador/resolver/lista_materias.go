package resolver

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/indexador"
)

type listaMateriasModel struct {
	materias []indexador.Materia
	widget   list.Model
}

type materiaItem indexador.Materia

func (i materiaItem) Title() string {
	return i.MateriaSiu.Nombre
}

func (i materiaItem) Description() string {
	return ""
}

func (i materiaItem) FilterValue() string {
	return i.MateriaSiu.Nombre
}

func newListaMaterias(materias []indexador.Materia) listaMateriasModel {
	l := newDefaultList()
	l.Title = "Materias"

	items := make([]list.Item, len(materias))
	for i, m := range materias {
		items[i] = materiaItem(m)
	}
	l.SetItems(items)

	return listaMateriasModel{
		materias: materias,
		widget:   l,
	}
}

func (m listaMateriasModel) Init() tea.Cmd {
	materia := m.materias[0]
	docente := materia.Catedras[0].Docentes[0]
	return tea.Batch(setMateriaCmd(materia), setDocenteCmd(docente))
}

func (m listaMateriasModel) Update(msg tea.Msg) (listaMateriasModel, tea.Cmd) {
	iAnterior := m.widget.GlobalIndex()

	var cmd tea.Cmd
	m.widget, cmd = m.widget.Update(msg)

	if iActual := m.widget.GlobalIndex(); iActual != iAnterior {
		return m, tea.Batch(cmd, setMateriaCmd(m.materias[iActual]))
	}

	return m, cmd
}

func (m listaMateriasModel) View() string {
	return m.widget.View()
}

type setMateriaMsg indexador.Materia

func setMateriaCmd(m indexador.Materia) tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle(fmt.Sprintf(
			"fiuba-reviews • %s • %s",
			m.MateriaSiu.Codigo,
			m.MateriaDb.Nombre,
		)),
		func() tea.Msg {
			return setMateriaMsg(m)
		},
	)
}
