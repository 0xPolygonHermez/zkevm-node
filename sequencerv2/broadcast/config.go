package broadcast

// ServerConfig represents the configuration of the broadcast server.
type ServerConfig struct {
	Host string `mapstructure:"Host"`
	Port int    `mapstructure:"Port"`
}

// ClientConfig represents the configuration of the broadcast client.
type ClientConfig struct {
	URI string `mapstructure:"URI"`
}
