package instrumentation

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// FullTrace contents executor call trace steps.
type FullTrace struct {
	Context Context `json:"context"`
	Steps   []Step  `json:"steps"`
}

// Context is the trace context.
type Context struct {
	Type         string      `json:"type"`
	From         string      `json:"from"`
	To           string      `json:"to"`
	Input        []byte      `json:"input"`
	Gas          uint64      `json:"gas"`
	Value        *big.Int    `json:"value"`
	Output       []byte      `json:"output"`
	Nonce        uint64      `json:"nonce"`
	GasPrice     string      `json:"gasPrice"`
	OldStateRoot common.Hash `json:"oldStateRoot"`
	Time         uint64      `json:"time"`
	GasUsed      uint64      `json:"gasUsed"`
}

// Step is a trace step.
type Step struct {
	StateRoot    common.Hash                 `json:"stateRoot"`
	Depth        int                         `json:"depth"`
	Pc           uint64                      `json:"pc"`
	Gas          uint64                      `json:"gas"`
	OpCode       string                      `json:"opcode"`
	Refund       uint64                      `json:"refund"`
	Op           uint64                      `json:"op"`
	Error        error                       `json:"error"`
	Contract     Contract                    `json:"contract"`
	GasCost      uint64                      `json:"gasCost"`
	Stack        []*big.Int                  `json:"stack"`
	Memory       []byte                      `json:"memory"`
	MemorySize   uint32                      `json:"memorySize"`
	MemoryOffset uint32                      `json:"memoryOffset"`
	ReturnData   []byte                      `json:"returnData"`
	Storage      map[common.Hash]common.Hash `json:"storage"`
}

// Contract represents a contract in the trace.
type Contract struct {
	Address common.Address `json:"address"`
	Caller  common.Address `json:"caller"`
	Value   *big.Int       `json:"value"`
	Input   []byte         `json:"input"`
	Gas     uint64         `json:"gas"`
}

// Tracer represents the executor tracer.
type Tracer struct {
	Code string `json:"tracer"`
}

type InternalTxContext struct {
	OpCode       string
	RemainingGas uint64
}
