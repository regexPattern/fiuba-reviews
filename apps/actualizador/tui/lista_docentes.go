package tui

import (
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patch"
)

type listaDocentesModel struct {
	patch             patch.Patch
	docentesOrdenados []string
	lista             listaModel
}

func newListaDocentes() listaDocentesModel {
	lista := newListaModel("Docentes")

	return listaDocentesModel{
		lista: lista,
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

func (m *listaDocentesModel) SetPatch(patch patch.Patch) {
	m.patch = patch

	docentes := make(map[string]bool)
	for _, c := range m.patch.Catedras {
		for _, d := range c.Docentes {
			docentes[d.Nombre] = true
		}
	}

	m.docentesOrdenados = make([]string, 0, len(docentes))
	for nombre := range docentes {
		m.docentesOrdenados = append(m.docentesOrdenados, nombre)
	}

	sort.Strings(m.docentesOrdenados)

	m.lista.setItems(m.docentesOrdenados)
}

func (m listaDocentesModel) GetSelectedDocente() string {
	return m.docentesOrdenados[m.lista.globalIndex()]
}
