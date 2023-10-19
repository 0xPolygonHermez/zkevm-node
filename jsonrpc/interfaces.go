package jsonrpc

import (
	"sync/atomic"

	"github.com/gorilla/websocket"
)

// storageInterface json rpc internal storage to persist data
type storageInterface interface {
	GetAllBlockFiltersWithWSConn() ([]*Filter, error)
	GetAllLogFiltersWithWSConn() ([]*Filter, error)
	GetFilter(filterID string) (*Filter, error)
	NewBlockFilter(wsConn *atomic.Pointer[websocket.Conn]) (string, error)
	NewLogFilter(wsConn *atomic.Pointer[websocket.Conn], filter LogFilter) (string, error)
	NewPendingTransactionFilter(wsConn *atomic.Pointer[websocket.Conn]) (string, error)
	UninstallFilter(filterID string) error
	UninstallFilterByWSConn(wsConn *atomic.Pointer[websocket.Conn]) error
	UpdateFilterLastPoll(filterID string) error
}
