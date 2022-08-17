package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

const repoURL = "https://github.com/0xPolygonHermez/zkevm-commonjs"

// GenesisAction struct
type GenesisAction struct {
	Balance  string            `json:"balanace"`
	Nonce    string            `json:"nonce"`
	Address  string            `json:"address"`
	Bytecode string            `json:"bytecode"`
	Storage  map[string]string `json:"storage"`
}

// GenesisData struct
type GenesisData struct {
	Root    string          `json:"root"`
	Genesis []GenesisAction `json:"genesis"`
}

func main() {
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

	var genesisData GenesisData
	err = json.Unmarshal(genesis, &genesisData)
	if err != nil {
		panic(fmt.Errorf("error json unmarshal: %v", err))
	}

	genesisActions := make([]state.GenesisAction, 0)

	for _, g := range genesisData.Genesis {
		if len(g.Nonce) != 0 {
			genesisActions = append(genesisActions, state.GenesisAction{
				Address: g.Address,
				Type:    int(merkletree.LeafTypeNonce),
				Value:   g.Nonce,
			})
		}
		if len(g.Bytecode) != 0 {
			genesisActions = append(genesisActions, state.GenesisAction{
				Address:  g.Address,
				Type:     int(merkletree.LeafTypeCode),
				Bytecode: g.Bytecode,
			})
		}
		for key, value := range g.Storage {
			genesisActions = append(genesisActions, state.GenesisAction{
				Address:         g.Address,
				Type:            int(merkletree.LeafTypeStorage),
				StoragePosition: key,
				Value:           value,
			})
		}
	}

	gJson, _ := json.MarshalIndent(genesisActions, "", " ")
	gString := string(gJson)
	gString = strings.Replace(gString, "[\n", "", -1)
	gString = strings.Replace(gString, "]", "", -1)
	gString = `//nolint
package config

import "github.com/0xPolygonHermez/zkevm-node/state" 

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

	err = ioutil.WriteFile("./config/genesis.go", []byte(gString), 0600) //nolint:gomnd
	if err != nil {
		panic(fmt.Errorf("error writing file: %v", err))
	}
}
