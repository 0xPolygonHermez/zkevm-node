package tree

import (
	"errors"
	"sync"

	"github.com/dgraph-io/ristretto"
)

const (
	cacheNumCounters = 1e7     // number of keys to track frequency of (10M).
	cacheMaxCost     = 1 << 30 // maximum cost of cache (1GB).
	cacheDefaultCost = 1       // cost for regular kv items
	cacheBufferItems = 64      // number of keys per Get buffer.
)

// NewStoreCache creates a new cache object to be used in MT backends.
func NewStoreCache() (*ristretto.Cache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: cacheNumCounters,
		MaxCost:     cacheMaxCost,
		BufferItems: cacheBufferItems,
	})
	if err != nil {
		return nil, err
	}
	return cache, err
}

// nodeCache represents a cache layer to store nodes during a db transaction in
// the merkletree.
type nodeCache struct {
	lock   *sync.Mutex
	data   map[string][]uint64
	active bool
}

const (
	maxMTNodeCacheEntries = 256
)

var (
	errMTNodeCacheItemNotFound = errors.New("MT node cache item not found")
)

// newNodeCache is the nodeCache constructor.
func newNodeCache() *nodeCache {
	return &nodeCache{
		lock: &sync.Mutex{},
		data: make(map[string][]uint64),
	}
}

// get reads a MT node cache entry.
func (nc *nodeCache) get(key []uint64) ([]uint64, error) {
	keyStr := h4ToString(key)

	item, ok := nc.data[keyStr]
	if !ok {
		return nil, errMTNodeCacheItemNotFound
	}
	return item, nil
}

// set inserts a new MT node cache entry.
func (nc *nodeCache) set(key []uint64, value []uint64) error {
	if len(nc.data) >= maxMTNodeCacheEntries {
		return errors.New("MT node cache is full")
	}
	keyStr := h4ToString(key)

	nc.lock.Lock()
	defer nc.lock.Unlock()

	nc.data[keyStr] = value

	return nil
}

// clear removes all the entries of the MT node cache.
func (nc *nodeCache) clear() {
	nc.lock.Lock()
	defer nc.lock.Unlock()

	for k := range nc.data {
		delete(nc.data, k)
	}
}

// isActive is the active field getter.
func (nc *nodeCache) isActive() bool {
	return nc.active
}

// setActive is the active field setter.
func (nc *nodeCache) setActive(active bool) {
	nc.active = active
}

// init initializes the MT node cache.
func (nc *nodeCache) init() {
	nc.clear()
	nc.setActive(true)
}

// teardown resets the MT node cache.
func (nc *nodeCache) teardown() {
	nc.clear()
	nc.setActive(false)
}
