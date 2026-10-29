package pbsync

import "sync"

// Safe provides mutex-guarded access to a value of type T.
type Safe[T any] struct {
	mtx sync.Mutex
	val T
}

// NewSafe returns a new [Safe] initialized with val.
func NewSafe[T any](val T) *Safe[T] {
	return &Safe[T]{
		val: val,
	}
}

// Get returns the current value.
func (s *Safe[T]) Get() T {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.val
}

// Set stores val and returns the previous value.
func (s *Safe[T]) Set(val T) T {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	prev := s.val
	s.val = val
	return prev
}
