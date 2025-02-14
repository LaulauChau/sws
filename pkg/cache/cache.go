package cache

import (
	"sync"
	"time"
)

// Cache is a generic cache implementation that can store any type
type Cache[T any] struct {
	mu            sync.RWMutex
	data          T
	lastUpdated   time.Time
	updateTimeout time.Duration
	initialized   bool
}

// NewCache creates a new cache instance with the specified timeout
func NewCache[T any](updateTimeout time.Duration) *Cache[T] {
	return &Cache[T]{
		updateTimeout: updateTimeout,
	}
}

// Get retrieves the cached data and whether it's valid
func (c *Cache[T]) Get() (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var zero T
	if !c.initialized || time.Since(c.lastUpdated) > c.updateTimeout {
		return zero, false
	}

	return c.data, true
}

// Set updates the cached data
func (c *Cache[T]) Set(data T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = data
	c.lastUpdated = time.Now()
	c.initialized = true
}
