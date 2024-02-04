package pool

import (
	"github.com/ethereum/go-ethereum/common"
)

const (
	// BridgeClaimMethodSignature for tracking BridgeClaimMethodSignature method
	BridgeClaimMethodSignature = "0x2cffd02e"
)

func contains(s []string, ele common.Address) bool {
	for _, e := range s {
		if common.HexToAddress(e) == ele {
			return true
		}
	}
	return false
}
