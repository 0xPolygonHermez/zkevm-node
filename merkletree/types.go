package merkletree

type ResultCode int64

const (
	Unspecified ResultCode = iota
	Success
	KeyNotFound
	DBError
	InternalError
)

// Proof is a proof generated on Get operation.
type Proof struct {
	Root  []uint64
	Key   []uint64
	Value []uint64
}

// UpdateProof is a proof generated on Set operation.
type UpdateProof struct {
	OldRoot  []uint64
	NewRoot  []uint64
	Key      []uint64
	NewValue []uint64
}

// ProgramProof is a proof generated on GetProgram operation.
type ProgramProof struct {
	Data []byte
}
