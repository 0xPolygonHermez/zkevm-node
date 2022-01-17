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
	configFilePath := ctx.String(flagCfg)
	network := ctx.String(flagNetwork)
	toName := ctx.String(flagAddress)
	a := ctx.String(flagAmount)

	c, err := config.Load(configFilePath, network)
	if err != nil {
		return err
	}
	// log.Warnf("%+v", c)
	var toAddress common.Address
	switch toName {
	case "poe":
		toAddress = c.NetworkConfig.PoEAddr
	case "bridge":
		toAddress = c.NetworkConfig.BridgeAddr
	}

	fmt.Print("*WARNING* Are you sure you want to approve " + a +
		" tokens to be spent by the smc <Name: " + toName + ". Address: " + toAddress.String() + ">? [y/N]: ")
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

	//Check if it is already registered
	etherman, err := newEtherman(*c)
	if err != nil {
		log.Fatal(err)
		return err
	}
	var amount *big.Int
	amount, _ = new(big.Int).SetString(a, 10)

	tx, err := etherman.ApproveMatic(amount, toAddress)
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
