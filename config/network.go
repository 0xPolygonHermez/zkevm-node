package config

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/encoding"
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
	mainnetConfig = NetworkConfig{
		Arity:            4,
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
		Arity:            4,
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
		GenBlockNumber:   6269039,
		PoEAddr:          common.HexToAddress("0xBfAB746dB7fA9ab1F10a4f240F5bB8Ca0924DB56"),
		BridgeAddr:       common.HexToAddress("0x7e102c6AeBA4465089FE3509d3FE85f9F791B0f8"),
		MaticAddr:        common.HexToAddress("0xCaA128d9176CD2afAaA6Af5E739227C1A445c82B"),
		L1ChainID:        5, //Goerli
		L2DefaultChainID: 1000,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x4d5Cf5032B2a844602278b01199ED191A86c93ff"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xA67CD3f603E42dcBF674ffBa511872Bd397EB895"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xbAe5deBDDf9381686ec18a8A2B99E09ADa982adf"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xfcFC415D6D21660b90c0545CA0e91E68172B8650"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x999b52bE90FA59fCaEf59d7243FD874aF3E43E04"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x2536C2745Ac4A584656A830f7bdCd329c94e8F30"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x380ed8Bd696c78395Fb1961BDa42739D2f5242a1"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xd873F6DC68e3057e4B7da74c6b304d0eF0B484C7"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x1EA2EBB132aBD1157831feE038F31A39674b9992"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xb48cA794d49EeC406A5dD2c547717e37b5952a83"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xCF7A13951c6F804E334C39F2eD81D79317e65093"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x56b2118d90cCA76E4683EfECEEC35662372d64Cd"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xd66d09242faa9b3beae711f89d8fff0946974a21"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x615031554479128d65f30Ffa721791D6441d9727"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x890C6F9dAa205C93FeD6546F9ecb4d8D71cfC250"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x8faF6d5f53cD2459aDB7D4cF8682db024dCdCD26"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x8277F27d66BAB5902bde22Fcf0A13452932Ca347"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x35f165ae573dcF9275Cd0923320950cD82D8813E"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xA09E79B2c7F5aB5fdF0C282BfbF646821f4df720"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x66Fd7Bd3FB5CAC3A802D655E2C02A22513ce981f"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x4D3d1e505348cA92bFC2eef44F45C2F89244c5F7"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x38B23d0b9AE34bE5e46C7f1Ce9a9035323A1d027"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xf41832D434405840880a4c4fDB2032a4B243EA35"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xA4e030860e039C83e46F2b62A6136B9DB51f839D"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x5b49b856136360474993Aa35282eabF087848022"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x6224D3aa2B3DCfc70D49a5B5a425CdE5A1812A70"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x0E2B728ddB680a7CE1bFD10464373e66C918399C"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x16654572A1CC90BA6B58626736F9c49f384C3cC6"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x3B0f36ee2dbC45Da5F1370C619c680eA765A25B3"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x0e35d2eAc7E9C05Ce6b8647232e3FF9d5EC47D7a"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x609aC1DFFec23719Cdf1d79fC8c6bA0Dd4D38aD6"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xD6035d50189c13673E789Ab8B14DA92186c4d6b0"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xFF5D77112F7c1E4EF50179129Db6aB2b0A49685F"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x531Df3e16AF889a1d009fF9e49baDc84BB615C19"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x0E7020134410931C9eC16c4dFB251d78E9fC3cAB"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x5A2A939c7D30F24912C97F93EbA321cDe25Dcc26"): bigIntFromBase10String("1000000000000000000000"),
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

func bigIntFromBase10String(s string) *big.Int {
	i, ok := big.NewInt(0).SetString(s, encoding.Base10)
	if !ok {
		return big.NewInt(0)
	}
	return i
}
