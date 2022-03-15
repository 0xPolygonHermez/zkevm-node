package main

import (
	"github.com/hermeznetwork/hermez-core/scripts/cmd/compilesc"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

func compileSC(ctx *cli.Context) error {
	aferoFs := afero.NewOsFs()

	return compilesc.NewManager(aferoFs).Run(ctx.String(flagInput))
}
