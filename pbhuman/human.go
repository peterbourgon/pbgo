package pbhuman

import (
	"fmt"
	"strings"
	"time"
)

func CalcRate[T int | int64 | uint64 | float32 | float64](n T, d time.Duration) float64 {
	return float64(n) / d.Seconds()
}

func CalcRateStr[T int | int64 | uint64 | float32 | float64](n T, d time.Duration) string {
	return RateStr(CalcRate(n, d))
}

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

func DurationStr(d time.Duration) string {
	dt := DurationTrunc(d)
	ds := dt.String()
	if dt >= time.Hour && strings.HasSuffix(ds, "0s") {
		ds = strings.TrimSuffix(ds, "0s")
	}
	return ds
}

func FloatStr(f float64) (s string) {
	defer func() {
		s, _ = strings.CutSuffix(s, ".0")
	}()
	switch {
	case f > 1_000_000_000_000:
		return fmt.Sprintf("%.0fT+", f/1_000_000_000_000)
	case f > 10_000_000_000:
		return fmt.Sprintf("%.0fB", f/1_000_000_000) // 34_567_890_123 -> 35B
	case f > 1_000_000_000:
		return fmt.Sprintf("%.1fB", f/1_000_000_000) // 3_456_789_012 -> 3.5B
	case f > 10_000_000:
		return fmt.Sprintf("%.0fM", f/1_000_000) // 34_567_890 -> 34M
	case f > 1_000_000:
		return fmt.Sprintf("%.1fM", f/1_000_000) // 3_456_789 -> 3.5M
	case f > 10_000:
		return fmt.Sprintf("%.0fK", f/1000) // 32756 -> 32K
	case f > 1_000:
		return fmt.Sprintf("%.1fK", f/1000) // 5142 -> 5.1K
	case f >= 1:
		return fmt.Sprintf("%.1f", f) // 812.312 -> 812.3
	case f == 0:
		return "0"
	default:
		return fmt.Sprintf("%.2f", f) // 0.15845 -> 0.16
	}
}

func ByteStrSI[T interface {
	~int | ~uint | ~int64 | ~uint64
}](n T) string {
	var (
		kib = float64(1024)
		mib = float64(1024 * kib)
		fn  = float64(n)
	)
	switch {
	case fn < 1*kib:
		return fmt.Sprintf("%0.1fB", fn)
	case fn < 100*kib:
		return fmt.Sprintf("%.1fKiB", fn/kib)
	case fn < 1*mib:
		return fmt.Sprintf("%.0fKiB", fn/kib)
	case fn < 100*mib:
		return fmt.Sprintf("%.1fMiB", fn/mib)
	default:
		return fmt.Sprintf("%.0fMiB", fn/mib)
	}
}

func ByteStrStd[T interface {
	~int | ~uint | ~int64 | ~uint64
}](n T) string {
	var (
		kb = float64(1000)
		mb = float64(1000 * kb)
		fn = float64(n)
	)
	switch {
	case fn < 1*kb:
		return fmt.Sprintf("%0.1fB", fn)
	case fn < 100*kb:
		return fmt.Sprintf("%.1fKB", fn/kb)
	case fn < 1*mb:
		return fmt.Sprintf("%.0fKB", fn/kb)
	case fn < 100*mb:
		return fmt.Sprintf("%.1fMB", fn/mb)
	default:
		return fmt.Sprintf("%.0fMB", fn/mb)
	}
}
