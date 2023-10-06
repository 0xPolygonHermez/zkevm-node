package main

import (
	"encoding/binary"
	"os"
	"reflect"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/0xPolygonHermez/zkevm-data-streamer/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/tools/datastreamer/config"
	"github.com/0xPolygonHermez/zkevm-node/tools/datastreamer/db"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

const appName = "zkevm-data-streamer-tool"

var (
	configFileFlag = cli.StringFlag{
		Name:     config.FlagCfg,
		Aliases:  []string{"c"},
		Usage:    "Configuration `FILE`",
		Required: false,
	}
)

func main() {
	app := cli.NewApp()
	app.Name = appName

	flags := []cli.Flag{
		&configFileFlag,
	}

	app.Commands = []*cli.Command{
		{
			Name:    "generate",
			Aliases: []string{},
			Usage:   "Generate stream file form scratch",
			Action:  generate,
			Flags:   flags,
		},
		{
			Name:    "rebuild",
			Aliases: []string{},
			Usage:   "Rebuild state roots from a block",
			Action:  rebuild,
			Flags:   flags,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func initializeStreamServer(c *config.Config) (*datastreamer.StreamServer, error) {
	// Init logger
	log.Init(c.StreamServer.Log)
	log.Info("Starting tool")

	// Create a stream server
	streamServer, err := datastreamer.New(c.StreamServer.Port, state.StreamTypeSequencer, c.StreamServer.Filename, &c.StreamServer.Log)
	if err != nil {
		return nil, err
	}

	// Set entities definition
	entriesDefinition := map[datastreamer.EntryType]datastreamer.EntityDefinition{
		state.EntryTypeL2BlockStart: {
			Name:       "L2BlockStart",
			StreamType: state.StreamTypeSequencer,
			Definition: reflect.TypeOf(state.DSL2BlockStart{}),
		},
		state.EntryTypeL2Tx: {
			Name:       "L2Transaction",
			StreamType: state.StreamTypeSequencer,
			Definition: reflect.TypeOf(state.DSL2Transaction{}),
		},
		state.EntryTypeL2BlockEnd: {
			Name:       "L2BlockEnd",
			StreamType: state.StreamTypeSequencer,
			Definition: reflect.TypeOf(state.DSL2BlockEnd{}),
		},
	}

	streamServer.SetEntriesDef(entriesDefinition)
	err = streamServer.Start()
	if err != nil {
		return nil, err
	}

	return &streamServer, nil
}

func generate(cliCtx *cli.Context) error {
	c, err := config.Load(cliCtx)
	if err != nil {
		log.Fatal(err)

	}
	log.Infof("Loaded configuration: %+v", c)

	streamServer, err := initializeStreamServer(c)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the database
	stateSqlDB, err := db.NewSQLDB(c.StateDB)
	if err != nil {
		log.Fatal(err)
	}
	defer stateSqlDB.Close()
	stateDB := db.NewStateDB(stateSqlDB)
	log.Info("Connected to the database")

	header := streamServer.GetHeader()

	var currentL2Block uint64
	var currentTxIndex uint64

	if header.TotalEntries == 0 {
		// Get Genesis block
		genesisL2Block, err := stateDB.GetGenesisBlock(cliCtx.Context)
		if err != nil {
			log.Fatal(err)
		}

		err = streamServer.StartAtomicOp()
		if err != nil {
			log.Fatal(err)
		}

		genesisBlock := state.DSL2BlockStart{
			BatchNumber:    genesisL2Block.BatchNumber,
			L2BlockNumber:  genesisL2Block.L2BlockNumber,
			Timestamp:      genesisL2Block.Timestamp,
			GlobalExitRoot: genesisL2Block.GlobalExitRoot,
			Coinbase:       genesisL2Block.Coinbase,
			ForkID:         genesisL2Block.ForkID,
		}

		log.Infof("Genesis block: %+v", genesisBlock)

		_, err = streamServer.AddStreamEntry(1, genesisBlock.Encode())
		if err != nil {
			log.Fatal(err)
		}

		genesisBlockEnd := state.DSL2BlockEnd{
			L2BlockNumber: genesisL2Block.L2BlockNumber,
			BlockHash:     genesisL2Block.BlockHash,
			StateRoot:     genesisL2Block.StateRoot,
		}

		_, err = streamServer.AddStreamEntry(state.EntryTypeL2BlockEnd, genesisBlockEnd.Encode())
		if err != nil {
			log.Fatal(err)
		}

		err = streamServer.CommitAtomicOp()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		latestEntry, err := streamServer.GetEntry(header.TotalEntries - 1)
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("Latest entry: %+v", latestEntry)

		switch latestEntry.EntryType {
		case state.EntryTypeL2BlockStart:
			log.Info("Latest entry type is L2BlockStart")
			currentL2Block = binary.LittleEndian.Uint64(latestEntry.Data[8:16])
		case state.EntryTypeL2Tx:
			log.Info("Latest entry type is L2Tx")

			for latestEntry.EntryType == state.EntryTypeL2Tx {
				currentTxIndex++
				latestEntry, err = streamServer.GetEntry(header.TotalEntries - currentTxIndex)
				if err != nil {
					log.Fatal(err)
				}
			}

			if latestEntry.EntryType != state.EntryTypeL2BlockStart {
				log.Fatal("Latest entry is not a L2BlockStart")
			}
			currentL2Block = binary.LittleEndian.Uint64(latestEntry.Data[8:16])

		case state.EntryTypeL2BlockEnd:
			log.Info("Latest entry type is L2BlockEnd")
			currentL2Block = binary.LittleEndian.Uint64(latestEntry.Data[0:8])
		}
	}

	log.Infof("Current transaction index: %d", currentTxIndex)
	log.Infof("Current L2 block number: %d", currentL2Block)

	var limit uint64 = 1000
	var offset uint64 = currentL2Block
	var entry uint64 = header.TotalEntries
	var l2blocks []*state.DSL2Block

	if entry > 0 {
		entry--
	}

	for err == nil {
		log.Infof("Current entry number: %d", entry)

		l2blocks, err = stateDB.GetL2Blocks(cliCtx.Context, limit, offset)
		offset += limit
		if len(l2blocks) == 0 {
			break
		}
		// Get transactions for all the retrieved l2 blocks
		l2Transactions, err := stateDB.GetL2Transactions(cliCtx.Context, l2blocks[0].L2BlockNumber, l2blocks[len(l2blocks)-1].L2BlockNumber)
		if err != nil {
			log.Fatal(err)
		}

		err = streamServer.StartAtomicOp()
		if err != nil {
			log.Fatal(err)
		}

		for x, l2block := range l2blocks {
			if currentTxIndex > 0 {
				x += int(currentTxIndex)
				currentTxIndex = 0
			}

			blockStart := state.DSL2BlockStart{
				BatchNumber:    l2block.BatchNumber,
				L2BlockNumber:  l2block.L2BlockNumber,
				Timestamp:      l2block.Timestamp,
				GlobalExitRoot: l2block.GlobalExitRoot,
				Coinbase:       l2block.Coinbase,
				ForkID:         l2block.ForkID,
			}

			_, err = streamServer.AddStreamEntry(state.EntryTypeL2BlockStart, blockStart.Encode())
			if err != nil {
				log.Fatal(err)
			}

			entry, err = streamServer.AddStreamEntry(state.EntryTypeL2Tx, l2Transactions[x].Encode())
			if err != nil {
				log.Fatal(err)
			}

			blockEnd := state.DSL2BlockEnd{
				L2BlockNumber: l2block.L2BlockNumber,
				BlockHash:     l2block.BlockHash,
				StateRoot:     l2block.StateRoot,
			}

			_, err = streamServer.AddStreamEntry(state.EntryTypeL2BlockEnd, blockEnd.Encode())
			if err != nil {
				log.Fatal(err)
			}
		}
		err = streamServer.CommitAtomicOp()
		if err != nil {
			log.Fatal(err)
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Info("Finished tool")

	return nil
}

func rebuild(cliCtx *cli.Context) error {
	c, err := config.Load(cliCtx)
	if err != nil {
		log.Fatal(err)

	}
	log.Infof("Loaded configuration: %+v", c)

	ctx := cliCtx.Context

	streamServer, err := initializeStreamServer(c)
	if err != nil {
		log.Fatal(err)
	}

	oldEndEntry, err := streamServer.GetEntry(1)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("oldEndEntry: %+v", oldEndEntry)

	startEntry, err := streamServer.GetEntry(2)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("startEntry: %+v", startEntry)
	log.Infof("Length of data in startEntry: %d", len(startEntry.Data))

	txEntry, err := streamServer.GetEntry(3)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("txEntry: %+v", txEntry)

	endEntry, err := streamServer.GetEntry(4)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("endEntry: %+v", endEntry)

	tx, err := state.DecodeTx(string(txEntry.Data[6:]))
	if err != nil {
		log.Fatal(err)
	}

	/*
		log.Infof("tx nonce: %+v", tx.Nonce())

		sender, err := state.GetSender(*tx)
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("tx sender: %+v", sender)
	*/

	// RLP encode the transaction using the proper fork id
	batchL2Data, err := state.EncodeTransaction(*tx, 255, uint64(binary.LittleEndian.Uint16(startEntry.Data[76:78])))
	if err != nil {
		log.Fatal(err)
	}

	processBatchRequest := &executor.ProcessBatchRequest{
		OldBatchNum:      uint64(startEntry.Data[0]) - 1,
		Coinbase:         common.Bytes2Hex(startEntry.Data[56:76]),
		BatchL2Data:      batchL2Data,
		OldStateRoot:     oldEndEntry.Data[40:72],
		GlobalExitRoot:   startEntry.Data[23:55],
		OldAccInputHash:  []byte{},
		EthTimestamp:     binary.LittleEndian.Uint64(startEntry.Data[16:24]),
		UpdateMerkleTree: uint32(0),
		ChainId:          1001,
		ForkId:           uint64(binary.LittleEndian.Uint16(startEntry.Data[76:78])),
	}

	log.Infof("ProcessBatchRequest: %+v", processBatchRequest)

	log.Infof("Old state root:    %s", common.Bytes2Hex(processBatchRequest.OldStateRoot))
	log.Infof("Expected new root: %s", common.Bytes2Hex(endEntry.Data[40:72]))
	log.Infof("Batch L2 data:     %s", common.Bytes2Hex(processBatchRequest.BatchL2Data))
	log.Infof("Coinbase:          %s", processBatchRequest.Coinbase)
	log.Infof("Timestamp:         %d", processBatchRequest.EthTimestamp)
	log.Infof("Fork id:           %d", processBatchRequest.ForkId)

	/*
		if processBatchRequest.ForkId >= 5 {
			processBatchRequest.BatchL2Data = append(processBatchRequest.BatchL2Data, byte(255))
		}
	*/

	// Connect to the executor
	executorClient, executorClientConn, executorCancel := executor.NewExecutorClient(ctx, c.Executor)
	s := executorClientConn.GetState()
	log.Infof("executorClientConn state: %s", s.String())
	defer func() {
		executorCancel()
		executorClientConn.Close()
	}()

	/*
		mtDBServerConfig := merkletree.Config{URI: fmt.Sprintf("%s:50061", "toni-prover")}
		var mtDBCancel context.CancelFunc
		mtDBServiceClient, mtDBClientConn, mtDBCancel := merkletree.NewMTDBServiceClient(ctx, mtDBServerConfig)
		s = mtDBClientConn.GetState()
		log.Infof("stateDbClientConn state: %s", s.String())
		defer func() {
			mtDBCancel()
			mtDBClientConn.Close()
		}()

		stateTree := merkletree.NewStateTree(mtDBServiceClient)

		// Connect to MT
		nonce, err := stateTree.GetNonce(ctx, sender, processBatchRequest.OldStateRoot)
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("Nonce in MT: %d", nonce)
	*/

	// Process batch
	processBatchResponse, err := executorClient.ProcessBatch(ctx, processBatchRequest)
	if err != nil {
		log.Fatal(err)
	}

	if processBatchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
		log.Fatal(processBatchResponse.Error)
	}

	log.Infof("ProcessBatchResponse: %+v", processBatchResponse)

	log.Infof("New root: %s", common.Bytes2Hex(processBatchResponse.NewStateRoot))

	return nil
}
