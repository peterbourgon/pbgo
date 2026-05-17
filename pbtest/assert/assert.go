// Package assert provides minimal assertion helpers for tests.
package assert

import (
	"cmp"
	"errors"
	"math"
	"reflect"
	"slices"
	"testing"
)

// NoErr fails the test if err is non-nil.
func NoErr(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatalf("NoErr: %v", err)
	}
}

// Err fails the test if err is nil.
func Err(tb testing.TB, err error) {
	tb.Helper()
	if err == nil {
		tb.Fatalf("Err: want error, have nil")
	}
}

// ErrIs fails the test if [errors.Is](have, want) is false.
func ErrIs(tb testing.TB, want, have error) {
	tb.Helper()
	if !errors.Is(have, want) {
		tb.Fatalf("ErrIs: want error %v, have %v", want, have)
	}
}

// Equal fails the test if want != have.
func Equal[T comparable](tb testing.TB, want, have T) {
	tb.Helper()
	if want != have {
		tb.Fatalf("Equal: want %v, have %v", want, have)
	}
}

// DeepEqual fails the test if want and have are not [reflect.DeepEqual].
func DeepEqual(tb testing.TB, want, have any) {
	tb.Helper()
	if !reflect.DeepEqual(want, have) {
		tb.Fatalf("DeepEqual: want %v, have %v", want, have)
	}
}

// True fails the test if b is false.
func True(tb testing.TB, b bool) {
	tb.Helper()
	if !b {
		tb.Fatalf("True: want true, have false")
	}
}

// False fails the test if b is true.
func False(tb testing.TB, b bool) {
	tb.Helper()
	if b {
		tb.Fatalf("False: want false, have true")
	}
}

// SlicesEqual fails the test if want and have are not [slices.Equal].
func SlicesEqual[T comparable](tb testing.TB, want, have []T) {
	tb.Helper()
	if !slices.Equal(want, have) {
		tb.Fatalf("SlicesEqual: want %v, have %v", want, have)
	}
}

// Between fails the test if have is outside the inclusive range [lo, hi].
func Between[T cmp.Ordered](tb testing.TB, lo, hi T, have T) {
	tb.Helper()
	if have < lo || hi < have {
		tb.Fatalf("Between: want between %v and %v, have %v", lo, hi, have)
	}
}

// Within fails the test if have differs from want by more than tol as a
// fraction of want (e.g. tol=0.05 allows ±5%). When want is zero, Within
// requires have to be exactly zero.
func Within[T ~int | ~uint64 | ~int64 | ~float64](tb testing.TB, want T, tol float64, have T) {
	tb.Helper()
	if want == 0 {
		if have != 0 {
			tb.Fatalf("Within: want 0 +/- %.1f%%, have %v", tol*100, have)
		}
		return
	}
	off := float64(have-want) / float64(want)
	if math.Abs(off) > math.Abs(tol) {
		tb.Fatalf("Within: want %v +/- %.1f%%, have %v (%.1f%%)", want, tol*100, have, off*100)
	}
}
