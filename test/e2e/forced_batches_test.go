package e2e

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/matic"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevmglobalexitroot"
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

	var err error
	nTxs := 10
	ctx := context.Background()
	opsman, auth, clientL2, amount, gasLimit, gasPrice, nonce := setupEnvironment(ctx, t)
	l1 := setupEnvironmentL1(ctx, t)

	txs := make([]*types.Transaction, 0, nTxs)
	for i := 0; i < nTxs; i++ {
		tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)
		nonce = nonce + 1
		txs = append(txs, tx)
	}

	var l2BlockNumbers []*big.Int
	l2BlockNumbers, err = operations.ApplyL2Txs(ctx, txs, auth, clientL2, operations.VerifiedConfirmationLevel)
	require.NoError(t, err)

	time.Sleep(2 * time.Second)
	amount = big.NewInt(0).Add(amount, big.NewInt(10))
	unsignedTx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)
	signedTx, err := auth.Signer(auth.From, unsignedTx)
	require.NoError(t, err)
	encodedTxs, err := state.EncodeTransactions([]types.Transaction{*signedTx}, constants.EffectivePercentage, forkID6)
	require.NoError(t, err)
	forcedBatch, err := sendForcedBatch(ctx, t, encodedTxs, opsman, l1)
	require.NoError(t, err)

	// Checking if all txs sent before the forced batch were processed within previous closed batch
	for _, l2blockNum := range l2BlockNumbers {
		batch, err := opsman.State().GetBatchByL2BlockNumber(ctx, l2blockNum.Uint64(), nil)
		require.NoError(t, err)
		require.Less(t, batch.BatchNumber, forcedBatch.BatchNumber)
	}
}

type L1Stuff struct {
	ethClient       *ethclient.Client
	authSequencer   *bind.TransactOpts
	authForcedBatch *bind.TransactOpts
	zkEvmAddr       common.Address
	zkEvm           *polygonzkevm.Polygonzkevm
}

func setupEnvironmentL1(ctx context.Context, t *testing.T) *L1Stuff {
	// Connect to ethereum node
	ethClient, err := ethclient.Dial(operations.DefaultL1NetworkURL)
	require.NoError(t, err)
	authSequencer, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL1ChainID)
	require.NoError(t, err)
	authForcedBatch, err := operations.GetAuth(operations.DefaultForcedBatchesPrivateKey, operations.DefaultL1ChainID)
	require.NoError(t, err)
	maticSmc, err := matic.NewMatic(common.HexToAddress(operations.DefaultL1MaticSmartContract), ethClient)
	require.NoError(t, err)
	maticAmount, _ := big.NewInt(0).SetString("9999999999999999999999", 0)
	txValue, err := maticSmc.Transfer(authSequencer, common.HexToAddress(operations.DefaultForcedBatchesAddress), maticAmount)
	require.NoError(t, err)
	log.Debug(txValue)
	log.Debugf("Waiting for tx %s to be mined (transfer of matic from sequencer -> forcedBatches)", txValue.Hash().String())
	err = operations.WaitTxToBeMined(ctx, ethClient, txValue, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)
	balance, err := maticSmc.BalanceOf(&bind.CallOpts{Pending: false}, common.HexToAddress(operations.DefaultSequencerAddress))
	require.NoError(t, err)
	log.Debugf("Account (sequencer) %s MATIC balance %s", operations.DefaultSequencerAddress, balance.String())

	balance, err = maticSmc.BalanceOf(&bind.CallOpts{Pending: false}, common.HexToAddress(operations.DefaultForcedBatchesAddress))
	require.NoError(t, err)
	log.Debugf("Account (force_batches) %s MATIC balance %s", operations.DefaultForcedBatchesAddress, balance.String())
	log.Debugf("Approve to zkEVM SMC to spend %s MATIC", maticAmount.String())
	_, err = maticSmc.Approve(authForcedBatch, common.HexToAddress(operations.DefaultL1ZkEVMSmartContract), maticAmount)
	require.NoError(t, err)

	zkEvmAddr := common.HexToAddress(operations.DefaultL1ZkEVMSmartContract)
	zkEvm, err := polygonzkevm.NewPolygonzkevm(zkEvmAddr, ethClient)
	require.NoError(t, err)
	return &L1Stuff{ethClient: ethClient, authSequencer: authSequencer, authForcedBatch: authForcedBatch, zkEvmAddr: zkEvmAddr, zkEvm: zkEvm}
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
	require.NoError(t, opsman.SetForkID(genesisConfig.Genesis.GenesisBlockNum, forkID6))
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

func sendForcedBatch(ctx context.Context, t *testing.T, txs []byte, opsman *operations.Manager, l1 *L1Stuff) (*state.Batch, error) {
	st := opsman.State()

	initialGer, _, err := st.GetLatestGer(ctx, gerFinalityBlocks)
	require.NoError(t, err)

	log.Info("Using address: ", l1.authForcedBatch.From)

	num, err := l1.zkEvm.LastForceBatch(&bind.CallOpts{Pending: false})
	require.NoError(t, err)

	log.Info("Number of forceBatches in the smc: ", num)

	// Get tip
	tip, err := l1.zkEvm.GetForcedBatchFee(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	log.Infof("Foced Batch Fee:%s", tip.String())
	managerAddress, err := l1.zkEvm.GlobalExitRootManager(&bind.CallOpts{Pending: false})
	require.NoError(t, err)

	manager, err := polygonzkevmglobalexitroot.NewPolygonzkevmglobalexitroot(managerAddress, l1.ethClient)
	require.NoError(t, err)

	rootInContract, err := manager.GetLastGlobalExitRoot(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	rootInContractHash := common.BytesToHash(rootInContract[:])

	disallowed, err := l1.zkEvm.IsForcedBatchDisallowed(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	if disallowed {
		tx, err := l1.zkEvm.ActivateForceBatches(l1.authSequencer)
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, l1.ethClient, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)
	}

	currentBlock, err := l1.ethClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	log.Debug("currentBlock.Time(): ", currentBlock.Time())

	// Send forceBatch
	tx, err := l1.zkEvm.ForceBatch(l1.authForcedBatch, txs, tip)
	require.NoError(t, err)

	log.Info("TxHash: ", tx.Hash())
	time.Sleep(1 * time.Second)

	err = operations.WaitTxToBeMined(ctx, l1.ethClient, tx, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)

	query := ethereum.FilterQuery{
		FromBlock: currentBlock.Number(),
		Addresses: []common.Address{l1.zkEvmAddr},
	}
	logs, err := l1.ethClient.FilterLogs(ctx, query)
	require.NoError(t, err)

	var forcedBatch *state.Batch
	for _, vLog := range logs { // TODO check if that make sense
		if vLog.Topics[0] != constants.ForcedBatchSignatureHash {
			logs, err = l1.ethClient.FilterLogs(ctx, query)
			require.NoError(t, err)
			continue
		}
		fb, err := l1.zkEvm.ParseForceBatch(vLog)
		if err != nil {
			log.Errorf("failed to parse force batch log event, err: ", err)
		}
		log.Debugf("log decoded: %+v", fb)
		ger := fb.LastGlobalExitRoot
		log.Info("GlobalExitRoot: ", ger)
		log.Info("Transactions: ", common.Bytes2Hex(fb.Transactions))
		fullBlock, err := l1.ethClient.BlockByHash(ctx, vLog.BlockHash)
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
