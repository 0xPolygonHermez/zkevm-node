package db

import "os"

// Config provide fields to configure the pool
type Config struct {
	// Database name
	Database string

	// User name
	User string

	// Password of the user
	Password string

	// Host address
	Host string

	// Port Number
	Port string
}

// NewConfigFromEnv creates config from standard postgres environment variables,
// see https://www.postgresql.org/docs/11/libpq-envars.html for details
func NewConfigFromEnv() Config {
	return Config{
		Database: getEnv("PGDATABASE", "polygon-hermez"),
		User:     getEnv("PGUSER", "hermez"),
		Password: getEnv("PGPASSWORD", "polygon"),
		Host:     getEnv("PGHOST", "localhost"),
		Port:     getEnv("PGPORT", "5432"),
	}
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
