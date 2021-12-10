package main

import (
	"fmt"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/urfave/cli/v2"
)

var (
	version = "v1.0.1"
	commit  = "dev"
	date    = ""
)

func cmdVersion(*cli.Context) error {
	log.Debugf("version command executed")

	fmt.Printf("Version = \"%v\"\n", version)
	fmt.Printf("Build = \"%v\"\n", commit)
	fmt.Printf("Date = \"%v\"\n", date)
	return nil
}
