package main

import (
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/scripts/cmd/compilesc"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

func compileSC(ctx *cli.Context) error {
	aferoFs := afero.NewOsFs()

	err := compilesc.NewManager(aferoFs).Run(ctx.String(flagInput))
	if err != nil {
		log.Error("Failed to compile SCs, err: ", err)
	}
	return err
}
