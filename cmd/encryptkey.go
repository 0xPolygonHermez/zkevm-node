package main

import (
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli/v2"
)

const (
	encryptKeyFlagPrivateKey = "privateKey"
	encryptKeyFlagPassword   = "password"
	encryptKeyFlagOutput     = "output"
)

var encryptKeyFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     encryptKeyFlagPrivateKey,
		Aliases:  []string{"pk"},
		Usage:    "Private key hash",
		Required: true,
	},
	&cli.StringFlag{
		Name:     encryptKeyFlagPassword,
		Aliases:  []string{"pw"},
		Usage:    "Password to encrypt the private key",
		Required: true,
	},
	&cli.StringFlag{
		Name:     encryptKeyFlagOutput,
		Aliases:  []string{"o"},
		Usage:    "Output directory to save the encrypted private key file",
		Required: true,
	},
}

func encryptKey(ctx *cli.Context) error {
	privateKeyHash := ctx.String(encryptKeyFlagPrivateKey)
	password := ctx.String(encryptKeyFlagPassword)
	outputDir := ctx.String(encryptKeyFlagOutput)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHash, "0x"))
	if err != nil {
		log.Fatal("Invalid private key: ", err)
	}

	ks := keystore.NewKeyStore(outputDir, keystore.StandardScryptN, keystore.StandardScryptP)
	if _, err := ks.ImportECDSA(privateKey, password); err != nil {
		log.Fatal("Failed to encrypt private key: ", err)
	}
	return nil
}
