package pbsync

import (
	"iter"
	"maps"
	"sync"
)

// AtomicMap is a map of K to V that is safe for concurrent use by multiple
// goroutines. Operations are synchronized with a [sync.RWMutex]. The zero value
// is a valid empty map.
type AtomicMap[K comparable, V any] struct {
	mtx sync.RWMutex
	m   map[K]V
}

// Len returns the number of entries in the map.
func (am *AtomicMap[K, V]) Len() int {
	am.mtx.RLock()
	defer am.mtx.RUnlock()
	return len(am.m)
}

// Has returns true if k exists in the map.
func (am *AtomicMap[K, V]) Has(k K) bool {
	_, ok := am.Get(k)
	return ok
}

// Get returns the value for k and if it exists.
func (am *AtomicMap[K, V]) Get(k K) (V, bool) {
	am.mtx.RLock()
	defer am.mtx.RUnlock()
	if am.m == nil {
		return *new(V), false
	}
	v, ok := am.m[k]
	return v, ok
}

// GetOrCreate returns the value for k and true, if k exists. Otherwise, it
// calls the constructor, stores the result as k, and returns the constructed
// value and false.
func (am *AtomicMap[K, V]) GetOrCreate(k K, constructor func() V) (V, bool) {
	am.mtx.Lock()
	defer am.mtx.Unlock()
	if am.m == nil {
		am.m = make(map[K]V)
	}
	v, ok := am.m[k]
	if !ok {
		v = constructor()
		am.m[k] = v
	}
	return v, ok
}

// Set key k to value v, returning the previous value for k and true if it
// existed, or a zero value and false otherwise.
func (am *AtomicMap[K, V]) Set(k K, v V) (previous V, existed bool) {
	am.mtx.Lock()
	defer am.mtx.Unlock()
	if am.m == nil {
		am.m = make(map[K]V)
	}
	prev, ok := am.m[k]
	am.m[k] = v
	return prev, ok
}

// Del deletes k from the map, returning the deleted value and true if it
// existed, or a zero value and false otherwise.
func (am *AtomicMap[K, V]) Del(k K) (deleted V, existed bool) {
	am.mtx.Lock()
	defer am.mtx.Unlock()
	if am.m == nil {
		return *new(V), false
	}
	prev, ok := am.m[k]
	delete(am.m, k)
	return prev, ok
}

// Range calls f for each entry in the map under a single read lock. If f
// returns a non-nil error, range stops and returns that error immediately.
func (am *AtomicMap[K, V]) Range(f func(K, V) error) error {
	am.mtx.RLock()
	defer am.mtx.RUnlock()
	for k, v := range am.m {
		if err := f(k, v); err != nil {
			return err
		}
	}
	return nil
}

// All returns an iterator over a snapshot of the map.
func (am *AtomicMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		am.mtx.RLock()
		snapshot := maps.Clone(am.m)
		am.mtx.RUnlock()

		for k, v := range snapshot {
			if !yield(k, v) {
				return
			}
		}
	}
}

// Apply calls f with the current value for k and true if k exists in the map,
// or the zero value of V and false if k doesn't exist. If f returns a value and
// true, that value is stored in the map for k; otherwise, nothing happens.
func (am *AtomicMap[K, V]) Apply(k K, f func(V, bool) (V, bool)) {
	am.mtx.Lock()
	defer am.mtx.Unlock()
	v, exists := am.m[k]
	vv, ok := f(v, exists)
	if !ok {
		return
	}
	if am.m == nil {
		am.m = make(map[K]V)
	}
	am.m[k] = vv
}

//
//
//

// AtomicMapMap is a two-level map of K1 to K2 to V that is safe for concurrent
// use by multiple goroutines. The zero value is an empty map of maps, ready to
// use. Inner maps are pruned once they become empty.
type AtomicMapMap[K1, K2 comparable, V any] struct {
	mtx sync.RWMutex
	mm  map[K1]*AtomicMap[K2, V]
}

// Len returns the number of outer keys in the map.
func (amm *AtomicMapMap[K1, K2, V]) Len() int {
	amm.mtx.RLock()
	defer amm.mtx.RUnlock()
	return len(amm.mm)
}

// LenOf returns the number of inner keys stored under k1.
func (amm *AtomicMapMap[K1, K2, V]) LenOf(k1 K1) int {
	amm.mtx.RLock()
	defer amm.mtx.RUnlock()
	m, ok := amm.mm[k1]
	if !ok {
		return 0
	}
	return m.Len()
}

// Has returns true if k1 and k2 exist in the map.
func (amm *AtomicMapMap[K1, K2, V]) Has(k1 K1, k2 K2) bool {
	_, ok := amm.Get(k1, k2)
	return ok
}

// Get returns the value for k1 and k2, and true if it exists.
func (amm *AtomicMapMap[K1, K2, V]) Get(k1 K1, k2 K2) (V, bool) {
	m, ok := func() (*AtomicMap[K2, V], bool) {
		amm.mtx.RLock()
		defer amm.mtx.RUnlock()
		if amm.mm == nil {
			return nil, false
		}
		m, ok := amm.mm[k1]
		return m, ok
	}()
	if !ok {
		return *new(V), false
	}
	return m.Get(k2)
}

// GetOrCreate returns the existing value for k1 and k2, and true, if it exists.
// Otherwise, it stores and returns the result of constructor, and false.
func (amm *AtomicMapMap[K1, K2, V]) GetOrCreate(k1 K1, k2 K2, constructor func() V) (V, bool) {
	amm.mtx.Lock()
	defer amm.mtx.Unlock()
	if amm.mm == nil {
		amm.mm = make(map[K1]*AtomicMap[K2, V])
	}
	m, ok := amm.mm[k1]
	if !ok {
		m = &AtomicMap[K2, V]{}
		amm.mm[k1] = m
	}
	return m.GetOrCreate(k2, constructor)
}

// Set stores v for k1 and k2, and returns the previous value and if it existed.
func (amm *AtomicMapMap[K1, K2, V]) Set(k1 K1, k2 K2, v V) (previous V, existed bool) {
	amm.mtx.Lock()
	defer amm.mtx.Unlock()
	if amm.mm == nil {
		amm.mm = make(map[K1]*AtomicMap[K2, V])
	}
	m, ok := amm.mm[k1]
	if !ok {
		m = &AtomicMap[K2, V]{}
		amm.mm[k1] = m
	}
	return m.Set(k2, v)
}

// Del deletes k2 from the inner map at k1, pruning the inner map if it
// becomes empty, and returns the previous value and whether it existed.
func (amm *AtomicMapMap[K1, K2, V]) Del(k1 K1, k2 K2) (previous V, existed bool) {
	amm.mtx.Lock()
	defer amm.mtx.Unlock()
	if amm.mm == nil {
		return *new(V), false
	}
	m, ok := amm.mm[k1]
	if !ok {
		return *new(V), false
	}
	prev, ok := m.Del(k2)
	if m.Len() <= 0 {
		delete(amm.mm, k1)
	}
	return prev, ok
}

// All returns an iterator over a snapshot of the outer map taken when the
// iterator is invoked; each inner sequence is in turn a snapshot of its
// [AtomicMap], taken when that inner sequence is iterated.
func (amm *AtomicMapMap[K1, K2, V]) All() iter.Seq2[K1, iter.Seq2[K2, V]] {
	return func(yield func(K1, iter.Seq2[K2, V]) bool) {
		amm.mtx.RLock()
		snapshot := maps.Clone(amm.mm)
		amm.mtx.RUnlock()
		for k1, m := range snapshot {
			if !yield(k1, m.All()) {
				return
			}
		}
	}
}
