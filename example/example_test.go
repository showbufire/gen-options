package example

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestFoo(t *testing.T) {
	foo := NewFoo(
		OptionFirst(3),
		OptionSecond(&Bar{}),
		OptionTrd([]string{"yo"}),
		OptionFourth(nil),
	)
	ensure.DeepEqual(t, foo.fst, 3)
	ensure.DeepEqual(t, foo.trd, []string{"yo"})
}

func TestBaz(t *testing.T) {
	baz := &Baz{}
	OptionXYZ(42)(baz)
	ensure.DeepEqual(t, baz.xyz, 42)
}
