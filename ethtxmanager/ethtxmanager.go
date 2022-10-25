// Package ethtxmanager handles ethereum transactions:  It makes
// calls to send and to aggregate batch, checks possible errors, like wrong nonce or gas limit too low
// and make correct adjustments to request according to it. Also, it tracks transaction receipt and status
// of tx in case tx is rejected and send signals to sequencer/aggregator to resend sequence/batch
package ethtxmanager

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
)

// Client for eth tx manager
type Client struct {
	cfg    Config
	state  stateInterface
	ethMan etherman
}

// New creates new eth tx manager
func New(cfg Config, st stateInterface, ethMan etherman) *Client {
	return &Client{
		cfg:    cfg,
		state:  st,
		ethMan: ethMan,
	}
}

// SyncPendingSequences loads pending sequences from the state and
// sync them with PoE on L1
func (c *Client) SyncPendingSequences() {
	c.groupSequences()
	c.syncSequences()
}

// groupSequences build sequence groups with sequences without group
func (c *Client) groupSequences() {
	ctx := context.Background()

	// get sequences without group
	sequencesWithoutGroup, err := c.state.GetSequencesWithoutGroup(ctx, nil)
	if err != nil {
		log.Errorf("failed to get sequences without group: %v", err)
		return
	}

	// if there is no sequence without group, returns
	if len(sequencesWithoutGroup) == 0 {
		return
	}

	// send the sequences to create the tx
	var tx *types.Transaction
	for {
		// TODO: force the low gas in order to make this tx to be discarded
		// and replaced by a good one by the tx review processing
		gas := uint64(1)
		tx, err = c.ethMan.SequenceBatches(sequencesWithoutGroup, gas, nil)
		if err != nil && err.Error() == core.ErrOversizedData.Error() {
			sequencesWithoutGroup = sequencesWithoutGroup[:len(sequencesWithoutGroup)-1]
		} else if err != nil {
			log.Errorf("failed to send sequence batches: %v", err)
			return
		} else {
			break
		}
	}

	// create a pending sequence group with sequences and tx
	sequenceGroup := state.SequenceGroup{
		TxHash:       tx.Hash(),
		TxNonce:      tx.Nonce(),
		Status:       state.SequenceGroupStatusPending,
		CreatedAt:    time.Now(),
		BatchNumbers: make([]uint64, 0, len(sequencesWithoutGroup)),
	}
	for _, sequence := range sequencesWithoutGroup {
		sequenceGroup.BatchNumbers = append(sequenceGroup.BatchNumbers, sequence.BatchNumber)
	}

	// persist sequence group to start monitoring this tx
	err = c.state.AddSequenceGroup(ctx, sequenceGroup, nil)
	if err != nil {
		log.Errorf("failed to create sequence group: %v", err)
		return
	}
}

func (c *Client) syncSequences() {
	ctx := context.Background()

	pendingSequenceGroups, err := c.state.GetPendingSequenceGroups(ctx, nil)
	if err != nil {
		log.Errorf("failed to get pending sequence groups: %v", err)
		return
	}

	for _, pendingSequenceGroup := range pendingSequenceGroups {
		// check if the tx was already mined
		receipt, err := c.ethMan.GetTxReceipt(ctx, pendingSequenceGroup.TxHash)
		if err != nil && !errors.Is(err, ethereum.NotFound) {
			log.Errorf("failed to get send batch tx receipt, hash %v: %v", pendingSequenceGroup.TxHash.String(), err)
			return
		}
		if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
			err := c.state.SetSequenceGroupAsConfirmed(ctx, pendingSequenceGroup.TxHash, nil)
			if err != nil {
				log.Errorf("failed to set sequence group as confirmed for tx %v: %v", pendingSequenceGroup.TxHash.String(), err)
			}
			return
		}

		// if it was not mined yet, check it against the rules to improve the tx in order to get it mined
		//
		// - check if the timeout since the last time the group was update has expired, if so, update the tx
		lastTimeSequenceWasUpdated := pendingSequenceGroup.CreatedAt
		if pendingSequenceGroup.UpdatedAt != nil {
			lastTimeSequenceWasUpdated = *pendingSequenceGroup.UpdatedAt
		}
		if time.Since(lastTimeSequenceWasUpdated) >= c.cfg.IntervalToReviewSendBatchTx.Duration {
			sequences, err := c.state.GetSequencesByBatchNums(ctx, pendingSequenceGroup.BatchNumbers, nil)
			if err != nil {
				log.Errorf("failed to get sequences by batch numbers: %v", err)
				return
			}

			// use the same nonce
			nonce := big.NewInt(0).SetUint64(pendingSequenceGroup.TxNonce)

			// gets the current nonce to check if the tx nonce is already used
			currentNonce, err := c.ethMan.CurrentNonce()
			if err != nil {
				log.Errorf("failed to get current nonce to check the sequence group tx nonce: %v", err)
				return
			}

			// if the tx nonce was already used, we update it to the current one, so the txs can get mined
			if pendingSequenceGroup.TxNonce <= currentNonce {
				nonce = big.NewInt(0).SetUint64(currentNonce)
			}

			// use 0 gas to allow the system to get the suggestion from the network
			gas := uint64(0)
			tx, err := c.ethMan.SequenceBatches(sequences, gas, nonce)
			if err != nil {
				log.Errorf("failed to send sequence batches: %v", err)
				return
			}

			log.Infof("updating tx for sequence group related to batches %v from %v to %v",
				pendingSequenceGroup.BatchNumbers, pendingSequenceGroup.TxHash.String(), tx.Hash().String())

			err = c.state.UpdateSequenceGroupTx(ctx, pendingSequenceGroup.TxHash, *tx, nil)
			if err != nil {
				log.Errorf("failed to update sequence group from %v to %v: %v", pendingSequenceGroup.TxHash.String(), tx.Hash().String(), err)
				return
			}
		}
	}
}

// VerifyBatch TODO REMOVE
func (c *Client) VerifyBatch(ctx context.Context, batchNum uint64, proof *pb.GetProofResponse) error {
	return fmt.Errorf("not implemented yet")
}
