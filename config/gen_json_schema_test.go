package config

import (
	"errors"
	"testing"

	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

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

func TestGenerateJsonSchemaInjectDefaultValue1stLevel(t *testing.T) {
	cli := cli.NewContext(nil, nil, nil)
	generator := NewConfigJsonSchemaGenerater()
	generator.pathSourceCode = "../../"
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
	generator.pathSourceCode = "../../"
	generator.defaultValues.Log.Level = "mylevel"
	schema, err := generator.GenerateJsonSchema(cli)
	require.NoError(t, err)
	require.NotNil(t, schema)
	v, err := getValueFromSchema(schema, []string{"Log", "Level"})
	require.NoError(t, err)
	require.EqualValues(t, "mylevel", v.Default)

}
