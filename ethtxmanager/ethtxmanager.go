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
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

const failureIntervalInSeconds = 5

var (
	// ErrNotFound when the object is not found
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists when the object already exists
	ErrAlreadyExists = errors.New("already exists")
)

// Client for eth tx manager
type Client struct {
	ctx    context.Context
	cancel context.CancelFunc

	cfg      Config
	etherman ethermanInterface
	storage  storageInterface
	state    stateInterface
}

// New creates new eth tx manager
func New(cfg Config, ethMan ethermanInterface, storage storageInterface, state stateInterface) *Client {
	c := &Client{
		cfg:      cfg,
		etherman: ethMan,
		storage:  storage,
		state:    state,
	}

	return c
}

// VerifyBatches deprecated
func (c *Client) VerifyBatches(ctx context.Context, lastVerifiedBatch uint64, batchNum uint64, inputs *ethmanTypes.FinalProofInputs) error {
	panic("deprecated")
}

// Add a transaction to be sent and monitored
func (c *Client) Add(ctx context.Context, owner, id string, from common.Address, to *common.Address, value *big.Int, data []byte, dbTx pgx.Tx) error {
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
		owner: owner, id: id, from: from, to: to,
		nonce: nonce, value: value, data: data,
		gas: gas, gasPrice: gasPrice,
		status: MonitoredTxStatusCreated,
	}

	// add to storage
	err = c.storage.Add(ctx, mTx, dbTx)
	if err != nil {
		err := fmt.Errorf("failed to add tx to get monitored: %w", err)
		log.Errorf(err.Error())
		return err
	}

	return nil
}

// ResultsByStatus returns all the results for all the monitored txs related to the owner and matching the provided statuses
// if the statuses are empty, all the statuses are considered.
//
// the slice is returned is in order by created_at field ascending
func (c *Client) ResultsByStatus(ctx context.Context, owner string, statuses []MonitoredTxStatus, dbTx pgx.Tx) ([]MonitoredTxResult, error) {
	mTxs, err := c.storage.GetByStatus(ctx, &owner, statuses, dbTx)
	if err != nil {
		return nil, err
	}

	results := make([]MonitoredTxResult, 0, len(mTxs))

	for _, mTx := range mTxs {
		result, err := c.buildResult(ctx, mTx)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

// Result returns the current result of the transaction execution with all the details
func (c *Client) Result(ctx context.Context, owner, id string, dbTx pgx.Tx) (MonitoredTxResult, error) {
	mTx, err := c.storage.Get(ctx, owner, id, dbTx)
	if err != nil {
		return MonitoredTxResult{}, err
	}

	return c.buildResult(ctx, mTx)
}

// SetStatusDone sets the status of a monitored tx to MonitoredStatusDone.
// this method is provided to the callers to decide when a monitored tx should be
// considered done, so they can start to ignore it when querying it by Status.
func (c *Client) SetStatusDone(ctx context.Context, owner, id string, dbTx pgx.Tx) error {
	mTx, err := c.storage.Get(ctx, owner, id, nil)
	if err != nil {
		return err
	}

	mTx.status = MonitoredTxStatusDone

	return c.storage.Update(ctx, mTx, dbTx)
}

func (c *Client) buildResult(ctx context.Context, mTx monitoredTx) (MonitoredTxResult, error) {
	history := mTx.historyHashSlice()
	txs := make(map[common.Hash]TxResult, len(history))

	for _, txHash := range history {
		tx, _, err := c.etherman.GetTx(ctx, txHash)
		if !errors.Is(err, ethereum.NotFound) && err != nil {
			return MonitoredTxResult{}, err
		}

		receipt, err := c.etherman.GetTxReceipt(ctx, txHash)
		if !errors.Is(err, ethereum.NotFound) && err != nil {
			return MonitoredTxResult{}, err
		}

		revertMessage, err := c.etherman.GetRevertMessage(ctx, tx)
		if err != nil {
			return MonitoredTxResult{}, err
		}

		txs[txHash] = TxResult{
			Tx:            tx,
			Receipt:       receipt,
			RevertMessage: revertMessage,
		}
	}

	result := MonitoredTxResult{
		ID:     mTx.id,
		Status: mTx.status,
		Txs:    txs,
	}

	return result, nil
}

// Start will start the tx management, reading txs from storage,
// send then to the blockchain and keep monitoring them until they
// get mined
func (c *Client) Start() {
	// infinite loop to manage txs as they arrive
	c.ctx, c.cancel = context.WithCancel(context.Background())

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-time.After(c.cfg.FrequencyToMonitorTxs.Duration):
			err := c.processMonitoredTxs(context.Background())
			if err != nil {
				c.logErrorAndWait("failed to process created monitored txs: %v", err)
			}
		}
	}
}

// Stop will stops the monitored tx management
func (c *Client) Stop() {
	c.cancel()
}

// ProcessReorg updates all monitored txs from provided block number until the last one to
// Reorged status, allowing it to be reprocessed by the tx monitoring
func (c *Client) ProcessReorg(ctx context.Context, fromBlockNumber uint64, dbTx pgx.Tx) error {
	mTxs, err := c.storage.GetByBlock(ctx, &fromBlockNumber, nil, dbTx)
	if err != nil {
		return err
	}
	for _, mTx := range mTxs {
		mTx.blockNumber = nil
		mTx.status = MonitoredTxStatusReorged

		err = c.storage.Update(ctx, mTx, dbTx)
		if err != nil {
			log.Errorf("failed to update monitored tx to reorg status: %v", err)
			return err
		}
	}
	return nil
}

// processMonitoredTxs process all monitored tx with Created and Sent status
func (c *Client) processMonitoredTxs(ctx context.Context) error {
	statusesFilter := []MonitoredTxStatus{MonitoredTxStatusCreated, MonitoredTxStatusSent, MonitoredTxStatusReorged}
	mTxs, err := c.storage.GetByStatus(ctx, nil, statusesFilter, nil)
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
			// if is a reorged, move to the next
			if mTx.status == MonitoredTxStatusReorged {
				continue
			}

			// review tx and increase gas and gas price if needed
			if mTx.status == MonitoredTxStatusSent {
				err := c.ReviewMonitoredTx(ctx, mTx)
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
				mTxLog.Debugf("signed tx added to the monitored tx history")
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
					err = c.storage.Update(ctx, mTx, nil)
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

		mTx.blockNumber = receipt.BlockNumber

		// if mined, check receipt and mark as Failed or Confirmed
		if receipt.Status == types.ReceiptStatusSuccessful {
			receiptBlockNum := receipt.BlockNumber.Uint64()

			// check block synced
			block, err := c.state.GetLastBlock(ctx, nil)
			if errors.Is(err, state.ErrStateNotSynchronized) {
				mTxLog.Debugf("state not synchronized yet, waiting for L1 block %v to be synced", receiptBlockNum)
				continue
			} else if err != nil {
				mTxLog.Errorf("failed to check if L1 block %v is already synced: %v", receiptBlockNum, err)
				continue
			} else if block.BlockNumber < receiptBlockNum {
				mTxLog.Debugf("L1 block %v not synchronized yet, waiting for L1 block to be synced in order to confirm monitored tx", receiptBlockNum)
				continue
			} else {
				mTx.status = MonitoredTxStatusConfirmed
			}
		} else {
			mTx.status = MonitoredTxStatusFailed
		}

		// update monitored tx changes into storage
		err = c.storage.Update(ctx, mTx, nil)
		if err != nil {
			mTxLog.Errorf("failed to update monitored tx: %v", err)
			continue
		}
		mTxLog.Debugf("storage updated")
	}

	return nil
}

// ReviewMonitoredTx checks if some field needs to be updated
// accordingly to the current information stored and the current
// state of the network
func (c *Client) ReviewMonitoredTx(ctx context.Context, mTx monitoredTx) error {
	// get gas
	gas, err := c.etherman.EstimateGas(ctx, mTx.from, mTx.to, mTx.value, mTx.data)
	if err != nil {
		err := fmt.Errorf("failed to estimate gas: %w", err)
		log.Errorf(err.Error())
		return err
	}

	// check gas
	if gas > mTx.gas {
		mTx.gas = gas
	}

	// get gas price
	gasPrice, err := c.etherman.SuggestedGasPrice(ctx)
	if err != nil {
		err := fmt.Errorf("failed to get suggested gas price: %w", err)
		log.Errorf(err.Error())
		return err
	}

	// check gas price
	if gasPrice.Cmp(mTx.gasPrice) == 1 {
		mTx.gasPrice = gasPrice
	}
	return nil
}

// logErrorAndWait used when an error is detected before trying again
func (c *Client) logErrorAndWait(msg string, err error) {
	log.Errorf(msg, err)
	time.Sleep(failureIntervalInSeconds * time.Second)
}
