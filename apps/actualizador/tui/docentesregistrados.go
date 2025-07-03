package tui

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/regexPattern/fiuba-reviews/actualizador-ofertas/actualizador"
)

type docentesRegistradosView struct {
	docentes            []docenteDb
	matches             [][]docenteSiu
	currDocente         int
	currMatchesDocentes []int
}

type docenteDb struct {
	codigo string
	nombre string
}

type docenteSiu struct {
	nombre string
	rol    string
}

func (sd docentesRegistradosView) Update(msg tea.Msg) (docentesRegistradosView, tea.Cmd) {
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

func newDocentesViews(pd actualizador.PatchDocentesOutput) (docentesRegistradosView, docentesNuevosView) {
	dsDb := make([]docenteDb, 0, len(pd.Registrados))
	matchesDb := make([][]docenteSiu, 0, len(pd.Registrados))

	nombresDsDb := slices.Collect(maps.Keys(pd.Registrados))
	nombresDsSiu := slices.Collect(maps.Keys(pd.Nuevos))

	slices.Sort(nombresDsDb)

	for _, nom := range nombresDsDb {
		cod := pd.Registrados[nom]
		matches := fuzzy.RankFind(nom, nombresDsSiu)

		// No tiene mucho sentido vincular un docente del SIU a un docente de
		// la base de datos cuyo nombre no tiene similitud con este. El
		// algoritmos de la librería de fuzzy matching usada es bastante
		// permisivo, así que si no hay matches es porque definitivamente no
		// hay matches.

		if len(matches) == 0 {
			continue
		}

		sort.Sort(matches)

		matchesDocentesSiu := make([]docenteSiu, 0, len(matches))
		for _, m := range matches {
			m := docenteSiu{
				nombre: m.Target,
				rol:    pd.Nuevos[m.Target],
			}
			matchesDocentesSiu = append(matchesDocentesSiu, m)
		}

		d := docenteDb{
			codigo: cod,
			nombre: nom,
		}
		dsDb = append(dsDb, d)
		matchesDb = append(matchesDb, matchesDocentesSiu)
	}

	regs := docentesRegistradosView{
		docentes:            dsDb,
		matches:             matchesDb,
		currDocente:         0,
		currMatchesDocentes: make([]int, len(dsDb)),
	}

	nuevos := docentesNuevosView{}

	return regs, nuevos
}

var (
	styleTabDocente         = lipgloss.NewStyle().Padding(0, 1)
	styleTabDocenteActivo   = styleTabDocente.Foreground(lipgloss.ANSIColor(0)).Background(colorFIUBA)
	styleTabDocenteInactivo = styleTabDocente.Faint(true)
)

func (v docentesRegistradosView) View(screenWidth int) string {
	s := strings.Builder{}
	currLWidth := 0

	_, rightPad, _, leftPad := styleTabDocente.GetPadding()
	pad := leftPad + rightPad

	for i, d := range v.docentes {
		var style lipgloss.Style

		if i == v.currDocente {
			style = styleTabDocenteActivo
		} else {
			style = styleTabDocenteInactivo
		}

		if currLWidth+len(d.nombre)+pad > screenWidth {
			s.WriteString("\n")
			currLWidth = 0
		}

		s.WriteString(style.Padding(0, 1).Render(d.nombre))
		currLWidth += pad + len(d.nombre)
	}

	s.WriteString("\n")

	for i, d := range v.matches[v.currDocente] {
		styles := lipgloss.NewStyle()
		s.WriteString("[ ] ")

		if i == v.currMatchesDocentes[v.currDocente] {
			styles = styles.Underline(true)
		}

		s.WriteString(styles.Render(fmt.Sprintf("%v - %v", d.nombre, d.rol)))
		s.WriteString("\n")
	}

	return lipgloss.NewStyle().Width(screenWidth).Render(s.String())
}
