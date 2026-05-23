package pbreflect_test

import (
	"fmt"
	"testing"

	"github.com/peterbourgon/pbgo/pbreflect"
	"github.com/peterbourgon/pbgo/pbtest/assert"
)

func TestDirectlyImplements(t *testing.T) {
	for _, tc := range []struct {
		x    any
		want bool
	}{
		{PtrRecvr{}, false},
		{&PtrRecvr{}, true},
		{ValRecvr{}, true},
		{&ValRecvr{}, true},
		{EmbedPtrRecvr{}, false},
		{&EmbedPtrRecvr{}, false},
		{EmbedPtrPtrRecvr{}, false},
		{&EmbedPtrPtrRecvr{}, false},
		{NoMethod{}, false},
		{&NoMethod{}, false},
	} {
		t.Run(fmt.Sprintf("%T", tc.x), func(t *testing.T) {
			_, ok := pbreflect.DirectlyImplements[Fooer](tc.x)
			assert.Equal(t, tc.want, ok)
		})
	}
}

type Fooer interface{ Foo() }

type PtrRecvr struct{}

func (*PtrRecvr) Foo() {}

type ValRecvr struct{}

func (ValRecvr) Foo() {}

type EmbedPtrRecvr struct{ PtrRecvr }

type EmbedPtrPtrRecvr struct{ *PtrRecvr }

type NoMethod struct{}
