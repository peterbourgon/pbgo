package pbreflect_test

import (
	"testing"

	"github.com/peterbourgon/pbgo/pbreflect"
	"github.com/peterbourgon/pbgo/pbtest/assert"
)

func TestDirectlyImplements(t *testing.T) {
	_, ok := pbreflect.DirectlyImplements[Fooer](MyStruct{})
	assert.True(t, ok)

	_, ok = pbreflect.DirectlyImplements[Fooer](MyEmbed{})
	assert.False(t, ok)

	_, ok = pbreflect.DirectlyImplements[Fooer](NoMethod{})
	assert.False(t, ok)

	_, ok = pbreflect.DirectlyImplements[Fooer](MyEmbedPtr{})
	assert.False(t, ok)

	_, ok = pbreflect.DirectlyImplements[Fooer](PtrRecvr{})
	assert.False(t, ok)

	_, ok = pbreflect.DirectlyImplements[Fooer](&PtrRecvr{})
	assert.True(t, ok)
}

type Fooer interface{ Foo() }

type MyStruct struct{}

func (MyStruct) Foo() {}

type MyEmbed struct{ MyStruct }

type NoMethod struct{}

type MyEmbedPtr struct{ *MyStruct }

type PtrRecvr struct{ MyStruct }

func (*PtrRecvr) Foo() {}
