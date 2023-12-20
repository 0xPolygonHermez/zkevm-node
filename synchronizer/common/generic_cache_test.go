package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheGet(t *testing.T) {
	timerProvider := &MockTimerProvider{}
	cache := NewCache[string, string](timerProvider, time.Hour)

	// Add an item to the cache
	cache.Set("key1", "value1")

	// Test that the item can be retrieved from the cache
	value, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", value)

	// Test that an item that doesn't exist in the cache returns false
	_, ok = cache.Get("key2")
	assert.False(t, ok)

	// Test that an item that has expired is removed from the cache
	timerProvider.now = time.Now().Add(2 * time.Hour)
	_, ok = cache.Get("key1")
	assert.False(t, ok)
}

func TestCacheGetOrDefault(t *testing.T) {
	noExistsString := "no_exists"
	timerProvider := &MockTimerProvider{}
	cache := NewCache[string, string](timerProvider, time.Hour)

	// Add an item to the cache
	cache.Set("key1", "value1")

	// Test that the item can be retrieved from the cache
	value := cache.GetOrDefault("key1", noExistsString)
	assert.Equal(t, "value1", value)

	// Test that an item that doesn't exist in the cache returns false
	value = cache.GetOrDefault("key2", noExistsString)
	assert.Equal(t, noExistsString, value)

	// Test that an item that has expired is removed from the cache
	timerProvider.now = time.Now().Add(2 * time.Hour)
	value = cache.GetOrDefault("key1", noExistsString)
	assert.Equal(t, noExistsString, value)
}

func TestCacheSet(t *testing.T) {
	timerProvider := &MockTimerProvider{}
	cache := NewCache[string, string](timerProvider, time.Hour)

	// Add an item to the cache
	cache.Set("key1", "value1")

	// Test that the item can be retrieved from the cache
	value, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", value)

	// Test that an item that doesn't exist in the cache returns false
	_, ok = cache.Get("key2")
	assert.False(t, ok)

	// Test that an item that has expired is removed from the cache
	timerProvider.now = time.Now().Add(2 * time.Hour)
	_, ok = cache.Get("key1")
	assert.False(t, ok)

	// Test that an item can be updated in the cache
	cache.Set("key1", "value2")
	value, ok = cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value2", value)

	// Test that a new item can be added to the cache
	cache.Set("key2", "value3")
	value, ok = cache.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, "value3", value)
}

func TestCacheDelete(t *testing.T) {
	timerProvider := &MockTimerProvider{}
	cache := NewCache[string, string](timerProvider, time.Hour)

	// Add an item to the cache
	cache.Set("key1", "value1")

	// Delete the item from the cache
	cache.Delete("key1")

	// Test that the item has been removed from the cache
	_, ok := cache.Get("key1")
	assert.False(t, ok)

	// Test that deleting a non-existent item does not cause an error
	cache.Delete("key2")
}
func TestCacheClear(t *testing.T) {
	timerProvider := &MockTimerProvider{}
	cache := NewCache[string, string](timerProvider, time.Hour)

	// Add some items to the cache
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// Clear the cache
	cache.Clear()

	// Test that all items have been removed from the cache
	_, ok := cache.Get("key1")
	assert.False(t, ok)
	_, ok = cache.Get("key2")
	assert.False(t, ok)
	_, ok = cache.Get("key3")
	assert.False(t, ok)
}

func TestCacheDeleteOutdated(t *testing.T) {
	timerProvider := &MockTimerProvider{}
	cache := NewCache[string, string](timerProvider, time.Hour)
	now := time.Now()
	timerProvider.now = now
	// Add some items to the cache
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	timerProvider.now = now.Add(2 * time.Hour)
	cache.Set("key3", "value3")

	// Call DeleteOutdated to remove the outdated items
	cache.DeleteOutdated()
	assert.Equal(t, 1, cache.Len())

	// Test that key1 and key2 have been removed, but key3 is still present
	_, ok := cache.Get("key1")
	assert.False(t, ok)
	_, ok = cache.Get("key2")
	assert.False(t, ok)
	_, ok = cache.Get("key3")
	assert.True(t, ok)
}

func TestCacheGetDoesntReturnsOutdatedValues(t *testing.T) {
	timerProvider := &MockTimerProvider{}
	cache := NewCache[string, string](timerProvider, time.Hour)
	now := time.Now()
	timerProvider.now = now
	// Add some items to the cache
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	timerProvider.now = now.Add(2 * time.Hour)
	cache.Set("key3", "value3")

	// Test that key1 and key2 are outdated, but key3 is still present
	_, ok := cache.Get("key1")
	assert.False(t, ok)
	_, ok = cache.Get("key2")
	assert.False(t, ok)
	_, ok = cache.Get("key3")
	assert.True(t, ok)
}

func TestCacheGetExtendsTimeOfLiveOfItems(t *testing.T) {
	timerProvider := &MockTimerProvider{}
	cache := NewCache[string, string](timerProvider, time.Hour)
	now := time.Now()
	timerProvider.now = now
	// Add some items to the cache
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	timerProvider.now = now.Add(59 * time.Minute)
	_, ok := cache.Get("key1")
	assert.True(t, ok)
	timerProvider.now = now.Add(61 * time.Minute)
	cache.Set("key3", "value3")

	// Test that key1 have been extended,  key2 are outdated, and key3 is still present
	_, ok = cache.Get("key1")
	assert.True(t, ok)
	_, ok = cache.Get("key2")
	assert.False(t, ok)
	_, ok = cache.Get("key3")
	assert.True(t, ok)
}
