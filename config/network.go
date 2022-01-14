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
	BridgeAddr       common.Address
	MaticAddr        common.Address
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
	balance, _    = new(big.Int).SetString("1000000000000000000000", 10)
	mainnetConfig = NetworkConfig{
		Arity:            16,
		GenBlockNumber:   13808430,
		PoEAddr:          common.HexToAddress("0x11D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"),
		BridgeAddr:       common.HexToAddress("0x11D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"),
		MaticAddr:        common.HexToAddress("0x37AffAf737C3683aB73F6E1B0933b725Ab9796Aa"),
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
		BridgeAddr:       common.HexToAddress("0x21D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"),
		MaticAddr:        common.HexToAddress("0x37AffAf737C3683aB73F6E1B0933b725Ab9796Aa"),
		L1ChainID:        4, //Rinkeby
		L2DefaultChainID: 40000,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"): big.NewInt(1000),
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB"): big.NewInt(2000),
		},
	}
	internalTestnetConfig = NetworkConfig{
		Arity:            4,
		GenBlockNumber:   6195997,
		PoEAddr:          common.HexToAddress("0x516A728856e0F87F99578E4Ff69F4908Db2ac669"),
		BridgeAddr:       common.HexToAddress("0xdEF035490aCb289548F34f0d4Aac488d7314Bf33"),
		MaticAddr:        common.HexToAddress("0x37AffAf737C3683aB73F6E1B0933b725Ab9796Aa"),
		L1ChainID:        5, //Goerli
		L2DefaultChainID: 1000,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"): balance,
			common.HexToAddress("0x4d5Cf5032B2a844602278b01199ED191A86c93ff"): balance,
			common.HexToAddress("0xA67CD3f603E42dcBF674ffBa511872Bd397EB895"): balance,
			common.HexToAddress("0xbAe5deBDDf9381686ec18a8A2B99E09ADa982adf"): balance,
			common.HexToAddress("0xfcFC415D6D21660b90c0545CA0e91E68172B8650"): balance,
			common.HexToAddress("0x999b52bE90FA59fCaEf59d7243FD874aF3E43E04"): balance,
			common.HexToAddress("0x2536C2745Ac4A584656A830f7bdCd329c94e8F30"): balance,
			common.HexToAddress("0x380ed8Bd696c78395Fb1961BDa42739D2f5242a1"): balance,
			common.HexToAddress("0xd873F6DC68e3057e4B7da74c6b304d0eF0B484C7"): balance,
			common.HexToAddress("0x1EA2EBB132aBD1157831feE038F31A39674b9992"): balance,
			common.HexToAddress("0xb48cA794d49EeC406A5dD2c547717e37b5952a83"): balance,
			common.HexToAddress("0xCF7A13951c6F804E334C39F2eD81D79317e65093"): balance,
			common.HexToAddress("0x56b2118d90cCA76E4683EfECEEC35662372d64Cd"): balance,
			common.HexToAddress("0xd66d09242faa9b3beae711f89d8fff0946974a21"): balance,
			common.HexToAddress("0x615031554479128d65f30Ffa721791D6441d9727"): balance,
		},
	}
	localConfig = NetworkConfig{
		Arity:            4,
		GenBlockNumber:   1,
		PoEAddr:          common.HexToAddress("0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"),
		BridgeAddr:       common.HexToAddress("0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"),
		MaticAddr:        common.HexToAddress("0x37AffAf737C3683aB73F6E1B0933b725Ab9796Aa"),
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
