package patch

import (
	"context"
)

type PatchProposal struct {}

type PatchResolution struct{}

func ApplyPatches(ctx context.Context, resolved []PatchResolution) error {
	return nil
}
