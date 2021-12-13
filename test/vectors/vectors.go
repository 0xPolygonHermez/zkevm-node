package vectors

import (
	"math/big"
)

// StateTransitionTestCase holds the metadata needed to run a state transition test
type StateTransitionTestCase struct {
	ID                  uint   `json:"id"`
	Description         string `json:"description"`
	Arity               uint   `json:"arity"`
	ChanIDSequencer     uint64 `json:"chainIdSequencer"`
	SequencerAddress    string `json:"sequencerAddress"`
	SequencerPrivateKey string `json:"sequencerPvtKey"`

	GenesisAccounts  []GenesisAccount `json:"genesis"`
	ExpectedOldRoot  string           `json:"expectedOldRoot"`
	Txs              []Tx             `json:"txs"`
	ExpectedNewRoot  string           `json:"expectedNewRoot"`
	ExpectedNewLeafs map[string]Leaf  `json:"expectedNewLeafs"`
}

// GenesisAccount represents the state of an account when the network
// starts
type GenesisAccount struct {
	Address string    `json:"address"`
	PvtKey  string    `json:"pvtKey"`
	Balance argBigInt `json:"balance"`
	Nonce   string    `json:"nonce"`
}

// Tx represents a transactions that will be applied during the test
type Tx struct {
	From     string     `json:"from"`
	To       string     `json:"to"`
	Nonce    uint64     `json:"nonce"`
	Value    *big.Float `json:"value"`
	GasLimit uint64     `json:"gasLimit"`
	GasPrice *big.Float `json:"gasPrice"`
	ChainID  uint64     `json:"chainId"`
	RawTx    string     `json:"rawTx"`
}

// Leaf represents the state of a leaf in the merkle tree
type Leaf struct {
	Balance argBigInt `json:"balance"`
	Nonce   string    `json:"nonce"`
}
