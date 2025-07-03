package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/regexPattern/fiuba-reviews/actualizador-ofertas/actualizador"
)

var (
	colorFIUBA = lipgloss.Color("#4eacd4")
)

type model struct {
	materias    []materiaView
	currMateria int
	screenWidth int
}

func Run(patches []actualizador.PatchMateriaOutput) error {
	m := newModel(patches)
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

func newModel(patches []actualizador.PatchMateriaOutput) model {
	mats := newMateriaViews(patches)
	m := model{
		materias: mats,
	}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.screenWidth = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case ">":
			if m.currMateria < len(m.materias)-1 {
				m.currMateria++
			} else {
				m.currMateria = 0
			}
			return m, nil
		case "<":
			if m.currMateria > 0 {
				m.currMateria--
			} else {
				m.currMateria = len(m.materias) - 1
			}
			return m, nil
		}
	}

	m.materias[m.currMateria], cmd =
		m.materias[m.currMateria].Update(msg)

	return m, cmd
}

var (
	styleTituloMateria = lipgloss.NewStyle()
)

func (m model) View() string {
	s := strings.Builder{}
	mat := m.materias[m.currMateria]

	s.WriteString(styleTituloMateria.
		Render(fmt.Sprintf("%v • %v", mat.codigo, strings.ToUpper(mat.nombre))))
	s.WriteString("\n")
	s.WriteString(mat.View(m.screenWidth))

	return s.String()
}
