package state

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// CheckSignature checks a transaction signature
func CheckSignature(tx *types.Transaction) error {
	// Check Signature
	v, r, s := tx.RawSignatureValues()
	plainV := byte(v.Uint64() - 35 - 2*(tx.ChainId().Uint64()))

	if !crypto.ValidateSignatureValues(plainV, r, s, false) {
		return ErrInvalidSig
	}

	return nil
}
