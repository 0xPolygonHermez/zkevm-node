package dataavailability

import "fmt"

// DABackendType is the data availability protocol for the CDK
type DABackendType string

const (
	// DataAvailabilityCommittee is the DAC protocol backend
	DataAvailabilityCommittee DABackendType = "DataAvailabilityCommittee"
)

// Config represents the configuration of the data availability
type Config struct {
	// Backend is the data availability protocol for the CDK
	Backend DABackendType `mapstructure:"Backend"`
}

// Validate validates that the configuration is fine
func (c *Config) Validate() error {
	switch c.Backend {
	case DataAvailabilityCommittee:
		return nil
	default:
		return fmt.Errorf("unsupported backend %s", c.Backend)
	}
}
