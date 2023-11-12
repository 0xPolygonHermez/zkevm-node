package jsonrpc

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/google/uuid"
)

// ErrNotFound represent a not found error.
var ErrNotFound = errors.New("object not found")

// ErrFilterInvalidPayload indicates there is an invalid payload when creating a filter
var ErrFilterInvalidPayload = errors.New("invalid argument 0: cannot specify both BlockHash and FromBlock/ToBlock, choose one or the other")

// Storage uses memory to store the data
// related to the json rpc server
type Storage struct {
	LogFilters       map[string]*Filter
	BlockFilters     map[string]*Filter
	PendingTxFilters map[string]*Filter
	AllFilters       map[string]*Filter
	mutex            *sync.Mutex
}

// NewStorage creates and initializes an instance of Storage
func NewStorage() *Storage {
	return &Storage{
		LogFilters:       make(map[string]*Filter),
		BlockFilters:     make(map[string]*Filter),
		PendingTxFilters: make(map[string]*Filter),
		AllFilters:       make(map[string]*Filter),
		mutex:            &sync.Mutex{},
	}
}

// Lock() locks the access to storage data
func (s *Storage) Lock() {
	s.mutex.Lock()
}

// Lock() unlocks the access to storage data
func (s *Storage) Unlock() {
	s.mutex.Unlock()
}

// NewLogFilter persists a new log filter
func (s *Storage) NewLogFilter(wsConn *concurrentWsConn, filter LogFilter) (string, error) {
	shouldFilterByBlockHash := filter.BlockHash != nil
	shouldFilterByBlockRange := filter.FromBlock != nil || filter.ToBlock != nil

	if shouldFilterByBlockHash && shouldFilterByBlockRange {
		return "", ErrFilterInvalidPayload
	}

	return s.createFilter(FilterTypeLog, filter, wsConn)
}

// NewBlockFilter persists a new block log filter
func (s *Storage) NewBlockFilter(wsConn *concurrentWsConn) (string, error) {
	return s.createFilter(FilterTypeBlock, nil, wsConn)
}

// NewPendingTransactionFilter persists a new pending transaction filter
func (s *Storage) NewPendingTransactionFilter(wsConn *concurrentWsConn) (string, error) {
	return s.createFilter(FilterTypePendingTx, nil, wsConn)
}

// create persists the filter to the memory and provides the filter id
func (s *Storage) createFilter(t FilterType, parameters interface{}, wsConn *concurrentWsConn) (string, error) {
	lastPoll := time.Now().UTC()
	id, err := s.generateFilterID()
	if err != nil {
		return "", fmt.Errorf("failed to generate filter ID: %w", err)
	}

	f := &Filter{
		ID:          id,
		Type:        t,
		Parameters:  parameters,
		LastPoll:    lastPoll,
		WsConn:      wsConn,
		wsDataQueue: state.NewQueue[[]byte](),
		mutex:       &sync.Mutex{},
	}

	s.mutex.Lock()

	s.AllFilters[id] = f

	if t == FilterTypeLog {
		s.LogFilters[id] = f
	} else if t == FilterTypeBlock {
		s.BlockFilters[id] = f
	} else if t == FilterTypePendingTx {
		s.PendingTxFilters[id] = f
	}

	s.mutex.Unlock()

	return id, nil
}

func (s *Storage) generateFilterID() (string, error) {
	r, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	b, err := r.MarshalBinary()
	if err != nil {
		return "", err
	}

	id := hex.EncodeToHex(b)
	return id, nil
}

// SendBlockFilterSubscriptionData sends the data to all block filter subscriptions
func (s *Storage) SendBlockFilterSubscriptionData(data []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, filter := range s.BlockFilters {
		if filter.WsConn != nil {
			filter.EnqueueSubscriptionDataToBeSent(data)
			go filter.SendEnqueuedSubscriptionData()
		}
	}
}

// GetFilter gets a filter by its id
func (s *Storage) GetFilter(filterID string) (*Filter, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	filter, found := s.AllFilters[filterID]
	if !found {
		return nil, ErrNotFound
	}

	return filter, nil
}

// UpdateFilterLastPoll updates the last poll to now
func (s *Storage) UpdateFilterLastPoll(filterID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	filter, found := s.AllFilters[filterID]
	if !found {
		return ErrNotFound
	}
	filter.LastPoll = time.Now().UTC()
	return nil
}

// deleteFilter deletes a filter from all the maps
func (s *Storage) deleteFilter(filter *Filter) {
	if filter.Type == FilterTypeLog {
		delete(s.LogFilters, filter.ID)
	} else if filter.Type == FilterTypeBlock {
		delete(s.BlockFilters, filter.ID)
	} else if filter.Type == FilterTypePendingTx {
		delete(s.PendingTxFilters, filter.ID)
	}

	delete(s.AllFilters, filter.ID)
}

// UninstallFilter deletes a filter by its id
func (s *Storage) UninstallFilter(filterID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	filter, found := s.AllFilters[filterID]
	if !found {
		return ErrNotFound
	}

	s.deleteFilter(filter)
	return nil
}

// UninstallFilterByWSConn deletes all filters connected to the provided web socket connection
func (s *Storage) UninstallFilterByWSConn(wsConn *concurrentWsConn) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	filtersToDelete := []*Filter{}

	for _, filter := range s.AllFilters {
		if filter.WsConn == wsConn {
			filtersToDelete = append(filtersToDelete, filter)
		}
	}

	for _, filterToDelete := range filtersToDelete {
		s.deleteFilter(filterToDelete)
	}

	return nil
}
