package jsonrpc

import (
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// ErrNotFound represent a not found error.
var ErrNotFound = errors.New("object not found")

// ErrFilterInvalidPayload indicates there is an invalid payload when creating a filter
var ErrFilterInvalidPayload = errors.New("invalid argument 0: cannot specify both BlockHash and FromBlock/ToBlock, choose one or the other")

// Storage uses memory to store the data
// related to the json rpc server
type Storage struct {
	filters map[string]*Filter
}

// NewStorage creates and initializes an instance of Storage
func NewStorage() *Storage {
	return &Storage{
		filters: make(map[string]*Filter),
	}
}

// NewLogFilter persists a new log filter
func (s *Storage) NewLogFilter(wsConn *websocket.Conn, filter LogFilter) (string, error) {
	if filter.BlockHash != nil && (filter.FromBlock != nil || filter.ToBlock != nil) {
		return "", ErrFilterInvalidPayload
	}

	return s.createFilter(FilterTypeLog, filter, wsConn)
}

// NewBlockFilter persists a new block log filter
func (s *Storage) NewBlockFilter(wsConn *websocket.Conn) (string, error) {
	return s.createFilter(FilterTypeBlock, nil, wsConn)
}

// NewPendingTransactionFilter persists a new pending transaction filter
func (s *Storage) NewPendingTransactionFilter(wsConn *websocket.Conn) (string, error) {
	return s.createFilter(FilterTypePendingTx, nil, wsConn)
}

// create persists the filter to the memory and provides the filter id
func (s *Storage) createFilter(t FilterType, parameters interface{}, wsConn *websocket.Conn) (string, error) {
	lastPoll := time.Now().UTC()
	id, err := s.generateFilterID()
	if err != nil {
		return "", fmt.Errorf("failed to generate filter ID: %w", err)
	}
	s.filters[id] = &Filter{
		ID:         id,
		Type:       t,
		Parameters: parameters,
		LastPoll:   lastPoll,
		WsConn:     wsConn,
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
func (s *Storage) GetAllBlockFiltersWithWSConn() ([]*Filter, error) {
	filtersWithWSConn := []*Filter{}
	for _, filter := range s.filters {
		if filter.WsConn == nil || filter.Type != FilterTypeBlock {
			continue
		}

		f := filter
		filtersWithWSConn = append(filtersWithWSConn, f)
	}

	return filtersWithWSConn, nil
}

// GetAllLogFiltersWithWSConn returns an array with all filter that have
// a web socket connection and are filtering by new logs
func (s *Storage) GetAllLogFiltersWithWSConn() ([]*Filter, error) {
	filtersWithWSConn := []*Filter{}
	for _, filter := range s.filters {
		if filter.WsConn == nil || filter.Type != FilterTypeLog {
			continue
		}

		f := filter
		filtersWithWSConn = append(filtersWithWSConn, f)
	}

	return filtersWithWSConn, nil
}

// GetFilter gets a filter by its id
func (s *Storage) GetFilter(filterID string) (*Filter, error) {
	filter, found := s.filters[filterID]
	if !found {
		return nil, ErrNotFound
	}

	return filter, nil
}

// UpdateFilterLastPoll updates the last poll to now
func (s *Storage) UpdateFilterLastPoll(filterID string) error {
	filter, found := s.filters[filterID]
	if !found {
		return ErrNotFound
	}

	filter.LastPoll = time.Now().UTC()
	s.filters[filterID] = filter
	return nil
}

// UninstallFilter deletes a filter by its id
func (s *Storage) UninstallFilter(filterID string) error {
	_, found := s.filters[filterID]
	if !found {
		return ErrNotFound
	}

	delete(s.filters, filterID)
	return nil
}

// UninstallFilterByWSConn deletes all filters connected to the provided web socket connection
func (s *Storage) UninstallFilterByWSConn(wsConn *websocket.Conn) error {
	filterIDsToDelete := []string{}

	for id, filter := range s.filters {
		if filter.WsConn == wsConn {
			filterIDsToDelete = append(filterIDsToDelete, id)
		}
	}

	for _, filterID := range filterIDsToDelete {
		delete(s.filters, filterID)
	}

	return nil
}
