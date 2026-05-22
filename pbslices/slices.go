// Package pbslices provides utilities for working with slices.
package pbslices

import (
	"slices"
)

// Filter returns a new slice containing the elements of s for which fn returns
// true, in their original order.
func Filter[T any](s []T, fn func(T) bool) []T {
	var res []T
	for _, v := range s {
		if fn(v) {
			res = append(res, v)
		}
	}
	return res
}

// First returns the first element of s for which fn returns true, along with
// true. If no element matches, it returns the zero value and false.
func First[T any](s []T, fn func(T) bool) (T, bool) {
	for _, v := range s {
		if fn(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// FirstDefault returns the first element of s for which fn returns true, or
// defaultValue if no element matches.
func FirstDefault[T any](s []T, fn func(T) bool, defaultValue T) T {
	if v, ok := First(s, fn); ok {
		return v
	}
	return defaultValue
}

// Partition s into two slices, the first containing the elements for which fn
// returns true, and the second containing the elements for which fn returns
// false. The order of elements in the output slices is the same as in the input
// slice.
func Partition[T any](s []T, fn func(T) bool) (a, b []T) {
	for _, v := range s {
		if fn(v) {
			a = append(a, v)
		} else {
			b = append(b, v)
		}
	}
	return a, b
}

// Sort sorts s in place using cmp and returns it.
func Sort[T any](s []T, cmp func(a, b T) int) []T {
	slices.SortFunc(s, cmp)
	return s
}

// SortUniq sorts s in place using cmp, removes consecutive duplicates as
// determined by cmp, and returns the resulting slice.
func SortUniq[T any](s []T, cmp func(a, b T) int) []T {
	s = Sort(s, cmp)
	s = Uniq(s, func(a, b T) bool { return cmp(a, b) == 0 })
	return s
}

// Transform returns a new slice produced by applying fn to each element of s.
func Transform[T1 any, T2 any](s []T1, fn func(T1) T2) []T2 {
	res := make([]T2, len(s))
	for i, v := range s {
		res[i] = fn(v)
	}
	return res
}

// Uniq removes consecutive runs of equal elements from s, in place, and
// returns the resulting slice. To dedupe an unsorted slice, sort it first or
// use SortUniq.
func Uniq[T any](s []T, equal func(a, b T) bool) []T {
	s = slices.CompactFunc(s, equal)
	return s
}
