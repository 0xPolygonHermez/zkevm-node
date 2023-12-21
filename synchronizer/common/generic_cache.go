package common

import (
	"time"
)

type cacheItem[T any] struct {
	value     T
	validTime time.Time
}

// Cache is a generic cache implementation with TOL (time of live) for each item
type Cache[K comparable, T any] struct {
	data            map[K]cacheItem[T] // map[K]T is a map with key type K and value type T
	timeOfLiveItems time.Duration
	timerProvider   TimeProvider
}

// NewCache creates a new cache
func NewCache[K comparable, T any](timerProvider TimeProvider, timeOfLiveItems time.Duration) *Cache[K, T] {
	return &Cache[K, T]{
		data:            make(map[K]cacheItem[T]),
		timeOfLiveItems: timeOfLiveItems,
		timerProvider:   timerProvider}
}

// Get returns the value of the key and true if the key exists and is not outdated
func (c *Cache[K, T]) Get(key K) (T, bool) {
	item, ok := c.data[key]
	if !ok {
		var zeroValue T
		return zeroValue, false
	}
	// If the item is outdated, return zero value and remove from cache
	if item.validTime.Before(c.timerProvider.Now()) {
		delete(c.data, key)
		var zeroValue T
		return zeroValue, false
	}
	// We extend the life of the item if it is used
	item.validTime = c.timerProvider.Now().Add(c.timeOfLiveItems)
	c.data[key] = item
	return item.value, true
}

// GetOrDefault returns the value of the key and defaultValue if the key does not exist or is outdated
func (c *Cache[K, T]) GetOrDefault(key K, defaultValue T) T {
	item, ok := c.Get(key)
	if !ok {
		return defaultValue
	}
	return item
}

// Set sets the value of the key
func (c *Cache[K, T]) Set(key K, value T) {
	c.data[key] = cacheItem[T]{value: value, validTime: c.timerProvider.Now().Add(c.timeOfLiveItems)}
}

// Delete deletes the key from the cache
func (c *Cache[K, T]) Delete(key K) {
	delete(c.data, key)
}

// Len returns the number of items in the cache
func (c *Cache[K, T]) Len() int {
	return len(c.data)
}

// Keys returns the keys of the cache
func (c *Cache[K, T]) Keys() []K {
	keys := make([]K, 0, len(c.data))
	for k := range c.data {
		keys = append(keys, k)
	}
	return keys
}

// Values returns the values of the cache
func (c *Cache[K, T]) Values() []T {
	values := make([]T, 0, len(c.data))
	for _, v := range c.data {
		values = append(values, v.value)
	}
	return values
}

// Clear clears the cache
func (c *Cache[K, T]) Clear() {
	c.data = make(map[K]cacheItem[T])
}

// DeleteOutdated deletes the outdated items from the cache
func (c *Cache[K, T]) DeleteOutdated() {
	for k, v := range c.data {
		if isOutdated(v.validTime, c.timerProvider.Now()) {
			delete(c.data, k)
		}
	}
}

func isOutdated(validTime time.Time, now time.Time) bool {
	return validTime.Before(now)
}

// RenewEntry renews the entry of the key
func (c *Cache[K, T]) RenewEntry(key K, validTime time.Time) {
	item, ok := c.data[key]
	if ok {
		item.validTime = c.timerProvider.Now().Add(c.timeOfLiveItems)
		c.data[key] = item
	}
}
