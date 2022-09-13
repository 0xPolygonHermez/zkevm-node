package etherman

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/bridge"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/proofofefficiency"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	})
}

//This function prepare the blockchain, the wallet with funds and deploy the smc
func newTestingEnv() (ethman *Client, ethBackend *backends.SimulatedBackend, maticAddr common.Address, br *bridge.Bridge) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		log.Fatal(err)
	}
	ethman, ethBackend, maticAddr, br, err = NewSimulatedEtherman(Config{}, auth)
	if err != nil {
		log.Fatal(err)
	}
	return ethman, ethBackend, maticAddr, br
}

func TestGEREvent(t *testing.T) {
	// Set up testing environment
	etherman, ethBackend, _, br := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	amount := big.NewInt(1000000000000000)
	a := etherman.auth
	a.Value = amount
	_, err = br.Bridge(a, common.Address{}, 1, etherman.auth.From, amount, []byte{})
	require.NoError(t, err)

	// Mine the tx in a block
	ethBackend.Commit()

	// Now read the event
	finalBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, _, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)

	assert.Equal(t, uint64(2), blocks[0].GlobalExitRoots[0].BlockNumber)
	assert.Equal(t, big.NewInt(1), blocks[0].GlobalExitRoots[0].GlobalExitRootNum)
	assert.NotEqual(t, common.Hash{}, blocks[0].GlobalExitRoots[0].MainnetExitRoot)
	assert.Equal(t, common.Hash{}, blocks[0].GlobalExitRoots[0].RollupExitRoot)
}

func TestForcedBatchEvent(t *testing.T) {
	// Set up testing environment
	etherman, ethBackend, _, _ := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	amount, err := etherman.PoE.CalculateForceProverFee(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	rawTxs := "f84901843b9aca00827b0c945fbdb2315678afecb367f032d93f642f64180aa380a46057361d00000000000000000000000000000000000000000000000000000000000000048203e9808073efe1fa2d3e27f26f32208550ea9b0274d49050b816cadab05a771f4275d0242fd5d92b3fb89575c070e6c930587c520ee65a3aa8cfe382fcad20421bf51d621c"
	data, err := hex.DecodeString(rawTxs)
	require.NoError(t, err)
	_, err = etherman.PoE.ForceBatch(etherman.auth, data, amount)
	require.NoError(t, err)

	// Mine the tx in a block
	ethBackend.Commit()

	// Now read the event
	finalBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, _, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), blocks[0].BlockNumber)
	assert.Equal(t, uint64(2), blocks[0].ForcedBatches[0].BlockNumber)
	assert.NotEqual(t, common.Hash{}, blocks[0].ForcedBatches[0].GlobalExitRoot)
	assert.NotEqual(t, time.Time{}, blocks[0].ForcedBatches[0].ForcedAt)
	assert.Equal(t, uint64(1), blocks[0].ForcedBatches[0].ForcedBatchNumber)
	dataFromSmc := "eaeb077b00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000008cf84901843b9aca00827b0c945fbdb2315678afecb367f032d93f642f64180aa380a46057361d00000000000000000000000000000000000000000000000000000000000000048203e9808073efe1fa2d3e27f26f32208550ea9b0274d49050b816cadab05a771f4275d0242fd5d92b3fb89575c070e6c930587c520ee65a3aa8cfe382fcad20421bf51d621c0000000000000000000000000000000000000000"
	assert.Equal(t, dataFromSmc, hex.EncodeToString(blocks[0].ForcedBatches[0].RawTxsData))
	assert.Equal(t, etherman.auth.From, blocks[0].ForcedBatches[0].Sequencer)
}

func TestSequencedBatchesEvent(t *testing.T) {
	// Set up testing environment
	etherman, ethBackend, _, br := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	// Make a bridge tx
	a := etherman.auth
	a.Value = big.NewInt(1000000000000000)
	_, err = br.Bridge(a, common.Address{}, 1, a.From, a.Value, []byte{})
	require.NoError(t, err)
	ethBackend.Commit()
	a.Value = big.NewInt(0)

	// Get the last ger
	ger, err := etherman.GlobalExitRootManager.GetLastGlobalExitRoot(nil)
	require.NoError(t, err)

	amount, err := etherman.PoE.CalculateForceProverFee(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	rawTxs := "f84901843b9aca00827b0c945fbdb2315678afecb367f032d93f642f64180aa380a46057361d00000000000000000000000000000000000000000000000000000000000000048203e9808073efe1fa2d3e27f26f32208550ea9b0274d49050b816cadab05a771f4275d0242fd5d92b3fb89575c070e6c930587c520ee65a3aa8cfe382fcad20421bf51d621c"
	data, err := hex.DecodeString(rawTxs)
	require.NoError(t, err)
	_, err = etherman.PoE.ForceBatch(etherman.auth, data, amount)
	require.NoError(t, err)
	require.NoError(t, err)
	ethBackend.Commit()

	currentBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	var sequences []proofofefficiency.ProofOfEfficiencyBatchData
	sequences = append(sequences, proofofefficiency.ProofOfEfficiencyBatchData{
		GlobalExitRoot:        ger,
		Timestamp:             currentBlock.Time() - 1,
		ForceBatchesTimestamp: []uint64{30},
		Transactions:          common.Hex2Bytes(rawTxs),
	})
	sequences = append(sequences, proofofefficiency.ProofOfEfficiencyBatchData{
		GlobalExitRoot:        ger,
		Timestamp:             currentBlock.Time() + 1,
		ForceBatchesTimestamp: []uint64{},
		Transactions:          common.Hex2Bytes(rawTxs),
	})
	_, err = etherman.PoE.SequenceBatches(etherman.auth, sequences)
	require.NoError(t, err)

	// Mine the tx in a block
	ethBackend.Commit()

	// Now read the event
	finalBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, order, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	assert.Equal(t, 3, len(blocks))
	assert.Equal(t, 1, len(blocks[2].SequencedBatches))
	assert.Equal(t, common.Hex2Bytes(rawTxs), blocks[2].SequencedBatches[0][1].Transactions)
	assert.Equal(t, currentBlock.Time()-1, blocks[2].SequencedBatches[0][0].Timestamp)
	assert.Equal(t, ger, blocks[2].SequencedBatches[0][0].GlobalExitRoot)
	assert.Equal(t, []uint64{currentBlock.Time()}, blocks[2].SequencedBatches[0][0].ForceBatchesTimestamp)
	assert.Equal(t, 0, order[blocks[2].BlockHash][0].Pos)
}

func TestVerifyBatchEvent(t *testing.T) {
	// Set up testing environment
	etherman, ethBackend, _, _ := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()

	initBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	rawTxs := "f84901843b9aca00827b0c945fbdb2315678afecb367f032d93f642f64180aa380a46057361d00000000000000000000000000000000000000000000000000000000000000048203e9808073efe1fa2d3e27f26f32208550ea9b0274d49050b816cadab05a771f4275d0242fd5d92b3fb89575c070e6c930587c520ee65a3aa8cfe382fcad20421bf51d621c"
	tx := proofofefficiency.ProofOfEfficiencyBatchData{
		GlobalExitRoot:        common.Hash{},
		Timestamp:             initBlock.Time(),
		ForceBatchesTimestamp: []uint64{},
		Transactions:          common.Hex2Bytes(rawTxs),
	}
	_, err = etherman.PoE.SequenceBatches(etherman.auth, []proofofefficiency.ProofOfEfficiencyBatchData{tx})
	require.NoError(t, err)

	// Mine the tx in a block
	ethBackend.Commit()

	var (
		proofA = [2]*big.Int{big.NewInt(1), big.NewInt(1)}
		proofC = [2]*big.Int{big.NewInt(1), big.NewInt(1)}
		proofB = [2][2]*big.Int{proofC, proofC}
	)
	_, err = etherman.PoE.VerifyBatch(etherman.auth, common.Hash{}, common.Hash{}, 1, proofA, proofB, proofC)
	require.NoError(t, err)

	// Mine the tx in a block
	ethBackend.Commit()

	// Now read the event
	finalBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, order, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)

	assert.Equal(t, uint64(3), blocks[1].BlockNumber)
	assert.Equal(t, uint64(1), blocks[1].VerifiedBatches[0].BatchNumber)
	assert.NotEqual(t, common.Address{}, blocks[1].VerifiedBatches[0].Aggregator)
	assert.NotEqual(t, common.Hash{}, blocks[1].VerifiedBatches[0].TxHash)
	assert.Equal(t, GlobalExitRootsOrder, order[blocks[1].BlockHash][0].Name)
	assert.Equal(t, VerifyBatchOrder, order[blocks[1].BlockHash][1].Name)
	assert.Equal(t, 0, order[blocks[1].BlockHash][0].Pos)
	assert.Equal(t, 0, order[blocks[1].BlockHash][1].Pos)
}

func TestSequenceForceBatchesEvent(t *testing.T) {
	// Set up testing environment
	etherman, ethBackend, _, _ := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	amount, err := etherman.PoE.CalculateForceProverFee(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	rawTxs := "f84901843b9aca00827b0c945fbdb2315678afecb367f032d93f642f64180aa380a46057361d00000000000000000000000000000000000000000000000000000000000000048203e9808073efe1fa2d3e27f26f32208550ea9b0274d49050b816cadab05a771f4275d0242fd5d92b3fb89575c070e6c930587c520ee65a3aa8cfe382fcad20421bf51d621c"
	data, err := hex.DecodeString(rawTxs)
	require.NoError(t, err)
	_, err = etherman.PoE.ForceBatch(etherman.auth, data, amount)
	require.NoError(t, err)
	ethBackend.Commit()

	err = ethBackend.AdjustTime((24*7 + 1) * time.Hour)
	require.NoError(t, err)
	ethBackend.Commit()

	_, err = etherman.PoE.SequenceForceBatches(etherman.auth, 1)
	require.NoError(t, err)
	ethBackend.Commit()

	// Now read the event
	finalBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, order, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	assert.Equal(t, uint64(4), blocks[1].BlockNumber)
	assert.Equal(t, uint64(1), blocks[1].SequencedForceBatches[0].LastBatchSequenced)
	assert.Equal(t, uint64(1), blocks[1].SequencedForceBatches[0].ForceBatchNumber)
	assert.Equal(t, 0, order[blocks[1].BlockHash][0].Pos)
}

func TestSendSequences(t *testing.T) {
	// Set up testing environment
	etherman, ethBackend, _, br := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	// Make a bridge tx
	a := etherman.auth
	a.Value = big.NewInt(1000000000000000)
	_, err = br.Bridge(a, common.Address{}, 1, a.From, a.Value, []byte{})
	require.NoError(t, err)
	ethBackend.Commit()
	a.Value = big.NewInt(0)

	// Get the last ger
	ger, err := etherman.GlobalExitRootManager.GetLastGlobalExitRoot(nil)
	require.NoError(t, err)

	currentBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	sequence := ethmanTypes.Sequence{
		GlobalExitRoot: ger,
		Timestamp:      int64(currentBlock.Time() - 1),
		Txs:            []types.Transaction{*tx1},
	}
	tx, err := etherman.sequenceBatches(etherman.auth, []ethmanTypes.Sequence{sequence})
	require.NoError(t, err)
	log.Debug("TX: ", tx.Hash())
	ethBackend.Commit()

	// Now read the event
	finalBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, order, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	assert.Equal(t, 2, len(blocks))
	assert.Equal(t, 1, len(blocks[1].SequencedBatches))
	assert.Equal(t, currentBlock.Time()-1, blocks[1].SequencedBatches[0][0].Timestamp)
	assert.Equal(t, ger, blocks[1].SequencedBatches[0][0].GlobalExitRoot)
	assert.Equal(t, []uint64{}, blocks[1].SequencedBatches[0][0].ForceBatchesTimestamp)
	assert.Equal(t, 0, order[blocks[1].BlockHash][0].Pos)
}
