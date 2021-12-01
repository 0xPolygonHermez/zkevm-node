package vectors

import (
	"github.com/ethereum/go-ethereum/common"
)

type StateTransitionVector struct {
	StateTests []StateTest `json:"tests"`
}

type StateTest struct {
	ID               uint                    `json:"id"`
	Description      string                  `json:"description"`
	Arity            uint                    `json:"arity"`
	ChanIDSequencer  uint64                  `json:"chainIdSequencer"`
	SequencerAddress common.MixedcaseAddress `json:"sequencerAddress"`

	Genesis          []Genesis               `json:"genesis"`
	ExpectedOldRoot  []byte                  `json:"expectedOldRoot"`
	Txs              []Tx                    `json:"txs"`
	ExpectedNewRoot  []byte                  `json:"expectedNewRoot"`
	ExpectedNewLeafs map[common.Address]Leaf `json:"expectedNewLeafs"`
}

type Genesis struct {
	Address common.MixedcaseAddress `json:"address"`
	PvtKey  string                  `json:"pvtKey"`
	Balance argBigInt               `json:"balance"`
	Nonce   uint64                  `json:"nonce"`
}

type Tx struct {
	From     common.MixedcaseAddress `json:"from"`
	To       common.MixedcaseAddress `json:"to"`
	Nonce    uint64                  `json:"nonce"`
	Value    argBigInt               `json:"value"`
	GasLimit uint64                  `json:"gasLimit"`
	GasPrice uint64                  `json:"gasPrice"`
	ChainID  uint64                  `json:"chainId"`
	RawTx    string                  `json:"rawTx"`
}

type Leaf struct {
	Balance argBigInt `json:"balance"`
	Nonce   uint64    `json:"nonce"`
}
