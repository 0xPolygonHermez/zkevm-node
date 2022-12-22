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

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const failureIntervalInSeconds = 5

var (
	// ErrNotFound when the object is not found
	ErrNotFound = errors.New("Not Found")
	// ErrAlreadyExists when the object already exists
	ErrAlreadyExists = errors.New("Already Exists")
)

// Client for eth tx manager
type Client struct {
	cfg      Config
	etherman ethermanInterface
	storage  storageInterface
}

// New creates new eth tx manager
func New(cfg Config, ethMan ethermanInterface, storage storageInterface) *Client {
	c := &Client{
		cfg:      cfg,
		etherman: ethMan,
		storage:  storage,
	}

	go c.manageTxs()

	return c
}

// SequenceBatches deprecated
func (c *Client) SequenceBatches(ctx context.Context, sequences []ethmanTypes.Sequence) error {
	panic("deprecated")
}

// VerifyBatches deprecated
func (c *Client) VerifyBatches(ctx context.Context, lastVerifiedBatch uint64, batchNum uint64, inputs *ethmanTypes.FinalProofInputs) error {
	panic("deprecated")
}

// Add a transaction to be sent and monitored
func (c *Client) Add(ctx context.Context, id string, from common.Address, to *common.Address, value *big.Int, data []byte) error {
	// get next nonce
	nonce, err := c.etherman.CurrentNonce(ctx)
	if err != nil {
		err := fmt.Errorf("failed to get current nonce: %w", err)
		log.Errorf(err.Error())
		return err
	}
	// get gas
	gas, err := c.etherman.EstimateGas(ctx, from, to, nil, data)
	if err != nil {
		err := fmt.Errorf("failed to estimate gas: %w", err)
		log.Errorf(err.Error())
		return err
	}
	// get gas price
	gasPrice, err := c.etherman.SuggestedGasPrice(ctx)
	if err != nil {
		err := fmt.Errorf("failed to get suggested gas price: %w", err)
		log.Errorf(err.Error())
		return err
	}

	// create monitored tx
	mTx := monitoredTx{
		id: id, from: from, to: to,
		nonce: nonce, value: value, data: data,
		gas: gas, gasPrice: gasPrice,
		status: MonitoredTxStatusCreated,
	}

	// add to storage
	err = c.storage.Add(ctx, mTx)
	if err != nil {
		err := fmt.Errorf("failed to add tx to get monitored: %w", err)
		log.Errorf(err.Error())
		return err
	}

	return nil
}

// Status returns the current status of the transaction
func (c *Client) Status(ctx context.Context, id string) (MonitoredTxStatus, error) {
	mTx, err := c.storage.Get(ctx, id)
	if err != nil {
		return MonitoredTxStatus(""), err
	}

	return mTx.status, nil
}

// manageTxs will read txs from storage, send then to L1
// and keep monitoring tham until they get mined, so the
// next one can be sent
func (c *Client) manageTxs() {
	// infinite loop to manage txs as they arrive
	for {
		ctx := context.Background()
		err := c.processMonitoredTxs(ctx)
		if err != nil {
			c.logErrorAndWait("failed to process created monitored txs: %v", err)
		}
	}
}

// processMonitoredTxs process all monitored tx with Created and Sent status
func (c *Client) processMonitoredTxs(ctx context.Context) error {
	mTxs, err := c.storage.GetByStatus(ctx, MonitoredTxStatusCreated, MonitoredTxStatusSent)
	if err != nil {
		return fmt.Errorf("failed to get created monitored txs: %v", err)
	}

	log.Debugf("found %v monitored tx to process", len(mTxs))

	for _, mTx := range mTxs {
		mTxLog := log.WithFields("monitored tx", mTx.id)
		mTxLog.Debug("processing")

		// check if any of the txs in the history was mined
		mined := false
		var receipt *types.Receipt
		for txHash := range mTx.history {
			mined, receipt, err = c.etherman.CheckTxWasMined(ctx, txHash)
			if err != nil {
				mTxLog.Errorf("failed to check if tx %v was mined: %v", txHash.String(), err)
				continue
			}
		}

		if !mined {
			// review tx and increase gas and gas price if needed
			if mTx.status == MonitoredTxStatusSent {
				err := c.ReviewMonitoredTx(mTx)
				if err != nil {
					mTxLog.Errorf("failed to review monitored tx: %v", err)
					continue
				}
			}

			// rebuild transaction
			tx := mTx.Tx()
			mTxLog.Debugf("unsigned tx %v created", tx.Hash().String(), mTx.id)

			// sign tx
			signedTx, err := c.etherman.SignTx(ctx, tx)
			if err != nil {
				mTxLog.Errorf("failed to sign tx %v created from monitored tx %v: %v", tx.Hash().String(), mTx.id, err)
				continue
			}
			mTxLog.Debugf("signed tx %v created", signedTx.Hash().String())

			// add tx to monitored tx history
			err = mTx.AddHistory(signedTx)
			if errors.Is(err, ErrAlreadyExists) {
				mTxLog.Debugf("signed tx already existed in the history")
			} else if err != nil {
				mTxLog.Errorf("failed to add signed tx to monitored tx %v history: %v", mTx.id, err)
				continue
			} else {
				mTxLog.Debugf("signed tx added to the history")
			}

			// check if the tx is already in the network, if not, send it
			_, _, err = c.etherman.GetTx(ctx, signedTx.Hash())
			// if not found, send it tx to the network
			if errors.Is(err, ethereum.NotFound) {
				mTxLog.Debugf("signed tx not found in the network")
				err := c.etherman.SendTx(ctx, signedTx)
				if err != nil {
					mTxLog.Errorf("failed to send tx %v to network: %v", signedTx.Hash().String(), err)
					continue
				}
				mTxLog.Debugf("signed tx sent to the network")
				if mTx.status == MonitoredTxStatusCreated {
					// update tx status to sent
					mTx.status = MonitoredTxStatusSent
					mTxLog.Debugf("status changed to %v", string(mTx.status))
					// update monitored tx changes into storage
					err = c.storage.Update(ctx, mTx)
					if err != nil {
						mTxLog.Errorf("failed to update monitored tx changes: %v", err)
						continue
					}
					mTxLog.Debugf("storage updated")
				}
			} else {
				mTxLog.Debugf("signed tx already found in the network")
			}

			// wait tx to get mined
			err = c.etherman.WaitTxToBeMined(ctx, signedTx, c.cfg.WaitTxToBeMined.Duration)
			if err != nil {
				mTxLog.Errorf("failed to wait tx to be mined: %v", err)
				continue
			}

			// get tx receipt
			receipt, err = c.etherman.GetTxReceipt(ctx, signedTx.Hash())
			if err != nil {
				mTxLog.Errorf("failed to get tx receipt for tx %v: %v", signedTx.Hash().String(), err)
				continue
			}
		}

		// if mined, check receipt and mark as Failed or Confirmed
		if receipt.Status == types.ReceiptStatusSuccessful {
			mTx.status = MonitoredTxStatusConfirmed
		} else {
			mTx.status = MonitoredTxStatusFailed
		}

		// update monitored tx changes into storage
		err = c.storage.Update(ctx, mTx)
		if err != nil {
			mTxLog.Errorf("failed to update monitored tx changes: %v", err)
			continue
		}
		mTxLog.Debugf("storage updated")
	}

	return nil
}

// ReviewMonitoredTx checks if some field needs to be updated
// accordingly to the current information stored and the current
// state of the network
func (c *Client) ReviewMonitoredTx(mTx monitoredTx) error {
	panic("not implemented yet")
}

// logErrorAndWait used when an error is detected before trying again
func (c *Client) logErrorAndWait(msg string, err error) {
	log.Errorf(msg, err)
	time.Sleep(failureIntervalInSeconds * time.Second)
}
