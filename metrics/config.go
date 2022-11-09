package metrics

// Config represents the configuration of the metrics
type Config struct {
	Host    string `mapstructure:"Host"`
	Port    int    `mapstructure:"Port"`
	Enabled bool   `mapstructure:"Enabled"`
}
