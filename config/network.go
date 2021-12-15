package config

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/log"
)

//NetworkConfig is the configuration struct for the different environments
type NetworkConfig struct {
	Arity          uint8
	GenBlockNumber uint64
	PoEAddr        common.Address
	L1ChainID      uint64
	L2ChainID      uint64
	Balances       map[common.Address]*big.Int
}

const (
	mainnet         = "mainnet"
	testnet         = "testnet"
	internalTestnet = "internal"
	local           = "local"
)

//nolint:gomnd
var (
	mainnetGen = NetworkConfig{
		Arity:          16,
		GenBlockNumber: 13808430,
		PoEAddr:        common.HexToAddress("0x11D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"),
		L1ChainID:      1, //Mainnet
		L2ChainID:      10000,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"): big.NewInt(1000),
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB"): big.NewInt(2000),
		},
	}
	testnetGen = NetworkConfig{
		Arity:          16,
		GenBlockNumber: 9817974,
		PoEAddr:        common.HexToAddress("0x21D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"),
		L1ChainID:      4, //Rinkeby
		L2ChainID:      40000,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"): big.NewInt(1000),
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB"): big.NewInt(2000),
		},
	}
	internalTestnetGen = NetworkConfig{
		Arity:          16,
		GenBlockNumber: 6025263,
		PoEAddr:        common.HexToAddress("0x31D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"),
		L1ChainID:      5, //Goerli
		L2ChainID:      50000,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"): big.NewInt(1000),
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB"): big.NewInt(2000),
		},
	}
	localGen = NetworkConfig{
		Arity:          4,
		GenBlockNumber: 0,
		PoEAddr:        common.HexToAddress("0x41D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"),
		L1ChainID:      1337,
		L2ChainID:      50000,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"): big.NewInt(1000),
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB"): big.NewInt(2000),
		},
	}
)

func (cfg *Config) loadNetworkConfig(network string) {
	switch network {
	case mainnet:
		log.Debug("Mainnet network selected")
		cfg.NetworkConfig = mainnetGen
	case testnet:
		log.Debug("Testnet network selected")
		cfg.NetworkConfig = testnetGen
	case internalTestnet:
		log.Debug("InternalTestnet network selected")
		cfg.NetworkConfig = internalTestnetGen
	case local:
		log.Debug("Local network selected")
		cfg.NetworkConfig = localGen
	default:
		log.Warn("Unknown network selected")
	}
}
