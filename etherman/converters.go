package etherman

import (
	"fmt"
	"math/big"

	"github.com/hermeznetwork/hermez-core/proverclient"
)

const (
	minProofLen = 2
	maxProofLen = 3
)

func stringToFixedByteArray(str string) ([32]byte, error) {
	var res [32]byte
	copy(res[:], str)
	return res, nil
}

func strSliceToBigIntArray(data []string) ([2]*big.Int, error) {
    if len(data) < minProofLen || len(data) >  maxProofLen {
		return [2]*big.Int{}, fmt.Errorf("wrong slice length, current %d, expected between %d or %d", len(data), minProofLen, maxProofLen)
	}
	var res [2]*big.Int
	for i, v := range data {
		if i < minProofLen {
			bigInt, ok := new(big.Int).SetString(v, 16)
			if !ok {
				return [2]*big.Int{}, fmt.Errorf("failed to convert string to big int, str: %s", v)
			}
			res[i] = bigInt
		}
	}
	return res, nil
}

func proofSlcToIntArray(proofs []*proverclient.ProofX) ([2][2]*big.Int, error) {
	if len(proofs) != minProofLen {
		return [2][2]*big.Int{}, fmt.Errorf("wrong proof slice length, current %d, expected 2", len(proofs))
	}

	var res [2][2]*big.Int
	for i, v := range proofs {
		for j, b := range proofs[i].Proof {
			if j < minProofLen {
				bigInt, ok := new(big.Int).SetString(b, 16)
				if !ok {
					return [2][2]*big.Int{}, fmt.Errorf("failed to convert string to big int, str: %s", v)
				}
				res[i][1-j] = bigInt
			}
		}
	}

	return res, nil
}
