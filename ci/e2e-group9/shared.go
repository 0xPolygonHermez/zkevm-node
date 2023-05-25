package e2e

import "github.com/ethereum/go-ethereum/common"

const (
	toAddressHex      = "0x4d5Cf5032B2a844602278b01199ED191A86c93ff"
	gerFinalityBlocks = uint64(250)
)

var (
	toAddress = common.HexToAddress(toAddressHex)
)
