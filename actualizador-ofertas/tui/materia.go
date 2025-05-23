package tui

import (
	"cmp"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/regexPattern/fiuba-reviews/actualizador-ofertas/actualizador"
)

type materiaView struct {
	codigo                string
	nombre                string
	docentesRegistrados   docentesRegistradosView
	docentesNuevos        docentesNuevosView
	viewRegistradosActiva bool
}

func newMateriaViews(patches []actualizador.PatchMateriaOutput) []materiaView {
	mats := make([]materiaView, 0, len(patches))

	for _, p := range patches {
		regs, nuevos := newDocentesViews(p.Docentes)
		m := materiaView{
			codigo:                p.Codigo,
			nombre:                p.Nombre,
			docentesRegistrados:   regs,
			docentesNuevos:        nuevos,
			viewRegistradosActiva: true,
		}
		mats = append(mats, m)
	}

	slices.SortFunc(mats, func(m1, m2 materiaView) int {
		return cmp.Compare(m1.codigo, m2.codigo)
	})

	return mats
}

func (v materiaView) Update(msg tea.Msg) (materiaView, tea.Cmd) {
	var cmd tea.Cmd

	if v.viewRegistradosActiva {
		v.docentesRegistrados, cmd = v.docentesRegistrados.Update(msg)
	} else {
		v.docentesNuevos, cmd = v.docentesNuevos.Update(msg)
	}

	return v, cmd
}

func (v materiaView) View(screenWidth int) string {
	s := strings.Builder{}

	if v.viewRegistradosActiva {
		s.WriteString(v.docentesRegistrados.View(screenWidth))
	} else {
		s.WriteString(v.docentesNuevos.View(screenWidth))
	}

	return s.String()
}
