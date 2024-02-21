package jsonrpc

import (
	"context"
	"errors"
	"math/big"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

type DynamicGPConfig struct {

	// Enabled defines if the dynamic gas price is enabled or disabled
	Enabled bool `mapstructure:"Enabled"`

	// CongestionTxThreshold defines the tx threshold to measure whether there is congestion
	CongestionTxThreshold uint64 `mapstructure:"CongestionTxThreshold"`

	// CheckBatches defines the number of recent Batches used to sample gas price
	CheckBatches int `mapstructure:"CheckBatches"`

	// SampleTxNumer defines the number of sampled gas prices in each batch
	SampleNumer int `mapstructure:"SampleTxNum"`

	// Percentile defines the sampling weight of all sampled gas prices
	Percentile int `mapstructure:"Percentile"`

	// MaxPrice defines the dynamic gas price upper limit
	MaxPrice uint64 `mapstructure:"MaxPrice"`

	// MinPrice defines the dynamic gas price lower limit
	MinPrice uint64 `mapstructure:"MinPrice"`
}

type DynamicGPManager struct {
	lastL2BatchNumber uint64
	lastPrice         *big.Int
	cacheLock         sync.RWMutex
	fetchLock         sync.Mutex
}

func (e *EthEndpoints) calcDynamicGP(ctx context.Context) {
	l2BatchNumber, err := e.state.GetLastBatchNumber(ctx, nil)
	if err != nil {
		log.Errorf("failed to get last l2 batch number, err: %v", err)
	}

	e.dgpMan.cacheLock.RLock()
	lastL2BatchNumber, lastPrice := e.dgpMan.lastL2BatchNumber, new(big.Int).Set(e.dgpMan.lastPrice)
	e.dgpMan.cacheLock.RUnlock()
	if l2BatchNumber == lastL2BatchNumber {
		log.Debug("Batch is still the same, no need to update the gas price at the moment, lastL2BatchNumber: ", lastL2BatchNumber)
		return
	}

	e.dgpMan.fetchLock.Lock()
	defer e.dgpMan.fetchLock.Unlock()

	var (
		sent, exp int
		number    = lastL2BatchNumber
		result    = make(chan results, e.cfg.DynamicGP.CheckBatches)
		quit      = make(chan struct{})
		results   []*big.Int
	)

	for sent < e.cfg.DynamicGP.CheckBatches && number > 0 {
		go e.getL2BatchTxsTips(ctx, number, e.cfg.DynamicGP.SampleNumer, result, quit)
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
		price = results[(len(results)-1)*e.cfg.DynamicGP.Percentile/100]
	}

	e.dgpMan.cacheLock.Lock()
	e.dgpMan.lastPrice = price
	e.dgpMan.lastL2BatchNumber = l2BatchNumber
	e.dgpMan.cacheLock.Unlock()
}

func (e *EthEndpoints) getL2BatchTxsTips(ctx context.Context, l2BlockNumber uint64, limit int, result chan results, quit chan struct{}) {
	txs, _, err := e.state.GetTransactionsByBatchNumber(ctx, l2BlockNumber, nil)
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
	var lowPrices []*big.Int
	var highPrices []*big.Int
	for _, tx := range sorter.txs {
		tip := tx.GasTipCap()

		lowPrices = append(lowPrices, tip)
		if len(lowPrices) >= limit {
			break
		}
	}

	sorter.Reverse()
	for _, tx := range sorter.txs {
		tip := tx.GasTipCap()

		highPrices = append(highPrices, tip)
		if len(highPrices) >= limit {
			break
		}
	}

	if len(highPrices) != len(lowPrices) {
		err := errors.New("len(highPrices) != len(lowPrices)")
		log.Errorf("getL2BlockTxsTips err: %v", err)
		select {
		case result <- results{nil, err}:
		case <-quit:
		}
		return
	}

	for i := 0; i < len(lowPrices); i++ {
		price := getAvgPrice(lowPrices[i], highPrices[i])
		prices = append(prices, price)
	}

	select {
	case result <- results{prices, nil}:
	case <-quit:
	}
}

func (e *EthEndpoints) isCongested(ctx context.Context) (bool, error) {
	txCount, err := e.pool.CountPendingTransactions(ctx)
	if err != nil {
		return false, err
	}
	if txCount >= e.cfg.DynamicGP.CongestionTxThreshold {
		return true, nil
	}
	return false, nil
}

type results struct {
	values []*big.Int
	err    error
}

type txSorter struct {
	txs []types.Transaction
}

func newSorter(txs []types.Transaction) *txSorter {
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

func (s *txSorter) Reverse() {
	for i := 0; i < len(s.txs)/2; i++ {
		j := len(s.txs) - i - 1
		s.txs[i], s.txs[j] = s.txs[j], s.txs[i]
	}
}

type bigIntArray []*big.Int

func (s bigIntArray) Len() int           { return len(s) }
func (s bigIntArray) Less(i, j int) bool { return s[i].Cmp(s[j]) < 0 }
func (s bigIntArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func getAvgPrice(low *big.Int, high *big.Int) *big.Int {
	avg := new(big.Int).Add(low, high)
	avg = avg.Quo(avg, big.NewInt(2)) //nolint:gomnd
	return avg
}
