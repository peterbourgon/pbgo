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

func TestConcatStopPropagates(t *testing.T) {
	var bVisited bool
	a := pbseq.FromSlice([]int{1, 2, 3})
	b := func(yield func(int) bool) {
		bVisited = true
		yield(99)
	}
	have := slices.Collect(pbseq.Limit(pbseq.Concat(a, b), 1))
	assert.SlicesEqual(t, []int{1}, have)
	assert.False(t, bVisited)
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

func FuzzMerge(f *testing.F) {
	f.Add([]byte{1, 2, 3}, []byte{4, 5, 6}, []byte{0, 7})
	f.Add([]byte{}, []byte{1}, []byte{2, 3})
	f.Fuzz(func(t *testing.T, a, b, c []byte) {
		toSorted := func(s []byte) []int {
			out := make([]int, len(s))
			for i, v := range s {
				out[i] = int(v)
			}
			slices.Sort(out)
			return out
		}
		ai, bi, ci := toSorted(a), toSorted(b), toSorted(c)

		got := slices.Collect(pbseq.Merge(cmp.Less[int],
			pbseq.FromSlice(ai),
			pbseq.FromSlice(bi),
			pbseq.FromSlice(ci),
		))

		assert.True(t, slices.IsSorted(got))
		assert.Equal(t, len(ai)+len(bi)+len(ci), len(got))
	})
}
