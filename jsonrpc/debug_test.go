package jsonrpc_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	mtDBclientpb "github.com/0xPolygonHermez/zkevm-node/merkletree/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	executorclientpb "github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const (
	l1URL     = "http://localhost:8545"
	l2URL     = "http://localhost:8123"
	txTimeout = 60 * time.Second * 2
)

var (
	testState  *state.State
	stateTree  *merkletree.StateTree
	stateDb    *pgxpool.Pool
	err        error
	stateDBCfg = dbutils.NewStateConfigFromEnv()
	ctx        = context.Background()
	stateCfg   = state.Config{
		MaxCumulativeGasUsed: 800000,
	}
	executorClient                     executorclientpb.ExecutorServiceClient
	mtDBServiceClient                  mtDBclientpb.StateDBServiceClient
	executorClientConn, mtDBClientConn *grpc.ClientConn
)

func TestMain(m *testing.M) {
	// initOrResetDB()

	stateDb, err = db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}
	defer stateDb.Close()

	zkProverURI := testutils.GetEnv("ZKPROVER_URI", "34.245.216.26")

	executorServerConfig := executor.Config{URI: fmt.Sprintf("%s:51071", zkProverURI)}
	var executorCancel context.CancelFunc
	executorClient, executorClientConn, executorCancel = executor.NewExecutorClient(ctx, executorServerConfig)
	s := executorClientConn.GetState()
	log.Infof("executorClientConn state: %s", s.String())
	defer func() {
		executorCancel()
		executorClientConn.Close()
	}()

	mtDBServerConfig := merkletree.Config{URI: fmt.Sprintf("%s:51061", zkProverURI)}
	var mtDBCancel context.CancelFunc
	mtDBServiceClient, mtDBClientConn, mtDBCancel = merkletree.NewMTDBServiceClient(ctx, mtDBServerConfig)
	s = mtDBClientConn.GetState()
	log.Infof("stateDbClientConn state: %s", s.String())
	defer func() {
		mtDBCancel()
		mtDBClientConn.Close()
	}()

	stateTree = merkletree.NewStateTree(mtDBServiceClient)
	testState = state.NewState(stateCfg, state.NewPostgresStorage(stateDb), executorClient, stateTree)

	result := m.Run()
	os.Exit(result)
}

func initOrResetDB() {
	if err := dbutils.InitOrResetState(stateDBCfg); err != nil {
		panic(err)
	}
}

func TestTraceTransaction(t *testing.T) {
	var senderAddress = common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")
	var senderPvtKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	l1Client, err := ethclient.Dial(l1URL)
	require.NoError(t, err)
	l1ChainID, err := l1Client.ChainID(ctx)
	require.NoError(t, err)
	log.Debugf("L1 ChainID = %v", l1ChainID)

	l2Client, err := ethclient.Dial(l2URL)
	require.NoError(t, err)
	l2ChainID, err := l2Client.ChainID(ctx)
	require.NoError(t, err)
	log.Debugf("L2 ChainID = %v", l2ChainID)

	nonceL1, err := l1Client.PendingNonceAt(ctx, senderAddress)
	require.NoError(t, err)
	log.Debugf("L1 nonce = %v", nonceL1)

	nonceL2, err := l2Client.PendingNonceAt(ctx, senderAddress)
	require.NoError(t, err)
	log.Debugf("L2 nonce = %v", nonceL2)

	/*
			// Set Genesis
			block := state.Block{
				BlockNumber: 0,
				BlockHash:   state.ZeroHash,
				ParentHash:  state.ZeroHash,
				ReceivedAt:  time.Now(),
			}

			genesis := state.Genesis{
				Actions: []*state.GenesisAction{
					{
						Address: senderAddress.String(),
						Type:    int(merkletree.LeafTypeBalance),
						Value:   "999977784139164252154",
					},
					{
						Address: senderAddress.String(),
						Type:    int(merkletree.LeafTypeNonce),
						Value:   fmt.Sprintf("%v", nonceL1),
					},
				},
			}

			dbTx, err := testState.BeginStateTransaction(ctx)
			require.NoError(t, err)
			stateRoot, err := testState.SetGenesis(ctx, block, genesis, dbTx)
			require.NoError(t, err)
			require.NoError(t, dbTx.Commit(ctx))

		nonceL2, err := l2Client.PendingNonceAt(ctx, senderAddress)
		require.NoError(t, err)
		log.Debugf("L2 nonce = %v", nonceL2)

		// Create transaction
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    nonceL2,
			To:       &senderAddress,
			Value:    new(big.Int).SetUint64(2),
			Gas:      uint64(30000),
			GasPrice: new(big.Int).SetUint64(1),
			Data:     nil,
		})
	*/
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPvtKey, "0x"))
	require.NoError(t, err)
	authL1, err := bind.NewKeyedTransactorWithChainID(privateKey, l1ChainID)
	require.NoError(t, err)
	authL2, err := bind.NewKeyedTransactorWithChainID(privateKey, l2ChainID)
	require.NoError(t, err)
	/*
		signedTxL2, err := authL2.Signer(authL2.From, tx)
		require.NoError(t, err)

		batchL2Data, err := state.EncodeTransactions([]types.Transaction{*signedTxL2})
		require.NoError(t, err)

		// Create Batch
		processBatchRequest := &executorclientpb.ProcessBatchRequest{
			BatchNum:         1,
			Coinbase:         senderAddress.String(),
			BatchL2Data:      batchL2Data,
			OldStateRoot:     common.HexToHash("0x2f6faa6d4df6548625caca49d4b474e7283173bcedd37480c7c88a221e739399").Bytes(),
			GlobalExitRoot:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
			OldLocalExitRoot: common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
			EthTimestamp:     uint64(0),
			UpdateMerkleTree: 1,
			ChainId:          l2ChainID.Uint64(),
		}

		// Process batch
		_, err = executorClient.ProcessBatch(ctx, processBatchRequest)
		require.NoError(t, err)
	*/

	/*
		// L1
		balance, err := l1Client.BalanceAt(context.Background(), senderAddress, nil)
		require.NoError(t, err)
		require.Equal(t, "999977784139164252154", balance.String())
	*/

	// L2
	balance, err := l2Client.BalanceAt(context.Background(), senderAddress, nil)
	require.NoError(t, err)
	require.NotEqual(t, "0", balance.String())

	// Deploy SC
	scRevertByteCode, err := testutils.ReadBytecode("Revert2/Revert2.bin")
	require.NoError(t, err)
	/*
		nonceL2, err = l2Client.NonceAt(ctx, senderAddress, nil)
		require.NoError(t, err)
		log.Debugf("L2 nonce = %v", nonceL2)

		require.Equal(t, nonceL1, nonceL2)
	*/
	gasPrice, err := l1Client.SuggestGasPrice(ctx)
	require.NoError(t, err)

	// Deploy revert.sol
	txL1 := types.NewTx(&types.LegacyTx{
		Nonce:    nonceL1,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(4000000),
		GasPrice: gasPrice,
		Data:     common.Hex2Bytes(scRevertByteCode),
	})

	signedTxL1, err := authL1.Signer(authL1.From, txL1)
	require.NoError(t, err)

	err = l1Client.SendTransaction(ctx, signedTxL1)
	require.NoError(t, err)

	log.Debugf("txHash L1 = %v", signedTxL1.Hash())

	// Wait fot the transaction to be mined
	err = operations.WaitTxToBeMined(ctx, l1Client, signedTxL1, txTimeout)
	require.NoError(t, err)

	gasPrice, err = l2Client.SuggestGasPrice(ctx)
	require.NoError(t, err)

	// Deploy revert.sol
	txL2 := types.NewTx(&types.LegacyTx{
		Nonce:    nonceL2,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(4000000),
		GasPrice: gasPrice,
		Data:     common.Hex2Bytes(scRevertByteCode),
	})

	signedTxL2, err := authL2.Signer(authL2.From, txL2)
	require.NoError(t, err)

	err = l2Client.SendTransaction(ctx, signedTxL2)
	require.NoError(t, err)

	log.Debugf("txHash L2 = %v", signedTxL2.Hash())

	// Wait fot the transaction to be mined
	err = operations.WaitTxToBeMined(ctx, l2Client, signedTxL2, txTimeout)
	require.NoError(t, err)

	// Debug Transaction
	_, err = debugTransaction(t, l1URL, signedTxL1.Hash().String())
	require.NoError(t, err)

	_, err = debugTransaction(t, l2URL, signedTxL2.Hash().String())
	require.NoError(t, err)

	// Try it directly against state
	// _, err = testState.DebugTransaction(ctx, tx.Hash(), tracer, nil)
	// require.NoError(t, err)
}

func debugTransaction(t *testing.T, nodeURL string, hash string) (string, error) {
	var tracer instrumentation.Tracer

	client := http.Client{
		Timeout: 60 * time.Second,
	}

	tracerFile, err := os.Open("../test/tracers/tracer.json")
	require.NoError(t, err)
	defer tracerFile.Close()

	byteCode, err := ioutil.ReadAll(tracerFile)
	require.NoError(t, err)

	err = json.Unmarshal(byteCode, &tracer)
	require.NoError(t, err)

	// json := `{"jsonrpc": "2.0", "id": 1, "method": "debug_traceTransaction", "params": ["` + hash + `", {"tracer":"` + strings.ReplaceAll(tracer.Code, "\n", " ") + `", "disableStack": true, "disableMemory": true, "disableStorage": true}]}`
	json := `{"jsonrpc": "2.0", "id": 1, "method": "debug_traceTransaction", "params": ["` + hash + `", {"disableStack": true, "disableMemory": true, "disableStorage": true}]}`

	// log.Debugf("Request with:", json)
	jsonByte := []byte(json)
	req, err := http.NewRequest("POST", nodeURL, bytes.NewBuffer(jsonByte))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	require.NoError(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	log.Debugf("Response info: " + resp.Status + " " + string(body))
	return string(body), err
}
