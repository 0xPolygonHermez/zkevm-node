package main

import (
	"github.com/0xPolygonHermez/zkevm-node/test/scripts/cmd/dependencies"
	"github.com/urfave/cli/v2"
)

func updateDeps(ctx *cli.Context) error {
	cfg := &dependencies.Config{
		Images: &dependencies.ImagesConfig{
			Names:          []string{"hermeznetwork/geth-zkevm-contracts", "hermeznetwork/zkprover-local"},
			TargetFilePath: "../../../docker-compose.yml",
		},
		PB: &dependencies.PBConfig{
			TargetDirPath: "../../../proto/src",
			SourceRepo:    "https://github.com/0xPolygonHermez/zkevm-comms-protocol.git",
		},
		TV: &dependencies.TVConfig{
			TargetDirPath: "../../../test/vectors/src",
			SourceRepo:    "https://github.com/0xPolygonHermez/zkevm-testvectors.git",
		},
	}

	return dependencies.NewManager(cfg).Run()
}
