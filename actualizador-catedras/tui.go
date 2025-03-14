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
	"github.com/charmbracelet/log"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type model struct {
	materias    []viewMateria
	currMateria int
	screenWidth int
}

type viewMateria struct {
	codigo               string
	nombre               string
	docentesDb           viewSincronDocentesDb
	docentesNuevos       viewDocentesNuevos
	viewDocentesDbActiva bool
}

type viewSincronDocentesDb struct {
	docentes            []modeloDocenteDb
	matches             [][]modeloDocenteSiu
	currDocente         int
	currMatchesDocentes []int
}

type viewDocentesNuevos struct {
	docentes    []modeloDocenteSiu
	currDocente int
}

type modeloDocenteDb struct {
	codigo string
	nombre string
}

type modeloDocenteSiu struct {
	nombre string
	rol    string
}

func newModeloApp(patches []patch) model {
	pagsMaterias := make([]viewMateria, 0, len(patches))

	for _, p := range patches {
		docentesDb, docentesSiu := newViewsDocentes(*p.docentes)

		pm := viewMateria{
			codigo:               p.codigoMateria,
			nombre:               p.nombreMateria,
			docentesDb:           docentesDb,
			docentesNuevos:       docentesSiu,
			viewDocentesDbActiva: true,
		}

		pagsMaterias = append(pagsMaterias, pm)
	}

	slices.SortFunc(pagsMaterias, func(m1, m2 viewMateria) int {
		return cmp.Compare(m1.codigo, m2.codigo)
	})

	log.Debugf("sincronizando docentes de %v materias", len(pagsMaterias))

	modelo := model{
		materias:    pagsMaterias,
		currMateria: 0,
	}

	return modelo
}

func newViewsDocentes(pd patchDocentes) (viewSincronDocentesDb, viewDocentesNuevos) {
	docentesDb := make([]modeloDocenteDb, 0, len(pd.db))
	matchesDocentesDb := make([][]modeloDocenteSiu, 0, len(pd.db))

	nombresDocentesDb := slices.Collect(maps.Keys(pd.db))
	nombresDocentesSiu := slices.Collect(maps.Keys(pd.siu))

	slices.Sort(nombresDocentesDb)

	for _, nombre := range nombresDocentesDb {
		cod := pd.db[nombre]

		matches := fuzzy.RankFind(nombre, nombresDocentesSiu)

		// No tiene mucho sentido vincular un docente del SIU a un docente de
		// la base de datos cuyo nombre no tiene similitud con este. El
		// algoritmos de la librería de fuzzy matching usada es bastante
		// permisivo, así que si no hay matches es porque definitivamente no
		// hay matches.

		if len(matches) == 0 {
			continue
		}

		sort.Sort(matches)
		matchesDocentesSiu := make([]modeloDocenteSiu, 0, len(matches))

		for _, m := range matches {
			m := modeloDocenteSiu{
				nombre: m.Target,
				rol:    pd.siu[m.Target],
			}

			matchesDocentesSiu = append(matchesDocentesSiu, m)
		}

		d := modeloDocenteDb{
			codigo: cod,
			nombre: nombre,
		}

		docentesDb = append(docentesDb, d)
		matchesDocentesDb = append(matchesDocentesDb, matchesDocentesSiu)
	}

	viewSincronDocentesDb := viewSincronDocentesDb{
		docentes:            docentesDb,
		matches:             matchesDocentesDb,
		currDocente:         0,
		currMatchesDocentes: make([]int, len(docentesDb)),
	}

	viewDocentesNuevos := viewDocentesNuevos{}

	return viewSincronDocentesDb, viewDocentesNuevos
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

	m.materias[m.currMateria], cmd = m.materias[m.currMateria].Update(msg)

	return m, cmd
}

func (m viewMateria) Update(msg tea.Msg) (viewMateria, tea.Cmd) {
	var cmd tea.Cmd

	if m.viewDocentesDbActiva {
		m.docentesDb, cmd = m.docentesDb.Update(msg)
	} else {
		m.docentesNuevos, cmd = m.docentesNuevos.Update(msg)
	}

	return m, cmd
}

func (sd viewSincronDocentesDb) Update(msg tea.Msg) (viewSincronDocentesDb, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "right":
			if sd.currDocente+1 < len(sd.docentes) {
				sd.currDocente++
			} else {
				sd.currDocente = 0
			}
			return sd, nil
		case "left":
			if sd.currDocente-1 < 0 {
				sd.currDocente = len(sd.docentes) - 1
			} else {
				sd.currDocente--
			}
			return sd, nil
		case "up":
			matchesDocente := len(sd.matches[sd.currDocente])
			if sd.currMatchesDocentes[sd.currDocente]-1 >= 0 {
				sd.currMatchesDocentes[sd.currDocente]--
			} else {
				sd.currMatchesDocentes[sd.currDocente] = matchesDocente - 1
			}
			return sd, nil
		case "down":
			matchesDocente := len(sd.matches[sd.currDocente])
			if sd.currMatchesDocentes[sd.currDocente]+1 < matchesDocente {
				sd.currMatchesDocentes[sd.currDocente]++
			} else {
				sd.currMatchesDocentes[sd.currDocente] = 0
			}
			return sd, nil
		}
	}

	return sd, cmd
}

func (m viewDocentesNuevos) Update(msg tea.Msg) (viewDocentesNuevos, tea.Cmd) {
	return m, nil
}

var (
	colorFIUBA              = lipgloss.Color("#4eacd4")
	styleTituloMateria      = lipgloss.NewStyle()
	styleTabDocente         = lipgloss.NewStyle().Padding(0, 1)
	styleTabDocenteActivo   = styleTabDocente.Foreground(lipgloss.ANSIColor(0)).Background(colorFIUBA)
	styleTabDocenteInactivo = styleTabDocente.Faint(true)
)

func (m model) View() string {
	pm := m.materias[m.currMateria]

	s := strings.Builder{}

	s.WriteString(styleTituloMateria.
		Render(fmt.Sprintf("%v • %v", pm.codigo, strings.ToUpper(pm.nombre))))
	s.WriteString("\n")

	s.WriteString(pm.View(m.screenWidth))
	s.WriteString("\n")

	return s.String()
}

func (m viewMateria) View(screenWidth int) string {
	s := strings.Builder{}

	if m.viewDocentesDbActiva {
		s.WriteString(m.docentesDb.View(screenWidth))
	} else {
		s.WriteString(m.docentesNuevos.View(screenWidth))
	}

	return s.String()
}

func (sd viewSincronDocentesDb) View(screenWidth int) string {
	s := strings.Builder{}

	lineWidth := 0

	_, rightPad, _, leftPad := styleTabDocente.GetPadding()
	pad := leftPad + rightPad

	for i, d := range sd.docentes {
		var style lipgloss.Style

		if i == sd.currDocente {
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

	for i, d := range sd.matches[sd.currDocente] {
		styles := lipgloss.NewStyle()
		s.WriteString("[ ] ")

		if i == sd.currMatchesDocentes[sd.currDocente] {
			styles = styles.Underline(true)
		}

		s.WriteString(styles.Render(fmt.Sprintf("%v - %v", d.nombre, d.rol)))
		s.WriteString("\n")
	}

	return lipgloss.NewStyle().Width(screenWidth).Render(s.String())
}

func (m viewDocentesNuevos) View(_ int) string {
	return ""
}
