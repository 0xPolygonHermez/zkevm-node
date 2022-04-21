package evm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/bits"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/state/helper"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/hermeznetwork/hermez-core/state/runtime/instrumentation"
)

type instruction func(ctx context.Context, s *state)

const (
	wSize = 32
)

var (
	zero     = big.NewInt(0)
	one      = big.NewInt(1)
	wordSize = big.NewInt(wSize)
)

func opAdd(ctx context.Context, s *state) {
	a := s.pop()
	b := s.top()

	b.Add(a, b)
	toU256(b)
}

func opMul(ctx context.Context, s *state) {
	a := s.pop()
	b := s.top()

	b.Mul(a, b)
	toU256(b)
}

func opSub(ctx context.Context, s *state) {
	a := s.pop()
	b := s.top()

	b.Sub(a, b)
	toU256(b)
}

func opDiv(ctx context.Context, s *state) {
	a := s.pop()
	b := s.top()

	if b.Sign() == 0 {
		// division by zero
		b.Set(zero)
	} else {
		b.Div(a, b)
		toU256(b)
	}
}

func opSDiv(ctx context.Context, s *state) {
	a := to256(s.pop())
	b := to256(s.top())

	if b.Sign() == 0 {
		// division by zero
		b.Set(zero)
	} else {
		neg := a.Sign() != b.Sign()
		b.Div(a.Abs(a), b.Abs(b))
		if neg {
			b.Neg(b)
		}
		toU256(b)
	}
}

func opMod(ctx context.Context, s *state) {
	a := s.pop()
	b := s.top()

	if b.Sign() == 0 {
		// division by zero
		b.Set(zero)
	} else {
		b.Mod(a, b)
		toU256(b)
	}
}

func opSMod(ctx context.Context, s *state) {
	a := to256(s.pop())
	b := to256(s.top())

	if b.Sign() == 0 {
		// division by zero
		b.Set(zero)
	} else {
		neg := a.Sign() < 0
		b.Mod(a.Abs(a), b.Abs(b))
		if neg {
			b.Neg(b)
		}
		toU256(b)
	}
}

var bigPool = sync.Pool{
	New: func() interface{} {
		return new(big.Int)
	},
}

func acquireBig() *big.Int {
	return bigPool.Get().(*big.Int)
}

func releaseBig(b *big.Int) {
	bigPool.Put(b)
}

func opExp(ctx context.Context, s *state) {
	x := s.pop()
	y := s.top()

	var gas uint64
	if s.config.EIP158 {
		gas = 50
	} else {
		gas = 10
	}
	gasCost := uint64((y.BitLen()+7)/8) * gas //nolint:gomnd
	if !s.consumeGas(gasCost) {
		return
	}

	z := acquireBig().Set(one)

	// https://www.programminglogic.com/fast-exponentiation-algorithms/
	for _, d := range y.Bits() {
		for i := 0; i < _W; i++ {
			if d&1 == 1 {
				toU256(z.Mul(z, x))
			}
			d >>= 1
			toU256(x.Mul(x, x))
		}
	}
	y.Set(z)
	releaseBig(z)
}

func opAddMod(ctx context.Context, s *state) {
	a := s.pop()
	b := s.pop()
	z := s.top()

	if z.Sign() == 0 {
		// division by zero
		z.Set(zero)
	} else {
		a = a.Add(a, b)
		z = z.Mod(a, z)
		toU256(z)
	}
}

func opMulMod(ctx context.Context, s *state) {
	a := s.pop()
	b := s.pop()
	z := s.top()

	if z.Sign() == 0 {
		// division by zero
		z.Set(zero)
	} else {
		a = a.Mul(a, b)
		z = z.Mod(a, z)
		toU256(z)
	}
}

func opAnd(ctx context.Context, s *state) {
	a := s.pop()
	b := s.top()

	b.And(a, b)
}

func opOr(ctx context.Context, s *state) {
	a := s.pop()
	b := s.top()

	b.Or(a, b)
}

func opXor(ctx context.Context, s *state) {
	a := s.pop()
	b := s.top()

	b.Xor(a, b)
}

var opByteMask = big.NewInt(255) //nolint:gomnd

func opByte(ctx context.Context, s *state) {
	x := s.pop()
	y := s.top()

	indx := x.Int64()
	if indx > 31 { //nolint:gomnd
		y.Set(zero)
	} else {
		sh := (31 - indx) * 8 //nolint:gomnd
		y.Rsh(y, uint(sh))
		y.And(y, opByteMask)
	}
}

func opNot(ctx context.Context, s *state) {
	a := s.top()

	a.Not(a)
	toU256(a)
}

func opIsZero(ctx context.Context, s *state) {
	a := s.top()

	if a.Sign() == 0 {
		a.Set(one)
	} else {
		a.Set(zero)
	}
}

func opEq(ctx context.Context, s *state) {
	a := s.pop()
	b := s.top()

	if a.Cmp(b) == 0 {
		b.Set(one)
	} else {
		b.Set(zero)
	}
}

func opLt(ctx context.Context, s *state) {
	a := s.pop()
	b := s.top()

	if a.Cmp(b) < 0 {
		b.Set(one)
	} else {
		b.Set(zero)
	}
}

func opGt(ctx context.Context, s *state) {
	a := s.pop()
	b := s.top()

	if a.Cmp(b) > 0 {
		b.Set(one)
	} else {
		b.Set(zero)
	}
}

func opSlt(ctx context.Context, s *state) {
	a := to256(s.pop())
	b := to256(s.top())

	if a.Cmp(b) < 0 {
		b.Set(one)
	} else {
		b.Set(zero)
	}
}

func opSgt(ctx context.Context, s *state) {
	a := to256(s.pop())
	b := to256(s.top())

	if a.Cmp(b) > 0 {
		b.Set(one)
	} else {
		b.Set(zero)
	}
}

func opSignExtension(ctx context.Context, s *state) {
	ext := s.pop()
	x := s.top()

	if ext.Cmp(wordSize) > 0 {
		return
	}
	if x == nil {
		return
	}

	bit := uint(ext.Uint64()*8 + 7) //nolint:gomnd

	mask := acquireBig().Set(one)
	mask.Lsh(mask, bit)
	mask.Sub(mask, one)

	if x.Bit(int(bit)) > 0 {
		mask.Not(mask)
		x.Or(x, mask)
	} else {
		x.And(x, mask)
	}

	toU256(x)
	releaseBig(mask)
}

func equalOrOverflowsUint256(b *big.Int) bool {
	return b.BitLen() > 8 //nolint:gomnd
}

func opShl(ctx context.Context, s *state) {
	if !s.config.Constantinople {
		s.exit(errOpCodeNotFound)
		return
	}

	shift := s.pop()
	value := s.top()

	if equalOrOverflowsUint256(shift) {
		value.Set(zero)
	} else {
		value.Lsh(value, uint(shift.Uint64()))
		toU256(value)
	}
}

func opShr(ctx context.Context, s *state) {
	if !s.config.Constantinople {
		s.exit(errOpCodeNotFound)
		return
	}

	shift := s.pop()
	value := s.top()

	if equalOrOverflowsUint256(shift) {
		value.Set(zero)
	} else {
		value.Rsh(value, uint(shift.Uint64()))
		toU256(value)
	}
}

func opSar(ctx context.Context, s *state) {
	if !s.config.Constantinople {
		s.exit(errOpCodeNotFound)
		return
	}

	shift := s.pop()
	value := to256(s.top())

	if equalOrOverflowsUint256(shift) {
		if value.Sign() >= 0 {
			value.Set(zero)
		} else {
			value.Set(tt256m1)
		}
	} else {
		value.Rsh(value, uint(shift.Uint64()))
		toU256(value)
	}
}

// memory operations

var bufPool = sync.Pool{
	New: func() interface{} {
		// Store pointer to avoid heap allocation in caller
		// Please check SA6002 in StaticCheck for details
		buf := make([]byte, 128)
		return &buf
	},
}

func opMload(ctx context.Context, s *state) {
	offset := s.pop()

	var ok bool
	s.tmp, ok = s.get2(s.tmp[:0], offset, wordSize)
	if !ok {
		return
	}
	s.push1().SetBytes(s.tmp)
}

var (
	_W = bits.UintSize
	_S = _W / 8 //nolint:gomnd
)

func opMStore(ctx context.Context, s *state) {
	offset := s.pop()
	val := s.pop()

	if s.instrumented {
		diff := instrumentation.MemoryDiff{
			Offset: offset.Uint64(),
			Data:   val.Uint64(),
		}
		s.memDiff = &diff
	}

	if !s.checkMemory(offset, wordSize) {
		return
	}

	o := offset.Uint64()
	buf := s.memory[o : o+32]

	i := 32

	// convert big.int to bytes
	// https://golang.org/src/math/big/nat.go#L1284
	for _, d := range val.Bits() {
		for j := 0; j < _S; j++ {
			i--
			buf[i] = byte(d)
			d >>= 8
		}
	}

	// fill the rest of the slot with zeros
	for i > 0 {
		i--
		buf[i] = 0
	}
}

func opMStore8(ctx context.Context, s *state) {
	offset := s.pop()
	val := s.pop()

	if s.instrumented {
		diff := instrumentation.MemoryDiff{
			Offset: offset.Uint64(),
			Data:   val.Uint64(),
		}
		s.memDiff = &diff
	}

	if !s.checkMemory(offset, one) {
		return
	}
	s.memory[offset.Uint64()] = byte(val.Uint64() & 0xff) //nolint:gomnd
}

// --- storage ---

func opSload(ctx context.Context, s *state) {
	loc := s.top()

	var gas uint64
	if s.config.Istanbul {
		// eip-1884
		gas = 800
	} else if s.config.EIP150 {
		gas = 200
	} else {
		gas = 50
	}
	if !s.consumeGas(gas) {
		return
	}

	val := s.host.GetStorage(ctx, s.msg.Address, loc)
	loc.SetBytes(val.Bytes())
}

func opSStore(ctx context.Context, s *state) {
	if s.inStaticCall() {
		s.exit(errWriteProtection)
		return
	}

	if s.config.Istanbul && s.gas <= 2300 { //nolint:gomnd
		s.exit(errOutOfGas)
		return
	}

	key := s.pop()
	val := s.pop()

	legacyGasMetering := !s.config.Istanbul && (s.config.Petersburg || !s.config.Constantinople)

	status := s.host.SetStorage(ctx, s.msg.Address, key, val, s.config)
	cost := uint64(0)

	if s.instrumented {
		s.storeDiff = &instrumentation.StoreDiff{Location: key.Uint64(), Value: uint(val.Uint64())}
	}

	switch status {
	case runtime.StorageUnchanged:
		if s.config.Istanbul {
			// eip-2200
			cost = 800
		} else if legacyGasMetering {
			cost = 5000
		} else {
			cost = 200
		}

	case runtime.StorageModified:
		cost = 5000

	case runtime.StorageModifiedAgain:
		if s.config.Istanbul {
			// eip-2200
			cost = 800
		} else if legacyGasMetering {
			cost = 5000
		} else {
			cost = 200
		}

	case runtime.StorageAdded:
		cost = 20000

	case runtime.StorageDeleted:
		cost = 5000
	}
	if !s.consumeGas(cost) {
		return
	}
}

const sha3WordGas uint64 = 6

func opSha3(ctx context.Context, s *state) {
	offset := s.pop()
	length := s.pop()

	var ok bool
	if s.tmp, ok = s.get2(s.tmp[:0], offset, length); !ok {
		return
	}

	size := length.Uint64()
	if !s.consumeGas(((size + 31) / 32) * sha3WordGas) { //nolint:gomnd
		return
	}

	s.tmp = helper.Keccak256To(s.tmp[:0], s.tmp)

	v := s.push1()
	v.SetBytes(s.tmp)
}

func opPop(ctx context.Context, s *state) {
	s.pop()
}

// context operations

func opAddress(ctx context.Context, s *state) {
	s.push1().SetBytes(s.msg.Address.Bytes())
}

func opBalance(ctx context.Context, s *state) {
	addr, _ := s.popAddr()

	var gas uint64
	if s.config.Istanbul {
		// eip-1884
		gas = 700
	} else if s.config.EIP150 {
		gas = 400
	} else {
		gas = 20
	}

	if !s.consumeGas(gas) {
		return
	}

	s.push1().Set(s.host.GetBalance(ctx, addr))
}

func opSelfBalance(ctx context.Context, s *state) {
	if !s.config.Istanbul {
		s.exit(errOpCodeNotFound)
		return
	}

	s.push1().Set(s.host.GetBalance(ctx, s.msg.Address))
}

func opChainID(ctx context.Context, s *state) {
	if !s.config.Istanbul {
		s.exit(errOpCodeNotFound)
		return
	}

	s.push1().SetUint64(uint64(s.host.GetTxContext().ChainID))
}

func opOrigin(ctx context.Context, s *state) {
	s.push1().SetBytes(s.host.GetTxContext().Origin.Bytes())
}

func opCaller(ctx context.Context, s *state) {
	s.push1().SetBytes(s.msg.Caller.Bytes())
}

func opCallValue(ctx context.Context, s *state) {
	v := s.push1()
	if value := s.msg.Value; value != nil {
		v.Set(value)
	} else {
		v.Set(zero)
	}
}

func min(i, j uint64) uint64 {
	if i < j {
		return i
	}
	return j
}

func opCallDataLoad(ctx context.Context, s *state) {
	offset := s.top()

	bufPtr := bufPool.Get().(*[]byte)
	buf := *bufPtr
	s.setBytes(buf[:32], s.msg.Input, 32, offset)
	offset.SetBytes(buf[:32])
	bufPool.Put(bufPtr)
}

func opCallDataSize(ctx context.Context, s *state) {
	s.push1().SetUint64(uint64(len(s.msg.Input)))
}

func opCodeSize(ctx context.Context, s *state) {
	s.push1().SetUint64(uint64(len(s.code)))
}

func opExtCodeSize(ctx context.Context, s *state) {
	addr, _ := s.popAddr()

	var gas uint64
	if s.config.EIP150 {
		gas = 700
	} else {
		gas = 20
	}
	if !s.consumeGas(gas) {
		return
	}

	s.push1().SetUint64(uint64(s.host.GetCodeSize(ctx, addr)))
}

func opGasPrice(ctx context.Context, s *state) {
	s.push1().SetBytes(s.host.GetTxContext().GasPrice.Bytes())
}

func opReturnDataSize(ctx context.Context, s *state) {
	if !s.config.Byzantium {
		s.exit(errOpCodeNotFound)
	} else {
		s.push1().SetUint64(uint64(len(s.returnData)))
	}
}

func opExtCodeHash(ctx context.Context, s *state) {
	if !s.config.Constantinople {
		s.exit(errOpCodeNotFound)
		return
	}

	address, _ := s.popAddr()

	var gas uint64
	if s.config.Istanbul {
		gas = 700
	} else {
		gas = 400
	}
	if !s.consumeGas(gas) {
		return
	}

	v := s.push1()
	if s.host.Empty(ctx, address) {
		v.Set(zero)
	} else {
		v.SetBytes(s.host.GetCodeHash(ctx, address).Bytes())
	}
}

func opPC(ctx context.Context, s *state) {
	s.push1().SetUint64(uint64(s.ip))
}

func opMSize(ctx context.Context, s *state) {
	s.push1().SetUint64(uint64(len(s.memory)))
}

func opGas(ctx context.Context, s *state) {
	s.push1().SetUint64(s.gas)
}

func (s *state) setBytes(dst, input []byte, size uint64, dataOffset *big.Int) {
	if !dataOffset.IsUint64() {
		// overflow, copy 'size' 0 bytes to dst
		for i := uint64(0); i < size; i++ {
			dst[i] = 0
		}
		return
	}

	inputSize := uint64(len(input))
	begin := min(dataOffset.Uint64(), inputSize)

	copySize := min(size, inputSize-begin)
	if copySize > 0 {
		copy(dst, input[begin:begin+copySize])
	}
	if size-copySize > 0 {
		dst = dst[copySize:]
		for i := uint64(0); i < size-copySize; i++ {
			dst[i] = 0
		}
	}
}

const copyGas uint64 = 3

func opExtCodeCopy(ctx context.Context, s *state) {
	address, _ := s.popAddr()
	memOffset := s.pop()
	codeOffset := s.pop()
	length := s.pop()

	if !s.checkMemory(memOffset, length) {
		return
	}

	size := length.Uint64()
	if !s.consumeGas(((size + 31) / 32) * copyGas) { //nolint:gomnd
		return
	}

	var gas uint64
	if s.config.EIP150 {
		gas = 700
	} else {
		gas = 20
	}
	if !s.consumeGas(gas) {
		return
	}
	code := s.host.GetCode(ctx, address)
	if size != 0 {
		s.setBytes(s.memory[memOffset.Uint64():], code, size, codeOffset)
	}
}

func opCallDataCopy(ctx context.Context, s *state) {
	memOffset := s.pop()
	dataOffset := s.pop()
	length := s.pop()

	if !s.checkMemory(memOffset, length) {
		return
	}

	size := length.Uint64()
	if !s.consumeGas(((size + 31) / 32) * copyGas) { //nolint:gomnd
		return
	}

	if size != 0 {
		s.setBytes(s.memory[memOffset.Uint64():], s.msg.Input, size, dataOffset)
	}
}

func opReturnDataCopy(ctx context.Context, s *state) {
	if !s.config.Byzantium {
		s.exit(errOpCodeNotFound)
		return
	}

	memOffset := s.pop()
	dataOffset := s.pop()
	length := s.pop()

	if !s.checkMemory(memOffset, length) {
		return
	}

	size := length.Uint64()
	if !s.consumeGas(((size + 31) / 32) * copyGas) { //nolint:gomnd
		return
	}

	end := length.Add(dataOffset, length)
	if !end.IsUint64() {
		s.exit(errReturnDataOutOfBounds)
		return
	}
	size = end.Uint64()
	if uint64(len(s.returnData)) < size {
		s.exit(errReturnDataOutOfBounds)
		return
	}

	data := s.returnData[dataOffset.Uint64():size]
	copy(s.memory[memOffset.Uint64():], data)
}

func opCodeCopy(ctx context.Context, s *state) {
	memOffset := s.pop()
	dataOffset := s.pop()
	length := s.pop()

	if !s.checkMemory(memOffset, length) {
		return
	}

	size := length.Uint64()
	if !s.consumeGas(((size + 31) / 32) * copyGas) { //nolint:gomnd
		return
	}
	if size != 0 {
		s.setBytes(s.memory[memOffset.Uint64():], s.code, size, dataOffset)
	}
}

// block information

func opBlockHash(ctx context.Context, s *state) {
	num := s.top()

	if !num.IsInt64() {
		num.Set(zero)
		return
	}

	n := num.Int64()
	lastBlock := s.host.GetTxContext().Number

	if lastBlock-257 < n && n < lastBlock {
		num.SetBytes(s.host.GetBlockHash(n).Bytes())
	} else {
		num.Set(zero)
	}
}

func opCoinbase(ctx context.Context, s *state) {
	s.push1().SetBytes(s.host.GetTxContext().Coinbase.Bytes())
}

func opTimestamp(ctx context.Context, s *state) {
	s.push1().SetInt64(s.host.GetTxContext().Timestamp)
}

func opNumber(ctx context.Context, s *state) {
	s.push1().SetInt64(s.host.GetTxContext().Number)
}

func opDifficulty(ctx context.Context, s *state) {
	s.push1().SetBytes(s.host.GetTxContext().Difficulty.Bytes())
}

func opGasLimit(ctx context.Context, s *state) {
	s.push1().SetInt64(s.host.GetTxContext().GasLimit)
}

func opSelfDestruct(ctx context.Context, s *state) {
	if s.inStaticCall() {
		s.exit(errWriteProtection)
		return
	}

	address, _ := s.popAddr()

	// try to remove the gas first
	var gas uint64
	// EIP150 reprice fork
	if s.config.EIP150 {
		gas = 5000
		if s.config.EIP158 {
			// if empty and transfers value
			if s.host.Empty(ctx, address) && s.host.GetBalance(ctx, s.msg.Address).Sign() != 0 {
				gas += 25000
			}
		} else if !s.host.AccountExists(ctx, address) {
			gas += 25000
		}
	}

	if !s.consumeGas(gas) {
		return
	}

	s.host.Selfdestruct(ctx, s.msg.Address, address)
	s.halt()
}

func opJump(ctx context.Context, s *state) {
	dest := s.pop()

	if s.validJumpdest(dest) {
		s.ip = int(dest.Uint64() - 1)
	} else {
		s.exit(errInvalidJump)
	}
}

func opJumpi(ctx context.Context, s *state) {
	dest := s.pop()
	cond := s.pop()

	if cond.Sign() != 0 {
		if s.validJumpdest(dest) {
			s.ip = int(dest.Uint64() - 1)
		} else {
			s.exit(errInvalidJump)
		}
	}
}

func opJumpDest(ctx context.Context, s *state) {
}

func opPush(n int) instruction {
	return func(ctx context.Context, s *state) {
		ins := s.code
		ip := s.ip

		v := s.push1()
		if ip+1+n > len(ins) {
			v.SetBytes(append(ins[ip+1:], make([]byte, n)...))
		} else {
			v.SetBytes(ins[ip+1 : ip+1+n])
		}

		s.ip += n
	}
}

func opDup(n int) instruction {
	return func(ctx context.Context, s *state) {
		if !s.stackAtLeast(n) {
			s.exit(errStackUnderflow)
		} else {
			val := s.peekAt(n)
			s.push1().Set(val)
		}
	}
}

func opSwap(n int) instruction {
	return func(ctx context.Context, s *state) {
		if !s.stackAtLeast(n + 1) {
			s.exit(errStackUnderflow)
		} else {
			s.swap(n)
		}
	}
}

func opLog(size int) instruction {
	size = size - 1
	return func(ctx context.Context, s *state) {
		if s.inStaticCall() {
			s.exit(errWriteProtection)
			return
		}

		if !s.stackAtLeast(2 + size) { //nolint:gomnd
			s.exit(errStackUnderflow)
			return
		}

		mStart := s.pop()
		mSize := s.pop()

		topics := make([]common.Hash, size)
		for i := 0; i < size; i++ {
			topics[i] = bigToHash(s.pop())
		}

		var ok bool
		s.tmp, ok = s.get2(s.tmp[:0], mStart, mSize)
		if !ok {
			return
		}

		s.host.EmitLog(s.msg.Address, topics, s.tmp)

		if !s.consumeGas(uint64(size) * 375) { //nolint:gomnd
			return
		}
		if !s.consumeGas(mSize.Uint64() * 8) { //nolint:gomnd
			return
		}
	}
}

func opStop(ctx context.Context, s *state) {
	s.halt()
}

func opCreate(op OpCode) instruction {
	return func(ctx context.Context, s *state) {
		if s.inStaticCall() {
			s.exit(errWriteProtection)
			return
		}

		if op == CREATE2 {
			if !s.config.Constantinople {
				s.exit(errOpCodeNotFound)
				return
			}
		}

		// reset the return data
		s.resetReturnData()

		contract, err := s.buildCreateContract(ctx, op)
		if err != nil {
			s.push1().Set(zero)
			if contract != nil {
				s.gas += contract.Gas
			}
			return
		}
		if contract == nil {
			return
		}

		contract.Type = runtime.Create

		// Correct call
		result := s.host.Callx(ctx, contract, s.host)

		v := s.push1()
		if op == CREATE && s.config.Homestead && errors.Is(result.Err, runtime.ErrCodeStoreOutOfGas) {
			v.Set(zero)
		} else if result.Failed() && result.Err != runtime.ErrCodeStoreOutOfGas {
			v.Set(zero)
		} else {
			v.SetBytes(contract.Address.Bytes())
		}

		s.gas += result.GasLeft

		if result.Reverted() {
			s.returnData = append(s.returnData[:0], result.ReturnValue...)
		}

		s.returnVMTrace = result.VMTrace
	}
}

func opCall(op OpCode) instruction {
	return func(ctx context.Context, s *state) {
		s.resetReturnData()

		if op == CALL && s.inStaticCall() {
			if val := s.peekAt(3); val != nil && val.BitLen() > 0 { //nolint:gomnd
				s.exit(errWriteProtection)
				return
			}
		}

		if op == DELEGATECALL && !s.config.Homestead {
			s.exit(errOpCodeNotFound)
			return
		}
		if op == STATICCALL && !s.config.Byzantium {
			s.exit(errOpCodeNotFound)
			return
		}

		var callType runtime.CallType
		switch op {
		case CALL:
			callType = runtime.Call

		case CALLCODE:
			callType = runtime.CallCode

		case DELEGATECALL:
			callType = runtime.DelegateCall

		case STATICCALL:
			callType = runtime.StaticCall

		default:
			panic("not expected")
		}

		contract, offset, size, err := s.buildCallContract(ctx, op)
		if err != nil {
			s.push1().Set(zero)
			if contract != nil {
				s.gas += contract.Gas
			}
			return
		}
		if contract == nil {
			return
		}

		contract.Type = callType

		result := s.host.Callx(ctx, contract, s.host)

		v := s.push1()
		if result.Succeeded() {
			v.Set(one)
		} else {
			v.Set(zero)
		}

		if result.Succeeded() || result.Reverted() {
			if len(result.ReturnValue) != 0 {
				copy(s.memory[offset:offset+size], result.ReturnValue)
				if s.instrumented {
					diff := instrumentation.MemoryDiff{
						Offset: offset,
						Data:   new(big.Int).SetBytes(result.ReturnValue).Uint64(),
					}
					s.memDiff = &diff
				}
			}
		}

		s.gas += result.GasLeft
		s.returnData = append(s.returnData[:0], result.ReturnValue...)
		s.returnVMTrace = result.VMTrace
	}
}

func (s *state) buildCallContract(ctx context.Context, op OpCode) (*runtime.Contract, uint64, uint64, error) {
	// Pop input arguments
	initialGas := s.pop()
	addr, _ := s.popAddr()

	var value *big.Int
	if op == CALL || op == CALLCODE {
		value = s.pop()
	}

	// input range
	inOffset := s.pop()
	inSize := s.pop()

	// output range
	retOffset := s.pop()
	retSize := s.pop()

	// Get the input arguments
	args, ok := s.get2(nil, inOffset, inSize)
	if !ok {
		return nil, 0, 0, nil
	}
	// Check if the memory return offsets are out of bounds
	if !s.checkMemory(retOffset, retSize) {
		return nil, 0, 0, nil
	}

	var gasCost uint64
	if s.config.EIP150 {
		gasCost = 700
	} else {
		gasCost = 40
	}

	eip158 := s.config.EIP158
	transfersValue := (op == CALL || op == CALLCODE) && value != nil && value.Sign() != 0

	if op == CALL {
		if eip158 {
			if transfersValue && s.host.Empty(ctx, addr) {
				gasCost += 25000
			}
		} else if !s.host.AccountExists(ctx, addr) {
			gasCost += 25000
		}
	}
	if transfersValue {
		gasCost += 9000
	}

	var gas uint64

	ok = initialGas.IsUint64()
	if s.config.EIP150 {
		availableGas := s.gas - gasCost
		availableGas = availableGas - availableGas/64 //nolint:gomnd

		if !ok || availableGas < initialGas.Uint64() {
			gas = availableGas
		} else {
			gas = initialGas.Uint64()
		}
	} else {
		if !ok {
			s.exit(errOutOfGas)
			return nil, 0, 0, nil
		}
		gas = initialGas.Uint64()
	}

	gasCost = gasCost + gas

	// Consume gas cost
	if !s.consumeGas(gasCost) {
		return nil, 0, 0, nil
	}
	if transfersValue {
		gas += 2300
	}

	parent := s

	contract := runtime.NewContractCall(s.msg.Depth+1, parent.msg.Origin, parent.msg.Address, addr, value, gas, s.host.GetCode(ctx, addr), args)

	if op == STATICCALL || parent.msg.Static {
		contract.Static = true
	}
	if op == CALLCODE || op == DELEGATECALL {
		contract.Address = parent.msg.Address
		if op == DELEGATECALL {
			contract.Value = parent.msg.Value
			contract.Caller = parent.msg.Caller
		}
	}

	if transfersValue {
		if s.host.GetBalance(ctx, s.msg.Address).Cmp(value) < 0 {
			return contract, 0, 0, fmt.Errorf("bad")
		}
	}
	return contract, retOffset.Uint64(), retSize.Uint64(), nil
}

func (s *state) buildCreateContract(ctx context.Context, op OpCode) (*runtime.Contract, error) {
	// Pop input arguments
	value := s.pop()
	offset := s.pop()
	length := s.pop()

	var salt *big.Int
	if op == CREATE2 {
		salt = s.pop()
	}

	// check if the value can be transferred
	hasTransfer := value != nil && value.Sign() != 0

	// Calculate and consume gas cost

	// var overflow bool
	var gasCost uint64

	// Both CREATE and CREATE2 use memory
	var input []byte
	var ok bool

	input, ok = s.get2(input[:0], offset, length) // Does the memory check
	if !ok {
		return nil, nil
	}

	// Consume memory resize gas (TODO, change with get2)
	if !s.consumeGas(gasCost) {
		return nil, nil
	}

	if hasTransfer {
		if s.host.GetBalance(ctx, s.msg.Address).Cmp(value) < 0 {
			return nil, fmt.Errorf("bad")
		}
	}

	if op == CREATE2 {
		// Consume sha3 gas cost
		size := length.Uint64()
		if !s.consumeGas(((size + 31) / 32) * sha3WordGas) { //nolint:gomnd
			return nil, nil
		}
	}

	// Calculate and consume gas for the call
	gas := s.gas

	// CREATE2 uses by default EIP150
	if s.config.EIP150 || op == CREATE2 { //nolint:gomnd
		gas -= gas / 64 //nolint:gomnd
	}

	if !s.consumeGas(gas) {
		return nil, nil
	}

	// Calculate address
	var address common.Address
	if op == CREATE {
		address = helper.CreateAddress(s.msg.Address, s.host.GetNonce(ctx, s.msg.Address))
	} else {
		address = helper.CreateAddress2(s.msg.Address, bigToHash(salt), input)
	}
	contract := runtime.NewContractCreation(s.msg.Depth+1, s.msg.Origin, s.msg.Address, address, value, gas, input)
	return contract, nil
}

func opHalt(op OpCode) instruction {
	return func(ctx context.Context, s *state) {
		if op == REVERT && !s.config.Byzantium {
			s.exit(errOpCodeNotFound)
			return
		}

		offset := s.pop()
		size := s.pop()

		var ok bool
		s.ret, ok = s.get2(s.ret[:0], offset, size)
		if !ok {
			return
		}

		if op == REVERT {
			s.exit(errRevert)
		} else {
			s.halt()
		}
	}
}

var (
	tt256   = new(big.Int).Lsh(big.NewInt(1), 256)   // 2 ** 256
	tt256m1 = new(big.Int).Sub(tt256, big.NewInt(1)) // 2 ** 256 - 1
)

func toU256(x *big.Int) *big.Int {
	if x.Sign() < 0 || x.BitLen() > 256 {
		x.And(x, tt256m1)
	}
	return x
}

func to256(x *big.Int) *big.Int {
	if x.BitLen() > 255 { //nolint:gomnd
		x.Sub(x, tt256)
	}
	return x
}
