package testvector

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/afero"
)

// E2E contains all the test vectors.
type E2E struct {
	Items []*E2EItem
}

// E2EItem contains an end-to-end test vector.
type E2EItem struct {
	BatchL2Data    string
	GlobalExitRoot string
	Traces         *Traces
	GenesisRaw     []*state.GenesisAction
}

// Traces represents executor processing traces.
type Traces struct {
	BatchHash     string
	OldStateRoot  string `json:"old_state_root"`
	GlobalHash    string
	NumBatch      uint64
	Timestamp     uint64
	SequencerAddr string

	*ProcessBatchResponse
}

// ProcessBatchResponse contains information about the executor response to a
// ProcessBatch request.
type ProcessBatchResponse struct {
	CumulativeGasUsed   string                        `json:"cumulative_gas_used,omitempty"`
	Responses           []*ProcessTransactionResponse `json:"responses,omitempty"`
	NewStateRoot        string                        `json:"new_state_root,omitempty"`
	NewLocalExitRoot    string                        `json:"new_local_exit_root,omitempty"`
	CntKeccakHashes     uint32                        `json:"cnt_keccak_hashes,omitempty"`
	CntPoseidonHashes   uint32                        `json:"cnt_poseidon_hashes,omitempty"`
	CntPoseidonPaddings uint32                        `json:"cnt_poseidon_paddings,omitempty"`
	CntMemAligns        uint32                        `json:"cnt_mem_aligns,omitempty"`
	CntArithmetics      uint32                        `json:"cnt_arithmetics,omitempty"`
	CntBinaries         uint32                        `json:"cnt_binaries,omitempty"`
	CntSteps            uint32                        `json:"cnt_steps,omitempty"`
}

// ProcessTransactionResponse contains information about the executor response to a
// transaction execution.
type ProcessTransactionResponse struct {
	TxHash                 string `json:"tx_hash"`
	Type                   uint32
	GasLeft                string `json:"gas_left"`
	GasUsed                string `json:"gas_used"`
	GasRefunded            string `json:"gas_refunded"`
	StateRoot              string `json:"state_root"`
	Logs                   []*Log
	UnprocessedTransaction bool       `json:"unprocessed_transaction,omitempty"`
	CallTrace              *CallTrace `json:"call_trace,omitempty"`
}

// Log represent logs emitted by LOG opcode.
type Log struct {
	Data        []string
	Topics      []string
	Address     string
	BatchNumber uint64 `json:"batch_number"`
	TxHash      string `json:"tx_hash"`
	TxIndex     uint32 `json:"tx_index"`
	BatchHash   string `json:"batch_hash"`
	Index       uint32
}

// CallTrace represents the batch call trace.
type CallTrace struct {
	Context *TransactionContext `json:"context,omitempty"`
	Steps   []*TransactionStep  `json:"steps,omitempty"`
}

// TransactionContext represents a transaction's context.
type TransactionContext struct {
	From          string `json:"from,omitempty"`
	To            string `json:"to,omitempty"`
	Type          string `json:"type,omitempty"`
	Data          string `json:"data,omitempty"`
	Gas           string `json:"gas,omitempty"`
	Value         string `json:"value,omitempty"`
	Batch         string `json:"batch,omitempty"`
	Output        string `json:"output,omitempty"`
	GasUsed       string `json:"gas_used,omitempty"`
	ExecutionTime string `json:"execution_time,omitempty"`
	OldStateRoot  string `json:"old_state_root,omitempty"`
	GasPrice      string `json:"gasPrice,omitempty"`
}

// TransactionStep represents a transaction's step.
type TransactionStep struct {
	Depth        uint32 `json:"depth,omitempty"`
	Pc           uint64 `json:"pc,omitempty"`
	RemainingGas string `json:"remaining_gas,omitempty"`
	OpCode       string
	GasRefund    string `json:"gas_refund,omitempty"`
	Op           string `json:"op,omitempty"`
	Error        string
	StateRoot    string    `json:"state_root"`
	Contract     *Contract `json:"contract,omitempty"`
	ReturnData   []string  `json:"return_data,omitempty"`
	GasCost      string    `json:"gas_cost"`
	Stack        []string
	Memory       []string
}

// Contract contains information about SCs executed in a batch.
type Contract struct {
	Address string `json:"address,omitempty"`
	Caller  string `json:"caller,omitempty"`
	Value   string `json:"value,omitempty"`
	Data    string `json:"data,omitempty"`
	Gas     string `json:"gas,omitempty"`
}

// Container is a wrapper for test vectors.
type Container struct {
	E2E *E2E
}

// NewContainer is the Container constructor.
func NewContainer(testVectorPath string, aferoFs afero.Fs) (*Container, error) {
	e2e, err := getE2E(testVectorPath, aferoFs)
	if err != nil {
		return nil, err
	}

	return &Container{
		E2E: e2e,
	}, nil
}

// FindSMTValue searches for the given key on all the genesisRaw items present,
// checking also that the given root was the root returned by the previous item.
// If both the value and the root of the previous item match it returns the
// associated value and new root.
func (c *Container) FindSMTValue(inputKey, oldRoot string) (value, newRoot string, err error) {
	zero := common.HexToHash("").String()
	var lastValue string
	for _, item := range c.E2E.Items {
		for index, action := range item.GenesisRaw {
			if action.Key == inputKey {
				if index > 0 && oldRoot == item.GenesisRaw[index-1].Root ||
					index == 0 && oldRoot == zero {
					return item.GenesisRaw[index].Value, item.GenesisRaw[index].Root, nil
				} else {
					lastValue = item.GenesisRaw[index].Value
				}
			}
		}
		if len(item.GenesisRaw) > 0 &&
			oldRoot == item.GenesisRaw[len(item.GenesisRaw)-1].Root &&
			lastValue != "" {
			return lastValue, oldRoot, nil
		}
	}
	return "", "", fmt.Errorf("key %q not found for oldRoot %q", inputKey, oldRoot)
}

// FindBytecode searches for the given key on the value fields of all the
// genesisRaw items present and returns the associated bytecode field on match.
func (c *Container) FindBytecode(inputKey string) (bytecode string, err error) {
	for _, item := range c.E2E.Items {
		for index, action := range item.GenesisRaw {
			if action.Value == inputKey && action.Bytecode != "" {
				return item.GenesisRaw[index].Bytecode, nil
			}
		}
	}
	return "", fmt.Errorf("bytecode for key %q not found", inputKey)
}

// FindProcessBatchResponse searches for the responses to a process batch
// request identified by tge batch L2 data.
func (c *Container) FindProcessBatchResponse(batchL2Data string) (*ProcessBatchResponse, error) {
	for _, item := range c.E2E.Items {
		if strings.Replace(item.BatchL2Data, "0x", "", -1) == strings.Replace(batchL2Data, "0x", "", -1) {
			return item.Traces.ProcessBatchResponse, nil
		}
	}
	return nil, fmt.Errorf("ProcessBatchResponse for batchL2Data %q not found", batchL2Data)
}

func getE2E(testVectorPath string, aferoFs afero.Fs) (*E2E, error) {
	e2e := &E2E{}

	err := afero.Walk(aferoFs, testVectorPath, func(wpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info == nil || info.IsDir() {
			return nil
		}
		e2eFile, err := getE2EFile(wpath, aferoFs)
		if err != nil {
			return err
		}
		e2e.Items = append(e2e.Items, e2eFile.Items...)

		return nil
	})

	return e2e, err
}

func getE2EFile(filePath string, aferoFs afero.Fs) (*E2E, error) {
	contents, err := afero.ReadFile(aferoFs, filePath)
	if err != nil {
		return nil, err
	}
	var e2e E2E
	if err := json.Unmarshal(contents, &e2e.Items); err != nil {
		return nil, err
	}
	return &e2e, nil
}
