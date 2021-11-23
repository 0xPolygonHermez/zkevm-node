package evm

import "github.com/holiman/uint256"

// Memory represents the memory
type Memory struct {
	data []byte
}

// NewMemory is the constructor
func NewMemory() *Memory {
	return &Memory{}
}

// Set sets a value at memory offset
func (m *Memory) Set(offset uint64, size uint64, value []byte) {
	if size > 0 {
		if offset+size > uint64(len(m.data)) {
			panic("invalid memory")
		}
		copy(m.data[offset:offset+size], value)
	}
}

// Set32 sets the 32 bytes starting at offset to the value of val, left-padded with zeroes to
// 32 bytes.
func (m *Memory) Set32(offset uint64, val *uint256.Int) {
	// length of store may never be less than offset + size.
	// The store should be resized PRIOR to setting the memory
	if offset+32 > uint64(len(m.data)) {
		panic("invalid memory: store empty")
	}
	// Zero the memory area
	copy(m.data[offset:offset+32], []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	// Fill in relevant bits
	val.WriteToSlice(m.data[offset:])
}

// GetCopy returns size bytes from memory offset as a new slice
func (m *Memory) GetCopy(offset, size int64) (cpy []byte) {
	if size == 0 {
		return nil
	}

	if len(m.data) > int(offset) {
		cpy = make([]byte, size)
		copy(cpy, m.data[offset:offset+size])

		return
	}

	return
}

// GetPtr returns size bytes from memory offset
func (m *Memory) GetPtr(offset, size int64) []byte {
	if size == 0 {
		return nil
	}

	if len(m.data) > int(offset) {
		return m.data[offset : offset+size]
	}

	return nil
}

// Len returns the length of the backing slice
func (m *Memory) Len() int {
	return len(m.data)
}

// Data returns the backing slice
func (m *Memory) Data() []byte {
	return m.data
}
