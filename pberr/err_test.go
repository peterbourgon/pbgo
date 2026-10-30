package pberr_test

import (
	"errors"
	"testing"

	"github.com/peterbourgon/pbgo/pberr"
	"github.com/peterbourgon/pbgo/pbtest/assert"
)

func TestAnnotatef(t *testing.T) {
	assert.NoErr(t, pberr.Annotatef(nil, "whatever"))

	err := errors.New("boom")
	assert.ErrIs(t, err, pberr.Annotatef(err, "doing thing"))
	assert.Equal(t, "doing thing: boom", pberr.Annotatef(err, "doing thing").Error())

	assert.Equal(t, err, pberr.Annotatef(err, ""))
}

func TestAnnotatefArgs(t *testing.T) {
	err := errors.New("boom")
	have := pberr.Annotatef(err, "doing %s for %d", "thing", 42)
	assert.Equal(t, "doing thing for 42: boom", have.Error())
	assert.ErrIs(t, err, have)
}

func TestSquelch(t *testing.T) {
	errA := errors.New("a")
	errB := errors.New("b")

	assert.NoErr(t, pberr.Squelch(errA, errA))
	assert.NoErr(t, pberr.Squelch(errA, errB, errA))
	assert.Equal(t, errA, pberr.Squelch(errA, errB))
	assert.NoErr(t, pberr.Squelch(nil, errA))
	assert.Equal(t, errA, pberr.Squelch(errA))
}
