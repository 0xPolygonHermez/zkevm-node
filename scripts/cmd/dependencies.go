package main

import (
	"github.com/hermeznetwork/hermez-core/scripts/cmd/dependencies"
	"github.com/urfave/cli/v2"
)

func updateDeps(ctx *cli.Context) error {
	cfg := &dependencies.Config{
		Images: &dependencies.ImagesConfig{
			Names:          []string{"hermeznetwork/geth-zkevm-contracts"},
			TargetFilePath: "../../../docker-compose.yml",
		},
		PB: &dependencies.PBConfig{
			TargetDirPath: "../../../proto/src",
			SourceRepo:    "git@github.com:hermeznetwork/comms-protocol.git",
		},
		TV: &dependencies.TVConfig{
			TargetDirPath: "../../../test/vectors/src",
			SourceRepo:    "git@github.com:hermeznetwork/test-vectors.git",
		},
	}

	return dependencies.NewManager(cfg).Run()
}
