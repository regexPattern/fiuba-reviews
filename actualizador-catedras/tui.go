package main

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type modeloApp struct {
	pagsMaterias []modeloPagMateria
	currMateria  int
	screenWidth  int
	help         help.Model
}

type modeloPagMateria struct {
	codigo       string
	nombre       string
	tabsDocentes modeloTabsDocentes
}

type modeloTabsDocentes struct {
	docentesDb    []modeloDocenteDb
	currDocenteDb int
	docentesSiu   []modeloDocenteSiu
}

type modeloDocenteDb struct {
	codigo  string
	nombre  string
	matches []modeloDocenteSiu
}

type modeloDocenteSiu struct {
	nombre string
	rol    string
}

func (m modeloPagMateria) tituloMateria() string {
	return fmt.Sprintf("%v • %v", m.codigo, strings.ToUpper(m.nombre))
}

func newModeloApp(patches []patch) modeloApp {
	materias := make([]modeloPagMateria, 0, len(patches))

	for _, p := range patches {
		docentesDb := newModeloDocentesDb(*p.docentes)

		pm := modeloPagMateria{
			codigo: p.codigoMateria,
			nombre: p.nombreMateria,
			tabsDocentes: modeloTabsDocentes{
				docentesDb:    docentesDb,
				currDocenteDb: 0,
			},
		}

		materias = append(materias, pm)
	}

	slices.SortFunc(materias, func(m1, m2 modeloPagMateria) int {
		return cmp.Compare(m1.codigo, m2.codigo)
	})

	modelo := modeloApp{
		pagsMaterias: materias,
		currMateria:  0,
	}

	return modelo
}

func newModeloDocentesDb(pd patchDocentes) []modeloDocenteDb {
	docentes := make([]modeloDocenteDb, 0, len(pd.db))
	nombresDocentesSiu := slices.Collect(maps.Keys(pd.siu))

	for nombre, cod := range pd.db {
		matches := fuzzy.RankFind(nombre, nombresDocentesSiu)
		if len(matches) == 0 {
			continue
		}

		sort.Sort(matches)

		matchesConRol := make([]modeloDocenteSiu, 0, len(matches))
		for _, m := range matches {
			m := modeloDocenteSiu{
				nombre: m.Target,
				rol:    pd.siu[m.Target],
			}

			matchesConRol = append(matchesConRol, m)
		}

		d := modeloDocenteDb{
			codigo:  cod,
			nombre:  nombre,
			matches: matchesConRol,
		}

		docentes = append(docentes, d)
	}

	slices.SortFunc(docentes, func(d1, d2 modeloDocenteDb) int {
		return cmp.Compare(d1.nombre, d2.nombre)
	})

	return docentes
}

func (m modeloApp) Init() tea.Cmd {
	log.Default().WithPrefix("🎨").Info("iniciando interfaz gráfica")
	return nil
}

func (td modeloTabsDocentes) Init() tea.Cmd {
	return nil
}

func (m modeloApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.screenWidth = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case ">":
			if m.currMateria < len(m.pagsMaterias)-1 {
				m.currMateria++
			} else {
				m.currMateria = 0
			}
			return m, nil
		case "<":
			if m.currMateria > 0 {
				m.currMateria--
			} else {
				m.currMateria = len(m.pagsMaterias) - 1
			}
			return m, nil
		}
	}

	pagMateria := m.pagsMaterias[m.currMateria]
	pagMateria.tabsDocentes, cmd = pagMateria.tabsDocentes.Update(msg)

	return m, cmd
}

func (m *modeloTabsDocentes) avanzarDocente() {
	m.currDocenteDb++
}

func (m modeloTabsDocentes) Update(msg tea.Msg) (modeloTabsDocentes, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "right", "ctrl+f":
			m.avanzarDocente()
			return m, nil
		case "left", "ctrl+b":
			m.currDocenteDb--
			return m, nil
		}
	}

	return m, cmd
}

var (
	colorFIUBA              = lipgloss.Color("#4eacd4")
	styleTituloMateria      = lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.NormalBorder())
	styleTabDocente         = lipgloss.NewStyle().Padding(0, 1)
	styleTabDocenteActivo   = styleTabDocente.Foreground(lipgloss.ANSIColor(0)).Background(colorFIUBA)
	styleTabDocenteInactivo = styleTabDocente.Faint(true)
)

func (m modeloApp) View() string {
	pagMateria := m.pagsMaterias[m.currMateria]

	s := strings.Builder{}

	s.WriteString(styleTituloMateria.Render(pagMateria.tituloMateria()))
	s.WriteString("\n")

	border := lipgloss.NormalBorder()
	border.BottomLeft = "├"
	border.BottomRight = "┴"

	if m.screenWidth > 0 {
		s.WriteString(lipgloss.
			NewStyle().
			Padding(0, 1).
			Border(border).
			Render("DOCENTES REGISTRADOS"))
		s.WriteString(strings.Repeat(lipgloss.NormalBorder().Top, m.screenWidth-len("DOCENTES REGISTRADOS")-5))
		s.WriteString("┐")
		s.WriteString("\n")
	}

	tabsDocentes := strings.Builder{}
	lineWidth := 0

	for i, d := range pagMateria.tabsDocentes.docentesDb {
		var style lipgloss.Style

		if i == pagMateria.tabsDocentes.currDocenteDb {
			style = styleTabDocenteActivo
		} else {
			style = styleTabDocenteInactivo
		}

		if lineWidth+len(d.nombre)+2 >= m.screenWidth-2 {
			tabsDocentes.WriteString("\n")
			lineWidth = 0
		}

		tabsDocentes.WriteString(style.Padding(0, 1).Render(d.nombre))
		lineWidth += 2 + len(d.nombre)
	}

	s.WriteString(lipgloss.NewStyle().Width(m.screenWidth-2).Border(lipgloss.NormalBorder(), false, true, true, true).Render(tabsDocentes.String()))
	s.WriteString("\n")

	return pagMateria.tabsDocentes.View()
}

func (td modeloTabsDocentes) View() string {
	return fmt.Sprintf("%v", td.currDocenteDb)
}
