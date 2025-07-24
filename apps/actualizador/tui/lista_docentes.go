package tui

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patch"
)

type itemDocente string

func (i itemDocente) FilterValue() string {
	return string(i)
}

type itemDocenteDelegate struct{}

func (d itemDocenteDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(itemDocente)
	if !ok {
		return
	}

	label := fmt.Sprint(i)

	// Truncar con elipsis si excede el ancho máximo considerando el padding
	maxLen := maxItemWidth - 4 // Restamos el padding left de 4
	if index == m.Index() {
		maxLen = maxItemWidth - 4 // Para el item activo también consideramos "> " (2 chars) + padding (2)
	}

	if len(label) > maxLen {
		if maxLen > 3 {
			label = label[:maxLen-3] + "..."
		} else {
			label = "..."
		}
	}

	styleFn := styleItemLista.Render
	if index == m.Index() {
		styleFn = func(s ...string) string {
			return styleItemActivoLista.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, styleFn(label))
}

func (d itemDocenteDelegate) Height() int {
	return 1
}

func (d itemDocenteDelegate) Spacing() int {
	return 0
}

func (d itemDocenteDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

type docentesModel struct {
	patch             patch.Patch
	docentesOrdenados []string
	lista             list.Model
}

func newListaDocentes() docentesModel {
	l := list.New([]list.Item{}, itemDocenteDelegate{}, listWidth, listHeight)

	l.Title = "Docentes SIU"
	l.SetWidth(listWidth)
	l.SetHeight(listHeight)
	l.SetShowHelp(false)

	l.KeyMap.CloseFullHelp.Unbind()
	l.KeyMap.ShowFullHelp.Unbind()
	l.KeyMap.Quit.Unbind()
	l.KeyMap.ForceQuit.Unbind()

	return docentesModel{
		lista: l,
	}
}

func (m docentesModel) Init() tea.Cmd {
	return nil
}

func (m docentesModel) Update(msg tea.Msg) (docentesModel, tea.Cmd) {
	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

func (m docentesModel) View() string {
	return m.lista.View()
}

func (m *docentesModel) SetPatch(patch patch.Patch) {
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

	items := make([]list.Item, len(m.docentesOrdenados))
	for i, nombre := range m.docentesOrdenados {
		items[i] = itemDocente(nombre)
	}
	m.lista.SetItems(items)
}

func (m docentesModel) GetSelectedDocente() string {
	return m.docentesOrdenados[m.lista.GlobalIndex()]
}
