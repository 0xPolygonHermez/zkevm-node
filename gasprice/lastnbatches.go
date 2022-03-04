package gasprice

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
)

const sampleNumber = 3 // Number of transactions sampled in a batch.

// LastNBatches struct for gas price estimator last n batches.
type LastNBatches struct {
	lastBatchNumber uint64
	lastPrice       *big.Int

	cfg Config

	cacheLock sync.RWMutex
	fetchLock sync.Mutex

	state localState
}

// UpdateGasPriceAvg for last n bathes strategy is not needed to implement this function.
func (g *LastNBatches) UpdateGasPriceAvg(newValue *big.Int) {}

// NewEstimatorLastNBatches init gas price estimator for last n batches strategy.
func NewEstimatorLastNBatches(cfg Config, state localState) *LastNBatches {
	return &LastNBatches{
		cfg:   cfg,
		state: state,
	}
}

// GetAvgGasPrice calculate avg gas price from last n batches.
func (g *LastNBatches) GetAvgGasPrice(ctx context.Context) (*big.Int, error) {
	batchNumber, err := g.state.GetLastBatchNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get last batch number, err: %v", err)
	}
	g.cacheLock.RLock()
	lastBatchNumber, lastPrice := g.lastBatchNumber, g.lastPrice
	g.cacheLock.RUnlock()
	if batchNumber == lastBatchNumber {
		return lastPrice, nil
	}

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

// getBatchTxsTips calculates batch transaction gas fees.
func (g *LastNBatches) getBatchTxsTips(ctx context.Context, batchNum uint64, limit int, ignorePrice *big.Int, result chan results, quit chan struct{}) {
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
