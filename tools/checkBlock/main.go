package main

import (
	"context"
	"math/big"
	"os"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "BlockTool"
	app.Version = "v0.0.1"
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "url",
			Aliases: []string{"u"},
			Value:   "http://localhost:8123",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "getBlockHash",
			Aliases: []string{},
			Usage:   "gets Block by Hash",
			Action:  getBlockHash,
			Flags: append(flags, &cli.StringFlag{
				Name:    "hash",
				Aliases: []string{"bh"},
				Value:   common.Hash{}.String(),
			}),
		},
		{
			Name:    "readBlockHash",
			Aliases: []string{},
			Usage:   "reads the block hash",
			Action:  readBlockHash,
			Flags: append(flags, &cli.Uint64Flag{
				Name:    "blocknumber",
				Aliases: []string{"n"},
				Value:   1,
			}),
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("\nError: %v\n", err)
		os.Exit(1)
	}
}

func getBlockHash(ctx *cli.Context) error {
	blockHash := common.HexToHash(ctx.String("hash"))
	url := ctx.String("url")
	getBlockByHash(ctx.Context, url, blockHash)
	return nil
}

func readBlockHash(ctx *cli.Context) error {
	blockNumber := ctx.Int64("blocknumber")
	url := ctx.String("url")
	block, err := getBlockByNumber(ctx.Context, url, blockNumber)
	if err != nil {
		log.Error("error in getBlockByNumber. Error: ", err)
		return err
	}
	log.Infof("Block full header read (block %d): %+v", block.NumberU64(), block.Header())
	log.Infof("Block stateRoot read (block %d): %+v", block.NumberU64(), block.Header().Root)
	log.Infof("Block computed hash (block %d): %+v", block.NumberU64(), block.Header().Hash())
	log.Infof("Block ReceiptHash read (block %d): %+v", block.NumberU64(), block.Header().ReceiptHash)
	return nil
}

func getBlockByNumber(ctx context.Context, url string, blockNumber int64) (*types.Block, error) {
	// Connect to ethereum node
	ethClient, err := ethclient.Dial(url)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", url, err)
		return nil, err
	}
	bn := big.NewInt(blockNumber)
	block, err := ethClient.BlockByNumber(ctx, bn)
	if err != nil {
		log.Errorf("error getting block number: %d. Error: %v", blockNumber, err)
		return nil, err
	}
	return block, nil
}

func getBlockByHash(ctx context.Context, url string, hash common.Hash) {
	// Connect to ethereum node
	ethClient, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("error connecting to %s: %+v", url, err)
	}

	blockByHash, err := ethClient.BlockByHash(ctx, hash)
	if err != nil {
		log.Errorf("error getting block by computed hash: %s. Error: %v", hash, err)
	} else {
		log.Info("blockByHash:", hash, blockByHash.Hash())
		log.Info("Block StateRoot: ", blockByHash.Header().Root)
	}
}
