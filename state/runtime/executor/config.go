package executor

// Config represents the configuration of the executor server
type Config struct {
	URI                string `mapstructure:"URI"`
	MaxGRPCMessageSize int    `mapstructure:"MaxGRPCMessageSize"`
}
