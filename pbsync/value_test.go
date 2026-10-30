package pbsync_test

import (
	"sync"
	"testing"

	"github.com/peterbourgon/pbgo/pbsync"
	"github.com/peterbourgon/pbgo/pbtest/assert"
)

func TestValue(t *testing.T) {
	v := pbsync.NewValue(1)
	assert.Equal(t, 1, v.Get())

	prev := v.Set(2)
	assert.Equal(t, 1, prev)
	assert.Equal(t, 2, v.Get())
}

func TestValueZeroValue(t *testing.T) {
	var v pbsync.Value[int]
	assert.Equal(t, 0, v.Get())

	prev := v.Set(42)
	assert.Equal(t, 0, prev)
	assert.Equal(t, 42, v.Get())
}

func TestValueConcurrent(t *testing.T) {
	v := pbsync.NewValue(0)
	var wg sync.WaitGroup
	for i := range 10 {
		wg.Go(func() {
			for range 1000 {
				v.Set(i)
			}
		})
	}
	wg.Wait()
}
