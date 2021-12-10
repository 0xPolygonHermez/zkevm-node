package config

import (
	"time"

	"github.com/hermeznetwork/hermez-core/aggregator"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/jsonrpc"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/proverclient"
	"github.com/hermeznetwork/hermez-core/sequencer"
	"github.com/hermeznetwork/hermez-core/synchronizer"
)

// Config represents the configuration of the entire Hermez Node
type Config struct {
	Log          log.Config
	Database     db.Config
	Etherman     etherman.Config
	Synchronizer synchronizer.Config
	RPC          jsonrpc.Config
	Sequencer    sequencer.Config
	Aggregator   aggregator.Config
	Prover       proverclient.Config
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
		Database: db.Config{
			Database: "polygon-hermez",
			User:     "hermez",
			Password: "polygon",
			Host:     "localhost",
			Port:     "5432",
		},
		Etherman: etherman.Config{
			PrivateKeyPath:     "../test/test.keystore",
			PrivateKeyPassword: "testonly"},
		RPC: jsonrpc.Config{
			Host: "",
			Port: 8123,
		},
		Synchronizer: synchronizer.Config{},
		Sequencer: sequencer.Config{
			IntervalToProposeBatch: 15 * time.Second,
			URL:                    "http://localhost",
		},
		Aggregator: aggregator.Config{
			IntervalToConsolidateState: 3 * time.Second,
		},
		Prover: proverclient.Config{
			ProverURI: "0.0.0.0:50051",
		},
	}
}
