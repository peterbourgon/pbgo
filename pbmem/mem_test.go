package pbmem

import (
	"testing"

	"github.com/peterbourgon/pbgo/pbtest/assert"
)

func TestSize(t *testing.T) {
	t.Run("*[1024]byte", func(t *testing.T) {
		constr := func() *[1024]byte { return new([1024]byte{}) }
		want := uint64(1024)
		have := Size(constr, 32)
		assert.Within(t, want, 0.05, have)
	})

	t.Run("*fixed128", func(t *testing.T) {
		type fixed128 struct{ words [16]uint64 } // exactly 128 bytes
		constr := func() *fixed128 { return new(fixed128) }
		want := uint64(128)
		have := Size(constr, 32)
		assert.Within(t, want, 0.05, have)
	})

	t.Run("map[int]fixed1024", func(t *testing.T) {
		type fixed1024 struct{ words [128]uint64 } // exactly 1024 bytes
		constrFor := func(entries int) func() map[int]fixed1024 {
			return func() map[int]fixed1024 {
				m := make(map[int]fixed1024)
				for i := range entries {
					m[i] = fixed1024{}
				}
				return m
			}
		}
		t.Run("n=0", func(t *testing.T) {
			want := uint64(48) // empty map with initial buckets
			have := Size(constrFor(0), 32)
			assert.Within(t, want, 0.10, have)
		})
		t.Run("n=100", func(t *testing.T) {
			want := uint64(48) + 100*(8+1024)
			have := Size(constrFor(100), 32)
			assert.Within(t, want, 0.10, have)
		})
		t.Run("n=1000", func(t *testing.T) {
			want := uint64(48) + 1000*(8+1024)
			have := Size(constrFor(1000), 32)
			assert.Within(t, want, 0.10, have)
		})
		t.Run("n=2500", func(t *testing.T) {
			want := uint64(48) + 2500*(8+1024)
			have := Size(constrFor(2500), 32)
			assert.Within(t, want, 0.10, have)
		})
	})
}
