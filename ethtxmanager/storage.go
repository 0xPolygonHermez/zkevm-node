package ethtxmanager

import (
	"context"
	"encoding/json"
	"reflect"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/ethereum/go-ethereum/core/types"
)

// txsBufferSize txs channel limit
const txsBufferSize = 1000

type persistedEnqueuedTx struct {
	Hash         string
	RawTx        string
	Data         json.RawMessage
	Type         string
	Status       enqueuedTxStatus
	WaitDuration int64
}

// storage hold txs to be managed
type storage struct {
	// txs enqueued to be monitored
	txs chan persistedEnqueuedTx
}

// newStorage creates a new instance of storage
func newStorage() *storage {
	return &storage{
		txs: make(chan persistedEnqueuedTx, txsBufferSize),
	}
}

// enqueueSequences adds a tx to the enqueued txs to sequence batches
func (s *storage) enqueueSequences(ctx context.Context, e etherman, cfg Config, sequences []ethmanTypes.Sequence) (*types.Transaction, error) {
	tx, err := e.SequenceBatches(ctx, sequences, 0, nil, nil, true)
	if err != nil {
		return nil, err
	}

	etx := &enqueuedSequencesTx{
		baseEnqueuedTx: baseEnqueuedTx{tx: tx, waitDuration: cfg.FrequencyForResendingFailedSendBatches.Duration},
		sequences:      sequences,
	}

	p, err := s.ToPersistedEnqueuedTx(etx)
	if err != nil {
		return nil, err
	}

	err = s.enqueue(p)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// enqueueVerifyBatches adds a tx to the enqueued txs to verify batches
func (s *storage) enqueueVerifyBatches(ctx context.Context, e etherman, cfg Config, lastVerifiedBatch uint64, finalBatchNum uint64, inputs *ethmanTypes.FinalProofInputs) (*types.Transaction, error) {
	tx, err := e.TrustedVerifyBatches(ctx, lastVerifiedBatch, finalBatchNum, inputs, 0, nil, nil, true)
	if err != nil {
		return nil, err
	}

	etx := &enqueuedVerifyBatchesTx{
		baseEnqueuedTx:    baseEnqueuedTx{tx: tx, waitDuration: cfg.FrequencyForResendingFailedVerifyBatch.Duration},
		lastVerifiedBatch: lastVerifiedBatch,
		finalBatchNum:     finalBatchNum,
		inputs:            inputs,
	}

	p, err := s.ToPersistedEnqueuedTx(etx)
	if err != nil {
		return nil, err
	}

	err = s.enqueue(p)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *storage) enqueue(pEtx persistedEnqueuedTx) error {
	s.txs <- pEtx
	return nil
}

// Next returns the next enqueued tx
func (s *storage) Next() (enqueuedTx, error) {
	pEtx := <-s.txs
	var etx enqueuedTx
	if pEtx.Type == reflect.TypeOf(&enqueuedSequencesTx{}).String() {
		etx = &enqueuedSequencesTx{}
	} else if pEtx.Type == reflect.TypeOf(&enqueuedVerifyBatchesTx{}).String() {
		etx = &enqueuedVerifyBatchesTx{}
	}

	err := etx.Load(pEtx)
	if err != nil {
		return nil, err
	}

	return etx, nil
}

// Confirm updates the enqueued tx to confirmed
func (s *storage) Confirm(ctx context.Context, etx enqueuedTx) error {
	return nil
}

// ToPersistedEnqueuedTx converts an instance of enqueuedTx into
// a struct that can be used to store the data into the storage
func (s *storage) ToPersistedEnqueuedTx(etx enqueuedTx) (persistedEnqueuedTx, error) {
	result := persistedEnqueuedTx{}

	b, err := etx.Tx().MarshalBinary()
	if err != nil {
		return result, err
	}
	rawTx := hex.EncodeToHex(b)

	data, err := etx.Data()
	if err != nil {
		return result, err
	}

	result = persistedEnqueuedTx{
		Hash:         etx.Tx().Hash().String(),
		RawTx:        rawTx,
		Type:         reflect.TypeOf(etx).String(),
		Data:         data,
		Status:       etx.Status(),
		WaitDuration: etx.WaitDuration().Nanoseconds(),
	}

	return result, nil
}
