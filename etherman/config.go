package etherman

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/0xPolygonHermez/zkevm-node/etherman/etherscan"
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

	Etherscan etherscan.Config
}
