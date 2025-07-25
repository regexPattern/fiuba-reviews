package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patch"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/tui/lista"
)

type listaMateriasModel struct {
	patches []patch.PatchMateria
	lista   lista.Model
}

func newListaMaterias(patches []patch.PatchMateria) listaMateriasModel {
	l := lista.New("Materias")

	nombres := make([]string, len(patches))
	for i, p := range patches {
		nombres[i] = p.Nombre
	}

	l.SetItems(nombres)

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

func (m listaMateriasModel) getPatchSeleccionado() patch.PatchMateria {
	return m.patches[m.lista.GlobalIndex()]
}
