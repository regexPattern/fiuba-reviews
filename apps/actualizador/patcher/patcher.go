package patcher

import (
	"context"
	"fmt"
)

type PatchProposal struct {
	MateriasSiu []LatestMateria
	MateriasDb  []Materia
}

type PatchResolution struct{}

func GeneratePatches(ctx context.Context) (*PatchProposal, error) {
	if err := initS3Client(ctx); err != nil {
		return nil, err
	}

	ofertas, err := getOfertasFromSiu(ctx)
	if err != nil {
		return nil, err
	}

	materiasSiu := mergeLatestOfertas(ofertas)
	materiasDb, err := getMateriasFromDb(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println(materiasDb)

	pp := PatchProposal{
		MateriasSiu: materiasSiu,
		MateriasDb:  materiasDb,
	}

	return &pp, nil
}

func ApplyPatches(ctx context.Context, resolved []PatchResolution) error {
	return nil
}
