package syncinterfaces

import "github.com/jackc/pgx/v4"

// SynchronizerFlushIDManager is a interface with the methods to manage the flushID
type SynchronizerFlushIDManager interface {
	PendingFlushID(flushID uint64, proverID string)
	CheckFlushID(dbTx pgx.Tx) error
}

type SynchronizerIsTrustedSequencer interface {
	IsTrustedSequencer() bool
}

type SynchronizerCleanTrustedState interface {
	CleanTrustedState()
}

type SynchronizerFullInterface interface {
	SynchronizerFlushIDManager
	SynchronizerIsTrustedSequencer
	SynchronizerCleanTrustedState
}
