package jsonrpc

// storageInterface json rpc internal storage to persist data
type storageInterface interface {
	Lock()
	Unlock()
	GetAllBlockFiltersWithWSConn() map[string]*Filter
	GetAllLogFiltersWithWSConn() map[string]*Filter
	GetFilter(filterID string) (*Filter, error)
	NewBlockFilter(wsConn *concurrentWsConn) (string, error)
	NewLogFilter(wsConn *concurrentWsConn, filter LogFilter) (string, error)
	NewPendingTransactionFilter(wsConn *concurrentWsConn) (string, error)
	UninstallFilter(filterID string) error
	UninstallFilterByWSConn(wsConn *concurrentWsConn) error
	UpdateFilterLastPoll(filterID string) error
}
