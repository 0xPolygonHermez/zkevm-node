package metrics

// Config represents the configuration of the metrics
type Config struct {
	// Host is the address to bind the metrics server
	Host string `mapstructure:"Host"`
	// Port is the port to bind the metrics server
	Port int `mapstructure:"Port"`
	// Enabled is the flag to enable/disable the metrics server
	Enabled bool `mapstructure:"Enabled"`
	// ProfilingHost is the address to bind the profiling server
	ProfilingHost string `mapstructure:"ProfilingHost"`
	// ProfilingPort is the port to bind the profiling server
	ProfilingPort int `mapstructure:"ProfilingPort"`
	// ProfilingEnabled is the flag to enable/disable the profiling server
	ProfilingEnabled bool `mapstructure:"ProfilingEnabled"`
}
