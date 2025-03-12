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
	pagsMaterias []modeloPagMateria
	currMateria  int
	screenWidth  int
}

type modeloPagMateria struct {
	codigo            string
	nombre            string
	docentesDb        modeloDocentesDb
	docentesSiu       modeloDocentesSiu
	docentesDbActivos bool
}

type modeloDocentesDb struct {
	docentes    []docenteDb
	matches     [][]docenteSiu
	currDocente int
}

type modeloDocentesSiu struct {
	docentes    []docenteSiu
	currDocente int
}

type docenteDb struct {
	codigo string
	nombre string
}

type docenteSiu struct {
	nombre string
	rol    string
}

func newModeloApp(patches []patch) modeloApp {
	pagsMaterias := make([]modeloPagMateria, 0, len(patches))

	for _, p := range patches {
		docentesDb, docentesSiu := newModelosDocentes(*p.docentes)

		pm := modeloPagMateria{
			codigo:            p.codigoMateria,
			nombre:            p.nombreMateria,
			docentesDb:        docentesDb,
			docentesSiu:       docentesSiu,
			docentesDbActivos: true,
		}

		pagsMaterias = append(pagsMaterias, pm)
	}

	slices.SortFunc(pagsMaterias, func(m1, m2 modeloPagMateria) int {
		return cmp.Compare(m1.codigo, m2.codigo)
	})

	modelo := modeloApp{
		pagsMaterias: pagsMaterias,
		currMateria:  0,
	}

	return modelo
}

func newModelosDocentes(pd patchDocentes) (modeloDocentesDb, modeloDocentesSiu) {
	docentesDb := make([]docenteDb, 0, len(pd.db))
	matches := make([][]docenteSiu, 0, len(pd.db))

	nombresDocentesDb := slices.Collect(maps.Keys(pd.db))
	nombresDocentesSiu := slices.Collect(maps.Keys(pd.siu))

	slices.Sort(nombresDocentesDb)

	for _, nombre := range nombresDocentesDb {
		cod := pd.db[nombre]

		rankedMatches := fuzzy.RankFind(nombre, nombresDocentesSiu)

		// No tiene mucho sentido vincular un docente del SIU a un docente de
		// la base de datos cuyo nombre no tiene similitud con este. El
		// algoritmos de la librería de fuzzy matching usada es bastante
		// permisivo, así que si no hay matches es porque definitivamente no
		// hay matches.

		if len(rankedMatches) == 0 {
			continue
		}

		sort.Sort(rankedMatches)
		matchesDocentesSiu := make([]docenteSiu, 0, len(rankedMatches))

		for _, m := range rankedMatches {
			m := docenteSiu{
				nombre: m.Target,
				rol:    pd.siu[m.Target],
			}

			matchesDocentesSiu = append(matchesDocentesSiu, m)
		}

		d := docenteDb{
			codigo: cod,
			nombre: nombre,
		}

		docentesDb = append(docentesDb, d)
		matches = append(matches, matchesDocentesSiu)
	}

	modeloDb := modeloDocentesDb{
		docentes:    docentesDb,
		matches:     matches,
		currDocente: 0,
	}

	modeloSiu := modeloDocentesSiu{}

	return modeloDb, modeloSiu
}

func (m modeloApp) Init() tea.Cmd {
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

	m.pagsMaterias[m.currMateria], cmd = m.pagsMaterias[m.currMateria].Update(msg)

	return m, cmd
}

func (m modeloPagMateria) Update(msg tea.Msg) (modeloPagMateria, tea.Cmd) {
	var cmd tea.Cmd

	if m.docentesDbActivos {
		m.docentesDb, cmd = m.docentesDb.Update(msg)
	} else {
		m.docentesSiu, cmd = m.docentesSiu.Update(msg)
	}

	return m, cmd
}

func (m modeloDocentesDb) Update(msg tea.Msg) (modeloDocentesDb, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "right":
			if m.currDocente+1 < len(m.docentes) {
				m.currDocente++
			} else {
				m.currDocente = 0
			}
			return m, nil
		case "left":
			if m.currDocente-1 < 0 {
				m.currDocente = len(m.docentes) - 1
			} else {
				m.currDocente--
			}
			return m, nil
		}
	}

	return m, cmd
}

func (m modeloDocentesSiu) Update(msg tea.Msg) (modeloDocentesSiu, tea.Cmd) {
	return m, nil
}

var (
	colorFIUBA              = lipgloss.Color("#4eacd4")
	styleTituloMateria      = lipgloss.NewStyle()
	styleTabDocente         = lipgloss.NewStyle().Padding(0, 1)
	styleTabDocenteActivo   = styleTabDocente.Foreground(lipgloss.ANSIColor(0)).Background(colorFIUBA)
	styleTabDocenteInactivo = styleTabDocente.Faint(true)
)

func (m modeloApp) View() string {
	pm := m.pagsMaterias[m.currMateria]

	s := strings.Builder{}

	s.WriteString(styleTituloMateria.
		Render(fmt.Sprintf("%v • %v", pm.codigo, strings.ToUpper(pm.nombre))))
	s.WriteString("\n")

	s.WriteString(pm.View(m.screenWidth))
	s.WriteString("\n")

	return s.String()
}

func (m modeloPagMateria) View(screenWidth int) string {
	s := strings.Builder{}

	if m.docentesDbActivos {
		s.WriteString(m.docentesDb.View(screenWidth))
	} else {
		s.WriteString(m.docentesSiu.View(screenWidth))
	}

	return s.String()
}

func (m modeloDocentesDb) View(screenWidth int) string {
	s := strings.Builder{}

	lineWidth := 0

	_, rightPad, _, leftPad := styleTabDocente.GetPadding()
	pad := leftPad + rightPad

	for i, d := range m.docentes {
		var style lipgloss.Style

		if i == m.currDocente {
			style = styleTabDocenteActivo
		} else {
			style = styleTabDocenteInactivo
		}

		if lineWidth+len(d.nombre)+pad > screenWidth {
			s.WriteString("\n")
			lineWidth = 0
		}

		s.WriteString(style.Padding(0, 1).Render(d.nombre))
		lineWidth += pad + len(d.nombre)
	}

	s.WriteString("\n")

	for _, d := range m.matches[m.currDocente] {
		s.WriteString(fmt.Sprintf("• %v - %v", d.nombre, d.rol))
		s.WriteString("\n")
	}

	return lipgloss.NewStyle().Width(screenWidth).Render(s.String())
}

func (m modeloDocentesSiu) View(_ int) string {
	return ""
}
