package main

import (
	"github.com/hermeznetwork/hermez-core/scripts/cmd/dependencies"
	"github.com/urfave/cli/v2"
)

func updateDeps(ctx *cli.Context) error {
	return dependencies.NewManager().Run()
}
