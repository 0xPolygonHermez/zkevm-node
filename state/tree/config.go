package tree

// Config represents the configuration of the MTService.
type Config struct {
	Host string `mapstructure:"Host"`
	Port int    `mapstructure:"Port"`
}
