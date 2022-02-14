package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/config"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/urfave/cli/v2"
)

func approveTokens(ctx *cli.Context) error {
	toName := ctx.String(flagAddress)
	a := ctx.String(flagAmount)
	amount, _ := new(big.Float).SetString(a)
	if amount == nil {
		fmt.Println("Please, introduce a valid amount. Use '.' instead of ',' if it is a decimal number")
		return nil
	}
	c, err := config.Load(ctx)
	if err != nil {
		return err
	}
	var toAddress common.Address
	switch toName {
	case "poe":
		toAddress = c.NetworkConfig.PoEAddr
	case "bridge":
		toAddress = c.NetworkConfig.BridgeAddr
	}

	fmt.Print("*WARNING* Are you sure you want to approve ", amount,
		" tokens to be spent by the smc <Name: "+toName+". Address: "+toAddress.String()+">? [y/N]: ")
	var input string
	if _, err := fmt.Scanln(&input); err != nil {
		return err
	}
	input = strings.ToLower(input)
	if !(input == "y" || input == "yes") {
		return nil
	}

	setupLog(c.Log)

	runMigrations(c.Database)

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
	tx, err := etherman.ApproveMatic(amountB, toAddress)
	if err != nil {
		return err
	}
	const (
		mainnet = 1
		rinkeby = 4
		goerli  = 5
	)
	switch c.NetworkConfig.L1ChainID {
	case mainnet:
		fmt.Println("Check tx status: https://etherscan.io/tx/" + tx.Hash().String())
	case rinkeby:
		fmt.Println("Check tx status: https://rinkeby.etherscan.io/tx/" + tx.Hash().String())
	case goerli:
		fmt.Println("Check tx status: https://goerli.etherscan.io/tx/" + tx.Hash().String())
	}
	return nil
}
