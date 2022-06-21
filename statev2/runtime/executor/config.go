package executor

// ServerConfig represents the configuration of the executor server
type ServerConfig struct {
	Host string `mapstructure:"Host"`
	Port int    `mapstructure:"Port"`
}
