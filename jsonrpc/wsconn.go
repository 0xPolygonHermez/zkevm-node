package jsonrpc

import (
	"sync"

	"github.com/gorilla/websocket"
)

// concurrentWsConn is a wrapped web socket connection
// that provide methods to deal with concurrency
type concurrentWsConn struct {
	wsConn *websocket.Conn
	mutex  *sync.Mutex
}

// NewConcurrentWsConn creates a new instance of concurrentWsConn
func newConcurrentWsConn(wsConn *websocket.Conn) *concurrentWsConn {
	return &concurrentWsConn{
		wsConn: wsConn,
		mutex:  &sync.Mutex{},
	}
}

// ReadMessage reads a message from the inner web socket connection
func (c *concurrentWsConn) ReadMessage() (messageType int, p []byte, err error) {
	return c.wsConn.ReadMessage()
}

// WriteMessage writes a message to the inner web socket connection
func (c *concurrentWsConn) WriteMessage(messageType int, data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.wsConn.WriteMessage(messageType, data)
}

// Close closes the inner web socket connection
func (c *concurrentWsConn) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.wsConn.Close()
}

// SetReadLimit sets the read limit to the inner web socket connection
func (c *concurrentWsConn) SetReadLimit(limit int64) {
	c.wsConn.SetReadLimit(limit)
}
