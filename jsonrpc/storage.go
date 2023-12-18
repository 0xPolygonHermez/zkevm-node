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
	allFilters                 map[string]*Filter
	allFiltersWithWSConn       map[*concurrentWsConn]map[string]*Filter
	blockFiltersWithWSConn     map[string]*Filter
	logFiltersWithWSConn       map[string]*Filter
	pendingTxFiltersWithWSConn map[string]*Filter

	blockMutex     *sync.Mutex
	logMutex       *sync.Mutex
	pendingTxMutex *sync.Mutex
}

// NewStorage creates and initializes an instance of Storage
func NewStorage() *Storage {
	return &Storage{
		allFilters:                 make(map[string]*Filter),
		allFiltersWithWSConn:       make(map[*concurrentWsConn]map[string]*Filter),
		blockFiltersWithWSConn:     make(map[string]*Filter),
		logFiltersWithWSConn:       make(map[string]*Filter),
		pendingTxFiltersWithWSConn: make(map[string]*Filter),
		blockMutex:                 &sync.Mutex{},
		logMutex:                   &sync.Mutex{},
		pendingTxMutex:             &sync.Mutex{},
	}
}

// NewLogFilter persists a new log filter
func (s *Storage) NewLogFilter(wsConn *concurrentWsConn, filter LogFilter) (string, error) {
	if err := filter.Validate(); err != nil {
		return "", err
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

	s.blockMutex.Lock()
	s.logMutex.Lock()
	s.pendingTxMutex.Lock()
	defer s.blockMutex.Unlock()
	defer s.logMutex.Unlock()
	defer s.pendingTxMutex.Unlock()

	f := &Filter{
		ID:            id,
		Type:          t,
		Parameters:    parameters,
		LastPoll:      lastPoll,
		WsConn:        wsConn,
		wsQueue:       state.NewQueue[[]byte](),
		wsQueueSignal: sync.NewCond(&sync.Mutex{}),
	}

	go state.InfiniteSafeRun(f.SendEnqueuedSubscriptionData, fmt.Sprintf("failed to send enqueued subscription data to filter %v", id), time.Second)

	s.allFilters[id] = f
	if f.WsConn != nil {
		if _, found := s.allFiltersWithWSConn[f.WsConn]; !found {
			s.allFiltersWithWSConn[f.WsConn] = make(map[string]*Filter)
		}

		s.allFiltersWithWSConn[f.WsConn][id] = f
		if t == FilterTypeBlock {
			s.blockFiltersWithWSConn[id] = f
		} else if t == FilterTypeLog {
			s.logFiltersWithWSConn[id] = f
		} else if t == FilterTypePendingTx {
			s.pendingTxFiltersWithWSConn[id] = f
		}
	}
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

// GetAllBlockFiltersWithWSConn returns an array with all filter that have
// a web socket connection and are filtering by new blocks
func (s *Storage) GetAllBlockFiltersWithWSConn() []*Filter {
	s.blockMutex.Lock()
	defer s.blockMutex.Unlock()

	filters := []*Filter{}
	for _, filter := range s.blockFiltersWithWSConn {
		f := filter
		filters = append(filters, f)
	}
	return filters
}

// GetAllLogFiltersWithWSConn returns an array with all filter that have
// a web socket connection and are filtering by new logs
func (s *Storage) GetAllLogFiltersWithWSConn() []*Filter {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	filters := []*Filter{}
	for _, filter := range s.logFiltersWithWSConn {
		f := filter
		filters = append(filters, f)
	}
	return filters
}

// GetFilter gets a filter by its id
func (s *Storage) GetFilter(filterID string) (*Filter, error) {
	s.blockMutex.Lock()
	s.logMutex.Lock()
	s.pendingTxMutex.Lock()
	defer s.blockMutex.Unlock()
	defer s.logMutex.Unlock()
	defer s.pendingTxMutex.Unlock()

	filter, found := s.allFilters[filterID]
	if !found {
		return nil, ErrNotFound
	}

	return filter, nil
}

// UpdateFilterLastPoll updates the last poll to now
func (s *Storage) UpdateFilterLastPoll(filterID string) error {
	s.blockMutex.Lock()
	s.logMutex.Lock()
	s.pendingTxMutex.Lock()
	defer s.blockMutex.Unlock()
	defer s.logMutex.Unlock()
	defer s.pendingTxMutex.Unlock()

	filter, found := s.allFilters[filterID]
	if !found {
		return ErrNotFound
	}
	filter.LastPoll = time.Now().UTC()
	s.allFilters[filterID] = filter
	return nil
}

// UninstallFilter deletes a filter by its id
func (s *Storage) UninstallFilter(filterID string) error {
	s.blockMutex.Lock()
	s.logMutex.Lock()
	s.pendingTxMutex.Lock()
	defer s.blockMutex.Unlock()
	defer s.logMutex.Unlock()
	defer s.pendingTxMutex.Unlock()

	filter, found := s.allFilters[filterID]
	if !found {
		return ErrNotFound
	}

	s.deleteFilter(filter)
	return nil
}

// UninstallFilterByWSConn deletes all filters connected to the provided web socket connection
func (s *Storage) UninstallFilterByWSConn(wsConn *concurrentWsConn) error {
	s.blockMutex.Lock()
	s.logMutex.Lock()
	s.pendingTxMutex.Lock()
	defer s.blockMutex.Unlock()
	defer s.logMutex.Unlock()
	defer s.pendingTxMutex.Unlock()

	filters, found := s.allFiltersWithWSConn[wsConn]
	if !found {
		return nil
	}

	for _, filter := range filters {
		s.deleteFilter(filter)
	}

	return nil
}

// deleteFilter deletes a filter from all the maps
func (s *Storage) deleteFilter(filter *Filter) {
	if filter.Type == FilterTypeBlock {
		delete(s.blockFiltersWithWSConn, filter.ID)
	} else if filter.Type == FilterTypeLog {
		delete(s.logFiltersWithWSConn, filter.ID)
	} else if filter.Type == FilterTypePendingTx {
		delete(s.pendingTxFiltersWithWSConn, filter.ID)
	}

	if filter.WsConn != nil {
		delete(s.allFiltersWithWSConn[filter.WsConn], filter.ID)
		if len(s.allFiltersWithWSConn[filter.WsConn]) == 0 {
			delete(s.allFiltersWithWSConn, filter.WsConn)
		}
	}

	delete(s.allFilters, filter.ID)
}
