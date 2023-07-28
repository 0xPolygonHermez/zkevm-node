package sequencer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/didip/tollbooth/v6"
)

const (
	stopAfterCurrentBatchEndpoint = "/stopAfterCurrentBatch"
	stopAtBatchEndpoint           = "/stopAtBatch"
	resumeProcessingEndpoint      = "/resumeProcessing"
	getCurrentBatchNumberEndpoint = "/getCurrentBatchNumber"
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
	switch req.Method {
	case "GET":
		if req.URL.Path != getCurrentBatchNumberEndpoint {
			response := map[string]string{
				"message": "zkEVM Sequencer",
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Error(err)
			}
		} else {
			s.getCurrentBatchNumber(w)
		}

	case "OPTIONS":
		// No body response for OPTIONS
		return

	case "POST":
		switch req.URL.Path {
		case stopAfterCurrentBatchEndpoint:
			s.stopAfterCurrentBatch(w)
		case stopAtBatchEndpoint:
			s.stopAtBatch(w, req)
		case resumeProcessingEndpoint:
			s.resumeProcessing(w)
		default:
			err := errors.New("invalid path " + req.URL.Path)
			s.handleInvalidRequest(w, err)
		}

	default:
		err := errors.New("method " + req.Method + " not allowed")
		s.handleInvalidRequest(w, err)
	}
}

func (s *Sequencer) handleInvalidRequest(w http.ResponseWriter, err error) {
	log.Error(err)

	w.WriteHeader(http.StatusBadRequest)
	response := map[string]string{
		"error": err.Error(),
	}
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error(err)
	}

	fmt.Println(w)
}

func (s *Sequencer) stopAfterCurrentBatch(w http.ResponseWriter) {
	s.finalizer.stopAfterCurrentBatch()

	response := map[string]string{
		"message": "Stopping after current batch",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error(err)
	}
}

type batchRequest struct {
	BatchNumber uint64 `json:"batchNumber"`
}

func (s *Sequencer) stopAtBatch(w http.ResponseWriter, req *http.Request) {
	var batchReq batchRequest
	err := json.NewDecoder(req.Body).Decode(&batchReq)
	if err != nil {
		s.handleInvalidRequest(w, errors.New("invalid request body"))
		log.Error(err)
		return
	}

	s.finalizer.stopAtBatch(batchReq.BatchNumber)

	response := map[string]string{
		"message": "Stopping at specific batch",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error(err)
	}
}

func (s *Sequencer) resumeProcessing(w http.ResponseWriter) {
	s.finalizer.resumeProcessing()

	response := map[string]string{
		"message": "Resuming processing",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error(err)
	}
}

func (s *Sequencer) getCurrentBatchNumber(w http.ResponseWriter) {
	currBatchNumber := s.finalizer.getCurrentBatchNumber()

	response := map[string]string{
		"currentBatchNumber": strconv.FormatUint(currBatchNumber, 10),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error(err)
	}
}
