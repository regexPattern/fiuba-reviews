package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/actualizador"
)

type listaDocentesModel struct {
	patch *actualizador.PatchActualizacionMateria
	lista listaModel
}

func newListaDocentes() listaDocentesModel {
	l := NewLista("Docentes")

	return listaDocentesModel{
		patch: nil,
		lista: l,
	}
}

func (m listaDocentesModel) Init() tea.Cmd {
	return nil
}

func (m listaDocentesModel) Update(msg tea.Msg) (listaDocentesModel, tea.Cmd) {
	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

func (m listaDocentesModel) View() string {
	return m.lista.View()
}

func (m *listaDocentesModel) SetPatch(patch *actualizador.PatchActualizacionMateria) {
	m.patch = patch

	docentes := make(map[string]bool)
	for _, c := range m.patch.Catedras {
		for _, d := range c.Docentes {
			docentes[d.Nombre] = true
		}
	}

	// TODO: realmente deberia hacer este sort cuando creo los patches
	// al menos desde el TUI
	// m.lista.setItems([]{|)
}
