package tree

// MTStoreBackendType is an alias for the merkletree store backends.
type MTStoreBackendType string

const (
	// PgMTStoreBackend is a store backend composed of PostgreSQL without cache.
	PgMTStoreBackend MTStoreBackendType = "PostgreSQL"

	// PgRistrettoMTStoreBackend is a store backend composed of PostgreSQL with
	// Ristretto as cache.
	PgRistrettoMTStoreBackend MTStoreBackendType = "PostgreSQLRistretto"

	// BadgerRistrettoMTStoreBackend is a store backend composed of BadgerDB with
	// Ristretto as cache.
	BadgerRistrettoMTStoreBackend MTStoreBackendType = "BadgerDBRistretto"
)

// ServerConfig represents the configuration of the MT server.
type ServerConfig struct {
	Host         string             `mapstructure:"Host"`
	Port         int                `mapstructure:"Port"`
	StoreBackend MTStoreBackendType `mapstructure:"StoreBackend"`
}

// ClientConfig represents the configuration of the MT client.
type ClientConfig struct {
	// values for the client
	URI string `mapstructure:"URI"`
}
