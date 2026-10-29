package pbsync

import "sync"

// Value provides mutex-guarded access to a value of type T.
type Value[T any] struct {
	mtx sync.Mutex
	val T
}

// NewValue returns a new [Value] initialized with val.
func NewValue[T any](val T) *Value[T] {
	return &Value[T]{
		val: val,
	}
}

// Get returns the current value.
func (s *Value[T]) Get() T {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.val
}

// Set stores val and returns the previous value.
func (s *Value[T]) Set(val T) T {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	prev := s.val
	s.val = val
	return prev
}
