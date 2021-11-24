package jsonrpc

import "github.com/hermeznetwork/hermez-core/pool"

// Config represents the configuration of the json rpc
type Config struct {
	Host string
	Port int

	ChainID uint64

	Pool pool.Config
}
