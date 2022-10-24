// Package ethtxmanager handles ethereum transactions:  It makes
// calls to send and to aggregate batch, checks possible errors, like wrong nonce or gas limit too low
// and make correct adjustments to request according to it. Also, it tracks transaction receipt and status
// of tx in case tx is rejected and send signals to sequencer/aggregator to resend sequence/batch
package ethtxmanager

import (
	"context"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum"
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
	tx, err := c.ethMan.SequenceBatches(sequencesWithoutGroup, 0, nil)
	if err != nil {
		log.Errorf("failed to send sequence batches: %v", err)
		return
	}

	// create a pending sequence group with sequences and tx
	sequenceGroup := state.SequenceGroup{
		Tx:           *tx,
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
		receipt, err := c.ethMan.GetTxReceipt(ctx, pendingSequenceGroup.Tx.Hash())
		if err != nil && !errors.Is(err, ethereum.NotFound) {
			log.Errorf("failed to get send batch tx receipt, hash %v: %v", pendingSequenceGroup.Tx.Hash().String(), err)
			return
		}
		if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
			err := c.state.SetSequenceGroupAsConfirmed(ctx, pendingSequenceGroup.Tx.Hash(), nil)
			if err != nil {
				log.Errorf("failed to set sequence group as confirmed for tx %v: %v", pendingSequenceGroup.Tx.Hash().String(), err)
			}
			return
		}

		// if it was not mined yet, check it against the rules to improve the tx in order to get it mined
		//
		// check if the timeout since the last time the group was update has expired, if so, update the tx
		lastTimeSequenceWasUpdated := pendingSequenceGroup.CreatedAt
		if pendingSequenceGroup.UpdatedAt != nil {
			lastTimeSequenceWasUpdated = *pendingSequenceGroup.UpdatedAt
		}
		if time.Since(lastTimeSequenceWasUpdated) >= c.cfg.IntervalToReviewSendBatchTx.Duration {
			// TODO improve the current tx to make it mineable
		}
	}
}

// VerifyBatch TODO REMOVE
func (c *Client) VerifyBatch(batchNum uint64, proof *pb.GetProofResponse) {

}
