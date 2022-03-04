package benchmarks

import "context"

// Consumer interfaces required by the package.

// localState gathers the methods required to interact with the state.
type localState interface {
	GetLastBatchNumber(ctx context.Context) (uint64, error)
}
