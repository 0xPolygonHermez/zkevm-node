package config

import (
	"flag"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/test/testutils"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestLoadCustomNetworkConfig(t *testing.T) {
	tcs := []struct {
		description      string
		inputConfigStr   string
		expectedConfig   NetworkConfig
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			description: "happy path",
			inputConfigStr: `{
  "deploymentBlockNumber": 6934972,
  "proofOfEfficiencyAddress": "0x2f612dc8fB986E7976AEfc13d8bB0Eb18488a4C9",
  "maticTokenAddress": "0xEa2f9aC0cd926C92923355e88Af73Ee83F2D9C67",

  "arity": 4,
	"globalExitRootStoragePosition": 0,
  "localExitRootStoragePosition":  1,
  "oldStateRootPosition":          0,
	"l1ChainID":                     5,

  "defaultChainID": 1000,
  "genesis": [
    {
      "balance": "0",
      "nonce": "2",
      "address": "0xc949254d682d8c9ad5682521675b8f43b102aec4"
     },
    {
      "balance": "0",
      "nonce": "1",
      "address": "0xae4bb80be56b819606589de61d5ec3b522eeb032",
      "bytecode": "0xbeef1",
      "storage": {
        "0x0000000000000000000000000000000000000000000000000000000000000002": "0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988"
      },
      "contractName": "GlobalExitRootManagerL2"
    },
    {
      "balance": "100000000000000000000000",
      "nonce": "2",
      "address": "0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988",
      "bytecode": "0xbeef2",
      "storage": {
        "0x0000000000000000000000000000000000000000000000000000000000000000": "0xc949254d682d8c9ad5682521675b8f43b102aec4"
      },
      "contractName": "Bridge"
    },
    {
      "balance": "0",
      "nonce": "1",
      "address": "0x61ba0248b0986c2480181c6e76b6adeeaa962483",
      "bytecode": "0xbeef3",
      "storage": {
        "0x0000000000000000000000000000000000000000000000000000000000000000": "0x01"
      }
    }
  ],
  "maxCumulativeGasUsed": 300000
}`,
			expectedConfig: NetworkConfig{
				Arity:          4,
				GenBlockNumber: 6934972,
				PoEAddr:        common.HexToAddress("0x2f612dc8fB986E7976AEfc13d8bB0Eb18488a4C9"),
				MaticAddr:      common.HexToAddress("0xEa2f9aC0cd926C92923355e88Af73Ee83F2D9C67"),

				L2GlobalExitRootManagerAddr:   common.HexToAddress("0xae4bb80be56b819606589de61d5ec3b522eeb032"),
				SystemSCAddr:                  common.Address{},
				GlobalExitRootStoragePosition: 0,
				LocalExitRootStoragePosition:  1,
				OldStateRootPosition:          0,
				L1ChainID:                     5,
				L2DefaultChainID:              1000,
				Genesis: Genesis{
					Balances: map[common.Address]*big.Int{
						common.HexToAddress("0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988"): bigIntFromBase10String("100000000000000000000000"),
					},
					SmartContracts: map[common.Address][]byte{
						common.HexToAddress("0xae4bb80be56b819606589de61d5ec3b522eeb032"): common.FromHex("0xbeef1"),
						common.HexToAddress("0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988"): common.FromHex("0xbeef2"),
						common.HexToAddress("0x61ba0248b0986c2480181c6e76b6adeeaa962483"): common.FromHex("0xbeef3"),
					},
					Storage: map[common.Address]map[*big.Int]*big.Int{
						common.HexToAddress("0xae4bb80be56b819606589de61d5ec3b522eeb032"): {
							new(big.Int).SetBytes(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000002")): new(big.Int).SetBytes(common.Hex2Bytes("9d98deabc42dd696deb9e40b4f1cab7ddbf55988")),
						},
						common.HexToAddress("0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988"): {
							new(big.Int).SetBytes(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000")): new(big.Int).SetBytes(common.Hex2Bytes("c949254d682d8c9ad5682521675b8f43b102aec4")),
						},
						common.HexToAddress("0x61ba0248b0986c2480181c6e76b6adeeaa962483"): {
							new(big.Int).SetBytes(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000")): new(big.Int).SetBytes(common.Hex2Bytes("01")),
						},
					},
					Nonces: map[common.Address]*big.Int{
						common.HexToAddress("0xc949254d682d8c9ad5682521675b8f43b102aec4"): bigIntFromBase10String("2"),
						common.HexToAddress("0xae4bb80be56b819606589de61d5ec3b522eeb032"): bigIntFromBase10String("1"),
						common.HexToAddress("0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988"): bigIntFromBase10String("2"),
						common.HexToAddress("0x61ba0248b0986c2480181c6e76b6adeeaa962483"): bigIntFromBase10String("1"),
					},
				},
				MaxCumulativeGasUsed: 300000,
			},
		},
		{
			description: "imported from network-config.example.json",
			inputConfigStr: `{
  "arity":            4,
  "deploymentBlockNumber":   1,
  "proofOfEfficiencyAddress":          "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9",
  "maticTokenAddress":        "0x37AffAf737C3683aB73F6E1B0933b725Ab9796Aa",
  "systemSCAddr": "0x0000000000000000000000000000000000000000",
  "globalExitRootStoragePosition": 2,
  "localExitRootStoragePosition": 2,
  "oldStateRootPosition": 0,
  "l1ChainID":        1337,
  "defaultChainID": 1000,
  "genesis": [
    {
      "address": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
      "balance": "1000000000000000000000"
    },
    {
      "address": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
      "balance": "2000000000000000000000"
    },
    {
      "address": "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC",
      "balance": "3000000000000000000000"
    }
  ],
  "maxCumulativeGasUsed": 123456
}`,
			expectedConfig: NetworkConfig{
				Arity:          4,
				GenBlockNumber: 1,
				PoEAddr:        common.HexToAddress("0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"),
				MaticAddr:      common.HexToAddress("0x37AffAf737C3683aB73F6E1B0933b725Ab9796Aa"),

				SystemSCAddr:                  common.Address{},
				GlobalExitRootStoragePosition: 2,
				LocalExitRootStoragePosition:  2,
				OldStateRootPosition:          0,
				L1ChainID:                     1337,
				L2DefaultChainID:              1000,
				Genesis: Genesis{
					Balances: map[common.Address]*big.Int{
						common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"): bigIntFromBase10String("1000000000000000000000"),
						common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"): bigIntFromBase10String("2000000000000000000000"),
						common.HexToAddress("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"): bigIntFromBase10String("3000000000000000000000"),
					},
					SmartContracts: map[common.Address][]byte{},
					Storage:        map[common.Address]map[*big.Int]*big.Int{},
					Nonces:         map[common.Address]*big.Int{},
				},
				MaxCumulativeGasUsed: 123456,
			},
		},
		{
			description:      "not valid JSON gives error",
			inputConfigStr:   "not a valid json",
			expectedError:    true,
			expectedErrorMsg: "invalid character",
		},
		{
			description:      "empty JSON gives error",
			expectedError:    true,
			expectedErrorMsg: "unexpected end of JSON input",
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			file, err := ioutil.TempFile("", "loadCustomNetworkConfig")
			require.NoError(t, err)
			defer func() {
				require.NoError(t, os.Remove(file.Name()))
			}()
			require.NoError(t, os.WriteFile(file.Name(), []byte(tc.inputConfigStr), 0600))

			flagSet := flag.NewFlagSet("test", flag.ExitOnError)
			flagSet.String(FlagNetworkCfg, file.Name(), "")
			ctx := cli.NewContext(nil, flagSet, nil)

			actualConfig, err := loadCustomNetworkConfig(ctx)
			require.NoError(t, testutils.CheckError(err, tc.expectedError, tc.expectedErrorMsg))

			require.Equal(t, tc.expectedConfig.Arity, actualConfig.Arity)
			require.Equal(t, tc.expectedConfig.GenBlockNumber, actualConfig.GenBlockNumber)
			require.Equal(t, tc.expectedConfig.PoEAddr, actualConfig.PoEAddr)
			require.Equal(t, tc.expectedConfig.MaticAddr, actualConfig.MaticAddr)
			require.Equal(t, tc.expectedConfig.L2GlobalExitRootManagerAddr, actualConfig.L2GlobalExitRootManagerAddr)
			require.Equal(t, tc.expectedConfig.SystemSCAddr, actualConfig.SystemSCAddr)
			require.Equal(t, tc.expectedConfig.GlobalExitRootStoragePosition, actualConfig.GlobalExitRootStoragePosition)
			require.Equal(t, tc.expectedConfig.LocalExitRootStoragePosition, actualConfig.LocalExitRootStoragePosition)
			require.Equal(t, tc.expectedConfig.OldStateRootPosition, actualConfig.OldStateRootPosition)
			require.Equal(t, tc.expectedConfig.L1ChainID, actualConfig.L1ChainID)
			require.Equal(t, tc.expectedConfig.L2DefaultChainID, actualConfig.L2DefaultChainID)

			require.Equal(t, tc.expectedConfig.Genesis.Balances, actualConfig.Genesis.Balances)
			require.Equal(t, tc.expectedConfig.Genesis.SmartContracts, actualConfig.Genesis.SmartContracts)
			require.Equal(t, tc.expectedConfig.Genesis.Nonces, actualConfig.Genesis.Nonces)
			require.Equal(t, len(tc.expectedConfig.Genesis.Storage), len(tc.expectedConfig.Genesis.Storage))
			for address, expectedStoragePair := range tc.expectedConfig.Genesis.Storage {
				actualStoragePair := actualConfig.Genesis.Storage[address]
				require.NotNil(t, actualStoragePair)
				for expectedStorageKey := range expectedStoragePair {
					expectedStorageKeyStr := expectedStorageKey.String()
					expectedStorageValueStr := tc.expectedConfig.Genesis.Storage[address][expectedStorageKey].String()
					found := false
					for actualStorageKey := range actualStoragePair {
						actualStorageKeyStr := actualStorageKey.String()
						actualStorageValueStr := actualConfig.Genesis.Storage[address][actualStorageKey].String()
						if expectedStorageKeyStr == actualStorageKeyStr && expectedStorageValueStr == actualStorageValueStr {
							found = true
							break
						}
					}
					require.True(t, found)
				}
			}
		})
	}
}

func TestMergeNetworkConfig(t *testing.T) {
	tcs := []struct {
		description          string
		inputCustomConfig    NetworkConfig
		inputBaseConfig      NetworkConfig
		expectedOutputConfig NetworkConfig
	}{
		{
			description:          "empty",
			inputCustomConfig:    NetworkConfig{},
			inputBaseConfig:      NetworkConfig{},
			expectedOutputConfig: NetworkConfig{},
		},
		{
			description: "matching keys",
			inputCustomConfig: NetworkConfig{
				GenBlockNumber: 300,
				PoEAddr:        common.HexToAddress("0xc949254d682d8c9ad5682521675b8f43b102aec4"),
				MaticAddr:      common.HexToAddress("0x1D217d81831009a5fE44C9a1Ee2480e48830CbD4"),
			},
			inputBaseConfig: NetworkConfig{
				GenBlockNumber: 100,
				PoEAddr:        common.HexToAddress("0xb1Fe4a65D3392df68F96daC8eB4df56B2411afBf"),
				MaticAddr:      common.HexToAddress("0x6bad17aC92f0E9313E8c7c3B80E902f1c4D5255F"),
			},
			expectedOutputConfig: NetworkConfig{
				GenBlockNumber: 300,
				PoEAddr:        common.HexToAddress("0xc949254d682d8c9ad5682521675b8f43b102aec4"),
				MaticAddr:      common.HexToAddress("0x1D217d81831009a5fE44C9a1Ee2480e48830CbD4"),
			},
		},
		{
			description: "non-matching keys",
			inputCustomConfig: NetworkConfig{
				GenBlockNumber: 300,
				PoEAddr:        common.HexToAddress("0xc949254d682d8c9ad5682521675b8f43b102aec4"),
				MaticAddr:      common.HexToAddress("0x1D217d81831009a5fE44C9a1Ee2480e48830CbD4"),
			},
			inputBaseConfig: NetworkConfig{
				PoEAddr:              common.HexToAddress("0xb1Fe4a65D3392df68F96daC8eB4df56B2411afBf"),
				Arity:                4,
				L1ChainID:            5,
				MaxCumulativeGasUsed: 300,
			},
			expectedOutputConfig: NetworkConfig{
				Arity:                4,
				L1ChainID:            5,
				MaxCumulativeGasUsed: 300,
				GenBlockNumber:       300,
				PoEAddr:              common.HexToAddress("0xc949254d682d8c9ad5682521675b8f43b102aec4"),
				MaticAddr:            common.HexToAddress("0x1D217d81831009a5fE44C9a1Ee2480e48830CbD4"),
			},
		},
		{
			description: "nested keys",
			inputCustomConfig: NetworkConfig{
				GenBlockNumber: 300,
				Genesis: Genesis{
					Balances: map[common.Address]*big.Int{
						common.HexToAddress("0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988"): bigIntFromBase10String("100000000000000000000000"),
						common.HexToAddress("0x1D217d81831009a5fE44C9a1Ee2480e48830CbD4"): bigIntFromBase10String("900000000000000000000000"),
					},
				},
			},
			inputBaseConfig: NetworkConfig{
				GenBlockNumber: 10,
				Genesis: Genesis{
					Balances: map[common.Address]*big.Int{
						common.HexToAddress("0xb1Fe4a65D3392df68F96daC8eB4df56B2411afBf"): bigIntFromBase10String("200000000000000000000000"),
						common.HexToAddress("0x1D217d81831009a5fE44C9a1Ee2480e48830CbD4"): bigIntFromBase10String("700000000000000000000000"),
					},
				},
			},
			expectedOutputConfig: NetworkConfig{
				GenBlockNumber: 300,
				Genesis: Genesis{
					Balances: map[common.Address]*big.Int{
						common.HexToAddress("0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988"): bigIntFromBase10String("100000000000000000000000"),
						common.HexToAddress("0x1D217d81831009a5fE44C9a1Ee2480e48830CbD4"): bigIntFromBase10String("900000000000000000000000"),
						common.HexToAddress("0xb1Fe4a65D3392df68F96daC8eB4df56B2411afBf"): bigIntFromBase10String("200000000000000000000000"),
					},
				},
			},
		},
		{
			description: "zero address doesn't overwrite destination",
			inputCustomConfig: NetworkConfig{
				PoEAddr: common.Address{},
			},
			inputBaseConfig: NetworkConfig{
				PoEAddr: common.HexToAddress("0xc949254d682d8c9ad5682521675b8f43b102aec4"),
			},
			expectedOutputConfig: NetworkConfig{
				PoEAddr: common.HexToAddress("0xc949254d682d8c9ad5682521675b8f43b102aec4"),
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			actualOutputConfig, err := mergeNetworkConfigs(tc.inputCustomConfig, tc.inputBaseConfig)
			require.NoError(t, err)

			require.Equal(t, tc.expectedOutputConfig, actualOutputConfig)
		})
	}
}
