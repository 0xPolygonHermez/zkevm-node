package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

func approveTokens(ctx *cli.Context) error {
	const bitSize uint = 256
	useMaxAmountArg := ctx.Bool(config.FlagMaxAmount)
	var amount *big.Int
	if !useMaxAmountArg {
		amountArg := ctx.String(config.FlagAmount)
		amount, _ = new(big.Int).SetString(amountArg, encoding.Base10)
		if amount == nil {
			fmt.Println("Please, introduce a valid amount in wei")
			return nil
		}
	} else {
		amount = new(big.Int).Sub(new(big.Int).Lsh(common.Big1, bitSize), common.Big1)
	}

	addrKeyStorePath := ctx.String(config.FlagKeyStorePath)
	addrPassword := ctx.String(config.FlagPassword)

	c, err := config.Load(ctx, true)
	if err != nil {
		return err
	}

	if !ctx.Bool(config.FlagYes) {
		fmt.Print("*WARNING* Are you sure you want to approve ", amount.String(),
			" tokens (in wei) for the smc <Name: PoE. Address: "+c.NetworkConfig.L1Config.ZkEVMAddr.String()+">? [y/N]: ")
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			return err
		}
		input = strings.ToLower(input)
		if !(input == "y" || input == "yes") {
			return nil
		}
	}

	setupLog(c.Log)

	// Check if it is already registered
	etherman, err := newEtherman(*c, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// load auth from keystore file
	auth, _, err := etherman.LoadAuthFromKeyStore(addrKeyStorePath, addrPassword)
	if err != nil {
		log.Fatal(err)
		return err
	}

	tx, err := etherman.ApprovePol(ctx.Context, auth.From, amount, c.NetworkConfig.L1Config.ZkEVMAddr)
	if err != nil {
		return err
	}
	const (
		mainnet = 1
		rinkeby = 4
		goerli  = 5
		local   = 1337
	)
	switch c.NetworkConfig.L1Config.L1ChainID {
	case mainnet:
		fmt.Println("Check tx status: https://etherscan.io/tx/" + tx.Hash().String())
	case rinkeby:
		fmt.Println("Check tx status: https://rinkeby.etherscan.io/tx/" + tx.Hash().String())
	case goerli:
		fmt.Println("Check tx status: https://goerli.etherscan.io/tx/" + tx.Hash().String())
	case local:
		fmt.Println("Local network. Tx Hash: " + tx.Hash().String())
	default:
		fmt.Println("Unknown network. Tx Hash: " + tx.Hash().String())
	}
	return nil
}
