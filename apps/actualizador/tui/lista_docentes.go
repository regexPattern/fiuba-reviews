package tui

import (
	"fmt"
	"maps"
	"slices"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
)

type listaDocentesModel struct {
	docentesMaterias      map[string][]patcher.DocenteSiu
	contextosMaterias     map[string]patcher.ContextoMateriaBD
	indicesSeleccionados  map[string]int
	nombreMateriaActual   string
	widgetLista           list.Model
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
		contextosMaterias:    make(map[string]patcher.ContextoMateriaBD),
		indicesSeleccionados: make(map[string]int),
		widgetLista:          l,
	}
}

func (m listaDocentesModel) Init() tea.Cmd {
	return nil
}

func (m listaDocentesModel) Update(msg tea.Msg) (listaDocentesModel, tea.Cmd) {
	iAnterior := m.widgetLista.GlobalIndex()

	var cmd tea.Cmd
	m.widgetLista, cmd = m.widgetLista.Update(msg)

	if iActual := m.widgetLista.GlobalIndex(); iActual != iAnterior {
		// Guardar el índice seleccionado para esta materia
		m.indicesSeleccionados[m.nombreMateriaActual] = iActual
		
		docente := m.docentesMaterias[m.nombreMateriaActual][iActual]
		contextoMateria := m.contextosMaterias[m.nombreMateriaActual]
		return m, tea.Batch(cmd, seleccionarDocenteCmd(docente, contextoMateria))
	}

	return m, cmd
}

func (m listaDocentesModel) View() string {
	return m.widgetLista.View()
}

func (m *listaDocentesModel) setDocentes(msg materiaSeleccionadaMsg) tea.Cmd {
	nombreMateria := msg.Materia.Nombre
	m.nombreMateriaActual = nombreMateria

	m.contextosMaterias[nombreMateria] = msg.ContextoMateriaBD

	if _, ok := m.docentesMaterias[nombreMateria]; !ok {
		docentesUnicos := make(map[patcher.DocenteSiu]bool)
		for _, c := range msg.Materia.Catedras {
			for _, d := range c.Docentes {
				docentesUnicos[d] = true
			}
		}

		docentesOrdenados := slices.Collect(maps.Keys(docentesUnicos))
		sort.Slice(docentesOrdenados, func(i, j int) bool {
			return docentesOrdenados[i].Nombre < docentesOrdenados[j].Nombre
		})

		m.docentesMaterias[nombreMateria] = docentesOrdenados
		// Si es una materia nueva, inicializar índice en 0
		m.indicesSeleccionados[nombreMateria] = 0
	}
	
	items := make([]list.Item, len(m.docentesMaterias[nombreMateria]))
	for i, d := range m.docentesMaterias[nombreMateria] {
		items[i] = docenteItem(d)
	}
	m.widgetLista.SetItems(items)

	// Restaurar la selección previa para esta materia
	indiceSeleccionado := m.indicesSeleccionados[nombreMateria]
	if indiceSeleccionado < len(m.docentesMaterias[nombreMateria]) {
		m.widgetLista.Select(indiceSeleccionado)
		
		docenteSeleccionado := m.docentesMaterias[nombreMateria][indiceSeleccionado]
		contextoMateria := m.contextosMaterias[nombreMateria]
		return seleccionarDocenteCmd(docenteSeleccionado, contextoMateria)
	}
	
	return nil
}

type patchDocente struct {
	patcher.DocenteSiu
	patcher.ContextoMateriaBD
}

type docenteSeleccionadoMsg patchDocente

func seleccionarDocenteCmd(
	docente patcher.DocenteSiu,
	contexto patcher.ContextoMateriaBD,
) tea.Cmd {
	return func() tea.Msg {
		return docenteSeleccionadoMsg{
			docente,
			contexto,
		}
	}
}
