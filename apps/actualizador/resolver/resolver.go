package resolver

import (
	"cmp"
	"context"
	"log/slog"
	"os"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackc/pgx/v5"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/indexador"
)

var conn *pgx.Conn

func ResolverActualizaciones(inputConn *pgx.Conn, materias []indexador.Materia) error {
	conn, _ = pgx.Connect(context.TODO(), os.Getenv("DATABASE_URL"))

	if len(materias) == 0 {
		slog.Info("no hay materias por actualizar")
		return nil
	}

	sortSegunPrioridad(materias)

	p := tea.NewProgram(newModel(materias))
	_, err := p.Run()
	return err
}

func sortSegunPrioridad(materias []indexador.Materia) {
	nDocentes := make(map[string]int, len(materias))
	for _, m := range materias {
		docentesUnicos := make(map[string]bool)
		for _, c := range m.MateriaSiu.Catedras {
			for _, d := range c.Docentes {
				docentesUnicos[d.Nombre] = true
			}
		}
		nDocentes[m.MateriaSiu.Codigo] = len(docentesUnicos)
	}

	slices.SortFunc(materias, func(a, b indexador.Materia) int {
		return cmp.Compare(
			nDocentes[b.MateriaDb.Codigo],
			nDocentes[a.MateriaDb.Codigo],
		)
	})
}
