package results

import (
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Print prints the results of the benchmark
func Print(client *ethclient.Client, elapsed time.Duration, txs []*types.Transaction) {
	// calculate the total gas used
	var totalGas uint64
	for _, tx := range txs {
		// Fetch the transaction receipt
		receipt, err := client.TransactionReceipt(params.Ctx, tx.Hash())
		if err != nil {
			log.Error("Unable to fetch transaction receipt", "error", err)
			continue
		}

		totalGas += receipt.GasUsed
	}

	// calculate the average gas used per transaction
	var avgGas uint64
	if len(txs) > 0 {
		avgGas = totalGas / uint64(len(txs))
	}

	// calculate the gas per second
	gasPerSecond := float64(len(txs)*int(avgGas)) / elapsed.Seconds()

	// Print results
	fmt.Println("##########")
	fmt.Println("# Result #")
	fmt.Println("##########")
	fmt.Printf("Total time took for the sequencer to select all txs from the pool: %v\n", elapsed)
	fmt.Printf("Number of operations sent: %d\n", params.NumberOfOperations)
	fmt.Printf("Txs per second: %f\n", float64(params.NumberOfOperations)/elapsed.Seconds())
	fmt.Printf("Average gas used per transaction: %d\n", avgGas)
	fmt.Printf("Total Gas: %d\n", totalGas)
	fmt.Printf("Gas per second: %f\n", gasPerSecond)
}

func PrintUniswapDeployments(deployments time.Duration, count uint64) {
	fmt.Println("#######################")
	fmt.Println("# Uniswap Deployments #")
	fmt.Println("#######################")
	fmt.Printf("Total time took for the sequencer to deploy all contracts: %v\n", deployments)
	fmt.Printf("Number of txs sent: %d\n", count)
	fmt.Printf("Txs per second: %f\n", float64(count)/deployments.Seconds())
}
