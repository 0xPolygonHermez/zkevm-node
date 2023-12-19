package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
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
	ByteGasCost               uint64  // gas cost of 1 byte
	ZeroGasCost               uint64  // gas cost of 1 byte zero
	NetProfitFactor           float64 // L2 network profit factor
	L1GasPriceFactor          float64 // L1 gas price factor
	L2GasPriceSugFactor       float64 // L2 gas price suggester factor
	FinalDeviationPct         uint64  // max final deviation percentage
	MinGasPriceAllowed        uint64  // min gas price allowed
	L2GasPriceSugFactorPreEGP float64 // L2 gas price suggester factor (pre EGP)
}

type egpLogRecord struct {
	l2BlockNum        uint64
	l2BlockReceived   time.Time
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
	sumGasFinal      uint64 // Accumulated sum of final gas (to get average)
	countGasFinal    uint64 // Count number of accumulated (to get average)
	sumGasPreEGP     uint64 // Accumulated sum of gas without EGP
	countGasPreEGP   uint64 // Count number of accumulated pre EGP (to get average)
	sumFee           uint64
	sumFeePreEGP     uint64
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
		ByteGasCost:               16,         // nolint:gomnd
		ZeroGasCost:               4,          // nolint:gomnd
		NetProfitFactor:           1.0,        // nolint:gomnd
		L1GasPriceFactor:          0.25,       // nolint:gomnd
		L2GasPriceSugFactor:       0.5,        // nolint:gomnd
		FinalDeviationPct:         10,         // nolint:gomnd
		MinGasPriceAllowed:        1000000000, // nolint:gomnd
		L2GasPriceSugFactorPreEGP: 0.1,        // nolint:gomnd
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
		fromBlock = 7920355 // nolint:gomnd
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
		select lb.received_at, t.l2_block_num, coalesce(t.egp_log::varchar,'') as egp_log, t.encoded
		from state.transaction t 
			join state.l2block lb on lb.block_num = t.l2_block_num 
		where t.l2_block_num >= %d and t.l2_block_num <= %d`, fromBlock, toBlock)

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		logf("Error executing query: %v", err)
		return err
	}
	defer rows.Close()

	// Loop data rows
	logf("Starting from L2 block %d...", fromBlock)
	var blockReceived time.Time
	var l2Block uint64
	var egpLog, encoded string
	var stats, simulateStats egpStats
	var timeFirst, timeLast time.Time

	i := uint64(0)
	for rows.Next() {
		// Fetch row
		err = rows.Scan(&blockReceived, &l2Block, &egpLog, &encoded)
		if err != nil {
			logf("Error fetching row: %v", err)
			return err
		}

		// First and last txs timestamp
		if i == 0 {
			timeFirst = blockReceived
			timeLast = blockReceived
		}
		if blockReceived.Before(timeFirst) {
			timeFirst = blockReceived
		}
		if blockReceived.After(timeLast) {
			timeLast = blockReceived
		}

		// Work in progress
		if i%100000 == 0 {
			logf("Working txs #%d (L2 block [%d] %v)...", i, l2Block, blockReceived)
		}
		i++

		// Transaction info
		egpData := egpLogRecord{
			l2BlockReceived: blockReceived,
			l2BlockNum:      l2Block,
			encoded:         encoded,
			missingLogInfo:  egpLog == "",
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
		countStats(i, l2Block, &egpData, &stats, nil)

		// Simulate using alternative config
		if egpCfg != nil {
			egpSimData := egpData
			simulateConfig(&egpSimData, egpCfg)
			countStats(i, l2Block, &egpSimData, &simulateStats, egpCfg)
		}
	}

	logf("Finished txs #%d (L2 block [%d] %v).", i, l2Block, blockReceived)

	// Print stats results
	diff := timeLast.Sub(timeFirst).Hours()
	logf("\nPERIOD [%.2f days]: %v ... %v", diff/24, timeFirst, timeLast)
	logf("EGP REAL STATS:")
	printStats(&stats)

	// Print simulation stats results
	if egpCfg != nil {
		logf("\nEGP SIMULATION STATS:")
		printStats(&simulateStats)
		logf("PARAMS: byte[%d] zero[%d] netFactor[%.2f] L1factor[%.2f] L2sugFactor[%.2f] devPct[%d] minGas[%d] L2sugPreEGP[%.2f]", egpCfg.ByteGasCost,
			egpCfg.ZeroGasCost, egpCfg.NetProfitFactor, egpCfg.L1GasPriceFactor, egpCfg.L2GasPriceSugFactor, egpCfg.FinalDeviationPct, egpCfg.MinGasPriceAllowed, egpCfg.L2GasPriceSugFactorPreEGP)
	}

	return nil
}

// max calculates the maximum between 2 numbers
// func max(a uint64, b uint64) uint64 {
// 	if a > b {
// 		return a
// 	} else {
// 		return b
// 	}
// }

// countStats calculates and counts statistics for an EGP record
func countStats(i uint64, block uint64, egp *egpLogRecord, stats *egpStats, cfg *egpConfig) {
	// printEgpLogRecord(egp, false)
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

		// Tx Fee
		stats.sumFee += egp.LogValueFinal * egp.LogGasUsedSecond

		// Gas total and average
		stats.countGasFinal++
		stats.sumGasFinal += egp.LogValueFinal

		// Gas total and average without EGP
		var l2SugPreEGP float64
		if cfg != nil {
			l2SugPreEGP = cfg.L2GasPriceSugFactorPreEGP
		} else {
			l2SugPreEGP = 0.1
		}

		stats.countGasPreEGP++
		stats.sumGasPreEGP += uint64(float64(egp.LogL1GasPrice) * l2SugPreEGP)
		stats.sumFeePreEGP += uint64(float64(egp.LogL1GasPrice) * l2SugPreEGP * float64(egp.LogGasUsedSecond))

		// Loss
		if egp.LogValueFinal == egp.LogGasPrice {
			loss := uint64(0)
			if egp.LogReprocess {
				if (egp.LogValueSecond-egp.LogValueFinal > 0) && (egp.LogValueFinal < egp.LogValueSecond) {
					loss = egp.LogValueSecond - egp.LogValueFinal
					stats.totalLossCount++
				}
			} else {
				if egp.LogValueFirst-egp.LogValueFinal > 0 && (egp.LogValueFinal < egp.LogValueFirst) {
					loss = egp.LogValueFirst - egp.LogValueFinal
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
	fmt.Printf("L2BlockNum: [%d]\n", record.l2BlockNum)
	fmt.Printf("  timestamp: [%v]\n", record.l2BlockReceived)
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
	if record.LogReprocess {
		fmt.Printf("block %d reprocessed!", record.l2BlockNum)
	}
	fmt.Println()
}

// printStats prints EGP statistics
func printStats(stats *egpStats) {
	const (
		GWEI_DIV = 1000000000
		ETH_DIV  = 1000000000000000000
	)

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
		if stats.totalReprocessed > 0 {
			fmt.Printf("        Suspicious Tx....: [%d] (%.2f%%)\n", stats.totalShady, float64(stats.totalShady)/float64(stats.totalReprocessed)*100)
		} else {
			fmt.Printf("        Suspicious Tx....: [%d] (0.00%%)\n", stats.totalShady)
		}
		fmt.Printf("    Final gas:\n")
		fmt.Printf("        Used EGP1........: [%d] (%.2f%%)\n", stats.totalUsedFirst, float64(stats.totalUsedFirst)/float64(statsCount)*100)
		fmt.Printf("        Used EGP2........: [%d] (%.2f%%)\n", stats.totalUsedSecond, float64(stats.totalUsedSecond)/float64(statsCount)*100)
		fmt.Printf("        Used User Gas....: [%d] (%.2f%%)\n", stats.totalUsedUser, float64(stats.totalUsedUser)/float64(statsCount)*100)
		fmt.Printf("        Used Weird Gas...: [%d] (%.2f%%)\n", stats.totalUsedWeird, float64(stats.totalUsedWeird)/float64(statsCount)*100)
		if stats.countGasFinal > 0 {
			fmt.Printf("    Gas price avg........: [%d] (%.3f GWei) (%.9f ETH)\n", stats.sumGasFinal/stats.countGasFinal,
				float64(stats.sumGasFinal/stats.countGasFinal)/GWEI_DIV, float64(stats.sumGasFinal/stats.countGasFinal)/ETH_DIV)
		}
		if stats.countGasFinal > 0 {
			fmt.Printf("    Tx fee avg...........: [%d] (%.3f GWei) (%.9f ETH)\n", stats.sumFee/stats.countGasFinal,
				float64(stats.sumFee/stats.countGasFinal)/GWEI_DIV, float64(stats.sumFee/stats.countGasFinal)/ETH_DIV)
		}
		if stats.countGasPreEGP > 0 {
			fmt.Printf("    Gas pri.avg preEGP...: [%d] (%.3f GWei) (%.9f ETH)\n", stats.sumGasPreEGP/stats.countGasPreEGP,
				float64(stats.sumGasPreEGP/stats.countGasPreEGP)/GWEI_DIV, float64(stats.sumGasPreEGP/stats.countGasPreEGP)/ETH_DIV)
		}
		if stats.countGasPreEGP > 0 {
			fmt.Printf("    Tx fee avg preEGP....: [%d] (%.3f GWei) (%.9f ETH)\n", stats.sumFeePreEGP/stats.countGasPreEGP,
				float64(stats.sumFeePreEGP/stats.countGasPreEGP)/GWEI_DIV, float64(stats.sumFeePreEGP/stats.countGasPreEGP)/ETH_DIV)
		}
		fmt.Printf("    Diff fee EGP-preEGP..: [%d] (%.3f Gwei) (%.9f ETH)\n", int64(stats.sumFee-stats.sumFeePreEGP),
			float64(int64(stats.sumFee-stats.sumFeePreEGP))/GWEI_DIV, float64(int64(stats.sumFee-stats.sumFeePreEGP))/ETH_DIV)
		fmt.Printf("    Loss count.......: [%d] (%.2f%%)\n", stats.totalLossCount, float64(stats.totalLossCount)/float64(statsCount)*100)
		if stats.totalLoss < GWEI_DIV*10 {
			fmt.Printf("    Loss total.......: [%d] (%d KWei)\n", stats.totalLoss, stats.totalLoss/1000)
		} else {
			fmt.Printf("    Loss total.......: [%d] (%d GWei) (%.9f ETH)\n", stats.totalLoss, stats.totalLoss/GWEI_DIV, float64(stats.totalLoss)/ETH_DIV)
		}
		if stats.totalLossCount > 0 {
			if stats.totalLoss/stats.totalLossCount < GWEI_DIV*10 {
				fmt.Printf("    Loss average.....: [%d] (%d KWei)\n", stats.totalLoss/stats.totalLossCount, stats.totalLoss/stats.totalLossCount/1000)
			} else {
				fmt.Printf("    Loss average.....: [%d] (%d GWei) (%.9f ETH)\n", stats.totalLoss/stats.totalLossCount, stats.totalLoss/stats.totalLossCount/GWEI_DIV,
					float64(stats.totalLoss/stats.totalLossCount)/ETH_DIV)
			}
		}
	}
}

// simulateConfig simulates scenario using received config
func simulateConfig(egp *egpLogRecord, cfg *egpConfig) {
	// L2 and user gas price
	egp.LogL2GasPrice = uint64(float64(egp.LogL1GasPrice) * cfg.L2GasPriceSugFactor)
	egp.LogGasPrice = egp.LogL2GasPrice

	// Compute EGP
	var err error
	egp.LogReprocess = false
	egp.LogValueFirst, err = calcEffectiveGasPrice(egp.LogGasUsedFirst, egp, cfg)
	if err != nil {
		logf("Simulation error in L2 block [%d], EGP failed, error: %v", egp.l2BlockNum, err)
		os.Exit(1)
	}

	if egp.LogValueFirst < egp.LogGasPrice {
		// Recompute NEGP
		egp.LogValueSecond, err = calcEffectiveGasPrice(egp.LogGasUsedSecond, egp, cfg)
		if err != nil {
			logf("Simulation error in L2 block [%d], NEGP failed, error: %v", egp.l2BlockNum, err)
			os.Exit(2)
		}

		// Gas price deviation
		egp.LogFinalDeviation = uint64(math.Abs(float64(int64(egp.LogValueSecond) - int64(egp.LogValueFirst))))
		egp.LogMaxDeviation = egp.LogValueFirst * cfg.FinalDeviationPct / 100

		if egp.LogFinalDeviation < egp.LogMaxDeviation {
			// Final gas: EGP
			egp.LogValueFinal = egp.LogValueFirst
		} else {
			egp.LogReprocess = true
			if (egp.LogValueSecond < egp.LogGasPrice) && !egp.LogGasPriceOC && !egp.LogBalanceOC {
				// Final gas: NEGP
				egp.LogValueFinal = egp.LogValueSecond
			} else {
				// Final gas: price signed
				egp.LogValueFinal = egp.LogGasPrice
			}
		}
	} else {
		egp.LogValueSecond = 0

		// Final gas: price signed
		egp.LogValueFinal = egp.LogGasPrice
	}

	// Gas price effective percentage
	if egp.LogGasPrice > 0 {
		egp.LogPercentage = ((egp.LogValueFinal*256)+egp.LogGasPrice-1)/egp.LogGasPrice - 1
	} else {
		egp.LogPercentage = 0
	}
}

// calcEffectiveGasPrice calculates the effective gas price
func calcEffectiveGasPrice(gasUsed uint64, tx *egpLogRecord, cfg *egpConfig) (uint64, error) {
	// Calculate break even gas price
	var breakEvenGasPrice float64
	if gasUsed == 0 {
		breakEvenGasPrice = float64(tx.LogGasPrice)
	} else {
		// Decode tx
		rawBytes, err := decodeTx(tx)
		if err != nil {
			return 0, err
		}

		// Zero and non zero bytes
		txZeroBytes := uint64(bytes.Count(rawBytes, []byte{0}))
		txNonZeroBytes := uint64(len(rawBytes)) - txZeroBytes
		// logf("size: %d", len(rawBytes))

		// Calculates break even gas price
		l2MinGasPrice := float64(tx.LogL1GasPrice) * cfg.L1GasPriceFactor
		if l2MinGasPrice < float64(cfg.MinGasPriceAllowed) {
			l2MinGasPrice = float64(cfg.MinGasPriceAllowed)
		}
		totalTxPrice := float64(gasUsed)*l2MinGasPrice + float64(((fixedBytesTx+txNonZeroBytes)*cfg.ByteGasCost+txZeroBytes*cfg.ZeroGasCost)*tx.LogL1GasPrice)
		breakEvenGasPrice = totalTxPrice / float64(gasUsed) * cfg.NetProfitFactor
	}

	// Calculate effective gas price
	var ratioPriority float64
	if tx.LogGasPrice > tx.LogL2GasPrice {
		ratioPriority = math.Round(float64(tx.LogGasPrice / tx.LogL2GasPrice))
	} else {
		ratioPriority = 1
	}
	effectiveGasPrice := breakEvenGasPrice * ratioPriority

	return uint64(effectiveGasPrice), nil
}

// decodeTx decodes the encoded tx
func decodeTx(record *egpLogRecord) ([]byte, error) {
	tx, err := state.DecodeTx(record.encoded)
	if err != nil {
		return nil, err
	}

	binaryTx, err := prepareRLPTxData(*tx)
	if err != nil {
		return nil, err
	}

	return binaryTx, nil
}

// prepareRLPTxData prepares RLP raw transaction data
func prepareRLPTxData(tx types.Transaction) ([]byte, error) {
	const ether155V = 27

	v, r, s := tx.RawSignatureValues()
	sign := 1 - (v.Uint64() & 1)

	nonce, gasPrice, gas, to, value, data, chainID := tx.Nonce(), tx.GasPrice(), tx.Gas(), tx.To(), tx.Value(), tx.Data(), tx.ChainId()

	rlpFieldsToEncode := []interface{}{
		nonce,
		gasPrice,
		gas,
		to,
		value,
		data,
	}

	if !IsPreEIP155Tx(tx) {
		rlpFieldsToEncode = append(rlpFieldsToEncode, chainID)
		rlpFieldsToEncode = append(rlpFieldsToEncode, uint(0))
		rlpFieldsToEncode = append(rlpFieldsToEncode, uint(0))
	}

	txCodedRlp, err := rlp.EncodeToBytes(rlpFieldsToEncode)
	if err != nil {
		return nil, err
	}

	newV := new(big.Int).Add(big.NewInt(ether155V), big.NewInt(int64(sign)))
	newRPadded := fmt.Sprintf("%064s", r.Text(hex.Base))
	newSPadded := fmt.Sprintf("%064s", s.Text(hex.Base))
	newVPadded := fmt.Sprintf("%02s", newV.Text(hex.Base))
	txData, err := hex.DecodeString(hex.EncodeToString(txCodedRlp) + newRPadded + newSPadded + newVPadded)
	if err != nil {
		return nil, err
	}
	return txData, nil
}

// IsPreEIP155Tx checks if tx is previous EIP155
func IsPreEIP155Tx(tx types.Transaction) bool {
	v, _, _ := tx.RawSignatureValues()
	return tx.ChainId().Uint64() == 0 && (v.Uint64() == 27 || v.Uint64() == 28)
}
