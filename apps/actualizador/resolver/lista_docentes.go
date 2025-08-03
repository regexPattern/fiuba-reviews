package resolver

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
)

type listaDocentesModel struct {
	docentesMaterias     map[string][]patcher.DocenteSiu
	contextosMaterias    map[string]patcher.ContextoMateriaDb
	indicesSeleccionados map[string]int
	nombreMateriaActual  string
	lista                list.Model
}

type docenteItem patcher.DocenteSiu

func (i docenteItem) Title() string {
	return fmt.Sprintf("%s (%s)", i.Nombre, i.Rol)
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
		docentesMaterias:     make(map[string][]patcher.DocenteSiu),
		contextosMaterias:    make(map[string]patcher.ContextoMateriaDb),
		indicesSeleccionados: make(map[string]int),
		lista:                l,
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
		m.indicesSeleccionados[m.nombreMateriaActual] = iActual

		docente := m.docentesMaterias[m.nombreMateriaActual][iActual]
		contextoMateria := m.contextosMaterias[m.nombreMateriaActual]
		return m, tea.Batch(cmd, seleccionarDocenteCmd(docente, contextoMateria))
	}

	return m, cmd
}

func (m listaDocentesModel) View() string {
	return m.lista.View()
}

func (m *listaDocentesModel) setDocentes(msg setMateriaMsg) tea.Cmd {
	return nil
}

type patchDocente struct {
	patcher.DocenteSiu
	patcher.ContextoMateriaDb
}

type setDocenteMsg patchDocente

func seleccionarDocenteCmd(
	docente patcher.DocenteSiu,
	contexto patcher.ContextoMateriaDb,
) tea.Cmd {
	return func() tea.Msg {
		return setDocenteMsg{
			docente,
			contexto,
		}
	}
}
