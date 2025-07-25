package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patch"
)

type listaMateriasModel struct {
	patches []patch.Patch
	lista   listaModel
}

func newListaMaterias(patches []patch.Patch) listaMateriasModel {
	l := newListaModel("Materias")

	nombres := make([]string, len(patches))
	for i, p := range patches {
		nombres[i] = p.Nombre
	}
	
	l.setItems(nombres)

	return listaMateriasModel{
		patches: patches,
		lista:   l,
	}
}

func (m listaMateriasModel) Init() tea.Cmd {
	return nil
}

func (m listaMateriasModel) Update(msg tea.Msg) (listaMateriasModel, tea.Cmd) {
	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

func (m listaMateriasModel) View() string {
	return m.lista.View()
}

func (m listaMateriasModel) GetSelectedPatch() patch.Patch {
	return m.patches[m.lista.globalIndex()]
}
