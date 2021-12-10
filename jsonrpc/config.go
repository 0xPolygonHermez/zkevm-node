package jsonrpc

// Config represents the configuration of the json rpc
type Config struct {
	Host string `env:"HERMEZCORE_RPC_HOST"`
	Port int    `env:"HERMEZCORE_RPC_PORT"`

	ChainID uint64 `env:"HERMEZCORE_RPC_CHAINID"`
}
