package pbslices

import (
	"slices"
)

func Filter[T any](s []T, fn func(T) bool) []T {
	var res []T
	for _, v := range s {
		if fn(v) {
			res = append(res, v)
		}
	}
	return res
}

func First[T any](s []T, fn func(T) bool) (T, bool) {
	if x := Filter(s, fn); len(x) > 0 {
		return x[0], true
	}
	var zero T
	return zero, false
}

func FirstDefault[T any](s []T, fn func(T) bool, defaultValue T) T {
	if x := Filter(s, fn); len(x) > 0 {
		return x[0]
	}
	return defaultValue
}

func Sort[T any](s []T, cmp func(a, b T) int) []T {
	slices.SortFunc(s, cmp)
	return s
}

func SortUniq[T any](s []T, cmp func(a, b T) int) []T {
	s = Sort(s, cmp)
	s = Uniq(s, func(a, b T) bool { return cmp(a, b) == 0 })
	return s
}

func Transform[T1 any, T2 any](s []T1, fn func(T1) T2) []T2 {
	res := make([]T2, len(s))
	for i, v := range s {
		res[i] = fn(v)
	}
	return res
}

func Uniq[T any](s []T, equal func(a, b T) bool) []T {
	s = slices.CompactFunc(s, equal)
	return s
}
