package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/metrics"
)

func main() {
	fmt.Println("Starting the program...")
	fmt.Println("-----------------------")

	// Command line flags
	tType := flag.String("type", "", "The type of transactions to test: erc20, uniswap, or eth.")
	sequencerIP := flag.String("sequencer-ip", "", "The IP address of the sequencer.")
	numOps := flag.Uint64("num-ops", 200, "The number of operations to run. Default is 200.")
	help := flag.Bool("help", false, "Display help message")
	flag.Parse()

	if *help {
		fmt.Println("Usage: go run exec_erc20_transfers.go --type TRANSACTIONS_TYPE --sequencer-ip SEQUENCER_IP [--num-ops NUMBER_OF_OPERATIONS]")
		flag.PrintDefaults()
		return
	}

	// Check if .env file exists
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		panic(fmt.Sprintf("Error: .env file does not exist. Please create it and add all environment variables from the Deployment Docs." +
			"\n ** check env.exmaple for an example. **"))
	}

	fmt.Println("Loading .env file...")
	fmt.Println("--------------------")
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}

	fmt.Println("Validating TYPE...")
	fmt.Println("------------------")
	// Validate TYPE
	if *tType == "" || (*tType != "erc20" && *tType != "uniswap" && *tType != "eth") {
		panic(fmt.Sprintf("Error: Invalid TYPE argument. Accepted values are 'erc20', 'uniswap', or 'eth'."))
	}

	fmt.Println("Validating SEQUENCER_IP...")
	fmt.Println("---------------------------")
	// Validate SEQUENCER_IP
	if *sequencerIP == "" {
		panic(fmt.Sprintf("Error: SEQUENCER_IP argument is missing. Please provide it as the second argument."))
	}

	fmt.Println("Checking environment variables...")
	fmt.Println("---------------------------------")
	// Check environment variables
	checkEnvVar("BASTION_HOST")
	checkEnvVar("POOLDB_LOCALPORT")
	checkEnvVar("POOLDB_EP")
	checkEnvVar("RPC_URL")
	checkEnvVar("CHAIN_ID")
	checkEnvVar("PRIVATE_KEY")

	// Forward BASTION Ports
	sshArgs := []string{"-fN",
		"-L", os.Getenv("POOLDB_LOCALPORT") + ":" + os.Getenv("POOLDB_EP") + ":5432",
		"ubuntu@" + os.Getenv("BASTION_HOST")}
	_, err = runCmd("ssh", sshArgs...)
	if err != nil {
		panic(fmt.Sprintf("Failed to forward BASTION ports: %v", err))
	}
	defer killSSHProcess(err)

	// ExecuteERC20Transfers wget to get metrics from the BASTION HOST
	fmt.Println("Fetching start metrics...")
	fmt.Println("--------------------------")
	output, err := runCmd("ssh", "ubuntu@"+os.Getenv("BASTION_HOST"), "wget", "-qO-", "http://"+*sequencerIP+":9091/metrics")
	if err != nil {
		panic(fmt.Sprintf("Failed to collect start metrics from BASTION HOST: %v", err))
	}
	err = os.WriteFile("start-metrics.txt", []byte(output), 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to write start metrics to file: %v", err))
	}

	// Run transfers script
	fmt.Println("Running transfers script...")
	fmt.Println("---------------------------")
	var totalGas uint64
	switch *tType {
	case "erc20":
		totalGas = ExecuteERC20Transfers(*numOps)
	case "uniswap":
		totalGas = ExecuteUniswapTransfers(*numOps)
	case "eth":
		totalGas = ExecuteEthTransfers(*numOps)
	}

	// ExecuteERC20Transfers wget to get metrics from the BASTION HOST
	fmt.Println("Fetching end metrics...")
	fmt.Println("------------------------")
	output, err = runCmd("ssh", "ubuntu@"+os.Getenv("BASTION_HOST"), "wget", "-qO-", "http://"+*sequencerIP+":9091/metrics")
	if err != nil {
		panic(fmt.Sprintf("Failed to collect end metrics from BASTION HOST: %v", err))
	}
	err = os.WriteFile("end-metrics.txt", []byte(output), 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to write end metrics to file: %v", err))
	}

	// Calc and Print Results
	fmt.Println("Calculating and printing results...")
	fmt.Printf("------------------------------------\n\n")
	calculateAndPrintResults(*tType, totalGas, *numOps)

	fmt.Println("Done!")
}

func runCmd(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func runCmdRealTime(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	cmd.Start()

	var lastLine string
	go func() {
		scanner := bufio.NewScanner(stdoutIn)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			lastLine = line
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrIn)
		for scanner.Scan() {
			m := scanner.Text()
			_, err := fmt.Fprintln(os.Stderr, m)
			if err != nil {
				fmt.Println("Error printing stderr: ", err)
				return
			}
		}
	}()

	err := cmd.Wait()
	if err != nil {
		return "", err
	}
	return lastLine, nil
}

func checkEnvVar(varName string) {
	if os.Getenv(varName) == "" {
		panic(fmt.Sprintf("Error: %s is not set. Please export all environment variables from the Deployment Docs.", varName))
	}
}

func killSSHProcess(err error) {
	fmt.Println("Killing SSH process...")
	_, err = runCmd("pkill", "-f", "ssh -fN -L "+os.Getenv("POOLDB_LOCALPORT"))
	if err != nil {
		panic(fmt.Sprintf("Failed to kill the SSH process: %v", err))
	}
}

func calculateAndPrintResults(txsType string, totalGas uint64, numberOfOperations uint64) {
	totalTransactionsSent := numberOfOperations

	startData := parseFile("start-metrics.txt")
	endData := parseFile("end-metrics.txt")

	totalTxs := uint64(endData["sequencer_processing_time"].processingTimeCount - startData["sequencer_processing_time"].processingTimeCount)

	processingTimeSequencer := endData["sequencer_processing_time"].processingTimeSum - startData["sequencer_processing_time"].processingTimeSum
	processingTimeExecutor := endData["state_executor_processing_time{caller=\"sequencer\"}"].processingTimeSum - startData["state_executor_processing_time{caller=\"sequencer\"}"].processingTimeSum

	fmt.Println("########################")
	fmt.Println("#        Results       #")
	fmt.Printf("########################\n\n")

	metrics.PrintSummary(
		txsType,
		totalTransactionsSent,
		totalTxs,
		processingTimeSequencer,
		processingTimeExecutor,
		totalGas,
	)
}

type timeData struct {
	processingTimeSum   float64
	processingTimeCount int
}

func parseLine(line string) (key string, value float64) {
	parts := strings.Split(line, " ")
	key = parts[0]
	value, _ = strconv.ParseFloat(parts[1], 64)
	return
}

func parseFile(filename string) map[string]timeData {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	result := map[string]timeData{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		key, value := parseLine(line)
		if strings.Contains(key, "sum") {
			key = strings.Replace(key, "_sum", "", -1)
			if data, ok := result[key]; ok {
				data.processingTimeSum = value
				result[key] = data
			} else {
				result[key] = timeData{processingTimeSum: value}
			}
		} else if strings.Contains(key, "count") {
			key = strings.Replace(key, "_count", "", -1)
			if data, ok := result[key]; ok {
				data.processingTimeCount = int(value)
				result[key] = data
			} else {
				result[key] = timeData{processingTimeCount: int(value)}
			}
		}
	}

	return result
}
