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
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

const (
	failureIntervalInSeconds = 5
	// maxHistorySize           = 10
)

var (
	// ErrNotFound when the object is not found
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists when the object already exists
	ErrAlreadyExists = errors.New("already exists")

	// ErrExecutionReverted returned when trying to get the revert message
	// but the call fails without revealing the revert reason
	ErrExecutionReverted = errors.New("execution reverted")

	// gasOffsets for aggregator and sequencer
	gasOffsets = map[string]uint64{
		"sequencer":  80000, //nolint:gomnd
		"aggregator": 0,
	}
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

// Add a transaction to be sent and monitored
func (c *Client) Add(ctx context.Context, owner, id string, from common.Address, to *common.Address, value *big.Int, data []byte, dbTx pgx.Tx) error {
	// get next nonce
	nonce, err := c.etherman.CurrentNonce(ctx, from)
	if err != nil {
		err := fmt.Errorf("failed to get current nonce: %w", err)
		log.Errorf(err.Error())
		return err
	}
	// get gas
	gas, err := c.etherman.EstimateGas(ctx, from, to, value, data)
	if err != nil {
		err := fmt.Errorf("failed to estimate gas: %w, data: %v", err, common.Bytes2Hex(data))
		log.Error(err.Error())
		if c.cfg.ForcedGas > 0 {
			gas = c.cfg.ForcedGas
		} else {
			return err
		}
	} else {
		offset := gasOffsets[owner]
		gas += offset
		log.Debugf("Applying gasOffset: %d. Final Gas: %d, Owner: %s", offset, gas, owner)
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
func (c *Client) setStatusDone(ctx context.Context, owner, id string, dbTx pgx.Tx) error {
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
		if !errors.Is(err, ethereum.NotFound) && err != nil && err.Error() != ErrExecutionReverted.Error() {
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
			err := c.monitorTxs(context.Background())
			if err != nil {
				c.logErrorAndWait("failed to monitor txs: %v", err)
			}
		}
	}
}

// Stop will stops the monitored tx management
func (c *Client) Stop() {
	c.cancel()
}

// Reorg updates all monitored txs from provided block number until the last one to
// Reorged status, allowing it to be reprocessed by the tx monitoring
func (c *Client) Reorg(ctx context.Context, fromBlockNumber uint64, dbTx pgx.Tx) error {
	log.Infof("processing reorg from block: %v", fromBlockNumber)
	mTxs, err := c.storage.GetByBlock(ctx, &fromBlockNumber, nil, dbTx)
	if err != nil {
		log.Errorf("failed to monitored tx by block: %v", err)
		return err
	}
	log.Infof("updating %v monitored txs to reorged", len(mTxs))
	for _, mTx := range mTxs {
		mTx.blockNumber = nil
		mTx.status = MonitoredTxStatusReorged

		err = c.storage.Update(ctx, mTx, dbTx)
		if err != nil {
			log.Errorf("failed to update monitored tx to reorg status: %v", err)
			return err
		}
		log.Infof("monitored tx %v status updated to reorged", mTx.id)
	}
	log.Infof("reorg from block %v processed successfully", fromBlockNumber)
	return nil
}

// monitorTxs process all pending monitored tx
func (c *Client) monitorTxs(ctx context.Context) error {
	statusesFilter := []MonitoredTxStatus{MonitoredTxStatusCreated, MonitoredTxStatusSent, MonitoredTxStatusReorged}
	mTxs, err := c.storage.GetByStatus(ctx, nil, statusesFilter, nil)
	if err != nil {
		return fmt.Errorf("failed to get created monitored txs: %v", err)
	}

	log.Infof("found %v monitored tx to process", len(mTxs))

	for _, mTx := range mTxs {
		mTx := mTx // force variable shadowing to avoid pointer conflicts
		mTxLog := log.WithFields("monitoredTx", mTx.id)
		mTxLog.Info("processing")

		// check if any of the txs in the history was mined
		mined := false
		var receipt *types.Receipt
		hasFailedReceipts := false
		allHistoryTxMined := true
		for txHash := range mTx.history {
			mined, receipt, err = c.etherman.CheckTxWasMined(ctx, txHash)
			if err != nil {
				mTxLog.Errorf("failed to check if tx %v was mined: %v", txHash.String(), err)
				continue
			}

			// if the tx is not mined yet, check that not all the tx were mined and go to the next
			if !mined {
				allHistoryTxMined = false
				continue
			}

			// if the tx was mined successfully we can break the loop and proceed
			if receipt.Status == types.ReceiptStatusSuccessful {
				break
			}

			// if the tx was mined but failed, we continue to consider it was not mined
			// and store the failed receipt to be used to check if nonce needs to be reviewed
			mined = false
			hasFailedReceipts = true
		}

		// we need to check if we need to review the nonce carefully, to avoid sending
		// duplicated data to the block chain.
		//
		// if we have failed receipts, this means at least one of the generated txs was mined
		// so maybe the current nonce was already consumed, then we need to check if there are
		// tx that were not mined yet, if so, we just need to wait, because maybe one of them
		// will get mined successfully
		//
		// in case of all tx were mined and none of them were mined successfully, we need to
		// review the nonce
		if hasFailedReceipts && allHistoryTxMined {
			mTxLog.Infof("nonce needs to be updated")
			err := c.ReviewMonitoredTxNonce(ctx, &mTx)
			if err != nil {
				mTxLog.Errorf("failed to review monitored tx nonce: %v", err)
				continue
			}
			err = c.storage.Update(ctx, mTx, nil)
			if err != nil {
				mTxLog.Errorf("failed to update monitored tx nonce change: %v", err)
				continue
			}
		}

		// if the history size reaches the max history size, this means something is really wrong with
		// this Tx and we are not able to identify automatically, so we mark this as failed to let the
		// caller know something is not right and needs to be review and to avoid to monitor this
		// tx infinitely
		// if len(mTx.history) == maxHistorySize {
		// 	mTx.status = MonitoredTxStatusFailed
		// 	mTxLog.Infof("marked as failed because reached the history size limit: %v", err)
		// 	// update monitored tx changes into storage
		// 	err = c.storage.Update(ctx, mTx, nil)
		// 	if err != nil {
		// 		mTxLog.Errorf("failed to update monitored tx when max history size limit reached: %v", err)
		// 		continue
		// 	}
		// }

		var signedTx *types.Transaction
		if !mined {
			// if is a reorged, move to the next
			if mTx.status == MonitoredTxStatusReorged {
				continue
			}

			// review tx and increase gas and gas price if needed
			if mTx.status == MonitoredTxStatusSent {
				err := c.ReviewMonitoredTx(ctx, &mTx)
				if err != nil {
					mTxLog.Errorf("failed to review monitored tx: %v", err)
					continue
				}
				err = c.storage.Update(ctx, mTx, nil)
				if err != nil {
					mTxLog.Errorf("failed to update monitored tx review change: %v", err)
					continue
				}
			}

			// rebuild transaction
			tx := mTx.Tx()
			mTxLog.Debugf("unsigned tx %v created", tx.Hash().String(), mTx.id)

			// sign tx
			signedTx, err = c.etherman.SignTx(ctx, mTx.from, tx)
			if err != nil {
				mTxLog.Errorf("failed to sign tx %v created from monitored tx %v: %v", tx.Hash().String(), mTx.id, err)
				continue
			}
			mTxLog.Debugf("signed tx %v created", signedTx.Hash().String())

			// add tx to monitored tx history
			err = mTx.AddHistory(signedTx)
			if errors.Is(err, ErrAlreadyExists) {
				mTxLog.Infof("signed tx already existed in the history")
			} else if err != nil {
				mTxLog.Errorf("failed to add signed tx to monitored tx %v history: %v", mTx.id, err)
				continue
			} else {
				// update monitored tx changes into storage
				err = c.storage.Update(ctx, mTx, nil)
				if err != nil {
					mTxLog.Errorf("failed to update monitored tx: %v", err)
					continue
				}
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
				mTxLog.Infof("signed tx sent to the network: %v", signedTx.Hash().String())
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
				}
			} else {
				mTxLog.Infof("signed tx already found in the network")
			}

			log.Infof("waiting signedTx to be mined...")

			// wait tx to get mined
			mined, err = c.etherman.WaitTxToBeMined(ctx, signedTx, c.cfg.WaitTxToBeMined.Duration)
			if err != nil {
				mTxLog.Errorf("failed to wait tx to be mined: %v", err)
				continue
			}
			if !mined {
				log.Infof("signedTx not mined yet and timeout has been reached")
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
				mTxLog.Info("confirmed")
				mTx.status = MonitoredTxStatusConfirmed
			}
		} else {
			// if we should continue to monitor, we move to the next one and this will
			// be reviewed in the next monitoring cycle
			if c.shouldContinueToMonitorThisTx(ctx, receipt) {
				continue
			}
			mTxLog.Info("failed")
			// otherwise we understand this monitored tx has failed
			mTx.status = MonitoredTxStatusFailed
		}

		// update monitored tx changes into storage
		err = c.storage.Update(ctx, mTx, nil)
		if err != nil {
			mTxLog.Errorf("failed to update monitored tx: %v", err)
			continue
		}
	}

	return nil
}

// shouldContinueToMonitorThisTx checks the the tx receipt and decides if it should
// continue or not to monitor the monitored tx related to the tx from this receipt
func (c *Client) shouldContinueToMonitorThisTx(ctx context.Context, receipt *types.Receipt) bool {
	// if the receipt has a is successful result, stop monitoring
	if receipt.Status == types.ReceiptStatusSuccessful {
		return false
	}

	tx, _, err := c.etherman.GetTx(ctx, receipt.TxHash)
	if err != nil {
		log.Errorf("failed to get tx when monitored tx identified as failed, tx : %v", receipt.TxHash.String(), err)
		return false
	}
	_, err = c.etherman.GetRevertMessage(ctx, tx)
	if err != nil {
		// if the error when getting the revert message is not identified, continue to monitor
		if err.Error() == ErrExecutionReverted.Error() {
			return true
		} else {
			log.Errorf("failed to get revert message for monitored tx identified as failed, tx %v: %v", receipt.TxHash.String(), err)
		}
	}
	// if nothing weird was found, stop monitoring
	return false
}

// ReviewMonitoredTx checks if some field needs to be updated
// accordingly to the current information stored and the current
// state of the blockchain
func (c *Client) ReviewMonitoredTx(ctx context.Context, mTx *monitoredTx) error {
	mTxLog := log.WithFields("monitoredTx", mTx.id)
	mTxLog.Debug("reviewing")
	// get gas
	gas, err := c.etherman.EstimateGas(ctx, mTx.from, mTx.to, mTx.value, mTx.data)
	if err != nil {
		err := fmt.Errorf("failed to estimate gas: %w", err)
		mTxLog.Errorf(err.Error())
		return err
	}

	// check gas
	if gas > mTx.gas {
		mTxLog.Infof("monitored tx gas updated from %v to %v", mTx.gas, gas)
		mTx.gas = gas
	}

	// get gas price
	gasPrice, err := c.etherman.SuggestedGasPrice(ctx)
	if err != nil {
		err := fmt.Errorf("failed to get suggested gas price: %w", err)
		mTxLog.Errorf(err.Error())
		return err
	}

	// check gas price
	if gasPrice.Cmp(mTx.gasPrice) == 1 {
		mTxLog.Infof("monitored tx gas price updated from %v to %v", mTx.gasPrice.String(), gasPrice.String())
		mTx.gasPrice = gasPrice
	}
	return nil
}

// ReviewMonitoredTxNonce checks if the nonce needs to be updated accordingly to
// the current nonce of the sender account.
//
// IMPORTANT: Nonce is reviewed apart from the other fields because it is a very
// sensible information and can make duplicated data to be sent to the blockchain,
// causing possible side effects and wasting resources on taxes.
func (c *Client) ReviewMonitoredTxNonce(ctx context.Context, mTx *monitoredTx) error {
	mTxLog := log.WithFields("monitoredTx", mTx.id)
	mTxLog.Debug("reviewing nonce")
	nonce, err := c.etherman.CurrentNonce(ctx, mTx.from)
	if err != nil {
		err := fmt.Errorf("failed to estimate gas: %w", err)
		mTxLog.Errorf(err.Error())
		return err
	}

	if nonce > mTx.nonce {
		mTxLog.Infof("monitored tx nonce updated from %v to %v", mTx.nonce, nonce)
		mTx.nonce = nonce
	}

	return nil
}

// logErrorAndWait used when an error is detected before trying again
func (c *Client) logErrorAndWait(msg string, err error) {
	log.Errorf(msg, err)
	time.Sleep(failureIntervalInSeconds * time.Second)
}

// ResultHandler used by the caller to handle results
// when processing monitored txs
type ResultHandler func(MonitoredTxResult, pgx.Tx)

// ProcessPendingMonitoredTxs will check all monitored txs of this owner
// and wait until all of them are either confirmed or failed before continuing
//
// for the confirmed and failed ones, the resultHandler will be triggered
func (c *Client) ProcessPendingMonitoredTxs(ctx context.Context, owner string, resultHandler ResultHandler, dbTx pgx.Tx) {
	statusesFilter := []MonitoredTxStatus{
		MonitoredTxStatusCreated,
		MonitoredTxStatusSent,
		MonitoredTxStatusFailed,
		MonitoredTxStatusConfirmed,
		MonitoredTxStatusReorged,
	}
	// keep running until there are pending monitored txs
	for {
		results, err := c.ResultsByStatus(ctx, owner, statusesFilter, dbTx)
		if err != nil {
			// if something goes wrong here, we log, wait a bit and keep it in the infinite loop to not unlock the caller.
			log.Errorf("failed to get results by statuses from eth tx manager to monitored txs err: ", err)
			time.Sleep(time.Second)
			continue
		}

		if len(results) == 0 {
			// if there are not pending monitored txs, stop
			return
		}

		for _, result := range results {
			resultLog := log.WithFields("owner", owner, "id", result.ID)

			// if the result is confirmed, we set it as done do stop looking into this monitored tx
			if result.Status == MonitoredTxStatusConfirmed {
				err := c.setStatusDone(ctx, owner, result.ID, dbTx)
				if err != nil {
					resultLog.Errorf("failed to set monitored tx as done, err: %v", err)
					// if something goes wrong at this point, we skip this result and move to the next.
					// this result is going to be handled again in the next cycle by the outer loop.
					continue
				} else {
					resultLog.Info("monitored tx confirmed")
				}
				resultHandler(result, dbTx)
				continue
			}

			// if the result is failed, we need to go around it and rebuild a batch verification
			if result.Status == MonitoredTxStatusFailed {
				resultHandler(result, dbTx)
				continue
			}

			// if the result is either not confirmed or failed, it means we need to wait until it gets confirmed of failed.
			for {
				// wait before refreshing the result info
				time.Sleep(time.Second)

				// refresh the result info
				result, err := c.Result(ctx, owner, result.ID, dbTx)
				if err != nil {
					resultLog.Errorf("failed to get monitored tx result, err: %v", err)
					continue
				}

				// if the result status is confirmed or failed, breaks the wait loop
				if result.Status == MonitoredTxStatusConfirmed || result.Status == MonitoredTxStatusFailed {
					break
				}

				resultLog.Infof("waiting for monitored tx to get confirmed, status: %v", result.Status.String())
			}
		}
	}
}
