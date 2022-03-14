package tree

import "github.com/dgraph-io/ristretto"

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
