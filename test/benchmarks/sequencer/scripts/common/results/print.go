package results

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
)

// Print prints the results of the benchmark
func Print(elapsed time.Duration) {
	// Print results
	log.Info("##########")
	log.Info("# Result #")
	log.Info("##########")
	log.Infof("Total time took for the sequencer to select all txs from the pool: %v", elapsed)
	log.Infof("Number of txs sent: %d", params.NumberOfOperations)
	log.Infof("Txs per second: %f", float64(params.NumberOfOperations)/elapsed.Seconds())
}

func PrintUniswapDeployments(deployments time.Duration, count uint64) {
	log.Info("#######################")
	log.Info("# Uniswap Deployments #")
	log.Info("#######################")
	log.Infof("Total time took for the sequencer to deploy all contracts: %v", deployments)
	log.Infof("Number of txs sent: %d", count)
	log.Infof("Txs per second: %f", float64(count)/deployments.Seconds())
}
