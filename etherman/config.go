package etherman

import "github.com/0xPolygonHermez/zkevm-node/etherman/etherscan"

// Config represents the configuration of the etherman
type Config struct {
	// URL is the URL of the Ethereum node for L1
	URL string `mapstructure:"URL"`

	//PrivateKeyPath     string `mapstructure:"PrivateKeyPath"`
	//PrivateKeyPassword string `mapstructure:"PrivateKeyPassword"`

	// allow that L1 gas price calculation use multiples sources
	MultiGasProvider bool `mapstructure:"MultiGasProvider"`
	// Configuration for use Etherscan as used as gas provider, basically it needs the API-KEY
	Etherscan etherscan.Config
}
