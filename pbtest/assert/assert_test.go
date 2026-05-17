package assert_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/peterbourgon/pbgo/pbtest/assert"
)

// recorder implements testing.TB enough to capture Fatalf calls without
// aborting the outer test.
type recorder struct {
	testing.TB
	failed bool
	msg    string
}

func (r *recorder) Helper() {}

func (r *recorder) Fatalf(format string, args ...any) {
	r.failed = true
	r.msg = fmt.Sprintf(format, args...)
}

func TestNoErr(t *testing.T) {
	r := &recorder{}
	assert.NoErr(r, nil)
	assert.False(t, r.failed)

	r = &recorder{}
	assert.NoErr(r, errors.New("boom"))
	assert.True(t, r.failed)
}

func TestErr(t *testing.T) {
	r := &recorder{}
	assert.Err(r, errors.New("x"))
	assert.False(t, r.failed)

	r = &recorder{}
	assert.Err(r, nil)
	assert.True(t, r.failed)
}

func TestErrIs(t *testing.T) {
	sentinel := errors.New("sentinel")
	wrapped := errors.Join(errors.New("other"), sentinel)

	r := &recorder{}
	assert.ErrIs(r, sentinel, wrapped)
	assert.False(t, r.failed)

	r = &recorder{}
	assert.ErrIs(r, sentinel, errors.New("nope"))
	assert.True(t, r.failed)
}

func TestEqual(t *testing.T) {
	r := &recorder{}
	assert.Equal(r, 1, 1)
	assert.False(t, r.failed)

	r = &recorder{}
	assert.Equal(r, 1, 2)
	assert.True(t, r.failed)
}

func TestDeepEqual(t *testing.T) {
	r := &recorder{}
	assert.DeepEqual(r, []int{1, 2}, []int{1, 2})
	assert.False(t, r.failed)

	r = &recorder{}
	assert.DeepEqual(r, []int{1, 2}, []int{1, 3})
	assert.True(t, r.failed)
}

func TestTrueFalse(t *testing.T) {
	r := &recorder{}
	assert.True(r, true)
	assert.False(r, false)
	assert.False(t, r.failed)

	r = &recorder{}
	assert.True(r, false)
	assert.True(t, r.failed)

	r = &recorder{}
	assert.False(r, true)
	assert.True(t, r.failed)
}

func TestSlicesEqual(t *testing.T) {
	r := &recorder{}
	assert.SlicesEqual(r, []string{"a", "b"}, []string{"a", "b"})
	assert.False(t, r.failed)

	r = &recorder{}
	assert.SlicesEqual(r, []string{"a"}, []string{"a", "b"})
	assert.True(t, r.failed)
}

func TestBetween(t *testing.T) {
	r := &recorder{}
	assert.Between(r, 1, 10, 5)
	assert.Between(r, 1, 10, 1)
	assert.Between(r, 1, 10, 10)
	assert.False(t, r.failed)

	r = &recorder{}
	assert.Between(r, 1, 10, 0)
	assert.True(t, r.failed)

	r = &recorder{}
	assert.Between(r, 1, 10, 11)
	assert.True(t, r.failed)
}

func TestWithin(t *testing.T) {
	r := &recorder{}
	assert.Within(r, 100, 0.1, 105)
	assert.False(t, r.failed)

	r = &recorder{}
	assert.Within(r, 100, 0.1, 120)
	assert.True(t, r.failed)

	r = &recorder{}
	assert.Within(r, 0, 0.1, 0)
	assert.False(t, r.failed)

	r = &recorder{}
	assert.Within(r, 0, 0.1, 5)
	assert.True(t, r.failed)
}

