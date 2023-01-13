package sequencer

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

var (
	f             *finalizer
	dbManagerMock = new(DbManagerMock)
	executorMock  = new(StateMock)
	workerMock    = new(WorkerMock)
	bc            = batchConstraints{
		MaxTxsPerBatch:       150,
		MaxBatchBytesSize:    150000,
		MaxCumulativeGasUsed: 30000000,
		MaxKeccakHashes:      468,
		MaxPoseidonHashes:    279620,
		MaxPoseidonPaddings:  149796,
		MaxMemAligns:         262144,
		MaxArithmetics:       262144,
		MaxBinaries:          262144,
		MaxSteps:             8388608,
	}
	txsStore = TxsStore{
		Ch: make(chan *txToStore, 1),
		Wg: new(sync.WaitGroup),
	}
	closingSignalCh = ClosingSignalCh{
		ForcedBatchCh:        make(chan state.ForcedBatch),
		GERCh:                make(chan common.Hash),
		L2ReorgCh:            make(chan L2ReorgEvent),
		SendingToL1TimeoutCh: make(chan bool),
	}
	cfg = FinalizerCfg{
		GERDeadlineTimeoutInSec: types.Duration{
			Duration: 60,
		},
		ForcedBatchDeadlineTimeoutInSec: types.Duration{
			Duration: 60,
		},
		SendingToL1DeadlineTimeoutInSec: types.Duration{
			Duration: 60,
		},
		SleepDurationInMs: types.Duration{
			Duration: 60,
		},
		ResourcePercentageToCloseBatch: 90,
	}
	seqAddr  = common.Address{}
	ctx      = context.Background()
	hash     = common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f2")
	hash2    = common.HexToHash("0xe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	sender   = common.HexToAddress("0x3445324")
	isSynced = func(ctx context.Context) bool {
		return true
	}
	tx1 = ethTypes.NewTransaction(0, common.HexToAddress("0"), big.NewInt(0), 0, big.NewInt(0), []byte("aaa"))
	tx2 = ethTypes.NewTransaction(1, common.HexToAddress("1"), big.NewInt(1), 0, big.NewInt(1), []byte("bbb"))
)

func testNow() time.Time {
	return time.Unix(0, 0)
}

func TestNewFinalizer(t *testing.T) {
	// arrange and act
	f = newFinalizer(cfg, workerMock, dbManagerMock, executorMock, seqAddr, isSynced, closingSignalCh, txsStore, bc)

	// assert
	assert.NotNil(t, f)
	assert.Equal(t, f.cfg, cfg)
	assert.Equal(t, f.worker, workerMock)
	assert.Equal(t, f.dbManager, dbManagerMock)
	assert.Equal(t, f.executor, executorMock)
	assert.Equal(t, f.sequencerAddress, seqAddr)
	assert.Equal(t, f.closingSignalCh, closingSignalCh)
	assert.Equal(t, f.txsStore, txsStore)
	assert.Equal(t, f.batchConstraints, bc)
}

func Test_reprocessBatch(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	n := uint(2)
	dbManagerMock.On("GetLastNBatches", ctx, n).Return([]*state.Batch{
		{
			StateRoot:    hash,
			AccInputHash: hash,
		},
	}, nil)
	processRequest := state.ProcessRequest{
		BatchNumber:     f.batch.batchNumber,
		GlobalExitRoot:  hash,
		OldStateRoot:    hash,
		OldAccInputHash: hash,
		Coinbase:        seqAddr,
		Timestamp:       f.batch.timestamp,
		Caller:          state.SequencerCallerLabel,
	}
	executorMock.On("ProcessBatch", ctx, processRequest).Return(&state.ProcessBatchResponse{
		Error:            nil,
		IsBatchProcessed: true,
		TouchedAddresses: nil,
	}, nil)

	// act
	err := f.reprocessBatch(ctx)

	// assert
	assert.NoError(t, err)
	dbManagerMock.AssertExpectations(t)
	executorMock.AssertExpectations(t)
}

func Test_prepareProcessRequestFromState(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	n := uint(2)
	dbManagerMock.On("GetLastNBatches", ctx, n).Return([]*state.Batch{
		{
			StateRoot:    hash,
			AccInputHash: hash,
		},
	}, nil)
	expected := state.ProcessRequest{
		BatchNumber:     f.batch.batchNumber,
		GlobalExitRoot:  hash,
		OldStateRoot:    hash,
		OldAccInputHash: hash,
		Coinbase:        seqAddr,
		Timestamp:       f.batch.timestamp,
		Caller:          state.SequencerCallerLabel,
	}

	// act
	actual, err := f.prepareProcessRequestFromState(ctx)
	if err != nil {
		return
	}

	// assert
	assert.Equal(t, expected, actual)
}

func Test_isCurrBatchAboveLimitWindow_Is(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	f.batch.remainingResources.zKCounters.CumulativeGasUsed = f.getConstraintThresholdUint64(bc.MaxCumulativeGasUsed) + 1

	// act
	result := f.isCurrBatchAboveLimitWindow()

	// assert
	assert.True(t, result)
}

func Test_isCurrBatchAboveLimitWindow_Not(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	f.batch.remainingResources.bytes = 1
	f.batch.remainingResources.zKCounters.CumulativeGasUsed = 1
	f.batch.remainingResources.zKCounters.UsedKeccakHashes = 1
	f.batch.remainingResources.zKCounters.UsedPoseidonHashes = 1
	f.batch.remainingResources.zKCounters.UsedPoseidonPaddings = 1
	f.batch.remainingResources.zKCounters.UsedMemAligns = 1
	f.batch.remainingResources.zKCounters.UsedArithmetics = 1
	f.batch.remainingResources.zKCounters.UsedBinaries = 1
	f.batch.remainingResources.zKCounters.UsedSteps = 1

	// act
	result := f.isCurrBatchAboveLimitWindow()

	// assert
	assert.False(t, result)
}

func Test_setNextForcedBatchDeadline(t *testing.T) {
	// arrange
	f = setupFinalizer(false)
	now = testNow
	expected := now().Unix() + int64(f.cfg.ForcedBatchDeadlineTimeoutInSec.Duration)

	// act
	f.setNextForcedBatchDeadline()

	// assert
	assert.Equal(t, expected, f.nextForcedBatchDeadline)
}

func Test_setNextGERDeadline(t *testing.T) {
	// arrange
	f = setupFinalizer(false)
	now = testNow
	expected := now().Unix() + int64(f.cfg.GERDeadlineTimeoutInSec.Duration)

	// act
	f.setNextGERDeadline()

	// assert
	assert.Equal(t, expected, f.nextGERDeadline)
}

func Test_setNextSendingToL1Deadline(t *testing.T) {
	// arrange
	f = setupFinalizer(false)
	now = testNow
	expected := now().Unix() + int64(f.cfg.SendingToL1DeadlineTimeoutInSec.Duration)

	// act
	f.setNextSendingToL1Deadline()

	// assert
	assert.Equal(t, expected, f.nextSendingToL1Deadline)
}

func Test_getConstraintThresholdUint64(t *testing.T) {
	// arrange
	f = setupFinalizer(false)
	input := uint64(100)
	expect := input * uint64(f.cfg.ResourcePercentageToCloseBatch) / 100

	// act
	result := f.getConstraintThresholdUint64(input)

	// assert
	assert.Equal(t, result, expect)
}

func Test_getConstraintThresholdUint32(t *testing.T) {
	// arrange
	f = setupFinalizer(false)
	input := uint32(100)
	expect := uint32(input * f.cfg.ResourcePercentageToCloseBatch / 100)

	// act
	result := f.getConstraintThresholdUint32(input)

	// assert
	assert.Equal(t, result, expect)
}

func Test_getRemainingResources(t *testing.T) {
	// act
	remainingResources := getMaxRemainingResources(bc)

	// assert
	assert.Equal(t, remainingResources.zKCounters.CumulativeGasUsed, bc.MaxCumulativeGasUsed)
	assert.Equal(t, remainingResources.zKCounters.UsedKeccakHashes, bc.MaxKeccakHashes)
	assert.Equal(t, remainingResources.zKCounters.UsedPoseidonHashes, bc.MaxPoseidonHashes)
	assert.Equal(t, remainingResources.zKCounters.UsedPoseidonPaddings, bc.MaxPoseidonPaddings)
	assert.Equal(t, remainingResources.zKCounters.UsedMemAligns, bc.MaxMemAligns)
	assert.Equal(t, remainingResources.zKCounters.UsedArithmetics, bc.MaxArithmetics)
	assert.Equal(t, remainingResources.zKCounters.UsedBinaries, bc.MaxBinaries)
	assert.Equal(t, remainingResources.zKCounters.UsedSteps, bc.MaxSteps)
	assert.Equal(t, remainingResources.bytes, bc.MaxBatchBytesSize)
}

func setupFinalizer(withWipBatch bool) *finalizer {
	wipBatch := new(WipBatch)
	if withWipBatch {
		wipBatch = &WipBatch{
			batchNumber:        1,
			coinbase:           seqAddr,
			initialStateRoot:   hash,
			stateRoot:          hash2,
			timestamp:          uint64(time.Now().Unix()),
			globalExitRoot:     hash,
			remainingResources: getMaxRemainingResources(bc),
		}
	}
	return &finalizer{
		cfg:                cfg,
		txsStore:           txsStore,
		closingSignalCh:    closingSignalCh,
		isSynced:           isSynced,
		sequencerAddress:   seqAddr,
		worker:             workerMock,
		dbManager:          dbManagerMock,
		executor:           executorMock,
		sharedResourcesMux: new(sync.RWMutex),
		batch:              wipBatch,
		batchConstraints:   bc,
		processRequest:     state.ProcessRequest{},
		// closing signals
		nextGER:                   common.Hash{},
		nextGERDeadline:           0,
		nextGERMux:                new(sync.RWMutex),
		nextForcedBatches:         make([]state.ForcedBatch, 0),
		nextForcedBatchDeadline:   0,
		nextForcedBatchesMux:      new(sync.RWMutex),
		nextSendingToL1Deadline:   0,
		nextSendingToL1TimeoutMux: new(sync.RWMutex),
	}
}
