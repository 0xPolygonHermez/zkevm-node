package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func versionCmd(*cli.Context) error {
	fmt.Printf("Version = \"%v\"\n", version)
	fmt.Printf("Build = \"%v\"\n", commit)
	fmt.Printf("Date = \"%v\"\n", date)
	return nil
}
