package e2e

// import (
// 	"context"
// 	"math/big"
// 	"sync"
// 	"testing"
// 	"time"

// 	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
// 	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevmglobalexitroot"
// 	"github.com/0xPolygonHermez/zkevm-node/log"
// 	"github.com/0xPolygonHermez/zkevm-node/state"
// 	"github.com/0xPolygonHermez/zkevm-node/test/operations"
// 	"github.com/ethereum/go-ethereum"
// 	"github.com/ethereum/go-ethereum/accounts/abi/bind"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/ethclient"
// 	"github.com/stretchr/testify/require"
// )

// func TestForcedBatches(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip()
// 	}

// 	ctx := context.Background()
// 	//defer func() { require.NoError(t, operations.Teardown()) }()

// 	err := operations.Teardown()
// 	require.NoError(t, err)
// 	opsCfg := operations.GetDefaultOperationsConfig()
// 	opsCfg.State.MaxCumulativeGasUsed = 80000000000
// 	opsman, err := operations.NewManager(ctx, opsCfg)
// 	require.NoError(t, err)
// 	err = opsman.Setup()
// 	require.NoError(t, err)
// 	time.Sleep(5 * time.Second)
// 	// Load account with balance on local genesis
// 	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL2ChainID)
// 	require.NoError(t, err)
// 	// Load eth client
// 	client, err := ethclient.Dial(operations.DefaultL2NetworkURL)
// 	require.NoError(t, err)
// 	// Send txs
// 	amount := big.NewInt(10000)
// 	senderBalance, err := client.BalanceAt(ctx, auth.From, nil)
// 	require.NoError(t, err)
// 	senderNonce, err := client.PendingNonceAt(ctx, auth.From)
// 	require.NoError(t, err)

// 	log.Infof("Receiver Addr: %v", toAddress.String())
// 	log.Infof("Sender Addr: %v", auth.From.String())
// 	log.Infof("Sender Balance: %v", senderBalance.String())
// 	log.Infof("Sender Nonce: %v", senderNonce)

// 	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{From: auth.From, To: &toAddress, Value: amount})
// 	require.NoError(t, err)

// 	gasPrice, err := client.SuggestGasPrice(ctx)
// 	require.NoError(t, err)

// 	nonce, err := client.PendingNonceAt(ctx, auth.From)
// 	require.NoError(t, err)

// 	txs := make([]*types.Transaction, 0, nTxs)
// 	for i := 0; i < nTxs; i++ {
// 		tx := types.NewTransaction(nonce+uint64(i), toAddress, amount, gasLimit, gasPrice, nil)
// 		txs = append(txs, tx)
// 	}

// 	wgNormalL2Transfers := new(sync.WaitGroup)
// 	wgNormalL2Transfers.Add(1)
// 	var l2BlockNumbers []*big.Int
// 	go func() {
// 		defer wgNormalL2Transfers.Done()
// 		l2BlockNumbers, err = operations.ApplyL2Txs(ctx, txs, auth, client)
// 		require.NoError(t, err)
// 	}()

// 	time.Sleep(2 * time.Second)
// 	forcedBatch, err := sendForcedBatch(t, opsman)
// 	require.NoError(t, err)
// 	wgNormalL2Transfers.Wait()

// 	// Checking if all txs sent before the forced batch were processed within previous closed batch
// 	for _, l2blockNum := range l2BlockNumbers {
// 		batch, err := opsman.State().GetBatchByL2BlockNumber(ctx, l2blockNum.Uint64(), nil)
// 		require.NoError(t, err)
// 		require.Less(t, batch.BatchNumber, forcedBatch.BatchNumber)
// 	}
// }

// func sendForcedBatch(t *testing.T, opsman *operations.Manager) (*state.Batch, error) {
// 	ctx := context.Background()
// 	st := opsman.State()
// 	// Connect to ethereum node
// 	ethClient, err := ethclient.Dial(operations.DefaultL1NetworkURL)
// 	require.NoError(t, err)

// 	initialGer, _, err := st.GetLatestGer(ctx, gerFinalityBlocks)
// 	require.NoError(t, err)

// 	// Create smc client
// 	zkEvmAddr := common.HexToAddress(operations.DefaultL1ZkEVMSmartContract)
// 	zkEvm, err := polygonzkevm.NewPolygonzkevm(zkEvmAddr, ethClient)
// 	require.NoError(t, err)

// 	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL1ChainID)
// 	require.NoError(t, err)

// 	log.Info("Using address: ", auth.From)

// 	num, err := zkEvm.LastForceBatch(&bind.CallOpts{Pending: false})
// 	require.NoError(t, err)

// 	log.Info("Number of forceBatches in the smc: ", num)

// 	currentBlock, err := ethClient.BlockByNumber(ctx, nil)
// 	require.NoError(t, err)

// 	log.Debug("currentBlock.Time(): ", currentBlock.Time())

// 	// Get tip
// 	tip, err := zkEvm.GetCurrentBatchFee(&bind.CallOpts{Pending: false})
// 	require.NoError(t, err)

// 	managerAddress, err := zkEvm.GlobalExitRootManager(&bind.CallOpts{Pending: false})
// 	require.NoError(t, err)

// 	manager, err := polygonzkevmglobalexitroot.NewPolygonzkevmglobalexitroot(managerAddress, ethClient)
// 	require.NoError(t, err)

// 	rootInContract, err := manager.GetLastGlobalExitRoot(&bind.CallOpts{Pending: false})
// 	require.NoError(t, err)
// 	rootInContractHash := common.BytesToHash(rootInContract[:])

// 	// Send forceBatch
// 	tx, err := zkEvm.ForceBatch(auth, []byte{}, tip)
// 	require.NoError(t, err)

// 	log.Info("TxHash: ", tx.Hash())

// 	time.Sleep(1 * time.Second)

// 	err = operations.WaitTxToBeMined(ctx, ethClient, tx, operations.DefaultTimeoutTxToBeMined)
// 	require.NoError(t, err)

// 	query := ethereum.FilterQuery{
// 		FromBlock: currentBlock.Number(),
// 		Addresses: []common.Address{zkEvmAddr},
// 	}
// 	logs, err := ethClient.FilterLogs(ctx, query)
// 	require.NoError(t, err)

// 	var forcedBatch *state.Batch
// 	for _, vLog := range logs {
// 		fb, err := zkEvm.ParseForceBatch(vLog)
// 		if err != nil {
// 			log.Fatal("failed to parse force batch log event, err: ", err)
// 		}
// 		log.Debugf("log decoded: %+v", fb)
// 		ger := fb.LastGlobalExitRoot
// 		log.Info("GlobalExitRoot: ", ger)
// 		log.Info("Transactions: ", common.Bytes2Hex(fb.Transactions))
// 		fullBlock, err := ethClient.BlockByHash(ctx, vLog.BlockHash)
// 		if err != nil {
// 			log.Errorf("error getting hashParent. BlockNumber: %d. Error: %v", vLog.BlockNumber, err)
// 			return nil, err
// 		}
// 		log.Info("MinForcedTimestamp: ", fullBlock.Time())
// 		forcedBatch, err = st.GetBatchByForcedBatchNum(ctx, fb.ForceBatchNum, nil)
// 		for err == state.ErrStateNotSynchronized {
// 			time.Sleep(1 * time.Second)
// 			forcedBatch, err = st.GetBatchByForcedBatchNum(ctx, fb.ForceBatchNum, nil)
// 		}
// 		require.NoError(t, err)
// 		require.NotNil(t, forcedBatch)

// 		err = operations.WaitBatchToBeVirtualized(forcedBatch.BatchNumber, 4*time.Minute, st)
// 		require.NoError(t, err)

// 		err = operations.WaitBatchToBeConsolidated(forcedBatch.BatchNumber, 4*time.Minute, st)
// 		require.NoError(t, err)

// 		if rootInContractHash != initialGer.GlobalExitRoot {
// 			finalGer, _, err := st.GetLatestGer(ctx, gerFinalityBlocks)
// 			require.NoError(t, err)
// 			if finalGer.GlobalExitRoot != rootInContractHash {
// 				log.Fatal("global exit root is not updated")
// 			}
// 		}
// 	}

// 	return forcedBatch, nil
// }
