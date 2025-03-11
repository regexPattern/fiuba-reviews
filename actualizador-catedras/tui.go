package main

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type modeloApp struct {
	materias         []modeloMateria
	idxMateriaActual int
	width            int
}

type modeloMateria struct {
	codigo             string
	nombre             string
	docentesDb         []modeloDocenteDb
	idxDocenteDbActual int
	docentesSiu        []modeloDocenteSiu
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

func (m modeloMateria) tituloMateria() string {
	return fmt.Sprintf("%v • %v", m.codigo, strings.ToUpper(m.nombre))
}

func newModeloApp(patches []patch) modeloApp {
	materias := make([]modeloMateria, 0, len(patches))

	for _, p := range patches {
		docentesDb := newModelosDocentesDb(*p.docentes)

		m := modeloMateria{
			codigo:             p.codigoMateria,
			nombre:             p.nombreMateria,
			docentesDb:         docentesDb,
			idxDocenteDbActual: 0,
		}

		materias = append(materias, m)
	}

	slices.SortFunc(materias, func(r1, r2 modeloMateria) int {
		return cmp.Compare(r1.codigo, r2.codigo)
	})

	model := modeloApp{
		materias:         materias,
		idxMateriaActual: 0,
	}

	return model
}

func newModelosDocentesDb(pd patchDocentes) []modeloDocenteDb {
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
	// log.Default().WithPrefix("🎨").Info("iniciando interfaz gráfica")
	return nil
}

func (m modeloApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "right", "ctrl+n":
			if m.idxMateriaActual < len(m.materias)-1 {
				m.idxMateriaActual++
			} else {
				m.idxMateriaActual = 0
			}
			return m, nil
		case "left", "ctrl+p":
			if m.idxMateriaActual > 0 {
				m.idxMateriaActual--
			} else {
				m.idxMateriaActual = len(m.materias) - 1
			}
			return m, nil
		case "tab", "ctrl+f":
			m.materias[m.idxMateriaActual].idxDocenteDbActual++
			return m, nil
		case "shift+tab", "ctrl+b":
			m.materias[m.idxMateriaActual].idxDocenteDbActual--
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
	materia := m.materias[m.idxMateriaActual]

	s := strings.Builder{}

	s.WriteString(styleTituloMateria.Render(materia.tituloMateria()))
	s.WriteString("\n")

	border := lipgloss.NormalBorder()
	border.BottomLeft = "├"
	border.BottomRight = "┴"

	if m.width > 0 {
		s.WriteString(lipgloss.
			NewStyle().
			Padding(0, 1).
			Border(border).
			Render("DOCENTES REGISTRADOS"))
		s.WriteString(strings.Repeat(lipgloss.NormalBorder().Top, m.width-len("DOCENTES REGISTRADOS")-5))
		s.WriteString("┐")
		s.WriteString("\n")
	}

	tabsDocentes := strings.Builder{}
	lineWidth := 0

	for i, d := range materia.docentesDb {
		var style lipgloss.Style

		if i == materia.idxDocenteDbActual {
			style = styleTabDocenteActivo
		} else {
			style = styleTabDocenteInactivo
		}

		if lineWidth+len(d.nombre)+2 >= m.width-2 {
			tabsDocentes.WriteString("\n")
			lineWidth = 0
		}

		tabsDocentes.WriteString(style.Padding(0, 1).Render(d.nombre))
		lineWidth += 2 + len(d.nombre)
	}

	s.WriteString(lipgloss.NewStyle().Width(m.width-2).Border(lipgloss.NormalBorder(), false, true, true, true).Render(tabsDocentes.String()))
	s.WriteString("\n")

	return s.String()
}
