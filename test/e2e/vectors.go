package e2e

import (
	"github.com/ethereum/go-ethereum/common"
)

// StateTransitionVector contains test cases for state transitions
type StateTransitionVector struct {
	StateTests []StateTest `json:"tests"`
}

// StateTest holds the metadata needed to run a state transition test
type StateTest struct {
	ID               uint                    `json:"id"`
	Description      string                  `json:"description"`
	Arity            uint                    `json:"arity"`
	ChanIDSequencer  uint64                  `json:"chainIdSequencer"`
	SequencerAddress common.MixedcaseAddress `json:"sequencerAddress"`

	GenesisAccounts  []GenesisAccount        `json:"genesis"`
	ExpectedOldRoot  []byte                  `json:"expectedOldRoot"`
	Txs              []Tx                    `json:"txs"`
	ExpectedNewRoot  []byte                  `json:"expectedNewRoot"`
	ExpectedNewLeafs map[common.Address]Leaf `json:"expectedNewLeafs"`
}

// GenesisAccount represents the state of an account when the network
// starts
type GenesisAccount struct {
	Address common.MixedcaseAddress `json:"address"`
	PvtKey  string                  `json:"pvtKey"`
	Balance argBigInt               `json:"balance"`
	Nonce   uint64                  `json:"nonce"`
}

// Tx represents a transactions that will be applied during the test
type Tx struct {
	From     common.MixedcaseAddress `json:"from"`
	To       common.MixedcaseAddress `json:"to"`
	Nonce    uint64                  `json:"nonce"`
	Value    argBigInt               `json:"value"`
	GasLimit uint64                  `json:"gasLimit"`
	GasPrice argBigInt               `json:"gasPrice"`
	ChainID  uint64                  `json:"chainId"`
	RawTx    string                  `json:"rawTx"`
}

// Leaf represents the state of a leaf in the merkle tree
type Leaf struct {
	Balance argBigInt `json:"balance"`
	Nonce   uint64    `json:"nonce"`
}
