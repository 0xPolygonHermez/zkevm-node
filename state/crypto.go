package state

import (
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// ErrInvalidSig indicates the signature of the transaction is not valid
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")
)

// CheckSignature checks a transaction signature
func CheckSignature(tx types.Transaction) error {
	// Check Signature
	v, r, s := tx.RawSignatureValues()
	plainV := byte(v.Uint64() - 35 - 2*(tx.ChainId().Uint64()))

	if !crypto.ValidateSignatureValues(plainV, r, s, false) {
		return ErrInvalidSig
	}

	return nil
}
