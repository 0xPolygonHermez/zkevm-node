package config

import (
	"github.com/hermeznetwork/hermez-core/aggregator"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/jsonrpc"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/sequencer"
)

// Config represents the configuration of the entire Hermez Node
type Config struct {
	Log        log.Config
	RPC        jsonrpc.Config
	Sequencer  sequencer.Config
	Aggregator aggregator.Config
}

// Load loads the configuration
func Load() Config {
	// TODO: load from config file
	//nolint:gomnd
	return Config{
		Log: log.Config{
			Level:   "debug",
			Outputs: []string{"stdout"},
		},
		RPC: jsonrpc.Config{
			Host: "",
			Port: 8123,

			ChainID: 2576980377, // 0x99999999,
			Pool: pool.Config{
				Database: "polygon-hermez",
				User:     "hermez",
				Password: "polygon",
				Host:     "localhost",
				Port:     "5432",
			},
		},
		Sequencer: sequencer.Config{
			Etherman: etherman.Config{},
		},
		Aggregator: aggregator.Config{
			Etherman: etherman.Config{},
		},
	}
}
