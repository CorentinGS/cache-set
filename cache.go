// Package cacheset contains a thread-safe map with expiration times.
//
// Path: cache.go
//
// Description: cache.go contains the Cache type and its methods.
//
// Usage:
//
//	// Create a new cache that cleans every 5 minutes
//	cache := New(5 * time.Minute)
//
//	// Add an element to the cache with a 1 minute expiration time
//	cache.Add("foo", 1 * time.Minute)
//
//	// Check if an element is in the cache
//	if cache.Contains("foo") {
//		// ...
//	}
//
//	// Delete an element from the cache
//	cache.Delete("foo")
//
//	// Get a copy of the cache's set
//	set := cache.CopySet()
//
//	// Close the cache's cleaning goroutine
//	cache.Close()
package cacheset

import (
	"sync"
	"time"
)

// Cache is a thread-safe map with expiration times.
type Cache[T comparable] struct {
	set[T]                     // set is a map with expiration times
	close        chan struct{} // close is a channel that stops the cache's cleaning goroutine
	sync.RWMutex               // RWMutex is a mutex that can be locked for reading or writing
}

// New creates a new cache that asynchronously cleans
func New[T comparable](cleanInterval time.Duration) *Cache[T] {
	c := &Cache[T]{
		set:   newSet[T](),
		close: make(chan struct{}),
	}

	ticker := time.NewTicker(cleanInterval) // ticker is a ticker that cleans the cache every cleanInterval
	defer ticker.Stop()                     // defer ticker.Stop() stops the ticker when the function returns

	go func() {
		for {
			select {
			case <-c.close: // c.close is a channel that stops the cache's cleaning goroutine
				return
			case <-ticker.C: // ticker.C is a channel that sends a value every time the ticker ticks
				c.Lock()
				c.ExpireAll() // ExpireAll expires all elements in the cache
				c.Unlock()
			}
		}
	}()

	return c
}

// CopySet returns a copy of the cache's set
//
// Description: CopySet returns a copy of the cache's set. The returned set is a map of elements to their expiration times.
// In go maps are passed by reference, so this function returns a copy of the map.
func (c *Cache[T]) CopySet() map[T]int64 {
	c.RLock()
	defer c.RUnlock()

	return c.set.Copy()
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
