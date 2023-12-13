package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var (
	showErrors bool
	showLosses bool
	showDetail bool
)

const (
	signatureBytes    = 65
	effectivePctBytes = 1
	fixedBytesTx      = signatureBytes + effectivePctBytes
)

type egpConfig struct {
	byteGasCost         uint64  // gas cost of 1 byte
	zeroGasCost         uint64  // gas cost of 1 byte zero
	netProfitFactor     float64 // L2 network profit factor
	l1GasPriceFactor    float64 // L1 gas price factor
	l2GasPriceSugFactor float64 // L2 gas price suggester factor
	breakEvenFactor     float64 // break even gas price factor
	finalDeviationPct   uint64  // max final deviation percentage
}

type egpLogRecord struct {
	encoded           string
	missingLogInfo    bool   // Flag if egp_log field is empty
	LogError          string `json:"Error"`
	LogEnabled        bool   `json:"Enabled"`
	LogL1GasPrice     uint64 `json:"L1GasPrice"`     // L1 gas price
	LogBalanceOC      bool   `json:"BalanceOC"`      // uses opcode to query balance
	LogGasPriceOC     bool   `json:"GasPriceOC"`     // uses opcode to query gas price
	LogGasUsedFirst   uint64 `json:"GasUsedFirst"`   // execute estimate gas
	LogGasUsedSecond  uint64 `json:"GasUsedSecond"`  // after execute gas
	LogL2GasPrice     uint64 `json:"L2GasPrice"`     // L2 gas price = LogL1GasPrice * l2GasPriceSugFactor
	LogGasPrice       uint64 `json:"GasPrice"`       // user gas price (signed) = L2 gas price
	LogValueFirst     uint64 `json:"ValueFirst"`     // effective gas price using LogGasUsedFirst (EGP)
	LogValueSecond    uint64 `json:"ValueSecond"`    // effective gas price using LogGasUsedSecond (NEGP)
	LogValueFinal     uint64 `json:"ValueFinal"`     // final gas price
	LogReprocess      bool   `json:"Reprocess"`      // reprocessed (executed 2 times)
	LogPercentage     uint64 `json:"Percentage"`     // user gas/final gas, coded percentage (0:not used, 1..255)
	LogMaxDeviation   uint64 `json:"MaxDeviation"`   // max allowed deviation = LogValueFirst * finalDeviationPct
	LogFinalDeviation uint64 `json:"FinalDeviation"` // final gas deviation = abs(LogValueSecond - LogValueFirst)
}

type egpStats struct {
	totalTx          uint64 // Analyzed tx count
	totalError       uint64 // EGP error tx count
	totalNoInfo      uint64 // Empty egp_log tx count
	totalEgp         uint64 // EGP enabled tx count
	totalReprocessed uint64 // Reprocessed tx count
	totalShady       uint64 // Suspicious tx count (used balance or gasprice opcodes)
	totalUsedFirst   uint64 // Used final gas is the first EGP computed
	totalUsedSecond  uint64 // Used final gas is the new EGP recomputed
	totalUsedUser    uint64 // Used final gas is the user gas price signed
	totalUsedWeird   uint64 // Used final gas is different from EGP, new EGP, and user
	totalLossCount   uint64 // Loss gas tx count
	totalLoss        uint64 // Total loss gas amount
}

func main() {
	// Create CLI app
	app := cli.NewApp()
	app.Usage = "Analyze stats for EGP"
	app.Flags = []cli.Flag{
		&cli.Uint64Flag{
			Name:  "from",
			Usage: "stats from L2 block onwards",
			Value: ^uint64(0),
		},
		&cli.Uint64Flag{
			Name:  "to",
			Usage: "stats until L2 block (optional)",
			Value: ^uint64(0),
		},
		&cli.BoolFlag{
			Name:  "showerror",
			Usage: "show transactions with EGP errors",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "showloss",
			Usage: "show transactions with losses",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "showdetail",
			Usage: "show full detail record when show loss/error",
			Value: false,
		},
		&cli.StringFlag{
			Name:     "cfg",
			Aliases:  []string{"c"},
			Usage:    "configuration file",
			Required: false,
		},
		&cli.StringFlag{
			Name:  "db",
			Usage: "DB connection string: \"host=xxx port=xxx user=xxx dbname=xxx password=xxx\"",
			Value: "",
		},
	}
	app.Action = runStats

	// Run CLI app
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// defaultConfig parses the default configuration values
func defaultConfig() (*egpConfig, error) {
	cfg := egpConfig{
		byteGasCost:         16,   // nolint:gomnd
		zeroGasCost:         4,    // nolint:gomnd
		netProfitFactor:     1,    // nolint:gomnd
		l1GasPriceFactor:    0.25, // nolint:gomnd
		l2GasPriceSugFactor: 0.5,  // nolint:gomnd
		breakEvenFactor:     1,    // nolint:gomnd
		finalDeviationPct:   10,   // nolint:gomnd
	}

	viper.SetConfigType("toml")
	return &cfg, nil
}

// loadConfig loads the configuration
func loadConfig(ctx *cli.Context) (*egpConfig, error) {
	cfg, err := defaultConfig()
	if err != nil {
		return nil, err
	}

	configFilePath := ctx.String("cfg")
	if configFilePath != "" {
		dirName, fileName := filepath.Split(configFilePath)

		fileExtension := strings.TrimPrefix(filepath.Ext(fileName), ".")
		fileNameWithoutExtension := strings.TrimSuffix(fileName, "."+fileExtension)

		viper.AddConfigPath(dirName)
		viper.SetConfigName(fileNameWithoutExtension)
		viper.SetConfigType(fileExtension)
	}

	err = viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			return nil, errors.New("config file not found")
		} else {
			return nil, err
		}
	}

	decodeHooks := []viper.DecoderConfigOption{
		// this allows arrays to be decoded from env var separated by ",", example: MY_VAR="value1,value2,value3"
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(mapstructure.TextUnmarshallerHookFunc(), mapstructure.StringToSliceHookFunc(","))),
	}

	err = viper.Unmarshal(&cfg, decodeHooks...)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// runStats calculates statistics based on EGP log
func runStats(ctx *cli.Context) error {
	// CLI params
	dbConn := ctx.String("db")
	if dbConn == "" {
		return errors.New("missing required parameter --db")
	} else {
		dbConn = dbConn + " sslmode=disable"
	}

	fromBlock := ctx.Uint64("from")
	if fromBlock == ^uint64(0) {
		// Default value if param not present
		fromBlock = 7950000 // nolint:gomnd
	}
	toBlock := ctx.Uint64("to")
	showErrors = ctx.Bool("showerror")
	showLosses = ctx.Bool("showloss")
	showDetail = ctx.Bool("showdetail")

	// Load simulation config file
	var err error
	var egpCfg *egpConfig
	if ctx.String("cfg") != "" {
		egpCfg, err = loadConfig(ctx)
		if err != nil {
			return err
		}
	}

	// Set DB connection
	config, err := pgx.ParseConfig(dbConn)
	if err != nil {
		logf("Error setting connection to db: %v", err)
		return err
	}

	// Connect to DB
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		logf("Error connecting to db: %v", err)
		return err
	}
	defer conn.Close(context.Background())

	// Query data
	query := fmt.Sprintf(`
		select t.l2_block_num, coalesce(t.egp_log::varchar,'') as egp_log, t.encoded
		from state.transaction t 
		where t.l2_block_num >= %d and t.l2_block_num <= %d`, fromBlock, toBlock)

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		logf("Error executing query: %v", err)
		return err
	}
	defer rows.Close()

	// Loop data rows
	logf("Starting from L2 block %d...", fromBlock)
	var l2Block uint64
	var egpLog, encoded string
	var stats, simulateStats egpStats

	i := uint64(0)
	for rows.Next() {
		// Fetch row
		err = rows.Scan(&l2Block, &egpLog, &encoded)
		if err != nil {
			logf("Error fetching row: %v", err)
			return err
		}

		// Work in progress
		if i%100000 == 0 {
			logf("Working txs #%d (L2 block %d)...", i, l2Block)
		}
		i++

		// Transaction info
		egpData := egpLogRecord{
			encoded:        encoded,
			missingLogInfo: egpLog == "",
		}

		// Check if EGP info is present
		if egpLog != "" {
			// Decode EGP log json
			err = json.Unmarshal([]byte(egpLog), &egpData)
			if err != nil {
				logf("Error decoding json from egp_log field: %v", err)
				return err
			}
		}

		// Calculate stats
		countStats(i, l2Block, &egpData, &stats)

		// Simulate using alternative config
		if egpCfg != nil {
			egpSimData := egpData
			simulateConfig(&egpSimData, egpCfg)
			countStats(i, l2Block, &egpSimData, &simulateStats)
		}
	}

	logf("Finished txs #%d (L2 block %d).", i, l2Block)

	// Print stats results
	logf("\nEGP STATS:")
	printStats(&stats)

	// Print simulation stats results
	if egpCfg != nil {
		logf("\nEGP SIMULATION STATS:")
		printStats(&simulateStats)
	}

	return nil
}

// countStats calculates and counts statistics for an EGP record
func countStats(i uint64, block uint64, egp *egpLogRecord, stats *egpStats) {
	printEgpLogRecord(egp, false)

	// Total transactions
	stats.totalTx++

	// Error transactions
	if egp.LogError != "" {
		stats.totalError++
		if showErrors {
			fmt.Printf("egp-error:#%d:(L2 block %d):%s\n", i, block, egp.LogError)
			if showDetail {
				printEgpLogRecord(egp, false)
			}
		}
	}

	// Field egp_log is empty
	if egp.missingLogInfo {
		stats.totalNoInfo++
	} else {
		// Analyze info
		if egp.LogReprocess {
			stats.totalReprocessed++

			// Suspicious
			if (egp.LogValueSecond < egp.LogGasPrice) && (egp.LogBalanceOC || egp.LogGasPriceOC) {
				stats.totalShady++
			}
		}

		if egp.LogEnabled {
			stats.totalEgp++
		}

		// Gas used
		if egp.LogValueFinal == egp.LogValueFirst {
			stats.totalUsedFirst++
		} else if egp.LogValueFinal == egp.LogValueSecond {
			stats.totalUsedSecond++
		} else if egp.LogValueFinal == egp.LogGasPrice {
			stats.totalUsedUser++
		} else {
			stats.totalUsedWeird++
		}

		// Loss
		if egp.LogValueFinal == egp.LogGasPrice {
			loss := uint64(0)
			if egp.LogReprocess {
				if (egp.LogGasUsedSecond-egp.LogValueFinal > 0) && (egp.LogValueFinal < egp.LogGasUsedSecond) {
					loss = egp.LogGasUsedSecond - egp.LogValueFinal
					stats.totalLossCount++
				}
			} else {
				if egp.LogGasUsedFirst-egp.LogValueFinal > 0 && (egp.LogValueFinal < egp.LogGasUsedFirst) {
					loss = egp.LogGasUsedFirst - egp.LogValueFinal
					stats.totalLossCount++
				}
			}
			stats.totalLoss += loss

			if showLosses {
				info := fmt.Sprintf("reprocess=%t, final=%d, egp1=%d, egp2=%d, user=%d", egp.LogReprocess, egp.LogValueFinal, egp.LogGasUsedFirst, egp.LogGasUsedSecond, egp.LogGasPrice)
				fmt.Printf("egp-loss:#%d:(L2 block %d):loss=%d:info:%s\n", i, block, loss, info)
				if showDetail {
					printEgpLogRecord(egp, false)
				}
			}
		}
	}
}

// logf prints log message
func logf(format string, args ...any) {
	msg := fmt.Sprintf(format+"\n", args...)
	fmt.Printf("%v", msg)
}

// printEgpLogRecord prints values of egpLogRecord struct
func printEgpLogRecord(record *egpLogRecord, showTxInfo bool) {
	fmt.Printf("  Error: [%s]\n", record.LogError)
	fmt.Printf("  Enabled: [%t]\n", record.LogEnabled)
	fmt.Printf("  L1GasPrice: [%d]\n", record.LogL1GasPrice)
	fmt.Printf("  BalanceOC: [%t]\n", record.LogBalanceOC)
	fmt.Printf("  GasPriceOC: [%t]\n", record.LogGasPriceOC)
	fmt.Printf("  GasUsedFirst: [%d]\n", record.LogGasUsedFirst)
	fmt.Printf("  GasUsedSecond: [%d]\n", record.LogGasUsedSecond)
	fmt.Printf("  L2GasPrice: [%d]\n", record.LogL2GasPrice)
	fmt.Printf("  GasPrice: [%d]\n", record.LogGasPrice)
	fmt.Printf("  ValueFirst: [%d]\n", record.LogValueFirst)
	fmt.Printf("  ValueSecond: [%d]\n", record.LogValueSecond)
	fmt.Printf("  ValueFinal: [%d]\n", record.LogValueFinal)
	fmt.Printf("  Reprocess: [%t]\n", record.LogReprocess)
	fmt.Printf("  Percentage: [%d]\n", record.LogPercentage)
	fmt.Printf("  MaxDeviation: [%d]\n", record.LogMaxDeviation)
	fmt.Printf("  FinalDeviation: [%d]\n", record.LogFinalDeviation)
	if showTxInfo {
		fmt.Printf("  encoded: [%s]\n", record.encoded)
	}
}

// printStats prints EGP statistics
func printStats(stats *egpStats) {
	fmt.Printf("Total Tx.........: [%d]\n", stats.totalTx)
	if stats.totalTx == 0 {
		return
	}

	fmt.Printf("Error Tx.........: [%d] (%.2f%%)\n", stats.totalError, float64(stats.totalError)/float64(stats.totalTx)*100)
	fmt.Printf("Total No EGP info: [%d] (%.2f%%)\n", stats.totalNoInfo, float64(stats.totalNoInfo)/float64(stats.totalTx)*100)

	statsCount := stats.totalTx - stats.totalNoInfo
	fmt.Printf("Total Tx EGP info: [%d] (%.2f%%)\n", statsCount, float64(statsCount)/float64(stats.totalTx)*100)
	if statsCount > 0 {
		fmt.Printf("    EGP enable.......: [%d] (%.2f%%)\n", stats.totalEgp, float64(stats.totalEgp)/float64(statsCount)*100)
		fmt.Printf("    Reprocessed Tx...: [%d] (%.2f%%)\n", stats.totalReprocessed, float64(stats.totalReprocessed)/float64(statsCount)*100)
		fmt.Printf("        Suspicious Tx....: [%d] (%.2f%%)\n", stats.totalShady, float64(stats.totalShady)/float64(stats.totalReprocessed)*100)
		fmt.Printf("    Final gas:\n")
		fmt.Printf("        Used EGP1........: [%d] (%.2f%%)\n", stats.totalUsedFirst, float64(stats.totalUsedFirst)/float64(statsCount)*100)
		fmt.Printf("        Used EGP2........: [%d] (%.2f%%)\n", stats.totalUsedSecond, float64(stats.totalUsedSecond)/float64(statsCount)*100)
		fmt.Printf("        Used User Gas....: [%d] (%.2f%%)\n", stats.totalUsedUser, float64(stats.totalUsedUser)/float64(statsCount)*100)
		fmt.Printf("        Used Weird Gas...: [%d] (%.2f%%)\n", stats.totalUsedWeird, float64(stats.totalUsedWeird)/float64(statsCount)*100)
		fmt.Printf("    Loss count.......: [%d] (%.2f%%)\n", stats.totalLossCount, float64(stats.totalLossCount)/float64(statsCount)*100)
		fmt.Printf("    Loss total.......: [%d]\n", stats.totalLoss)
		if stats.totalLossCount > 0 {
			fmt.Printf("    Loss average.....: [%d]\n", stats.totalLoss/stats.totalLossCount)
		}
	}
}

// simulateConfig simulates scenario using received config
func simulateConfig(egp *egpLogRecord, cfg *egpConfig) {
	// L2 and user gas price
	egp.LogL2GasPrice = uint64(float64(egp.LogL1GasPrice) * cfg.l2GasPriceSugFactor)
	egp.LogGasPrice = egp.LogL2GasPrice

	// Compute EGP
	egp.LogReprocess = false
	egp.LogValueFirst = uint64(calcEffectiveGasPrice(egp.LogGasUsedFirst, egp, cfg))

	if egp.LogValueFirst < egp.LogGasPrice {
		// Recompute NEGP
		egp.LogValueSecond = uint64(calcEffectiveGasPrice(egp.LogGasUsedSecond, egp, cfg))

		// Gas price deviation
		egp.LogFinalDeviation = uint64(math.Abs(float64(egp.LogValueSecond - egp.LogValueFirst)))
		egp.LogMaxDeviation = egp.LogValueFirst * cfg.finalDeviationPct / 100

		if egp.LogFinalDeviation < egp.LogMaxDeviation {
			// Final gas: EGP
			egp.LogValueFinal = egp.LogValueFirst
		} else {
			egp.LogReprocess = true
			if (egp.LogValueSecond < egp.LogGasPrice) || (egp.LogGasPriceOC || egp.LogBalanceOC) {
				// Final gas: price signed
				egp.LogValueFinal = egp.LogGasPrice
			} else {
				// Final gas: NEGP
				egp.LogValueFinal = egp.LogValueSecond
			}
		}
	} else {
		egp.LogValueSecond = 0

		// Final gas: price signed
		egp.LogValueFinal = egp.LogGasPrice
	}

	// Gas price effective percentage
	egp.LogPercentage = ((egp.LogValueFinal*256)+egp.LogGasPrice-1)/egp.LogGasPrice - 1
}

// calcEffectiveGasPrice calculates the effective gas price
func calcEffectiveGasPrice(gasUsed uint64, tx *egpLogRecord, cfg *egpConfig) float64 {
	// Calculate break even gas price
	var breakEvenGasPrice float64
	if gasUsed == 0 {
		breakEvenGasPrice = float64(tx.LogGasPrice)
	} else {
		// String 0x format to raw bytes
		rawBytes, err := hex.DecodeString(tx.encoded[2:])
		if err != nil {
			logf("Error converting encoded string to slice bytes: %v", err)
		}

		// Zero and non zero bytes
		txZeroBytes := uint64(bytes.Count(rawBytes, []byte{0}))
		txNonZeroBytes := uint64(len(rawBytes)) - txZeroBytes

		// Calculates break even gas price
		l2MinGasPrice := float64(tx.LogL1GasPrice) * cfg.l1GasPriceFactor
		totalTxPrice := float64(gasUsed)*l2MinGasPrice + float64(((fixedBytesTx+txNonZeroBytes)*cfg.byteGasCost+txZeroBytes*cfg.zeroGasCost)*tx.LogL1GasPrice)
		breakEvenGasPrice = totalTxPrice / float64(gasUsed) * cfg.netProfitFactor
	}

	// Calculate effective gas price
	var ratioPriority float64
	if gasUsed > tx.LogL2GasPrice {
		ratioPriority = float64(gasUsed/tx.LogL2GasPrice) - 1
	} else {
		ratioPriority = 0
	}
	effectiveGasPrice := breakEvenGasPrice * (1 + ratioPriority)

	// logf("zBytes=%d | nzBytes: %d | l2min: %f | txPrice: %f | breakEven: %f | gasPriceRPC: %f | prio: %f | EGP: %f",
	// 	txZeroBytes, txNonZeroBytes, l2MinGasPrice, totalTxPrice, breakEvenGasPrice, gasPriceRPC, ratioPriority, effectiveGasPrice)
	return effectiveGasPrice
}
