package gasprice

import (
	"context"
	"math/big"
	"sort"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

const sampleNumber = 3 // Number of transactions sampled in a batch.

// LastNL2BlocksGasPrice struct for gas price estimator last n l2 blocks.
type LastNL2BlocksGasPrice struct {
	lastL2BlockNumber uint64
	lastPrice         *big.Int

	cfg Config
	ctx context.Context

	cacheLock sync.RWMutex
	fetchLock sync.Mutex

	state stateInterface
	pool  poolInterface
}

// newLastNL2BlocksGasPriceSuggester init gas price suggester for last n l2 blocks strategy.
func newLastNL2BlocksGasPriceSuggester(ctx context.Context, cfg Config, state stateInterface, pool poolInterface) *LastNL2BlocksGasPrice {
	return &LastNL2BlocksGasPrice{
		cfg:   cfg,
		ctx:   ctx,
		state: state,
		pool:  pool,
	}
}

// UpdateGasPriceAvg for last n bathes strategy is not needed to implement this function.
func (g *LastNL2BlocksGasPrice) UpdateGasPriceAvg() {
	l2BlockNumber, err := g.state.GetLastL2BlockNumber(g.ctx, nil)
	if err != nil {
		log.Errorf("failed to get last l2 block number, err: %v", err)
	}
	g.cacheLock.RLock()
	lastL2BlockNumber, lastPrice := g.lastL2BlockNumber, g.lastPrice
	g.cacheLock.RUnlock()
	if l2BlockNumber == lastL2BlockNumber {
		log.Debug("Block is still the same, no need to update the gas price at the moment")
		return
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
		go g.getL2BlockTxsTips(g.ctx, number, sampleNumber, g.cfg.IgnorePrice, result, quit)
		sent++
		exp++
		number--
	}

	for exp > 0 {
		res := <-result
		if res.err != nil {
			close(quit)
			return
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

	// Store gasPrices
	factorAsPercentage := int64(g.cfg.Factor * 100) // nolint:gomnd
	factor := big.NewInt(factorAsPercentage)
	l1GasPriceDivBy100 := new(big.Int).Div(g.lastPrice, factor)
	l1GasPrice := l1GasPriceDivBy100.Mul(l1GasPriceDivBy100, big.NewInt(100)) // nolint:gomnd
	err = g.pool.SetGasPrices(g.ctx, g.lastPrice.Uint64(), l1GasPrice.Uint64())
	if err != nil {
		log.Errorf("failed to update gas price in poolDB, err: %v", err)
	}
}

// getL2BlockTxsTips calculates l2 block transaction gas fees.
func (g *LastNL2BlocksGasPrice) getL2BlockTxsTips(ctx context.Context, l2BlockNumber uint64, limit int, ignorePrice *big.Int, result chan results, quit chan struct{}) {
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
