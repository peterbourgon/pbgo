package pbrand_test

import (
	"encoding"
	"io"
	"math/rand/v2"
	"testing"

	"github.com/peterbourgon/pbgo/pbrand"
	"github.com/peterbourgon/pbgo/pbtest/assert"
)

func TestChaCha8RandDeterministic(t *testing.T) {
	a := pbrand.NewChaCha8Rand(42)
	b := pbrand.NewChaCha8Rand(42)
	for range 100 {
		assert.Equal(t, a.Uint64(), b.Uint64())
	}
}

func TestChaCha8RandDistinctSeeds(t *testing.T) {
	a := pbrand.NewChaCha8Rand(1)
	b := pbrand.NewChaCha8Rand(2)
	assert.True(t, a.Uint64() != b.Uint64())
}

func TestChaCha8RandRead(t *testing.T) {
	r := pbrand.NewChaCha8Rand(7)
	buf := make([]byte, 64)
	n, err := r.Read(buf)
	assert.NoErr(t, err)
	assert.Equal(t, len(buf), n)

	var nonzero int
	for _, b := range buf {
		if b != 0 {
			nonzero++
		}
	}
	assert.True(t, nonzero > 0)
}

func TestChaCha8RandRandMethods(t *testing.T) {
	r := pbrand.NewChaCha8Rand(99)
	assert.Between(t, 0, 9, r.IntN(10))
	f := r.Float64()
	assert.True(t, f >= 0 && f < 1)
}

func TestChaCha8RandInterfaces(t *testing.T) {
	r := pbrand.NewChaCha8Rand(1)
	var (
		_ rand.Source                = r
		_ io.Reader                  = r
		_ encoding.BinaryUnmarshaler = r
	)
}
