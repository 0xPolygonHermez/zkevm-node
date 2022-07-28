package state

import (
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

const ether155V = 27

// EncodeTransactions RLP encodes the given transactions.
func EncodeTransactions(txs []types.Transaction) ([]byte, error) {
	var batchL2Data []byte

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

// EncodeUnsignedTransaction RLP encodes the given unsigned transaction
func EncodeUnsignedTransaction(tx types.Transaction) ([]byte, error) {
	v, _ := new(big.Int).SetString("0x1c", 0)
	r, _ := new(big.Int).SetString("0xa54492cfacf71aef702421b7fbc70636537a7b2fbe5718c5ed970a001bb7756b", 0)
	s, _ := new(big.Int).SetString("0x2e9fb27acc75955b898f0b12ec52aa34bf08f01db654374484b80bf12f0d841e", 0)

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

	return txData, nil
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
