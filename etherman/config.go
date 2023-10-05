package etherman

import "github.com/0xPolygonHermez/zkevm-node/etherman/etherscan"

// Config represents the configuration of the etherman
type Config struct {
	// IsRollup indicates if the sequence sender is supposed to use a rollup consensus (if false it asumes validium)
	IsRollup bool `mapstructure:"IsRollup"`
	// URL is the URL of the Ethereum node for L1
	URL string `mapstructure:"URL"`

	// ForkIDChunkSize is the max interval for each call to L1 provider to get the forkIDs
	ForkIDChunkSize uint64 `mapstructure:"ForkIDChunkSize"`

	// allow that L1 gas price calculation use multiples sources
	MultiGasProvider bool `mapstructure:"MultiGasProvider"`
	// Configuration for use Etherscan as used as gas provider, basically it needs the API-KEY
	Etherscan etherscan.Config
}
