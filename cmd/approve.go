package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/urfave/cli/v2"
)

func approveTokens(ctx *cli.Context) error {
	a := ctx.String(config.FlagAmount)
	amount, _ := new(big.Float).SetString(a)
	if amount == nil {
		fmt.Println("Please, introduce a valid amount. Use '.' instead of ',' if it is a decimal number")
		return nil
	}
	c, err := config.Load(ctx)
	if err != nil {
		return err
	}

	if !ctx.Bool(config.FlagYes) {
		fmt.Print("*WARNING* Are you sure you want to approve ", amount,
			" tokens to be spent by the smc <Name: PoE. Address: "+c.Etherman.PoEAddr.String()+">? [y/N]: ")
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

	runStateMigrations(c.StateDB)
	runPoolMigrations(c.PoolDB)
	runRPCMigrations(c.RPC.DB)

	// Check if it is already registered
	etherman, err := newEtherman(*c)
	if err != nil {
		log.Fatal(err)
		return err
	}

	const decimals = 1000000000000000000
	amountInWei := new(big.Float).Mul(amount, big.NewFloat(decimals))
	amountB := new(big.Int)
	amountInWei.Int(amountB)
	tx, err := etherman.ApproveMatic(amountB, c.Etherman.PoEAddr)
	if err != nil {
		return err
	}
	const (
		mainnet = 1
		rinkeby = 4
		goerli  = 5
		local   = 1337
	)
	switch c.Etherman.L1ChainID {
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
