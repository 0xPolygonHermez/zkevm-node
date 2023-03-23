package pool

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// bridgeClaimMethodSignature for tracking bridgeClaimMethodSignature method
	bridgeClaimMethodSignature = "0x2cffd02e"
	retryInterval              = 2 * time.Second
)

var (
	// ErrAlreadyKnown is returned if the transactions is already contained
	// within the pool.
	ErrAlreadyKnown = errors.New("already known")

	// ErrReplaceUnderpriced is returned if a transaction is attempted to be replaced
	// with a different one without the required price bump.
	ErrReplaceUnderpriced = errors.New("replacement transaction underpriced")
)

// Pool is an implementation of the Pool interface
// that uses a postgres database to store the data
type Pool struct {
	storage
	state                   stateInterface
	l2BridgeAddr            common.Address
	chainID                 uint64
	cfg                     Config
	minSuggestedGasPrice    *big.Int
	minSuggestedGasPriceMux *sync.RWMutex
}

// NewPool creates and initializes an instance of Pool
func NewPool(cfg Config, s storage, st stateInterface, l2BridgeAddr common.Address, chainID uint64) *Pool {
	return &Pool{
		cfg:                     cfg,
		storage:                 s,
		state:                   st,
		l2BridgeAddr:            l2BridgeAddr,
		chainID:                 chainID,
		minSuggestedGasPriceMux: new(sync.RWMutex),
	}
}

// StartPollingMinSuggestedGasPrice starts polling the minimum suggested gas price
func (p *Pool) StartPollingMinSuggestedGasPrice(ctx context.Context) {
	err := p.pollMinSuggestedGasPrice(ctx)
	if err != nil && err != state.ErrNotFound {
		log.Fatalf("Error polling min suggested gas price: %v", err)
	}

	// Retrying until we have a valid minSuggestedGasPrice
	for err == state.ErrNotFound {
		err = p.pollMinSuggestedGasPrice(ctx)
		log.Infof("Retrying to poll min suggested gas price ...")
		time.Sleep(retryInterval)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(p.cfg.PollMinSuggestedGasPriceInterval.Duration):
				err = p.pollMinSuggestedGasPrice(ctx)
				if err != nil {
					log.Errorf("Error polling min suggested gas price: %v", err)
				}
			}
		}
	}()
}

// AddTx adds a transaction to the pool with the pending state
func (p *Pool) AddTx(ctx context.Context, tx types.Transaction, ip string) error {
	if err := p.validateTx(ctx, tx); err != nil {
		return err
	}

	return p.StoreTx(ctx, tx, ip, false)
}

// StoreTx adds a transaction to the pool with the pending state
func (p *Pool) StoreTx(ctx context.Context, tx types.Transaction, ip string, isWIP bool) error {
	poolTx := Transaction{
		Transaction: tx,
		Status:      TxStatusPending,
		IsClaims:    false,
		ReceivedAt:  time.Now(),
		IsWIP:       isWIP,
		IP:          ip,
	}

	poolTx.IsClaims = poolTx.IsClaimTx(p.l2BridgeAddr, p.cfg.FreeClaimGasLimit)

	// Execute transaction to calculate its zkCounters
	zkCounters, err, isOOC := p.PreExecuteTx(ctx, tx)
	if err != nil {
		log.Debugf("PreExecuteTx error (this can be ignored): %v", err)

		if isOOC {
			event := &state.Event{
				EventType: state.EventType_Prexecution_OOC,
				Timestamp: time.Now(),
				IP:        ip,
				TxHash:    tx.Hash(),
			}

			err := p.state.AddEvent(ctx, event, nil)
			if err != nil {
				log.Errorf("Error adding event: %v", err)
			}
		}
	}
	poolTx.ZKCounters = zkCounters

	return p.storage.AddTx(ctx, poolTx)
}

// PreExecuteTx executes a transaction to calculate its zkCounters
func (p *Pool) PreExecuteTx(ctx context.Context, tx types.Transaction) (state.ZKCounters, error, bool) {
	processBatchResponse, err := p.state.PreProcessTransaction(ctx, &tx, nil)
	if err != nil {
		return state.ZKCounters{}, err, false
	}
	return processBatchResponse.UsedZkCounters, processBatchResponse.ExecutorError, !processBatchResponse.IsBatchProcessed
}

// GetPendingTxs from the pool
// limit parameter is used to limit amount of pending txs from the db,
// if limit = 0, then there is no limit
func (p *Pool) GetPendingTxs(ctx context.Context, isClaims bool, limit uint64) ([]Transaction, error) {
	return p.storage.GetTxsByStatus(ctx, TxStatusPending, isClaims, limit)
}

// GetNonWIPPendingTxs from the pool
// limit parameter is used to limit amount of pending txs from the db,
// if limit = 0, then there is no limit
func (p *Pool) GetNonWIPPendingTxs(ctx context.Context, isClaims bool, limit uint64) ([]Transaction, error) {
	return p.storage.GetNonWIPTxsByStatus(ctx, TxStatusPending, isClaims, limit)
}

// GetSelectedTxs gets selected txs from the pool db
func (p *Pool) GetSelectedTxs(ctx context.Context, limit uint64) ([]Transaction, error) {
	return p.storage.GetTxsByStatus(ctx, TxStatusSelected, false, limit)
}

// GetPendingTxHashesSince returns the hashes of pending tx since the given date.
func (p *Pool) GetPendingTxHashesSince(ctx context.Context, since time.Time) ([]common.Hash, error) {
	return p.storage.GetPendingTxHashesSince(ctx, since)
}

// UpdateTxStatus updates a transaction state accordingly to the
// provided state and hash
func (p *Pool) UpdateTxStatus(ctx context.Context, hash common.Hash, newStatus TxStatus, isWIP bool) error {
	return p.storage.UpdateTxStatus(ctx, hash, newStatus, isWIP)
}

// SetGasPrice allows an external component to define the gas price
func (p *Pool) SetGasPrice(ctx context.Context, gasPrice uint64) error {
	return p.storage.SetGasPrice(ctx, gasPrice)
}

// GetGasPrice returns the current gas price
func (p *Pool) GetGasPrice(ctx context.Context) (uint64, error) {
	return p.storage.GetGasPrice(ctx)
}

// CountPendingTransactions get number of pending transactions
// used in bench tests
func (p *Pool) CountPendingTransactions(ctx context.Context) (uint64, error) {
	return p.storage.CountTransactionsByStatus(ctx, TxStatusPending)
}

// IsTxPending check if tx is still pending
func (p *Pool) IsTxPending(ctx context.Context, hash common.Hash) (bool, error) {
	return p.storage.IsTxPending(ctx, hash)
}

func (p *Pool) validateTx(ctx context.Context, tx types.Transaction) error {
	// check chain id
	txChainID := tx.ChainId().Uint64()
	if txChainID != p.chainID && txChainID != 0 {
		return ErrInvalidChainID
	}

	// Accept only legacy transactions until EIP-2718/2930 activates.
	if tx.Type() != types.LegacyTxType {
		return ErrTxTypeNotSupported
	}
	// Reject transactions over defined size to prevent DOS attacks
	if tx.Size() > p.cfg.MaxTxBytesSize {
		return ErrOversizedData
	}
	// Reject transactions with a gas price lower than the minimum gas price
	p.minSuggestedGasPriceMux.RLock()
	gasPriceCmp := tx.GasPrice().Cmp(p.minSuggestedGasPrice)
	p.minSuggestedGasPriceMux.RUnlock()
	if gasPriceCmp == -1 {
		return ErrGasPrice
	}
	// Transactions can't be negative. This may never happen using RLP decoded
	// transactions but may occur if you create a transaction using the RPC.
	if tx.Value().Sign() < 0 {
		return ErrNegativeValue
	}
	// Make sure the transaction is signed properly.
	if err := state.CheckSignature(tx); err != nil {
		return ErrInvalidSender
	}
	from, err := state.GetSender(tx)
	if err != nil {
		return ErrInvalidSender
	}

	lastL2BlockNumber, err := p.state.GetLastL2BlockNumber(ctx, nil)
	if err != nil {
		return err
	}

	nonce, err := p.state.GetNonce(ctx, from, lastL2BlockNumber, nil)
	if err != nil {
		return err
	}
	// Ensure the transaction adheres to nonce ordering
	if nonce > tx.Nonce() {
		return ErrNonceTooLow
	}

	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	balance, err := p.state.GetBalance(ctx, from, lastL2BlockNumber, nil)
	if err != nil {
		return err
	}

	if balance.Cmp(tx.Cost()) < 0 {
		return ErrInsufficientFunds
	}

	// Ensure the transaction has more gas than the basic tx fee.
	intrGas, err := IntrinsicGas(tx)
	if err != nil {
		return err
	}
	if tx.Gas() < intrGas {
		return ErrIntrinsicGas
	}

	// try to get a transaction from the pool with the same nonce to check
	// if the new one has a price bump
	oldTxs, err := p.storage.GetTxsByFromAndNonce(ctx, from, tx.Nonce())
	if err != nil {
		return err
	}

	// check if the new transaction has more gas than all the other txs in the pool
	// with the same from and nonce to be able to replace the current txs by the new
	// when being selected
	for _, oldTx := range oldTxs {
		// discard invalid txs
		if oldTx.Status == TxStatusInvalid || oldTx.Status == TxStatusFailed {
			continue
		}

		oldTxPrice := new(big.Int).Mul(oldTx.GasPrice(), new(big.Int).SetUint64(oldTx.Gas()))
		txPrice := new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(tx.Gas()))

		if oldTx.Hash() == tx.Hash() {
			return ErrAlreadyKnown
		}

		// if old Tx Price is higher than the new tx price, it returns an error
		if oldTxPrice.Cmp(txPrice) > 0 {
			return ErrReplaceUnderpriced
		}
	}

	// Executor field size requirements check
	if err := p.checkTxFieldCompatibilityWithExecutor(ctx, tx); err != nil {
		return err
	}

	return nil
}

func (p *Pool) pollMinSuggestedGasPrice(ctx context.Context) error {
	fromTimestamp := time.Now().UTC().Add(-p.cfg.MinSuggestedGasPriceInterval.Duration)
	gasPrice, err := p.storage.MinGasPriceSince(ctx, fromTimestamp)
	if err == state.ErrNotFound {
		log.Warnf("No suggested min gas price since: %v", fromTimestamp)
		return err
	} else if err != nil {
		log.Errorf("Error getting min gas price since: %v", fromTimestamp)
		return err
	} else {
		p.minSuggestedGasPriceMux.Lock()
		p.minSuggestedGasPrice = big.NewInt(0).SetUint64(gasPrice)
		p.minSuggestedGasPriceMux.Unlock()
		log.Infof("Min suggested gas price updated to: %d", gasPrice)
	}

	return nil
}

// checkTxFieldCompatibilityWithExecutor checks the field sizes of the transaction to make sure
// they ar compatible with the Executor needs
// GasLimit: 256 bits
// GasPrice: 256 bits
// Value: 256 bits
// Data: 30000 bytes
// Nonce: 64 bits
// To: 160 bits
// ChainId: 64 bits
func (p *Pool) checkTxFieldCompatibilityWithExecutor(ctx context.Context, tx types.Transaction) error {
	maxUint64BigInt := big.NewInt(0).SetUint64(math.MaxUint64)

	// GasLimit, Nonce and To fields are limited by their types, no need to check
	// Gas Price and Value are checked against the balance, and the max balance allowed
	// by the merkletree service is uint256, in this case, if the transaction has a
	// gas price or value bigger than uint256, the check against the balance will
	// reject the transaction

	dataSize := len(tx.Data())
	if dataSize > p.cfg.MaxTxDataBytesSize {
		return fmt.Errorf("data size bigger than allowed, current size is %v bytes and max allowed is %v bytes", dataSize, p.cfg.MaxTxDataBytesSize)
	}

	if tx.ChainId().Cmp(maxUint64BigInt) == 1 {
		return fmt.Errorf("chain id higher than allowed, max allowed is %v", uint64(math.MaxUint64))
	}

	return nil
}

// DeleteReorgedTransactions deletes transactions from the pool
func (p *Pool) DeleteReorgedTransactions(ctx context.Context, transactions []*types.Transaction) error {
	hashes := []common.Hash{}

	for _, tx := range transactions {
		hashes = append(hashes, tx.Hash())
	}

	return p.storage.DeleteTransactionsByHashes(ctx, hashes)
}

// UpdateTxWIPStatus updates a transaction wip status accordingly to the
// provided WIP status and hash
func (p *Pool) UpdateTxWIPStatus(ctx context.Context, hash common.Hash, isWIP bool) error {
	return p.storage.UpdateTxWIPStatus(ctx, hash, isWIP)
}

const (
	txDataNonZeroGas      uint64 = 16
	txGasContractCreation uint64 = 53000
	txGas                 uint64 = 21000
	txDataZeroGas         uint64 = 4
)

// IntrinsicGas computes the 'intrinsic gas' for a given transaction.
func IntrinsicGas(tx types.Transaction) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if tx.To() == nil {
		gas = txGasContractCreation
	} else {
		gas = txGas
	}
	dataLen := uint64(len(tx.Data()))
	// Bump the required gas by the amount of transactional data
	if dataLen > 0 {
		// Zero and non-zero bytes are priced differently
		var nz uint64
		for _, byt := range tx.Data() {
			if byt != 0 {
				nz++
			}
		}
		// Make sure we don't exceed uint64 for all data combinations
		nonZeroGas := txDataNonZeroGas
		if (math.MaxUint64-gas)/nonZeroGas < nz {
			return 0, ErrGasUintOverflow
		}
		gas += nz * nonZeroGas

		z := dataLen - nz
		if (math.MaxUint64-gas)/txDataZeroGas < z {
			return 0, ErrGasUintOverflow
		}
		gas += z * txDataZeroGas
	}
	return gas, nil
}
