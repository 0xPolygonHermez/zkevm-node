package state

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/proofofefficiency"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

const ether155V = 27

// EncodeTransactions RLP encodes the given transactions.
func EncodeTransactions(txs []types.Transaction) ([]byte, error) {
	var batchL2Data []byte

	// TODO: Check how to encode unsigned transactions

	for _, tx := range txs {
		v, r, s := tx.RawSignatureValues()
		sign := 1 - (v.Uint64() & 1)

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
			return nil, err
		}

		newV := new(big.Int).Add(big.NewInt(ether155V), big.NewInt(int64(sign)))
		newRPadded := fmt.Sprintf("%064s", r.Text(hex.Base))
		newSPadded := fmt.Sprintf("%064s", s.Text(hex.Base))
		newVPadded := fmt.Sprintf("%02s", newV.Text(hex.Base))
		txData, err := hex.DecodeString(hex.EncodeToString(txCodedRlp) + newRPadded + newSPadded + newVPadded)
		if err != nil {
			return nil, err
		}

		batchL2Data = append(batchL2Data, txData...)
	}

	return batchL2Data, nil
}

// DecodeTxs extracts Tansactions for its encoded form
func DecodeTxs(txsData []byte) ([]types.Transaction, []byte, error) {
	// The first 4 bytes are the function hash bytes. These bytes has to be ripped.
	// After that, the unpack method is used to read the call data.
	// The txs data is a chunk of concatenated rawTx. This rawTx is the encoded tx information in rlp + the signature information (v, r, s).
	//So, txs data will look like: txRLP+r+s+v+txRLP2+r2+s2+v2

	// Extract coded txs.
	// Load contract ABI
	abi, err := abi.JSON(strings.NewReader(proofofefficiency.ProofofefficiencyMetaData.ABI))
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

	// Process coded txs
	var pos int64
	var txs []types.Transaction
	const (
		headerByteLength = 1
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
	if txDataLength == 0 {
		return txs, txsData, nil
	}
	for pos < int64(txDataLength) {
		num, err := strconv.ParseInt(hex.EncodeToString(txsData[pos:pos+1]), hex.Base, encoding.BitSize64)
		if err != nil {
			log.Debug("error parsing header length: ", err)
			return []types.Transaction{}, []byte{}, err
		}
		// First byte is the length and must be ignored
		len := num - c0
		if len > shortRlp { // If rlp is bigger than length 55
			// n is the length of the rlp data without the header (1 byte) for example "0xf7"
			n, err := strconv.ParseInt(hex.EncodeToString(txsData[pos+1:pos+1+num-f7]), hex.Base, encoding.BitSize64) // +1 is the header. For example 0xf7
			if err != nil {
				log.Debug("error parsing length: ", err)
				return []types.Transaction{}, []byte{}, err
			}
			len = n + num - f7 // num - f7 is the header. For example 0xf7
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
			log.Debug("error decoding tx bytes: ", err, ". fullDataTx: ", hex.EncodeToString(fullDataTx), "\n tx: ", hex.EncodeToString(txInfo), "\n Txs received: ", hex.EncodeToString(txsData))
			return []types.Transaction{}, []byte{}, err
		}

		//tx.V = v-27+chainId*2+35
		tx.V = new(big.Int).Add(new(big.Int).Sub(new(big.Int).SetBytes(v), big.NewInt(ether155V)), new(big.Int).Add(new(big.Int).Mul(tx.V, big.NewInt(mul2)), big.NewInt(etherNewV)))
		tx.R = new(big.Int).SetBytes(r)
		tx.S = new(big.Int).SetBytes(s)

		txs = append(txs, *types.NewTx(&tx))
	}
	return txs, txsData, nil
}

// DecodeTx decodes a string rlp tx representation into a types.Transaction instance
func DecodeTx(encodedTx string) (*types.Transaction, error) {
	b, err := hex.DecodeHex(encodedTx)
	if err != nil {
		return nil, err
	}

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return tx, nil
}

func generateReceipt(block *types.Block, processedTx *ProcessTransactionResponse) *types.Receipt {
	receipt := &types.Receipt{
		Type:              uint8(processedTx.Type),
		PostState:         processedTx.StateRoot.Bytes(),
		CumulativeGasUsed: processedTx.GasUsed,
		BlockNumber:       block.Number(),
		BlockHash:         block.Hash(),
		GasUsed:           processedTx.GasUsed,
		TxHash:            processedTx.Tx.Hash(),
		TransactionIndex:  0,
		ContractAddress:   processedTx.CreateAddress,
		Logs:              processedTx.Logs,
	}

	if processedTx.Error == "" {
		receipt.Status = types.ReceiptStatusSuccessful
	} else {
		receipt.Status = types.ReceiptStatusFailed
	}

	return receipt
}
