// Package pbsync provides supplementary synchronization primitives.
package pbsync

import (
	"math"
	"sync/atomic"
)

// Float64 provides atomic operations for float64 values.
type Float64 struct{ u atomic.Uint64 }

// Load atomically loads and returns the value stored in x.
func (x *Float64) Load() float64 {
	bits := x.u.Load()
	return math.Float64frombits(bits)
}

// Store atomically stores val into x.
func (x *Float64) Store(val float64) {
	valbits := math.Float64bits(val)
	x.u.Store(valbits)
}

// Swap atomically stores new into x and returns the previous value.
func (x *Float64) Swap(new float64) (old float64) {
	newbits := math.Float64bits(new)
	oldbits := x.u.Swap(newbits)
	return math.Float64frombits(oldbits)
}

// CompareAndSwap executes the compare-and-swap operation for x.
func (x *Float64) CompareAndSwap(old, new float64) (swapped bool) {
	oldbits := math.Float64bits(old)
	newbits := math.Float64bits(new)
	return x.u.CompareAndSwap(oldbits, newbits)
}

// Add atomically adds delta to x and returns the new value.
func (x *Float64) Add(delta float64) (new float64) {
	for {
		oldval := x.Load()
		newval := oldval + delta
		if x.CompareAndSwap(oldval, newval) {
			return newval
		}
	}
}
