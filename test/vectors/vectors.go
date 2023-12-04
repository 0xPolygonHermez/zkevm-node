package vectors

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// StateTransitionTestCase holds the metadata needed to run a state transition test
type StateTransitionTestCase struct {
	ID                  uint   `json:"id"`
	Description         string `json:"description"`
	ChainIDSequencer    uint64 `json:"chainIdSequencer"`
	SequencerAddress    string `json:"sequencerAddress"`
	SequencerPrivateKey string `json:"sequencerPvtKey"`

	GenesisAccounts       []GenesisAccount       `json:"genesis"`
	GenesisSmartContracts []GenesisSmartContract `json:"genesisSC"`
	ExpectedOldRoot       string                 `json:"expectedOldRoot"`
	Txs                   []Tx                   `json:"txs"`
	ExpectedNewRoot       string                 `json:"expectedNewRoot"`
	ExpectedNewLeafs      map[string]Leaf        `json:"expectedNewLeafs"`
	Receipts              []TestReceipt          `json:"receipts"`
	GlobalExitRoot        string                 `json:"globalExitRoot"`
}

// GenesisAccount represents the state of an account when the network
// starts
type GenesisAccount struct {
	Address string    `json:"address"`
	PvtKey  string    `json:"pvtKey"`
	Balance argBigInt `json:"balance"`
	Nonce   string    `json:"nonce"`
}

// GenesisSmartContract represents the smart contract to init when the network starts
type GenesisSmartContract struct {
	Address string `json:"address"`
	Code    string `json:"bytecode"`
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

// TxEventsSendBatchTestCase holds the metadata needed to run a etherman test
type TxEventsSendBatchTestCase struct {
	ID  uint `json:"id"`
	Txs []Tx `json:"txs"`

	BatchL2Data   string      `json:"batchL2Data"`
	BatchHashData common.Hash `json:"batchHashData"`
	PolAmount   string        `json:"polAmount"`
	FullCallData  string      `json:"fullCallData"`
}

// TestReceipt holds the metadata needed to run the receipt tests
type TestReceipt struct {
	TxID    uint    `json:"txId"`
	Receipt Receipt `json:"receipt"`
}

// Receipt is the receipt used for receipts tests
type Receipt struct {
	TransactionHash    string `json:"transactionHash"`
	TransactionIndex   uint   `json:"transactionIndex"`
	BlockNumber        uint64 `json:"blockNumber"`
	From               string `json:"from"`
	To                 string `json:"to"`
	CumulativeGastUsed uint64 `json:"cumulativeGasUsed"`
	GasUsedForTx       uint64 `json:"gasUsedForTx"`
	ContractAddress    string `json:"contractAddress"`
	Logs               uint64 `json:"logs"`
	LogsBloom          uint64 `json:"logsBloom"`
	Status             uint64 `json:"status"`
	BlockHash          string `json:"blockHash"`
}
