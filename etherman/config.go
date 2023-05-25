package etherman

import (
	"github.com/0xPolygonHermez/zkevm-node/etherman/etherscan"
)

// Config represents the configuration of the etherman
type Config struct {
	URL string `mapstructure:"URL"`

	PrivateKeyPath     string `mapstructure:"PrivateKeyPath"`
	PrivateKeyPassword string `mapstructure:"PrivateKeyPassword"`

	MultiGasProvider bool `mapstructure:"MultiGasProvider"`
	Etherscan        etherscan.Config
}
