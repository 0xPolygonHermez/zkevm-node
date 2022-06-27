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

type ProgramProof struct {
	Data []byte
}
