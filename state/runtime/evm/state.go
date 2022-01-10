package evm

import (
	"errors"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/state/runtime"
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

// Instructions is the code of instructions

type state struct {
	ip   int
	code []byte
	tmp  []byte

	// host   runtime.Host
	msg    *runtime.Contract // change with msg
	config *runtime.ForksInTime

	// memory
	memory      []byte
	lastGasCost uint64

	// stack
	stack []*big.Int
	sp    int

	// remove later
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

func (c *state) push(val *big.Int) {
	c.push1().Set(val)
}

func (c *state) push1() *big.Int {
	if len(c.stack) > c.sp {
		c.sp++
		return c.stack[c.sp-1]
	}
	v := big.NewInt(0)
	c.stack = append(c.stack, v)
	c.sp++
	return v
}

func (c *state) stackAtLeast(n int) bool {
	return c.sp >= n
}

func (c *state) popHash() common.Hash {
	return common.BytesToHash(c.pop().Bytes())
}

func (c *state) popAddr() (common.Address, bool) {
	b := c.pop()
	if b == nil {
		return common.Address{}, false
	}

	return common.BytesToAddress(b.Bytes()), true
}

func (c *state) stackSize() int {
	return c.sp
}

func (c *state) top() *big.Int {
	if c.sp == 0 {
		return nil
	}
	return c.stack[c.sp-1]
}

func (c *state) pop() *big.Int {
	if c.sp == 0 {
		return nil
	}
	o := c.stack[c.sp-1]
	c.sp--
	return o
}

func (c *state) peekAt(n int) *big.Int {
	return c.stack[c.sp-n]
}

func (c *state) swap(n int) {
	c.stack[c.sp-1], c.stack[c.sp-n-1] = c.stack[c.sp-n-1], c.stack[c.sp-1]
}

func (c *state) get2(dst []byte, offset, length *big.Int) ([]byte, bool) {
	if length.Sign() == 0 {
		return nil, true
	}

	if !c.checkMemory(offset, length) {
		return nil, false
	}

	o := offset.Uint64()
	l := length.Uint64()

	dst = append(dst, c.memory[o:o+l]...)
	return dst, true
}

func (c *state) checkMemory(offset, size *big.Int) bool {
	if size.Sign() == 0 {
		return true
	}

	if !offset.IsUint64() || !size.IsUint64() {
		c.exit(errGasUintOverflow)
		return false
	}

	o := offset.Uint64()
	s := size.Uint64()

	if o > 0xffffffffe0 || s > 0xffffffffe0 {
		c.exit(errGasUintOverflow)
		return false
	}

	m := uint64(len(c.memory))
	newSize := o + s

	if m < newSize {
		w := (newSize + 31) / 32
		newCost := uint64(3*w + w*w/512)
		cost := newCost - c.lastGasCost
		c.lastGasCost = newCost

		if !c.consumeGas(cost) {
			c.exit(errOutOfGas)
			return false
		}

		// resize the memory
		c.memory = extendByteSlice(c.memory, int(w*32))
	}
	return true
}

// Run executes the virtual machine
func (s *state) Run() ([]byte, error) {
	var vmerr error

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
		inst.inst(s)

		// check if stack size exceeds the max size
		if s.sp > stackSize {
			s.exit(errStackOverflow)
			break
		}
		s.ip++
	}

	if err := s.err; err != nil {
		vmerr = err
	}
	return s.ret, vmerr
}

func extendByteSlice(b []byte, needLen int) []byte {
	b = b[:cap(b)]
	if n := needLen - cap(b); n > 0 {
		b = append(b, make([]byte, n)...)
	}
	return b[:needLen]
}
