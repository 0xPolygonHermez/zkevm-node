package log

// Config for log
type Config struct {
	// Level of log, e.g. INFO, WARN, ...
	Level string
	// Outputs
	Outputs []string
}
