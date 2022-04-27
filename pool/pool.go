package pool

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/helper"
)

const (
	// txSlotSize is used to calculate how many data slots a single transaction
	// takes up based on its size. The slots are used as DoS protection, ensuring
	// that validating a new transaction remains a constant operation (in reality
	// O(maxslots), where max slots are 4 currently).
	txSlotSize = 32 * 1024

	// txMaxSize is the maximum size a single transaction can have. This field has
	// non-trivial consequences: larger transactions are significantly harder and
	// more expensive to propagate; larger transactions also take more resources
	// to validate whether they fit into the pool or not.
	txMaxSize = 4 * txSlotSize // 128KB

	// bridgeClaimMethodSignature for tracking bridgeClaimMethodSignature method
	bridgeClaimMethodSignature = "0x122650ff"
)

// Pool is an implementation of the Pool interface
// that uses a postgres database to store the data
type Pool struct {
	storage                     storage
	state                       stateInterface
	l2GlobalExitRootManagerAddr common.Address
}

// NewPool creates and initializes an instance of Pool
func NewPool(s storage, st stateInterface, l2GlobalExitRootManagerAddr common.Address) *Pool {
	return &Pool{
		storage:                     s,
		state:                       st,
		l2GlobalExitRootManagerAddr: l2GlobalExitRootManagerAddr,
	}
}

// AddTx adds a transaction to the pool with the pending state
func (p *Pool) AddTx(ctx context.Context, tx types.Transaction) error {
	if err := p.validateTx(ctx, tx); err != nil {
		return err
	}

	poolTx := Transaction{
		Transaction: tx,
		State:       TxStatePending,
		IsClaims:    false,
		ReceivedAt:  time.Now(),
	}

	poolTx.IsClaims = poolTx.IsClaimTx(p.l2GlobalExitRootManagerAddr)

	return p.storage.AddTx(ctx, poolTx)
}

// GetPendingTxs from the pool
// limit parameter is used to limit amount of pending txs from the db,
// if limit = 0, then there is no limit
func (p *Pool) GetPendingTxs(ctx context.Context, isClaims bool, limit uint64) ([]Transaction, error) {
	return p.storage.GetTxsByState(ctx, TxStatePending, isClaims, limit)
}

// UpdateTxState updates a transaction state accordingly to the
// provided state and hash
func (p *Pool) UpdateTxState(ctx context.Context, hash common.Hash, newState TxState) error {
	return p.storage.UpdateTxState(ctx, hash, newState)
}

// UpdateTxsState updates transactions state accordingly to the provided state and hashes
func (p *Pool) UpdateTxsState(ctx context.Context, hashes []common.Hash, newState TxState) error {
	return p.storage.UpdateTxsState(ctx, hashes, newState)
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
	return p.storage.CountTransactionsByState(ctx, TxStatePending)
}

func (p *Pool) validateTx(ctx context.Context, tx types.Transaction) error {
	// Accept only legacy transactions until EIP-2718/2930 activates.
	if tx.Type() != types.LegacyTxType {
		return ErrTxTypeNotSupported
	}
	// Reject transactions over defined size to prevent DOS attacks
	if uint64(tx.Size()) > txMaxSize {
		return ErrOversizedData
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
	from, err := helper.GetSender(tx)
	if err != nil {
		return ErrInvalidSender
	}

	lastBatch, err := p.state.GetLastBatch(ctx, true, "")
	if err != nil {
		return err
	}
	lastBatchNumber := lastBatch.Number().Uint64()

	nonce, err := p.state.GetNonce(ctx, from, lastBatchNumber, "")
	if err != nil {
		return err
	}
	// Ensure the transaction adheres to nonce ordering
	if nonce > tx.Nonce() {
		return ErrNonceTooLow
	}

	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	balance, err := p.state.GetBalance(ctx, from, lastBatchNumber, "")
	if err != nil {
		return err
	}
	if balance.Cmp(tx.Cost()) < 0 {
		return ErrInsufficientFunds
	}

	return nil
}
