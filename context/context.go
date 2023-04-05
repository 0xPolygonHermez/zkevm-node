package context

import (
	"context"
	"net/http"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// DeadlineExceeded context.DeadlineExceeded
var DeadlineExceeded error = context.DeadlineExceeded

// CancelFunc context.CancelFunc
type CancelFunc context.CancelFunc

// contextKey is required by the context.WithValue func to
// avoid conflict among the context keys
type contextKey string

const (
	idKey          contextKey = "contextID"
	loggerKey      contextKey = "logger"
	wsConnKey      contextKey = "wsConn"
	httpRequestKey contextKey = "httpRequest"
)

// RequestContext is a context used by the jRPC requests
type RequestContext struct {
	context.Context
}

// Background creates and initializes an instance of RequestContext
// based on the context.Background
func Background() *RequestContext {
	id := uuid.NewString()
	return NewRequestContext(context.Background(), id)
}

// Wrap wraps a context into an instance of RequestContext
func Wrap(ctx context.Context) *RequestContext {
	id := uuid.NewString()
	return NewRequestContext(ctx, id)
}

// NewRequestContext creates and initializes an instance of RequestContext
func NewRequestContext(ctx context.Context, id string) *RequestContext {
	c := context.WithValue(ctx, idKey, id)
	logger := log.WithFields(string(idKey), id)
	c = context.WithValue(c, loggerKey, logger)
	return &RequestContext{Context: c}
}

// InnerContext return the inner context
func (r *RequestContext) InnerContext() context.Context {
	return r
}

// ID returns the context ID
func (r *RequestContext) ID() string {
	return r.Value(idKey).(string)
}

// Logger returns a contextualized instance of Logger
func (r *RequestContext) Logger() *log.Logger {
	return r.Value(idKey).(*log.Logger)
}

// SetWsConn sets the websocket connection attached to this request
func (r *RequestContext) SetWsConn(wsConn *websocket.Conn) {
	r.Context = context.WithValue(r.Context, wsConnKey, wsConn)
}

// WsConn returns the websocket connection attached to this request
func (r *RequestContext) WsConn() *websocket.Conn {
	return r.Value(wsConnKey).(*websocket.Conn)
}

// SetHttpRequest sets the http request to this request
func (r *RequestContext) SetHttpRequest(httpRequest *http.Request) {
	r.Context = context.WithValue(r.Context, httpRequestKey, httpRequest)
}

// HttpRequest returns the http request to this request
func (r *RequestContext) HttpRequest() *http.Request {
	return r.Value(httpRequestKey).(*http.Request)
}

// SetTimeout sets the context timeout
func (r *RequestContext) SetTimeout(timeout time.Duration) CancelFunc {
	c, cancel := context.WithTimeout(r, timeout)
	r.Context = c
	return CancelFunc(cancel)
}

// WithCancel applies the context.WithCancel to the context
func (r *RequestContext) WithCancel() CancelFunc {
	c, cancel := context.WithCancel(r)
	r.Context = c
	return CancelFunc(cancel)
}

// SetValue sets a value to the context
func (r *RequestContext) SetValue(key, value any) {
	r.Context = context.WithValue(r.Context, key, value)
}

// Clone creates a copy of the current instance
func (r *RequestContext) Clone() *RequestContext {
	return &RequestContext{
		Context: r.Context,
	}
}
