package vectors

import (
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

// StateTransitionTestCaseV2 holds the metadata needed to run a state transition test
type StateTransitionTestCaseEtrog struct {
	Description          string               `json:"Description"`
	Genesis              []GenesisEntityEtrog `json:"genesis"`
	ExpectedOldStateRoot string               `json:"expectedOldRoot"`
	ExpectedNewStateRoot string               `json:"expectedNewRoot"`
	ExpectedNewLeafs     map[string]LeafEtrog `json:"expectedNewLeafs"`
	Receipts             []TestReceipt        `json:"receipts"`
	GlobalExitRoot       string               `json:"globalExitRoot"`
	Txs                  []TxEtrog            `json:"txs"`
	OldAccInputHash      string               `json:"oldAccInputHash"`
	L1InfoRoot           string               `json:"l1InfoRoot"`
	TimestampLimit       string               `json:"timestampLimit"`
	BatchL2Data          string               `json:"batchL2Data"`
	BatchHashData        string               `json:"batchHashData"`
	ForkID               uint64               `json:"forkID"`
}

// LeafEtrog represents the state of a leaf in the merkle tree
type LeafEtrog struct {
	Balance         argBigInt         `json:"balance"`
	Nonce           string            `json:"nonce"`
	Storage         map[string]string `json:"storage"`
	IsSmartContract bool              `json:"isSmartContract"`
	Bytecode        string            `json:"bytecode"`
	HashBytecode    string            `json:"hashBytecode"`
	BytecodeLength  int               `json:"bytecodeLength"`
}

// GenesisEntity represents the state of an account or smart contract when the network
// starts
type GenesisEntityEtrog struct {
	Address         string            `json:"address"`
	PvtKey          *string           `json:"pvtKey"`
	Balance         argBigInt         `json:"balance"`
	Nonce           string            `json:"nonce"`
	Storage         map[string]string `json:"storage"`
	IsSmartContract bool              `json:"isSmartContract"`
	Bytecode        *string           `json:"bytecode"`
}

// TxEtrog represents a transactions that will be applied during the test
type TxEtrog struct {
	ID                uint       `json:"id"`
	From              string     `json:"from"`
	To                string     `json:"to"`
	Nonce             uint64     `json:"nonce"`
	Value             *big.Float `json:"value"`
	GasLimit          uint64     `json:"gasLimit"`
	GasPrice          *big.Float `json:"gasPrice"`
	ChainID           uint64     `json:"chainId"`
	RawTx             string     `json:"rawTx"`
	CustomRawTx       string     `json:"customRawTx"`
	Overwrite         Overwrite  `json:"overwrite"`
	EncodeInvalidData bool       `json:"encodeInvalidData"`
	Reason            string     `json:"reason"`
}

func GenerateGenesisActionsEtrog(genesis []GenesisEntityEtrog) []*state.GenesisAction {
	var genesisActions []*state.GenesisAction
	for _, genesisEntity := range genesis {

		if genesisEntity.Balance.String() != "0" {
			action := &state.GenesisAction{
				Address: genesisEntity.Address,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   genesisEntity.Balance.String(),
			}
			genesisActions = append(genesisActions, action)
		}

		if genesisEntity.Nonce != "" && genesisEntity.Nonce != "0" {
			action := &state.GenesisAction{
				Address: genesisEntity.Address,
				Type:    int(merkletree.LeafTypeNonce),
				Value:   genesisEntity.Nonce,
			}
			genesisActions = append(genesisActions, action)
		}

		if genesisEntity.IsSmartContract && genesisEntity.Bytecode != nil && *genesisEntity.Bytecode != "0x" {
			action := &state.GenesisAction{
				Address:  genesisEntity.Address,
				Type:     int(merkletree.LeafTypeCode),
				Bytecode: *genesisEntity.Bytecode,
			}
			genesisActions = append(genesisActions, action)
		}

		if genesisEntity.IsSmartContract && len(genesisEntity.Storage) > 0 {
			for storageKey, storageValue := range genesisEntity.Storage {
				action := &state.GenesisAction{
					Address:         genesisEntity.Address,
					Type:            int(merkletree.LeafTypeStorage),
					StoragePosition: storageKey,
					Value:           storageValue,
				}
				genesisActions = append(genesisActions, action)
			}
		}
	}

	return genesisActions
}
