package main

import (
	"os"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/urfave/cli/v2"
)

func genJSONSchema(cli *cli.Context) error {
	generator := config.NewConfigJsonSchemaGenerater()
	schema, err := generator.GenerateJsonSchema(cli)
	if err != nil {
		return err
	}
	file, err := generator.SerializeJsonSchema(schema)
	if err != nil {
		return err
	}
	output := cli.String(config.FlagOutputFile)
	err = os.WriteFile(output, file, 0600) //nolint:gomnd
	if err != nil {
		return err
	}
	return nil
}
