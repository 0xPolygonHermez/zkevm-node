package state

import "github.com/ethereum/go-ethereum/common"

// Genesis contains the information to populate state on creation
type Genesis struct {
	Root         common.Hash
	Actions      []*GenesisAction
	Transactions []GenesisTx
}

// GenesisAction represents one of the values set on the SMT during genesis.
type GenesisAction struct {
	Address         string `json:"address"`
	Type            int    `json:"type"`
	StoragePosition string `json:"storagePosition"`
	Bytecode        string `json:"bytecode"`
	Key             string `json:"key"`
	Value           string `json:"value"`
	Root            string `json:"root"`
}

// GenesisTx represents the txs of the genesis
type GenesisTx struct {
	RawTx         string         `json:"rawTx"`
	Receipt       GenesisReceipt `json:"receipt"`
	CreateAddress common.Address `json:"createAddress"`
}

// GenesisReceipt represents the genesis receipt
type GenesisReceipt struct {
	Status  uint8           `json:"status"`
	GasUsed uint64          `json:"gasUsed"`
	Logs    [][]interface{} `json:"logs"`
}
