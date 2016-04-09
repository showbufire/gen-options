package example

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestFoo(t *testing.T) {
	foo := NewFoo(
		MyOptionFirst(3),
		MyOptionSecond(&Bar{}),
		MyOptionTrd([]string{"yo"}),
		MyOptionFourth(nil),
	)
	ensure.DeepEqual(t, foo.fst, 3)
	ensure.DeepEqual(t, foo.trd, []string{"yo"})
}
