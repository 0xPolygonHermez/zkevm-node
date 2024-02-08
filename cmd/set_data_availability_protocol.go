package main

import (
	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

const flagDAAddress = "data-availability-address"

var setDataAvailabilityProtocolFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     flagDAAddress,
		Aliases:  []string{"da-addr"},
		Usage:    "address of the new data availibility protocol",
		Required: true,
	},
	&cli.StringFlag{
		Name:     config.FlagKeyStorePath,
		Aliases:  []string{"ksp"},
		Usage:    "the path of the key store file containing the private key of the account going to set new data availability protocol",
		Required: true,
	},
	&cli.StringFlag{
		Name:     config.FlagPassword,
		Aliases:  []string{"pw"},
		Usage:    "the password do decrypt the key store file",
		Required: true,
	},
	&configFileFlag,
	&networkFlag,
	&customNetworkFlag,
}

func setDataAvailabilityProtocol(ctx *cli.Context) error {
	c, err := config.Load(ctx, true)
	if err != nil {
		return err
	}

	setupLog(c.Log)

	daAddress := common.HexToAddress(ctx.String(flagDAAddress))
	addrKeyStorePath := ctx.String(config.FlagKeyStorePath)
	addrPassword := ctx.String(config.FlagPassword)

	etherman, err := newEtherman(*c, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	auth, _, err := etherman.LoadAuthFromKeyStore(addrKeyStorePath, addrPassword)
	if err != nil {
		log.Fatal(err)
		return err
	}

	tx, err := etherman.SetDataAvailabilityProtocol(auth.From, daAddress)
	if err != nil {
		return err
	}

	log.Infof("Transaction to set new data availability protocol sent. Hash: %s", tx.Hash())

	return nil
}
