package evm

// Ethereum Virtual Machine OpCode s
// https://ethervm.io/#opcodes

// OpCode is the EVM opcode
type OpCode byte

const (

	// STOP halts execution of the contract
	STOP OpCode = 0x00

	// ADD performs (u)int256 addition modulo 2**256
	ADD OpCode = 0x01

	// MUL performs (u)int256 multiplication modulo 2**256
	MUL OpCode = 0x02

	// SUB performs (u)int256 subtraction modulo 2**256
	SUB OpCode = 0x03

	// DIV performs uint256 division
	DIV OpCode = 0x04

	// SDIV performs int256 division
	SDIV OpCode = 0x05

	// MOD performs uint256 modulus
	MOD OpCode = 0x06

	// SMOD performs int256 modulus
	SMOD OpCode = 0x07

	// ADDMOD performs (u)int256 addition modulo N
	ADDMOD OpCode = 0x08

	// MULMOD performs (u)int256 multiplication modulo N
	MULMOD OpCode = 0x09

	// EXP performs uint256 exponentiation modulo 2**256
	EXP OpCode = 0x0A

	// SIGNEXTEND performs sign extends x from (b + 1) * 8 bits to 256 bits.
	SIGNEXTEND OpCode = 0x0B

	// LT performs int256 comparison
	LT OpCode = 0x10

	// GT performs int256 comparison
	GT OpCode = 0x11

	// SLT performs int256 comparison
	SLT OpCode = 0x12

	// SGT performs int256 comparison
	SGT OpCode = 0x13

	// EQ performs (u)int256 equality
	EQ OpCode = 0x14

	// ISZERO checks if (u)int256 is zero
	ISZERO OpCode = 0x15

	// AND performs 256-bit bitwise and
	AND OpCode = 0x16

	// OR performs 256-bit bitwise or
	OR OpCode = 0x17

	// XOR performs 256-bit bitwise xor
	XOR OpCode = 0x18

	// NOT performs 256-bit bitwise not
	NOT OpCode = 0x19

	// BYTE returns the ith byte of (u)int256 x counting from most significant byte
	BYTE OpCode = 0x1A

	// SHL performs a shift left
	SHL OpCode = 0x1B

	// SHR performs a logical shift right
	SHR OpCode = 0x1C

	// SAR performs an arithmetic shift right
	SAR OpCode = 0x1D

	// SHA3 performs the keccak256 hash function
	SHA3 OpCode = 0x20

	// ADDRESS returns the address of the executing contract
	ADDRESS OpCode = 0x30

	// BALANCE returns the address balance in wei
	BALANCE OpCode = 0x31

	// ORIGIN returns the transaction origin address
	ORIGIN OpCode = 0x32

	// CALLER returns the message caller address
	CALLER OpCode = 0x33

	// CALLVALUE returns the message funds in wei
	CALLVALUE OpCode = 0x34

	// CALLDATALOAD reads a (u)int256 from message data
	CALLDATALOAD OpCode = 0x35

	// CALLDATASIZE returns the message data length in bytes
	CALLDATASIZE OpCode = 0x36

	// CALLDATACOPY copies the message data
	CALLDATACOPY OpCode = 0x37

	// CODESIZE returns the length of the executing contract's code in bytes
	CODESIZE OpCode = 0x38

	// CODECOPY copies the executing contract bytecode
	CODECOPY OpCode = 0x39

	// GASPRICE returns the gas price of the executing transaction, in wei per unit of gas
	GASPRICE OpCode = 0x3A

	// EXTCODESIZE returns the length of the contract bytecode at addr
	EXTCODESIZE OpCode = 0x3B

	// EXTCODECOPY copies the contract bytecode
	EXTCODECOPY OpCode = 0x3C

	// RETURNDATASIZE returns the size of the returned data from the last external call in bytes
	RETURNDATASIZE OpCode = 0x3D

	// RETURNDATACOPY copies the returned data
	RETURNDATACOPY OpCode = 0x3E

	// EXTCODEHASH returns the hash of the specified contract bytecode
	EXTCODEHASH OpCode = 0x3F

	// BLOCKHASH returns the hash of the specific block. Only valid for the last 256 most recent blocks
	BLOCKHASH OpCode = 0x40

	// COINBASE returns the address of the current block's miner
	COINBASE OpCode = 0x41

	// TIMESTAMP returns the current block's Unix timestamp in seconds
	TIMESTAMP OpCode = 0x42

	// NUMBER returns the current block's number
	NUMBER OpCode = 0x43

	// DIFFICULTY returns the current block's difficulty
	DIFFICULTY OpCode = 0x44

	// GASLIMIT returns the current block's gas limit
	GASLIMIT OpCode = 0x45

	// CHAINID returns the id of the chain
	CHAINID OpCode = 0x46

	// SELFBALANCE returns the balance of the current account
	SELFBALANCE OpCode = 0x47

	// POP pops a (u)int256 off the stack and discards it
	POP OpCode = 0x50

	// MLOAD reads a (u)int256 from memory
	MLOAD OpCode = 0x51

	// MSTORE writes a (u)int256 to memory
	MSTORE OpCode = 0x52

	// MSTORE8 writes a uint8 to memory
	MSTORE8 OpCode = 0x53

	// SLOAD reads a (u)int256 from storage
	SLOAD OpCode = 0x54

	// SSTORE writes a (u)int256 to storage
	SSTORE OpCode = 0x55

	// JUMP performs an unconditional jump
	JUMP OpCode = 0x56

	// JUMPI performs a conditional jump if condition is truthy
	JUMPI OpCode = 0x57

	// PC returns the program counter
	PC OpCode = 0x58

	// MSIZE returns the size of memory for this contract execution, in bytes
	MSIZE OpCode = 0x59

	// GAS returns the remaining gas
	GAS OpCode = 0x5A

	// JUMPDEST corresponds to a possible jump destination
	JUMPDEST OpCode = 0x5B

	// PUSH1 pushes a 1-byte value onto the stack
	PUSH1 OpCode = 0x60

	// PUSH2 pushes a 2-bytes value onto the stack
	PUSH2 OpCode = 0x61

	// PUSH3 pushes a 3-bytes value onto the stack
	PUSH3 OpCode = 0x62

	// PUSH4 pushes a 4-bytes value onto the stack
	PUSH4 OpCode = 0x63

	// PUSH5 pushes a 5-bytes value onto the stack
	PUSH5 OpCode = 0x64

	// PUSH6 pushes a 6-bytes value onto the stack
	PUSH6 OpCode = 0x65

	// PUSH7 pushes a 7-bytes value onto the stack
	PUSH7 OpCode = 0x66

	// PUSH8 pushes a 8-bytes value onto the stack
	PUSH8 OpCode = 0x67

	// PUSH9 pushes a 9-bytes value onto the stack
	PUSH9 OpCode = 0x68

	// PUSH10 pushes a 10-bytes value onto the stack
	PUSH10 OpCode = 0x69

	// PUSH11 pushes a 11-bytes value onto the stack
	PUSH11 OpCode = 0x6A

	// PUSH12 pushes a 12-bytes value onto the stack
	PUSH12 OpCode = 0x6B

	// PUSH13 pushes a 13-bytes value onto the stack
	PUSH13 OpCode = 0x6C

	// PUSH14 pushes a 14-bytes value onto the stack
	PUSH14 OpCode = 0x6D

	// PUSH15 pushes a 15-bytes value onto the stack
	PUSH15 OpCode = 0x6E

	// PUSH16 pushes a 16-bytes value onto the stack
	PUSH16 OpCode = 0x6F

	// PUSH17 pushes a 17-bytes value onto the stack
	PUSH17 OpCode = 0x70

	// PUSH18 pushes a 18-bytes value onto the stack
	PUSH18 OpCode = 0x71

	// PUSH19 pushes a 19-bytes value onto the stack
	PUSH19 OpCode = 0x72

	// PUSH20 pushes a 20-bytes value onto the stack
	PUSH20 OpCode = 0x73

	// PUSH21 pushes a 21-bytes value onto the stack
	PUSH21 OpCode = 0x74

	// PUSH22 pushes a 22-bytes value onto the stack
	PUSH22 OpCode = 0x75

	// PUSH23 pushes a 23-bytes value onto the stack
	PUSH23 OpCode = 0x76

	// PUSH24 pushes a 24-bytes value onto the stack
	PUSH24 OpCode = 0x77

	// PUSH25 pushes a 25-bytes value onto the stack
	PUSH25 OpCode = 0x78

	// PUSH26 pushes a 26-bytes value onto the stack
	PUSH26 OpCode = 0x79

	// PUSH27 pushes a 27-bytes value onto the stack
	PUSH27 OpCode = 0x7A

	// PUSH28 pushes a 28-bytes value onto the stack
	PUSH28 OpCode = 0x7B

	// PUSH29 pushes a 29-bytes value onto the stack
	PUSH29 OpCode = 0x7C

	// PUSH30 pushes a 30-bytes value onto the stack
	PUSH30 OpCode = 0x7D

	// PUSH31 pushes a 31-bytes value onto the stack
	PUSH31 OpCode = 0x7E

	// PUSH32 pushes a 32-byte value onto the stack
	PUSH32 OpCode = 0x7F

	// DUP1 clones the last value on the stack
	DUP1 OpCode = 0x80

	// DUP2 clones the 2nd last value on the stack
	DUP2 OpCode = 0x81

	// DUP3 clones the 3rd last value on the stack
	DUP3 OpCode = 0x82

	// DUP4 clones the 4th last value on the stack
	DUP4 OpCode = 0x83

	// DUP5 clones the 5th last value on the stack
	DUP5 OpCode = 0x84

	// DUP6 clones the 6th last value on the stack
	DUP6 OpCode = 0x85

	// DUP7 clones the 7th last value on the stack
	DUP7 OpCode = 0x86

	// DUP8 clones the 8th last value on the stack
	DUP8 OpCode = 0x87

	// DUP9 clones the 9th last value on the stack
	DUP9 OpCode = 0x88

	// DUP10 clones the 10th last value on the stack
	DUP10 OpCode = 0x89

	// DUP11 clones the 11th last value on the stack
	DUP11 OpCode = 0x8A

	// DUP12 clones the 12th last value on the stack
	DUP12 OpCode = 0x8B

	// DUP13 clones the 13th last value on the stack
	DUP13 OpCode = 0x8C

	// DUP14 clones the 14th last value on the stack
	DUP14 OpCode = 0x8D

	// DUP15 clones the 15th last value on the stack
	DUP15 OpCode = 0x8E

	// DUP16 clones the 16th last value on the stack
	DUP16 OpCode = 0x8F

	// SWAP1 swaps the last two values on the stack
	SWAP1 OpCode = 0x90

	// SWAP2 swaps the top of the stack with the 3rd last element
	SWAP2 OpCode = 0x91

	// SWAP3 swaps the top of the stack with the 4th last element
	SWAP3 OpCode = 0x92

	// SWAP4 swaps the top of the stack with the 5th last element
	SWAP4 OpCode = 0x93

	// SWAP5 swaps the top of the stack with the 6th last element
	SWAP5 OpCode = 0x94

	// SWAP6 swaps the top of the stack with the 7th last element
	SWAP6 OpCode = 0x95

	// SWAP7 swaps the top of the stack with the 8th last element
	SWAP7 OpCode = 0x96

	// SWAP8 swaps the top of the stack with the 9th last element
	SWAP8 OpCode = 0x97

	// SWAP9 swaps the top of the stack with the 10th last element
	SWAP9 OpCode = 0x98

	// SWAP10 swaps the top of the stack with the 11th last element
	SWAP10 OpCode = 0x99

	// SWAP11 swaps the top of the stack with the 12th last element
	SWAP11 OpCode = 0x9A

	// SWAP12 swaps the top of the stack with the 13th last element
	SWAP12 OpCode = 0x9B

	// SWAP13 swaps the top of the stack with the 14th last element
	SWAP13 OpCode = 0x9C

	// SWAP14 swaps the top of the stack with the 15th last element
	SWAP14 OpCode = 0x9D

	// SWAP15 swaps the top of the stack with the 16th last element
	SWAP15 OpCode = 0x9E

	// SWAP16 swaps the top of the stack with the 17th last element
	SWAP16 OpCode = 0x9F

	// LOG0 fires an event without topics
	LOG0 OpCode = 0xA0

	// LOG1 fires an event with one topic
	LOG1 OpCode = 0xA1

	// LOG2 fires an event with two topics
	LOG2 OpCode = 0xA2

	// LOG3 fires an event with three topics
	LOG3 OpCode = 0xA3

	// LOG4 fires an event with four topics
	LOG4 OpCode = 0xA4

	// CREATE creates a child contract
	CREATE OpCode = 0xF0

	// CALL calls a method in another contract
	CALL OpCode = 0xF1

	// CALLCODE calls a method in another contract
	CALLCODE OpCode = 0xF2

	// RETURN returns from this contract call
	RETURN OpCode = 0xF3

	// DELEGATECALL calls a method in another contract using the storage of the current contract
	DELEGATECALL OpCode = 0xF4

	// CREATE2 creates a child contract with a salt
	CREATE2 OpCode = 0xF5

	// STATICCALL calls a method in another contract
	STATICCALL OpCode = 0xFA

	// REVERT reverts with return data
	REVERT OpCode = 0xFD

	// SELFDESTRUCT destroys the contract and sends all funds to addr
	SELFDESTRUCT OpCode = 0xFF
)
