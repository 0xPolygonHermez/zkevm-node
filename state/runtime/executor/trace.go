package executor

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Trace struct {
	Context Context `json:"context"`
	Steps   []Step  `json:"step"`
}

type Context struct {
	Type     string `json:"type"`
	From     string `json:"from"`
	To       string `json:"to"`
	Input    string `json:"input"`
	Gas      string `json:"gas"`
	Value    string `json:"value"`
	Output   string `json:"output"`
	Nonce    uint64 `json:"nonce"`
	GasPrice string `json:"gasPrice"`
	ChainID  uint64 `json:"chainId"`
}

type Step struct {
	StateRoot string                      `json:"staterRoot"`
	Depth     uint32                      `json:"depth"`
	Pc        uint32                      `json:"pc"`
	Gas       string                      `json:"gas"`
	OpCode    string                      `json:"opcode"`
	Refund    string                      `json:"refund"`
	Op        byte                        `json:"op"`
	Error     string                      `json:"error"`
	Storage   map[common.Hash]common.Hash `json:"storage"`
	Step      uint32                      `json:"step"`
	Contract  Contract                    `json:"contract"`
	GasCost   uint64                      `json:"gasCost"`
	Stack     []*big.Int                  `json:"stack"`
	Memory    []byte                      `json:"memory"`
}

type Contract struct {
	Address string `json:"address"`
	Caller  uint32 `json:"caller"`
	Value   uint32 `json:"value"`
	Input   string `json:"input"`
}

type Tracer struct {
	Code string `json:"tracer"`
}
