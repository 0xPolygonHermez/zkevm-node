package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "RlpTool"
	app.Version = "v0.0.1"
	app.Commands = []*cli.Command{
		{
			Name:    "decodeFull",
			Aliases: []string{},
			Usage:   "decode full callData rlp",
			Action:  decodeFull,
		},
		{
			Name:    "decode",
			Aliases: []string{},
			Usage:   "decode rlp",
			Action:  decode,
		},
		{
			Name:    "encode",
			Aliases: []string{},
			Usage:   "encode tx with rlp",
			Action:  encode,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("\nError: %v\n", err)
		os.Exit(1)
	}
}

func decodeFull(ctx *cli.Context) error {
	callData := ctx.Args().First()
	bytesCallData, err := hex.DecodeHex(callData)
	if err != nil {
		log.Error("error decoding callData: ", err)
		return err
	}
	txs, rawTxs, err := decodeFullCallDataToTxs(bytesCallData)
	if err != nil {
		return err
	}
	printTxs(txs, rawTxs)
	return nil
}

func decode(ctx *cli.Context) error {
	rawTxs := ctx.Args().First()
	bytesRawTxs, err := hex.DecodeHex(rawTxs)
	if err != nil {
		log.Error("error decoding rawTxs: ", err)
		return err
	}
	txs, _, err := state.DecodeTxs(bytesRawTxs)
	if err != nil {
		log.Error("error decoding tx callData: ", err)
		return err
	}
	printTxs(txs, bytesRawTxs)
	return nil
}

func encode(ctx *cli.Context) error {
	fmt.Print("Nonce : ")
	var nonceS string
	if _, err := fmt.Scanln(&nonceS); err != nil {
		return err
	}
	nonce, err := strconv.ParseUint(nonceS, encoding.Base10, 64) //nolint:gomnd
	if err != nil {
		log.Error("error decoding nonce: ", err)
		return err
	}
	log.Info("Nonce: ", nonce)

	fmt.Print("GasPrice : ")
	var gasPriceS string
	if _, err := fmt.Scanln(&gasPriceS); err != nil {
		return err
	}
	gasPrice, _ := new(big.Int).SetString(gasPriceS, encoding.Base10)
	log.Info("GasPrice: ", gasPrice)

	fmt.Print("Gas : ")
	var gasS string
	if _, err := fmt.Scanln(&gasS); err != nil {
		return err
	}
	gas, err := strconv.ParseUint(gasS, encoding.Base10, 64) //nolint:gomnd
	if err != nil {
		log.Error("error decoding gas: ", err)
		return err
	}
	log.Info("Gas: ", gas)

	fmt.Print("To : ")
	var toS string
	if _, err := fmt.Scanln(&toS); err != nil {
		return err
	}
	to := common.HexToAddress(toS)
	log.Info("To: ", to)

	fmt.Print("Value : ")
	var valueS string
	if _, err := fmt.Scanln(&valueS); err != nil {
		return err
	}
	value, _ := new(big.Int).SetString(valueS, encoding.Base10)
	log.Info("Value: ", value)

	fmt.Print("Data : ")
	var dataS string
	if _, err := fmt.Scanln(&dataS); err != nil {
		if err.Error() != "unexpected newline" {
			return err
		}
	}
	var data []byte
	if dataS != "" {
		data, err = hex.DecodeHex(dataS)
		if err != nil {
			log.Error("error decoding data: ", err)
			return err
		}
	}
	log.Info("Data: ", data)

	fmt.Print("V: ")
	var vS string
	if _, err := fmt.Scanln(&vS); err != nil {
		return err
	}
	v, _ := new(big.Int).SetString(vS, encoding.Base10)
	log.Info("V: ", v)

	fmt.Print("R: ")
	var rS string
	if _, err := fmt.Scanln(&rS); err != nil {
		return err
	}
	r, _ := new(big.Int).SetString(rS, encoding.Base10)
	log.Info("R: ", r)

	fmt.Print("S: ")
	var sS string
	if _, err := fmt.Scanln(&sS); err != nil {
		return err
	}
	s, _ := new(big.Int).SetString(sS, encoding.Base10)
	log.Info("S: ", s)

	var txLegacy = types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gas,
		To:       &to,
		Value:    value,
		Data:     data,
		V:        v,
		R:        r,
		S:        s,
	}
	tx := types.NewTx(&txLegacy)

	rawBytes, err := state.EncodeTransactions([]types.Transaction{*tx})
	if err != nil {
		log.Error("error encoding txs: ", err)
		return err
	}
	log.Info("encoded tx with signature using RLP in []byte: ", rawBytes)
	log.Info("rawtx with signature using RLP in hex: ", hex.EncodeToString(rawBytes))

	return nil
}

func printTxs(txs []types.Transaction, rawTxs []byte) {
	log.Info("RawTxs: ", hex.EncodeToHex(rawTxs))
	for _, tx := range txs {
		log.Info("#######################################################################")
		log.Info("#######################################################################")
		log.Infof("Decoded tx: %+v", tx)
		log.Info("ChainID: ", tx.ChainId())
		log.Info("Cost: ", tx.Cost())
		log.Info("Data: ", hex.EncodeToString(tx.Data()))
		log.Info("Gas: ", tx.Gas())
		log.Info("GasPrice: ", tx.GasPrice())
		log.Info("Hash: ", tx.Hash())
		log.Info("Nonce: ", tx.Nonce())
		v, r, s := tx.RawSignatureValues()
		log.Info("V: ", v)
		log.Info("R: ", r)
		log.Info("S: ", s)
		log.Info("To: ", tx.To())
		log.Info("Type: ", tx.Type())
		log.Info("Value: ", tx.Value())
	}
}

func decodeFullCallDataToTxs(txsData []byte) ([]types.Transaction, []byte, error) {
	// The first 4 bytes are the function hash bytes. These bytes has to be ripped.
	// After that, the unpack method is used to read the call data.
	// The txs data is a chunk of concatenated rawTx. This rawTx is the encoded tx information in rlp + the signature information (v, r, s).
	//So, txs data will look like: txRLP+r+s+v+txRLP2+r2+s2+v2

	// Extract coded txs.
	// Load contract ABI
	abi, err := abi.JSON(strings.NewReader(polygonzkevm.PolygonzkevmABI))
	if err != nil {
		log.Fatal("error reading smart contract abi: ", err)
	}

	// Recover Method from signature and ABI
	method, err := abi.MethodById(txsData[:4])
	if err != nil {
		log.Fatal("error getting abi method: ", err)
	}

	// Unpack method inputs
	data, err := method.Inputs.Unpack(txsData[4:])
	if err != nil {
		log.Fatal("error reading call data: ", err)
	}

	txsData = data[0].([]byte)

	txs, _, err := state.DecodeTxs(txsData)
	return txs, txsData, err
}
