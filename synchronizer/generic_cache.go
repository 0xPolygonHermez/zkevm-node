package synchronizer

import (
	"time"
)

type TimeProvider interface {
	Now() time.Time
}

type DefaultTimeProvider struct{}

func (d DefaultTimeProvider) Now() time.Time {
	return time.Now()
}

type CacheItem[T any] struct {
	value     T
	validTime time.Time
}

type Cache[K comparable, T any] struct {
	data            map[K]CacheItem[T] // map[K]T is a map with key type K and value type T
	timeOfLiveItems time.Duration
	timerProvider   TimeProvider
}

func NewCache[K comparable, T any](timerProvider TimeProvider, timeOfLiveItems time.Duration) *Cache[K, T] {
	return &Cache[K, T]{
		data:            make(map[K]CacheItem[T]),
		timeOfLiveItems: timeOfLiveItems,
		timerProvider:   timerProvider}
}

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
	return item.value, true
}

func (c *Cache[K, T]) Set(key K, value T) {
	c.data[key] = CacheItem[T]{value: value, validTime: c.timerProvider.Now().Add(c.timeOfLiveItems)}
}

func (c *Cache[K, T]) Delete(key K) {
	delete(c.data, key)
}

func (c *Cache[K, T]) Len() int {
	return len(c.data)
}

func (c *Cache[K, T]) Keys() []K {
	keys := make([]K, 0, len(c.data))
	for k := range c.data {
		keys = append(keys, k)
	}
	return keys
}

func (c *Cache[K, T]) Values() []T {
	values := make([]T, 0, len(c.data))
	for _, v := range c.data {
		values = append(values, v.value)
	}
	return values
}

func (c *Cache[K, T]) Clear() {
	c.data = make(map[K]CacheItem[T])
}

func (c *Cache[K, T]) DeleteOutdated() {
	for k, v := range c.data {
		if v.validTime.Before(c.timerProvider.Now()) {
			delete(c.data, k)
		}
	}
}

func (c *Cache[K, T]) RenewEntry(key K, validTime time.Time) {
	item, ok := c.data[key]
	if ok {
		item.validTime = c.timerProvider.Now().Add(c.timeOfLiveItems)
		c.data[key] = item
	}
}
