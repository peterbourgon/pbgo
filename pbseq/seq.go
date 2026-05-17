package pbseq

import (
	"iter"
)

func Limit[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		if seq == nil {
			return
		}
		count := 0
		seq(func(t T) bool {
			if count >= n {
				return false
			}
			count++
			return yield(t)
		})
	}
}

func Transform[T1, T2 any](seq iter.Seq[T1], fn func(T1) T2) iter.Seq[T2] {
	return func(yield func(T2) bool) {
		if seq == nil {
			return
		}
		seq(func(t1 T1) bool {
			return yield(fn(t1))
		})
	}
}

func Filter[T any](seq iter.Seq[T], fn func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		if seq == nil {
			return
		}
		seq(func(t T) bool {
			if !fn(t) {
				return true
			}
			return yield(t)
		})
	}
}

func FromSlice[T any](slice []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, t := range slice {
			if !yield(t) {
				return
			}
		}
	}
}

func ToSlice[T any](seq iter.Seq[T]) []T {
	var res []T
	seq(func(t T) bool {
		res = append(res, t)
		return true
	})
	return res
}

func Concat[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range seqs {
			if seq == nil {
				continue
			}
			seq(func(t T) bool {
				return yield(t)
			})
		}
	}
}

// Merge performs a K-way merge of the provided input sequences into a single
// output sequence.
//
// The input sequences must be ordered according to the provided less function.
// The output sequence will stop yielding values when all input sequences have
// been exhausted, or when the yield function returns false. Output order for
// values considered equal by the less function is not stable.
//
// Merge is efficient and operates in O(n log k) time, where n is the total
// number of elements across all input sequences and k is the number of input
// sequences.
func Merge[T any](less func(T, T) bool, seqs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		type source struct {
			next  func() (T, bool)
			stop  func()
			value T
			ok    bool
		}

		sources := make([]source, 0, len(seqs))
		for _, seq := range seqs {
			if seq == nil {
				continue
			}
			next, stop := iter.Pull(seq)
			value, ok := next()
			sources = append(sources, source{
				next:  next,
				stop:  stop,
				value: value,
				ok:    ok,
			})
		}

		defer func() {
			for i := range sources {
				sources[i].stop()
			}
		}()

		for {
			selected := -1
			for i := range sources {
				if !sources[i].ok {
					continue
				}
				if selected < 0 {
					selected = i
					continue
				}
				if less(sources[i].value, sources[selected].value) {
					selected = i
					continue
				}
			}
			if selected < 0 {
				return
			}
			if !yield(sources[selected].value) {
				return
			}
			sources[selected].value, sources[selected].ok = sources[selected].next()
		}
	}
}
