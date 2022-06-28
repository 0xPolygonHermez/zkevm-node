// Package ethtxmanager handles ethereum transactions:  It makes
// calls to send and to aggregate batch, checks possible errors, like wrong nonce or gas limit too low
// and make correct adjustments to request according to it. Also, it tracks transaction receipt and status
// of tx in case tx is rejected and send signals to sequencer/aggregator to resend sequence/batch
package ethtxmanager

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethmanTypes "github.com/hermeznetwork/hermez-core/ethermanv2/types"
	"github.com/hermeznetwork/hermez-core/log"
)

const (
	gasLimitIncrease  = 1.2
	sentEthTxsChanLen = 100
)

// Client for eth tx manager
type Client struct {
	cfg Config

	ethMan                 etherman
	sequenceBatchesTxsChan chan sequenceBatchesTx
}

type sequenceBatchesTx struct {
	sequences []ethmanTypes.Sequence
	hash      common.Hash
	gasLimit  uint64
}

// New creates new eth tx manager
func New(cfg Config, ethMan etherman) *Client {
	sequenceBatchesTxsChan := make(chan sequenceBatchesTx, sentEthTxsChanLen)
	return &Client{cfg: cfg, sequenceBatchesTxsChan: sequenceBatchesTxsChan, ethMan: ethMan}
}

// SequenceBatches send request to ethereum
func (c *Client) SequenceBatches(sequences []ethmanTypes.Sequence) error {
	gas, err := c.ethMan.EstimateGasSequenceBatches(sequences)
	if err != nil {
		return fmt.Errorf("failed to estimate gas for sending sequences batches, err: %v", err)
	}

	gasLimit := uint64(float64(gas) * gasLimitIncrease)
	tx, err := c.ethMan.SequenceBatches(sequences, gasLimit)
	if err != nil {
		return err
	}
	c.sequenceBatchesTxsChan <- sequenceBatchesTx{
		sequences: sequences,
		hash:      tx.Hash(),
		gasLimit:  gasLimit,
	}
	return nil
}

// TrackEthSentTransactions tracks sent txs to the ethereum
func (c *Client) TrackEthSentTransactions(ctx context.Context) {
	for {
		select {
		case tx := <-c.sequenceBatchesTxsChan:
			c.resendTxIfNeeded(tx)
		case <-ctx.Done():
			return
		}
	}
}

func (c *Client) resendTxIfNeeded(tx sequenceBatchesTx) {
	var (
		gasLimit       uint64
		counter        uint32
		isTxSuccessful bool
		err            error
	)
	hash := tx.hash
	ctx := context.Background()
	for !isTxSuccessful && counter <= c.cfg.MaxSendTxRetries {
		time.Sleep(time.Duration(c.cfg.FrequencyForResendingFailedTxs) * time.Millisecond)
		receipt := c.getTxReceipt(ctx, hash)
		if receipt == nil {
			continue
		}
		// tx is failed, so batch should be sent again
		if receipt.Status == 0 {
			gasLimit, hash, err = c.resendSequenceBatches(gasLimit, tx, hash, counter)
			if err != nil {
				log.Errorf("failed to resend sequence batches to the ethereum, err: %v", err)
			}
			counter++
			continue
		}

		log.Infof("sendBatch transaction %s is successful", hash.Hex())
		isTxSuccessful = true
	}
	if counter == c.cfg.MaxSendTxRetries {
		log.Fatalf("failed to send txs %v several times,"+
			" gas limit %d is too high, first tx hash %s, last tx hash %s",
			tx.sequences, gasLimit, tx.hash.Hex(), hash.Hex())
	}
}

func (c *Client) resendSequenceBatches(gasLimit uint64, tx sequenceBatchesTx, hash common.Hash, counter uint32) (uint64, common.Hash, error) {
	log.Warnf("increasing gas limit for the transaction sending, previous failed tx hash %v", hash)

	gasLimit = uint64(float64(gasLimit) * gasLimitIncrease)
	sentTx, err := c.ethMan.SequenceBatches(tx.sequences, gasLimit)
	if err != nil {
		log.Warnf("failed to send batch once again, err: %v", err)
		return gasLimit, hash, err
	}
	hash = sentTx.Hash()
	log.Infof("sent sendBatch transaction with hash %s and gas limit %d with try number %d",
		hash, gasLimit, counter)

	return gasLimit, hash, nil
}

func (c *Client) getTxReceipt(ctx context.Context, hash common.Hash) *types.Receipt {
	_, isPending, err := c.ethMan.GetTx(ctx, hash)
	if err != nil {
		log.Warnf("failed to get tx with hash %s, err %v", hash, err)
		return nil
	}
	if isPending {
		log.Debugf("sendBatch transaction %s is pending", hash)
		return nil
	}

	receipt, err := c.ethMan.GetTxReceipt(ctx, hash)
	if err != nil {
		log.Warnf("failed to get tx receipt with hash %v, err %v", hash.Hex(), err)
		return nil
	}
	return receipt
}
