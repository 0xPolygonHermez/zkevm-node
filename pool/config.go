package pool

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
