package cacheset

import (
	"sync"
	"time"
)

// Cache is a thread-safe map with expiration times.
type Cache[T comparable] struct {
	set[T]
	sync.RWMutex
	close chan struct{}
}

// New creates a new cache that asynchronously cleans
func New[T comparable](cleanInterval time.Duration) *Cache[T] {
	c := &Cache[T]{
		set:   newSet[T](),
		close: make(chan struct{}),
	}

	ticker := time.NewTicker(cleanInterval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-c.close:
				return
			case <-ticker.C:
				c.Lock()
				c.ExpireAll()
				c.Unlock()
			}
		}
	}()

	return c
}

// CopySet returns a copy of the cache's set
func (c *Cache[T]) CopySet() map[T]int64 {
	c.RLock()
	defer c.RUnlock()

	return c.set
}

// Delete removes the given element from the cache
func (c *Cache[T]) Delete(elem T) {
	c.Lock()
	defer c.Unlock()

	c.set.Delete(elem)
}

// Len returns the number of elements in the cache
func (c *Cache[T]) Len() int {
	c.RLock()
	defer c.RUnlock()

	return c.set.Len()
}

// Close stops the cache's cleaning goroutine
func (c *Cache[T]) Close() {
	c.close <- struct{}{}
	close(c.close)
	c.set = nil
}

// Add adds the given element to the cache
func (c *Cache[T]) Add(elem T, duration time.Duration) {
	c.Lock()
	defer c.Unlock()

	c.set.Add(elem, duration)
}

// Contains returns true if the given element is in the cache
func (c *Cache[T]) Contains(elem T) bool {
	c.RLock()
	defer c.RUnlock()

	return c.set.Contains(elem)
}

// ToSlice returns a slice of all elements in the cache
func (c *Cache[T]) ToSlice() []T {
	c.RLock()
	defer c.RUnlock()

	return c.set.ToSlice()
}

// Clear clears the cache
func (c *Cache[T]) Clear() {
	c.Lock()
	defer c.Unlock()

	c.set.Clear()
}

// Expire expires the given element
func (c *Cache[T]) Expire(elem T) {
	c.Lock()
	defer c.Unlock()

	c.set.Expire(elem)
}

// ExpireAll expires all elements in the cache
func (c *Cache[T]) ExpireAll() {
	c.Lock()
	defer c.Unlock()

	c.set.ExpireAll()
}

// Exists returns true if the given key exists
func (c *Cache[T]) Exists(elem T) bool {
	c.RLock()
	defer c.RUnlock()

	return c.set.Contains(elem)
}
