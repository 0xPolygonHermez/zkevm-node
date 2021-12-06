package etherman

import (
	"fmt"
	"github.com/hermeznetwork/hermez-core/proverclient"
	"math/big"
)

func byteSliceToFixedByteArray(byteSlc []byte) ([32]byte, error) {
	if len(byteSlc) != 32 {
		return [32]byte{}, fmt.Errorf("wrong slice length, current %d, expected 32", len(byteSlc))
	}
	var res [32]byte
	copy(res[:], byteSlc)
	return res, nil
}

func strSliceToBigIntArray(strSlc []string) ([2]*big.Int, error) {
	if len(strSlc) != 2 {
		return [2]*big.Int{}, fmt.Errorf("wrong slice length, current %d, expected 2", len(strSlc))
	}
	var res [2]*big.Int
	for i, v := range strSlc {
		bigInt, ok := new(big.Int).SetString(v, 16)
		if !ok {
			return [2]*big.Int{}, fmt.Errorf("failed to convert string to big int, str: %s", v)
		}
		res[i] = bigInt
	}
	return res, nil
}

func proofSlcToIntArray(proofs []*proverclient.ProofX) ([2][2]*big.Int, error) {
	if len(proofs) != 2 {
		return [2][2]*big.Int{}, fmt.Errorf("wrong proof slice length, current %d, expected 2", len(proofs))
	}

	var res [2][2]*big.Int
	for i, v := range proofs {
		for j, b := range proofs[i].Proof {
			bigInt, ok := new(big.Int).SetString(b, 16)
			if !ok {
				return [2][2]*big.Int{}, fmt.Errorf("failed to convert string to big int, str: %s", v)
			}
			res[i][j] = bigInt
		}
	}

	return res, nil
}
