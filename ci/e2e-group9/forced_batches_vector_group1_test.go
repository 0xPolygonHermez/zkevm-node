package e2e

import (
	"context"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/constants"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/vectors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	forkID5 uint64 = 5
)

func TestForcedBatchesVectorFiles(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}
	vectorFilesDir := "./../vectors/src/state-transition/forced-tx/group1"
	ctx := context.Background()
	err := filepath.Walk(vectorFilesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !strings.HasSuffix(info.Name(), "list.json") {

			t.Run(info.Name(), func(t *testing.T) {

				defer func() {
					require.NoError(t, operations.Teardown())
				}()

				// Load test vectors
				log.Info("=====================================================================")
				log.Info(path)
				log.Info("=====================================================================")
				testCase, err := vectors.LoadStateTransitionTestCaseV2(path)
				require.NoError(t, err)

				opsCfg := operations.GetDefaultOperationsConfig()
				opsCfg.State.MaxCumulativeGasUsed = 80000000000
				opsman, err := operations.NewManager(ctx, opsCfg)
				require.NoError(t, err)

				// Setting Genesis
				log.Info("###################")
				log.Info("# Setting Genesis #")
				log.Info("###################")
				genesisActions := vectors.GenerateGenesisActions(testCase.Genesis)
				require.NoError(t, opsman.SetGenesis(genesisActions))
				require.NoError(t, opsman.SetForkID(forkID5))
				require.NoError(t, opsman.Setup())

				// Check initial root
				log.Info("################################")
				log.Info("# Verifying initial state root #")
				log.Info("################################")
				actualOldStateRoot, err := opsman.State().GetLastStateRoot(ctx, nil)
				require.NoError(t, err)
				require.Equal(t, testCase.ExpectedOldStateRoot, actualOldStateRoot.Hex())
				decodedData, err := hex.DecodeHex(testCase.BatchL2Data)
				require.NoError(t, err)
				_, txBytes, _, err := state.DecodeTxs(decodedData, forkID5)
				require.NoError(t, err)
				forcedBatch, err := sendForcedBatchForVector(t, txBytes, opsman)
				require.NoError(t, err)
				actualNewStateRoot := forcedBatch.StateRoot
				isClosed, err := opsman.State().IsBatchClosed(ctx, forcedBatch.BatchNumber, nil)
				require.NoError(t, err)

				// wait until is closed
				for !isClosed {
					time.Sleep(1 * time.Second)
					isClosed, err = opsman.State().IsBatchClosed(ctx, forcedBatch.BatchNumber, nil)
					require.NoError(t, err)
				}

				log.Info("#######################")
				log.Info("# Verifying new leafs #")
				log.Info("#######################")
				merkleTree := opsman.State().GetTree()
				for _, expectedNewLeaf := range testCase.ExpectedNewLeafs {
					if expectedNewLeaf.IsSmartContract {
						log.Info("Smart Contract Address: ", expectedNewLeaf.Address)
					} else {
						log.Info("Account Address: ", expectedNewLeaf.Address)
					}
					log.Info("Verifying Balance...")
					actualBalance, err := merkleTree.GetBalance(ctx, common.HexToAddress(expectedNewLeaf.Address), actualNewStateRoot.Bytes())
					require.NoError(t, err)
					require.Equal(t, expectedNewLeaf.Balance.String(), actualBalance.String())

					log.Info("Verifying Nonce...")
					actualNonce, err := merkleTree.GetNonce(ctx, common.HexToAddress(expectedNewLeaf.Address), actualNewStateRoot.Bytes())
					require.NoError(t, err)
					require.Equal(t, expectedNewLeaf.Nonce, actualNonce.String())
					if expectedNewLeaf.IsSmartContract {
						log.Info("Verifying Storage...")
						for positionHex, expectedNewStorageHex := range expectedNewLeaf.Storage {
							position, ok := big.NewInt(0).SetString(positionHex[2:], 16)
							require.True(t, ok)
							expectedNewStorage, ok := big.NewInt(0).SetString(expectedNewStorageHex[2:], 16)
							require.True(t, ok)
							actualStorage, err := merkleTree.GetStorageAt(ctx, common.HexToAddress(expectedNewLeaf.Address), position, actualNewStateRoot.Bytes())
							require.NoError(t, err)
							require.Equal(t, expectedNewStorage, actualStorage)
						}

						log.Info("Verifying HashBytecode...")
						actualHashByteCode, err := merkleTree.GetCodeHash(ctx, common.HexToAddress(expectedNewLeaf.Address), actualNewStateRoot.Bytes())
						require.NoError(t, err)
						require.Equal(t, expectedNewLeaf.HashBytecode, common.BytesToHash(actualHashByteCode).String())
					}
				}
				return
			})

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

	// Send forceBatch
	tx, err := zkEvm.ForceBatch(auth, txs, tip)
	require.NoError(t, err)

	log.Info("Forced Batch Submit to L1 TxHash: ", tx.Hash())
	time.Sleep(1 * time.Second)

	err = operations.WaitTxToBeMined(ctx, ethClient, tx, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)

	currentBlock, err := ethClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	log.Debug("currentBlock.Time(): ", currentBlock.Time())

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
		require.NoError(t, err)
		require.NotNil(t, forcedBatch)

		log.Info("Waiting Forced Batch to be virtualized ...")
		err = operations.WaitBatchToBeVirtualized(forcedBatch.BatchNumber, 4*time.Minute, st)
		require.NoError(t, err)

		log.Info("Waiting Forced Batch to be consolidated ...")
		err = operations.WaitBatchToBeConsolidated(forcedBatch.BatchNumber, 4*time.Minute, st)
		require.NoError(t, err)
	}

	return forcedBatch, nil
}
