package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/metrics"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/didip/tollbooth/v6"
	"github.com/gorilla/websocket"
)

const (
	// APIEth represents the eth API prefix.
	APIEth = "eth"
	// APINet represents the net API prefix.
	APINet = "net"
	// APIDebug represents the debug API prefix.
	APIDebug = "debug"
	// APIZKEVM represents the zkevm API prefix.
	APIZKEVM = "zkevm"
	// APITxPool represents the txpool API prefix.
	APITxPool = "txpool"
	// APIWeb3 represents the web3 API prefix.
	APIWeb3 = "web3"

	wsBufferSizeLimitInBytes = 1024
	maxRequestContentLength  = 1024 * 1024 * 5
	contentType              = "application/json"
)

// https://www.jsonrpc.org/historical/json-rpc-over-http.html#http-header
var acceptedContentTypes = []string{contentType, "application/json-rpc", "application/jsonrequest"}

// Server is an API backend to handle RPC requests
type Server struct {
	config     Config
	chainID    uint64
	handler    *Handler
	srv        *http.Server
	wsSrv      *http.Server
	wsUpgrader websocket.Upgrader
}

// Service implementation of a service an it's name
type Service struct {
	Name    string
	Service interface{}
}

// NewServer returns the JsonRPC server
func NewServer(
	cfg Config,
	chainID uint64,
	p types.PoolInterface,
	s types.StateInterface,
	storage storageInterface,
	services []Service,
) *Server {
	s.PrepareWebSocket()
	handler := newJSONRpcHandler()

	for _, service := range services {
		handler.registerService(service)
	}

	srv := &Server{
		config:  cfg,
		handler: handler,
		chainID: chainID,
	}
	return srv
}

// Start initializes the JSON RPC server to listen for request
func (s *Server) Start() error {
	metrics.Register()

	if s.config.WebSockets.Enabled {
		go s.startWS()
	}

	return s.startHTTP()
}

// startHTTP starts a server to respond http requests
func (s *Server) startHTTP() error {
	if s.srv != nil {
		return fmt.Errorf("server already started")
	}

	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("failed to create tcp listener: %v", err)
		return err
	}

	mux := http.NewServeMux()

	lmt := tollbooth.NewLimiter(s.config.MaxRequestsPerIPAndSecond, nil)
	mux.Handle("/", tollbooth.LimitFuncHandler(lmt, s.handle))

	s.srv = &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: s.config.ReadTimeout.Duration,
		ReadTimeout:       s.config.ReadTimeout.Duration,
		WriteTimeout:      s.config.WriteTimeout.Duration,
	}
	log.Infof("http server started: %s", address)
	if err := s.srv.Serve(lis); err != nil {
		if err == http.ErrServerClosed {
			log.Infof("http server stopped")
			return nil
		}
		log.Errorf("closed http connection: %v", err)
		return err
	}
	return nil
}

// startWS starts a server to respond WebSockets connections
func (s *Server) startWS() {
	log.Infof("starting websocket server")

	if s.wsSrv != nil {
		log.Errorf("websocket server already started")
		return
	}

	address := fmt.Sprintf("%s:%d", s.config.WebSockets.Host, s.config.WebSockets.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("failed to create tcp listener: %v", err)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleWs)

	s.wsSrv = &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: s.config.ReadTimeout.Duration,
		ReadTimeout:       s.config.ReadTimeout.Duration,
		WriteTimeout:      s.config.WriteTimeout.Duration,
	}
	s.wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  wsBufferSizeLimitInBytes,
		WriteBufferSize: wsBufferSizeLimitInBytes,
	}
	log.Infof("websocket server started: %s", address)
	if err := s.wsSrv.Serve(lis); err != nil {
		if err == http.ErrServerClosed {
			log.Infof("websocket server stopped")
			return
		}
		log.Errorf("closed websocket connection: %v", err)
		return
	}
}

// Stop shutdown the rpc server
func (s *Server) Stop() error {
	if s.srv != nil {
		if err := s.srv.Shutdown(context.Background()); err != nil {
			return err
		}

		if err := s.srv.Close(); err != nil {
			return err
		}
		s.srv = nil
	}

	if s.wsSrv != nil {
		if err := s.wsSrv.Shutdown(context.Background()); err != nil {
			return err
		}

		if err := s.wsSrv.Close(); err != nil {
			return err
		}
		s.wsSrv = nil
	}

	return nil
}

func (s *Server) handle(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		return
	}

	if req.Method == http.MethodGet {
		_, err := w.Write([]byte("zkEVM JSON RPC Server"))
		if err != nil {
			log.Error(err)
		}
		return
	}

	if code, err := validateRequest(req); err != nil {
		handleInvalidRequest(w, err, code)
		return
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	single, err := s.isSingleRequest(data)
	if err != nil {
		handleInvalidRequest(w, err, http.StatusBadRequest)
		return
	}

	start := time.Now()
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	var respLen int
	if single {
		respLen = s.handleSingleRequest(req, w, data)
	} else {
		respLen = s.handleBatchRequest(req, w, data)
	}
	metrics.RequestDuration(start)
	combinedLog(req, start, http.StatusOK, respLen)
}

// validateRequest returns a non-zero response code and error message if the
// request is invalid.
func validateRequest(req *http.Request) (int, error) {
	if req.Method != http.MethodPost {
		err := errors.New("method " + req.Method + " not allowed")
		return http.StatusMethodNotAllowed, err
	}

	if req.ContentLength > maxRequestContentLength {
		err := fmt.Errorf("content length too large (%d>%d)", req.ContentLength, maxRequestContentLength)
		return http.StatusRequestEntityTooLarge, err
	}

	// Check content-type
	if mt, _, err := mime.ParseMediaType(req.Header.Get("content-type")); err == nil {
		for _, accepted := range acceptedContentTypes {
			if accepted == mt {
				return 0, nil
			}
		}
	}
	// Invalid content-type
	err := fmt.Errorf("invalid content type, only %s is supported", contentType)
	return http.StatusUnsupportedMediaType, err
}

func (s *Server) isSingleRequest(data []byte) (bool, error) {
	x := bytes.TrimLeft(data, " \t\r\n")

	if len(x) == 0 {
		return false, fmt.Errorf("empty request body")
	}

	return x[0] != '[', nil
}

func (s *Server) handleSingleRequest(httpRequest *http.Request, w http.ResponseWriter, data []byte) int {
	defer metrics.RequestHandled(metrics.RequestHandledLabelSingle)
	request, err := s.parseRequest(data)
	if err != nil {
		handleInvalidRequest(w, err, http.StatusBadRequest)
		return 0
	}
	req := handleRequest{Request: request, HttpRequest: httpRequest}
	response := s.handler.Handle(req)

	respBytes, err := json.Marshal(response)
	if err != nil {
		handleError(w, err)
		return 0
	}

	_, err = w.Write(respBytes)
	if err != nil {
		handleError(w, err)
		return 0
	}
	return len(respBytes)
}

func (s *Server) handleBatchRequest(httpRequest *http.Request, w http.ResponseWriter, data []byte) int {
	// Checking if batch requests are enabled
	if !s.config.BatchRequestsEnabled {
		handleInvalidRequest(w, types.ErrBatchRequestsDisabled, http.StatusBadRequest)
		return 0
	}

	defer metrics.RequestHandled(metrics.RequestHandledLabelBatch)
	requests, err := s.parseRequests(data)
	if err != nil {
		handleInvalidRequest(w, err, http.StatusBadRequest)
		return 0
	}

	// Checking if batch requests limit is exceeded
	if s.config.BatchRequestsLimit > 0 {
		if len(requests) > int(s.config.BatchRequestsLimit) {
			handleInvalidRequest(w, types.ErrBatchRequestsLimitExceeded, http.StatusRequestEntityTooLarge)
			return 0
		}
	}

	responses := make([]types.Response, 0, len(requests))

	for _, request := range requests {
		req := handleRequest{Request: request, HttpRequest: httpRequest}
		response := s.handler.Handle(req)
		responses = append(responses, response)
	}

	respBytes, _ := json.Marshal(responses)
	_, err = w.Write(respBytes)
	if err != nil {
		log.Error(err)
		return 0
	}
	return len(respBytes)
}

func (s *Server) parseRequest(data []byte) (types.Request, error) {
	var req types.Request

	if err := json.Unmarshal(data, &req); err != nil {
		return types.Request{}, fmt.Errorf("invalid json object request body")
	}

	return req, nil
}

func (s *Server) parseRequests(data []byte) ([]types.Request, error) {
	var requests []types.Request

	if err := json.Unmarshal(data, &requests); err != nil {
		return nil, fmt.Errorf("invalid json array request body")
	}

	return requests, nil
}

func (s *Server) handleWs(w http.ResponseWriter, req *http.Request) {
	// CORS rule - Allow requests from anywhere
	s.wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// Upgrade the connection to a WS one
	wsConn, err := s.wsUpgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to upgrade to a WS connection, %s", err.Error()))

		return
	}

	// Set read limit
	wsConn.SetReadLimit(s.config.WebSockets.ReadLimit)

	// Defer WS closure
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			log.Error(fmt.Sprintf("Unable to gracefully close WS connection, %s", err.Error()))
		}
	}(wsConn)

	log.Info("Websocket connection established")
	var mu sync.Mutex
	for {
		msgType, message, err := wsConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
				log.Info("Closing WS connection gracefully")
			} else if errors.Is(err, websocket.ErrReadLimit) {
				log.Info("Closing WS connection due to read limit exceeded")
			} else {
				log.Error(fmt.Sprintf("Unable to read WS message, %s", err.Error()))
				log.Info("Closing WS connection with error")
			}

			s.handler.RemoveFilterByWsConn(wsConn)

			break
		}

		if msgType == websocket.TextMessage || msgType == websocket.BinaryMessage {
			go func() {
				mu.Lock()
				defer mu.Unlock()
				resp, err := s.handler.HandleWs(message, wsConn, req)
				if err != nil {
					log.Error(fmt.Sprintf("Unable to handle WS request, %s", err.Error()))
					_ = wsConn.WriteMessage(msgType, []byte(fmt.Sprintf("WS Handle error: %s", err.Error())))
				} else {
					_ = wsConn.WriteMessage(msgType, resp)
				}
			}()
		}
	}
}

func handleInvalidRequest(w http.ResponseWriter, err error, code int) {
	defer metrics.RequestHandled(metrics.RequestHandledLabelInvalid)
	log.Infof("Invalid Request: %v", err.Error())
	http.Error(w, err.Error(), code)
}

func handleError(w http.ResponseWriter, err error) {
	defer metrics.RequestHandled(metrics.RequestHandledLabelError)
	log.Errorf("Error processing request: %v", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// RPCErrorResponse formats error to be returned through RPC
func RPCErrorResponse(code int, message string, err error) (interface{}, types.Error) {
	return RPCErrorResponseWithData(code, message, nil, err)
}

// RPCErrorResponseWithData formats error to be returned through RPC
func RPCErrorResponseWithData(code int, message string, data *[]byte, err error) (interface{}, types.Error) {
	if err != nil {
		log.Errorf("%v: %v", message, err.Error())
	} else {
		log.Error(message)
	}
	return nil, types.NewRPCErrorWithData(code, message, data)
}

func combinedLog(r *http.Request, start time.Time, httpStatus, dataLen int) {
	log.Infof("%s - - %s \"%s %s %s\" %d %d \"%s\" \"%s\"",
		r.RemoteAddr,
		start.Format("[02/Jan/2006:15:04:05 -0700]"),
		r.Method,
		r.URL.Path,
		r.Proto,
		httpStatus,
		dataLen,
		r.Host,
		r.UserAgent(),
	)
}
