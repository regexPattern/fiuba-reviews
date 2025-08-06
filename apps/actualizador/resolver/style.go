package resolver

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const (
	listWidth  = 30
	listHeight = 20
)

var fiubaColor = lipgloss.Color("#4EACD4")

var (
	panelStyle          = lipgloss.NewStyle().Border(lipgloss.BlockBorder())
	focusedPanelStyle   = panelStyle.BorderForeground(fiubaColor)
	unfocusedPanelStyle = panelStyle.BorderForeground(lipgloss.Color("240"))
)

func newDefaultList() list.Model {
	d := list.NewDefaultDelegate()
	d.SetHeight(1)
	d.SetSpacing(0)
	d.ShowDescription = false

	l := list.New([]list.Item{}, d, listWidth, listHeight)

	l.DisableQuitKeybindings()
	l.KeyMap.ShowFullHelp.Unbind()
	l.KeyMap.CloseFullHelp.Unbind()

	l.SetShowHelp(false)

	l.Styles.TitleBar = l.Styles.TitleBar.UnsetPaddingLeft()
	l.Styles.TitleBar = l.Styles.TitleBar.UnsetPaddingRight()
	l.Styles.StatusBar = l.Styles.StatusBar.PaddingLeft(1)
	l.Styles.NoItems = l.Styles.NoItems.PaddingLeft(1)
	l.Styles.PaginationStyle = l.Styles.PaginationStyle.PaddingLeft(1)

	return l
}
