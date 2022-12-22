package sequencer

// Pool Loader and DB Updater
type dbManager struct {
	txPool txPool
	state  stateInterface
	worker workerInterface
}

func newDBManager(txPool txPool, state stateInterface, worker *Worker) *dbManager {
	return &dbManager{txPool: txPool, state: state, worker: worker}
}

func (d *dbManager) Start() {
	go d.loadFromPool()
}

func (d *dbManager) loadFromPool() {
	// TODO: Endless loop that keeps loading tx from the DB into the worker
}
