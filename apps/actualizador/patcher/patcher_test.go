package patcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapNombreCodigoMateriasDb(t *testing.T) {
	materias := []MateriaBD{
		{Codigo: "COD001", Nombre: "Análisis Matemático II"},
		{Codigo: "COD002", Nombre: "Álgebra Lineal"},
	}

	codigos := mapNombreCodigo(materias)

	assert.Len(t, codigos, 2)
	assert.Equal(t, codigos["analisis matematico ii"], "COD001")
	assert.Equal(t, codigos["algebra lineal"], "COD002")
}
