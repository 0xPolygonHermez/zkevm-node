package gasprice

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
)

const sampleNumber = 3 // Number of transactions sampled in a batch.

// LastNL2Blocks struct for gas price estimator last n l2 blocks.
type LastNL2Blocks struct {
	lastL2BlockNumber uint64
	lastPrice         *big.Int

	cfg Config

	cacheLock sync.RWMutex
	fetchLock sync.Mutex

	state stateInterface
}

// UpdateGasPriceAvg for last n bathes strategy is not needed to implement this function.
func (g *LastNL2Blocks) UpdateGasPriceAvg(newValue *big.Int) {}

// NewEstimatorLastNL2Blocks init gas price estimator for last n l2 blocks strategy.
func NewEstimatorLastNL2Blocks(cfg Config, state stateInterface) *LastNL2Blocks {
	return &LastNL2Blocks{
		cfg:   cfg,
		state: state,
	}
}

// GetAvgGasPrice calculate avg gas price from last n l2 blocks.
func (g *LastNL2Blocks) GetAvgGasPrice(ctx context.Context) (*big.Int, error) {
	l2BlockNumber, err := g.state.GetLastL2BlockNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get last l2 block number, err: %v", err)
	}
	g.cacheLock.RLock()
	lastL2BlockNumber, lastPrice := g.lastL2BlockNumber, g.lastPrice
	g.cacheLock.RUnlock()
	if l2BlockNumber == lastL2BlockNumber {
		return lastPrice, nil
	}

	g.fetchLock.Lock()
	defer g.fetchLock.Unlock()

	var (
		sent, exp int
		number    = lastL2BlockNumber
		result    = make(chan results, g.cfg.CheckBlocks)
		quit      = make(chan struct{})
		results   []*big.Int
	)

	for sent < g.cfg.CheckBlocks && number > 0 {
		go g.getL2BlockTxsTips(ctx, number, sampleNumber, g.cfg.IgnorePrice, result, quit)
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
	g.lastL2BlockNumber = l2BlockNumber
	g.cacheLock.Unlock()

	return price, nil
}

// getL2BlockTxsTips calculates l2 block transaction gas fees.
func (g *LastNL2Blocks) getL2BlockTxsTips(ctx context.Context, l2BlockNumber uint64, limit int, ignorePrice *big.Int, result chan results, quit chan struct{}) {
	txs, err := g.state.GetTxsByBlockNumber(ctx, l2BlockNumber, nil)
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
