package pool

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	// ErrNotFound indicates an object has not been found for the search criteria used
	ErrNotFound = errors.New("object not found")

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
	chainID                 uint64
	cfg                     Config
	blockedAddresses        sync.Map
	minSuggestedGasPrice    *big.Int
	minSuggestedGasPriceMux *sync.RWMutex
	eventLog                *event.EventLog
	startTimestamp          time.Time
	gasPrices               GasPrices
	gasPricesMux            *sync.RWMutex
}

type preExecutionResponse struct {
	usedZkCounters       state.ZKCounters
	isExecutorLevelError bool
	isOOC                bool
	isOOG                bool
	isReverted           bool
	txResponse           *state.ProcessTransactionResponse
}

// GasPrices contains the gas prices for L2 and L1
type GasPrices struct {
	L2GasPrice uint64
	L1GasPrice uint64
}

// NewPool creates and initializes an instance of Pool
func NewPool(cfg Config, s storage, st stateInterface, chainID uint64, eventLog *event.EventLog) *Pool {
	startTimestamp := time.Now()
	p := &Pool{
		cfg:                     cfg,
		startTimestamp:          startTimestamp,
		storage:                 s,
		state:                   st,
		chainID:                 chainID,
		blockedAddresses:        sync.Map{},
		minSuggestedGasPriceMux: new(sync.RWMutex),
		eventLog:                eventLog,
		gasPrices:               GasPrices{0, 0},
		gasPricesMux:            new(sync.RWMutex),
	}

	p.refreshBlockedAddresses()
	go func(cfg *Config, p *Pool) {
		for {
			time.Sleep(cfg.IntervalToRefreshBlockedAddresses.Duration)
			p.refreshBlockedAddresses()
		}
	}(&cfg, p)

	go func(cfg *Config, p *Pool) {
		for {
			p.refreshGasPrices()
			time.Sleep(cfg.IntervalToRefreshGasPrices.Duration)
		}
	}(&cfg, p)

	return p
}

// refresGasPRices refreshes the gas price
func (p *Pool) refreshGasPrices() {
	gasPrices, err := p.GetGasPrices(context.Background())
	if err != nil {
		log.Error("failed to load gas prices")
		return
	}

	p.gasPricesMux.Lock()
	p.gasPrices = gasPrices
	p.gasPricesMux.Unlock()
}

// refreshBlockedAddresses refreshes the list of blocked addresses for the provided instance of pool
func (p *Pool) refreshBlockedAddresses() {
	blockedAddresses, err := p.storage.GetAllAddressesBlocked(context.Background())
	if err != nil {
		log.Error("failed to load blocked addresses")
		return
	}

	blockedAddressesMap := sync.Map{}
	for _, blockedAddress := range blockedAddresses {
		blockedAddressesMap.Store(blockedAddress.String(), 1)
		p.blockedAddresses.Store(blockedAddress.String(), 1)
	}

	unblockedAddresses := []string{}
	p.blockedAddresses.Range(func(key, value any) bool {
		addrHex := key.(string)
		_, found := blockedAddressesMap.Load(addrHex)
		if found {
			return true
		}

		unblockedAddresses = append(unblockedAddresses, addrHex)
		return true
	})

	for _, unblockedAddress := range unblockedAddresses {
		p.blockedAddresses.Delete(unblockedAddress)
	}
}

// StartPollingMinSuggestedGasPrice starts polling the minimum suggested gas price
func (p *Pool) StartPollingMinSuggestedGasPrice(ctx context.Context) {
	p.pollMinSuggestedGasPrice(ctx)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(p.cfg.PollMinAllowedGasPriceInterval.Duration):
				p.pollMinSuggestedGasPrice(ctx)
			}
		}
	}()
}

// AddTx adds a transaction to the pool with the pending state
func (p *Pool) AddTx(ctx context.Context, tx types.Transaction, ip string) error {
	poolTx := NewTransaction(tx, ip, false)
	if err := p.validateTx(ctx, *poolTx); err != nil {
		return err
	}

	return p.StoreTx(ctx, tx, ip, false)
}

// StoreTx adds a transaction to the pool with the pending state
func (p *Pool) StoreTx(ctx context.Context, tx types.Transaction, ip string, isWIP bool) error {
	// Execute transaction to calculate its zkCounters
	preExecutionResponse, err := p.preExecuteTx(ctx, tx)
	if errors.Is(err, runtime.ErrIntrinsicInvalidBatchGasLimit) {
		return ErrGasLimit
	} else if preExecutionResponse.isExecutorLevelError {
		// Do not add tx to the pool
		return err
	} else if err != nil {
		log.Errorf("Pre execution error: %v", err)
		return err
	}

	if preExecutionResponse.isOOC {
		event := &event.Event{
			ReceivedAt:  time.Now(),
			IPAddress:   ip,
			Source:      event.Source_Node,
			Component:   event.Component_Pool,
			Level:       event.Level_Warning,
			EventID:     event.EventID_PreexecutionOOC,
			Description: tx.Hash().String(),
		}

		err := p.eventLog.LogEvent(ctx, event)
		if err != nil {
			log.Errorf("error adding event: %v", err)
		}
		// Do not add tx to the pool
		return fmt.Errorf("out of counters")
	} else if preExecutionResponse.isOOG {
		event := &event.Event{
			ReceivedAt:  time.Now(),
			IPAddress:   ip,
			Source:      event.Source_Node,
			Component:   event.Component_Pool,
			Level:       event.Level_Warning,
			EventID:     event.EventID_PreexecutionOOG,
			Description: tx.Hash().String(),
		}

		err := p.eventLog.LogEvent(ctx, event)
		if err != nil {
			log.Errorf("error adding event: %v", err)
		}
	}

	poolTx := NewTransaction(tx, ip, isWIP)
	poolTx.ZKCounters = preExecutionResponse.usedZkCounters

	return p.storage.AddTx(ctx, *poolTx)
}

// preExecuteTx executes a transaction to calculate its zkCounters
func (p *Pool) preExecuteTx(ctx context.Context, tx types.Transaction) (preExecutionResponse, error) {
	response := preExecutionResponse{usedZkCounters: state.ZKCounters{}, isOOC: false, isOOG: false, isReverted: false}

	// TODO: Add effectivePercentage = 0xFF to the request (factor of 1) when gRPC message is updated
	processBatchResponse, err := p.state.PreProcessTransaction(ctx, &tx, nil)
	if err != nil {
		return response, err
	}

	if processBatchResponse.Responses != nil && len(processBatchResponse.Responses) > 0 {
		errorToCheck := processBatchResponse.Responses[0].RomError
		response.isReverted = errors.Is(errorToCheck, runtime.ErrExecutionReverted)
		response.isExecutorLevelError = processBatchResponse.IsExecutorLevelError
		response.isOOC = executor.IsROMOutOfCountersError(executor.RomErrorCode(errorToCheck))
		response.isOOG = errors.Is(errorToCheck, runtime.ErrOutOfGas)
		response.usedZkCounters = processBatchResponse.UsedZkCounters
		response.txResponse = processBatchResponse.Responses[0]
	}

	return response, nil
}

// GetPendingTxs from the pool
// limit parameter is used to limit amount of pending txs from the db,
// if limit = 0, then there is no limit
func (p *Pool) GetPendingTxs(ctx context.Context, limit uint64) ([]Transaction, error) {
	return p.storage.GetTxsByStatus(ctx, TxStatusPending, limit)
}

// GetNonWIPPendingTxs from the pool
func (p *Pool) GetNonWIPPendingTxs(ctx context.Context) ([]Transaction, error) {
	return p.storage.GetNonWIPPendingTxs(ctx)
}

// GetSelectedTxs gets selected txs from the pool db
func (p *Pool) GetSelectedTxs(ctx context.Context, limit uint64) ([]Transaction, error) {
	return p.storage.GetTxsByStatus(ctx, TxStatusSelected, limit)
}

// GetPendingTxHashesSince returns the hashes of pending tx since the given date.
func (p *Pool) GetPendingTxHashesSince(ctx context.Context, since time.Time) ([]common.Hash, error) {
	return p.storage.GetPendingTxHashesSince(ctx, since)
}

// UpdateTxStatus updates a transaction state accordingly to the
// provided state and hash
func (p *Pool) UpdateTxStatus(ctx context.Context, hash common.Hash, newStatus TxStatus, isWIP bool, failedReason *string) error {
	return p.storage.UpdateTxStatus(ctx, TxStatusUpdateInfo{
		Hash:         hash,
		NewStatus:    newStatus,
		IsWIP:        isWIP,
		FailedReason: failedReason,
	})
}

// SetGasPrices sets the current L2 Gas Price and L1 Gas Price
func (p *Pool) SetGasPrices(ctx context.Context, l2GasPrice uint64, l1GasPrice uint64) error {
	return p.storage.SetGasPrices(ctx, l2GasPrice, l1GasPrice)
}

// DeleteGasPricesHistoryOlderThan deletes gas prices older than a given date except the most recent one
func (p *Pool) DeleteGasPricesHistoryOlderThan(ctx context.Context, date time.Time) error {
	return p.storage.DeleteGasPricesHistoryOlderThan(ctx, date)
}

// GetGasPrices returns the current L2 Gas Price and L1 Gas Price
func (p *Pool) GetGasPrices(ctx context.Context) (GasPrices, error) {
	l2GasPrice, l1GasPrice, err := p.storage.GetGasPrices(ctx)
	return GasPrices{L1GasPrice: l1GasPrice, L2GasPrice: l2GasPrice}, err
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

func (p *Pool) validateTx(ctx context.Context, poolTx Transaction) error {
	// Make sure the transaction is signed properly.
	if err := state.CheckSignature(poolTx.Transaction); err != nil {
		return ErrInvalidSender
	}

	// check chain id
	txChainID := poolTx.ChainId().Uint64()
	if txChainID != p.chainID && txChainID != 0 {
		return ErrInvalidChainID
	}

	// Accept only legacy transactions until EIP-2718/2930 activates.
	if poolTx.Type() != types.LegacyTxType {
		return ErrTxTypeNotSupported
	}

	// gets tx sender for validations
	from, err := state.GetSender(poolTx.Transaction)
	if err != nil {
		return ErrInvalidSender
	}

	// Reject transactions over defined size to prevent DOS attacks
	if poolTx.Size() > p.cfg.MaxTxBytesSize {
		log.Infof("%v: %v", ErrOversizedData.Error(), from.String())
		return ErrOversizedData
	}

	// Transactions can't be negative. This may never happen using RLP decoded
	// transactions but may occur if you create a transaction using the RPC.
	if poolTx.Value().Sign() < 0 {
		return ErrNegativeValue
	}

	// check if sender is blocked
	_, blocked := p.blockedAddresses.Load(from.String())
	if blocked {
		log.Infof("%v: %v", ErrBlockedSender.Error(), from.String())
		return ErrBlockedSender
	}

	lastL2Block, err := p.state.GetLastL2Block(ctx, nil)
	if err != nil {
		return err
	}

	currentNonce, err := p.state.GetNonce(ctx, from, lastL2Block.Root())
	if err != nil {
		return err
	}
	// Ensure the transaction adheres to nonce ordering
	if poolTx.Nonce() < currentNonce {
		return ErrNonceTooLow
	}

	// check if sender has reached the limit of transactions in the pool
	if p.cfg.AccountQueue > 0 {
		// txCount, err := p.storage.CountTransactionsByFromAndStatus(ctx, from, TxStatusPending)
		// if err != nil {
		// 	return err
		// }
		// if txCount >= p.cfg.AccountQueue {
		// 	return ErrTxPoolAccountOverflow
		// }

		// Ensure the transaction does not jump out of the expected AccountQueue
		if poolTx.Nonce() > currentNonce+p.cfg.AccountQueue-1 {
			log.Infof("%v: %v", ErrNonceTooHigh.Error(), from.String())
			return ErrNonceTooHigh
		}
	}

	// check if the pool is full
	if p.cfg.GlobalQueue > 0 {
		txCount, err := p.storage.CountTransactionsByStatus(ctx, TxStatusPending)
		if err != nil {
			return err
		}
		if txCount >= p.cfg.GlobalQueue {
			return ErrTxPoolOverflow
		}
	}

	// Reject transactions with a gas price lower than the minimum gas price
	p.minSuggestedGasPriceMux.RLock()
	gasPriceCmp := poolTx.GasPrice().Cmp(p.minSuggestedGasPrice)
	p.minSuggestedGasPriceMux.RUnlock()
	if gasPriceCmp == -1 {
		return ErrGasPrice
	}

	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	balance, err := p.state.GetBalance(ctx, from, lastL2Block.Root())
	if err != nil {
		return err
	}

	if balance.Cmp(poolTx.Cost()) < 0 {
		return ErrInsufficientFunds
	}

	// Ensure the transaction has more gas than the basic poolTx fee.
	intrGas, err := IntrinsicGas(poolTx.Transaction)
	if err != nil {
		return err
	}
	if poolTx.Gas() < intrGas {
		return ErrIntrinsicGas
	}

	// try to get a transaction from the pool with the same nonce to check
	// if the new one has a price bump
	oldTxs, err := p.storage.GetTxsByFromAndNonce(ctx, from, poolTx.Nonce())
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
		txPrice := new(big.Int).Mul(poolTx.GasPrice(), new(big.Int).SetUint64(poolTx.Gas()))

		if oldTx.Hash() == poolTx.Hash() {
			return ErrAlreadyKnown
		}

		// if old Tx Price is higher than the new poolTx price, it returns an error
		if oldTxPrice.Cmp(txPrice) > 0 {
			return ErrReplaceUnderpriced
		}
	}

	// Executor field size requirements check
	if err := p.checkTxFieldCompatibilityWithExecutor(ctx, poolTx.Transaction); err != nil {
		return err
	}

	return nil
}

func (p *Pool) pollMinSuggestedGasPrice(ctx context.Context) {
	fromTimestamp := time.Now().UTC().Add(-p.cfg.MinAllowedGasPriceInterval.Duration)
	// Ensuring we don't use a timestamp before the pool start as it may be using older L1 gas price factor
	if fromTimestamp.Before(p.startTimestamp) {
		fromTimestamp = p.startTimestamp
	}

	l2GasPrice, err := p.storage.MinL2GasPriceSince(ctx, fromTimestamp)
	if err != nil {
		p.minSuggestedGasPriceMux.Lock()
		// Ensuring we always have suggested minimum gas price
		if p.minSuggestedGasPrice == nil {
			p.minSuggestedGasPrice = big.NewInt(0).SetUint64(p.cfg.DefaultMinGasPriceAllowed)
			log.Infof("Min allowed gas price updated to: %d", p.cfg.DefaultMinGasPriceAllowed)
		}
		p.minSuggestedGasPriceMux.Unlock()
		if err == state.ErrNotFound {
			log.Warnf("No suggested min gas price since: %v", fromTimestamp)
		} else {
			log.Errorf("Error getting min gas price since: %v", fromTimestamp)
		}
	} else {
		p.minSuggestedGasPriceMux.Lock()
		p.minSuggestedGasPrice = big.NewInt(0).SetUint64(l2GasPrice)
		p.minSuggestedGasPriceMux.Unlock()
		log.Infof("Min allowed gas price updated to: %d", l2GasPrice)
	}
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

// GetDefaultMinGasPriceAllowed return the configured DefaultMinGasPriceAllowed value
func (p *Pool) GetDefaultMinGasPriceAllowed() uint64 {
	return p.cfg.DefaultMinGasPriceAllowed
}

// GetL1GasPrice returns the L1 gas price
func (p *Pool) GetL1GasPrice() uint64 {
	p.gasPricesMux.RLock()
	gasPrices := p.gasPrices
	p.gasPricesMux.RUnlock()

	return gasPrices.L1GasPrice
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
