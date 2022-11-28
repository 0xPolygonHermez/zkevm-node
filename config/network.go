package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

// NetworkConfig is the configuration struct for the different environments
type NetworkConfig struct {
	L2GlobalExitRootManagerAddr common.Address
	L2BridgeAddr                common.Address
	Genesis                     state.Genesis
	MaxCumulativeGasUsed        uint64
}

type genesisFromJSON struct {
	Root    string                   `json:"root"`
	Genesis []genesisAccountFromJSON `json:"genesis"`
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
	config, err := loadGenesisFileConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load genesis configuration from file. Error:", err)
	}
	cfg.NetworkConfig = config
}

func loadGenesisFileConfig(ctx *cli.Context) (NetworkConfig, error) {
	cfgPath := ctx.String(FlagGenesisFile)
	var cfg NetworkConfig

	if cfgPath != "" {
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

		var cfgJSON genesisFromJSON
		err = json.Unmarshal([]byte(b), &cfgJSON)
		if err != nil {
			return NetworkConfig{}, err
		}

		if len(cfgJSON.Genesis) == 0 {
			return cfg, nil
		}

		cfg.Genesis = state.Genesis{
			Root:    common.HexToHash(cfgJSON.Root),
			Actions: []*state.GenesisAction{},
		}

		const l2GlobalExitRootManagerSCName = "GlobalExitRootManagerL2"
		const l2BridgeSCName = "Bridge"

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
	}
	return cfg, nil
}
