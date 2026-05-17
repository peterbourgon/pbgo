// Package memory provides test helpers for measuring runtime memory usage.
package memory

import (
	"runtime"
	"slices"
)

// Stats summarizes a series of heap-bytes measurements produced by Measure.
// All byte counts are heap bytes retained by one value returned from constr.
type Stats struct {
	N      int    // number of samples after discarding the first
	Min    uint64 // smallest sample
	Max    uint64 // largest sample
	Median uint64 // middle sample after sorting
	Mean   uint64 // arithmetic mean
}

// Measure samples the heap bytes retained by one value returned from constr,
// n times, and returns summary statistics. The first sample is discarded to
// exclude one-time bookkeeping (e.g. the internal retention slice's backing
// array, lazy first-call initialization inside constr). Panics if n <= 0.
//
// Each sample brackets a single constr() call with two runtime.GC() calls
// before and after reading runtime.MemStats. The second GC in each pair
// blocks until the first GC's sweep phase has completed, reducing bias from
// concurrent sweep activity that would otherwise leak between samples.
func Measure[T any](constr func() T, n int) Stats {
	n = max(1, n)

	objects := make([]T, 0, n+1)
	diffs := make([]uint64, 0, n+1)

	runtime.GC()
	runtime.GC()

	for range n + 1 {
		runtime.GC()
		runtime.GC()
		var m1 runtime.MemStats
		runtime.ReadMemStats(&m1)

		objects = append(objects, constr())

		runtime.GC()
		runtime.GC()
		var m2 runtime.MemStats
		runtime.ReadMemStats(&m2)

		var diff uint64
		if m2.HeapAlloc > m1.HeapAlloc {
			diff = m2.HeapAlloc - m1.HeapAlloc
		}
		diffs = append(diffs, diff)
	}

	runtime.KeepAlive(objects)

	samples := diffs[1:]
	slices.Sort(samples)

	var sum uint64
	for _, s := range samples {
		sum += s
	}

	return Stats{
		N:      len(samples),
		Min:    samples[0],
		Max:    samples[len(samples)-1],
		Median: samples[len(samples)/2],
		Mean:   sum / uint64(len(samples)),
	}
}

// Size returns the median heap bytes retained by one value returned from
// constr, sampled n times. It is a thin wrapper over Measure that returns
// only the median; use Measure directly to see min, max, mean, and sample
// count. Panics if n <= 0.
func Size[T any](constr func() T, n int) uint64 {
	return Measure(constr, n).Median
}
