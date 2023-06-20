package config

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/urfave/cli/v2"
)

type ConfigJsonSchemaGenerater[T any] struct {
	repoName string
	// It force to remove any required field in json-schema
	cleanRequiredField bool
	// It read the comments in the code and add as description in schema
	addCodeCommentsToSchema bool
	// source directories for extract comments
	pathSourceCode string
	// Struct with the default values to set
	defaultValues *T
}

func NewConfigJsonSchemaGenerater() ConfigJsonSchemaGenerater[Config] {
	res := ConfigJsonSchemaGenerater[Config]{}
	res.repoName = "github.com/0xPolygonHermez/zkevm-node"
	res.addCodeCommentsToSchema = true
	res.pathSourceCode = "./"
	res.cleanRequiredField = true
	config_default_values, err := Default()
	res.defaultValues = config_default_values
	if err != nil {
		panic("can't create default values for config file")
	}
	return res
}

func (s ConfigJsonSchemaGenerater[T]) GenerateJsonSchema(cli *cli.Context) (*jsonschema.Schema, error) {
	r := new(jsonschema.Reflector)
	repoName := s.repoName
	r.Namer = func(rt reflect.Type) string {
		return strings.TrimLeft(rt.PkgPath(), repoName) + "_" + rt.Name()
	}
	r.ExpandedStruct = true
	r.DoNotReference = true
	if s.addCodeCommentsToSchema {
		if err := r.AddGoComments(repoName, "./"); err != nil {
			return nil, err
		}
	}
	schema := r.Reflect(s.defaultValues)
	schema.ID = jsonschema.ID(s.repoName + "/config/config")

	if s.cleanRequiredField {
		cleanRequiredFields(schema)
	}

	if s.defaultValues != nil {
		fillDefaultValues(schema, s.defaultValues)
	}

	return schema, nil

}

func (s ConfigJsonSchemaGenerater[T]) SerializeJsonSchema(schema *jsonschema.Schema) ([]byte, error) {
	file, err := json.MarshalIndent(schema, "", "\t")
	if err != nil {
		return nil, err
	}
	return file, nil
}

func variantFieldIsSet(field *interface{}) bool {
	value := reflect.ValueOf(field)
	if value.Kind() == reflect.Ptr && value.IsNil() {
		return false
	} else {
		return true
	}
}

func fillDefaultValues(schema *jsonschema.Schema, default_config interface{}) {
	fillDefaultValuesPartial(schema, default_config)
}

func fillDefaultValuesPartial(schema *jsonschema.Schema, default_config interface{}) {
	if schema.Properties == nil {
		return
	}
	for _, key := range schema.Properties.Keys() {
		value, ok := schema.Properties.Get(key)
		if ok {
			value_schema, _ := value.(*jsonschema.Schema)
			var default_object reflect.Value
			if reflect.ValueOf(default_config).Kind() == reflect.Ptr {
				default_object = reflect.ValueOf(default_config).Elem()
			} else {
				default_object = reflect.ValueOf(default_config)
			}
			default_value := default_object.FieldByName(key)

			if default_value.IsValid() && variantFieldIsSet(&value_schema.Default) {
				switch value_schema.Type {
				case "array":
					//panic("type not supported")
				case "object":
					fillDefaultValuesPartial(value_schema, default_value.Interface())
				default: // string, number, integer, boolean
					value_schema.Default = default_value.Interface()
				}
			}
		}
	}
}

func cleanRequiredFields(schema *jsonschema.Schema) {

	schema.Required = []string{}
	if schema.Properties == nil {
		return
	}
	for _, key := range schema.Properties.Keys() {
		value, ok := schema.Properties.Get(key)
		if ok {
			value_schema, _ := value.(*jsonschema.Schema)
			value_schema.Required = []string{}
			switch value_schema.Type {
			case "object":
				cleanRequiredFields(value_schema)
			case "array":
				cleanRequiredFields(value_schema.Items)
			}
		}
	}
}
