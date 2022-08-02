// Package ethtxmanager handles ethereum transactions:  It makes
// calls to send and to aggregate batch, checks possible errors, like wrong nonce or gas limit too low
// and make correct adjustments to request according to it. Also, it tracks transaction receipt and status
// of tx in case tx is rejected and send signals to sequencer/aggregator to resend sequence/batch
package ethtxmanager

import (
	"context"
	"fmt"
	"time"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	gasLimitIncrease  = 1.2
	sentEthTxsChanLen = 100
)

// Client for eth tx manager
type Client struct {
	cfg Config

	ethMan              etherman
	sequencesToSendChan chan []ethmanTypes.Sequence
	verifyBatchTxsChan  chan verifyBatchTx
}

type verifyBatchTx struct {
	batchNumber uint64
	resGetProof *pb.GetProofResponse
	hash        common.Hash
	gasLimit    uint64
}

// New creates new eth tx manager
func New(cfg Config, ethMan etherman) *Client {
	verifyBatchTxsChan := make(chan verifyBatchTx, sentEthTxsChanLen)
	sequencesToSendChan := make(chan []ethmanTypes.Sequence, sentEthTxsChanLen)
	return &Client{
		cfg:                 cfg,
		sequencesToSendChan: sequencesToSendChan,
		verifyBatchTxsChan:  verifyBatchTxsChan,
		ethMan:              ethMan,
	}
}

// TrackSequenceBatchesSending tracks and send sequences, that should be sent
func (c *Client) TrackSequenceBatchesSending(ctx context.Context) {
	for {
		select {
		case sequences := <-c.sequencesToSendChan:
			var attempts uint32
			log.Info("sending sequence to L1")
			tx, err := c.ethMan.SequenceBatches(sequences, 0)
			for err != nil && attempts < c.cfg.MaxSendBatchTxRetries {
				log.Errorf("failed to sequence batches, trying once again, retry #%d, gasLimit: %d, err: %v",
					attempts, 0, err)
				time.Sleep(c.cfg.FrequencyForResendingFailedSendBatches.Duration)
				attempts++
				tx, err = c.ethMan.SequenceBatches(sequences, 0)
			}
			if err != nil {
				log.Fatalf("failed to sequence batches, maximum attempts exceeded, gasLimit: %d, err: %v",
					0, err)
			}
			// Wait for tx to be mined
			log.Infof("waiting for sequence to be mined. Tx hash: %s", tx.Hash())
			err = c.ethMan.WaitTxToBeMined(tx.Hash(), time.Minute*2) //nolint:gomnd
			if err != nil {
				log.Fatalf("tx %s failed, err: %v", tx.Hash(), err)
			}
			log.Infof("sequence sent to L1 successfully. Tx hash: %s", tx.Hash())
			// Check if success
		case <-ctx.Done():
			return
		}
	}
}

// SequenceBatches send sequences to the channel
func (c *Client) SequenceBatches(sequences []ethmanTypes.Sequence) {
	c.sequencesToSendChan <- sequences
}

// VerifyBatch send VerifyBatch request to ethereum
func (c *Client) VerifyBatch(batchNum uint64, resGetProof *pb.GetProofResponse) error {
	gas, err := c.ethMan.EstimateGasForVerifyBatch(batchNum, resGetProof)
	if err != nil {
		return fmt.Errorf("failed to estimate gas for sending sequences batches, err: %v", err)
	}

	gasLimit := uint64(float64(gas) * gasLimitIncrease)
	tx, err := c.ethMan.VerifyBatch(batchNum, resGetProof, gasLimit)
	if err != nil {
		return err
	}
	c.verifyBatchTxsChan <- verifyBatchTx{
		batchNumber: batchNum,
		resGetProof: resGetProof,
		hash:        tx.Hash(),
		gasLimit:    gasLimit,
	}
	return nil
}

// TrackEthSentTransactions tracks sent txs to the ethereum
func (c *Client) TrackEthSentTransactions(ctx context.Context) {
	for {
		select {
		case tx := <-c.verifyBatchTxsChan:
			c.resendVerifyBatchTxIfNeeded(ctx, tx)
		case <-ctx.Done():
			return
		}
	}
}

func (c *Client) resendVerifyBatchTxIfNeeded(ctx context.Context, tx verifyBatchTx) {
	var (
		gasLimit       uint64
		counter        uint32
		isTxSuccessful bool
		err            error
	)
	hash := tx.hash
	for !isTxSuccessful && counter <= c.cfg.MaxVerifyBatchTxRetries {
		time.Sleep(c.cfg.FrequencyForResendingFailedVerifyBatch.Duration)
		receipt := c.getTxReceipt(ctx, hash)
		if receipt == nil {
			continue
		}
		// tx is failed, so batch should be sent again
		if receipt.Status == 0 {
			gasLimit, hash, err = c.resendVerifyBatch(gasLimit, tx, hash, counter)
			if err != nil {
				log.Errorf("failed to resend verify batch to the ethereum, err: %v", err)
			}
			counter++
			continue
		}

		log.Infof("verifyBatch transaction %s is successful", hash.Hex())
		isTxSuccessful = true
	}
	if counter == c.cfg.MaxSendBatchTxRetries {
		log.Fatalf("failed to send verify batch several times,"+
			"batchNumber %d, gas limit %d is too high, first tx hash %s, last tx hash %s",
			tx.batchNumber, gasLimit, tx.hash.Hex(), hash.Hex())
	}
}

func (c *Client) resendVerifyBatch(gasLimit uint64, tx verifyBatchTx, hash common.Hash, counter uint32) (uint64, common.Hash, error) {
	log.Warnf("increasing gas limit for the transaction sending, previous failed tx hash %v", hash)

	gasLimit = uint64(float64(gasLimit) * gasLimitIncrease)
	sentTx, err := c.ethMan.VerifyBatch(tx.batchNumber, tx.resGetProof, gasLimit)
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
