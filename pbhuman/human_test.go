package pbhuman_test

import (
	"testing"
	"time"

	"github.com/peterbourgon/pbgo/pbhuman"
	"github.com/peterbourgon/pbgo/pbtest/assert"
)

func TestCalcRate(t *testing.T) {
	assert.Equal(t, 10.0, pbhuman.CalcRate(10, time.Second))
	assert.Equal(t, 2.0, pbhuman.CalcRate(1, 500*time.Millisecond))
	assert.Equal(t, 0.5, pbhuman.CalcRate(int64(30), time.Minute))
}

func TestCalcRateStr(t *testing.T) {
	assert.Equal(t, "10.0/s", pbhuman.CalcRateStr(10, time.Second))
}

func TestRateStr(t *testing.T) {
	for _, tc := range []struct {
		in   float64
		want string
	}{
		{0.0001, "<0.001/s"},
		{0.5, "0.500/s"},
		{5, "5.00/s"},
		{50, "50.0/s"},
		{500, "500/s"},
		{50_000, "50.0K/s"},
		{500_000, "500K/s"},
		{5_000_000, "5.0M/s"},
		{50_000_000, "50M/s"},
	} {
		t.Run(tc.want, func(t *testing.T) {
			assert.Equal(t, tc.want, pbhuman.RateStr(tc.in))
		})
	}
}

func TestDurationTrunc(t *testing.T) {
	for _, tc := range []struct {
		in   time.Duration
		want time.Duration
	}{
		{15 * 24 * time.Hour, 15 * 24 * time.Hour},
		{36 * time.Hour, 36 * time.Hour},
		{90*time.Minute + 30*time.Second, 90 * time.Minute},
		{90*time.Second + 250*time.Millisecond, 90 * time.Second},
		{2*time.Second + 123*time.Millisecond, 2*time.Second + 100*time.Millisecond},
		{500 * time.Nanosecond, 500 * time.Nanosecond},
	} {
		assert.Equal(t, tc.want, pbhuman.DurationTrunc(tc.in))
	}
}

func TestDurationStr(t *testing.T) {
	assert.Equal(t, "2h30m", pbhuman.DurationStr(2*time.Hour+30*time.Minute))
	assert.Equal(t, "1h1m", pbhuman.DurationStr(time.Hour+time.Minute+5*time.Second))
	assert.Equal(t, "1m30s", pbhuman.DurationStr(90*time.Second))
}

func TestFloatStr(t *testing.T) {
	for _, tc := range []struct {
		in   float64
		want string
	}{
		{0, "0"},
		{0.15845, "0.16"},
		{812.312, "812.3"},
		{5142, "5.1K"},
		{32100, "32K"},
		{3_456_789, "3.5M"},
		{34_100_000, "34M"},
		{3_456_789_012, "3.5B"},
		{34_567_890_123, "35B"},
		{5_000_000_000_000, "5T+"},
	} {
		t.Run(tc.want, func(t *testing.T) {
			assert.Equal(t, tc.want, pbhuman.FloatStr(tc.in))
		})
	}
}

func TestByteStrSI(t *testing.T) {
	assert.Equal(t, "512.0B", pbhuman.ByteStrSI(512))
	assert.Equal(t, "1.0KiB", pbhuman.ByteStrSI(1024))
	assert.Equal(t, "500KiB", pbhuman.ByteStrSI(500*1024))
	assert.Equal(t, "1.0MiB", pbhuman.ByteStrSI(1024*1024))
	assert.Equal(t, "500MiB", pbhuman.ByteStrSI(500*1024*1024))
}

func TestByteStrStd(t *testing.T) {
	assert.Equal(t, "512.0B", pbhuman.ByteStrStd(512))
	assert.Equal(t, "1.0KB", pbhuman.ByteStrStd(1000))
	assert.Equal(t, "500KB", pbhuman.ByteStrStd(500_000))
	assert.Equal(t, "1.0MB", pbhuman.ByteStrStd(1_000_000))
	assert.Equal(t, "500MB", pbhuman.ByteStrStd(500_000_000))
}

func FuzzRateStr(f *testing.F) {
	for _, seed := range []float64{0, 1e-9, 0.5, 1, 99, 1234, 1e9, 1e18} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, r float64) {
		_ = pbhuman.RateStr(r)
	})
}

func FuzzFloatStr(f *testing.F) {
	for _, seed := range []float64{0, 1e-9, 0.5, 1, 1000, 1e6, 1e12, 1e18} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, x float64) {
		_ = pbhuman.FloatStr(x)
	})
}

func FuzzByteStr(f *testing.F) {
	for _, seed := range []uint64{0, 1, 1023, 1024, 1_000_000, 1 << 40} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, n uint64) {
		_ = pbhuman.ByteStrSI(n)
		_ = pbhuman.ByteStrStd(n)
	})
}
