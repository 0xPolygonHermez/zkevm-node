package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/vectors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

var (
	forcedBatchSignatureHash = crypto.Keccak256Hash([]byte("ForceBatch(uint64,bytes32,address,bytes)"))
)

func TestForcedBatchesVector(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}
	vectorFilesDir := "./../vectors/src/state-transition/forced-tx"
	ctx := context.Background()
	err := filepath.Walk(vectorFilesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !strings.HasSuffix(info.Name(), "list.json") {
			//defer func() {
			//	require.NoError(t, operations.Teardown())
			//}()

			// Load test vectors
			fmt.Println(path)
			testCase, err := vectors.LoadStateTransitionTestCaseV2(path)
			require.NoError(t, err)

			opsCfg := operations.GetDefaultOperationsConfig()
			opsCfg.State.MaxCumulativeGasUsed = 80000000000
			opsman, err := operations.NewManager(ctx, opsCfg)
			require.NoError(t, err)

			// Setting Genesis
			genesisActions := vectors.GenerateGenesisActions(testCase.Genesis)
			require.NoError(t, opsman.SetGenesis(genesisActions))
			require.NoError(t, opsman.Setup())

			// Check initial root
			actualOldStateRoot, err := opsman.State().GetLastStateRoot(ctx, nil)
			require.NoError(t, err)
			require.Equal(t, testCase.ExpectedOldStateRoot, actualOldStateRoot.Hex())
			b, err := hex.DecodeHex(testCase.BatchL2Data)
			require.NoError(t, err)
			txs, txsBytes, err := state.DecodeTxs(b)
			require.NoError(t, err)
			fmt.Println(txs[0].ChainId())

			_, err = sendForcedBatchForVector(t, txsBytes, opsman)
			require.NoError(t, err)

			// Check new root
			actualNewStateRoot, err := opsman.State().GetLastStateRoot(ctx, nil)

			require.NoError(t, err)
			require.Equal(t, testCase.ExpectedNewStateRoot, actualNewStateRoot.Hex())

			return nil
		}
		return nil
	})
	require.NoError(t, err)
}

func sendForcedBatchForVector(t *testing.T, txs []byte, opsman *operations.Manager) (*state.Batch, error) {
	ctx := context.Background()
	st := opsman.State()
	// Connect to ethereum node
	ethClient, err := ethclient.Dial(operations.DefaultL1NetworkURL)
	require.NoError(t, err)

	// Create smc client
	zkEvmAddr := common.HexToAddress(operations.DefaultL1ZkEVMSmartContract)
	zkEvm, err := polygonzkevm.NewPolygonzkevm(zkEvmAddr, ethClient)
	require.NoError(t, err)

	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL1ChainID)
	require.NoError(t, err)

	log.Info("Using address: ", auth.From)

	num, err := zkEvm.LastForceBatch(&bind.CallOpts{Pending: false})
	require.NoError(t, err)

	log.Info("Number of forceBatches in the smc: ", num)

	// Get tip
	tip, err := zkEvm.GetForcedBatchFee(&bind.CallOpts{Pending: false})
	require.NoError(t, err)

	disallowed, err := zkEvm.IsForcedBatchDisallowed(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	if disallowed {
		tx, err := zkEvm.ActivateForceBatches(auth)
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, ethClient, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)
	}

	currentBlock, err := ethClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	log.Debug("currentBlock.Time(): ", currentBlock.Time())

	// Send forceBatch
	tx, err := zkEvm.ForceBatch(auth, txs, tip)
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
		if vLog.Topics[0] != forcedBatchSignatureHash {
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
		require.NoError(t, err)
		require.NotNil(t, forcedBatch)

		err = operations.WaitBatchToBeVirtualized(forcedBatch.BatchNumber, 4*time.Minute, st)
		require.NoError(t, err)
	}

	return forcedBatch, nil
}
