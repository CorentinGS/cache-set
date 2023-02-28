package cacheset

import "time"

// set is a map with expiration times
type set[T comparable] map[T]int64

// Expire removes the given element from the set if it has expired
func (s set[T]) Expire(elem T) {
	if s.Expired(elem) {
		s.Delete(elem)
	}
}

// ExpireAll removes all expired elements from the set
func (s set[T]) ExpireAll() {
	for k := range s {
		s.Expire(k)
	}
}

// Expired returns true if the given element has expired
func (s set[T]) Expired(elem T) bool {
	expires, ok := s[elem]
	if !ok {
		return false
	}
	if expires > 0 && expires < time.Now().UnixNano() {
		return true
	}
	return false
}

// ToSlice returns a slice of the set's elements
func (s set[T]) ToSlice() []T {
	slice := make([]T, 0, len(s))
	for k := range s {
		slice = append(slice, k)
	}
	return slice

}

// Add adds the given element to the set with the given expiration time
func (s set[T]) Add(elem T, duration time.Duration) {
	var expires int64
	if duration > 0 {
		expires = time.Now().Add(duration).UnixNano()
	} else {
		expires = 0
	}
	s[elem] = expires
}

// Clear removes all elements from the set
func (s set[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Contains returns true if the given element is in the set
func (s set[T]) Contains(elem T) bool {
	_, ok := s[elem]
	return ok
}

// Delete removes the given element from the set
func (s set[T]) Delete(elem T) {
	delete(s, elem)
}

// Len returns the number of elements in the set
func (s set[T]) Len() int {
	return len(s)
}

// New returns a new set
func newSet[T comparable]() set[T] {
	return make(set[T])
}
