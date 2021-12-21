package vectors

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// StateTransitionTestCase holds the metadata needed to run a state transition test
type StateTransitionTestCase struct {
	ID                  uint   `json:"id"`
	Description         string `json:"description"`
	Arity               uint8  `json:"arity"`
	ChainIDSequencer    uint64 `json:"chainIdSequencer"`
	DefaultChainID      uint64 `json:"defaultChainId"`
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
	ID                uint       `json:"id"`
	From              string     `json:"from"`
	To                string     `json:"to"`
	Nonce             uint64     `json:"nonce"`
	Value             *big.Float `json:"value"`
	GasLimit          uint64     `json:"gasLimit"`
	GasPrice          *big.Float `json:"gasPrice"`
	ChainID           uint64     `json:"chainId"`
	RawTx             string     `json:"rawTx"`
	Overwrite         Overwrite  `json:"overwrite"`
	EncodeInvalidData bool       `json:"encodeInvalidData"`
	Reason            string     `json:"reason"`
}

// Leaf represents the state of a leaf in the merkle tree
type Leaf struct {
	Balance argBigInt `json:"balance"`
	Nonce   string    `json:"nonce"`
}

// Overwrite is used by Protocol team for testing
type Overwrite struct {
	S string `json:"s"`
}

// CallDataTestCase holds the metadata needed to run a etherman test
type CallDataTestCase struct {
	ID  uint `json:"id"`
	Txs []Tx `json:"txs"`

	BatchL2Data   string      `json:"batchL2Data"`
	BatchHashData common.Hash `json:"batchHashData"`
	MaticAmount   string      `json:"maticAmount"`
	FullCallData  string      `json:"fullCallData"`
}
