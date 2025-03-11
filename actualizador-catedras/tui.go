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

type app struct {
	pagsMaterias []pagMateria
	currMateria  int
	screenWidth  int
	help         help.Model
}

type pagMateria struct {
	codigo       string
	nombre       string
	tabsDocentes tabsDocentes
}

type tabsDocentes struct {
	docentesDb    []docenteDb
	currDocenteDb int
	docentesSiu   []docenteSiu
}

type docenteDb struct {
	codigo  string
	nombre  string
	matches []docenteSiu
}

type docenteSiu struct {
	nombre string
	rol    string
}

func (pm pagMateria) tituloMateria() string {
	return fmt.Sprintf("%v • %v", pm.codigo, strings.ToUpper(pm.nombre))
}

func newApp(patches []patch) app {
	materias := make([]pagMateria, 0, len(patches))

	for _, p := range patches {
		docentesDb := newDocentesDb(*p.docentes)

		pm := pagMateria{
			codigo: p.codigoMateria,
			nombre: p.nombreMateria,
			tabsDocentes: tabsDocentes{
				docentesDb:    docentesDb,
				currDocenteDb: 0,
			},
		}

		materias = append(materias, pm)
	}

	slices.SortFunc(materias, func(r1, r2 pagMateria) int {
		return cmp.Compare(r1.codigo, r2.codigo)
	})

	model := app{
		pagsMaterias: materias,
		currMateria:  0,
	}

	return model
}

func newDocentesDb(pd patchDocentes) []docenteDb {
	docentes := make([]docenteDb, 0, len(pd.db))
	nombresDocentesSiu := slices.Collect(maps.Keys(pd.siu))

	for nombre, cod := range pd.db {
		matches := fuzzy.RankFind(nombre, nombresDocentesSiu)
		if len(matches) == 0 {
			continue
		}

		sort.Sort(matches)

		matchesConRol := make([]docenteSiu, 0, len(matches))
		for _, m := range matches {
			m := docenteSiu{
				nombre: m.Target,
				rol:    pd.siu[m.Target],
			}

			matchesConRol = append(matchesConRol, m)
		}

		d := docenteDb{
			codigo:  cod,
			nombre:  nombre,
			matches: matchesConRol,
		}

		docentes = append(docentes, d)
	}

	slices.SortFunc(docentes, func(d1, d2 docenteDb) int {
		return cmp.Compare(d1.nombre, d2.nombre)
	})

	return docentes
}

func (a app) Init() tea.Cmd {
	log.Default().WithPrefix("🎨").Info("iniciando interfaz gráfica")
	return nil
}

func (td tabsDocentes) Init() tea.Cmd {
	return nil
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.screenWidth = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case ">":
			if a.currMateria < len(a.pagsMaterias)-1 {
				a.currMateria++
			} else {
				a.currMateria = 0
			}
			return a, nil
		case "<":
			if a.currMateria > 0 {
				a.currMateria--
			} else {
				a.currMateria = len(a.pagsMaterias) - 1
			}
			return a, nil
		}
	}

	pagMateria := a.pagsMaterias[a.currMateria]
	pagMateria.tabsDocentes, cmd = pagMateria.tabsDocentes.Update(msg)

	return a, cmd
}

func (td *tabsDocentes) avanzarDocente() {
	td.currDocenteDb++
}

func (td tabsDocentes) Update(msg tea.Msg) (tabsDocentes, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "right", "ctrl+f":
			td.avanzarDocente()
			return td, nil
		case "left", "ctrl+b":
			td.currDocenteDb--
			return td, nil
		}
	}

	return td, cmd
}

var (
	colorFIUBA              = lipgloss.Color("#4eacd4")
	styleTituloMateria      = lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.NormalBorder())
	styleTabDocente         = lipgloss.NewStyle().Padding(0, 1)
	styleTabDocenteActivo   = styleTabDocente.Foreground(lipgloss.ANSIColor(0)).Background(colorFIUBA)
	styleTabDocenteInactivo = styleTabDocente.Faint(true)
)

func (a app) View() string {
	pagMateria := a.pagsMaterias[a.currMateria]

	s := strings.Builder{}

	s.WriteString(styleTituloMateria.Render(pagMateria.tituloMateria()))
	s.WriteString("\n")

	border := lipgloss.NormalBorder()
	border.BottomLeft = "├"
	border.BottomRight = "┴"

	if a.screenWidth > 0 {
		s.WriteString(lipgloss.
			NewStyle().
			Padding(0, 1).
			Border(border).
			Render("DOCENTES REGISTRADOS"))
		s.WriteString(strings.Repeat(lipgloss.NormalBorder().Top, a.screenWidth-len("DOCENTES REGISTRADOS")-5))
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

		if lineWidth+len(d.nombre)+2 >= a.screenWidth-2 {
			tabsDocentes.WriteString("\n")
			lineWidth = 0
		}

		tabsDocentes.WriteString(style.Padding(0, 1).Render(d.nombre))
		lineWidth += 2 + len(d.nombre)
	}

	s.WriteString(lipgloss.NewStyle().Width(a.screenWidth-2).Border(lipgloss.NormalBorder(), false, true, true, true).Render(tabsDocentes.String()))
	s.WriteString("\n")

	return pagMateria.tabsDocentes.View()
}

func (td tabsDocentes) View() string {
	return fmt.Sprintf("%v", td.currDocenteDb)
}
