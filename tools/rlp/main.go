package main

import (
	"math/big"
	"strings"
	"strconv"
	"os"

	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/proofofefficiency"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/urfave/cli/v2"
)

const (
	ether155V = 27
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
	txs, err := decodeRawTxs(bytesRawTxs)
	if err != nil {
		log.Error("error decoding tx callData: ", err)
		return err
	}
	printTxs(txs, bytesRawTxs)
	return nil
}

func printTxs(txs []*types.Transaction, rawTxs []byte) {
	log.Debug("RawTxs: ", hex.EncodeToHex(rawTxs))
	for _, tx := range txs {
		log.Debug("#######################################################################")
		log.Debug("#######################################################################")
		log.Debugf("Decoded tx: %+v", tx)
		log.Debug("ChainID: ", tx.ChainId())
		log.Debug("Cost: ", tx.Cost())
		log.Debug("Data: ", hex.EncodeToString(tx.Data()))
		log.Debug("Gas: ", tx.Gas())
		log.Debug("GasPrice: ", tx.GasPrice())
		log.Debug("Hash: ", tx.Hash())
		log.Debug("Nonce: ", tx.Nonce())
		v, r, s := tx.RawSignatureValues()
		log.Debug("V: ", v)
		log.Debug("R: ", r)
		log.Debug("S: ", s)
		log.Debug("To: ", tx.To())
		log.Debug("Type: ", tx.Type())
		log.Debug("Value: ", tx.Value())
	}
}

func decodeFullCallDataToTxs(txsData []byte) ([]*types.Transaction, []byte, error) {
	// The first 4 bytes are the function hash bytes. These bytes has to be ripped.
	// After that, the unpack method is used to read the call data.
	// The txs data is a chunk of concatenated rawTx. This rawTx is the encoded tx information in rlp + the signature information (v, r, s).
	//So, txs data will look like: txRLP+r+s+v+txRLP2+r2+s2+v2

	// Extract coded txs.
	// Load contract ABI
	abi, err := abi.JSON(strings.NewReader(proofofefficiency.ProofofefficiencyABI))
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

	txs, err := decodeRawTxs(txsData)
	return txs, txsData, err
}


func decodeRawTxs(txsData []byte) ([]*types.Transaction, error) {
	// Process coded txs
	var pos int64
	var txs []*types.Transaction
	const (
		headerByteLength = 2
		sLength          = 32
		rLength          = 32
		vLength          = 1
		c0               = 192 // 192 is c0. This value is defined by the rlp protocol
		ff               = 255 // max value of rlp header
		shortRlp         = 55 // length of the short rlp codification
		f7               = 247 // 192 + 55 = c0 + shortRlp
		etherNewV        = 35
		mul2             = 2
	)
	txDataLength := len(txsData)
	for pos < int64(txDataLength) {
		num, err := strconv.ParseInt(hex.EncodeToString(txsData[pos : pos+1]), hex.Base, encoding.BitSize64)
		if err != nil {
			log.Debug("error parsing header length: ", err)
			return []*types.Transaction{}, err
		}
		// First byte is the length and must be ignored
		len := num - c0 - 1

		if len > shortRlp { // If rlp is bigger than lenght 55
			// numH is the length of the bytes that give the length of the rlp
			numH, err := strconv.ParseInt(hex.EncodeToString(txsData[pos : pos+1]), hex.Base, encoding.BitSize64)
			if err != nil {
				log.Debug("error parsing length of the bytes: ", err)
				return []*types.Transaction{}, err
			}
			// n is the length of the rlp data without the header (1 byte) for example "0xf7"
			n, err := strconv.ParseInt(hex.EncodeToString(txsData[pos+1 : pos+1+numH-f7]), hex.Base, encoding.BitSize64) // +1 is the header. For example 0xf7
			if err != nil {
				log.Debug("error parsing length: ", err)
				return []*types.Transaction{}, err
			}
			len = n+1 // +1 is the header. For example 0xf7
		}

		fullDataTx := txsData[pos : pos+len+rLength+sLength+vLength+headerByteLength]
		txInfo := txsData[pos : pos+len+headerByteLength]
		r := txsData[pos+len+headerByteLength : pos+len+rLength+headerByteLength]
		s := txsData[pos+len+rLength+headerByteLength : pos+len+rLength+sLength+headerByteLength]
		v := txsData[pos+len+rLength+sLength+headerByteLength : pos+len+rLength+sLength+vLength+headerByteLength]

		pos = pos + len + rLength + sLength + vLength + headerByteLength

		// Decode tx
		var tx types.LegacyTx
		err = rlp.DecodeBytes(txInfo, &tx)
		if err != nil {
			log.Debug("error decoding tx bytes: ", err, ". fullDataTx: ", hex.EncodeToString(fullDataTx), "\n tx: ", hex.EncodeToString(txInfo))
			return []*types.Transaction{}, err
		}

		//tx.V = v-27+chainId*2+35
		tx.V = new(big.Int).Add(new(big.Int).Sub(new(big.Int).SetBytes(v), big.NewInt(ether155V)), new(big.Int).Add(new(big.Int).Mul(tx.V, big.NewInt(mul2)), big.NewInt(etherNewV)))
		tx.R = new(big.Int).SetBytes(r)
		tx.S = new(big.Int).SetBytes(s)

		txs = append(txs, types.NewTx(&tx))
	}
	return txs, nil
}