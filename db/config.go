package db

// Config provide fields to configure the pool
type Config struct {
	// Database name
	Name string `mapstructure:"Name"`

	// Database User name
	User string `mapstructure:"User"`

	// Database Password of the user
	Password string `mapstructure:"Password"`

	// Host address of database
	Host string `mapstructure:"Host"`

	// Port Number of database
	Port string `mapstructure:"Port"`

	// EnableLog
	EnableLog bool `mapstructure:"EnableLog"`

	// MaxConns is the maximum number of connections in the pool.
	MaxConns int `mapstructure:"MaxConns"`
}
