package lista

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/tui/color"
)

const listaWidth = 30
const listaHeight = 20

var (
	itemNormalTitleStyle = lipgloss.NewStyle().
				PaddingLeft(3).
				Foreground(lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#DDDDDD"})
	itemSelectedTitleStyle     = lipgloss.NewStyle().Inherit(itemNormalTitleStyle)
	itemSelectedIndicatorStyle = lipgloss.NewStyle().
					PaddingLeft(1).
					Foreground(color.FiubaColor)
	itemFilterMatchStyle = lipgloss.NewStyle().
				Background(color.FiubaColor).
				Foreground(lipgloss.AdaptiveColor{Light: "#DDDDDD", Dark: "#1A1A1A"})
	listTitleBarStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				PaddingBottom(1)
	listTitleStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Background(color.FiubaColor).
			Foreground(lipgloss.AdaptiveColor{Light: "#DDDDDD", Dark: "#1A1A1A"})
	listFilterPromptStyle = lipgloss.NewStyle().Foreground(color.FiubaColor)
	listFilterCursorStyle = lipgloss.NewStyle().Foreground(color.FiubaColor)
)

type item string

func (i item) FilterValue() string {
	return string(i)
}

type delegate struct{}

func (d delegate) Height() int {
	return 1
}

func (d delegate) Spacing() int {
	return 0
}

func (d delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d delegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
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

type Model struct {
	lista list.Model
}

func New(titulo string) Model {
	defaultKeys := list.DefaultKeyMap()
	l := list.New([]list.Item{}, delegate{}, listaWidth, listaHeight)

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

	l.SetShowHelp(false)

	return Model{l}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.lista.View()
}

func (m *Model) SetItems(items []string) {
	listItems := make([]list.Item, len(items))
	for idx, i := range items {
		listItems[idx] = item(i)
	}
	m.lista.SetItems(listItems)
}

func (m Model) GlobalIndex() int {
	return m.lista.GlobalIndex()
}
