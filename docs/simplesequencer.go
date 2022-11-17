package docs

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// STRUCTS

type TxStatus struct {
	CheckedAtRoot          common.Hash
	WasExecutedSuccessfuly bool
	WasNonceTooBig         bool
	DidntHadBalance        bool
	// WrittenStorage  map[common.Hash]*big.Int    // Modified storage positions. Not needed yet
	// ReadStorage  map[common.Hash]*big.Int       // Consumed storage positions. Not needed yet
}

type ModLog struct {
	NewRoot         common.Hash
	From            common.Address              // Account that send the tx
	Nonce           uint64                      // Nonce of the sender aftter executing the tx
	UpdatedBalances map[common.Address]*big.Int // Modified balances per account
	// WrittenStorage  map[common.Hash]*big.Int    // Modified storage positions. Not needed yet
}

var StateModLog map[common.Hash]ModLog // state root => modifications

type PoolTx struct {
	*types.Transaction // Use raw data instead?
	Status             TxStatus
	From               common.Address
}

func (tx *PoolTx) UpdateStatus() error {
	needsReExecution := false
	nextRoot := StateModLog[tx.Status.CheckedAtRoot].NewRoot
	for {
		if mod, ok := StateModLog[nextRoot]; !ok {
			break
		} else {
			// THIS NEEDS TO BE HUGELY IMPORVED, the goal should be to only execute txs once (when their nonce is ready), OOC due to state change canbe ignored for now
			// Another option could be force setting current nonce to simulate execution... this could be done already at jRPC level
			if mod.From == tx.From && mod.Nonce == tx.Nonce()+1 {
				needsReExecution = true
				break
			}
			if _, ok := mod.UpdatedBalances[tx.From]; ok {
				needsReExecution = true
				break
			}
			nextRoot = mod.NewRoot
		}
	}
	if needsReExecution {
		if err := tx.ExecuteAndSetStatus(); err != nil {
			return err
		}
	} else {
		tx.Status.CheckedAtRoot = nextRoot
	}
	return nil
}

func (tx *PoolTx) ExecuteAndSetStatus() error {
	// Use executor
	// Parse result
	// If nonce too low or OOC return "must delete tx"
	// Set tx.Status
	return nil
}

func (tx *PoolTx) ExecuteOnTopOfAndSetStatus(root common.Hash) error {
	// Use executor
	// Parse result
	// If nonce too low or OOC return "must delete tx"
	// Set tx.Status
	return nil
}

// JSONRPC

type JSONRPC struct {
	pool Pool
}

func (rpc *JSONRPC) AddTx(tx *types.Transaction) {
	// 1. tx pre-checks...
	poolTx := PoolTx{tx, TxStatus{}, common.Address{}}
	if err := poolTx.ExecuteAndSetStatus(); err != nil {
		rpc.pool.StoreTx(poolTx)
	}
}

// Pool

type Pool struct{}

func (pool *Pool) StoreTx(tx PoolTx) {
	// Store
}

func (pool *Pool) DeleteTx(txHash common.Hash) {
	// Store
}

func (pool *Pool) GetTxs(from common.Hash, limit int) []PoolTx {
	// Query storage
	return []PoolTx{}
}

// SEQUENCER

type Sequencer struct {
	pool             Pool
	lastPoolTxLoaded common.Hash
	brokerCh         chan PoolTx
	txGroups         []TxGroup // Only one group at the beginning
	updateGERCh      chan common.Hash
	forcedBatchCh    chan state.Batch
}

type TxGroup struct {
	Ready             []*PoolTx                   // Txs that should get successfuly executed given it's order. Order should maximize profit
	BlockedByNonce    map[common.Address][]PoolTx // Each account may have many txs with a future noce (when nonce gap)
	BlockedByBalance  map[common.Address]PoolTx   // Each account can only have one tx blocked by balance (current nonce assumed)
	BlockedByGasPrice []PoolTx                    // Txs that won't be executed due to gas price too low
}

func (g *TxGroup) Add(tx PoolTx) error {
	// Add tx in corrct queue [Ready,BlockedByNonce, BlockedByBalance, BlockedByGasPrice] in correct order
	// If same nonce exists with higer gas price return error "must delete tx"
	return nil
}

func (g *TxGroup) PopBestTx() *PoolTx {
	old := g.Ready
	n := len(old)
	tx := old[0]
	old[n-1] = nil // avoid memory leak
	g.Ready = old[0 : n-1]
	return tx
}

func (g *TxGroup) Start() {
	for {

	}
}

func (s *Sequencer) Start(ctx context.Context) {
	// Put group to work, when having multiple/dynamic groups this may change
	go s.txGroups[0].Start()
	// Start broker
	go s.groupBroker()
	// Loop load txs from pool
	tickerLoadFromPool := time.NewTicker(time.Second)
	go func() {
		for {
			s.loadFromPool(ctx, tickerLoadFromPool)
		}
	}()
}

func (s *Sequencer) loadFromPool(ctx context.Context, ticker *time.Ticker) {
	txs := s.pool.GetTxs(s.lastPoolTxLoaded, 10)
	// This could be done in paralel
	for _, tx := range txs {
		if err := tx.UpdateStatus(); err != nil {
			s.pool.DeleteTx(tx.Hash())
		}
		// Send tx to broker
		s.brokerCh <- tx
	}
}

func (s *Sequencer) groupBroker() {
	// Could have many broker routines but could lose precision
	for {
		tx := <-s.brokerCh
		// group selection logic in the future
		if err := s.txGroups[0].Add(tx); err != nil {
			s.pool.DeleteTx(tx.Hash())
		}
	}
}

func (s *Sequencer) finalizer() {
	for {
		select {
		case <-s.forcedBatchCh:
			// If batch WIP is not full, can we wait a bit? (deathline to execute forced batch update)
			// Process forced batch
		case <-s.updateGERCh:
			// If batch WIP is not full, can we wait a bit? (deathline to update GER/timestamp)
			// Close batch, open batch (process wit no txs and new GER / timestamp)
		default:
			// When having multiple groups this could be refactored to use a selector that keeps popping txs from the different groups
			// and comunicate with the finalizer through channel
			tx := s.txGroups[0].PopBestTx()
			if tx != nil {
				// execute tx
				// if failed, send to group broker through chan
				// if ok:
				// send to statedb manager through chan
				// send state mod log update through chan
			}
		}
	}
}
