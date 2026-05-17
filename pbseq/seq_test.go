package pbseq_test

import (
	"cmp"
	"fmt"
	"slices"
	"testing"

	"github.com/peterbourgon/pbgo/pbseq"
	"github.com/peterbourgon/pbgo/pbtest/assert"
)

func TestLimit(t *testing.T) {
	seq := pbseq.FromSlice([]int{1, 2, 3, 4, 5})
	limited := pbseq.Limit(seq, 3)

	want := []int{1, 2, 3}
	have := slices.Collect(limited)

	assert.SlicesEqual(t, want, have)
}

func TestTransform(t *testing.T) {
	seq := pbseq.FromSlice([]int{1, 2, 3})
	transformed := pbseq.Transform(seq, func(x int) string {
		return fmt.Sprint(x * 2)
	})

	have := slices.Collect(transformed)
	want := []string{"2", "4", "6"}

	assert.SlicesEqual(t, want, have)
}

func TestFilter(t *testing.T) {
	seq := pbseq.FromSlice([]int{1, 2, 3, 4, 5})
	filtered := pbseq.Filter(seq, func(x int) bool {
		return x%2 == 0
	})

	have := slices.Collect(filtered)
	want := []int{2, 4}

	assert.SlicesEqual(t, want, have)
}

func TestSlice(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}

	seq := pbseq.FromSlice(input)
	have := slices.Collect(seq)

	assert.SlicesEqual(t, input, have)
}

func TestConcat(t *testing.T) {
	seq1 := pbseq.FromSlice([]int{1, 2})
	seq2 := pbseq.FromSlice([]int{3, 4})
	seq3 := pbseq.FromSlice([]int{5})
	concatenated := pbseq.Concat(seq1, seq2, seq3)

	have := slices.Collect(concatenated)
	want := []int{1, 2, 3, 4, 5}

	assert.SlicesEqual(t, want, have)
}

func TestMerge(t *testing.T) {
	seq1 := pbseq.FromSlice([]int{1, 3, 5})
	seq2 := pbseq.FromSlice([]int{2, 4, 6})
	seq3 := pbseq.FromSlice([]int{999})
	merged := pbseq.Merge(cmp.Less, seq1, seq2, seq3)

	have := slices.Collect(merged)
	want := []int{1, 2, 3, 4, 5, 6, 999}

	assert.SlicesEqual(t, want, have)
}
