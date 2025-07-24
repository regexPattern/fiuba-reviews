package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patch"
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

	// Truncar con elipsis si excede el ancho máximo considerando el padding
	maxLen := maxItemWidth - 4 // Restamos el padding left de 4
	if index == m.Index() {
		maxLen = maxItemWidth - 4 // Para el item activo también consideramos "> " (2 chars) + padding (2)
	}

	if len(i) > maxLen {
		if maxLen > 3 {
			i = i[:maxLen-3] + "..."
		} else {
			i = "..."
		}
	}

	styleFn := styleItemLista.Render
	if index == m.Index() {
		styleFn = func(s ...string) string {
			return styleItemActivoLista.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, styleFn(string(i)))
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
	patches []patch.Patch
	lista   list.Model
}

func newListaMaterias(patches []patch.Patch) listaMateriasModel {
	items := make([]list.Item, len(patches))
	for i, p := range patches {
		items[i] = itemMateria(p.Nombre)
	}

	l := list.New(items, itemMateriaDelegate{}, listWidth, listHeight)

	l.KeyMap.CloseFullHelp.Unbind()
	l.KeyMap.ShowFullHelp.Unbind()
	l.KeyMap.Quit.Unbind()
	l.KeyMap.ForceQuit.Unbind()

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

func (m listaMateriasModel) GetSelectedPatch() patch.Patch {
	return m.patches[m.lista.GlobalIndex()]
}
