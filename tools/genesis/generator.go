package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os/exec"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/tools/genesis/genesisparser"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

const (
	repoURL    = "https://github.com/0xPolygonHermez/zkevm-commonjs"
	inputFile  = "tools/fill-genesis/genesis.json"
	outputFile = "../../config/genesis.go"
)

// genesisAccountReader struct
type genesisAccountReader struct {
	Balance  string            `json:"balance"`
	Nonce    string            `json:"nonce"`
	Address  string            `json:"address"`
	Bytecode string            `json:"bytecode"`
	Storage  map[string]string `json:"storage"`
}

// genesisReader struct
type genesisReader struct {
	Root     string                 `json:"root"`
	Accounts []genesisAccountReader `json:"genesis"`
}

func (gr genesisReader) GenesisAccountTest() []genesisparser.GenesisAccountTest {
	accs := []genesisparser.GenesisAccountTest{}
	for i := 0; i < len(gr.Accounts); i++ {
		accs = append(accs, genesisparser.GenesisAccountTest{
			Balance:  gr.Accounts[i].Balance,
			Nonce:    gr.Accounts[i].Nonce,
			Address:  gr.Accounts[i].Address,
			Bytecode: gr.Accounts[i].Bytecode,
			Storage:  gr.Accounts[i].Storage,
		})
	}
	return accs
}

func main() {
	rawGenesis := getLatestGenesisRaw()
	actions := genesisparser.GenesisTest2Actions(rawGenesis.GenesisAccountTest())
	genGoCode(actions)
	err := assertGenesis(rawGenesis.Root)
	if err != nil {
		panic(err)
	}
}

func getLatestGenesisRaw() genesisReader {
	fs := memfs.New()

	_, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		ReferenceName: plumbing.NewBranchReferenceName("main"),
		URL:           repoURL,
	})
	if err != nil {
		panic(fmt.Errorf("error when clone repo: %v", err))
	}

	file, err := fs.Open(inputFile)
	if err != nil {
		panic(fmt.Errorf("error when open file: %v", err))
	}

	scanner := bufio.NewScanner(file)

	genesis := make([]byte, 0)

	for scanner.Scan() {
		genesis = append(genesis, scanner.Bytes()...)
	}
	var genesisData genesisReader
	err = json.Unmarshal(genesis, &genesisData)
	if err != nil {
		panic(fmt.Errorf("error json unmarshal: %v", err))
	}
	return genesisData
}

func genGoCode(actions []*state.GenesisAction) {
	gJson, _ := json.MarshalIndent(actions, "", " ")
	gString := string(gJson)
	gString = strings.Replace(gString, "[\n", "", -1)
	gString = strings.Replace(gString, "]", "", -1)
	gString = `package config

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

	err := ioutil.WriteFile(outputFile, []byte(gString), 0600) //nolint:gomnd
	if err != nil {
		panic(fmt.Errorf("error writing file: %v", err))
	}

	// format code
	cmd := exec.Command("gofmt", "-s", "-w", outputFile)
	res, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Errorf("error formating file: %s.\n%w", string(res), err))
	}
}

func assertGenesis(expectedRoot string) (err error) {
	// Build node
	if err = operations.RunMakeTarget("build-docker"); err != nil {
		log.Error(err)
		return
	}
	// Start DB and executor
	if err = operations.RunMakeTarget("run-db"); err != nil {
		log.Error(err)
		return
	}
	if err = operations.RunMakeTarget("run-zkprover"); err != nil {
		log.Error(err)
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
		log.Error(err)
		return
	}

	// Run node
	err = opsman.Setup()
	if err != nil {
		log.Error(err)
		return
	}

	// Get Genesis root using jRPC
	client, err := ethclient.Dial("http://localhost:8123")
	if err != nil {
		log.Error(err)
		return
	}
	blockHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	if err != nil {
		log.Error(err)
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
