package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const (
	anchoLista  = 30
	alturaLista = 20
)

var (
	colorFiuba = lipgloss.Color("#4EACD4")
)

var (
	estiloPanelBase = lipgloss.NewStyle().
			Width(anchoLista+5).
			Height(alturaLista+1).
			Padding(0, 1).
			Border(lipgloss.ThickBorder())
	estiloPanelActivo   = estiloPanelBase.BorderForeground(colorFiuba)
	estiloPanelInactivo = estiloPanelBase.BorderForeground(lipgloss.Color("240"))
)

func newDefaultList() list.Model {
	d := list.NewDefaultDelegate()
	d.SetHeight(1)
	d.SetSpacing(0)
	d.ShowDescription = false

	l := list.New([]list.Item{}, d, anchoLista, alturaLista)
	l.DisableQuitKeybindings()
	l.SetShowHelp(false)

	l.Styles.TitleBar = l.Styles.TitleBar.UnsetPaddingLeft()
	l.Styles.TitleBar = l.Styles.TitleBar.UnsetPaddingRight()
	l.Styles.StatusBar = l.Styles.StatusBar.UnsetPaddingLeft()
	l.Styles.StatusBar = l.Styles.StatusBar.UnsetPaddingRight()
	l.Styles.PaginationStyle = l.Styles.PaginationStyle.PaddingLeft(1)

	return l
}
