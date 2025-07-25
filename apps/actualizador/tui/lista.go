package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

var (
	fiubaColor           = lipgloss.Color("#4EACD4")
	itemNormalTitleStyle = lipgloss.NewStyle().
				PaddingLeft(3).
				Foreground(lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#DDDDDD"})
	itemSelectedTitleStyle     = lipgloss.NewStyle().Inherit(itemNormalTitleStyle)
	itemSelectedIndicatorStyle = lipgloss.NewStyle().
					PaddingLeft(1).
					Foreground(fiubaColor)
	itemFilterMatchStyle = lipgloss.NewStyle().
				Background(fiubaColor).
				Foreground(lipgloss.AdaptiveColor{Light: "#DDDDDD", Dark: "#1A1A1A"})
	listTitleBarStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				PaddingBottom(1)
	listTitleStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Background(fiubaColor).
			Foreground(lipgloss.AdaptiveColor{Light: "#DDDDDD", Dark: "#1A1A1A"})
	listFilterPromptStyle = lipgloss.NewStyle().Foreground(fiubaColor)
	listFilterCursorStyle = lipgloss.NewStyle().Foreground(fiubaColor)
)

type item string

func (i item) FilterValue() string {
	return string(i)
}

type itemDelegate struct{}

func (d itemDelegate) Height() int {
	return 1
}

func (d itemDelegate) Spacing() int {
	return 0
}

func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	textwidth := m.Width() - itemNormalTitleStyle.GetPaddingLeft() - itemNormalTitleStyle.GetPaddingRight()
	title := ansi.Truncate(fmt.Sprintf("%s", i), textwidth, "…")

	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	if emptyFilter {
		title = itemNormalTitleStyle.Render(title)
	} else if isSelected && m.FilterState() != list.Filtering {
		if isFiltered {
			matchedRunes := m.MatchesForItem(index)
			unmatched := itemSelectedTitleStyle.Inline(true)
			matched := itemFilterMatchStyle.Inherit(unmatched)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = fmt.Sprintf("%s %s", itemSelectedIndicatorStyle.Render(">"), itemSelectedTitleStyle.Render(title))
	} else {
		if isFiltered {
			matchedRunes := m.MatchesForItem(index)
			unmatched := itemNormalTitleStyle.Inline(true)
			matched := itemFilterMatchStyle.Inherit(unmatched)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = itemNormalTitleStyle.Render(title)
	}

	fmt.Fprint(w, title)
}

type listaModel struct {
	lista list.Model
}

func newListaModel(titulo string) listaModel {
	defaultKeys := list.DefaultKeyMap()
	l := list.New([]list.Item{}, itemDelegate{}, 50, 20)

	l.Title = titulo

	l.Styles.TitleBar = listTitleBarStyle
	l.Styles.Title = listTitleStyle
	l.Styles.FilterPrompt = listFilterPromptStyle
	l.Styles.FilterCursor = listFilterCursorStyle

	l.KeyMap = list.KeyMap{
		CursorUp:             defaultKeys.CursorUp,
		CursorDown:           defaultKeys.CursorDown,
		NextPage:             defaultKeys.NextPage,
		PrevPage:             defaultKeys.PrevPage,
		GoToStart:            defaultKeys.GoToStart,
		GoToEnd:              defaultKeys.GoToEnd,
		Filter:               defaultKeys.Filter,
		ClearFilter:          defaultKeys.ClearFilter,
		CancelWhileFiltering: defaultKeys.CancelWhileFiltering,
		AcceptWhileFiltering: defaultKeys.AcceptWhileFiltering,
	}
	l.Paginator.Type = paginator.Arabic

	itemName := strings.ToLower(titulo)
	l.SetStatusBarItemName(itemName[:len(itemName)-1], itemName)

	return listaModel{
		lista: l,
	}
}

func (m listaModel) Init() tea.Cmd {
	return nil
}

func (m listaModel) Update(msg tea.Msg) (listaModel, tea.Cmd) {
	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

func (m listaModel) View() string {
	return m.lista.View()
}

func (m *listaModel) setItems(items []string) {
	listItems := make([]list.Item, len(items))
	for idx, i := range items {
		listItems[idx] = item(i)
	}
	m.lista.SetItems(listItems)
}

func (m listaModel) globalIndex() int {
	return m.lista.GlobalIndex()
}
