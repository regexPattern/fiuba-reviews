package actualizador

import (
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PatchCatedrasOutput struct{}

func newPatchCatedras(_ *log.Logger, _ *pgxpool.Conn, _ ultimaOfertaMateria) (PatchCatedrasOutput, error) {
	return PatchCatedrasOutput{}, nil
}

type PatchMateriaInput struct{}
