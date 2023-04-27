package types

import (
	"time"

	"github.com/invopop/jsonschema"
)

// Duration is a wrapper type that parses time duration from text.
type Duration struct {
	time.Duration `validate:"required"`
}

// UnmarshalText unmarshalls time duration from text.
func (d *Duration) UnmarshalText(data []byte) error {
	duration, err := time.ParseDuration(string(data))
	if err != nil {
		return err
	}
	d.Duration = duration
	return nil
}

// NewDuration returns Duration wrapper
func NewDuration(duration time.Duration) Duration {
	return Duration{time.Duration(duration)}
}

// JSONSchema returns a custom schema to be used for the JSON Schema generation of this type
func (Duration) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Title:       "Duration",
		Description: "Duration expressed in units: [ns, us, ms, s, m, h, d]",
		Examples: []interface{}{
			"1m",
			"300ms",
		},
	}
}
