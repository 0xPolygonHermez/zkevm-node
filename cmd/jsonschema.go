package main

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/invopop/jsonschema"
	"github.com/urfave/cli/v2"
)

func genJSONSchema(cli *cli.Context) error {
	r := new(jsonschema.Reflector)
	repoName := "github.com/0xPolygonHermez/zkevm-node"
	r.Namer = func(rt reflect.Type) string {
		return strings.TrimLeft(rt.PkgPath(), repoName) + "_" + rt.Name()
	}
	r.ExpandedStruct = true
	r.DoNotReference = true
	if err := r.AddGoComments(repoName, "./"); err != nil {
		return err
	}
	schema := r.Reflect(&config.Config{})
	schema.ID = jsonschema.ID(repoName + "config/config")
	file, err := json.MarshalIndent(schema, "", "\t")
	if err != nil {
		return err
	}
	output := cli.String(config.FlagOutputFile)
	err = ioutil.WriteFile(output, file, 0644)
	if err != nil {
		return err
	}
	return nil
}
