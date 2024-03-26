package sequencer

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	pmetric "github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
)

var countinterval = 10

func (s *Sequencer) countPendingTx() {
	for {
		<-time.After(time.Second * time.Duration(countinterval))
		transactions, err := s.pool.CountPendingTransactions(context.Background())
		if err != nil {
			log.Errorf("load pending tx from pool: %v", err)
			continue
		}
		pmetric.PendingTxCount(int(transactions))
	}
}
