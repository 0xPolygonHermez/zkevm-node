package config

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/log"
)

//NetworkConfig is the configuration struct for the different environments
type NetworkConfig struct {
	Arity            uint8
	GenBlockNumber   uint64
	PoEAddr          common.Address
	L1ChainID        uint64
	L2DefaultChainID uint64
	Balances         map[common.Address]*big.Int
}

const (
	testnet         = "testnet"
	internalTestnet = "internaltestnet"
	local           = "local"
)

//nolint:gomnd
var (
	mainnetConfig = NetworkConfig{
		Arity:            16,
		GenBlockNumber:   13808430,
		PoEAddr:          common.HexToAddress("0x11D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"),
		L1ChainID:        1, //Mainnet
		L2DefaultChainID: 10000,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"): big.NewInt(1000),
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB"): big.NewInt(2000),
		},
	}
	testnetConfig = NetworkConfig{
		Arity:            16,
		GenBlockNumber:   9817974,
		PoEAddr:          common.HexToAddress("0x21D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"),
		L1ChainID:        4, //Rinkeby
		L2DefaultChainID: 40000,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"): big.NewInt(1000),
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB"): big.NewInt(2000),
		},
	}
	internalTestnetConfig = NetworkConfig{
		Arity:            16,
		GenBlockNumber:   6025263,
		PoEAddr:          common.HexToAddress("0x31D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"),
		L1ChainID:        5, //Goerli
		L2DefaultChainID: 50000,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"): big.NewInt(1000),
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB"): big.NewInt(2000),
		},
	}
	localConfig = NetworkConfig{
		Arity:            4,
		GenBlockNumber:   1,
		PoEAddr:          common.HexToAddress("0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"),
		L1ChainID:        1337,
		L2DefaultChainID: 1000,
		Balances:         map[common.Address]*big.Int{},
	}
)

func (cfg *Config) loadNetworkConfig(network string) {
	switch network {
	case testnet:
		log.Debug("Testnet network selected")
		cfg.NetworkConfig = testnetConfig
	case internalTestnet:
		log.Debug("InternalTestnet network selected")
		cfg.NetworkConfig = internalTestnetConfig
	case local:
		log.Debug("Local network selected")
		cfg.NetworkConfig = localConfig
	default:
		log.Debug("Mainnet network selected")
		cfg.NetworkConfig = mainnetConfig
	}
}
