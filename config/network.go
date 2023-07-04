package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

// NetworkConfig is the configuration struct for the different environments
type NetworkConfig struct {
	L1Config                    etherman.L1Config `json:"l1Config"`
	L2GlobalExitRootManagerAddr common.Address
	L2BridgeAddr                common.Address
	Genesis                     state.Genesis
	MaxCumulativeGasUsed        uint64
}

type network string

const mainnet network = "mainnet"
const testnet network = "testnet"
const custom network = "custom"

// GenesisFromJSON is the config file for network_custom
type GenesisFromJSON struct {
	Root            string                   `json:"root"`
	GenesisBlockNum uint64                   `json:"genesisBlockNumber"`
	Genesis         []genesisAccountFromJSON `json:"genesis"`
	L1Config        etherman.L1Config
}

type genesisAccountFromJSON struct {
	Balance      string            `json:"balance"`
	Nonce        string            `json:"nonce"`
	Address      string            `json:"address"`
	Bytecode     string            `json:"bytecode"`
	Storage      map[string]string `json:"storage"`
	ContractName string            `json:"contractName"`
}

func (cfg *Config) loadNetworkConfig(ctx *cli.Context) {
	var networkJSON string
	switch ctx.String(FlagNetwork) {
	case string(mainnet):
		networkJSON = MainnetNetworkConfigJSON
	case string(testnet):
		networkJSON = TestnetNetworkConfigJSON
	case string(custom):
		var err error
		networkJSON, err = loadGenesisFileAsString(ctx)
		if err != nil {
			panic(err.Error())
		}
	default:
		log.Fatalf("unsupported --network value. Must be one of: [%s, %s, %s]", mainnet, testnet, custom)
	}
	config, err := loadGenesisFromJSONString(networkJSON)
	if err != nil {
		panic(fmt.Errorf("failed to load genesis configuration from file. Error: %v", err))
	}
	cfg.NetworkConfig = config
}

func loadGenesisFileAsString(ctx *cli.Context) (string, error) {
	cfgPath := ctx.String(FlagCustomNetwork)
	if cfgPath != "" {
		f, err := os.Open(cfgPath) //nolint:gosec
		if err != nil {
			return "", err
		}
		defer func() {
			err := f.Close()
			if err != nil {
				log.Error(err)
			}
		}()

		b, err := io.ReadAll(f)
		if err != nil {
			return "", err
		}
		return string(b), nil
	} else {
		return "", errors.New("custom netwrork file not provided. Please use the custom-network-file flag")
	}
}

func loadGenesisFromJSONString(jsonStr string) (NetworkConfig, error) {
	var cfg NetworkConfig

	var cfgJSON GenesisFromJSON
	if err := json.Unmarshal([]byte(jsonStr), &cfgJSON); err != nil {
		return NetworkConfig{}, err
	}

	if len(cfgJSON.Genesis) == 0 {
		return cfg, nil
	}

	cfg.L1Config = cfgJSON.L1Config
	cfg.Genesis = state.Genesis{
		GenesisBlockNum: cfgJSON.GenesisBlockNum,
		Root:            common.HexToHash(cfgJSON.Root),
		GenesisActions:  []*state.GenesisAction{},
	}

	const l2GlobalExitRootManagerSCName = "PolygonZkEVMGlobalExitRootL2 proxy"
	const l2BridgeSCName = "PolygonZkEVMBridge proxy"

	for _, account := range cfgJSON.Genesis {
		if account.ContractName == l2GlobalExitRootManagerSCName {
			cfg.L2GlobalExitRootManagerAddr = common.HexToAddress(account.Address)
		}
		if account.ContractName == l2BridgeSCName {
			cfg.L2BridgeAddr = common.HexToAddress(account.Address)
		}
		if account.Balance != "" && account.Balance != "0" {
			action := &state.GenesisAction{
				Address: account.Address,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   account.Balance,
			}
			cfg.Genesis.GenesisActions = append(cfg.Genesis.GenesisActions, action)
		}
		if account.Nonce != "" && account.Nonce != "0" {
			action := &state.GenesisAction{
				Address: account.Address,
				Type:    int(merkletree.LeafTypeNonce),
				Value:   account.Nonce,
			}
			cfg.Genesis.GenesisActions = append(cfg.Genesis.GenesisActions, action)
		}
		if account.Bytecode != "" {
			action := &state.GenesisAction{
				Address:  account.Address,
				Type:     int(merkletree.LeafTypeCode),
				Bytecode: account.Bytecode,
			}
			cfg.Genesis.GenesisActions = append(cfg.Genesis.GenesisActions, action)
		}
		if len(account.Storage) > 0 {
			for storageKey, storageValue := range account.Storage {
				action := &state.GenesisAction{
					Address:         account.Address,
					Type:            int(merkletree.LeafTypeStorage),
					StoragePosition: storageKey,
					Value:           storageValue,
				}
				cfg.Genesis.GenesisActions = append(cfg.Genesis.GenesisActions, action)
			}
		}
	}

	return cfg, nil
}
