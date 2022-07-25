package broadcast

// ServerConfig represents the configuration of the broadcast server.
type ServerConfig struct {
	Host string `mapstructure:"Host"`
	Port int    `mapstructure:"Port"`
}
