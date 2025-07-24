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

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4).MaxWidth(30)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).MaxWidth(30).Foreground(lipgloss.Color("170"))
)

type itemMateria string

func (i itemMateria) FilterValue() string {
	return string(i)
}

type itemMateriaDelegate struct{}

func (d itemMateriaDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(itemMateria)
	if !ok {
		return
	}

	label := fmt.Sprintf("%v", i)

	styleFn := itemStyle.Render
	if index == m.Index() {
		styleFn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, styleFn(label))
}

func (d itemMateriaDelegate) Height() int {
	return 1
}

func (d itemMateriaDelegate) Spacing() int {
	return 0
}

func (d itemMateriaDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

type listaMateriasModel struct {
	patches []patcher.PatchGenerado
	lista   list.Model
}

func newListaMaterias(patches []patcher.PatchGenerado) listaMateriasModel {
	items := make([]list.Item, len(patches))
	for i, p := range patches {
		items[i] = itemMateria(p.Nombre)
	}

	l := list.New(items, itemMateriaDelegate{}, 20, 20)

	l.Title = "Materias"

	m := listaMateriasModel{
		patches: patches,
		lista:   l,
	}

	return m
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

func (m listaMateriasModel) GetSelectedPatch() patcher.PatchGenerado {
	return m.patches[m.lista.GlobalIndex()]
}
