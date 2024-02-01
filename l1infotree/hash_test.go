package l1infotree

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestHashLeaf(t *testing.T) {
	expectedLeafHash := common.HexToHash("0xf62f487534b899b1c362242616725878188ca891fab60854b792ca0628286de7")

	prevBlockHash := common.HexToHash("0x24a5871d68723340d9eadc674aa8ad75f3e33b61d5a9db7db92af856a19270bb")
	var minTimestamp uint64 = 1697231573
	ger := common.HexToHash("0x16994edfddddb9480667b64174fc00d3b6da7290d37b8db3a16571b4ddf0789f")

	leaf := HashLeafData(ger, prevBlockHash, minTimestamp)

	assert.Equal(t, expectedLeafHash, common.BytesToHash(leaf[:]))
}
