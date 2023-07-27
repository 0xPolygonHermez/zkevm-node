package sequencer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/didip/tollbooth/v6"
)

// startHTTP starts a server to respond http requests
func (s *Sequencer) startHTTP() {
	if s.srv != nil {
		log.Fatal(fmt.Errorf("server already started"))
	}

	address := fmt.Sprintf("%s:%d", s.cfg.Http.Host, s.cfg.Http.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("failed to create tcp listener: %v", err)
	}

	mux := http.NewServeMux()

	lmt := tollbooth.NewLimiter(s.cfg.Http.MaxRequestsPerIPAndSecond, nil)
	mux.Handle("/", tollbooth.LimitFuncHandler(lmt, s.handle))

	s.srv = &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: s.cfg.Http.ReadTimeout.Duration,
		ReadTimeout:       s.cfg.Http.ReadTimeout.Duration,
		WriteTimeout:      s.cfg.Http.WriteTimeout.Duration,
	}
	log.Infof("http server started: %s", address)
	if err := s.srv.Serve(lis); err != nil {
		if err == http.ErrServerClosed {
			log.Infof("http server stopped")
		}
		log.Fatalf("closed http connection: %v", err)
	}
}

func (s *Sequencer) handle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if req.Method == "GET" {
		_, err := w.Write([]byte("zkEVM Sequencer Server"))
		if err != nil {
			log.Error(err)
		}
		return
	}

	if (*req).Method == "OPTIONS" {
		return
	}

	if req.Method != "POST" {
		err := errors.New("method " + req.Method + " not allowed")
		s.handleInvalidRequest(w, err)
		return
	}

	switch req.URL.Path {
	case "/stopAfterCurrentBatch":
		s.stopAfterCurrentBatch(w)
	case "/stopAtBatch":
		s.stopAtBatch(w, req)
	case "/resumeProcessing":
		s.resumeProcessing(w)
	default:
		err := errors.New("invalid path " + req.URL.Path)
		s.handleInvalidRequest(w, err)
	}
}

func (s *Sequencer) handleInvalidRequest(w http.ResponseWriter, err error) {
	log.Error(err)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		log.Error(err)
	}
}

func (s *Sequencer) stopAfterCurrentBatch(w http.ResponseWriter) {
	s.finalizer.stopAfterCurrentBatch()
	_, err := w.Write([]byte("Stopping after current batch"))
	if err != nil {
		log.Error(err)
	}
}

type BatchRequest struct {
	BatchNumber uint64 `json:"batchNumber"`
}

func (s *Sequencer) stopAtBatch(w http.ResponseWriter, req *http.Request) {
	var batchReq BatchRequest
	err := json.NewDecoder(req.Body).Decode(&batchReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Error(err)
		return
	}

	s.finalizer.stopAtBatch(batchReq.BatchNumber)
	_, err = w.Write([]byte("Stopping at specific batch"))
	if err != nil {
		log.Error(err)
	}
}

func (s *Sequencer) resumeProcessing(w http.ResponseWriter) {
	s.finalizer.resumeProcessing()
	_, err := w.Write([]byte("Resuming processing"))
	if err != nil {
		log.Error(err)
	}
}
