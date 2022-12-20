// Package ethtxmanager handles ethereum transactions:  It makes
// calls to send and to aggregate batch, checks possible errors, like wrong nonce or gas limit too low
// and make correct adjustments to request according to it. Also, it tracks transaction receipt and status
// of tx in case tx is rejected and send signals to sequencer/aggregator to resend sequence/batch
package ethtxmanager

import (
	"context"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

// Client for eth tx manager
type Client struct {
	cfg     Config
	ethMan  etherman
	state   state
	storage *storage
}

// New creates new eth tx manager
func New(cfg Config, ethMan etherman, state state) *Client {
	c := &Client{
		cfg:     cfg,
		ethMan:  ethMan,
		state:   state,
		storage: newStorage(),
	}

	go c.manageTxs()

	return c
}

// SequenceBatches send sequences to the channel
func (c *Client) SequenceBatches(ctx context.Context, sequences []ethmanTypes.Sequence) error {
	log.Info("creating L1 tx to sequence batches")
	tx, err := c.storage.enqueueSequences(ctx, c.state, c.ethMan, c.cfg, sequences)
	if err != nil {
		log.Errorf("failed to create L1 tx to sequence batches, err: %w", err)
		return nil
	}
	log.Infof("L1 tx to sequence batches added to channel, hash: %v", tx.Hash().String())
	return nil
}

// VerifyBatches sends the VerifyBatches request to Ethereum. It is also
// responsible for retrying up to MaxVerifyBatchTxRetries times, increasing the
// Gas price or Gas limit, depending on the error returned by Ethereum.
func (c *Client) VerifyBatches(ctx context.Context, lastVerifiedBatch uint64, finalBatchNum uint64, inputs *ethmanTypes.FinalProofInputs) error {
	log.Info("creating L1 tx to verify batches")
	tx, err := c.storage.enqueueVerifyBatches(ctx, c.state, c.ethMan, c.cfg, lastVerifiedBatch, finalBatchNum, inputs)
	if err != nil {
		log.Errorf("failed to create L1 tx to verify batches, err: %w", err)
		return nil
	}
	log.Infof("L1 tx to verify batches added to channel, hash: %v", tx.Hash().String())
	return nil

	// log.Infof("Final proof for batches [%d-%d] verified in transaction [%v]", proof.BatchNumber, proof.BatchNumberFinal, tx.Hash())
}

// manageTxs will read txs from storage, send then to L1
// and keep monitoring tham until they get mined, so the
// next one can be sent
func (c *Client) manageTxs() {
	// infinite loop to manage txs as they arrive
	for {
		// gets the next tx to send to L1
		etx := c.storage.Next()
		var lastSentTxHash string
		// monitor and retries tx until it gets mined
		for {
			ctx := context.Background()
			err := etx.RenewTxIfNeeded(ctx, c.ethMan)
			if err != nil {
				log.Errorf("failed to renew tx if needed: %w", err)
				etx.Wait()
				continue
			}

			// sends the tx if it was renewed
			if etx.Tx().Hash().String() != lastSentTxHash {
				err := c.ethMan.SendTx(ctx, etx.Tx())
				if err != nil {
					log.Errorf("failed to send tx: %w", err)
					etx.Wait()
					continue
				}
				lastSentTxHash = etx.Tx().Hash().String()
			}

			// waits the txs to be mined
			err = c.ethMan.WaitTxToBeMined(ctx, etx.Tx(), c.cfg.WaitTxToBeMined.Duration)
			if err != nil {
				log.Errorf("failed to wait tx to be mined: %w", err)
				etx.Wait()
				continue
			}

			log.Infof("L1 tx mined successfully: %v", etx.Tx().Hash().String())

			// waits the synchronizer to sync the data from the tx that was mined
			err = etx.WaitSync(ctx)
			if err != nil {
				log.Errorf("failed to wait sync: %w", err)
				etx.Wait()
				continue
			}

			log.Infof("L1 tx synced successfully: %v", etx.Tx().Hash().String())

			break
		}
	}
}
