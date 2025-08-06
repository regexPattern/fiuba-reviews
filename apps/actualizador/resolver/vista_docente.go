package resolver

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackc/pgx/v5"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/indexador"
)

type vistaDocenteModel struct {
	indexador.Materia
	docente      patchDocente
	infoDocentes []infoDocente
	err          error
}

func newVistaDocente() vistaDocenteModel {
	return vistaDocenteModel{}
}

func (m vistaDocenteModel) Init() tea.Cmd {
	return nil
}

func (m vistaDocenteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case setMateriaMsg:
		m.setMateria(msg)
	}
	return m, nil
}

func (m vistaDocenteModel) View() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("%v • %v\n", m.Materia.MateriaSiu.Codigo, m.Materia.MateriaDb.Nombre))

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

func (m *vistaDocenteModel) setMateria(materia setMateriaMsg) tea.Cmd {
	m.Materia = indexador.Materia(materia)
	return func() tea.Msg {
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

func (m *vistaDocenteModel) setDocente(docente setDocenteMsg) {
	m.docente = patchDocente(docente)
}

func (m *vistaDocenteModel) setInfoDocentes(msg infoDocentesMsg) {
	if msg.err != nil {
		m.err = msg.err
	} else {
		m.infoDocentes = msg.payload
	}
}
