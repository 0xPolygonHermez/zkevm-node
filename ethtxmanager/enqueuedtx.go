package ethtxmanager

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
// enqueuedTxStatusPending    = enqueuedTxStatus("pending")
// enqueuedTxStatusProcessing = enqueuedTxStatus("processing")
// enqueuedTxStatusDiscarded  = enqueuedTxStatus("discarded")
// enqueuedTxStatusConfirmed  = enqueuedTxStatus("confirmed")
)

type enqueuedTx interface {
	Tx() *types.Transaction
	RenewTxIfNeeded(context.Context, etherman) error
	Wait()
	WaitDuration() time.Duration
	WaitSync(ctx context.Context, st state, cfg Config) error
	Status() enqueuedTxStatus
	Data() (json.RawMessage, error)
	SetData(json.RawMessage) error
	Load(persistedEnqueuedTx) error
	OnConfirmation()
}

type enqueuedTxStatus string

type baseEnqueuedTx struct {
	tx           *types.Transaction
	waitDuration time.Duration
	status       enqueuedTxStatus
}

// Tx returns the internal tx
func (etx *baseEnqueuedTx) Tx() *types.Transaction {
	return etx.tx
}

// Wait waits for the WaitDuration
func (etx *baseEnqueuedTx) Wait() {
	time.Sleep(etx.waitDuration)
}

// Tx returns the time it needs to wait to try to reprocess
func (etx *baseEnqueuedTx) WaitDuration() time.Duration {
	return etx.waitDuration
}

// Status of this enqueued tx in the queue
func (etx *baseEnqueuedTx) Status() enqueuedTxStatus {
	return etx.status
}

// Load base data from persisted instance
func (etx *baseEnqueuedTx) Load(pEtx persistedEnqueuedTx) error {
	b, err := hex.DecodeHex(pEtx.RawTx)
	if err != nil {
		return err
	}

	var tx = &types.Transaction{}
	if err := tx.UnmarshalBinary(b); err != nil {
		return err
	}

	etx.tx = tx
	etx.waitDuration = time.Duration(pEtx.WaitDuration)
	etx.status = pEtx.Status
	return nil
}

// enqueuedSequencesTx represents a ethereum tx created to
// sequence batches that can be enqueued to be monitored
type enqueuedSequencesTx struct {
	baseEnqueuedTx
	sequences []ethmanTypes.Sequence
}

// RenewTxIfNeeded checks for information in the inner tx and renews it
// if needed, for example changes the nonce is it realizes the nonce was
// already used or updates the gas price if the network has changed the
// prices since the tx was created
func (etx *enqueuedSequencesTx) RenewTxIfNeeded(ctx context.Context, e etherman) error {
	// nonce, err := e.CurrentNonce(ctx)
	// if err != nil {
	// 	return err
	// }
	// if etx.Tx().Nonce() < nonce {
	// 	err = etx.renewNonce(ctx, e)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	tx, err := e.EstimateGasSequenceBatches(etx.sequences)
	if err != nil {
		return err
	}
	if tx.Gas() > etx.Tx().Gas() {
		err = etx.renewGas(ctx, e)
		if err != nil {
			return err
		}
	}
	return nil
}

// // RenewNonce renews the inner TX nonce
// func (etx *enqueuedSequencesTx) renewNonce(ctx context.Context, e etherman) error {
// 	oldTx := etx.Tx()
// 	tx, err := e.SequenceBatches(ctx, etx.sequences, oldTx.Gas(), oldTx.GasPrice(), nil, true)
// 	if err != nil {
// 		return err
// 	}
// 	etx.baseEnqueuedTx.tx = tx
// 	return nil
// }

// RenewGasPrice renews the inner TX Gas Price
func (etx *enqueuedSequencesTx) renewGas(ctx context.Context, e etherman) error {
	oldTx := etx.Tx()
	oldNonce := big.NewInt(0).SetUint64(oldTx.Nonce())
	tx, err := e.SequenceBatches(ctx, etx.sequences, oldTx.Gas(), nil, oldNonce, true)
	if err != nil {
		return err
	}
	etx.baseEnqueuedTx.tx = tx
	return nil
}

// WaitSync checks if the sequences were already synced into the state
func (etx *enqueuedSequencesTx) WaitSync(ctx context.Context, st state, cfg Config) error {
	return st.WaitSequencingTxToBeSynced(ctx, etx.Tx(), cfg.WaitTxToBeSynced.Duration)
}

// Data provides information about the sequences to be stored
// and retrieved from storage
func (etx *enqueuedSequencesTx) Data() (json.RawMessage, error) {
	rawData := make([]map[string]interface{}, 0, len(etx.sequences))

	for _, sequence := range etx.sequences {
		txs := make([]string, 0, len(sequence.Txs))
		for _, tx := range sequence.Txs {
			b, err := tx.MarshalBinary()
			if err != nil {
				return nil, err
			}
			txRlp := hex.EncodeToHex(b)
			txs = append(txs, txRlp)
		}

		seqData := map[string]interface{}{
			"global_exit_root":    sequence.GlobalExitRoot.String(),
			"state_root":          sequence.StateRoot.String(),
			"local_exit_root":     sequence.LocalExitRoot.String(),
			"acc_input_hash":      sequence.AccInputHash.String(),
			"timestamp":           sequence.Timestamp,
			"txs":                 txs,
			"is_sequence_too_big": sequence.IsSequenceTooBig,
		}

		rawData = append(rawData, seqData)
	}

	data, err := json.Marshal(rawData)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// SetData sets the data, used by storage to retrieve stored data
func (etx *enqueuedSequencesTx) SetData(data json.RawMessage) error {
	rawData := []map[string]interface{}{}
	err := json.Unmarshal(data, &rawData)
	if err != nil {
		return err
	}

	sequences := make([]ethmanTypes.Sequence, 0, len(rawData))

	for _, seqData := range rawData {
		rawTxs := seqData["txs"].([]interface{})
		txs := make([]types.Transaction, 0, len(rawTxs))
		for _, rawTx := range rawTxs {
			b, err := hex.DecodeHex(rawTx.(string))
			if err != nil {
				return err
			}

			var tx = &types.Transaction{}
			if err := tx.UnmarshalBinary(b); err != nil {
				return err
			}
			txs = append(txs, *tx)
		}

		sequence := ethmanTypes.Sequence{
			GlobalExitRoot:   common.HexToHash(seqData["global_exit_root"].(string)),
			StateRoot:        common.HexToHash(seqData["state_root"].(string)),
			LocalExitRoot:    common.HexToHash(seqData["local_exit_root"].(string)),
			AccInputHash:     common.HexToHash(seqData["acc_input_hash"].(string)),
			Timestamp:        int64(seqData["timestamp"].(float64)),
			Txs:              txs,
			IsSequenceTooBig: seqData["is_sequence_too_big"].(bool),
		}
		sequences = append(sequences, sequence)
	}

	etx.sequences = sequences

	return nil
}

// Loads persisted sequence data
func (etx *enqueuedSequencesTx) Load(pEtx persistedEnqueuedTx) error {
	err := etx.baseEnqueuedTx.Load(pEtx)
	if err != nil {
		return err
	}
	err = etx.SetData(pEtx.Data)
	if err != nil {
		return err
	}
	return nil
}

// OnConfirmation is called after the transaction is confirmed
func (etx *enqueuedSequencesTx) OnConfirmation() {
	log.Infof("sequence sent to L1 successfully. Tx hash: %s", etx.Tx().Hash().String())
}

// enqueuedVerifyBatchesTx represents a ethereum tx created to
// verify batches that can be enqueued to be monitored
type enqueuedVerifyBatchesTx struct {
	baseEnqueuedTx
	lastVerifiedBatch uint64
	finalBatchNum     uint64
	inputs            *ethmanTypes.FinalProofInputs
}

// RenewTxIfNeeded checks for information in the inner tx and renews it
// if needed, for example changes the nonce is it realizes the nonce was
// already used or updates the gas price if the network has changed the
// prices since the tx was created
func (etx *enqueuedVerifyBatchesTx) RenewTxIfNeeded(ctx context.Context, e etherman) error {
	// nonce, err := e.CurrentNonce(ctx)
	// if err != nil {
	// 	return err
	// }
	// if etx.Tx().Nonce() < nonce {
	// 	err = etx.renewNonce(ctx, e)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	estimatedGas, err := e.EstimateGasForVerifyBatches(etx.lastVerifiedBatch, etx.finalBatchNum, etx.inputs)
	if err != nil {
		return err
	}
	if estimatedGas > etx.Tx().Gas() {
		err = etx.renewGas(ctx, e)
		if err != nil {
			return err
		}
	}
	return nil
}

// // RenewNonce renews the inner TX nonce
// func (etx *enqueuedVerifyBatchesTx) renewNonce(ctx context.Context, e etherman) error {
// 	oldTx := etx.Tx()
// 	tx, err := e.TrustedVerifyBatches(ctx, etx.lastVerifiedBatch, etx.finalBatchNum, etx.inputs, oldTx.Gas(), oldTx.GasPrice(), nil, true)
// 	if err != nil {
// 		return err
// 	}
// 	etx.baseEnqueuedTx.tx = tx
// 	return nil
// }

// RenewGasPrice renews the inner TX Gas Price
func (etx *enqueuedVerifyBatchesTx) renewGas(ctx context.Context, e etherman) error {
	oldTx := etx.Tx()
	oldNonce := big.NewInt(0).SetUint64(oldTx.Nonce())
	tx, err := e.VerifyBatches(ctx, etx.lastVerifiedBatch, etx.finalBatchNum, etx.inputs, oldTx.Gas(), nil, oldNonce, true)
	if err != nil {
		return err
	}
	etx.baseEnqueuedTx.tx = tx
	return nil
}

// WaitSync checks if the sequences were already synced into the state
func (etx *enqueuedVerifyBatchesTx) WaitSync(ctx context.Context, st state, cfg Config) error {
	return st.WaitVerifiedBatchToBeSynced(ctx, etx.finalBatchNum, cfg.WaitTxToBeSynced.Duration)
}

// Data provides information about the proof to be stored
// and retrieved from storage
func (etx *enqueuedVerifyBatchesTx) Data() (json.RawMessage, error) {
	panic("not implemented yet")
}

// SetData sets the data, used by storage to retrieve stored data
func (etx *enqueuedVerifyBatchesTx) SetData(data json.RawMessage) error {
	panic("not implemented yet")
}

// Loads persisted sequence data
func (etx *enqueuedVerifyBatchesTx) Load(pEtx persistedEnqueuedTx) error {
	err := etx.baseEnqueuedTx.Load(pEtx)
	if err != nil {
		return err
	}
	err = etx.SetData(pEtx.Data)
	if err != nil {
		return err
	}
	return nil
}

// OnConfirmation is called after the transaction is confirmed
func (etx *enqueuedVerifyBatchesTx) OnConfirmation() {
	log.Infof("Final proof for batches [%d-%d] verified in transaction [%v]", etx.lastVerifiedBatch, etx.finalBatchNum, etx.Tx().Hash().String())
}
