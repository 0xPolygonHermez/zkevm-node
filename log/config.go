package log

// Config for log
type Config struct {
	// Environment defining the log format ("production" or "development").
	// In development mode enables development mode (which makes DPanicLevel logs panic), uses a console encoder, writes to standard error, and disables sampling. Stacktraces are automatically included on logs of WarnLevel and above.
	// Check [here](https://pkg.go.dev/go.uber.org/zap@v1.24.0#NewDevelopmentConfig)
	Environment LogEnvironment `mapstructure:"Environment" jsonschema:"enum=production,enum=development"`
	// Level of log. As lower value more logs are going to be generated
	Level string `mapstructure:"Level" jsonschema:"enum=debug,enum=info,enum=warn,enum=error,enum=dpanic,enum=panic,enum=fatal"`
	// Outputs
	Outputs []string `mapstructure:"Outputs"`
}
