package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
)

const (
	widthLista  = 20
	heightLista = 20
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type materiaItem string

func (i materiaItem) FilterValue() string {
	return string(i)
}

type materiaItemDelegate struct{}

func (d materiaItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(materiaItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%v", i)

	styleFn := itemStyle.Render
	if index == m.Index() {
		styleFn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, styleFn(str))
}

func (d materiaItemDelegate) Height() int {
	return 1
}

func (d materiaItemDelegate) Spacing() int {
	return 0
}

func (d materiaItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

type materiasModel struct {
	patches []patcher.PatchGenerado
	lista   list.Model
}

func newMateriasModel(patches []patcher.PatchGenerado) materiasModel {
	items := make([]list.Item, len(patches))
	for i, p := range patches {
		items[i] = materiaItem(p.Nombre)
	}

	l := list.New(items, materiaItemDelegate{}, widthLista, heightLista)

	l.Title = "Materias"

	return materiasModel{
		patches: patches,
		lista:   l,
	}
}

func (m materiasModel) Init() tea.Cmd {
	return nil
}

func (m materiasModel) Update(msg tea.Msg) (materiasModel, tea.Cmd) {
	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

func (m materiasModel) View() string {
	return m.lista.View()
}
