package config

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/invopop/jsonschema"
	"github.com/urfave/cli/v2"
)

// ConfigJsonSchemaGenerater are the parameters to generate a json-schema based on the T struct
// The parametrization of the function are used for unittest
type ConfigJsonSchemaGenerater[T any] struct {
	repoName       string
	repoNameSuffix string
	// It force to remove any required field in json-schema
	cleanRequiredField bool
	// It read the comments in the code and add as description in schema
	addCodeCommentsToSchema bool
	// Check if there are mapstructure that renames the fields
	checkNoMapStructureIsRenamingFields bool
	// source directories for extract comments
	pathSourceCode string
	// Struct with the default values to set
	defaultValues *T
	// NetworkConfig is read from Genesis file, so make sense to remove
	// from general config file
	removeNetworkConfig bool
}

// NewNodeConfigJsonSchemaGenerater returns a new class for generating json-schema of the node config file (.toml)
func NewNodeConfigJsonSchemaGenerater() ConfigJsonSchemaGenerater[Config] {
	res := ConfigJsonSchemaGenerater[Config]{}
	res.repoName = "github.com/0xPolygonHermez/zkevm-node"
	res.repoNameSuffix = "/config/config"
	res.addCodeCommentsToSchema = true
	res.pathSourceCode = "./"
	res.cleanRequiredField = true
	res.checkNoMapStructureIsRenamingFields = true
	config_default_values, err := Default()
	res.defaultValues = config_default_values
	if err != nil {
		panic("can't create default values for config file")
	}
	return res
}

// NewNetworkConfigJsonSchemaGenerater returns a new class for generating json-schema of the network-custom config file (.json)
func NewNetworkConfigJsonSchemaGenerater() ConfigJsonSchemaGenerater[GenesisFromJSON] {
	res := ConfigJsonSchemaGenerater[GenesisFromJSON]{}
	res.repoName = "github.com/0xPolygonHermez/zkevm-node"
	res.repoNameSuffix = "/config/config"
	res.addCodeCommentsToSchema = true
	res.pathSourceCode = "./"
	res.cleanRequiredField = true
	res.checkNoMapStructureIsRenamingFields = false
	res.defaultValues = nil
	return res
}

// GenerateJsonSchema launchs the generation, and returns the schema
func (s ConfigJsonSchemaGenerater[T]) GenerateJsonSchema(cli *cli.Context) (*jsonschema.Schema, error) {
	if s.checkNoMapStructureIsRenamingFields {
		checkNoMapStructureIsRenamingFields(s.defaultValues)
	}

	r := new(jsonschema.Reflector)
	repoName := s.repoName
	r.Namer = func(rt reflect.Type) string {
		return rt.PkgPath() + "_" + rt.Name()
	}
	r.ExpandedStruct = true
	r.DoNotReference = true
	if s.addCodeCommentsToSchema {
		if err := r.AddGoComments(repoName, s.pathSourceCode); err != nil {
			return nil, err
		}
	}

	schema := r.Reflect(s.defaultValues)
	schema.ID = jsonschema.ID(s.repoName + s.repoNameSuffix)

	if s.cleanRequiredField {
		cleanRequiredFields(schema)
	}

	if s.defaultValues != nil {
		fillDefaultValues(schema, s.defaultValues)
	}

	if s.removeNetworkConfig {
		schema.Properties.Delete("NetworkConfig")
	}

	return schema, nil
}

// SerializeJsonSchema serializes the schema in JSON to be stored
func (s ConfigJsonSchemaGenerater[T]) SerializeJsonSchema(schema *jsonschema.Schema) ([]byte, error) {
	file, err := json.MarshalIndent(schema, "", "\t")
	if err != nil {
		return nil, err
	}
	return file, nil
}

// GenerateJsonSchemaAndWriteToFile generate the schema and store in `output_filename` file
func (s ConfigJsonSchemaGenerater[T]) GenerateJsonSchemaAndWriteToFile(cli *cli.Context, output_filename string) error {
	schema, err := s.GenerateJsonSchema(cli)
	if err != nil {
		return err
	}
	file, err := s.SerializeJsonSchema(schema)
	if err != nil {
		return err
	}

	err = os.WriteFile(output_filename, file, 0600) //nolint:gomnd
	if err != nil {
		return err
	}
	return nil
}

// The tag `magstructure` is not supported by `jsonschema` module
// so, if you try to rename a field using that the documentation is going to incosistent
// For that reason is a good practice to check that is not present this situation in
// the config files
func checkNoMapStructureIsRenamingFields(data interface{}) {
	var reflected reflect.Value
	if reflect.ValueOf(data).Kind() == reflect.Ptr {
		reflected = reflect.ValueOf(data).Elem()
	} else {
		reflected = reflect.ValueOf(data)
	}

	for i := 0; i < reflected.NumField(); i++ {
		field := reflected.Type().Field(i)
		tag := field.Tag.Get("mapstructure")

		if len(tag) > 0 && tag != field.Name {
			panic("field [" + field.Name + "] is renamed using mapstructure to [" + tag + "]! that is not supported")
		}
		if field.Type.Kind() == reflect.Struct {
			checkNoMapStructureIsRenamingFields(reflected.FieldByName(field.Name).Interface())
		}
	}
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

func getFieldNameFromTag(data reflect.Value, key string, tagName string) (reflect.Value, bool) {
	v := data
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)

		if tag == key {
			return v.Field(i), true
		}
	}

	return reflect.Value{}, false
}

func getValueFromStruct(default_config interface{}, key string) reflect.Value {
	var default_object reflect.Value
	if reflect.ValueOf(default_config).Kind() == reflect.Ptr {
		default_object = reflect.ValueOf(default_config).Elem()
	} else {
		default_object = reflect.ValueOf(default_config)
	}
	default_value := default_object.FieldByName(key)
	if !default_value.IsValid() {
		mappedFieldName, found := getFieldNameFromTag(default_object, key, "json")
		if found {
			default_value = mappedFieldName
		}
	}
	return default_value
}

func fillDefaultValuesPartial(schema *jsonschema.Schema, default_config interface{}) {
	if schema.Properties == nil {
		return
	}
	for _, key := range schema.Properties.Keys() {
		value, ok := schema.Properties.Get(key)
		if ok {
			value_schema, _ := value.(*jsonschema.Schema)
			default_value := getValueFromStruct(default_config, key)
			if default_value.IsValid() && variantFieldIsSet(&value_schema.Default) {
				switch value_schema.Type {
				case "array":
					if !default_value.IsZero() && !default_value.IsNil() {
						def_value := default_value.Interface()
						value_schema.Default = def_value
					}
				case "object":
					fillDefaultValuesPartial(value_schema, default_value.Interface())
				default: // string, number, integer, boolean
					if default_value.Type() == reflect.TypeOf(types.Duration{}) {
						duration, ok := default_value.Interface().(types.Duration)
						if ok {
							value_schema.Default = duration.String()
						}
					} else {
						value_schema.Default = default_value.Interface()
					}
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
