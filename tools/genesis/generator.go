package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

const repoURL = "https://github.com/0xPolygonHermez/zkevm-commonjs"

// GenesisAccount struct
type GenesisAccount struct {
	Balance  string            `json:"balance"`
	Nonce    string            `json:"nonce"`
	Address  string            `json:"address"`
	Bytecode string            `json:"bytecode"`
	Storage  map[string]string `json:"storage"`
}

// GenesisReader struct
type GenesisReader struct {
	Root     string           `json:"root"`
	Accounts []GenesisAccount `json:"genesis"`
}

// Genesis struct
type Genesis struct {
	Root  string
	Leafs []state.GenesisAction
}

func main() {
	rawGenesis := getLatestGenesisRaw()
	genesis := raw2Struct(rawGenesis)
	genGoCode(genesis)
	err := assertGenesis(genesis.Root)
	if err != nil {
		panic(err)
	}
}

func getLatestGenesisRaw() []byte {
	fs := memfs.New()

	_, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		panic(fmt.Errorf("error when clone repo: %v", err))
	}

	file, err := fs.Open("tools/fill-genesis/genesis.json")
	if err != nil {
		panic(fmt.Errorf("error when open file: %v", err))
	}

	scanner := bufio.NewScanner(file)

	genesis := make([]byte, 0)

	for scanner.Scan() {
		genesis = append(genesis, scanner.Bytes()...)
	}
	return genesis
}

func raw2Struct(raw []byte) Genesis {
	var genesisData GenesisReader
	err := json.Unmarshal(raw, &genesisData)
	if err != nil {
		panic(fmt.Errorf("error json unmarshal: %v", err))
	}

	leafs := make([]state.GenesisAction, 0)

	for _, acc := range genesisData.Accounts {
		if len(acc.Balance) != 0 && acc.Balance != "0" {
			leafs = append(leafs, state.GenesisAction{
				Address: acc.Address,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   acc.Balance,
			})
		}
		if len(acc.Nonce) != 0 && acc.Nonce != "0" {
			leafs = append(leafs, state.GenesisAction{
				Address: acc.Address,
				Type:    int(merkletree.LeafTypeNonce),
				Value:   acc.Nonce,
			})
		}
		if len(acc.Bytecode) != 0 {
			leafs = append(leafs, state.GenesisAction{
				Address:  acc.Address,
				Type:     int(merkletree.LeafTypeCode),
				Bytecode: acc.Bytecode,
			})
		}
		for key, value := range acc.Storage {
			leafs = append(leafs, state.GenesisAction{
				Address:         acc.Address,
				Type:            int(merkletree.LeafTypeStorage),
				StoragePosition: key,
				Value:           value,
			})
		}
	}
	return Genesis{
		Root:  genesisData.Root,
		Leafs: leafs,
	}
}

func genGoCode(genesis Genesis) {
	gJson, _ := json.MarshalIndent(genesis.Leafs, "", " ")
	gString := string(gJson)
	gString = strings.Replace(gString, "[\n", "", -1)
	gString = strings.Replace(gString, "]", "", -1)
	gString = `//nolint
package config

import (
	"github.com/0xPolygonHermez/zkevm-node/merkletree" 
	"github.com/0xPolygonHermez/zkevm-node/state" 
)

var commonGenesisActions = []*state.GenesisAction{
` + gString + `
}`

	gString = strings.Replace(gString, `"address"`, "Address", -1)
	gString = strings.Replace(gString, `"type"`, "Type", -1)
	gString = strings.Replace(gString, `"storagePosition"`, "StoragePosition", -1)
	gString = strings.Replace(gString, `"value"`, "Value", -1)
	gString = strings.Replace(gString, `"bytecode"`, "Bytecode", -1)
	gString = strings.Replace(gString, `"key"`, "Key", -1)
	gString = strings.Replace(gString, `"root"`, "Root", -1)
	gString = strings.Replace(gString, "\"\n", "\",\n", -1)
	gString = strings.Replace(gString, "}\n", "},\n", -1)
	gString = strings.Replace(gString, "Type: 0,", "Type: int(merkletree.LeafTypeBalance),", -1)
	gString = strings.Replace(gString, "Type: 1,", "Type: int(merkletree.LeafTypeNonce),", -1)
	gString = strings.Replace(gString, "Type: 2,", "Type: int(merkletree.LeafTypeCode),", -1)
	gString = strings.Replace(gString, "Type: 3,", "Type: int(merkletree.LeafTypeStorage),", -1)

	err := ioutil.WriteFile("../../config/genesis.go", []byte(gString), 0600) //nolint:gomnd
	if err != nil {
		panic(fmt.Errorf("error writing file: %v", err))
	}
}

func assertGenesis(expectedRoot string) (err error) {
	// Build node
	if err = operations.RunMakeTarget("build-docker"); err != nil {
		return
	}
	// Start DB and executor
	if err = operations.RunMakeTarget("run-db"); err != nil {
		return
	}
	if err = operations.RunMakeTarget("run-zkprover"); err != nil {
		return
	}
	// Stop everything once done
	defer func() {
		if defErr := operations.Teardown(); defErr != nil {
			err = fmt.Errorf("Error tearing down components: %s", defErr.Error())
		}
	}()

	// Setup opsman
	opsCfg := operations.GetDefaultOperationsConfig()
	opsCfg.State.MaxCumulativeGasUsed = 80000000000
	opsman, err := operations.NewManager(context.Background(), opsCfg)
	if err != nil {
		return
	}

	// Run node
	err = opsman.Setup()
	if err != nil {
		return
	}

	// Get Genesis root using jRPC
	client, err := ethclient.Dial("http://localhost:8123")
	if err != nil {
		return
	}
	blockHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	if err != nil {
		return
	}
	actualRoot := "0x" + blockHeader.Root.Hex()
	if actualRoot != expectedRoot {
		err = fmt.Errorf("Root missmatch: expected: %s, actual %s", expectedRoot, actualRoot)
		return
	}
	fmt.Printf("SUCCESS: expected: %s, actual %s\n", expectedRoot, actualRoot)
	return
}
