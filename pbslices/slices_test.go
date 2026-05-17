package pbslices_test

import (
	"cmp"
	"fmt"
	"testing"

	"github.com/peterbourgon/pbgo/pbslices"
	"github.com/peterbourgon/pbgo/pbtest/assert"
)

func TestFilter(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	fn := func(x int) bool { return x%2 == 0 }
	want := []int{2, 4}
	have := pbslices.Filter(input, fn)
	assert.SlicesEqual(t, want, have)
}

func TestFirst(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		fn := func(x int) bool { return x%2 == 1 && x > 1 }
		want := 3
		have, ok := pbslices.First(input, fn)
		assert.True(t, ok)
		assert.Equal(t, want, have)
	})

	t.Run("missing", func(t *testing.T) {
		input := []int{1, 2, 3}
		fn := func(x int) bool { return x > 10 }
		_, ok := pbslices.First(input, fn)
		assert.False(t, ok)
	})
}

func TestFirstDefault(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		fn := func(x int) bool { return x%2 == 0 }
		defaultValue := 99
		want := 2
		have := pbslices.FirstDefault(input, fn, defaultValue)
		assert.Equal(t, want, have)
	})

	t.Run("missing", func(t *testing.T) {
		input := []int{1, 2, 3}
		fn := func(x int) bool { return x > 10 }
		defaultValue := 99
		want := defaultValue
		have := pbslices.FirstDefault(input, fn, defaultValue)
		assert.Equal(t, want, have)
	})
}

func TestSort(t *testing.T) {
	input := []int{3, 1, 4, 1, 5, 9}
	cmpFn := cmp.Compare[int]
	want := []int{1, 1, 3, 4, 5, 9}
	have := pbslices.Sort(input, cmpFn)
	assert.SlicesEqual(t, want, have)
	assert.SlicesEqual(t, want, input)
}

func TestSortUniq(t *testing.T) {
	input := []int{3, 1, 4, 1, 5, 9, 2, 6, 5}
	cmpFn := cmp.Compare[int]
	want := []int{1, 2, 3, 4, 5, 6, 9}
	have := pbslices.SortUniq(input, cmpFn)
	assert.SlicesEqual(t, want, have)
}

func TestTransform(t *testing.T) {
	input := []int{1, 2, 3}
	fn := func(x int) string { return fmt.Sprint(x * 2) }
	want := []string{"2", "4", "6"}
	have := pbslices.Transform(input, fn)
	assert.SlicesEqual(t, want, have)
}

func TestUniq(t *testing.T) {
	input := []int{1, 1, 2, 2, 3, 3}
	equal := func(a, b int) bool { return a == b }
	want := []int{1, 2, 3}
	have := pbslices.Uniq(input, equal)
	assert.SlicesEqual(t, want, have)
}
