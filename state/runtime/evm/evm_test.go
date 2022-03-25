package evm

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"gotest.tools/assert"
)

func newMockContract(value *big.Int, gas uint64, code []byte) *runtime.Contract {
	return runtime.NewContract(
		1,
		common.BytesToAddress([]byte("0x0000000000000000000000000000000000000000")),
		common.BytesToAddress([]byte("0x0000000000000000000000000000000000000000")),
		common.BytesToAddress([]byte("0x0000000000000000000000000000000000000000")),
		value,
		gas,
		code,
	)
}

// mockHost is a struct which meets the requirements of runtime.Host interface but throws panic in each methods
// we don't test all opcodes in this test
type mockHost struct{}

func (m *mockHost) AccountExists(ctx context.Context, addr common.Address) bool {
	panic("Not implemented in tests")
}

func (m *mockHost) GetStorage(ctx context.Context, addr common.Address, key common.Hash) common.Hash {
	panic("Not implemented in tests")
}

func (m *mockHost) SetStorage(cxt context.Context, addr common.Address, key *big.Int, value *big.Int, config *runtime.ForksInTime) runtime.StorageStatus {
	panic("Not implemented in tests")
}

func (m *mockHost) GetBalance(ctx context.Context, addr common.Address) *big.Int {
	panic("Not implemented in tests")
}

func (m *mockHost) GetCodeSize(ctx context.Context, addr common.Address) int {
	panic("Not implemented in tests")
}

func (m *mockHost) GetCodeHash(ctx context.Context, addr common.Address) common.Hash {
	panic("Not implemented in tests")
}

func (m *mockHost) GetCode(ctx context.Context, addr common.Address) []byte {
	panic("Not implemented in tests")
}

func (m *mockHost) Selfdestruct(ctx context.Context, addr common.Address, beneficiary common.Address) {
	panic("Not implemented in tests")
}

func (m *mockHost) GetTxContext() runtime.TxContext {
	panic("Not implemented in tests")
}

func (m *mockHost) GetBlockHash(number int64) common.Hash {
	panic("Not implemented in tests")
}

func (m *mockHost) EmitLog(addr common.Address, topics []common.Hash, data []byte) {
	panic("Not implemented in tests")
}

func (m *mockHost) Callx(context.Context, *runtime.Contract, runtime.Host) *runtime.ExecutionResult {
	panic("Not implemented in tests")
}

func (m *mockHost) Empty(ctx context.Context, addr common.Address) bool {
	panic("Not implemented in tests")
}

func (m *mockHost) GetNonce(ctx context.Context, addr common.Address) uint64 {
	panic("Not implemented in tests")
}

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		value    *big.Int
		gas      uint64
		code     []byte
		config   *runtime.ForksInTime
		expected *runtime.ExecutionResult
	}{
		{
			name:  "should succeed because of no codes",
			value: big.NewInt(0),
			gas:   5000,
			code:  []byte{},
			expected: &runtime.ExecutionResult{
				ReturnValue: nil,
				GasLeft:     5000,
			},
		},
		{
			name:  "should succeed and return result",
			value: big.NewInt(0),
			gas:   5000,
			code: []byte{
				PUSH1, 0x01, PUSH1, 0x02, ADD,
				PUSH1, 0x00, MSTORE8,
				PUSH1, 0x01, PUSH1, 0x00, RETURN,
			},
			expected: &runtime.ExecutionResult{
				ReturnValue: []uint8{0x03},
				GasLeft:     4976,
			},
		},
		{
			name:  "should fail and consume all gas by error",
			value: big.NewInt(0),
			gas:   5000,
			// ADD will be failed by stack underflow
			code: []byte{ADD},
			expected: &runtime.ExecutionResult{
				ReturnValue: nil,
				GasLeft:     0,
				Err:         errStackUnderflow,
			},
		},
		{
			name:  "should fail by REVERT and return remaining gas at that time",
			value: big.NewInt(0),
			gas:   5000,
			// Stack size and offset for return value first
			code: []byte{PUSH1, 0x00, PUSH1, 0x00, REVERT},
			config: &runtime.ForksInTime{
				Byzantium: true,
			},
			expected: &runtime.ExecutionResult{
				ReturnValue: nil,
				// gas consumed for 2 push1 ops
				GasLeft: 4994,
				Err:     errRevert,
			},
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evm := NewEVM(false)
			contract := newMockContract(tt.value, tt.gas, tt.code)
			host := &mockHost{}
			config := tt.config
			if config == nil {
				config = &runtime.ForksInTime{}
			}
			res := evm.Run(ctx, contract, host, config)
			assert.Equal(t, tt.expected.GasUsed, res.GasUsed)
			assert.Equal(t, tt.expected.GasLeft, res.GasLeft)
			assert.Equal(t, tt.expected.Err, res.Err)
		})
	}
}
