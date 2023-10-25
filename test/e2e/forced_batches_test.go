package e2e

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonrollupmanager"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevmglobalexitroot"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/pol"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/constants"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestForcedBatches(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	defer func() {
		require.NoError(t, operations.Teardown())
	}()

	ctx := context.Background()
	nTxs := 10
	opsman, auth, client, amount, gasLimit, gasPrice, nonce := setupEnvironment(ctx, t)

	// Connect to ethereum node
	ethClient, err := ethclient.Dial(operations.DefaultL1NetworkURL)
	require.NoError(t, err)
	authL1, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL1ChainID)
	require.NoError(t, err)
	zkEvmAddr := common.HexToAddress(operations.DefaultL1ZkEVMSmartContract)
	zkEvm, err := polygonzkevm.NewPolygonzkevm(zkEvmAddr, ethClient)
	require.NoError(t, err)
	pol, err := pol.NewPol(common.HexToAddress(operations.DefaultL1PolSmartContract), ethClient)
	require.NoError(t, err)
	polAmount, _ := big.NewInt(0).SetString("9999999999999999999999", 0)
	_, err = pol.Transfer(authL1, common.HexToAddress(operations.DefaultForcedBatchesAddress), polAmount)
	require.NoError(t, err)
	// Create smc client
	allowed, err := zkEvm.IsForcedBatchAllowed(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	if !allowed {
		log.Debug("Activating ForcedBatches...")
		// Load account with balance on local genesis
		authForcedBatch, err := operations.GetAuth(operations.DefaultForcedBatchesPrivateKey, operations.DefaultL1ChainID)
		require.NoError(t, err)
		_, err = pol.Approve(authForcedBatch, zkEvmAddr, polAmount)
		require.NoError(t, err)
		tx, err := zkEvm.ActivateForceBatches(authL1)
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, ethClient, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)
	}
	txs := make([]*types.Transaction, 0, nTxs)
	for i := 0; i < nTxs; i++ {
		tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)
		nonce = nonce + 1
		txs = append(txs, tx)
	}

	var l2BlockNumbers []*big.Int
	l2BlockNumbers, err = operations.ApplyL2Txs(ctx, txs, auth, client, operations.VerifiedConfirmationLevel)
	require.NoError(t, err)

	time.Sleep(2 * time.Second)
	amount = big.NewInt(0).Add(amount, big.NewInt(10))
	unsignedTx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)
	signedTx, err := auth.Signer(auth.From, unsignedTx)
	require.NoError(t, err)
	encodedTxs, err := state.EncodeTransactions([]types.Transaction{*signedTx}, constants.EffectivePercentage, forkID)
	require.NoError(t, err)
	forcedBatch, err := sendForcedBatch(t, encodedTxs, opsman, ethClient, zkEvm, zkEvmAddr)
	require.NoError(t, err)

	// Checking if all txs sent before the forced batch were processed within previous closed batch
	for _, l2blockNum := range l2BlockNumbers {
		batch, err := opsman.State().GetBatchByL2BlockNumber(ctx, l2blockNum.Uint64(), nil)
		require.NoError(t, err)
		require.Less(t, batch.BatchNumber, forcedBatch.BatchNumber)
	}
}

func setupEnvironment(ctx context.Context, t *testing.T) (*operations.Manager, *bind.TransactOpts, *ethclient.Client, *big.Int, uint64, *big.Int, uint64) {
	err := operations.Teardown()
	require.NoError(t, err)
	opsCfg := operations.GetDefaultOperationsConfig()
	opsCfg.State.MaxCumulativeGasUsed = 80000000000
	genesisFileAsStr, err := config.LoadGenesisFileAsString("../../test/config/test.genesis.config.json")
	require.NoError(t, err)
	genesisConfig, err := config.LoadGenesisFromJSONString(genesisFileAsStr)
	require.NoError(t, err)
	opsman, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	genesisConfig.Genesis.FirstBatchData.Timestamp = uint64(time.Now().Unix())
	require.NoError(t, opsman.SetGenesis(genesisConfig.Genesis.GenesisBlockNum, genesisConfig.Genesis.GenesisActions))
	require.NoError(t, opsman.SetForkID(genesisConfig.Genesis.GenesisBlockNum, forkID))
	dbTx, err := opsman.BeginStateTransaction()
	require.NoError(t, err)
	require.NoError(t, opsman.SetInitialBatch(genesisConfig.Genesis, dbTx))
	require.NoError(t, dbTx.Commit(ctx))
	err = opsman.Setup()
	require.NoError(t, err)
	time.Sleep(5 * time.Second)
	// Load account with balance on local genesis
	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL2ChainID)
	require.NoError(t, err)
	// Load eth client
	client, err := ethclient.Dial(operations.DefaultL2NetworkURL)
	require.NoError(t, err)
	// Send txs
	amount := big.NewInt(10000)
	senderBalance, err := client.BalanceAt(ctx, auth.From, nil)
	require.NoError(t, err)
	senderNonce, err := client.PendingNonceAt(ctx, auth.From)
	require.NoError(t, err)

	log.Infof("Receiver Addr: %v", toAddress.String())
	log.Infof("Sender Addr: %v", auth.From.String())
	log.Infof("Sender Balance: %v", senderBalance.String())
	log.Infof("Sender Nonce: %v", senderNonce)

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{From: auth.From, To: &toAddress, Value: amount})
	require.NoError(t, err)

	gasPrice, err := client.SuggestGasPrice(ctx)
	require.NoError(t, err)

	nonce, err := client.PendingNonceAt(ctx, auth.From)
	require.NoError(t, err)

	return opsman, auth, client, amount, gasLimit, gasPrice, nonce
}

func sendForcedBatch(t *testing.T, txs []byte, opsman *operations.Manager, ethClient *ethclient.Client, zkEvm *polygonzkevm.Polygonzkevm, zkEvmAddr common.Address) (*state.Batch, error) {
	ctx := context.Background()
	st := opsman.State()

	initialGer, _, err := st.GetLatestGer(ctx, gerFinalityBlocks)
	require.NoError(t, err)

	rollupManagerAddr := common.HexToAddress(operations.DefaultL1RollupManagerSmartContract)
	rollupManager, err := polygonrollupmanager.NewPolygonrollupmanager(rollupManagerAddr, ethClient)
	require.NoError(t, err)

	authForcedBatch, err := operations.GetAuth(operations.DefaultForcedBatchesPrivateKey, operations.DefaultL1ChainID)
	require.NoError(t, err)

	log.Info("Using address: ", authForcedBatch.From)

	num, err := zkEvm.LastForceBatch(&bind.CallOpts{Pending: false})
	require.NoError(t, err)

	log.Info("Number of forceBatches in the smc: ", num)

	// Get tip
	tip, err := rollupManager.GetForcedBatchFee(&bind.CallOpts{Pending: false})
	require.NoError(t, err)

	managerAddress, err := zkEvm.GlobalExitRootManager(&bind.CallOpts{Pending: false})
	require.NoError(t, err)

	manager, err := polygonzkevmglobalexitroot.NewPolygonzkevmglobalexitroot(managerAddress, ethClient)
	require.NoError(t, err)

	rootInContract, err := manager.GetLastGlobalExitRoot(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	rootInContractHash := common.BytesToHash(rootInContract[:])

	currentBlock, err := ethClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	log.Debug("currentBlock.Time(): ", currentBlock.Time())

	allowed, err := zkEvm.IsForcedBatchAllowed(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	t.Log("IsForcedBatchAllowed and tip: ", allowed, tip)
	// Send forceBatch
	tx, err := zkEvm.ForceBatch(authForcedBatch, txs, tip)
	require.NoError(t, err)

	log.Info("TxHash: ", tx.Hash())
	time.Sleep(1 * time.Second)

	err = operations.WaitTxToBeMined(ctx, ethClient, tx, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)

	query := ethereum.FilterQuery{
		FromBlock: currentBlock.Number(),
		Addresses: []common.Address{zkEvmAddr},
	}
	logs, err := ethClient.FilterLogs(ctx, query)
	require.NoError(t, err)

	var forcedBatch *state.Batch
	for _, vLog := range logs {
		if vLog.Topics[0] != constants.ForcedBatchSignatureHash {
			logs, err = ethClient.FilterLogs(ctx, query)
			require.NoError(t, err)
			continue
		}
		fb, err := zkEvm.ParseForceBatch(vLog)
		if err != nil {
			log.Errorf("failed to parse force batch log event, err: ", err)
		}
		log.Debugf("log decoded: %+v", fb)
		ger := fb.LastGlobalExitRoot
		log.Info("GlobalExitRoot: ", ger)
		log.Info("Transactions: ", common.Bytes2Hex(fb.Transactions))
		fullBlock, err := ethClient.BlockByHash(ctx, vLog.BlockHash)
		if err != nil {
			log.Errorf("error getting hashParent. BlockNumber: %d. Error: %v", vLog.BlockNumber, err)
			return nil, err
		}
		log.Info("MinForcedTimestamp: ", fullBlock.Time())
		forcedBatch, err = st.GetBatchByForcedBatchNum(ctx, fb.ForceBatchNum, nil)
		for err == state.ErrStateNotSynchronized {
			time.Sleep(1 * time.Second)
			forcedBatch, err = st.GetBatchByForcedBatchNum(ctx, fb.ForceBatchNum, nil)
		}
		log.Info("ForcedBatchNum: ", forcedBatch.BatchNumber)
		require.NoError(t, err)
		require.NotNil(t, forcedBatch)

		log.Info("Waiting for batch to be virtualized...")
		err = operations.WaitBatchToBeVirtualized(forcedBatch.BatchNumber, 4*time.Minute, st)
		require.NoError(t, err)

		log.Info("Waiting for batch to be consolidated...")
		err = operations.WaitBatchToBeConsolidated(forcedBatch.BatchNumber, 4*time.Minute, st)
		require.NoError(t, err)

		if rootInContractHash != initialGer.GlobalExitRoot {
			finalGer, _, err := st.GetLatestGer(ctx, gerFinalityBlocks)
			require.NoError(t, err)
			if finalGer.GlobalExitRoot != rootInContractHash {
				log.Fatal("global exit root is not updated")
			}
		}
	}

	return forcedBatch, nil
}
