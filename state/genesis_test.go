package state_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/tools/genesis/genesisparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// genesisAccountReader struct
type genesisAccountReader struct {
	Balance  string            `json:"balance"`
	Nonce    string            `json:"nonce"`
	Address  string            `json:"address"`
	Bytecode string            `json:"bytecode"`
	Storage  map[string]string `json:"storage"`
}

// genesisTestVectorReader struct
type genesisTestVectorReader struct {
	Root     string                 `json:"expectedRoot"`
	Accounts []genesisAccountReader `json:"addresses"`
}

func (gr genesisTestVectorReader) GenesisAccountTest() []genesisparser.GenesisAccountTest {
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

func init() {
	// Change dir to project root
	// This is important because we have relative paths to files containing test vectors
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestGenesisVectors(t *testing.T) {
	// Load test vectors
	var testVectors []genesisTestVectorReader
	files := []string{
		"test/vectors/src/merkle-tree/smt-full-genesis.json",
		"test/vectors/src/merkle-tree/smt-genesis.json",
	}
	for _, f := range files {
		var tv []genesisTestVectorReader
		data, err := os.ReadFile(f)
		require.NoError(t, err)
		err = json.Unmarshal(data, &tv)
		require.NoError(t, err)
		testVectors = append(testVectors, tv...)
	}
	// Run vectors
	for ti, testVector := range testVectors {
		t.Run(fmt.Sprintf("Test vector %d", ti), func(t *testing.T) {
			genesisCase(t, testVector)
		})
	}
}

func genesisCase(t *testing.T, tv genesisTestVectorReader) {
	// Init database instance
	err := dbutils.InitOrResetState(stateDBCfg)
	require.NoError(t, err)
	actions := genesisparser.GenesisTest2Actions(tv.GenesisAccountTest())
	genesis := state.Genesis{
		Actions: actions,
	}
	ctx := context.Background()
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	root, err := testState.SetGenesis(ctx, state.Block{}, genesis, dbTx)
	require.NoError(t, err)
	err = dbTx.Commit(ctx)
	require.NoError(t, err)
	expectedRoot, _ := big.NewInt(0).SetString(tv.Root, 10)
	actualRoot, _ := big.NewInt(0).SetString(hex.EncodeToString(root), 16)
	assert.Equal(t, expectedRoot, actualRoot)
}
