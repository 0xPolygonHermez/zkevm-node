package log

// Config for log
type Config struct {
	// Level of log, e.g. INFO, WARN, ...
	Level string `mapstructure:"Level"`
	// Encoding of the logs ("json" or "console")
	Encoding string `mapstructure:"Encoding"`
	// Outputs
	Outputs []string `mapstructure:"Outputs"`
}
