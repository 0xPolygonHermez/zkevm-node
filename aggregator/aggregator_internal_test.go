package aggregator

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestCheckEthTxNotFound(t *testing.T) {
	intervalAfterWhichBatchConsolidateAnyway := new(Duration)
	_ = intervalAfterWhichBatchConsolidateAnyway.UnmarshalText([]byte("1s"))

	cfg := Config{
		IntervalAfterWhichBatchConsolidateAnyway: *intervalAfterWhichBatchConsolidateAnyway,
		TxProfitabilityCheckerType:               ProfitabilityAcceptAll,
		TxProfitabilityMinReward:                 TokenAmountWithDecimals{},
	}
	aggr, err := NewAggregator(cfg, nil, nil, nil)
	assert.NoError(t, err)
	hash := common.HexToHash("0x125714bb4db48757007fff2671b37637bbfd6d47b3a4757ebbd0c5222984f905")
	aggr.batchesSent[1] = &hash
	aggr.checkEthTxNotFound(hash, 1)

	assert.Equal(t, hash.Hex(), aggr.batchesSent[1].Hex())
	assert.Equal(t, 1, aggr.txNotFoundCounter)

	for i := 1; i <= 10; i++ {
		aggr.checkEthTxNotFound(hash, 1)
	}

	assert.Equal(t, 0, aggr.txNotFoundCounter)
	assert.Nil(t, aggr.batchesSent[1])
}
