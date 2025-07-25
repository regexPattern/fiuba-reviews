package tui

import "github.com/charmbracelet/lipgloss"

var (
	listWidth            = 50
	listHeight           = 20
	maxItemWidth         = 46
	styleItemLista       = lipgloss.NewStyle().PaddingLeft(2).MaxWidth(maxItemWidth)
	styleItemActivoLista = lipgloss.NewStyle().MaxWidth(maxItemWidth).Foreground(fiubaColor)
)
