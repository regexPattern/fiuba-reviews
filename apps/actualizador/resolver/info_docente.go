package resolver

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jackc/pgx/v5"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/indexador"
)

type infoDocenteModel struct {
	indexador.Materia
	docente      patchDocente
	infoDocentes []infoDocente
	spinner      spinner.Model
	err          error
}

func newInfoDocente() infoDocenteModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	return infoDocenteModel{
		spinner: s,
	}
}

func (m infoDocenteModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m infoDocenteModel) Update(msg tea.Msg) (infoDocenteModel, tea.Cmd) {
	switch msg := msg.(type) {
	case setMateriaMsg:
		m.setMateria(msg)
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m infoDocenteModel) View() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("%v • %v\n", m.Materia.MateriaSiu.Codigo, m.Materia.MateriaDb.Nombre))
	s.WriteString(fmt.Sprintf("%s\n", m.spinner.View()))

	if m.err != nil {
		s.WriteString(m.err.Error())
	} else {
		for i, d := range m.infoDocentes {
			if i > 10 {
				break
			}
			s.WriteString(fmt.Sprintf("%v - %v\n", d.Codigo, d.Nombre))
		}
	}

	return s.String()
}

type infoDocente struct {
	Codigo             string  `db:"codigo"`
	Nombre             string  `db:"nombre"`
	ResumenComentarios *string `db:"resumen_comentarios"`
}

type infoDocentesMsg struct {
	payload []infoDocente
	err     error
}

func (m *infoDocenteModel) setMateria(materia setMateriaMsg) tea.Cmd {
	m.Materia = indexador.Materia(materia)
	return func() tea.Msg {
		time.Sleep(time.Second * 3)
		rows, _ := conn.Query(
			context.Background(),
			"SELECT codigo, nombre, resumen_comentarios FROM docente WHERE codigo_materia = $1",
			materia.MateriaSiu.Codigo,
		)
		docentes, err := pgx.CollectRows(rows, pgx.RowToStructByName[infoDocente])
		return infoDocentesMsg{
			payload: docentes,
			err:     err,
		}
	}
}

func (m *infoDocenteModel) setDocente(docente setDocenteMsg) {
	m.docente = patchDocente(docente)
}

func (m *infoDocenteModel) setInfoDocentes(msg infoDocentesMsg) {
	if msg.err != nil {
		m.err = msg.err
	} else {
		m.infoDocentes = msg.payload
	}
}
