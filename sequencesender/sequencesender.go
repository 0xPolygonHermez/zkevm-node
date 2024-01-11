package sequencesender

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

const (
	ethTxManagerOwner = "sequencer"
	monitoredIDFormat = "sequence-from-%v-to-%v"
)

var (
	// ErrOversizedData is returned if the input data of a transaction is greater
	// than some meaningful limit a user might use. This is not a consensus error
	// making the transaction invalid, rather a DOS protection.
	ErrOversizedData = errors.New("oversized data")
)

// SequenceSender represents a sequence sender
type SequenceSender struct {
	cfg          Config
	state        stateInterface
	ethTxManager ethTxManager
	etherman     etherman
	eventLog     *event.EventLog
	da           dataAbilitier
}

// New inits sequence sender
func New(cfg Config, state stateInterface, etherman etherman, manager ethTxManager, eventLog *event.EventLog, da dataAbilitier) (*SequenceSender, error) {
	return &SequenceSender{
		cfg:          cfg,
		state:        state,
		etherman:     etherman,
		ethTxManager: manager,
		eventLog:     eventLog,
		da:           da,
	}, nil
}

// Start starts the sequence sender
func (s *SequenceSender) Start(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.WaitPeriodSendSequence.Duration)
	for {
		s.tryToSendSequence(ctx, ticker)
	}
}

func (s *SequenceSender) tryToSendSequence(ctx context.Context, ticker *time.Ticker) {
	retry := false
	// process monitored sequences before starting a next cycle
	s.ethTxManager.ProcessPendingMonitoredTxs(ctx, ethTxManagerOwner, func(result ethtxmanager.MonitoredTxResult, dbTx pgx.Tx) {
		if result.Status == ethtxmanager.MonitoredTxStatusFailed {
			retry = true
			mTxResultLogger := ethtxmanager.CreateMonitoredTxResultLogger(ethTxManagerOwner, result)
			mTxResultLogger.Error("failed to send sequence, TODO: review this fatal and define what to do in this case")
		}
	}, nil)

	if retry {
		return
	}

	// Check if synchronizer is up to date
	if !s.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		waitTick(ctx, ticker)
		return
	}

	// Check if should send sequence to L1
	log.Infof("getting sequences to send")
	sequences, err := s.getSequencesToSend(ctx)
	if err != nil || len(sequences) == 0 {
		if err != nil {
			log.Errorf("error getting sequences: %v", err)
		} else {
			log.Info("waiting for sequences to be worth sending to L1")
		}
		waitTick(ctx, ticker)
		return
	}

	lastVirtualBatchNum, err := s.state.GetLastVirtualBatchNum(ctx, nil)
	if err != nil {
		log.Errorf("failed to get last virtual batch num, err: %v", err)
		return
	}

	// Send sequences to L1
	sequenceCount := len(sequences)
	log.Infof(
		"sending sequences to L1. From batch %d to batch %d",
		lastVirtualBatchNum+1, lastVirtualBatchNum+uint64(sequenceCount),
	)
	metrics.SequencesSentToL1(float64(sequenceCount))

	// add sequence to be monitored
	dataAvailabilityMessage, err := s.da.PostSequence(ctx, sequences)
	if err != nil {
		log.Error("error posting sequences to the data availability protocol: ", err)
		return
	}
	to, data, err := s.etherman.BuildSequenceBatchesTxData(s.cfg.SenderAddress, sequences, s.cfg.L2Coinbase, dataAvailabilityMessage)
	if err != nil {
		log.Error("error estimating new sequenceBatches to add to eth tx manager: ", err)
		return
	}
	firstSequence := sequences[0]
	lastSequence := sequences[len(sequences)-1]
	monitoredTxID := fmt.Sprintf(monitoredIDFormat, firstSequence.BatchNumber, lastSequence.BatchNumber)
	err = s.ethTxManager.Add(ctx, ethTxManagerOwner, monitoredTxID, s.cfg.SenderAddress, to, nil, data, s.cfg.GasOffset, nil)
	if err != nil {
		mTxLogger := ethtxmanager.CreateLogger(ethTxManagerOwner, monitoredTxID, s.cfg.SenderAddress, to)
		mTxLogger.Errorf("error to add sequences tx to eth tx manager: ", err)
		return
	}
}

// getSequencesToSend generates an array of sequences to be send to L1.
// If the array is empty, it doesn't necessarily mean that there are no sequences to be sent,
// it could be that it's not worth it to do so yet.
func (s *SequenceSender) getSequencesToSend(ctx context.Context) ([]types.Sequence, error) {
	lastVirtualBatchNum, err := s.state.GetLastVirtualBatchNum(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get last virtual batch num, err: %w", err)
	}

	currentBatchNumToSequence := lastVirtualBatchNum + 1
	sequences := []types.Sequence{}

	// Add sequences until too big for a single L1 tx or last batch is reached
	for {
		//Check if the next batch belongs to a new forkid, in this case we need to stop sequencing as we need to
		//wait the upgrade of forkid is completed and s.cfg.NumBatchForkIdUpgrade is disabled (=0) again
		if (s.cfg.ForkUpgradeBatchNumber != 0) && (currentBatchNumToSequence == (s.cfg.ForkUpgradeBatchNumber + 1)) {
			return nil, fmt.Errorf("aborting sequencing process as we reached the batch %d where a new forkid is applied (upgrade)", s.cfg.ForkUpgradeBatchNumber+1)
		}

		// Check if batch is closed
		isClosed, err := s.state.IsBatchClosed(ctx, currentBatchNumToSequence, nil)
		if err != nil {
			return nil, err
		}
		if !isClosed {
			// Reached current (WIP) batch
			break
		}
		// Add new sequence
		batch, err := s.state.GetBatchByNumber(ctx, currentBatchNumToSequence, nil)
		if err != nil {
			return nil, err
		}

		seq := types.Sequence{
			GlobalExitRoot: batch.GlobalExitRoot,   //TODO: set empty for regular batches
			Timestamp:      batch.Timestamp.Unix(), //TODO: set empty for regular batches
			BatchL2Data:    batch.BatchL2Data,
			BatchNumber:    batch.BatchNumber,
		}

		if batch.ForcedBatchNum != nil {
			//TODO: Assign GER, timestamp(forcedAt) and l1block.parentHash to seq
			forcedBatch, err := s.state.GetForcedBatch(ctx, *batch.ForcedBatchNum, nil)
			if err != nil {
				return nil, err
			}

			// Get L1 block for the forced batch
			fbL1Block, err := s.state.GetBlockByNumber(ctx, forcedBatch.BlockNumber, nil)
			if err != nil {
				return nil, err
			}

			seq.GlobalExitRoot = forcedBatch.GlobalExitRoot
			seq.ForcedBatchTimestamp = forcedBatch.ForcedAt.Unix()
			seq.PrevBlockHash = fbL1Block.ParentHash
		}

		sequences = append(sequences, seq)
		// Check if can be send
		if len(sequences) == int(s.cfg.MaxBatchesForL1) {
			log.Info(
				"sequence should be sent to L1, because MaxBatchesForL1 (%d) has been reached",
				s.cfg.MaxBatchesForL1,
			)
			return sequences, nil
		}

		// Increase batch num for next iteration
		currentBatchNumToSequence++
	}

	// Reached latest batch. Decide if it's worth to send the sequence, or wait for new batches
	if len(sequences) == 0 {
		log.Info("no batches to be sequenced")
		return nil, nil
	}

	lastBatchVirtualizationTime, err := s.state.GetTimeForLatestBatchVirtualization(ctx, nil)
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		log.Warnf("failed to get last l1 interaction time, err: %v. Sending sequences as a conservative approach", err)
		return sequences, nil
	}
	if lastBatchVirtualizationTime.Before(time.Now().Add(-s.cfg.LastBatchVirtualizationTimeMaxWaitPeriod.Duration)) {
		// TODO: implement check profitability
		// if s.checker.IsSendSequencesProfitable(new(big.Int).SetUint64(estimatedGas), sequences) {
		log.Info("sequence should be sent to L1, because too long since didn't send anything to L1")
		return sequences, nil
		//}
	}

	log.Info("not enough time has passed since last batch was virtualized, and the sequence could be bigger")
	return nil, nil
}

func waitTick(ctx context.Context, ticker *time.Ticker) {
	select {
	case <-ticker.C:
		// nothing
	case <-ctx.Done():
		return
	}
}

func (s *SequenceSender) isSynced(ctx context.Context) bool {
	lastSyncedBatchNum, err := s.state.GetLastVirtualBatchNum(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get last isSynced batch, err: %v", err)
		return false
	}
	lastBatchNum, err := s.state.GetLastBatchNumber(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get last batch num, err: %v", err)
		return false
	}
	if lastBatchNum > lastSyncedBatchNum {
		return true
	}
	lastEthBatchNum, err := s.etherman.GetLatestBatchNumber()
	if err != nil {
		log.Errorf("failed to get last eth batch, err: %v", err)
		return false
	}
	if lastSyncedBatchNum < lastEthBatchNum {
		log.Infof("waiting for the state to be isSynced, lastSyncedBatchNum: %d, lastEthBatchNum: %d", lastSyncedBatchNum, lastEthBatchNum)
		return false
	}

	return true
}
