// Package pbhuman provides human-readable formatting helpers for numbers,
// rates, durations, and byte sizes.
package pbhuman

import (
	"fmt"
	"strings"
	"time"
)

// CalcRate returns n divided by d expressed in units per second.
func CalcRate[T ~int | ~uint | ~int64 | ~uint64 | ~float32 | ~float64](n T, d time.Duration) float64 {
	return float64(n) / d.Seconds()
}

// CalcRateStr returns the result of [CalcRate] formatted by [RateStr].
func CalcRateStr[T ~int | ~uint | ~int64 | ~uint64 | ~float32 | ~float64](n T, d time.Duration) string {
	return RateStr(CalcRate(n, d))
}

// RateStr formats r as a human-readable rate string, suffixed with "/s".
func RateStr(r float64) string {
	switch {
	case r < 0.001:
		return "<0.001/s"
	case r < 1:
		return fmt.Sprintf("%.3f/s", r)
	case r < 10:
		return fmt.Sprintf("%.2f/s", r)
	case r < 100:
		return fmt.Sprintf("%.1f/s", r)
	case r < 10_000:
		return fmt.Sprintf("%.0f/s", r)
	case r < 100_000:
		return fmt.Sprintf("%.1fK/s", r/1000)
	case r < 1_000_000:
		return fmt.Sprintf("%.0fK/s", r/1000)
	case r < 10_000_000:
		return fmt.Sprintf("%.1fM/s", r/1_000_000)
	default:
		return fmt.Sprintf("%.0fM/s", r/1_000_000)
	}
}

// DurationTrunc rounds d down to a unit appropriate for its magnitude: days
// for d ≥ 10 days, hours for ≥ 1 day, and so on down to microseconds.
func DurationTrunc(d time.Duration) time.Duration {
	switch {
	case d >= 10*24*time.Hour:
		return d.Truncate(24 * time.Hour)
	case d >= 24*time.Hour:
		return d.Truncate(time.Hour)
	case d >= time.Hour:
		return d.Truncate(time.Minute)
	case d >= time.Minute:
		return d.Truncate(time.Second)
	case d >= time.Second:
		return d.Truncate(100 * time.Millisecond)
	case d >= 10*time.Millisecond:
		return d.Truncate(1000 * time.Microsecond)
	case d >= 1*time.Millisecond:
		return d.Truncate(100 * time.Microsecond)
	case d >= 1*time.Microsecond:
		return d.Truncate(1 * time.Microsecond)
	default:
		return d
	}
}

// DurationStr formats d via [DurationTrunc] and trims a trailing "0s" from
// durations of an hour or more.
func DurationStr(d time.Duration) string {
	dt := DurationTrunc(d)
	ds := dt.String()
	if dt >= time.Hour && strings.HasSuffix(ds, "0s") {
		ds = strings.TrimSuffix(ds, "0s")
	}
	return ds
}

// FloatStr formats f compactly with a magnitude suffix (K, M, B, T+) for
// large values and two decimal places for values below 1.
func FloatStr(f float64) (s string) {
	defer func() {
		s, _ = strings.CutSuffix(s, ".0")
	}()
	switch {
	case f >= 1_000_000_000_000:
		return fmt.Sprintf("%.0fT+", f/1_000_000_000_000)
	case f >= 10_000_000_000:
		return fmt.Sprintf("%.0fB", f/1_000_000_000) // 34_567_890_123 -> 35B
	case f >= 1_000_000_000:
		return fmt.Sprintf("%.1fB", f/1_000_000_000) // 3_456_789_012 -> 3.5B
	case f >= 10_000_000:
		return fmt.Sprintf("%.0fM", f/1_000_000) // 34_567_890 -> 34M
	case f >= 1_000_000:
		return fmt.Sprintf("%.1fM", f/1_000_000) // 3_456_789 -> 3.5M
	case f >= 10_000:
		return fmt.Sprintf("%.0fK", f/1000) // 32756 -> 32K
	case f >= 1_000:
		return fmt.Sprintf("%.1fK", f/1000) // 5142 -> 5.1K
	case f >= 1:
		return fmt.Sprintf("%.1f", f) // 812.312 -> 812.3
	case f == 0:
		return "0"
	default:
		return fmt.Sprintf("%.2f", f) // 0.15845 -> 0.16
	}
}

// ByteStrSI formats n as a binary (1024-based) byte count using IEC suffixes
// (B, KiB, MiB).
func ByteStrSI[T ~int | ~uint | ~int64 | ~uint64 | ~float32 | ~float64](n T) string {
	return byteStr(float64(n), 1024, "B", "KiB", "MiB")
}

// ByteStrStd formats n as a decimal (1000-based) byte count using SI-style
// suffixes (B, KB, MB).
func ByteStrStd[T ~int | ~uint | ~int64 | ~uint64 | ~float32 | ~float64](n T) string {
	return byteStr(float64(n), 1000, "B", "KB", "MB")
}

func byteStr(n, base float64, b, k, m string) string {
	var (
		kUnit = base
		mUnit = base * base
	)
	switch {
	case n < 1*kUnit:
		return fmt.Sprintf("%0.1f%s", n, b)
	case n < 100*kUnit:
		return fmt.Sprintf("%.1f%s", n/kUnit, k)
	case n < 1*mUnit:
		return fmt.Sprintf("%.0f%s", n/kUnit, k)
	case n < 100*mUnit:
		return fmt.Sprintf("%.1f%s", n/mUnit, m)
	default:
		return fmt.Sprintf("%.0f%s", n/mUnit, m)
	}
}
