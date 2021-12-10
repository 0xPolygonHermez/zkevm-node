package db

// Config provide fields to configure the pool
type Config struct {
	// Database name
	Name string `env:"HERMEZCORE_DB_NAME"`

	// User name
	User string `env:"HERMEZCORE_DB_USER"`

	// Password of the user
	Password string `env:"HERMEZCORE_DB_PASSWORD"`

	// Host address
	Host string `env:"HERMEZCORE_DB_HOST"`

	// Port Number
	Port string `env:"HERMEZCORE_DB_PORT"`
}
