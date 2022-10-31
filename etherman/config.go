package etherman

import (
	"github.com/0xPolygonHermez/zkevm-node/etherman/etherscan"
	"github.com/ethereum/go-ethereum/common"
)

// Config represents the configuration of the etherman
type Config struct {
	URL       string `mapstructure:"URL"`
	L1ChainID uint64 `mapstructure:"L1ChainID"`

	PoEAddr                   common.Address `mapstructure:"PoEAddr"`
	MaticAddr                 common.Address `mapstructure:"MaticAddr"`
	GlobalExitRootManagerAddr common.Address `mapstructure:"GlobalExitRootManagerAddr"`

	PrivateKeyPath     string `mapstructure:"PrivateKeyPath"`
	PrivateKeyPassword string `mapstructure:"PrivateKeyPassword"`

	MultiGasProvider bool `mapstructure:"MultiGasProvider"`
	Etherscan        etherscan.Config
}
