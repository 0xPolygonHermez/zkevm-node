package log

// Config for log
type Config struct {
	// Level of log, e.g. INFO, WARN, ...
	Level string `mapstructure:"Level"`
	// Outputs
	Outputs []string `mapstructure:"Outputs"`
}
