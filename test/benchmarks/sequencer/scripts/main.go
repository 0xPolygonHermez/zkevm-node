package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/metrics"
)

const (
	maxRetryAttempts = 5
	retryDelay       = 1 * time.Second
)

func main() {
	fmt.Println("Starting the program...")
	fmt.Println("-----------------------")

	// Command line flags
	tType := flag.String("type", "", "The type of transactions to test: erc20, uniswap, or eth.")
	numOps := flag.Int("num-ops", 200, "The number of operations to run. Default is 200.")
	help := flag.Bool("help", false, "Display help message")
	flag.Parse()

	if *help {
		fmt.Println("Usage: go run main.go --type TRANSACTIONS_TYPE --sequencer-ip SEQUENCER_IP [--num-ops NUMBER_OF_OPERATIONS]")
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

	fmt.Println("Checking environment variables...")
	fmt.Println("---------------------------------")
	// Check environment variables
	checkEnvVar("BASTION_HOST")
	checkEnvVar("POOLDB_PORT")
	checkEnvVar("POOLDB_EP")
	checkEnvVar("RPC_URL")
	checkEnvVar("CHAIN_ID")
	checkEnvVar("PRIVATE_KEY")
	checkEnvVar("SEQUENCER_IP")

	// Forward BASTION Ports
	fmt.Println("Forwarding BASTION ports...")
	fmt.Println("---------------------------")
	sshArgs := []string{"-fN",
		"-L", os.Getenv("POOLDB_PORT") + ":" + os.Getenv("POOLDB_EP") + ":5432",
		"ubuntu@" + os.Getenv("BASTION_HOST")}
	_, err = runCmd("ssh", sshArgs...)
	if err != nil {
		panic(fmt.Sprintf("Failed to forward BASTION ports: %v", err))
	}
	defer killSSHProcess(err)

	// Execute wget to get metrics from the BASTION HOST
	sequencerIP := os.Getenv("SEQUENCER_IP")
	fmt.Println("Fetching start metrics...")
	fmt.Println("--------------------------")

	output, err := retryCmd("ssh", "ubuntu@"+os.Getenv("BASTION_HOST"), "wget", "-qO-", "http://"+sequencerIP+":9091/metrics")
	if err != nil {
		panic(fmt.Sprintf("Failed to collect start metrics from BASTION HOST: %v", err))
	}
	retryTimes := 0
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to collect start metrics from BASTION HOST: %v", err))
		fmt.Println("Retrying...")
		time.Sleep(1 * time.Second)
		output, err = runCmd("ssh", "ubuntu@"+os.Getenv("BASTION_HOST"), "wget", "-qO-", "http://"+sequencerIP+":9091/metrics")
		retryTimes++
		if retryTimes == 5 {
			panic(fmt.Sprintf("Failed to collect start metrics from BASTION HOST: %v", err))
		}
	}

	err = os.WriteFile("start-metrics.txt", []byte(output), 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to write start metrics to file: %v", err))
	}

	// Run the Go script depending on the type argument
	var goScript string
	switch *tType {
	case "erc20":
		goScript = "erc20-transfers"
	case "uniswap":
		goScript = "uniswap-transfers"
	case "eth":
		goScript = "eth-transfers"
	}

	// Run transfers script
	fmt.Println("Running transfers script...")
	fmt.Println("---------------------------")
	lastLine, err := runCmdRealTime("go", "run", "./"+goScript+"/main.go", "--num-ops", strconv.Itoa(*numOps))
	if err != nil {
		panic(fmt.Sprintf("Failed to run Go script for %s transactions: %v", *tType, err))
	}

	// Extract Total Gas
	fmt.Println("Extracting Total Gas...")
	fmt.Println("-----------------------")
	var totalGas string
	if strings.Contains(lastLine, "Total Gas") {
		parts := strings.Split(lastLine, " ")
		totalGas = parts[len(parts)-1]
	}
	if totalGas == "" {
		fmt.Println("Warning: Failed to extract Total Gas from Go script output.")
	}

	// Execute wget to get metrics from the BASTION HOST
	fmt.Println("Fetching end metrics...")
	fmt.Println("------------------------")
	output, err = retryCmd("ssh", "ubuntu@"+os.Getenv("BASTION_HOST"), "wget", "-qO-", "http://"+sequencerIP+":9091/metrics")
	if err != nil {
		panic(fmt.Sprintf("Failed to collect end metrics from BASTION HOST: %v", err))
	}
	err = os.WriteFile("end-metrics.txt", []byte(output), 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to write end metrics to file: %v", err))
	}

	// Run the Go script that calculates the metrics and prints the results
	totalGasInt, err := strconv.ParseUint(totalGas, 10, 64)
	if err != nil {
		fmt.Printf("Failed to convert totalGas to int: %v\n", err)
	}

	// Calc and Print Results
	fmt.Println("Calculating and printing results...")
	fmt.Printf("------------------------------------\n\n")
	calculateAndPrintResults(*tType, totalGasInt, uint64(*numOps))

	fmt.Println("Done!")
}

func runCmd(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// runCmdWithRetry executes the specified command with arguments and returns the combined output.
// It includes a retryCmd mechanism controlled by the enableRetry flag.
func runCmdWithRetry(enableRetry bool, command string, args ...string) (string, error) {
	var output string
	var err error

	if enableRetry {
		for attempt := 1; attempt <= maxRetryAttempts; attempt++ {
			cmd := exec.Command(command, args...)
			cmd.Stderr = os.Stderr
			result, runErr := cmd.CombinedOutput()
			output = string(result)
			err = runErr

			if err == nil {
				// Command succeeded, no need to retryCmd.
				break
			}

			fmt.Printf("Attempt %d: Command failed: %v\n", attempt, err)

			if attempt < maxRetryAttempts {
				fmt.Println("Retrying...")
				time.Sleep(time.Second) // Add a delay between retries (you can adjust the duration).
			}
		}
	} else {
		cmd := exec.Command(command, args...)
		cmd.Stderr = os.Stderr
		result, runErr := cmd.CombinedOutput()
		output = string(result)
		err = runErr
	}

	return output, err
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
	_, err = runCmd("pkill", "-f", "ssh -fN -L "+os.Getenv("POOLDB_PORT"))
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

func retryCmd(command string, args ...string) (string, error) {

	for i := 0; i < maxRetryAttempts; i++ {
		result, err := runCmd(command, args...)
		if err == nil {
			return result, nil // If the function succeeded, return its result.
		}

		// If it failed and it's not the last attempt, wait for the specified delay before retrying.
		if i < maxRetryAttempts-1 {
			time.Sleep(retryDelay)
		}
	}

	return "", errors.New("maximum retryCmd attempts reached")
}
