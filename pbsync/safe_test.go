package pbsync_test

import (
	"sync"
	"testing"

	"github.com/peterbourgon/pbgo/pbsync"
	"github.com/peterbourgon/pbgo/pbtest/assert"
)

func TestSafe(t *testing.T) {
	s := pbsync.NewSafe(1)
	assert.Equal(t, 1, s.Get())

	prev := s.Set(2)
	assert.Equal(t, 1, prev)
	assert.Equal(t, 2, s.Get())
}

func TestSafeZeroValue(t *testing.T) {
	var s pbsync.Safe[int]
	assert.Equal(t, 0, s.Get())

	prev := s.Set(42)
	assert.Equal(t, 0, prev)
	assert.Equal(t, 42, s.Get())
}

func TestSafeConcurrent(t *testing.T) {
	s := pbsync.NewSafe(0)
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
