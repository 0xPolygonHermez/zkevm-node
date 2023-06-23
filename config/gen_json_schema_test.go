package config

import (
	"bytes"
	"errors"
	"testing"

	"github.com/invopop/jsonschema"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

type MySectionConfig struct {
}

/*
Config represents the configuration of the entire Hermez Node
The file is [TOML format](https://en.wikipedia.org/wiki/TOML#)

You could find some examples:
- `config/environments/local/local.node.config.toml`: running a permisionless node
- `config/environments/mainnet/public.node.config.toml`
- `config/environments/public/public.node.config.toml`
- `test/config/test.node.config.toml`: configuration for a trusted node used in CI
*/
type MyTestConfig struct {
	// F1 field description
	F1 string
	// F2 field description
	F2 int
}

type MyTestConfigWithJsonRenaming struct {
	F1 string `json:"f1_another_name"`
	F2 int    `json:"f2_another_name"`
}

type MyTestConfigWithMapstructureRenaming struct {
	F1 string `mapstructure:"f1_another_name"`
	F2 int    `mapstructure:"f2_another_name"`
}

type MyTestConfigWithMapstructureRenamingInSubStruct struct {
	F1 string
	F2 int
	F3 MyTestConfigWithMapstructureRenaming
}

func checkDefaultValue(t *testing.T, schema *jsonschema.Schema, key []string, expectedValue interface{}) {
	v, err := getValueFromSchema(schema, key)
	require.NoError(t, err)
	require.EqualValues(t, expectedValue, v.Default)
}

const MyTestConfigTomlFile = `
f1_another_name="value_f1"
f2_another_name=5678
`

func TestGenerateJsonSchemaCommentsWithMultiplesLines(t *testing.T) {
	cli := cli.NewContext(nil, nil, nil)
	generator := ConfigJsonSchemaGenerater[MyTestConfig]{
		repoName:                "github.com/0xPolygonHermez/zkevm-node/config/",
		cleanRequiredField:      true,
		addCodeCommentsToSchema: true,
		pathSourceCode:          "./",
		repoNameSuffix:          "config/",
		defaultValues: &MyTestConfig{
			F1: "defaultf1",
			F2: 1234,
		},
	}
	schema, err := generator.GenerateJsonSchema(cli)
	require.NoError(t, err)
	require.NotNil(t, schema)
	v, err := getValueFromSchema(schema, []string{"F2"})
	require.NoError(t, err)
	require.EqualValues(t, 1234, v.Default)
	require.NotEmpty(t, v.Description)

}

// This test is just to check what is the behaviour of reading a file
// when using tags `mapstructure` and `json`
func TestExploratoryForCheckReadFromFile(t *testing.T) {
	t.Skip("Is not a real test, just an exploratory one")
	viper.SetConfigType("toml")
	err := viper.ReadConfig(bytes.NewBuffer([]byte(MyTestConfigTomlFile)))
	require.NoError(t, err)

	var cfgJson MyTestConfigWithJsonRenaming
	err = viper.Unmarshal(&cfgJson, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	require.NoError(t, err)

	var cfgMapStructure MyTestConfigWithMapstructureRenaming
	err = viper.Unmarshal(&cfgMapStructure, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	require.NoError(t, err)

	require.EqualValues(t, cfgMapStructure.F1, cfgJson.F1)
	require.EqualValues(t, cfgMapStructure.F2, cfgJson.F2)
}

func TestGenerateJsonSchemaCustomWithNameChangingUsingMapsInSubFieldtrucutMustPanic(t *testing.T) {
	cli := cli.NewContext(nil, nil, nil)
	generator := ConfigJsonSchemaGenerater[MyTestConfigWithMapstructureRenamingInSubStruct]{
		repoName:                "mytest",
		cleanRequiredField:      true,
		addCodeCommentsToSchema: true,
		pathSourceCode:          "./",
		defaultValues: &MyTestConfigWithMapstructureRenamingInSubStruct{
			F1: "defaultf1",
			F2: 1234,
		},
	}
	//https://gophersnippets.com/how-to-test-a-function-that-panics
	t.Run("panics", func(t *testing.T) {
		// If the function panics, recover() will
		// return a non nil value.
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("function should panic")
			}
		}()

		generator.GenerateJsonSchema(cli)
	})
}

func TestGenerateJsonSchemaCustomWithNameChangingUsingMapstrucutMustPanic(t *testing.T) {
	cli := cli.NewContext(nil, nil, nil)
	generator := ConfigJsonSchemaGenerater[MyTestConfigWithMapstructureRenaming]{
		repoName:                "mytest",
		cleanRequiredField:      true,
		addCodeCommentsToSchema: true,
		pathSourceCode:          "./",
		defaultValues: &MyTestConfigWithMapstructureRenaming{
			F1: "defaultf1",
			F2: 1234,
		},
	}
	//https://gophersnippets.com/how-to-test-a-function-that-panics
	t.Run("panics", func(t *testing.T) {
		// If the function panics, recover() will
		// return a non nil value.
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("function should panic")
			}
		}()

		generator.GenerateJsonSchema(cli)
	})

}

// This case is a field that is mapped with another name in the json file
func TestGenerateJsonSchemaCustomWithNameChangingSetDefault(t *testing.T) {
	cli := cli.NewContext(nil, nil, nil)
	generator := ConfigJsonSchemaGenerater[MyTestConfigWithJsonRenaming]{
		repoName:                "mytest",
		cleanRequiredField:      true,
		addCodeCommentsToSchema: true,
		pathSourceCode:          "./",
		defaultValues: &MyTestConfigWithJsonRenaming{
			F1: "defaultf1",
			F2: 1234,
		},
	}

	schema, err := generator.GenerateJsonSchema(cli)
	require.NoError(t, err)
	require.NotNil(t, schema)

	checkDefaultValue(t, schema, []string{"f1_another_name"}, "defaultf1")
	checkDefaultValue(t, schema, []string{"f2_another_name"}, 1234)
}

func TestGenerateJsonSchemaCustomSetDefault(t *testing.T) {
	cli := cli.NewContext(nil, nil, nil)
	generator := ConfigJsonSchemaGenerater[MyTestConfig]{
		repoName:                "mytest",
		cleanRequiredField:      true,
		addCodeCommentsToSchema: true,
		pathSourceCode:          "./",
		defaultValues: &MyTestConfig{
			F1: "defaultf1",
			F2: 1234,
		},
	}

	schema, err := generator.GenerateJsonSchema(cli)
	require.NoError(t, err)
	require.NotNil(t, schema)
	checkDefaultValue(t, schema, []string{"F1"}, "defaultf1")

}

func TestGenerateJsonSchemaInjectDefaultValue1stLevel(t *testing.T) {
	cli := cli.NewContext(nil, nil, nil)
	generator := NewConfigJsonSchemaGenerater()
	generator.pathSourceCode = "../"
	generator.defaultValues.IsTrustedSequencer = false
	schema, err := generator.GenerateJsonSchema(cli)
	require.NoError(t, err)
	require.NotNil(t, schema)
	v, err := getValueFromSchema(schema, []string{"IsTrustedSequencer"})
	require.NoError(t, err)
	require.EqualValues(t, false, v.Default)
}

func TestGenerateJsonSchemaInjectDefaultValue2stLevel(t *testing.T) {
	cli := cli.NewContext(nil, nil, nil)
	generator := NewConfigJsonSchemaGenerater()
	generator.pathSourceCode = "../"
	// This is a hack, we are not at root folder, then to store the comment is joining .. with reponame
	// and doesn't find out the comment
	generator.repoName = "github.com/0xPolygonHermez/zkevm-node/config/"
	generator.repoNameSuffix = "/config"
	generator.defaultValues.Log.Level = "mylevel"
	schema, err := generator.GenerateJsonSchema(cli)
	require.NoError(t, err)
	require.NotNil(t, schema)
	v, err := getValueFromSchema(schema, []string{"Log", "Level"})
	require.NoError(t, err)
	require.EqualValues(t, "mylevel", v.Default)
	require.NotEmpty(t, v.Description)
}

func getValueFromSchema(schema *jsonschema.Schema, keys []string) (*jsonschema.Schema, error) {
	if schema == nil {
		return nil, errors.New("schema is null")
	}
	subschema := schema

	for _, key := range keys {
		v, exist := subschema.Properties.Get(key)

		if !exist {
			return nil, errors.New("key " + key + " doesnt exist in schema")
		}

		new_schema, ok := v.(*jsonschema.Schema)
		if !ok {
			return nil, errors.New("fails conversion for key " + key + " doesnt exist in schema")
		}
		subschema = new_schema
	}
	return subschema, nil
}
