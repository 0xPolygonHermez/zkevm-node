package tree

// ServerConfig represents the configuration of the MT server.
type ServerConfig struct {
	Host string `mapstructure:"Host"`
	Port int    `mapstructure:"Port"`
}

// ClientConfig represents the configuration of the MT client.
type ClientConfig struct {
	// values for the client
	URI string `mapstructure:"URI"`
}
