package evm

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/hermeznetwork/hermez-core/state/runtime/fakevm"
	"github.com/hermeznetwork/hermez-core/state/runtime/instrumentation"
)

var statePool = sync.Pool{
	New: func() interface{} {
		return new(state)
	},
}

func acquireState() *state {
	return statePool.Get().(*state)
}

func releaseState(s *state) {
	s.reset()
	statePool.Put(s)
}

const stackSize = 1024

var (
	errOutOfGas              = runtime.ErrOutOfGas
	errStackUnderflow        = runtime.ErrStackUnderflow
	errStackOverflow         = runtime.ErrStackOverflow
	errRevert                = runtime.ErrExecutionReverted
	errGasUintOverflow       = errors.New("gas uint64 overflow")
	errWriteProtection       = errors.New("write protection")
	errInvalidJump           = errors.New("invalid jump destination")
	errOpCodeNotFound        = errors.New("opcode not found")
	errReturnDataOutOfBounds = errors.New("return data out of bounds")
)

type state struct {
	ip   int
	code []byte
	tmp  []byte

	host   runtime.Host
	msg    *runtime.Contract // change with msg
	config *runtime.ForksInTime

	// Instrumentation
	instrumented     bool
	storeDiff        *instrumentation.StoreDiff
	returnStructLogs []instrumentation.StructLog

	// Memory
	memory      []byte
	lastGasCost uint64

	// Stack
	stack []*big.Int
	sp    int

	// EVM
	evm *EVM

	err  error
	stop bool

	gas uint64

	// bitvec bitvec
	bitmap bitmap

	returnData []byte
	ret        []byte
}

func (s *state) reset() {
	s.sp = 0
	s.ip = 0
	s.gas = 0
	s.lastGasCost = 0
	s.stop = false
	s.err = nil

	// reset bitmap
	s.bitmap.reset()

	// reset memory
	for i := range s.memory {
		s.memory[i] = 0
	}

	s.tmp = s.tmp[:0]
	s.ret = s.ret[:0]
	s.code = s.code[:0]
	s.memory = s.memory[:0]
}

func (s *state) resetReturnData() {
	s.returnData = s.returnData[:0]
}

func (s *state) halt() {
	s.stop = true
}

func (s *state) exit(err error) {
	if err == nil {
		panic("cannot stop with none")
	}
	s.stop = true
	s.err = err
}

func (s *state) consumeGas(gas uint64) bool {
	if s.gas < gas {
		s.exit(errOutOfGas)
		return false
	}

	s.gas -= gas
	return true
}

func (s *state) push1() *big.Int {
	if len(s.stack) > s.sp {
		s.sp++
		return s.stack[s.sp-1]
	}
	v := big.NewInt(0)
	s.stack = append(s.stack, v)
	s.sp++

	return v
}

func (s *state) stackAtLeast(n int) bool {
	return s.sp >= n
}

func (s *state) popAddr() (common.Address, bool) {
	b := s.pop()
	if b == nil {
		return common.Address{}, false
	}

	return common.BytesToAddress(b.Bytes()), true
}

func (s *state) top() *big.Int {
	if s.sp == 0 {
		return nil
	}
	return s.stack[s.sp-1]
}

func (s *state) pop() *big.Int {
	if s.sp == 0 {
		return nil
	}
	o := s.stack[s.sp-1]
	s.sp--
	return o
}

func (s *state) peekAt(n int) *big.Int {
	return s.stack[s.sp-n]
}

func (s *state) swap(n int) {
	s.stack[s.sp-1], s.stack[s.sp-n-1] = s.stack[s.sp-n-1], s.stack[s.sp-1]
}

func (s *state) get2(dst []byte, offset, length *big.Int) ([]byte, bool) {
	if length.Sign() == 0 {
		return nil, true
	}

	if !s.checkMemory(offset, length) {
		return nil, false
	}

	o := offset.Uint64()
	l := length.Uint64()

	dst = append(dst, s.memory[o:o+l]...)
	return dst, true
}

func (s *state) checkMemory(offset, size *big.Int) bool {
	if size.Sign() == 0 {
		return true
	}

	if !offset.IsUint64() || !size.IsUint64() {
		s.exit(errGasUintOverflow)
		return false
	}

	o := offset.Uint64()
	sz := size.Uint64()

	if o > 0xffffffffe0 || sz > 0xffffffffe0 {
		s.exit(errGasUintOverflow)
		return false
	}

	m := uint64(len(s.memory))
	newSize := o + sz

	if m < newSize {
		w := (newSize + 31) / 32 //nolint:gomnd
		newCost := uint64(3*w + w*w/512)
		cost := newCost - s.lastGasCost
		s.lastGasCost = newCost

		if !s.consumeGas(cost) {
			s.exit(errOutOfGas)
			return false
		}

		// resize the memory
		s.memory = extendByteSlice(s.memory, int(w*32)) //nolint:gomnd
	}
	return true
}

// Run executes the virtual machine
func (s *state) Run(ctx context.Context, contract instrumentation.Contract) ([]byte, []instrumentation.StructLog, instrumentation.ExecutorTrace, error) {
	var vmerr error
	var structLogs []instrumentation.StructLog
	var executorTrace instrumentation.ExecutorTrace
	var steps []instrumentation.Step

	codeSize := len(s.code)
	for !s.stop {
		if s.ip >= codeSize {
			s.halt()
			break
		}

		op := OpCode(s.code[s.ip])

		inst := dispatchTable[op]
		if inst.inst == nil {
			s.exit(errOpCodeNotFound)
			break
		}

		// check if the depth of the stack is enough for the instruction
		if s.sp < inst.stack {
			s.exit(errStackUnderflow)
			break
		}
		// consume the gas of the instruction
		if !s.consumeGas(inst.gas) {
			s.exit(errOutOfGas)
			break
		}

		// execute the instruction
		inst.inst(ctx, s)

		// check if stack size exceeds the max size
		if s.sp > stackSize {
			s.exit(errStackOverflow)
			break
		}

		if s.instrumented {
			// Debug
			structLog := instrumentation.StructLog{
				Pc:         uint64(s.ip),
				Op:         op.String(),
				Gas:        s.gas,
				GasCost:    inst.gas,
				Memory:     s.memory,
				MemorySize: len(s.memory),
				Stack:      s.stack,
				ReturnData: s.returnData,
				Depth:      s.msg.Depth,
				Err:        s.err,
			}

			structLogs = append(structLogs, structLog)

			if op == CREATE || op == CREATE2 || op == CALL || op == CALLCODE || op == DELEGATECALL || op == STATICCALL {
				for i := range s.returnStructLogs {
					structLogs = append(structLogs, s.returnStructLogs[i])
				}
			}

			// Executor trace
			stack := bigArrayToStringArray(s.stack)
			memory := memoryToStringArray(s.memory)

			step := instrumentation.Step{
				Contract:   contract,
				StateRoot:  "0x" + hex.EncodeToString(s.host.GetStateRoot(ctx)),
				Depth:      s.msg.Depth,
				Pc:         uint64(s.ip),
				Gas:        fmt.Sprint(s.gas),
				OpCode:     op.String(),
				GasCost:    fmt.Sprint(inst.gas),
				Refund:     "0",
				Op:         "0x" + hex.EncodeToString([]byte{byte(op)}),
				Stack:      stack,
				Memory:     memory,
				ReturnData: "0x" + hex.EncodeToString(s.ret),
			}

			if s.err != nil {
				step.Error = s.err.Error()
			}

			steps = append(steps, step)
		}

		s.ip++
	}

	if err := s.err; err != nil {
		vmerr = err
	}

	executorTrace.Steps = steps

	return s.ret, structLogs, executorTrace, vmerr
}

func (s *state) inStaticCall() bool {
	return s.msg.Static
}

func (s *state) validJumpdest(dest *big.Int) bool {
	udest := dest.Uint64()
	if dest.BitLen() >= 63 || udest >= uint64(len(s.code)) {
		return false
	}
	return s.bitmap.isSet(uint(udest))
}

func extendByteSlice(b []byte, needLen int) []byte {
	b = b[:cap(b)]
	if n := needLen - cap(b); n > 0 {
		b = append(b, make([]byte, n)...)
	}
	return b[:needLen]
}

func bigToHash(b *big.Int) common.Hash {
	return common.BytesToHash(b.Bytes())
}

func bigArrayToStringArray(b []*big.Int) []string {
	s := []string{}

	for _, bn := range b {
		s = append(s, fmt.Sprintf("0x%x", bn))
	}

	return s
}

func memoryToStringArray(memory []byte) []string {
	numElements := len(memory) / fakevm.MemoryItemSize
	s := []string{}

	for x := 0; x < numElements; x++ {
		s = append(s, hex.EncodeToString(memory[x*fakevm.MemoryItemSize:(x*fakevm.MemoryItemSize)+32]))
	}

	return s
}
