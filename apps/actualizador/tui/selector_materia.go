package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/actualizador"
)

type selectorMateriaModel struct {
	patches []actualizador.PatchActualizacionMateria
	lista   listaModel
}

func newSelectorMateria(patches []actualizador.PatchActualizacionMateria) selectorMateriaModel {
	l := NewLista("Materias")

	nombresMaterias := make([]string, len(patches))
	for i, p := range patches {
		nombresMaterias[i] = p.Nombre
	}

	l.setItems(nombresMaterias)

	return selectorMateriaModel{
		patches: patches,
		lista:   l,
	}
}

func (m selectorMateriaModel) Init() tea.Cmd {
	if len(m.patches) > 0 {
		m := &m.patches[0]
		return materiaSeleccionadaCmd(m)
	} else {
		return nil
	}
}

func (m selectorMateriaModel) Update(msg tea.Msg) (selectorMateriaModel, tea.Cmd) {
	iAnterior := m.lista.globalIndex()

	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)

	if iActual := m.lista.globalIndex(); iActual != iAnterior {
		return m, tea.Batch(cmd, materiaSeleccionadaCmd(&m.patches[iActual]))
	}

	return m, nil
}

func (m selectorMateriaModel) View() string {
	return m.lista.View()
}
