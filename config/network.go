package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/imdario/mergo"
	"github.com/urfave/cli/v2"
)

//NetworkConfig is the configuration struct for the different environments
type NetworkConfig struct {
	GenBlockNumber                uint64
	PoEAddr                       common.Address
	MaticAddr                     common.Address
	L2GlobalExitRootManagerAddr   common.Address
	GlobalExitRootManagerAddr     common.Address
	SystemSCAddr                  common.Address
	GlobalExitRootStoragePosition uint64
	LocalExitRootStoragePosition  uint64
	OldStateRootPosition          uint64
	L1ChainID                     uint64
	L2ChainID                     uint64
	Genesis                       state.Genesis
	MaxCumulativeGasUsed          uint64
}

type networkConfigFromJSON struct {
	PoEAddr                       string                   `json:"proofOfEfficiencyAddress"`
	MaticAddr                     string                   `json:"maticTokenAddress"`
	GlobalExitRootManagerAddr     string                   `json:"globalExitRootManagerAddress"`
	GenBlockNumber                uint64                   `json:"deploymentBlockNumber"`
	SystemSCAddr                  string                   `json:"systemSCAddr"`
	GlobalExitRootStoragePosition uint64                   `json:"globalExitRootStoragePosition"`
	LocalExitRootStoragePosition  uint64                   `json:"localExitRootStoragePosition"`
	OldStateRootPosition          uint64                   `json:"oldStateRootPosition"`
	L1ChainID                     uint64                   `json:"l1ChainID"`
	L2ChainID                     uint64                   `json:"l2ChainID"`
	Root                          string                   `json:"root"`
	Genesis                       []genesisAccountFromJSON `json:"genesis"`
}

type genesisAccountFromJSON struct {
	Balance      string            `json:"balance"`
	Nonce        string            `json:"nonce"`
	Address      string            `json:"address"`
	Bytecode     string            `json:"bytecode"`
	Storage      map[string]string `json:"storage"`
	ContractName string            `json:"contractName"`
}

const (
	testnet         = "testnet"
	internalTestnet = "internaltestnet"
	local           = "local"
	merge           = "merge"
	custom          = "custom"
)

//nolint:gomnd
var (
	mainnetConfig = NetworkConfig{
		GenBlockNumber:                13808430,
		PoEAddr:                       common.HexToAddress("0x11D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"),
		MaticAddr:                     common.HexToAddress("0x37AffAf737C3683aB73F6E1B0933b725Ab9796Aa"),
		L2GlobalExitRootManagerAddr:   common.HexToAddress("0x0000000000000000000000000000000000000000"),
		GlobalExitRootManagerAddr:     common.HexToAddress("0x0000000000000000000000000000000000000000"),
		SystemSCAddr:                  common.HexToAddress("0x0000000000000000000000000000000000000000"),
		GlobalExitRootStoragePosition: 0,
		LocalExitRootStoragePosition:  1,
		OldStateRootPosition:          0,
		L1ChainID:                     1, //Mainnet
		L2ChainID:                     1000,
		Genesis: state.Genesis{
			Actions: []*state.GenesisAction{
				{
					Address: "0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA",
					Type:    int(merkletree.LeafTypeBalance),
					Value:   "1000",
				},
				{
					Address: "0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB",
					Type:    int(merkletree.LeafTypeBalance),
					Value:   "2000",
				},
			},
		},
	}
	testnetConfig = NetworkConfig{
		GenBlockNumber:                9817974,
		PoEAddr:                       common.HexToAddress("0x21D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"),
		MaticAddr:                     common.HexToAddress("0x37AffAf737C3683aB73F6E1B0933b725Ab9796Aa"),
		L2GlobalExitRootManagerAddr:   common.HexToAddress("0x0000000000000000000000000000000000000000"),
		GlobalExitRootManagerAddr:     common.HexToAddress("0x0000000000000000000000000000000000000000"),
		SystemSCAddr:                  common.HexToAddress("0x0000000000000000000000000000000000000000"),
		GlobalExitRootStoragePosition: 0,
		LocalExitRootStoragePosition:  1,
		OldStateRootPosition:          0,
		L1ChainID:                     4, //Rinkeby
		L2ChainID:                     1000,
		Genesis: state.Genesis{
			Actions: []*state.GenesisAction{
				{
					Address: "0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA",
					Type:    int(merkletree.LeafTypeBalance),
					Value:   "1000",
				},
				{
					Address: "0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB",
					Type:    int(merkletree.LeafTypeBalance),
					Value:   "2000",
				},
			},
		},
	}

	internalTestnetConfig = NetworkConfig{
		GenBlockNumber:                7674348,
		PoEAddr:                       common.HexToAddress("0x159113e5560c9CC2d8c4e716228CCf92c72E9603"),
		MaticAddr:                     common.HexToAddress("0x94Ca2BbE1b469f25D3B22BDf17Fc80ad09E7F662"),
		L2GlobalExitRootManagerAddr:   common.HexToAddress("0xae4bb80be56b819606589de61d5ec3b522eeb032"),
		GlobalExitRootManagerAddr:     common.HexToAddress("0xA379Dd55Eb12e8FCdb467A814A15DE2b29677066"),
		SystemSCAddr:                  common.HexToAddress("0x0000000000000000000000000000000000000000"),
		GlobalExitRootStoragePosition: 0,
		LocalExitRootStoragePosition:  1,
		OldStateRootPosition:          0,
		L1ChainID:                     5, //Goerli
		L2ChainID:                     1000,
		Genesis: state.Genesis{
			Root:    common.HexToHash("0xb33635210b9f5d07769cf70bf5a3cbf241ecbaf79a9b66ef79b28d920da1f776"),
			Actions: []*state.GenesisAction{},
		},
	}

	localConfig = NetworkConfig{
		GenBlockNumber:                1,
		PoEAddr:                       common.HexToAddress("0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6"),
		MaticAddr:                     common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3"),
		L2GlobalExitRootManagerAddr:   common.HexToAddress("0xae4bb80be56b819606589de61d5ec3b522eeb032"),
		GlobalExitRootManagerAddr:     common.HexToAddress("0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"),
		SystemSCAddr:                  common.HexToAddress("0x0000000000000000000000000000000000000000"),
		GlobalExitRootStoragePosition: 0,
		LocalExitRootStoragePosition:  1,
		OldStateRootPosition:          0,
		L1ChainID:                     1337,
		L2ChainID:                     1000,
		Genesis: state.Genesis{
			Root: common.HexToHash("0x5e3d5372166e22ee23b4800aecb491de96f425aa5c7d56f35c96905cc5e12cb8"),
			Actions: []*state.GenesisAction{
				{
					Address: "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
					Type:    int(merkletree.LeafTypeBalance),
					Value:   "100000000000000000000000",
				},
				{
					Address: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
					Type:    int(merkletree.LeafTypeBalance),
					Value:   "100000000000000000000000",
				},
			},
		},
	}

	networkConfigByName = map[string]NetworkConfig{
		testnet:         testnetConfig,
		internalTestnet: internalTestnetConfig,
		local:           localConfig,
	}
)

func (cfg *Config) loadNetworkConfig(ctx *cli.Context) {
	network := ctx.String(FlagNetwork)

	switch network {
	case testnet:
		log.Debug("Testnet network selected")
		cfg.NetworkConfig = testnetConfig
	case internalTestnet:
		log.Debug("InternalTestnet network selected")
		internalTestnetConfig.Genesis.Actions = append(internalTestnetConfig.Genesis.Actions, commonGenesisActions...)
		cfg.NetworkConfig = internalTestnetConfig
	case local:
		log.Debug("Local network selected")
		localConfig.Genesis.Actions = append(localConfig.Genesis.Actions, commonGenesisActions...)
		cfg.NetworkConfig = localConfig
	case custom:
		customNetworkConfig, err := loadCustomNetworkConfig(ctx)
		if err != nil {
			log.Fatalf("Failed to load custom network configuration, err:", err)
		}
		cfg.NetworkConfig = customNetworkConfig
	case merge:
		customNetworkConfig, err := loadCustomNetworkConfig(ctx)
		if err != nil {
			log.Fatalf("Failed to load custom network configuration, err:", err)
		}
		baseNetworkConfigName := ctx.String(FlagNetworkBase)
		baseNetworkConfig, ok := networkConfigByName[baseNetworkConfigName]
		if !ok {
			log.Fatalf("Base network configuration %q not found:", baseNetworkConfigName)
		}
		mergedNetworkConfig, err := mergeNetworkConfigs(customNetworkConfig, baseNetworkConfig)
		if err != nil {
			log.Fatalf("Failed to merge network configurations network configuration, err:", err)
		}
		cfg.NetworkConfig = mergedNetworkConfig
	default:
		log.Debug("Mainnet network selected")
		cfg.NetworkConfig = mainnetConfig
	}
}

func loadCustomNetworkConfig(ctx *cli.Context) (NetworkConfig, error) {
	cfgPath := ctx.String(FlagNetworkCfg)

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
	cfg.GenBlockNumber = cfgJSON.GenBlockNumber
	cfg.PoEAddr = common.HexToAddress(cfgJSON.PoEAddr)
	cfg.MaticAddr = common.HexToAddress(cfgJSON.MaticAddr)
	cfg.GlobalExitRootManagerAddr = common.HexToAddress(cfgJSON.GlobalExitRootManagerAddr)
	cfg.SystemSCAddr = common.HexToAddress(cfgJSON.SystemSCAddr)
	cfg.GlobalExitRootStoragePosition = cfgJSON.GlobalExitRootStoragePosition
	cfg.LocalExitRootStoragePosition = cfgJSON.LocalExitRootStoragePosition
	cfg.OldStateRootPosition = cfgJSON.OldStateRootPosition
	cfg.L1ChainID = cfgJSON.L1ChainID
	cfg.L2ChainID = cfgJSON.L2ChainID

	if len(cfgJSON.Genesis) == 0 {
		return cfg, nil
	}

	cfg.Genesis = state.Genesis{
		Root:    common.HexToHash(cfgJSON.Root),
		Actions: []*state.GenesisAction{},
	}

	const l2GlobalExitRootManagerSCName = "GlobalExitRootManagerL2"

	for _, account := range cfgJSON.Genesis {
		if account.ContractName == l2GlobalExitRootManagerSCName {
			cfg.L2GlobalExitRootManagerAddr = common.HexToAddress(account.Address)
		}
		if account.Balance != "" && account.Balance != "0" {
			action := &state.GenesisAction{
				Address: account.Address,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   account.Balance,
			}
			cfg.Genesis.Actions = append(cfg.Genesis.Actions, action)
		}
		if account.Nonce != "" && account.Nonce != "0" {
			action := &state.GenesisAction{
				Address: account.Address,
				Type:    int(merkletree.LeafTypeNonce),
				Value:   account.Nonce,
			}
			cfg.Genesis.Actions = append(cfg.Genesis.Actions, action)
		}
		if account.Bytecode != "" {
			action := &state.GenesisAction{
				Address:  account.Address,
				Type:     int(merkletree.LeafTypeCode),
				Bytecode: account.Bytecode,
			}
			cfg.Genesis.Actions = append(cfg.Genesis.Actions, action)
		}
		if len(account.Storage) > 0 {
			for storageKey, storageValue := range account.Storage {
				action := &state.GenesisAction{
					Address:         account.Address,
					Type:            int(merkletree.LeafTypeStorage),
					StoragePosition: storageKey,
					Value:           storageValue,
				}
				cfg.Genesis.Actions = append(cfg.Genesis.Actions, action)
			}
		}
	}
	return cfg, nil
}

type addressTransformer struct{}

func (a addressTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ != reflect.TypeOf(common.Address{}) {
		return nil
	}
	return func(dst, src reflect.Value) error {
		if !dst.CanSet() {
			return nil
		}
		hex := src.MethodByName("Hex")
		result := hex.Call([]reflect.Value{})
		if result[0].Interface().(string) != "0x0000000000000000000000000000000000000000" {
			dst.Set(src)
		}
		return nil
	}
}

func mergeNetworkConfigs(custom, base NetworkConfig) (NetworkConfig, error) {
	actionsBack := append(base.Genesis.Actions, custom.Genesis.Actions...)

	if err := mergo.MergeWithOverwrite(&base, custom, mergo.WithTransformers(addressTransformer{})); err != nil {
		return NetworkConfig{}, err
	}

	base.Genesis.Actions = actionsBack

	return base, nil
}
