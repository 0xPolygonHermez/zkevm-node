// Package ethtxmanager handles ethereum transactions:  It makes
// calls to send and to aggregate batch, checks possible errors, like wrong nonce or gas limit too low
// and make correct adjustments to request according to it. Also, it tracks transaction receipt and status
// of tx in case tx is rejected and send signals to sequencer/aggregator to resend sequence/batch
package ethtxmanager

import (
	"strings"
	"time"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
)

// Client for eth tx manager
type Client struct {
	cfg    Config
	ethMan etherman
}

// New creates new eth tx manager
func New(cfg Config, ethMan etherman) *Client {
	return &Client{
		cfg:    cfg,
		ethMan: ethMan,
	}
}

// SequenceBatches send sequences to the channel
func (c *Client) SequenceBatches(sequences []ethmanTypes.Sequence) {
	var attempts uint32
	var gas uint64
	log.Info("sending sequence to L1")
	for attempts < c.cfg.MaxSendBatchTxRetries {
		tx, err := c.ethMan.SequenceBatches(sequences, gas)
		for err != nil && attempts < c.cfg.MaxSendBatchTxRetries {
			log.Errorf("failed to sequence batches, trying once again, retry #%d, gasLimit: %d, err: %v",
				attempts, 0, err)
			time.Sleep(c.cfg.FrequencyForResendingFailedSendBatches.Duration)
			attempts++
			tx, err = c.ethMan.SequenceBatches(sequences, gas)
		}
		if err != nil {
			log.Fatalf("failed to sequence batches, maximum attempts exceeded, gasLimit: %d, err: %v",
				0, err)
		}
		// Wait for tx to be mined
		log.Infof("waiting for sequence to be mined. Tx hash: %s", tx.Hash())
		// TODO: timeout via config file
		err = c.ethMan.WaitTxToBeMined(tx.Hash(), time.Minute*2) //nolint:gomnd
		if err != nil {
			attempts++
			if strings.Contains(err.Error(), "out of gas") {
				// TODO: percentage gas inncrease via config file
				gas = uint64(float64(tx.Gas()) * 1.1) //nolint:gomnd
				log.Infof("out of gas with %d, retrying with %d", tx.Gas(), gas)
				continue
			}
			// TODO: handle timeout by increasing gas price
			log.Fatalf("tx %s failed, err: %v", tx.Hash(), err)
		} else {
			log.Infof("sequence sent to L1 successfully. Tx hash: %s", tx.Hash())
			return
		}
	}
}

// VerifyBatch send VerifyBatch request to ethereum
func (c *Client) VerifyBatch(batchNum uint64, resGetProof *pb.GetProofResponse) {
	var attempts uint32
	// TODO: Remove this limit after Testnet testing
	var gas uint64 = 500000
	log.Infof("sending batch %d verification to L1", batchNum)
	for attempts < c.cfg.MaxVerifyBatchTxRetries {
		tx, err := c.ethMan.VerifyBatch(batchNum, resGetProof, gas)
		for err != nil && attempts < c.cfg.MaxSendBatchTxRetries {
			log.Errorf("failed to send batch verification, trying once again, retry #%d, gasLimit: %d, err: %v",
				attempts, 0, err)
			time.Sleep(c.cfg.FrequencyForResendingFailedSendBatches.Duration)
			attempts++
			tx, err = c.ethMan.VerifyBatch(batchNum, resGetProof, gas)
		}
		if err != nil {
			log.Fatalf("failed to send batch verification, maximum attempts exceeded, gasLimit: %d, err: %v",
				0, err)
		}
		// Wait for tx to be mined
		log.Infof("waiting for tx to be mined. Tx hash: %s", tx.Hash())
		// TODO: timeout via config file
		err = c.ethMan.WaitTxToBeMined(tx.Hash(), time.Minute*2) //nolint:gomnd
		if err != nil {
			attempts++
			if strings.Contains(err.Error(), "out of gas") {
				// TODO: percentage gas inncrease via config file
				gas = uint64(float64(tx.Gas()) * 1.1) //nolint:gomnd
				log.Infof("out of gas with %d, retrying with %d", tx.Gas(), gas)
				continue
			}
			// TODO: handle timeout by increasing gas price
			log.Fatalf("tx %s failed, err: %v", tx.Hash(), err)
		} else {
			log.Infof("batch verification sent to L1 successfully. Tx hash: %s", tx.Hash())
			return
		}
	}
}
