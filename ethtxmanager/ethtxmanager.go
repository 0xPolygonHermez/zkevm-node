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

	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/core/types"
)

// ErrMaxRetriesExceeded max retries exceeded error.
var ErrMaxRetriesExceeded = errors.New("Maximum number of retries exceeded")

// Client for eth tx manager
type Client struct {
	cfg    Config
	ethMan etherman
	state  state
}

// New creates new eth tx manager
func New(cfg Config, ethMan etherman, state state) *Client {
	return &Client{
		cfg:    cfg,
		ethMan: ethMan,
		state:  state,
	}
}

// SequenceBatches send sequences to the channel
func (c *Client) SequenceBatches(ctx context.Context, sequences []ethmanTypes.Sequence) error {
	var (
		attempts uint32
		gas      uint64
		gasPrice *big.Int
		nonce    = big.NewInt(0)
	)
	log.Info("sending sequence to L1")
	for attempts < c.cfg.MaxSendBatchTxRetries {
		var (
			tx  *types.Transaction
			err error
		)
		if nonce.Uint64() > 0 {
			tx, err = c.ethMan.SequenceBatches(ctx, sequences, gas, gasPrice, nonce)
		} else {
			tx, err = c.ethMan.SequenceBatches(ctx, sequences, gas, gasPrice, nil)
		}
		for err != nil && attempts < c.cfg.MaxSendBatchTxRetries {
			log.Errorf("failed to sequence batches, trying once again, retry #%d, err: %w", attempts, 0, err)
			time.Sleep(c.cfg.FrequencyForResendingFailedSendBatches.Duration)
			attempts++
			if nonce.Uint64() > 0 {
				tx, err = c.ethMan.SequenceBatches(ctx, sequences, gas, gasPrice, nonce)
			} else {
				tx, err = c.ethMan.SequenceBatches(ctx, sequences, gas, gasPrice, nil)
			}
		}
		if err != nil {
			log.Errorf("failed to sequence batches, maximum attempts exceeded, err: %w", err)
			return fmt.Errorf("failed to sequence batches, maximum attempts exceeded, err: %w", err)
		}
		// Wait for tx to be mined
		log.Infof("waiting for tx to be mined. Tx hash: %s, nonce: %d, gasPrice: %d", tx.Hash(), tx.Nonce(), tx.GasPrice().Int64())
		err = c.ethMan.WaitTxToBeMined(ctx, tx, c.cfg.WaitTxToBeMined.Duration)
		if err != nil {
			attempts++
			if errors.Is(err, runtime.ErrOutOfGas) {
				gas = increaseGasLimit(tx.Gas(), c.cfg.PercentageToIncreaseGasLimit)
				log.Infof("out of gas with %d, retrying with %d", tx.Gas(), gas)
				continue
			} else if errors.Is(err, operations.ErrTimeoutReached) {
				nonce = new(big.Int).SetUint64(tx.Nonce())
				gasPrice = increaseGasPrice(tx.GasPrice(), c.cfg.PercentageToIncreaseGasPrice)
				log.Infof("tx %s reached timeout, retrying with gas price = %d", tx.Hash(), gasPrice)
				continue
			}
			log.Errorf("tx %s failed, err: %w", tx.Hash(), err)
			return fmt.Errorf("tx %s failed, err: %w", tx.Hash(), err)
		}

		log.Infof("sequence sent to L1 successfully. Tx hash: %s", tx.Hash())
		return c.state.WaitSequencingTxToBeSynced(ctx, tx, c.cfg.WaitTxToBeSynced.Duration)
	}
	return ErrMaxRetriesExceeded
}

// VerifyBatches sends the VerifyBatches request to Ethereum. It is also
// responsible for retrying up to MaxVerifyBatchTxRetries times, increasing the
// Gas price or Gas limit, depending on the error returned by Ethereum.
func (c *Client) VerifyBatches(ctx context.Context, lastVerifiedBatch uint64, finalBatchNum uint64, resGetProof *pb.FinalProof, inputs ethmanTypes.FinalProofInputs) (*types.Transaction, error) {
	var (
		attempts uint32
		gas      uint64
		gasPrice *big.Int
		nonce    = big.NewInt(0)
		tx       *types.Transaction
		err      error
	)

	log.Infof("sending verification to L1 for batches %d-%d", lastVerifiedBatch+1, finalBatchNum)

	for attempts < c.cfg.MaxVerifyBatchTxRetries {
		if nonce.Uint64() > 0 {
			tx, err = c.ethMan.VerifyBatches(ctx, lastVerifiedBatch, finalBatchNum, resGetProof, inputs, gas, gasPrice, nonce)
		} else {
			tx, err = c.ethMan.VerifyBatches(ctx, lastVerifiedBatch, finalBatchNum, resGetProof, inputs, gas, gasPrice, nil)
		}
		for err != nil && attempts < c.cfg.MaxVerifyBatchTxRetries {
			log.Errorf("failed to send batch verification, trying once again, retry #%d, err: %w", attempts, err)
			time.Sleep(c.cfg.FrequencyForResendingFailedVerifyBatch.Duration)

			if nonce.Uint64() > 0 {
				tx, err = c.ethMan.VerifyBatches(ctx, lastVerifiedBatch, finalBatchNum, resGetProof, inputs, gas, gasPrice, nonce)
			} else {
				tx, err = c.ethMan.VerifyBatches(ctx, lastVerifiedBatch, finalBatchNum, resGetProof, inputs, gas, gasPrice, nil)
			}

			attempts++
		}
		if err != nil {
			log.Errorf("failed to send batch verification, maximum attempts exceeded, err: %w", err)
			return nil, fmt.Errorf("failed to send batch verification, maximum attempts exceeded, err: %w", err)
		}
		// Wait for tx to be mined
		log.Infof("waiting for tx to be mined. Tx hash: %s, nonce: %d, gasPrice: %d", tx.Hash(), tx.Nonce(), tx.GasPrice().Int64())
		err = c.ethMan.WaitTxToBeMined(ctx, tx, c.cfg.WaitTxToBeMined.Duration)
		if err != nil {
			if errors.Is(err, runtime.ErrOutOfGas) {
				gas = increaseGasLimit(tx.Gas(), c.cfg.PercentageToIncreaseGasLimit)
				log.Infof("out of gas with %d, retrying with %d", tx.Gas(), gas)
				continue
			} else if errors.Is(err, operations.ErrTimeoutReached) {
				nonce = new(big.Int).SetUint64(tx.Nonce())
				gasPrice = increaseGasPrice(tx.GasPrice(), c.cfg.PercentageToIncreaseGasPrice)
				log.Infof("tx %s reached timeout, retrying with gas price = %d", tx.Hash(), gasPrice)
				continue
			}
			log.Errorf("tx %s failed, err: %w", tx.Hash(), err)
			return nil, fmt.Errorf("tx %s failed, err: %w", tx.Hash(), err)
		}

		log.Infof("batch verification sent to L1 successfully. Tx hash: %s", tx.Hash())
		return tx, c.state.WaitVerifiedBatchToBeSynced(ctx, finalBatchNum, c.cfg.WaitTxToBeSynced.Duration)
	}
	return nil, ErrMaxRetriesExceeded
}

func increaseGasPrice(currentGasPrice *big.Int, percentageIncrease uint64) *big.Int {
	gasPrice := big.NewInt(0).Mul(currentGasPrice, new(big.Int).SetUint64(uint64(100)+percentageIncrease)) //nolint:gomnd
	return gasPrice.Div(gasPrice, big.NewInt(100))                                                         //nolint:gomnd
}

func increaseGasLimit(currentGasLimit uint64, percentageIncrease uint64) uint64 {
	return currentGasLimit * (100 + percentageIncrease) / 100 //nolint:gomnd
}
