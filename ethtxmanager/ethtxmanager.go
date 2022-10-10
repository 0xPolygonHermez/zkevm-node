// Package ethtxmanager handles ethereum transactions:  It makes
// calls to send and to aggregate batch, checks possible errors, like wrong nonce or gas limit too low
// and make correct adjustments to request according to it. Also, it tracks transaction receipt and status
// of tx in case tx is rejected and send signals to sequencer/aggregator to resend sequence/batch
package ethtxmanager

import (
	"math/big"
	"strings"
	"time"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
)

const oneHundred = 100

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
func (c *Client) SequenceBatches(sequences []ethmanTypes.Sequence, gasLimit uint64) {
	var (
		attempts uint32
		gasPrice *big.Int
	)
	log.Info("sending sequence to L1")
	for attempts < c.cfg.MaxSendBatchTxRetries {
		tx, err := c.ethMan.SequenceBatches(sequences, gasLimit, gasPrice)
		for err != nil && attempts < c.cfg.MaxSendBatchTxRetries {
			log.Errorf("failed to sequence batches, trying once again, retry #%d, gasLimit: %d, err: %w",
				attempts, 0, err)
			time.Sleep(c.cfg.FrequencyForResendingFailedSendBatches.Duration)
			attempts++
			tx, err = c.ethMan.SequenceBatches(sequences, gasLimit, gasPrice)
		}
		if err != nil {
			log.Fatalf("failed to sequence batches, maximum attempts exceeded, gasLimit: %d, err: %w",
				0, err)
		}
		// Wait for tx to be mined
		log.Infof("waiting for sequence to be mined. Tx hash: %s", tx.Hash())
		err = c.ethMan.WaitTxToBeMined(tx.Hash(), c.cfg.WaitTxToBeMined.Duration)
		if err != nil {
			attempts++
			if strings.Contains(err.Error(), "out of gas") {
				gasLimit = (tx.Gas() * (oneHundred + c.cfg.PercentageToIncreaseGasLimit)) / oneHundred
				log.Infof("out of gas with %d, retrying with %d", tx.Gas(), gasLimit)
				continue
			} else if strings.Contains(err.Error(), "timeout has been reached") {
				gasPrice.Mul(tx.GasPrice(), new(big.Int).SetUint64(uint64(oneHundred)+c.cfg.PercentageToIncreaseGasPrice))
				gasPrice.Div(gasPrice, big.NewInt(oneHundred))
				log.Infof("tx %s reached timeout, retrying with gas price = %d", tx.Hash(), gasPrice)
				continue
			}
			log.Fatalf("tx %s failed, err: %w", tx.Hash(), err)
		} else {
			log.Infof("sequence sent to L1 successfully. Tx hash: %s", tx.Hash())
			return
		}
	}
}

// VerifyBatch send VerifyBatch request to ethereum
func (c *Client) VerifyBatch(batchNum uint64, resGetProof *pb.GetProofResponse) {
	var (
		attempts uint32
		gas      uint64
		gasPrice *big.Int
	)
	log.Infof("sending batch %d verification to L1", batchNum)
	for attempts < c.cfg.MaxVerifyBatchTxRetries {
		tx, err := c.ethMan.VerifyBatch(batchNum, resGetProof, gas, gasPrice)
		for err != nil && attempts < c.cfg.MaxSendBatchTxRetries {
			log.Errorf("failed to send batch verification, trying once again, retry #%d, gasLimit: %d, err: %w",
				attempts, 0, err)
			time.Sleep(c.cfg.FrequencyForResendingFailedSendBatches.Duration)
			attempts++
			tx, err = c.ethMan.VerifyBatch(batchNum, resGetProof, gas, gasPrice)
		}
		if err != nil {
			log.Fatalf("failed to send batch verification, maximum attempts exceeded, gasLimit: %d, err: %w",
				0, err)
		}
		// Wait for tx to be mined
		log.Infof("waiting for tx to be mined. Tx hash: %s", tx.Hash())
		err = c.ethMan.WaitTxToBeMined(tx.Hash(), c.cfg.WaitTxToBeMined.Duration)
		if err != nil {
			attempts++
			if strings.Contains(err.Error(), "out of gas") {
				gas = (tx.Gas() * (oneHundred + c.cfg.PercentageToIncreaseGasLimit)) / oneHundred
				log.Infof("out of gas with %d, retrying with %d", tx.Gas(), gas)
				continue
			} else if strings.Contains(err.Error(), "timeout has been reached") {
				gasPrice.Mul(tx.GasPrice(), new(big.Int).SetUint64(uint64(oneHundred)+c.cfg.PercentageToIncreaseGasPrice))
				gasPrice.Div(gasPrice, big.NewInt(oneHundred))
				log.Infof("tx %s reached timeout, retrying with gas price = %d", tx.Hash(), gasPrice)
				continue
			}
			log.Fatalf("tx %s failed, err: %w", tx.Hash(), err)
		} else {
			log.Infof("batch verification sent to L1 successfully. Tx hash: %s", tx.Hash())
			return
		}
	}
}
