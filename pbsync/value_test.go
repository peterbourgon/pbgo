package pbsync_test

import (
	"sync"
	"testing"

	"github.com/peterbourgon/pbgo/pbsync"
	"github.com/peterbourgon/pbgo/pbtest/assert"
)

func TestValue(t *testing.T) {
	s := pbsync.NewValue(1)
	assert.Equal(t, 1, s.Get())

	prev := s.Set(2)
	assert.Equal(t, 1, prev)
	assert.Equal(t, 2, s.Get())
}

func TestValueZeroValue(t *testing.T) {
	var s pbsync.Value[int]
	assert.Equal(t, 0, s.Get())

	prev := s.Set(42)
	assert.Equal(t, 0, prev)
	assert.Equal(t, 42, s.Get())
}

func TestValueConcurrent(t *testing.T) {
	s := pbsync.NewValue(0)
	var wg sync.WaitGroup
	for i := range 10 {
		wg.Go(func() {
			for range 1000 {
				s.Set(i)
			}
		})
	}
	wg.Wait()
}
