package gaspriceestimator

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
)

const sampleNumber = 3 // Number of transactions sampled in a batch

type GasPriceEstimator interface {
	GetAvgGasPrice() (*big.Int, error)
	UpdateGasPriceAvg(newValue *big.Int)
}

func NewGasPriceEstimator(cfg Config, state state.State, pool *pool.PostgresPool) GasPriceEstimator {
	switch cfg.Type {
	case "allblocks":
		return NewGasPriceEstimatorAllBlocks()
	case "lastnblocks":
		return NewGasPriceEstimatorLastNBlocks(cfg, state)
	case "default":
		return NewDefaultGasPriceEstimator(cfg, pool)
	}
	return nil
}

type Default struct {
	cfg  Config
	pool *pool.PostgresPool
}

func (d *Default) GetAvgGasPrice() (*big.Int, error) {
	ctx := context.Background()
	gasPrice, err := d.pool.GetGasPrice(ctx)
	if errors.Is(err, state.ErrNotFound) {
		return big.NewInt(0), nil
	} else if err != nil {
		return nil, err
	}
	return new(big.Int).SetUint64(gasPrice), nil
}

func (d *Default) UpdateGasPriceAvg(newValue *big.Int) {
	panic("can't update gas price for default gas price estimator strategy")
}

func (d *Default) setDefaultGasPrice() {
	ctx := context.Background()
	err := d.pool.SetGasPrice(ctx, d.cfg.DefaultPriceWei)
	if err != nil {
		panic(fmt.Errorf("failed to set default gas price, err: %v", err))
	}
}

func NewDefaultGasPriceEstimator(cfg Config, pool *pool.PostgresPool) *Default {
	return &Default{pool: pool}
}

type AllBlocks struct {
	// Average gas price (rolling average)
	averageGasPrice      *big.Int // The average gas price that gets queried
	averageGasPriceCount *big.Int // Param used in the avg. gas price calculation

	agpMux sync.Mutex // Mutex for the averageGasPrice calculation
}

func NewGasPriceEstimatorAllBlocks() *AllBlocks {
	return &AllBlocks{
		averageGasPrice:      big.NewInt(0),
		averageGasPriceCount: big.NewInt(0),
	}
}

// UpdateGasPriceAvg Updates the rolling average value of the gas price
func (g *AllBlocks) UpdateGasPriceAvg(newValue *big.Int) {
	g.agpMux.Lock()

	g.averageGasPriceCount.Add(g.averageGasPriceCount, big.NewInt(1))

	differential := big.NewInt(0)
	differential.Div(newValue.Sub(newValue, g.averageGasPrice), g.averageGasPriceCount)

	g.averageGasPrice.Add(g.averageGasPrice, differential)

	g.agpMux.Unlock()
}

// GetAvgGasPrice get avg gas price from all the blocks
func (g *AllBlocks) GetAvgGasPrice() (*big.Int, error) {
	return g.averageGasPrice, nil
}

type LastNBlocks struct {
	lastBatchNumber uint64
	lastPrice       *big.Int

	cfg Config

	cacheLock sync.RWMutex
	fetchLock sync.Mutex

	state state.State
}

func (g *LastNBlocks) UpdateGasPriceAvg(newValue *big.Int) {
	panic("not used in this gas price estimation strategy")
}

func NewGasPriceEstimatorLastNBlocks(cfg Config, state state.State) *LastNBlocks {
	return &LastNBlocks{
		cfg:   cfg,
		state: state,
	}
}

func (g *LastNBlocks) GetAvgGasPrice() (*big.Int, error) {
	ctx := context.Background()

	batchNumber, err := g.state.GetLastBatchNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get last batch number, err: %v", err)
	}
	g.cacheLock.RLock()
	lastBatchNumber, lastPrice := g.lastBatchNumber, g.lastPrice
	if batchNumber == lastBatchNumber {
		return lastPrice, nil
	}
	g.cacheLock.RUnlock()

	g.fetchLock.Lock()
	defer g.fetchLock.Unlock()

	var (
		sent, exp int
		number    = lastBatchNumber
		result    = make(chan results, g.cfg.CheckBlocks)
		quit      = make(chan struct{})
		results   []*big.Int
	)

	for sent < g.cfg.CheckBlocks && number > 0 {
		go g.getBatchTxsTips(ctx, number, sampleNumber, g.cfg.IgnorePrice, result, quit)
		sent++
		exp++
		number--
	}

	for exp > 0 {
		res := <-result
		if res.err != nil {
			close(quit)
			return lastPrice, res.err
		}
		exp--

		if len(res.values) == 0 {
			res.values = []*big.Int{lastPrice}
		}
		results = append(results, res.values...)
	}

	price := lastPrice
	if len(results) > 0 {
		sort.Sort(bigIntArray(results))
		price = results[(len(results)-1)*g.cfg.Percentile/100]
	}
	if price.Cmp(g.cfg.MaxPrice) > 0 {
		price = g.cfg.MaxPrice
	}

	g.cacheLock.Lock()
	g.lastPrice = price
	g.lastBatchNumber = batchNumber
	g.cacheLock.Unlock()

	return price, nil
}

func (g *LastNBlocks) getBatchTxsTips(ctx context.Context, batchNum uint64, limit int, ignorePrice *big.Int, result chan results, quit chan struct{}) {
	txs, err := g.state.GetTxsByBatchNum(ctx, batchNum)
	if txs == nil {
		select {
		case result <- results{nil, err}:
		case <-quit:
		}
		return
	}
	sorter := newSorter(txs)
	sort.Sort(sorter)

	var prices []*big.Int
	for _, tx := range sorter.txs {
		tip := tx.GasTipCap()
		if ignorePrice != nil && tip.Cmp(ignorePrice) == -1 {
			continue
		}
		prices = append(prices, tip)
		if len(prices) >= limit {
			break
		}
	}
	select {
	case result <- results{prices, nil}:
	case <-quit:
	}
}

type results struct {
	values []*big.Int
	err    error
}

type txSorter struct {
	txs []*types.Transaction
}

func newSorter(txs []*types.Transaction) *txSorter {
	return &txSorter{
		txs: txs,
	}
}

func (s *txSorter) Len() int { return len(s.txs) }
func (s *txSorter) Swap(i, j int) {
	s.txs[i], s.txs[j] = s.txs[j], s.txs[i]
}
func (s *txSorter) Less(i, j int) bool {
	tip1 := s.txs[i].GasTipCap()
	tip2 := s.txs[j].GasTipCap()
	return tip1.Cmp(tip2) < 0
}

type bigIntArray []*big.Int

func (s bigIntArray) Len() int           { return len(s) }
func (s bigIntArray) Less(i, j int) bool { return s[i].Cmp(s[j]) < 0 }
func (s bigIntArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
