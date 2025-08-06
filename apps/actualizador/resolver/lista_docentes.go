package resolver

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/indexador"
)

type listaDocentesModel struct {
	lista list.Model
}

type docenteItem indexador.DocenteSiu

func (i docenteItem) Title() string {
	return i.Nombre
}

func (i docenteItem) Description() string {
	return ""
}

func (i docenteItem) FilterValue() string {
	return i.Nombre
}

func newListaDocentes() listaDocentesModel {
	l := newDefaultList()
	l.Title = "Docentes"

	return listaDocentesModel{
		lista: l,
	}
}

func (m listaDocentesModel) Init() tea.Cmd {
	return nil
}

func (m listaDocentesModel) Update(msg tea.Msg) (listaDocentesModel, tea.Cmd) {
	iAnterior := m.lista.GlobalIndex()

	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)

	if iActual := m.lista.GlobalIndex(); iActual != iAnterior {
		docente := m.lista.SelectedItem()
		if d, ok := docente.(docenteItem); ok {
			return m, tea.Batch(cmd, setDocenteCmd(indexador.DocenteSiu(d)))
		}
	}

	return m, cmd
}

func (m listaDocentesModel) View() string {
	return m.lista.View()
}

func (m *listaDocentesModel) setDocentes(materia setMateriaMsg, paginated bool) {
	items := []list.Item{}

	for _, c := range materia.Catedras {
		for _, d := range c.Docentes {
			items = append(items, docenteItem(d))
		}
	}

	m.lista.SetItems(items)
	m.lista.Select(0)

	height := listHeight
	if paginated && len(items) > m.lista.Paginator.PerPage {
		height = listHeight + 1
	}

	m.lista.SetHeight(height)
}

type patchDocente struct {
	indexador.DocenteSiu
}

type setDocenteMsg patchDocente

func setDocenteCmd(docente indexador.DocenteSiu) tea.Cmd {
	return func() tea.Msg {
		return setDocenteMsg{
			docente,
		}
	}
}
