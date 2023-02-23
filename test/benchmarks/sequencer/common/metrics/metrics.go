package metrics

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	metricsLib "github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	metricsState "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
)

const (
	oneHundred    = 100
	profilingPort = 6060
)

// CalculateAndPrint calculates and prints the results
func CalculateAndPrint(prometheusResp *http.Response, profilingResult string, elapsed time.Duration, sequencerTimeSub, executorTimeSub float64, nTxs int) {
	var (
		sequencerTime, executorTime, workerTime float64
		err                                     error
	)
	if prometheusResp != nil {
		sequencerTime, executorTime, workerTime, err = GetValues(prometheusResp)
		if err != nil {
			log.Fatalf("error getting prometheus metrics: %v", err)
		}
	}

	log.Info("##########")
	log.Info("# Result #")
	log.Info("##########")
	log.Infof("Total time (including setup of environment and starting containers): %v", elapsed)

	if prometheusResp != nil {
		log.Info("######################")
		log.Info("# Prometheus Metrics #")
		log.Info("######################")
		actualTotalTime := sequencerTime - sequencerTimeSub
		actualExecutorTime := executorTime - executorTimeSub
		PrintPrometheus(actualTotalTime, actualExecutorTime, workerTime)
		log.Infof("[Transactions per second]: %v", float64(nTxs)/actualTotalTime)
	}
	if profilingResult != "" {
		log.Info("#####################")
		log.Info("# Profiling Metrics #")
		log.Info("#####################")
		log.Infof("%v", profilingResult)
	}
}

// PrintPrometheus prints the prometheus metrics
func PrintPrometheus(totalTime float64, executorTime float64, workerTime float64) {
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
		metricsResponse, err = FetchPrometheus()
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

// FetchPrometheus fetches the prometheus metrics
func FetchPrometheus() (*http.Response, error) {
	log.Infof("Fetching prometheus metrics ...")
	return http.Get(fmt.Sprintf("http://localhost:%d%s", params.PrometheusPort, metricsLib.Endpoint))
}

// FetchProfiling fetches the profiling metrics
func FetchProfiling() (string, error) {
	fullUrl := fmt.Sprintf("http://localhost:%d%s", profilingPort, metricsLib.ProfileEndpoint)
	log.Infof("Fetching profiling metrics from: %s ...", fullUrl)
	cmd := exec.Command("go", "tool", "pprof", "-show=sequencer", "-top", fullUrl)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running pprof: %v\n%s", err, out)
	}
	return string(out), err
}
