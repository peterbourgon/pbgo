package assert

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"reflect"
	"slices"
	"testing"
)

func NoErr(tb testing.TB, err error, args ...any) {
	tb.Helper()
	if err != nil {
		tb.Fatalf("%s: %v", getPrefix("NoErr", args...), err)
	}
}

func Err(tb testing.TB, err error, args ...any) {
	tb.Helper()
	if err == nil {
		tb.Fatalf("%s: want error, have nil", getPrefix("Err", args...))
	}
}

func ErrIs(tb testing.TB, want, have error, args ...any) {
	tb.Helper()
	if !errors.Is(have, want) {
		tb.Fatalf("%s: want error %v, have %v", getPrefix("ErrIs", args...), want, have)
	}
}

func Equal[T comparable](tb testing.TB, want, have T, args ...any) {
	tb.Helper()
	if want != have {
		tb.Fatalf("%s: want %v, have %v", getPrefix("Equal", args...), want, have)
	}
}

func DeepEqual(tb testing.TB, want, have any, args ...any) {
	tb.Helper()
	if !reflect.DeepEqual(want, have) {
		tb.Fatalf("%s: want %v, have %v", getPrefix("DeepEqual", args...), want, have)
	}
}

func True(tb testing.TB, b bool, args ...any) {
	tb.Helper()
	if !b {
		tb.Fatalf("%s: want true, have false", getPrefix("True", args...))
	}
}

func False(tb testing.TB, b bool, args ...any) {
	tb.Helper()
	if b {
		tb.Fatalf("%s: want false, have true", getPrefix("False", args...))
	}
}

func SlicesEqual[T comparable](tb testing.TB, want, have []T, args ...any) {
	tb.Helper()
	if !slices.Equal(want, have) {
		tb.Fatalf("%s: want %v, have %v", getPrefix("SlicesEqual", args...), want, have)
	}
}

func Between[T cmp.Ordered](tb testing.TB, min, max T, have T, args ...any) {
	tb.Helper()
	if have < min || max < have {
		tb.Fatalf("%s: want between %v and %v, have %v", getPrefix("Between", args...), min, max, have)
	}
}

func Within[T ~int | ~uint64 | ~int64 | ~float64](tb testing.TB, want T, tol float64, have T, args ...any) {
	tb.Helper()
	off := float64(have-want) / float64(want)
	if math.Abs(off) > math.Abs(tol) {
		tb.Fatalf("%s: want %v +/- %.1f%%, have %v (%.1f%%)", getPrefix("Within", args...), want, tol*100, have, off*100)
	}
}

func getPrefix(def string, args ...any) (prefix string) {
	if len(args) <= 0 {
		return def
	}

	switch x := args[0].(type) {
	case string:
		return fmt.Sprintf(x, args[1:]...)
	default:
		return fmt.Sprint(args...)
	}
}
