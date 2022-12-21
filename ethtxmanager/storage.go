package ethtxmanager

import (
	"context"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/ethereum/go-ethereum/core/types"
)

// txsBufferSize txs channel limit
const txsBufferSize = 1000

// storage hold txs to be managed
type storage struct {
	// txs enqueued to be monitored
	txs chan enqueuedTx
}

// newStorage creates a new instance of storage
func newStorage() *storage {
	return &storage{
		txs: make(chan enqueuedTx, txsBufferSize),
	}
}

// enqueueSequences adds a tx to the enqueued txs to sequence batches
func (s *storage) enqueueSequences(ctx context.Context, st state, e etherman, cfg Config, sequences []ethmanTypes.Sequence) (*types.Transaction, error) {
	tx, err := e.SequenceBatches(ctx, sequences, 0, nil, nil, true)
	if err != nil {
		return nil, err
	}

	s.txs <- &enqueuedSequencesTx{
		baseEnqueuedTx: baseEnqueuedTx{tx: tx, waitDuration: cfg.FrequencyForResendingFailedSendBatches.Duration},
		state:          st,
		cfg:            cfg,
		sequences:      sequences,
	}

	return tx, nil
}

// enqueueVerifyBatches adds a tx to the enqueued txs to verify batches
func (s *storage) enqueueVerifyBatches(ctx context.Context, st state, e etherman, cfg Config, lastVerifiedBatch uint64, finalBatchNum uint64, inputs *ethmanTypes.FinalProofInputs) (*types.Transaction, error) {
	tx, err := e.TrustedVerifyBatches(ctx, lastVerifiedBatch, finalBatchNum, inputs, 0, nil, nil, true)
	if err != nil {
		return nil, err
	}

	s.txs <- &enqueuedVerifyBatchesTx{
		baseEnqueuedTx:    baseEnqueuedTx{tx: tx, waitDuration: cfg.FrequencyForResendingFailedSendBatches.Duration},
		state:             st,
		cfg:               cfg,
		lastVerifiedBatch: lastVerifiedBatch,
		finalBatchNum:     finalBatchNum,
		inputs:            inputs,
	}

	return tx, nil
}

// Next returns the next enqueued tx
func (s *storage) Next() enqueuedTx {
	return <-s.txs
}
