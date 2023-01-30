package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	metricsLib "github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	metricsState "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/shared"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
)

const (
	oneHundred = 100
)

// CalculateAndPrint calculates and prints the results
func CalculateAndPrint(response *http.Response, elapsed time.Duration, sequencerTimeSub, executorTimeSub float64, nTxs int) {
	seqeuncerTime, executorTime, workerTime, err := GetValues(response)
	if err != nil {
		log.Fatalf("error getting prometheus metrics: %v", err)
	}

	log.Info("##########")
	log.Info("# Result #")
	log.Info("##########")
	log.Infof("Total time took for the sequencer To select all txs from the pool: %v", elapsed)
	log.Info("######################")
	log.Info("# Prometheus Metrics #")
	log.Info("######################")
	actualTotalTime := seqeuncerTime - sequencerTimeSub
	actualExecutorTime := executorTime - executorTimeSub
	Print(actualTotalTime, actualExecutorTime, workerTime)
	log.Infof("[Transactions per second]: %v", nTxs/int(actualTotalTime))
}

// Print prints the prometheus metrics
func Print(totalTime float64, executorTime float64, workerTime float64) {
	log.Infof("[TOTAL Processing Time]: %v s", totalTime)
	log.Infof("[EXECUTOR Processing Time]: %v s", executorTime)
	log.Infof("[SEQUENCER Processing Time]: %v s", totalTime-executorTime)
	log.Infof("[WORKER Processing Time]: %v s", workerTime)
	log.Infof("[EXECUTOR Time Percentage from TOTAL]: %.2f %%", (executorTime/totalTime)*oneHundred)
	log.Infof("[WORKER Time Percentage from TOTAL]: %.2f %%", (workerTime/totalTime)*oneHundred)
}

// GetValues gets the prometheus metric values
func GetValues(metricsResponse *http.Response) (float64, float64, float64, error) {
	var err error
	if metricsResponse == nil {
		metricsResponse, err = Fetch()
		if err != nil {
			log.Fatalf("error getting prometheus metrics: %v", err)
		}
	}

	mf, err := testutils.ParseMetricFamilies(metricsResponse.Body)
	if err != nil {
		return 0, 0, 0, err
	}
	sequencerTotalProcessingTimeHisto := mf[metrics.ProcessingTimeName].Metric[0].Histogram
	sequencerTotalProcessingTime := sequencerTotalProcessingTimeHisto.GetSampleSum()

	workerTotalProcessingTimeHisto := mf[metrics.WorkerProcessingTimeName].Metric[0].Histogram
	workerTotalProcessingTime := workerTotalProcessingTimeHisto.GetSampleSum()

	executorTotalProcessingTimeHisto := mf[metricsState.ExecutorProcessingTimeName].Metric[0].Histogram
	executorTotalProcessingTime := executorTotalProcessingTimeHisto.GetSampleSum()
	return sequencerTotalProcessingTime, executorTotalProcessingTime, workerTotalProcessingTime, nil
}

// Fetch fetches the prometheus metrics
func Fetch() (*http.Response, error) {
	return http.Get(fmt.Sprintf("http://localhost:%d%s", shared.PrometheusPort, metricsLib.Endpoint))
}
