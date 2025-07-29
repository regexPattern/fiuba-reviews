package patch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparacionCuatrimestres(t *testing.T) {
	assert.True(t, Cuatri{1, 2025}.despuesDe(Cuatri{2, 2023}))
	assert.False(t, Cuatri{2, 2023}.despuesDe(Cuatri{1, 2025}))
	assert.False(t, Cuatri{1, 2025}.despuesDe(Cuatri{1, 2025}))
}
