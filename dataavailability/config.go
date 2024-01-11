package dataavailability

import "fmt"

// DABackendType is the data availabilty protocol for the CDK
type DABackendType string

const (
	// DataAvailabilityCommittee is the DAC protocol backend
	DataAvailabilityCommittee DABackendType = "DataAvailabilityCommittee"
)

// Config represents the configuration of the data availability
type Config struct {
	// Backend is the data availabilty protocol for the CDK
	Backend DABackendType `mapstructure:"Backend"`
}

func (c *Config) Validate() error {
	switch c.Backend {
	case DataAvailabilityCommittee:
		return nil
	default:
		return fmt.Errorf("unsupported backend %s", c.Backend)
	}
}
