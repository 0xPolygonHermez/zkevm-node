package etherman

// Config represents the configuration of the etherman
type Config struct {
	URL string `mapstructure:"URL"`

	PrivateKeyPath     string `mapstructure:"PrivateKeyPath"`
	PrivateKeyPassword string `mapstructure:"PrivateKeyPassword"`
}
