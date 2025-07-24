package tui

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
)

var (
	itemDocenteStyle         = lipgloss.NewStyle().PaddingLeft(4).MaxWidth(30)
	selectedDocenteItemStyle = lipgloss.NewStyle().PaddingLeft(2).MaxWidth(30).Foreground(lipgloss.Color("170"))
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

	label := fmt.Sprintf("%v", i)

	styleFn := itemDocenteStyle.Render
	if index == m.Index() {
		styleFn = func(s ...string) string {
			return selectedDocenteItemStyle.Render("> " + strings.Join(s, " "))
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
	patch             patcher.PatchGenerado
	docentesOrdenados []string
	lista             list.Model
}

func newListaDocentes() docentesModel {
	l := list.New([]list.Item{}, itemDocenteDelegate{}, 20, 20)
	l.Title = "Docentes SIU"
	
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

func (m *docentesModel) SetPatch(patch patcher.PatchGenerado) {
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
