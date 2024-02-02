package sequencer

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
)

var countinterval = 10

func (d *dbManager) countPendingTx() {
	for {
		<-time.After(time.Second * time.Duration(countinterval))
		transactions, err := d.txPool.CountPendingTransactions(d.ctx)
		if err != nil {
			log.Errorf("load pending tx from pool: %v", err)
			continue
		}
		metrics.PendingTxCount(int(transactions))
	}
}
