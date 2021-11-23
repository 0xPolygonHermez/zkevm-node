package jsonrpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
)

// Server is an API backend to handle RPC requests
type Server struct {
	config  Config
	handler *Handler
}

// NewServer returns the JsonRPC server
func NewServer(config Config, p pool.Pool, s state.State) *Server {
	ethEndpoints := &Eth{chainID: config.ChainID, pool: p, state: s}
	netEndpoints := &Net{chainID: config.ChainID}

	handler := newJSONRpcHandler(ethEndpoints, netEndpoints)

	srv := &Server{
		config:  config,
		handler: handler,
	}
	return srv
}

// Start initializes the JSON RPC server to listen for request
func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	log.Infof("http server started: %s", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("failed to create tcp listener: %v", err)
		return err
	}

	mux := http.DefaultServeMux
	mux.HandleFunc("/", s.handle)

	srv := http.Server{
		Handler: mux,
	}
	if err := srv.Serve(lis); err != nil {
		log.Errorf("closed http connection: %v", err)
		return err
	}
	return nil
}

func (s *Server) handle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if (*req).Method == "OPTIONS" {
		return
	}

	if req.Method == "GET" {
		_, err := w.Write([]byte("Hermez JSON-RPC"))
		if err != nil {
			log.Error(err)
		}
		return
	}

	if req.Method != "POST" {
		err := errors.New("method " + req.Method + " not allowed")
		handleError(w, err)
		return
	}

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	single, err := s.isSingleRequest(data)
	if err != nil {
		handleError(w, err)
		return
	}

	if single {
		s.handleSingleRequest(w, data)
	} else {
		s.handleBatchRequest(w, data)
	}
}

func (s *Server) isSingleRequest(data []byte) (bool, detailedError) {
	x := bytes.TrimLeft(data, " \t\r\n")

	if len(x) == 0 {
		return false, newInvalidRequestError("Invalid json request")
	}

	return x[0] == '{', nil
}

func (s *Server) handleSingleRequest(w http.ResponseWriter, data []byte) {
	request, err := s.parseRequest(data)
	if err != nil {
		handleError(w, err)
		return
	}

	response := s.handler.Handle(request)

	respBytes, _ := response.Bytes()
	_, err = w.Write(respBytes)
	if err != nil {
		log.Error(err)
	}
}

func (s *Server) handleBatchRequest(w http.ResponseWriter, data []byte) {
	requests, err := s.parseRequests(data)
	if err != nil {
		handleError(w, err)
		return
	}

	responses := make([]Response, 0, len(requests))

	for _, request := range requests {
		response := s.handler.Handle(request)
		responses = append(responses, response)
	}

	respBytes, _ := json.Marshal(responses)
	_, err = w.Write(respBytes)
	if err != nil {
		log.Error(err)
	}
}

func (s *Server) parseRequest(data []byte) (Request, error) {
	var req Request

	if err := json.Unmarshal(data, &req); err != nil {
		return Request{}, newInvalidRequestError("Invalid json request")
	}

	return req, nil
}

func (s *Server) parseRequests(data []byte) ([]Request, error) {
	var requests []Request

	if err := json.Unmarshal(data, &requests); err != nil {
		return nil, newInvalidRequestError("Invalid json request")
	}

	return requests, nil
}

func handleError(w http.ResponseWriter, err error) {
	log.Error(err)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		log.Error(err)
	}
}
