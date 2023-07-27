package sequencer

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

var (
	workerCfg = WorkerCfg{
		ResourceCostMultiplier: 1000,
	}
)

type workerAddTxTestCase struct {
	name   string
	from   common.Address
	txHash common.Hash
	nonce  uint64
	// isClaim                bool
	benefit                int64
	cost                   *big.Int
	counters               state.ZKCounters
	usedBytes              uint64
	expectedEfficiencyList []common.Hash
}

type workerAddrQueueInfo struct {
	from    common.Address
	nonce   *big.Int
	balance *big.Int
}

func processWorkerAddTxTestCases(t *testing.T, worker *Worker, testCases []workerAddTxTestCase) {
	totalWeight := float64(worker.batchResourceWeights.WeightArithmetics +
		worker.batchResourceWeights.WeightBatchBytesSize + worker.batchResourceWeights.WeightBinaries +
		worker.batchResourceWeights.WeightCumulativeGasUsed + worker.batchResourceWeights.WeightKeccakHashes +
		worker.batchResourceWeights.WeightMemAligns + worker.batchResourceWeights.WeightPoseidonHashes +
		worker.batchResourceWeights.WeightPoseidonPaddings + worker.batchResourceWeights.WeightSteps)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tx := TxTracker{}

			tx.WeightMultipliers = calculateWeightMultipliers(worker.batchResourceWeights, totalWeight)
			tx.Constraints = worker.batchConstraints
			tx.ResourceCostMultiplier = worker.cfg.ResourceCostMultiplier
			tx.Hash = testCase.txHash
			tx.HashStr = testCase.txHash.String()
			tx.From = testCase.from
			tx.FromStr = testCase.from.String()
			tx.Nonce = testCase.nonce
			tx.Benefit = new(big.Int).SetInt64(testCase.benefit)
			tx.Cost = testCase.cost
			tx.BatchResources.Bytes = testCase.usedBytes
			tx.updateZKCounters(testCase.counters, worker.batchConstraints, worker.batchResourceWeights)
			t.Logf("%s=%s", testCase.name, fmt.Sprintf("%.2f", tx.Efficiency))

			_, err := worker.AddTxTracker(ctx, &tx)
			if err != nil {
				return
			}

			el := worker.efficiencyList
			if el.len() != len(testCase.expectedEfficiencyList) {
				t.Fatalf("Error efficiencylist.len(%d) != expectedEfficiencyList.len(%d)", el.len(), len(testCase.expectedEfficiencyList))
			}
			for i := 0; i < el.len(); i++ {
				if el.getByIndex(i).HashStr != string(testCase.expectedEfficiencyList[i].String()) {
					t.Fatalf("Error efficiencylist(%d). Expected=%s, Actual=%s", i, testCase.expectedEfficiencyList[i].String(), el.getByIndex(i).HashStr)
				}
			}
		})
	}
}

func TestWorkerAddTx(t *testing.T) {
	var nilErr error

	// Init ZKEVM resourceCostWeight values
	rcWeigth := batchResourceWeights{}
	rcWeigth.WeightCumulativeGasUsed = 1
	rcWeigth.WeightArithmetics = 1
	rcWeigth.WeightBinaries = 1
	rcWeigth.WeightKeccakHashes = 1
	rcWeigth.WeightMemAligns = 1
	rcWeigth.WeightPoseidonHashes = 1
	rcWeigth.WeightPoseidonPaddings = 1
	rcWeigth.WeightSteps = 1
	rcWeigth.WeightBatchBytesSize = 2

	// Init ZKEVM resourceCostMax values
	rcMax := batchConstraints{}
	rcMax.MaxCumulativeGasUsed = 10
	rcMax.MaxArithmetics = 10
	rcMax.MaxBinaries = 10
	rcMax.MaxKeccakHashes = 10
	rcMax.MaxMemAligns = 10
	rcMax.MaxPoseidonHashes = 10
	rcMax.MaxPoseidonPaddings = 10
	rcMax.MaxSteps = 10
	rcMax.MaxBatchBytesSize = 10

	stateMock := NewStateMock(t)
	worker := initWorker(stateMock, rcMax, rcWeigth)

	ctx := context.Background()

	stateMock.On("GetLastStateRoot", ctx, nil).Return(common.Hash{0}, nilErr)

	addrQueueInfo := []workerAddrQueueInfo{
		{from: common.Address{1}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{2}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{3}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{4}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
	}

	for _, aq := range addrQueueInfo {
		stateMock.On("GetNonceByStateRoot", ctx, aq.from, common.Hash{0}).Return(aq.nonce, nilErr)
		stateMock.On("GetBalanceByStateRoot", ctx, aq.from, common.Hash{0}).Return(aq.balance, nilErr)
	}

	addTxsTC := []workerAddTxTestCase{
		{
			name: "Adding from:0x01, tx:0x01/ef:10", from: common.Address{1}, txHash: common.Hash{1}, nonce: 1,
			benefit: 1000, cost: new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 1, UsedKeccakHashes: 1, UsedPoseidonHashes: 1, UsedPoseidonPaddings: 1, UsedMemAligns: 1, UsedArithmetics: 1, UsedBinaries: 1, UsedSteps: 1},
			usedBytes: 1,
			expectedEfficiencyList: []common.Hash{
				{1},
			},
		},
		{
			name: "Adding from:0x02, tx:0x02/ef:20", from: common.Address{2}, txHash: common.Hash{2}, nonce: 1,
			benefit: 2000, cost: new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 1, UsedKeccakHashes: 1, UsedPoseidonHashes: 1, UsedPoseidonPaddings: 1, UsedMemAligns: 1, UsedArithmetics: 1, UsedBinaries: 1, UsedSteps: 1},
			usedBytes: 1,
			expectedEfficiencyList: []common.Hash{
				{2}, {1},
			},
		},
		{
			name: "Readding from:0x02, tx:0x02/ef:4", from: common.Address{2}, txHash: common.Hash{2}, nonce: 1,
			benefit: 2000, cost: new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 5, UsedKeccakHashes: 5, UsedPoseidonHashes: 5, UsedPoseidonPaddings: 5, UsedMemAligns: 5, UsedArithmetics: 5, UsedBinaries: 5, UsedSteps: 5},
			usedBytes: 5,
			expectedEfficiencyList: []common.Hash{
				{1}, {2},
			},
		},
		{
			name: "Readding from:0x03, tx:0x03/ef:25", from: common.Address{3}, txHash: common.Hash{3}, nonce: 1,
			benefit: 5000, cost: new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 2, UsedKeccakHashes: 2, UsedPoseidonHashes: 2, UsedPoseidonPaddings: 2, UsedMemAligns: 2, UsedArithmetics: 2, UsedBinaries: 2, UsedSteps: 2},
			usedBytes: 2,
			expectedEfficiencyList: []common.Hash{
				{3}, {1}, {2},
			},
		},
	}

	processWorkerAddTxTestCases(t, worker, addTxsTC)

	// Change counters fpr tx:0x03/ef:9.61
	counters := state.ZKCounters{CumulativeGasUsed: 6, UsedKeccakHashes: 6, UsedPoseidonHashes: 6, UsedPoseidonPaddings: 6, UsedMemAligns: 6, UsedArithmetics: 6, UsedBinaries: 6, UsedSteps: 6}
	worker.UpdateTx(common.Hash{3}, common.Address{3}, counters)

	addTxsTC = []workerAddTxTestCase{
		{
			name: "Adding from:0x04, tx:0x04/ef:100", from: common.Address{4}, txHash: common.Hash{4}, nonce: 1,
			benefit: 10000, cost: new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 1, UsedKeccakHashes: 1, UsedPoseidonHashes: 1, UsedPoseidonPaddings: 1, UsedMemAligns: 1, UsedArithmetics: 1, UsedBinaries: 1, UsedSteps: 1},
			usedBytes: 1,
			expectedEfficiencyList: []common.Hash{
				{4}, {1}, {3}, {2},
			},
		},
	}

	processWorkerAddTxTestCases(t, worker, addTxsTC)
}

func TestWorkerGetBestTx(t *testing.T) {
	var nilErr error

	// Init ZKEVM resourceCostWeight values
	rcWeight := batchResourceWeights{}
	rcWeight.WeightCumulativeGasUsed = 1
	rcWeight.WeightArithmetics = 1
	rcWeight.WeightBinaries = 1
	rcWeight.WeightKeccakHashes = 1
	rcWeight.WeightMemAligns = 1
	rcWeight.WeightPoseidonHashes = 1
	rcWeight.WeightPoseidonPaddings = 1
	rcWeight.WeightSteps = 1
	rcWeight.WeightBatchBytesSize = 2

	// Init ZKEVM resourceCostMax values
	rcMax := batchConstraints{}
	rcMax.MaxCumulativeGasUsed = 10
	rcMax.MaxArithmetics = 10
	rcMax.MaxBinaries = 10
	rcMax.MaxKeccakHashes = 10
	rcMax.MaxMemAligns = 10
	rcMax.MaxPoseidonHashes = 10
	rcMax.MaxPoseidonPaddings = 10
	rcMax.MaxSteps = 10
	rcMax.MaxBatchBytesSize = 10

	rc := state.BatchResources{
		ZKCounters: state.ZKCounters{CumulativeGasUsed: 10, UsedKeccakHashes: 10, UsedPoseidonHashes: 10, UsedPoseidonPaddings: 10, UsedMemAligns: 10, UsedArithmetics: 10, UsedBinaries: 10, UsedSteps: 10},
		Bytes:      10,
	}

	stateMock := NewStateMock(t)
	worker := initWorker(stateMock, rcMax, rcWeight)

	ctx := context.Background()

	stateMock.On("GetLastStateRoot", ctx, nil).Return(common.Hash{0}, nilErr)

	addrQueueInfo := []workerAddrQueueInfo{
		{from: common.Address{1}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{2}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{3}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{4}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
	}

	for _, aq := range addrQueueInfo {
		stateMock.On("GetNonceByStateRoot", ctx, aq.from, common.Hash{0}).Return(aq.nonce, nilErr)
		stateMock.On("GetBalanceByStateRoot", ctx, aq.from, common.Hash{0}).Return(aq.balance, nilErr)
	}

	addTxsTC := []workerAddTxTestCase{
		{
			name: "Adding from:0x01, tx:0x01/ef:10", from: common.Address{1}, txHash: common.Hash{1}, nonce: 1,
			benefit: 1000, cost: new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 1, UsedKeccakHashes: 1, UsedPoseidonHashes: 1, UsedPoseidonPaddings: 1, UsedMemAligns: 1, UsedArithmetics: 1, UsedBinaries: 1, UsedSteps: 1},
			usedBytes: 1,
			expectedEfficiencyList: []common.Hash{
				{1},
			},
		},
		{
			name: "Adding from:0x02, tx:0x02/ef:12", from: common.Address{2}, txHash: common.Hash{2}, nonce: 1,
			benefit: 6000, cost: new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 5, UsedKeccakHashes: 5, UsedPoseidonHashes: 5, UsedPoseidonPaddings: 5, UsedMemAligns: 5, UsedArithmetics: 5, UsedBinaries: 5, UsedSteps: 5},
			usedBytes: 5,
			expectedEfficiencyList: []common.Hash{
				{2}, {1},
			},
		},
		{
			name: "Readding from:0x03, tx:0x03/ef:25", from: common.Address{3}, txHash: common.Hash{3}, nonce: 1,
			benefit: 5000, cost: new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 2, UsedKeccakHashes: 2, UsedPoseidonHashes: 2, UsedPoseidonPaddings: 2, UsedMemAligns: 2, UsedArithmetics: 2, UsedBinaries: 2, UsedSteps: 2},
			usedBytes: 2,
			expectedEfficiencyList: []common.Hash{
				{3}, {2}, {1},
			},
		},
		{
			name: "Adding from:0x04, tx:0x04/ef:100", from: common.Address{4}, txHash: common.Hash{4}, nonce: 1,
			benefit: 40000, cost: new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 4, UsedKeccakHashes: 4, UsedPoseidonHashes: 4, UsedPoseidonPaddings: 4, UsedMemAligns: 4, UsedArithmetics: 4, UsedBinaries: 4, UsedSteps: 4},
			usedBytes: 4,
			expectedEfficiencyList: []common.Hash{
				{4}, {3}, {2}, {1},
			},
		},
	}

	processWorkerAddTxTestCases(t, worker, addTxsTC)

	expectedGetBestTx := []common.Hash{{4}, {3}, {1}}
	ct := 0

	for {
		tx := worker.GetBestFittingTx(rc)
		if tx != nil {
			if ct >= len(expectedGetBestTx) {
				t.Fatalf("Error getting more best tx than expected. Expected=%d, Actual=%d", len(expectedGetBestTx), ct+1)
			}
			if tx.HashStr != string(expectedGetBestTx[ct].String()) {
				t.Fatalf("Error GetBestFittingTx(%d). Expected=%s, Actual=%s", ct, expectedGetBestTx[ct].String(), tx.HashStr)
			}
			err := rc.Sub(tx.BatchResources)
			assert.NoError(t, err)

			touch := make(map[common.Address]*state.InfoReadWrite)
			var newNonce uint64 = tx.Nonce + 1
			touch[tx.From] = &state.InfoReadWrite{Address: tx.From, Nonce: &newNonce, Balance: new(big.Int).SetInt64(10)}
			worker.UpdateAfterSingleSuccessfulTxExecution(tx.From, touch)
			ct++
		} else {
			if ct < len(expectedGetBestTx) {
				t.Fatalf("Error expecting more best tx. Expected=%d, Actual=%d", len(expectedGetBestTx), ct)
			}
			break
		}
	}
}

func initWorker(stateMock *StateMock, rcMax batchConstraints, rcWeigth batchResourceWeights) *Worker {
	worker := NewWorker(workerCfg, stateMock, rcMax, rcWeigth)
	return worker
}
