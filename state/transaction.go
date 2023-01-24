package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// GetSender gets the sender from the transaction's signature
func GetSender(tx types.Transaction) (common.Address, error) {
	signer := types.NewEIP155Signer(tx.ChainId())
	sender, err := signer.Sender(&tx)
	if err != nil {
		return common.Address{}, err
	}
	return sender, nil
}

// RlpFieldsToLegacyTx parses the rlp fields slice into a type.LegacyTx
// in this specific order:
//
// required fields:
// [0] Nonce    uint64
// [1] GasPrice *big.Int
// [2] Gas      uint64
// [3] To       *common.Address
// [4] Value    *big.Int
// [5] Data     []byte
//
// optional fields:
// [6] V        *big.Int
// [7] R        *big.Int
// [8] S        *big.Int
func RlpFieldsToLegacyTx(fields [][]byte) (tx *types.LegacyTx, chainID *big.Int, err error) {
	const (
		fieldsSizeWithoutChainID = 6
		fieldsSizeWithV          = 7
		fieldsSizeWithVR         = 8
		fieldsSizeWithVRS        = 9
	)

	if len(fields) != fieldsSizeWithoutChainID && len(fields) != fieldsSizeWithV && len(fields) != fieldsSizeWithVRS {
		return nil, nil, types.ErrTxTypeNotSupported
	}

	nonce := big.NewInt(0).SetBytes(fields[0]).Uint64()
	gasPrice := big.NewInt(0).SetBytes(fields[1])
	gas := big.NewInt(0).SetBytes(fields[2]).Uint64()
	var to *common.Address
	if fields[3] != nil {
		tmp := common.BytesToAddress(fields[3])
		to = &tmp
	}
	value := big.NewInt(0).SetBytes(fields[4])
	data := fields[5]

	v := big.NewInt(0)
	if len(fields) >= fieldsSizeWithV {
		v = big.NewInt(0).SetBytes(fields[6])
		chainID = big.NewInt(0).Sub(v, big.NewInt(0).SetUint64(etherNewV))
		chainID = big.NewInt(0).Quo(chainID, big.NewInt(double))
	}

	r := big.NewInt(0)
	if len(fields) >= fieldsSizeWithVR {
		r = big.NewInt(0).SetBytes(fields[7])
	}

	s := big.NewInt(0)
	if len(fields) >= fieldsSizeWithVRS {
		s = big.NewInt(0).SetBytes(fields[8])
	}

	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gas,
		To:       to,
		Value:    value,
		Data:     data,
		V:        v,
		R:        r,
		S:        s,
	}, chainID, nil
}
