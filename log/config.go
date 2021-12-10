package log

// Config for log
type Config struct {
	// Level of log, e.g. INFO, WARN, ...
	Level string `env:"HERMEZCORE_LOG_LEVEL"`
	// Outputs
	Outputs []string `env:"HERMEZCORE_LOG_OUTPUTS" envSeparator:","`
}
