package main

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var (
	sTituloMateria    = lipgloss.NewStyle().Bold(true).Underline(true)
	sTabDocente       = lipgloss.NewStyle().Padding(0, 1)
	sTabDocenteActivo = sTabDocente.Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#4eacd4"))
)

type model struct {
	codigosMaterias  []string
	currMateria      int
	patches          map[string]patch
	currDocente      map[string]int
	listaDocentesSiu list.Model
}

func NewModel(patches []patch) model {
	numDocentes := 0
	patchesMap := make(map[string]patch, len(patches))

	for _, a := range patches {
		numDocentes += len(a.docentes.db)
		patchesMap[a.codigoMateria] = a
		if len(a.docentes.siu) == 0 {
			log.Error("WFT", "codigo", a.codigoMateria)
		}
	}

	currDocentes := make(map[string]int, numDocentes)

	return model{
		codigosMaterias: slices.Collect(maps.Keys(patchesMap)),
		currMateria:     0,
		patches:         patchesMap,
		currDocente:     currDocentes,
	}
}

func (m model) Init() tea.Cmd {
	log.Default().WithPrefix("🎨").Info("iniciando interfaz gráfica")
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "down", "ctrl+n":
			// TODO: manejar los indices
			m.currMateria++
			return m, nil
		case "up", "ctrl+p":
			m.currMateria++
			return m, nil
		case "right", "ctrl+f":
			m.currDocente[m.codigosMaterias[m.currMateria]]++
			return m, nil
		case "left", "ctrl+b":
			m.currDocente[m.codigosMaterias[m.currMateria]]--
			return m, nil
		}
	}

	return m, cmd
}

func (m model) View() string {
	codigoMateria := m.codigosMaterias[m.currMateria]
	patch := m.patches[codigoMateria]

	s := strings.Builder{}

	// PERF: claramente no es ideal hacer esto en cada render
	docentesMateria := slices.Collect(maps.Keys(patch.docentes.db))
	sort.Strings(docentesMateria)

	tabsDocentes := make([]string, 0, len(patch.docentes.db))

	for i, nombre := range docentesMateria {
		styles := sTabDocente

		if i == m.currDocente[codigoMateria] {
			styles = sTabDocenteActivo
		}

		tabsDocentes = append(tabsDocentes, styles.Render(nombre))
	}

	s.WriteString("\n")
	s.WriteString(sTituloMateria.Render(patch.codigoMateria))
	s.WriteString("\n")

	tabsRow := lipgloss.JoinHorizontal(lipgloss.Top, tabsDocentes...)

	s.WriteString("\n")
	s.WriteString(tabsRow)
	s.WriteString("\n")
	s.WriteString("\n")

	for nombre, rol := range patch.docentes.siu {
		s.WriteString(lipgloss.NewStyle().Padding(0, 1).Render(fmt.Sprintf("* %v - %v", nombre, rol)))
		s.WriteString("\n")
	}

	s.WriteString("\n")

	return s.String()
}
