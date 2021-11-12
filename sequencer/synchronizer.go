package sequencer

type SynchronizerClient struct {
	IsSynced chan int
}

func NewSynchronizerClient() SynchronizerClient {
	// connect to synch to get signal, that it's synced
	isSynced := make(chan int)
	return SynchronizerClient{IsSynced: isSynced}
}

func (sc *SynchronizerClient) subscrToSyncedSignal() {
	// get signal
	sc.IsSynced <- 1
}
