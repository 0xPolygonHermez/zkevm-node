package db

// Config provide fields to configure the pool
type Config struct {
	// Database name
	Name string `mapstructure:"Name"`

	// User name
	User string `mapstructure:"User"`

	// Password of the user
	Password string `mapstructure:"Password"`

	// Host address
	Host string `mapstructure:"Host"`

	// Port Number
	Port string `mapstructure:"Port"`

	// EnableLog
	EnableLog bool `mapstructure:"EnableLog"`

	// MaxConns is the maximum number of connections in the pool.
	MaxConns int `mapstructure:"MaxConns"`
}
