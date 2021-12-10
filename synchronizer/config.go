package synchronizer

import "github.com/hermeznetwork/hermez-core/etherman"

// Config represents the configuration of the synchronizer
type Config struct {
	// Etherman is the configuration required by etherman to interact with L1
	Etherman     etherman.Config
	GenesisBlock uint64 `env:"HERMEZCORE_SYNC_GENESISBLOCK"`
}
