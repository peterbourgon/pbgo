package pbsync_test

import (
	"math"
	"sync"
	"testing"

	"github.com/peterbourgon/pbgo/pbsync"
	"github.com/peterbourgon/pbgo/pbtest/assert"
)

func TestFloat64(t *testing.T) {
	var f pbsync.Float64
	assert.Equal(t, 0.0, f.Load())

	f.Store(1.5)
	assert.Equal(t, 1.5, f.Load())

	assert.Equal(t, 1.5, f.Swap(2.5))
	assert.Equal(t, 2.5, f.Load())

	assert.False(t, f.CompareAndSwap(1.5, 3.5))
	assert.True(t, f.CompareAndSwap(2.5, 3.5))
	assert.Equal(t, 4.0, f.Add(0.5))
}

func TestFloat64CompareAndSwapBits(t *testing.T) {
	var f pbsync.Float64
	nan := math.Float64frombits(0x7ff8000000000001)
	f.Store(nan)
	assert.True(t, f.CompareAndSwap(nan, 1))
}

func TestFloat64AddConcurrent(t *testing.T) {
	var f pbsync.Float64
	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			for range 1000 {
				f.Add(0.5)
			}
		})
	}
	wg.Wait()

	assert.Equal(t, 5000.0, f.Load())
}
