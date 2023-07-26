package main

import (
	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/urfave/cli/v2"
)

func genJSONSchema(cli *cli.Context) error {
	file_config := cli.String(config.FlagDocumentationFileType)
	output := cli.String(config.FlagOutputFile)
	switch file_config {
	case NODE_CONFIGFILE:
		{
			generator := config.NewNodeConfigJsonSchemaGenerater()
			return generator.GenerateJsonSchemaAndWriteToFile(cli, output)
		}
	case NETWORK_CONFIGFILE:
		{
			generator := config.NewNetworkConfigJsonSchemaGenerater()
			return generator.GenerateJsonSchemaAndWriteToFile(cli, output)
		}
	default:
		panic("Not supported this config file: " + file_config)
	}
}
