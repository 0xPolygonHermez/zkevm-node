package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/proofofefficiency"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
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
	txs, err := decodeRawTxs(bytesRawTxs)
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

	var rawTxHex string
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

	V, R, S := tx.RawSignatureValues()
	sign := 1 - (V.Uint64() & 1)

	txCodedRlp, err := rlp.EncodeToBytes([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		tx.ChainId(), uint(0), uint(0),
	})
	if err != nil {
		log.Error("error encoding rlp tx: ", err)
		return fmt.Errorf("error encoding rlp tx: " + err.Error())
	}
	newV := new(big.Int).Add(big.NewInt(ether155V), big.NewInt(int64(sign)))
	newRPadded := fmt.Sprintf("%064s", R.Text(hex.Base))
	newSPadded := fmt.Sprintf("%064s", S.Text(hex.Base))
	newVPadded := fmt.Sprintf("%02s", newV.Text(hex.Base))
	rawTxHex = rawTxHex + hex.EncodeToString(txCodedRlp) + newRPadded + newSPadded + newVPadded

	rawTx, err := hex.DecodeString(rawTxHex)
	if err != nil {
		log.Error("error coverting hex string to []byte. Error: ", err)
		return fmt.Errorf("error coverting hex string to []byte. Error: " + err.Error())
	}
	log.Info("encoded tx with signature using RLP in []byte: ", rawTx)
	log.Info("rawtx with signature using RLP in hex: ", hex.EncodeToString(rawTx))

	return nil
}

func printTxs(txs []*types.Transaction, rawTxs []byte) {
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
		shortRlp         = 55  // length of the short rlp codification
		f7               = 247 // 192 + 55 = c0 + shortRlp
		etherNewV        = 35
		mul2             = 2
	)
	txDataLength := len(txsData)
	for pos < int64(txDataLength) {
		num, err := strconv.ParseInt(hex.EncodeToString(txsData[pos:pos+1]), hex.Base, encoding.BitSize64)
		if err != nil {
			log.Error("error parsing header length: ", err)
			return []*types.Transaction{}, err
		}
		// First byte is the length and must be ignored
		len := num - c0 - 1

		if len > shortRlp { // If rlp is bigger than length 55
			// numH is the length of the bytes that give the length of the rlp
			numH, err := strconv.ParseInt(hex.EncodeToString(txsData[pos:pos+1]), hex.Base, encoding.BitSize64)
			if err != nil {
				log.Error("error parsing length of the bytes: ", err)
				return []*types.Transaction{}, err
			}
			// n is the length of the rlp data without the header (1 byte) for example "0xf7"
			n, err := strconv.ParseInt(hex.EncodeToString(txsData[pos+1:pos+1+numH-f7]), hex.Base, encoding.BitSize64) // +1 is the header. For example 0xf7
			if err != nil {
				log.Error("error parsing length: ", err)
				return []*types.Transaction{}, err
			}
			len = n + 1 // +1 is the header. For example 0xf7
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
			log.Error("error decoding tx bytes: ", err, ". fullDataTx: ", hex.EncodeToString(fullDataTx), "\n tx: ", hex.EncodeToString(txInfo))
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
