package pbsync_test

import (
	"errors"
	"maps"
	"slices"
	"testing"

	"github.com/peterbourgon/pbgo/pbsync"
	"github.com/peterbourgon/pbgo/pbtest/assert"
)

func TestAtomicMapSetGetDel(t *testing.T) {
	var am pbsync.AtomicMap[string, int]

	assert.Equal(t, 0, am.Len())
	assert.False(t, am.Has("x"))

	_, ok := am.Get("x")
	assert.False(t, ok)

	prev, existed := am.Set("x", 1)
	assert.False(t, existed)
	assert.Equal(t, 0, prev)
	assert.Equal(t, 1, am.Len())
	assert.True(t, am.Has("x"))

	have, ok := am.Get("x")
	assert.True(t, ok)
	assert.Equal(t, 1, have)

	prev, existed = am.Set("x", 2)
	assert.True(t, existed)
	assert.Equal(t, 1, prev)

	deleted, existed, size := am.Del("x")
	assert.True(t, existed)
	assert.Equal(t, 2, deleted)
	assert.Equal(t, 0, size)
	assert.False(t, am.Has("x"))
	assert.Equal(t, 0, am.Len())

	deleted, existed, size = am.Del("x")
	assert.False(t, existed)
	assert.Equal(t, 0, deleted)
	assert.Equal(t, 0, size)
}

func TestAtomicMapGetOrCreate(t *testing.T) {
	var am pbsync.AtomicMap[string, int]

	v, existed := am.GetOrCreate("x", func() int { return 42 })
	assert.False(t, existed)
	assert.Equal(t, 42, v)

	v, existed = am.GetOrCreate("x", func() int { return 99 })
	assert.True(t, existed)
	assert.Equal(t, 42, v)
}

func TestAtomicMapRange(t *testing.T) {
	var am pbsync.AtomicMap[string, int]
	am.Set("a", 1)
	am.Set("b", 2)
	am.Set("c", 3)

	seen := map[string]int{}
	err := am.Range(func(k string, v int) error {
		seen[k] = v
		return nil
	})
	assert.NoErr(t, err)
	assert.DeepEqual(t, map[string]int{"a": 1, "b": 2, "c": 3}, seen)
}

func TestAtomicMapRangeError(t *testing.T) {
	var am pbsync.AtomicMap[string, int]
	am.Set("a", 1)
	am.Set("b", 2)

	boom := errors.New("boom")
	count := 0
	err := am.Range(func(k string, v int) error {
		count++
		return boom
	})
	assert.ErrIs(t, boom, err)
	assert.Equal(t, 1, count)
}

func TestAtomicMapAll(t *testing.T) {
	var am pbsync.AtomicMap[string, int]
	am.Set("a", 1)
	am.Set("b", 2)

	have := maps.Collect(am.All())
	assert.DeepEqual(t, map[string]int{"a": 1, "b": 2}, have)
}

func TestAtomicMapApply(t *testing.T) {
	var am pbsync.AtomicMap[string, int]

	// Apply on a missing key: f sees zero value and false, declines to store.
	am.Apply("x", func(v int, exists bool) (int, bool) {
		assert.Equal(t, 0, v)
		assert.False(t, exists)
		return 0, false
	})
	assert.False(t, am.Has("x"))

	// Apply on a missing key: f stores a value.
	am.Apply("x", func(v int, exists bool) (int, bool) {
		assert.False(t, exists)
		return 10, true
	})
	have, ok := am.Get("x")
	assert.True(t, ok)
	assert.Equal(t, 10, have)

	// Apply on an existing key: f sees current value and true, updates it.
	am.Apply("x", func(v int, exists bool) (int, bool) {
		assert.True(t, exists)
		return v + 1, true
	})
	have, ok = am.Get("x")
	assert.True(t, ok)
	assert.Equal(t, 11, have)

	// Apply on an existing key: f declines to store, value is unchanged.
	am.Apply("x", func(v int, exists bool) (int, bool) {
		return 999, false
	})
	have, ok = am.Get("x")
	assert.True(t, ok)
	assert.Equal(t, 11, have)
}

func TestAtomicMapMapSetGetDel(t *testing.T) {
	var mm pbsync.AtomicMapMap[string, string, int]

	_, ok := mm.Get("a", "x")
	assert.False(t, ok)

	prev, existed := mm.Set("a", "x", 1)
	assert.False(t, existed)
	assert.Equal(t, 0, prev)

	have, ok := mm.Get("a", "x")
	assert.True(t, ok)
	assert.Equal(t, 1, have)

	prev, existed = mm.Set("a", "x", 2)
	assert.True(t, existed)
	assert.Equal(t, 1, prev)

	assert.True(t, mm.Has("a", "x"))

	prev, existed = mm.Del("a", "x")
	assert.True(t, existed)
	assert.Equal(t, 2, prev)

	assert.False(t, mm.Has("a", "x"))
}

func TestAtomicMapMapLen(t *testing.T) {
	var mm pbsync.AtomicMapMap[string, string, int]

	assert.Equal(t, 0, mm.Len())
	assert.Equal(t, 0, mm.LenOf("a"))

	mm.Set("a", "x", 1)
	mm.Set("a", "y", 2)
	mm.Set("b", "z", 3)

	assert.Equal(t, 2, mm.Len())
	assert.Equal(t, 2, mm.LenOf("a"))
	assert.Equal(t, 1, mm.LenOf("b"))
	assert.Equal(t, 0, mm.LenOf("missing"))

	mm.Del("a", "x")
	assert.Equal(t, 1, mm.LenOf("a"))

	mm.Del("a", "y")
	assert.Equal(t, 0, mm.LenOf("a"))
	assert.Equal(t, 1, mm.Len())
}

func TestAtomicMapMapGetOrCreate(t *testing.T) {
	var mm pbsync.AtomicMapMap[string, string, int]

	v, existed := mm.GetOrCreate("a", "x", func() int { return 42 })
	assert.False(t, existed)
	assert.Equal(t, 42, v)

	v, existed = mm.GetOrCreate("a", "x", func() int { return 99 })
	assert.True(t, existed)
	assert.Equal(t, 42, v)
}

func TestAtomicMapMapAll(t *testing.T) {
	var mm pbsync.AtomicMapMap[string, string, int]

	mm.Set("a", "x", 1)
	mm.Set("a", "y", 2)
	mm.Set("b", "z", 3)

	want := map[string]map[string]int{
		"a": {"x": 1, "y": 2},
		"b": {"z": 3},
	}

	have := map[string]map[string]int{}
	for k1, inner := range mm.All() {
		have[k1] = maps.Collect(inner)
	}

	assert.DeepEqual(t, want, have)
}

func TestAtomicMapMapKeys(t *testing.T) {
	var mm pbsync.AtomicMapMap[string, string, int]

	mm.Set("a", "x", 1)
	mm.Set("b", "y", 2)

	have := slices.Collect(mm.Keys())
	slices.Sort(have)
	assert.SlicesEqual(t, []string{"a", "b"}, have)
}

func TestAtomicMapMapKeysEmpty(t *testing.T) {
	var mm pbsync.AtomicMapMap[string, string, int]

	have := slices.Collect(mm.Keys())
	assert.Equal(t, 0, len(have))
}

func TestAtomicMapMapAllOf(t *testing.T) {
	var mm pbsync.AtomicMapMap[string, string, int]

	mm.Set("a", "x", 1)
	mm.Set("a", "y", 2)
	mm.Set("b", "z", 3)

	have := maps.Collect(mm.AllOf("a"))
	assert.DeepEqual(t, map[string]int{"x": 1, "y": 2}, have)
}

func TestAtomicMapMapAllOfMissing(t *testing.T) {
	var mm pbsync.AtomicMapMap[string, string, int]

	have := maps.Collect(mm.AllOf("missing"))
	assert.Equal(t, 0, len(have))
}

func TestAtomicMapMapAllSkipOuter(t *testing.T) {
	var mm pbsync.AtomicMapMap[string, string, int]

	mm.Set("a", "x", 1)
	mm.Set("b", "y", 2)

	var seen []string
	for k1, inner := range mm.All() {
		if k1 == "a" {
			continue
		}
		for k2 := range inner {
			seen = append(seen, k1+"/"+k2)
		}
	}

	assert.SlicesEqual(t, []string{"b/y"}, seen)
}
