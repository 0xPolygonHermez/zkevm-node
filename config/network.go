package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/urfave/cli/v2"
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

type networkConfigFromJSON struct {
	Arity            uint8             `json:"arity"`
	GenBlockNumber   uint64            `json:"genBlockNumber"`
	PoEAddr          string            `json:"poeAddr"`
	BridgeAddr       string            `json:"bridgeAddr"`
	MaticAddr        string            `json:"maticAddr"`
	L1ChainID        uint64            `json:"l1ChainID"`
	L2DefaultChainID uint64            `json:"l2DefaultChainID"`
	Balances         map[string]string `json:"balances"`
}

const (
	testnet         = "testnet"
	internalTestnet = "internaltestnet"
	local           = "local"
	custom          = "custom"
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
		GenBlockNumber:   6279130,
		PoEAddr:          common.HexToAddress("0xaD9d51A5B5237aC36cF9d5f78EA84F8a79d3a274"),
		BridgeAddr:       common.HexToAddress("0x9Fe3268dbD5977e98891528Aa882B7726Ef48118"),
		MaticAddr:        common.HexToAddress("0xA8d4b3CA3e49dCE738E5E29DfcF78499FE7312C9"),
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
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x90F79bf6EB2c4f870365E785982E1f101E93b906"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x14dC79964da2C08b23698B3D3cc7Ca32193d9955"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xa0Ee7A142d267C1f36714E4a8F75612F20a79720"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xBcd4042DE499D14e55001CcbB24a551F3b954096"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x71bE63f3384f5fb98995898A86B02Fb2426c5788"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xFABB0ac9d68B0B445fB7357272Ff202C5651694a"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x1CBd3b2770909D4e10f157cABC84C7264073C9Ec"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xdF3e18d64BC6A983f673Ab319CCaE4f1a57C7097"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xcd3B766CCDd6AE721141F452C550Ca635964ce71"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x2546BcD3c84621e976D8185a91A922aE77ECEc30"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xbDA5747bFD65F08deb54cb465eB87D40e51B197E"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0xdD2FD4581271e230360230F9337D5c0430Bf44C0"): bigIntFromBase10String("1000000000000000000000"),
			common.HexToAddress("0x8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199"): bigIntFromBase10String("1000000000000000000000"),
		},
	}
)

func (cfg *Config) loadNetworkConfig(ctx *cli.Context) {
	network := ctx.String(flagNetwork)
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
	case custom:
		customNetworkConfig, err := loadCustomNetworkConfig(ctx)
		if err != nil {
			log.Fatalf("Failed to load custom network configuration, err:", err)
		}
		cfg.NetworkConfig = customNetworkConfig
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

func loadCustomNetworkConfig(ctx *cli.Context) (NetworkConfig, error) {
	cfgPath := ctx.String(flagNetworkCfg)

	f, err := os.Open(cfgPath) //nolint:gosec
	if err != nil {
		return NetworkConfig{}, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return NetworkConfig{}, err
	}

	var cfgJSON networkConfigFromJSON
	err = json.Unmarshal([]byte(b), &cfgJSON)
	if err != nil {
		return NetworkConfig{}, err
	}

	var cfg NetworkConfig
	cfg.Arity = cfgJSON.Arity
	cfg.GenBlockNumber = cfgJSON.GenBlockNumber
	cfg.PoEAddr = common.HexToAddress(cfgJSON.PoEAddr)
	cfg.BridgeAddr = common.HexToAddress(cfgJSON.BridgeAddr)
	cfg.MaticAddr = common.HexToAddress(cfgJSON.MaticAddr)
	cfg.L1ChainID = cfgJSON.L1ChainID
	cfg.L2DefaultChainID = cfgJSON.L2DefaultChainID
	cfg.Balances = make(map[common.Address]*big.Int, len(cfgJSON.Balances))

	for k, v := range cfgJSON.Balances {
		addr := common.HexToAddress(k)
		balance, ok := big.NewInt(0).SetString(v, encoding.Base10)
		if !ok {
			return NetworkConfig{}, fmt.Errorf("Invalid balance for account %s", addr)
		}
		cfg.Balances[addr] = balance
	}

	return cfg, nil
}
